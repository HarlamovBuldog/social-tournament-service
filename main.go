package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/HarlamovBuldog/social-tournament-service/pkg/server"
	"github.com/HarlamovBuldog/social-tournament-service/pkg/storage"
	"gopkg.in/yaml.v3"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	ConnStr string `yaml:"connstr"`
	Port    string `yaml:"port"`
	DBName  string `yaml:"dbname"`
}

func main() {
	yamlConfigFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("error opening cofiguration yaml file: %v\n", err)
		return
	}
	var conf config
	err = yaml.Unmarshal(yamlConfigFile, &conf)
	if err != nil {
		log.Printf("error unmarshaling yaml configuration file: %v\n", err)
		return
	}
	if conf.ConnStr == "" {
		log.Println("error in yaml configuration file: connection string is not provided")
		return
	}
	clientOptions := options.Client().ApplyURI(conf.ConnStr)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	defer client.Disconnect(context.TODO())

	if err != nil {
		log.Printf("error connecting to mongo db: %v", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Printf("error connecting to mongo db: %v", err)
	}

	log.Println("Connected to MongoDB!")

	db := storage.CreateNew(client.Database("sts"))
	s := server.NewServer(db)
	if conf.Port == "" {
		log.Println("error in yaml configuration file: port is not provided")
		return
	}
	log.Fatal(http.ListenAndServe("localhost"+conf.Port, s))
}
