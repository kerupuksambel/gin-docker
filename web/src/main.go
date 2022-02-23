package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	// "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"

)

type Note struct {
	Title  	string 	`json:"title"`
	Content string  `json:"content"`
}

func db() *mongo.Client{
	clientOptions := options.Client().ApplyURI("mongodb://159.203.187.203:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

var noteCollection = db().Database("crud").Collection("notes")


func main() {

	// route := mux.NewRouter()
	// s := route.PathPrefix("/api").Subrouter() //Base Path

	// //Routes

	// s.HandleFunc("/note", get).Methods("GET")
	// // s.HandleFunc("/note/{id}", detail).Methods("GET")
	// // s.HandleFunc("/note", create).Methods("POST")
	// // s.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	// // s.HandleFunc("/note/{id}", delete).Methods("DELETE")

	// log.Fatal(http.ListenAndServe(":3004", s)) // Run Server

	router := gin.Default()
	router.GET("/notes", get)
	router.GET("/note/:id", detail)
	router.POST("/note", create)
	router.PUT("/note/:id", update)
	router.DELETE("/note/:id", delete)

	router.Run()
}


func get(c *gin.Context){
	var results []primitive.M

	opts := options.Find().SetProjection(bson.D{{"content", 0}})
	cur, err := noteCollection.Find(context.TODO(), bson.D{{}}, opts)


	if err != nil {
		fmt.Println(err)
	}
	for cur.Next(context.TODO()){
		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}

	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	c.JSON(http.StatusOK, gin.H{
		"results" : results,
	})
}

func detail(c *gin.Context){
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Print(err)
	}

	var result primitive.M

	err = noteCollection.FindOne(context.TODO(), bson.D{{"_id", _id}}).Decode(&result)

	if err != nil {
		fmt.Println(err)
	}
	// for cur.Next(context.TODO()){
	// 	var elem primitive.M
	// 	err := cur.Decode(&elem)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	results = append(results, elem)
	// }

	// cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	c.JSON(http.StatusOK, gin.H{
		"result" : result,
	})
}

func create(c *gin.Context){
	var n Note
	err := json.NewDecoder(c.Request.Body).Decode(&n) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
	}

	insertResult, err := noteCollection.InsertOne(context.TODO(), n)
	var status bool
	if err != nil {
		log.Fatal(err)
		status = false
	}else{
		status = true
	}

	c.JSON(http.StatusOK, gin.H{
		"success" : status,
		"id" : insertResult.InsertedID,
	})
}

func delete(c *gin.Context){
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Print(err)
	}

	opts := options.Delete().SetCollation(&options.Collation{})
	deleteResult, err := noteCollection.DeleteOne(context.TODO(), bson.D{{"_id", _id}}, opts)
	var status bool
	if err != nil {
		log.Fatal(err)
		status = false
	}else{
		status = true
		fmt.Printf("%v documents deleted", deleteResult.DeletedCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"success" : status,
	})
}

func update(c *gin.Context){
	id := c.Param("id")
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Print(err)
	}

	var body Note
	e := json.NewDecoder(c.Request.Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
	}

	update := bson.D{{"$set", bson.M{"title": body.Title, "content": body.Content}}}

	updateResult := noteCollection.FindOneAndUpdate(context.TODO(), bson.D{{"_id", _id}}, update)
	
	var status bool
	var result primitive.M
	_ = updateResult.Decode(&result)
	if updateResult == nil {
		log.Fatal(err)
		status = false
	}else{
		status = true
	}

	c.JSON(http.StatusOK, gin.H{
		"success" : status,
	})
}

// func GetProducts(c *gin.Context) {

// 	db, err := gorm.Open("sqlite3", "test.db")
// 	if err != nil {
// 		panic("failed to connect database")
// 	}
// 	defer db.Close()
// 	db.AutoMigrate(&Product{})

// 	var products []Product

// 	if err := db.Find(&products).Error; err != nil {
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		log.Println(err)
// 	} else {
// 		c.JSON(http.StatusOK, products)
// 		log.Println(products)
// 	}

// }
