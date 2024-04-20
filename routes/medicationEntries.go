package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/timorodr/go-react-final/server/models"
	"go.mongodb.org/mongo-driver/bson" // Binary JSON encodes type and length info which allows it to be traversed more quickly compared to JSON
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/go-playground/validator/v10"
	// "gopkg.in/mgo.v2/bson"
)

var validate = validator.New() // helps build strong and safe programs by checking info and making sure if follows the rule set
// Open Collection from connection.go
var entryCollection *mongo.Collection = OpenCollection(Client, "medications")

// Context is the most important part of gin. It allows us to pass variables between middleware, manage the flow, validate the JSON of a request and render a JSON response

func AddEntry(c *gin.Context) { // access to params and request through gin.Context
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var medication models.Medication

	if err := c.BindJSON(&medication); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return // bind JSON serializes basically?
	}
	validationErr := validate.Struct(medication)
	if validationErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	medication.ID = primitive.NewObjectID()
	result, insertErr := entryCollection.InsertOne(ctx, medication)
	if insertErr != nil {
		msg := fmt.Sprintf("item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result)
}

func GetEntries(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var entries []bson.M                               // M is an unordered representation of a BSON document. This type should be used when the order of the elements does not matter. This type is handled as a regular map[string]interface{} when encoding and decoding. Elements will be serialized in an undefined, random order.
	cursor, err := entryCollection.Find(ctx, bson.M{}) // passing through empty object you get all values if you want specific you must declare/specify

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &entries); err != nil {
		// c.JSON serializes the given struct as JSON into the response body - it also sets the Content-Type as "application/json"
		// process of converting a data structure or object into a format that can be easily stored or transmitted
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(entries)
	c.JSON(http.StatusOK, entries)
}

func GetEntryById(c *gin.Context) {
	EntryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(EntryID) // primistive BSON package helps us with ID's

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var entry bson.M

	if err := entryCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	fmt.Println(entry)
	c.JSON(http.StatusOK, entry)
}

// func GetEntriesByIngredient(c *gin.Context) {
// 	ingredient := c.Params.ByName("id") // gets ingredient id
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	var entries []bson.M // splice bringing all in this var

// 	cursor, err := entryCollection.Find(ctx, bson.M{"ingredients": ingredient}) // instead of passing empty object 
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}
// 	if err = cursor.All(ctx, &entries); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}
// 	defer cancel()
// 	fmt.Println(entries)


// 	c.JSON(http.StatusOK, entries)

// }

// func UpdateIngredient(c *gin.Context) {
// 	entryID := c.Params.ByName("id")
// 	docID, _ := primitive.ObjectIDFromHex(entryID)
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 	type Ingredient struct {
// 		Ingredients *string `json: "ingredients"` // sending this deferring with * - JSON wills look like this sending this
// 	}

// 	var ingredient Ingredient

// 	if err := c.BindJSON(&ingredient); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return // bind JSON serializes basically?
// 	}

// 	result, err := entryCollection.UpdateOne(ctx, bson.M{"_id": docID},
// 		bson.D{{"$set", bson.D{{"ingredients", ingredient.Ingredients}}}}, // $ mongoDB .D ordered rep of a BSON doc - with ingredients and the ingredient var that refers to the struct that holds Ingredients JSON *string
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}
// 	defer cancel()
// 	c.JSON(http.StatusOK, result.ModifiedCount) // number of Docs modified by the operation
// }

func UpdateEntry(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var medication models.Medication

	if err := c.BindJSON(&medication); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return // bind JSON serializes basically?
	}

	validationErr := validate.Struct(medication)
	if validationErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}

	result, err := entryCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docID},
		bson.M{
			"id":		   primitive.NewObjectID(),
			"name":        medication.Name,
			"dosage":      medication.Dosage,
			"description": medication.Description,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, result.ModifiedCount)
}

func DeleteEntry(c *gin.Context) {
	entryID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(entryID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := entryCollection.DeleteOne(ctx, bson.M{"_id": docID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
	}

	defer cancel()
	c.JSON(http.StatusOK, result.DeletedCount)
}
