package collections

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"ohmytech.io/platform/config"
)

var (
	conn   *mongo.Client
	err    error
	dbName string
)

// getConn :
func getConn() *mongo.Client {
	if nil == conn {
		conf := config.GetConfig().NoSQL

		dbName = conf.Schema

		// Set client options
		clientOptions := options.Client().ApplyURI(conf.Connector)

		// Connect to MongoDB
		conn, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			panic(err.Error())
		}

		// Check the connection
		err = conn.Ping(context.TODO(), nil)
		if err != nil {
			panic(err.Error())
		}
	}

	return conn
}
