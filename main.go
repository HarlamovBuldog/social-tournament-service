package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/HarlamovBuldog/social-tournament-service/internal/server"
	"github.com/HarlamovBuldog/social-tournament-service/internal/storage"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

type config struct {
	ConnStr    string `yaml:"conn_str"`
	ServerPort int32  `yaml:"server_port"`
	DBName     string `yaml:"db_name"`
}

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
		ctx, _ := context.WithTimeout(context.TODO(), time.Second*5)
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatalf("error disconnecting from mongo db: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("error connecting to mongo db: %v", err)
	}

	log.Println("Connected to MongoDB!")
	db := storage.CreateNew(client.Database(conf.DBName))

	servPort := strconv.FormatInt(int64(conf.ServerPort), 10)

	srv := &http.Server{
		Addr:    ":" + servPort,
		Handler: server.NewServer(db),
	}

	go func() {
		// returns ErrServerClosed on graceful close
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe(): %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Print("Server shutting down...")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("error shutdown server: %s", err)
	}

	log.Print("Server stopped")
}
