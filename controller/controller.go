package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Grimmer98/mongoapi/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const connectionString = "mongodb+srv://kernelwhiterose:yoJeASZLA8ePDee2@cluster0.wntxf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
const dbName = "netflix"
const collectionName = "watchlist"

// IMPORTANT
var collection *mongo.Collection

// connect with mongoDB
func init() {
	//client option
	clientOption := options.Client().ApplyURI(connectionString)

	//connect to mongoDB
	client, err := mongo.Connect(clientOption)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB connection success")
	collection = client.Database(dbName).Collection(collectionName)

	//collection instance
	fmt.Println("Collection instance is ready")
}

//MONGODB helpers - file

//insert 1 record

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted one movie in db with id: ", inserted.InsertedID)
}

func updateOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Printf("Invalid ObjectID: %v", movieId)

		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Update failed")
		return
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

func deleteOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Printf("Invalid ObjectID: %v", movieId)

		return
	}
	filter := bson.M{"_id": id}

	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("Deletion failed")
		return
	}

	fmt.Println("Movie got deleted with delete count: ", deleteCount.DeletedCount)

}

func deleteAllMovies() int64 {

	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Number of movies deleted: ", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

func getAllMovies() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie primitive.M
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	defer cur.Close(context.Background())
	return movies
}

func GetMyAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")

	var movie model.Netflix
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Fatal(err)
	}
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")

	count := deleteAllMovies()
	json.NewEncoder(w).Encode(count)
}
