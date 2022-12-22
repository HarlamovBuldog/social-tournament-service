package cmd

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/HarlamovBuldog/social-tournament-service/internal/pkg/protocol/grpc"
	v1 "github.com/HarlamovBuldog/social-tournament-service/internal/pkg/service/v1"
	"github.com/HarlamovBuldog/social-tournament-service/internal/pkg/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

// TODO: move config to separate pkg.
type config struct {
	ConnStr    string `yaml:"conn_str"`
	ServerPort int32  `yaml:"server_port"`
	DBName     string `yaml:"db_name"`
}

// Validate checks if all config values are set.
func (conf *config) Validate() error {
	if conf.ConnStr == "" {
		return errors.New("connection string is not provided")
	}
	if conf.DBName == "" {
		return errors.New("database name is not provided")
	}
	if conf.ServerPort < 0 || conf.ServerPort > 65535 {
		return errors.New("bad server port provided")
	}

	return nil
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	yamlConfigFile, err := os.ReadFile("config.yaml")
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
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("error connecting to mongo db: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("error disconnecting from mongo db: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("error connecting to mongo db: %v", err)
	}

	log.Println("Connected to MongoDB!")
	db := storage.CreateNew(client.Database(conf.DBName))

	servPort := strconv.FormatInt(int64(conf.ServerPort), 10)
	v1API := v1.NewToDoServiceServer(db)

	return grpc.RunServer(ctx, v1API, servPort)
}
