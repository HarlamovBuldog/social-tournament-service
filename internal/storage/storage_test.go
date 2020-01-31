package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client      *mongo.Client
	db          *DB
	users       *mongo.Collection
	tournaments *mongo.Collection
)

const (
	dbName = "sts"
	dbPort = "27017/tcp"
)

type Container struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
}

func NewContainer(pool, repository, tag string, env []string) (*Container, error) {
	p, err := dockertest.NewPool(pool)
	if err != nil {
		return nil, errors.Wrap(err, "error create new docker pool")
	}

	opts := &dockertest.RunOptions{
		Name:       "mongo",
		Repository: repository,
		Tag:        tag,
		Env:        env,
		Cmd: []string{
			"mongod", "--replSet", "rs0",
		},
	}
	r, err := p.RunWithOptions(opts)
	if err != nil {
		return nil, errors.Wrap(err, "error running docker container")
	}

	return &Container{
		pool:     p,
		resource: r,
	}, nil
}

func (c *Container) Purge() error {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
		defer cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			return errors.Wrap(err, "error disconnect from mongo client")
		}
	}
	return c.pool.Purge(c.resource)
}

func (c *Container) GetBindedPort(p string) string {
	return c.resource.GetPort(p)
}

func TestMain(m *testing.M) {
	c, err := setupDB()
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()

	if err := c.Purge(); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func setupDB() (*Container, error) {
	c, err := NewContainer("", "mongo", "4.0.12-xenial", []string{})
	if err != nil {
		return nil, err
	}

	exec, err := c.pool.Client.CreateExec(docker.CreateExecOptions{
		Cmd: []string{
			"mongo", "--eval", "rs.initiate({_id: 'rs0', members: [{ _id: 0, host: 'localhost:27017' }]});db = db.getSiblingDB('sts')",
		},
		Container:    c.resource.Container.ID,
		AttachStdout: true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Create exec")
	}

	time.Sleep(2 * time.Second)
	err = c.pool.Client.StartExec(exec.ID, docker.StartExecOptions{
		Detach:       false,
		OutputStream: os.Stdout,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Start exec")
	}

	posPort := c.GetBindedPort(dbPort)

	i := 0
	for {
		if i >= 60 {
			return nil, errors.New("docker start time-out")
		}
		i++

		time.Sleep(2 * time.Second)

		writeConnectionString := fmt.Sprintf("mongodb://localhost:%s/?connect=direct&replicaSet=rs0&retryWrites=true", posPort)

		clientOptions := options.Client().ApplyURI(writeConnectionString)
		client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Printf("error connecting to mongo client: %s", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Printf("error ping mongo client: %s", err)
			continue
		}

		log.Println("Connected to MongoDB!")
		db = CreateNew(client.Database(dbName))
		users = client.Database(dbName).Collection(usersCollectionName)
		tournaments = client.Database(dbName).Collection(tournamentsCollectionName)

		break
	}

	return c, nil
}

func cleanUp(t *testing.T) {
	err := users.Drop(context.TODO())
	require.NoError(t, err)

	err = tournaments.Drop(context.TODO())
	require.NoError(t, err)
}
