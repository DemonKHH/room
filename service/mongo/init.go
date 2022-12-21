package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	var err error
	// 设置客户端连接配置（connect=direct当数据库为集群时，可通过配置此项解决此驱动连接docker方式部署的mongodb集群由于docker方式部署的访问域名连接，而外面的网络访问不到的情况） mongodb://user:password@ip:port/database?connect=direct
	// &replicaSet=集群名称  authSource=数据库名称  此配置也可以指定数据库
	clientOptions := options.Client().ApplyURI("mongodb+srv://demon:ht19980706@cluster0.u2503.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

	// 连接到MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func GetMongoClient() *mongo.Client {
	return client
}
