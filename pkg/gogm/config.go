// Copyright (c) 2020 MindStand Technologies, Inc
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package gogm

import (
	"errors"
	"fmt"
	"github.com/cornelk/hashmap"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"strings"
)

var neoVersion float64

var log logrus.FieldLogger

func init() {
	makeDefaultLogger()
}

func makeDefaultLogger() {
	_log := logrus.New()
	log = _log.WithField("error", "not initialized")
}

func getIsV4() bool {
	return true
}

// Config Defined GoGM config
type Config struct {
	// Host is the neo4j host
	Host string `yaml:"host" json:"host" mapstructure:"host"`
	// Port is the neo4j port
	Port int `yaml:"port" json:"port" mapstructure:"port"`

	// IsCluster specifies whether GoGM is connecting to a casual cluster or not
	IsCluster bool `yaml:"is_cluster" json:"is_cluster" mapstructure:"is_cluster"`

	// Username is the GoGM username
	Username string `yaml:"username" json:"username" mapstructure:"username"`
	// Password is the GoGM password
	Password string `yaml:"password" json:"password" mapstructure:"password"`

	// PoolSize is the size of the connection pool for GoGM
	PoolSize int `yaml:"pool_size" json:"pool_size" mapstructure:"pool_size"`

	Realm string `yaml:"realm" json:"realm" mapstructure:"realm"`

	Encrypted bool `yaml:"encrypted" json:"encrypted" mapstructure:"encrypted"`

	// Index Strategy defines the index strategy for GoGM
	IndexStrategy IndexStrategy `yaml:"index_strategy" json:"index_strategy" mapstructure:"index_strategy"`
	TargetDbs     []string      `yaml:"target_dbs" json:"target_dbs" mapstructure:"target_dbs"`

	Logger logrus.FieldLogger `yaml:"-" json:"-" mapstructure:"-"`
	// if logger is not nil log level will be ignored
	LogLevel string `json:"log_level" yaml:"log_level" mapstructure:"log_level"`
}

// ConnectionString builds the neo4j bolt/bolt+s connection string
func (c *Config) ConnectionString() string {
	var protocol string

	if c.IsCluster {
		protocol = "bolt+s"
	} else {
		protocol = "bolt"
	}
	// In case of special characters in password string
	//password := url.QueryEscape(c.Password)
	return fmt.Sprintf("%s://%s:%v", protocol, c.Host, c.Port)
}

// Index Strategy typedefs int to define different index approaches
type IndexStrategy int

const (
	// ASSERT_INDEX ensures that all indices are set and sets them if they are not there
	ASSERT_INDEX IndexStrategy = 0
	// VALIDATE_INDEX ensures that all indices are set
	VALIDATE_INDEX IndexStrategy = 1
	// IGNORE_INDEX skips the index step of setup
	IGNORE_INDEX IndexStrategy = 2
)

//holds mapped types
var mappedTypes = &hashmap.HashMap{}

//thread pool
var driver neo4j.Driver

//relationship + label
var mappedRelations = &relationConfigs{}

func makeRelMapKey(start, edge, direction, rel string) string {
	return fmt.Sprintf("%s-%s-%v-%s", start, edge, direction, rel)
}

var isSetup = false

// Init sets up gogm. Takes in config object and variadic slice of gogm nodes to map.
// Note: Must pass pointers to nodes!
func Init(conf *Config, mapTypes ...interface{}) error {
	return setupInit(false, conf, mapTypes...)
}

// Resets GoGM configuration
func Reset() {
	mappedTypes = &hashmap.HashMap{}
	mappedRelations = &relationConfigs{}
	isSetup = false
}

var internalConfig *Config

// internal setup logic for gogm
func setupInit(isTest bool, conf *Config, mapTypes ...interface{}) error {
	if isSetup && !isTest {
		return errors.New("gogm has already been initialized")
	} else if isTest && isSetup {
		mappedRelations = &relationConfigs{}
	}

	if !isTest {
		if conf == nil {
			return errors.New("config can not be nil")
		}
	}

	if conf != nil {
		if conf.Logger != nil {
			log = conf.Logger
		} else {
			_log := logrus.New()

			// set info if nothing has been set
			if conf.LogLevel == "" {
				conf.LogLevel = "INFO"
			}
			lvl, err := logrus.ParseLevel(conf.LogLevel)
			if err != nil {
				return err
			}
			_log.SetLevel(lvl)
			log = _log
		}

		if conf.TargetDbs == nil || len(conf.TargetDbs) == 0 {
			conf.TargetDbs = []string{"neo4j"}
		}

		internalConfig = conf
	} else {
		internalConfig = &Config{
			TargetDbs: []string{"neo4j"},
		}
	}

	log.Debug("mapping types")
	for _, t := range mapTypes {
		name := reflect.TypeOf(t).Elem().Name()
		dc, err := getStructDecoratorConfig(t, mappedRelations)
		if err != nil {
			return err
		}

		log.Debugf("mapped type %s", name)
		mappedTypes.Set(name, *dc)
	}

	log.Debug("validating edges")
	if err := mappedRelations.Validate(); err != nil {
		log.WithError(err).Error("failed to validate edges")
		return err
	}

	if !isTest {
		log.Debug("opening connection to neo4j")
		// todo tls support
		config := func(neoConf *neo4j.Config) {
			//neoConf.Encrypted = conf.Encrypted
			neoConf.MaxConnectionPoolSize = conf.PoolSize
		}
		var err error
		driver, err = neo4j.NewDriver(conf.ConnectionString(), neo4j.BasicAuth(conf.Username, conf.Password, conf.Realm), config)
		if err != nil {
			return err
		}

		// get neoversion
		sess := driver.NewSession(neo4j.SessionConfig{
			AccessMode: neo4j.AccessModeRead,
		})

		defer func() {
			if e := sess.Close(); e != nil {
				fmt.Println("close the session err：", e)
			}
		}()

		res, err := sess.Run("return 1", nil)
		if err != nil {
			return err
		} else if err = res.Err(); err != nil {
			return err
		}

		sum, err := res.Consume()
		if err != nil {
			return err
		}

		// grab version
		version := strings.Split(strings.Replace(strings.ToLower(sum.Server().Version()), "neo4j/", "", -1), ".")
		neoVersion, err = strconv.ParseFloat(version[0], 64)
		if err != nil {
			return err
		}
	}

	log.Debug("starting index verification step")
	if !isTest {
		var err error
		if conf.IndexStrategy == ASSERT_INDEX {
			log.Debug("chose ASSERT_INDEX strategy")
			log.Debug("dropping all known indexes")
			err = dropAllIndexesAndConstraints()
			if err != nil {
				return err
			}

			log.Debug("creating all mapped indexes")
			err = createAllIndexesAndConstraints(mappedTypes)
			if err != nil {
				return err
			}

			log.Debug("verifying all indexes")
			err = verifyAllIndexesAndConstraints(mappedTypes)
			if err != nil {
				return err
			}
		} else if conf.IndexStrategy == VALIDATE_INDEX {
			log.Debug("chose VALIDATE_INDEX strategy")
			log.Debug("verifying all indexes")
			err = verifyAllIndexesAndConstraints(mappedTypes)
			if err != nil {
				return err
			}
		} else if conf.IndexStrategy == IGNORE_INDEX {
			log.Debug("ignoring indexes")
		} else {
			return errors.New("unknown index strategy")
		}
	}

	log.Debug("setup complete")

	isSetup = true

	return nil
}
