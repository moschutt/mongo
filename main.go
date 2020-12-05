package main

import (
	"context"
	"os"
	"strings"

	//   "fmt"
	"log"
	"time"

	errors "github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//   "go.mongodb.org/mongo-driver/mongo/readpref"
)

// Post is a struct
type Post struct {
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`
}

// InsertPost is a function
func (a myApp) InsertPost(ctx context.Context, title, body string) (*mongo.InsertOneResult, error) {
	post := Post{Title: title, Body: body}

	collection := a.client.Database("test").Collection("test")
	insertResult, err := collection.InsertOne(ctx, post)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert")
	}

	return insertResult, nil
}

// GetPost is a function
func (a myApp) GetPost(ctx context.Context, filter bson.D) (Post, error) {
	collection := a.client.Database("test").Collection("test")

	var post Post

	err := collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		return Post{}, errors.Wrap(err, "failed get of data")
	}

	return post, nil
}

type myApp struct {
	host       string
	port       string
	uri        string
	collection string
	db         string
	client     *mongo.Client
}

func (a myApp) makeMongoUri() myApp {
	if a.port != "" {
		a.uri = "mongodb://" + a.host + ":" + a.port
	} else {
		a.uri = "mongodb://" + a.host
	}

	return a
}
func appInit(host, port string) (myApp, error) {
	a := myApp{
		host: strings.TrimSpace(host),
		port: strings.TrimSpace(port),
	}
	a = a.makeMongoUri()

	client, err := mongo.NewClient(options.Client().ApplyURI(a.uri))
	if err != nil {
		return myApp{}, errors.Wrap(err, "failed to create mongo client")
	}
	a.client = client

	return a, nil
}

func main() {
	app, err := appInit(os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = app.client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer app.client.Disconnect(ctx)

	res, err := app.InsertPost(ctx, os.Args[0], os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)

	post, err := app.GetPost(ctx, bson.D{{"title", "Hello"}})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(post)

}
