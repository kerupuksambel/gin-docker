package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"

)

type Note struct {
	ID    	uint   	`json:"id"`
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
	r := gin.Default()

	r.GET("/note", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var body Note
		var result primitive.M
		
		e := json.NewDecoder(r.Body).Decode(&body)
		if e != nil{
			// Kalau tidak ada JSON, tampilkan semua
			err := noteCollection.Find(context.TODO(), bson.D{}).Decode(&result)
		}else{
			if body.ID == nil{
				fmt.Println("Request invalid")
			}else{
				// Kalau ada JSON valid, tampilkan hanya ID
				err := noteCollection.FindOne(context.TODO(), bson.D{{"id", body.ID}}).Decode(&result)
			}
		}

		if err != nil{
			fmt.Println(err)
		}

		json.NewEncoder(w).Encode(result)
	})

	r.GET("/echo/:echo", func(c *gin.Context) {
		echo := c.Param("echo")
		c.JSON(http.StatusOK, gin.H{
			"echo": echo,
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			// c.SaveUploadedFile(file, dst)
		}
		c.JSON(http.StatusOK, gin.H{
			"uploaded": len(files),
		})
	})

	// r.GET("/products", GetProducts)

	r.Run() // listen and serve on 0.0.0.0:8080
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
