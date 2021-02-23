package model

// Generator struct
type Generator struct {
	ID         int    `bson:"id"`
	Collection string `bson:"collection"`
}
