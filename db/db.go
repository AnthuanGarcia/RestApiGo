package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	model "github.com/AnthuanGarcia/IntegradoraII/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connTimeout        = 5
	connStringTemplate = "mongodb+srv://%s:%s@datos.fqmq7.mongodb.net/%s?retryWrites=true&w=majority"
	dataBase           = "PruebaGolang"
	collection         = "Datos"
)

// getConnection - Conexion a base de datos, por cliente de MongoDB
func getConnection() (*mongo.Client, context.Context, context.CancelFunc) {
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	clusterEndPoint := os.Getenv("MONGO_ENDPOINT")

	connectionURI := fmt.Sprintf(connStringTemplate, username, password, clusterEndPoint)

	ctx, cancel := context.WithTimeout(context.Background(), connTimeout*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionURI))

	if err != nil {
		log.Printf("Fallo para conectar al cliente: %v\n", err)
	}

	/*err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to cluster: %v\n", err)
	}*/

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Fallo ping al cluster: %v\n", err)
	}

	fmt.Printf("Conectado a MongoDB\n")

	return client, ctx, cancel
}

// GetAllTasks - Obtiene todos los documentos de la base de datos
func GetAllTasks() ([]*model.Task, error) {
	tasks := []*model.Task{}

	client, ctx, cancel := getConnection()

	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(dataBase)
	collection := db.Collection(collection)
	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	err = cursor.All(ctx, &tasks)

	if err != nil {
		log.Printf("Fallo Organizando: %v\n", err)
	}

	return tasks, nil
}

// GetTaskID - Obtiene un documento de la base de datos
func GetTaskID(id primitive.ObjectID) (*model.Task, error) {
	var task *model.Task

	client, ctx, cancel := getConnection()

	defer cancel()
	defer client.Disconnect(ctx)

	db := client.Database(dataBase)
	collection := db.Collection(collection)
	result := collection.FindOne(ctx, bson.D{})

	if result == nil {
		return nil, errors.New("No se Encontro la tarea")
	}

	err := result.Decode(&task)

	if err != nil {
		log.Printf("Fallo el ordenamiento %v", err)
		return nil, err
	}

	log.Printf("Tasks: %v", task)

	return task, nil
}

// Create - Crea un documento para la BD
func Create(task *model.Task) (primitive.ObjectID, error) {
	client, ctx, cancel := getConnection()

	defer cancel()
	defer client.Disconnect(ctx)

	task.ID = primitive.NewObjectID()

	result, err := client.Database(dataBase).Collection(collection).InsertOne(ctx, task)

	if err != nil {
		log.Printf("No se ha podido Crear la Tarea: %v", err)
		return primitive.NilObjectID, err
	}

	oid := result.InsertedID.(primitive.ObjectID)

	return oid, nil
}

// Update - Actualiza un documento de la BD
func Update(task *model.Task) (*model.Task, error) {
	var updatedTask *model.Task
	client, ctx, cancel := getConnection()

	defer cancel()
	defer client.Disconnect(ctx)

	update := bson.M{
		"$set": task,
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		Upsert:         &upsert,
		ReturnDocument: &after,
	}

	err := client.Database(dataBase).Collection(collection).FindOneAndUpdate(ctx, bson.M{"_id": task.ID}, update, &opt).Decode(&updatedTask)

	if err != nil {
		log.Printf("No se pudo almacenar la tarea: %v", err)
		return nil, err
	}

	return updatedTask, nil
}
