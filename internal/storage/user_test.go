package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	db     *DB
	users  *mongo.Collection
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

	r, err := p.Run(repository, tag, env)
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

	posPort := c.GetBindedPort(dbPort)

	i := 0
	for {
		if i >= 60 {
			return nil, errors.New("docker start time-out")
		}
		i++

		time.Sleep(5 * time.Second)

		writeConnectionString := fmt.Sprintf("mongodb://localhost:%s", posPort)

		clientOptions := options.Client().ApplyURI(writeConnectionString)
		client, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Printf("error connecting to mongo client: %s", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Printf("error ping mongo client: %s", err)
			continue
		}

		log.Println("Connected to MongoDB!")
		db = CreateNew(client.Database(dbName))
		users = client.Database(dbName).Collection(usersCollectionName)

		break
	}

	return c, nil
}

func cleanUp(t *testing.T) {
	err := client.Database(dbName).Collection(usersCollectionName).Drop(context.TODO())
	require.NoError(t, err)

	err = client.Database(dbName).Collection(tournamentsCollectionName).Drop(context.TODO())
	require.NoError(t, err)
}
func TestAddUser(t *testing.T) {
	userIDExpected, err := db.AddUser(context.TODO(), "gennadiy")
	if err != nil {
		t.Fatal(err)
	}

	actualUser, err := db.GetUser(context.TODO(), userIDExpected)
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)
	assert.Equal(userIDExpected, actualUser.ID.Hex(), "The two IDs should be the same.")

	userIDExpected2, err := primitive.ObjectIDFromHex(userIDExpected)
	if err != nil {
		t.Fatalf("error converting to primitive.ObjectID: %v", err)
	}

	assert.Equal(&User{ID: userIDExpected2, Name: "gennadiy"}, actualUser, "The two users should be the same.")

	cleanUp(t)
}

func TestGetUser(t *testing.T) {
	userIDExpected, err := db.AddUser(context.TODO(), "Vasya")
	if err != nil {
		t.Fatal(err)
	}

	actualUser, err := db.GetUser(context.TODO(), userIDExpected)
	if err != nil {
		t.Fatal(err)
	}
	assert := assert.New(t)
	assert.Equal(userIDExpected, actualUser.ID.Hex(), "The two IDs should be the same.")

	userIDExpected2, err := primitive.ObjectIDFromHex(userIDExpected)
	if err != nil {
		t.Fatalf("error converting to primitive.ObjectID: %v", err)
	}

	assert.Equal(&User{ID: userIDExpected2, Name: "Vasya"}, actualUser, "The two users should be the same.")

	badUserID := "safasf2412"
	_, err = db.GetUser(context.TODO(), badUserID)
	assert.EqualError(err,
		"convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+0073 's'",
		"The error should contain text")

	notExistUserID := primitive.NewObjectID()
	actualUser, err = db.GetUser(context.TODO(), notExistUserID.Hex())
	assert.EqualError(err, "decode returned doc: "+mongo.ErrNoDocuments.Error(), "The two errors should be the same")

	cleanUp(t)
}

func TestDeleteUser(t *testing.T) {
	userIDExpected, err := db.AddUser(context.TODO(), "Vasya")
	if err != nil {
		t.Fatal(err)
	}

	if err = db.DeleteUser(context.TODO(), userIDExpected); err != nil {
		t.Fatal(err)
	}

	badUserID := "safasf2412"
	err = db.DeleteUser(context.TODO(), badUserID)
	assert := assert.New(t)
	assert.EqualError(err,
		"convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+0073 's'",
		"The error should contain text")

	err = db.DeleteUser(context.TODO(), userIDExpected)
	assert.EqualError(err, "delete doc from collection: DeletedCount != 1", "The two errors should be the same")

	cleanUp(t)
}

func TestTakeUserBalance(t *testing.T) {
	generatedUserID := primitive.NewObjectID()
	amount := 100.0

	err := db.TakeUserBalance(context.TODO(), generatedUserID.Hex(), amount)
	assert := assert.New(t)
	assert.EqualError(err,
		"update doc in collection: ModifiedCount != 1",
		"The error should contain text")

	badUserID := "safasf2412"
	err = db.TakeUserBalance(context.TODO(), badUserID, amount)
	assert.EqualError(err,
		"convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+0073 's'",
		"The error should contain text")

	addedUserID, err := db.AddUser(context.TODO(), "Vasya")
	if err != nil {
		t.Fatal(err)
	}

	err = db.TakeUserBalance(context.TODO(), addedUserID, amount)
	if err != nil {
		t.Fatal(err)
	}

	addedUser, err := db.GetUser(context.TODO(), addedUserID)
	if err != nil {
		t.Fatal(err)
	}

	addedUserObjectID, err := primitive.ObjectIDFromHex(addedUserID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(&User{ID: addedUserObjectID, Name: "Vasya", Balance: -100.0}, addedUser,
		"The two users should be the same.")

}

func TestFundUserBalance(t *testing.T) {
	generatedUserID := primitive.NewObjectID()
	amount := 100.0

	err := db.FundUserBalance(context.TODO(), generatedUserID.Hex(), amount)
	assert := assert.New(t)
	assert.EqualError(err,
		"update doc in collection: ModifiedCount != 1",
		"The error should contain text")

	badUserID := "safasf2412"
	err = db.FundUserBalance(context.TODO(), badUserID, amount)
	assert.EqualError(err,
		"convert string value to primitive.ObjectID type: encoding/hex: invalid byte: U+0073 's'",
		"The error should contain text")

	addedUserID, err := db.AddUser(context.TODO(), "Vasya")
	if err != nil {
		t.Fatal(err)
	}

	err = db.FundUserBalance(context.TODO(), addedUserID, amount)
	if err != nil {
		t.Fatal(err)
	}

	addedUser, err := db.GetUser(context.TODO(), addedUserID)
	if err != nil {
		t.Fatal(err)
	}

	addedUserObjectID, err := primitive.ObjectIDFromHex(addedUserID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(&User{ID: addedUserObjectID, Name: "Vasya", Balance: 100.0}, addedUser,
		"The two users should be the same.")
}
