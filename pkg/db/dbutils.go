package db

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
	--------------
	| dbutils.go |
	--------------
	Utilities to connect to mongodb and keep a global variable to be accessed from anywhere in the file after calling the init and connect method once in the main function.
*/

// Global client
var mdb *MongoDB

type MongoDB struct {
	Url            string
	hasInitialised bool
	Client         *mongo.Client
}

// Function to initialize the MongoDB object for global access
func Init(url string) {
	mdb = &MongoDB{
		Url:            url,
		hasInitialised: true,
		Client:         nil,
	}
}

// Function to return the global MongoDB object
func GetMDB() *MongoDB {
	mdb.panicIfNotInitialised()
	return mdb
}

// Method to panic if the object has not been initialised
func (m *MongoDB) panicIfNotInitialised() {
	if !(m.hasInitialised) {
		panic(errors.New("MongoDB object not initialised"))
	}
}

// Method to connect to the database
func (m *MongoDB) Connect() {

	opt := options.Client()
	opt.ApplyURI(m.Url)

	c, err := mongo.Connect(context.Background(), opt)
	if err != nil {
		log.Fatalf("Error while connecting to mongodb: %v", err)
	}

	// Gets timeout from env in seconds
	timeout, err := strconv.Atoi(os.Getenv("MONGODB_CONNECTION_TIMEOUT"))
	if err != nil {
		timeout = 10
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	if err := c.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to connect to mongodb: %v", err)
		return
	}

	// Log if everything went successful (change into only on development mode)
	log.Printf("Successfully connected to mongodb: %s\n", m.Url)
	m.Client = c
}

// Method to get blog database
func (m *MongoDB) BlogDatabase() *mongo.Database {
	return m.Client.Database("go-blog")
}

// Method to get blogs collection
func (m *MongoDB) BlogsCollection() *mongo.Collection {
	return m.BlogDatabase().Collection("blogs")
}

// Method to get images collection
func (m *MongoDB) ImagesCollection() *mongo.Collection {
	return m.BlogDatabase().Collection("images")
}
