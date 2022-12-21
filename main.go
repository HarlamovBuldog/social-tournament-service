package main

import (
	"context"
	"log"
	"net"

	"github.com/HarlamovBuldog/social-tournament-service/internal/server/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

/*
func main() {
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
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("error connecting to mongo db: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
		defer cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("error disconnecting from mongo db: %v", err)
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
*/

type server struct {
	pb.UnimplementedTournamentServer
}

func (s *server) UserList(ctx context.Context, in *pb.GetUserListRequest) (*pb.GetUserListResponse, error) {
	return &pb.GetUserListResponse{Users: getUserList()}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterTournamentServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getUserList() []*pb.User {
	return []*pb.User{
		{
			Name: "Stas",
			Age:  12,
		},
		{
			Name: "Vlad",
			Age:  25,
		},
	}
}
