package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"sync"
)

type MongoDB struct {
	*mongo.Database
}

var mongodb mongo.Database

var mongodbMutex sync.Mutex

func GetMongoDB(name string) *MongoDB {
	uri := os.Getenv("mongodb_url")
	if uri == "" {
		uri = "mongodb://a_admin_c_rw:yameMongotest422@10.1.1.232:40005"
	}
	opts := options.Client().ApplyURI(uri)

	// 连接数据库
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Println(err)
	}

	// 判断服务是不是可用
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		log.Println(err)
	}

	mongodb = *client.Database(name)
	return &MongoDB{Database: &mongodb}
}

func GetDefaultMongoDB() *MongoDB {
	mongodbMutex.Lock()
	defer mongodbMutex.Unlock()

	if mongodb.Client() == nil {
		return GetMongoDB("cmdb")
	}

	return &MongoDB{Database: &mongodb}
}

func (d *MongoDB) Drop(collections ...string) {
	for _, collection := range collections {
		c := d.Database.Collection(collection)
		// 清空文档
		err := c.Drop(context.Background())
		if err != nil {
			log.Println(err)
		}
	}
}

func (d *MongoDB) InsertOne(collection string, document interface{}) {
	c := d.Database.Collection(collection)
	_, err := c.InsertOne(context.Background(), document)
	if err != nil {
		log.Println(err)
	}
}

func (d *MongoDB) FindOne(collection string, filter interface{}, document interface{}) {
	c := d.Database.Collection(collection)
	err := c.FindOne(context.Background(), filter).Decode(document)
	if err != nil {
		log.Println(err)
	}
}
