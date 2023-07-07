package database

import "context"

type Database interface {
	Collection(string) Collection
}

type Client interface {
	Database(string) Database
	Connect(context.Context) error
	Disconnect(context.Context) error
	Ping(context.Context) error
}

type Collection interface {
	UpsertOne(...UpsertionOption) (interface{}, error)
	FindMany(...SelectionOption) (MultiResult, error)
	FindOne(...SelectionOption) (SingleResult, error)
	DeleteOne(...DeletionOption) (interface{}, error)
}

type SingleResult interface {
	Decode(interface{}) error
}

type MultiResult interface {
	Decode(interface{}) ([]interface{}, error)
}
