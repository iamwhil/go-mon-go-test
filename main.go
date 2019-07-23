package main

import (
	"fmt"
	"log"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age int
	City string
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	handle_error(err)

	err = client.Ping(context.TODO(), nil)
	handle_error(err)

	fmt.Println("Connected to MONGO!")
	collection := client.Database("test").Collection("trainers")

	// ash := Trainer{"Ash", 10, "Pallet Town"}
	// misty := Trainer{"Misty", 10, "Cerulean City"}
	// brock := Trainer{"Brock", 15, "Pewter City"}

	// trainers := []interface{}{ash, misty, brock}

	// insertManyResults, err := collection.InsertMany(context.TODO(), trainers)
	// handle_error(err)

	// fmt.Println("Inserted multiple documents", insertManyResults.InsertedIDs)

	filter := buildFilter("Ash")
	update := buildUpdate()
	updateTrainer(collection, filter, update)

	var result Trainer 
	retrieveTrainer(collection, &result, filter)

	retrieveTrainers(collection, filter)

	fmt.Println(hasPermission(collection, filter))

	disconnectFromMongo(client)

}

func handle_error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func buildFilter(name string) bson.D {
	return bson.D{{"name", name}}
}

func buildUpdate() bson.D {
	return bson.D{
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
}

func updateTrainer(collection *mongo.Collection, filter, update bson.D) {
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	handle_error(err)
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}

func retrieveTrainer(collection *mongo.Collection, res *Trainer, filter bson.D) {
	err := collection.FindOne(context.TODO(), filter).Decode(&res)
	handle_error(err)
	fmt.Printf("Found a single document: %v\n", res)
}

func retrieveTrainers(collection *mongo.Collection, filter bson.D) {
	findOptions := options.Find()
	findOptions.SetLimit(2)

	var results[]*Trainer

	cursor, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	handle_error(err)

	for cursor.Next(context.TODO()){
		var element Trainer
		err := cursor.Decode(&element)
		handle_error(err)
		results = append(results, &element)
	}
	err = cursor.Err()
	handle_error(err)

	cursor.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %v\n", results)
	fmt.Println(len(results))
}

func hasPermission(collection *mongo.Collection, filter bson.D) bool {
	count, err := collection.CountDocuments(context.TODO(), filter)
	handle_error(err)
	if count > 0 {
		return true
	} 
	return false
}

func disconnectFromMongo(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	handle_error(err)
	fmt.Println("Disconnected from Mogno...")
}