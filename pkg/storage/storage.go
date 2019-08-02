package storage

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// DB is struct that holds database object
type DB struct {
	conn *mongo.Database
}

const (
	usersCollectionName       = "users"
	tournamentsCollectionName = "tournaments"
)

// CreateNew is constructor for db
func CreateNew(db *mongo.Database) *DB {
	return &DB{
		conn: db,
	}
}
