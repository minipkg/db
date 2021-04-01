package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IDB is the interface for a DB connection
type IDB interface {
	Collection(name string, opts ...*options.CollectionOptions) ICollection
	Close(ctx context.Context) error
}

type ICollection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) ISingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (ICursor, error)
	InsertOne(ctx context.Context, document interface{}) (interface{}, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (interface{}, error)
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
}

type ISingleResult interface {
	Decode(val interface{}) error
}

type ICursor interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
}

type DB struct {
	client *mongo.Client
	db     *mongo.Database
}

func (b DB) Collection(name string, opts ...*options.CollectionOptions) ICollection {
	return &Collection{
		collection: b.db.Collection(name, opts...),
	}
}

type Collection struct {
	collection	*mongo.Collection
}

func (c Collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) ISingleResult {
	return &SingleResult{
		singleResult:	c.collection.FindOne(ctx, filter, opts...),
	}
}

func (c Collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (ICursor, error) {
	cursor, err := c.collection.Find(ctx, filter, opts...)
	return &Cursor{
		cursor:	cursor,
	}, err
}

func (c Collection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	res, err := c.collection.InsertOne(ctx, document)
	return res.InsertedID, err
}

func (c Collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (interface{}, error) {
	res, err := c.collection.UpdateOne(ctx, filter, update)
	return res.ModifiedCount, err
}

func (c Collection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	res, err := c.collection.DeleteOne(ctx, filter)
	return res.DeletedCount, err
}

type SingleResult struct {
	singleResult	*mongo.SingleResult
}

func (r SingleResult) Decode(val interface{}) error {
	return r.singleResult.Decode(val)
}

type Cursor struct {
	cursor	*mongo.Cursor
}

func (c Cursor) Next(ctx context.Context) bool {
	return c.cursor.Next(ctx)
}

func (c Cursor) Decode(val interface{}) error {
	return c.cursor.Decode(val)
}


var _ IDB = (*DB)(nil)
var _ ICollection = (*Collection)(nil)
var _ ISingleResult = (*SingleResult)(nil)
var _ ICursor = (*Cursor)(nil)


type Config struct {
	DSN    string
	DBName string
}

// New creates a new DB connection
func New(conf Config) (IDB, error) {
	// Create client
	client, err := mongo.NewClient(options.Client().ApplyURI(conf.DSN))
	if err != nil {
		return nil, err
	}

	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(conf.DBName)

	dbobj := &DB{db: db}
	return dbobj, nil
}

func (d *DB) Close(ctx context.Context) error {
	if err := d.client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}
