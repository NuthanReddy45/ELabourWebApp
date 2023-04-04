package models

// import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	// ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
	PhnNo    string `json:"phnNo"`
}
