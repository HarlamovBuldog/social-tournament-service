package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/HarlamovBuldog/social-tournament-service/pkg/server"
	"github.com/HarlamovBuldog/social-tournament-service/pkg/storage"
	"gopkg.in/yaml.v3"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	ConnStr    string `yaml:"conn_str"`
	ServerPort string `yaml:"server_port"`
	DBName     string `yaml:"db_name"`
}

func (conf *config) Validate() error {
	if conf.ConnStr == "" {
		return errors.New("connection string is not provided")
	}
	if conf.DBName == "" {
		return errors.New("database name is not provided")
	}
	if conf.ServerPort == "" {
		return errors.New("server port is not provided")
	} else if _, err := strconv.ParseUint(conf.ServerPort, 0, 64); err != nil {
		return errors.Wrap(err, "bad server port provided")
	}
	return nil
}

func main() {
	yamlConfigFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error opening cofiguration yaml file: %v\n", err)
	}

	var conf config
	err = yaml.Unmarshal(yamlConfigFile, &conf)
	if err != nil {
		log.Fatalf("error unmarshaling yaml configuration file: %v\n", err)
	}

	if err = conf.Validate(); err != nil {
		log.Fatalf("error validating config file: %v", err)
	}

	clientOptions := options.Client().ApplyURI(conf.ConnStr)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("error connecting to mongo db: %v", err)
	}

	defer func() {
		err := client.Disconnect(context.TODO())
		if err != nil {
			log.Fatalf("error disconnecting from mongo db: %v", err)
		}
	}()

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("error connecting to mongo db: %v", err)
	}

	log.Println("Connected to MongoDB!")

	db := storage.CreateNew(client.Database(conf.DBName))
	s := server.NewServer(db)

	log.Fatal(http.ListenAndServe("localhost:"+conf.ServerPort, s))
}
