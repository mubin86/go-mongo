package handler

import (
	"context"
	"fmt"
  "os"
	"log"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string            `json:"title,omitempty" bson:"title,omitempty"`
	Description  string     `json:"description,omitempty" bson:"description,omitempty"`
}

//context
var ctx = func() context.Context {
	return context.Background()
//	return context.WithTimeout(context.Background(), 10*time.Second)
}()

var dbUserName = os.Getenv("DB_USERNAME")
var dbPassword = os.Getenv("DB_PASSWORD")

//connect database
func connect() (*mongo.Database, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://dbUserName:dbPassword@cluster0.4xaod.mongodb.net/gomongo"))
	if err != nil {
		return nil, err
	}
	
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("connect database successfully")
	
	return client.Database("gomongo"), nil
}

func CreateProduct(c *gin.Context)  {

	db, _ := connect()

	//var reqBody Product 
	product := new(Product)
	if err :=c.BindJSON(&product); err != nil {
			c.JSON(422, gin.H{
				"error":   true,
				"message": "invalid request body",
			})
			return
		}
		
		fmt.Println(product.Title)
		
		res,_ :=db.Collection("product").InsertOne(ctx, bson.M{"title": product.Title, "description": product.Description})

		fmt.Println(res.InsertedID)

	_ = db.Collection("product").FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&product)

		c.JSON(200, gin.H{
		"message": "successfully added product",
		"data": product,
		})
	}


	func GetProducts(c *gin.Context)  {
		db, _ := connect()
	
		cur, err := db.Collection("product").Find(ctx, bson.D{})
	
	fmt.Println(cur)
	 if err != nil {
		log.Fatal(err)
		c.JSON(404, gin.H{
			"error":   true,
			"message": "something went wrong",
		})
		return
	}	

	defer cur.Close(ctx)

	result := make([]Product, 0)
	for cur.Next(ctx) {
		var row Product
		err := cur.Decode(&row)
		if err != nil {
			c.JSON(404, gin.H{
				"error":   true,
				"message": "something went wrong",
			})
			return
		}
		result = append(result, row)
	}
	fmt.Println(result)
			c.JSON(200, gin.H{
			"message": "get all products",
			"data": result,
			})		
}
	
func SingleProduct(c *gin.Context)  {
	db, _ := connect()

	id := c.Param("id")
	fmt.Println(id)
	
	_id, _ := primitive.ObjectIDFromHex(id)
	fmt.Println(_id)

	product := new(Product)

  err := db.Collection("product").FindOne(ctx, bson.M{"_id": _id}).Decode(&product)

  fmt.Println(*product)
	
	if err != nil {
	log.Fatal(err)
	c.JSON(404, gin.H{
		"error":   true,
		"message": "not found",
	})
	return
}

c.JSON(200, gin.H{
	"message": "success",
	"data": product,
})
}

	
func UpdateProduct(c *gin.Context)  {
	db, _ := connect()

	product := new(Product)
	if err :=c.BindJSON(&product); err != nil {
			c.JSON(422, gin.H{
				"error":   true,
				"message": "invalid request body",
			})
			return
		}

	id := c.Param("id")
	_id, _ := primitive.ObjectIDFromHex(id)

	filter := bson.M{"_id": _id}

  _,err := db.Collection("product").UpdateOne(ctx, filter, bson.M{"$set": product})
	if err != nil {
		c.JSON(404, gin.H{
			"error":   true,
			"message": "something went wrong",
		})
		return
	}
  err2 := db.Collection("product").FindOne(ctx, filter).Decode(&product)
	
	if err2 != nil {
	c.JSON(404, gin.H{
		"error": true,
		"message": "not found",
	})
	return
}

c.JSON(200, gin.H{
	"message": "succesfully updated",
	"data": product,
})

}
