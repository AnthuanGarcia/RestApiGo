package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task - Modelo Basico de Informacion
type Task struct {
	ID	  primitive.ObjectID
	Title string
	Body  string
}
