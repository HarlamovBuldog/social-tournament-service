package storage

import (
	"context"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
)

type config struct {
	ConnStr string `yaml:"conn_str"`
	DBName  string `yaml:"db_name"`
}

var db *DB
var users *mongo.Collection

func TestMain(m *testing.M) {
	conf, err := populateConfig()
	if err != nil {
		log.Fatalf("error populating config: %v", err)
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
			return
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("error connecting to mongo db: %v", err)
		return
	}

	log.Println("Connected to MongoDB!")

	db = CreateNew(client.Database(conf.DBName))
	users = client.Database(conf.DBName).Collection(usersCollectionName)

	m.Run()
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

	if err := dropUsersCollection(); err != nil {
		t.Fatalf("error clear users collection: %v", err)
	}
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

	if err := dropUsersCollection(); err != nil {
		t.Fatalf("error clear users collection: %v", err)
	}
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

	if err := dropUsersCollection(); err != nil {
		t.Fatalf("error clear users collection: %v", err)
	}
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

func dropUsersCollection() error {
	if err := users.Drop(context.TODO()); err != nil {
		return errors.Wrap(err, "drop users collection")
	}

	return nil
}

func (conf *config) Validate() error {
	if conf.ConnStr == "" {
		return errors.New("connection string is not provided")
	}
	if conf.DBName == "" {
		return errors.New("database name is not provided")
	}

	return nil
}

func populateConfig() (*config, error) {
	yamlConfigFile, err := ioutil.ReadFile("../../test_data/config.yaml")
	if err != nil {
		return nil, errors.Wrap(err, "open config file")
	}

	var conf config
	err = yaml.Unmarshal(yamlConfigFile, &conf)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal config file")
	}

	if err = conf.Validate(); err != nil {
		return nil, errors.Wrap(err, "validate config file")
	}

	return &conf, nil
}
