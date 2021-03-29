package gogm

import (
	"errors"
	"fmt"
	"github.com/adam-hanna/arrayOperations"
	"github.com/cornelk/hashmap"
	dsl "github.com/mindstand/go-cypherdsl"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func resultToStringArrV3(res neo4j.Result) ([]string, error) {
	if res == nil {
		return nil, errors.New("result is nil")
	}

	var result []string

	for res.Next() {
		val := res.Record().Values
		// nothing to parse
		if val == nil || len(val) == 0 {
			continue
		}

		str, ok := val[0].(string)
		if !ok {
			return nil, fmt.Errorf("unable to parse [%T] to string. Value is %v: %w", val[0], val[0], ErrInternal)
		}

		result = append(result, str)
	}

	return result, nil
}

//drops all known indexes
func dropAllIndexesAndConstraintsV3() error {
	sess, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer sess.Close()

	res, err := sess.Run("CALL db.constraints", nil)
	if err != nil {
		return err
	}

	constraints, err := resultToStringArrV3(res)
	if err != nil {
		return err
	}

	//if there is anything, get rid of it
	if len(constraints) != 0 {
		tx, err := sess.BeginTransaction()
		if err != nil {
			return err
		}

		for _, constraint := range constraints {
			log.Debugf("dropping constraint '%s'", constraint)
			_, err := tx.Run(fmt.Sprintf("DROP %s", constraint), nil)
			if err != nil {
				oerr := err
				err = tx.Rollback()
				if err != nil {
					return fmt.Errorf("failed to rollback, original error was %s", oerr.Error())
				}

				return oerr
			}
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	res, err = sess.Run("CALL db.indexes()", nil)
	if err != nil {
		return err
	}

	indexes, err := resultToStringArrV3(res)
	if err != nil {
		return err
	}

	//if there is anything, get rid of it
	if len(indexes) != 0 {
		tx, err := sess.BeginTransaction()
		if err != nil {
			return err
		}

		for _, index := range indexes {
			if len(index) == 0 {
				return errors.New("invalid index config")
			}

			_, err := tx.Run(fmt.Sprintf("DROP %s", index), nil)
			if err != nil {
				oerr := err
				err = tx.Rollback()
				if err != nil {
					return fmt.Errorf("failed to rollback, original error was %s", oerr.Error())
				}

				return oerr
			}
		}

		return tx.Commit()
	} else {
		return nil
	}
}

//creates all indexes
func createAllIndexesAndConstraintsV3(mappedTypes *hashmap.HashMap) error {
	sess, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer sess.Close()

	//validate that we have to do anything
	if mappedTypes == nil || mappedTypes.Len() == 0 {
		return errors.New("must have types to map")
	}

	numIndexCreated := 0

	tx, err := sess.BeginTransaction()
	if err != nil {
		return err
	}

	//index and/or create unique constraints wherever necessary
	//for node, structConfig := range mappedTypes{
	for nodes := range mappedTypes.Iter() {
		node := nodes.Key.(string)
		structConfig := nodes.Value.(structDecoratorConfig)
		if structConfig.Fields == nil || len(structConfig.Fields) == 0 {
			continue
		}

		var indexFields []string

		for _, config := range structConfig.Fields {
			//pk is a special unique key
			if config.PrimaryKey || config.Unique {
				numIndexCreated++

				cyp, err := dsl.QB().Create(dsl.NewConstraint(&dsl.ConstraintConfig{
					Unique: true,
					Name:   node,
					Type:   structConfig.Label,
					Field:  config.Name,
				})).ToCypher()
				if err != nil {
					return err
				}

				_, err = tx.Run(cyp, nil)
				if err != nil {
					oerr := err
					err = tx.Rollback()
					if err != nil {
						return fmt.Errorf("failed to rollback, original error was %s", oerr.Error())
					}

					return oerr
				}
			} else if config.Index {
				indexFields = append(indexFields, config.Name)
			}
		}

		//create composite index
		if len(indexFields) > 0 {
			numIndexCreated++
			cyp, err := dsl.QB().Create(dsl.NewIndex(&dsl.IndexConfig{
				Type:   structConfig.Label,
				Fields: indexFields,
			})).ToCypher()
			if err != nil {
				return err
			}

			_, err = tx.Run(cyp, nil)
			if err != nil {
				oerr := err
				err = tx.Rollback()
				if err != nil {
					return fmt.Errorf("failed to rollback, original error was %s", oerr.Error())
				}

				return oerr
			}
		}
	}

	log.Debugf("created (%v) indexes", numIndexCreated)

	return tx.Commit()
}

//verifies all indexes
func verifyAllIndexesAndConstraintsV3(mappedTypes *hashmap.HashMap) error {
	sess, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return err
	}
	defer sess.Close()

	//validate that we have to do anything
	if mappedTypes == nil || mappedTypes.Len() == 0 {
		return errors.New("must have types to map")
	}

	var constraints []string
	var indexes []string

	//build constraint strings
	for nodes := range mappedTypes.Iter() {
		node := nodes.Key.(string)
		structConfig := nodes.Value.(structDecoratorConfig)

		if structConfig.Fields == nil || len(structConfig.Fields) == 0 {
			continue
		}

		fields := []string{}

		for _, config := range structConfig.Fields {

			if config.PrimaryKey || config.Unique {
				t := fmt.Sprintf("CONSTRAINT ON (%s:%s) ASSERT %s.%s IS UNIQUE", node, structConfig.Label, node, config.Name)
				constraints = append(constraints, t)

				indexes = append(indexes, fmt.Sprintf("INDEX ON :%s(%s)", structConfig.Label, config.Name))

			} else if config.Index {
				fields = append(fields, config.Name)
			}
		}

		f := "("
		for _, field := range fields {
			f += field
		}

		f += ")"

		indexes = append(indexes, fmt.Sprintf("INDEX ON :%s%s", structConfig.Label, f))

	}

	//get whats there now
	foundResult, err := sess.Run("CALL db.constraints", nil)
	if err != nil {
		return err
	}

	foundConstraints, err := resultToStringArrV3(foundResult)
	if err != nil {
		return err
	}

	foundInxdexResult, err := sess.Run("CALL db.indexes()", nil)
	if err != nil {
		return err
	}

	foundIndexes, err := resultToStringArrV3(foundInxdexResult)
	if err != nil {
		return err
	}

	//verify from there
	delta, found := arrayOperations.Difference(foundIndexes, indexes)
	if !found {
		return fmt.Errorf("found differences in remote vs ogm for found indexes, %v", delta)
	}

	log.Debug(delta)

	var founds []string

	for _, constraint := range foundConstraints {
		founds = append(founds, constraint)
	}

	delta, found = arrayOperations.Difference(founds, constraints)
	if !found {
		return fmt.Errorf("found differences in remote vs ogm for found constraints, %v", delta)
	}

	log.Debug(delta)

	return nil
}
