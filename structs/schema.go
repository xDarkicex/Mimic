package structs

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// MongoUser ...
type MongoUser struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Username  string        `bson:"username"`
	Password  string        `bson:"password"`
	Session   []Session     `bson:"session"`
	LastLogin time.Time     `bson:"lastLogin"`
	Files     []File
}

// MongoFile ...
type MongoFile struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Size     int64         `bson:"size"`
	GridID   string        `bson:"grid_id"`
	Modified time.Time     `bson:"modified"`
}

// MongoSession ...
type MongoSession struct {
	ID string `bson:"_id,omitempty"`
	IP string `bson:"ips,omitempty"`
}
