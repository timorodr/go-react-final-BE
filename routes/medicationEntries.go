package routes

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/timorodr/go-react-final/server/models"
	"go.mongodb.org/mongo-driver/bson" // Binary JSON encodes type and length info which allows it to be traversed more quickly compared to JSON
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// "gopkg.in/mgo.v2/bson"

	"log"

	routes "github.com/timorodr/go-react-final/server/database"
	helper "github.com/timorodr/go-react-final/server/helpers"

	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New() // helps build strong and safe programs by checking info and making sure if follows the rule set
// Open Collection from connection.go
var entryCollection *mongo.Collection = routes.OpenCollection(routes.Client, "medications")

// Context is the most important part of gin. It allows us to pass variables between middleware, manage the flow, validate the JSON of a request and render a JSON response
var userCollection *mongo.Collection = routes.OpenCollection(routes.Client, "user")

// var validate = validator.New()

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect")
		check = false
	}

	return check, msg
}

// CreateUser is the api used to tget a single user
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		// defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}
		
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken
		user.Medications = make([]models.Medication, 0)

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, resultInsertionNumber)

	}
}

// Login is the api used to tget a single user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
			return
		}
		// defer cancel()

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, foundUser.User_id)
		defer cancel()

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		c.JSON(http.StatusOK, foundUser)

	}
}
//** user auth attempt
// func AddEntry(c *gin.Context) { // access to params and request through gin.Context
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	// var medication models.Medication
// 	// userID := c.MustGet("user_id").(string)
//   	// userIDObject, err := primitive.ObjectIDFromHex(userID)
//   	// if err != nil {
//     // 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
//     // 	return
//   	// }

// 	var medication models.Medication
// 	// medication.UserID = userID
// 	if err := c.BindJSON(&medication); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return // bind JSON serializes basically?
// 	}
// 	validationErr := validate.Struct(medication)
// 	if validationErr != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
// 		fmt.Println(validationErr)
// 		return
// 	}
// 	_, insertErr := entryCollection.InsertOne(ctx, medication)
// 	if insertErr != nil {
// 		msg := fmt.Sprintf("item was not created")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 		fmt.Println(insertErr)
// 		return
// 	}
// 	medication.ID = primitive.NewObjectID()

// 	_, err = entryCollection.UpdateByID(ctx, userIDObject, bson.M{"$push": bson.M{"medications": medication}})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	  }
	  
// 	defer cancel()
// 	// c.JSON(http.StatusOK, result)
// 	c.JSON(http.StatusOK, gin.H{"message": "Medication added successfully"})
// }

func AddEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Params.ByName("_id")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		userObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		var medications models.Medication
		if err := c.ShouldBindJSON(&medications); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		medications.Medication_id = primitive.NewObjectID()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userObjectID}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "medications", Value: medications}}}}

		result, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "medication added successfully", "result": result})
	}
}
//** latest func
// func AddEntry() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		user_id := c.Query("id")
// 		if user_id == "" {
// 			c.Header("Content-Type", "application/json")
// 			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
// 			c.Abort()
// 			return
// 		}

// 		medication, err := primitive.ObjectIDFromHex(user_id)
// 		if err != nil {
// 			c.IndentedJSON(500, "Internal Server Error")
// 		}
// 		var medications models.Medication
// 		medications.Medication_id = primitive.NewObjectID()
// 		if err = c.BindJSON(&medications); err != nil {
// 			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
// 		}

// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

// 		// var entry models.Medication

// 		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: medication}}}}
// 		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$medication"}}}}
// 		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$medication_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}


// 		pointcursor, err := userCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
// 		if err != nil {
// 			c.IndentedJSON(500, "Internal Server Error")
// 		}

// 		var medicationinfo []bson.M
// 		if err = pointcursor.All(ctx, &medicationinfo); err != nil {
// 			panic(err)
// 		}

// 		// var size int32
// 		// for _, medication_no := range medicationinfo {
// 		// 	count := medication_no["count"]
// 		// 	size = count.(int32)
// 		// }
// 		filter := bson.D{primitive.E{Key: "_id", Value: medication}}
// 			update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "medication", Value: medications}}}}
// 			result, err := userCollection.UpdateOne(ctx, filter, update)
// 			if err != nil {
// 				fmt.Println(err)
// 			}




// 		// if err := c.BindJSON(&entry); err != nil {
// 		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		// 	fmt.Println(err)
// 		// 	return
// 		// }
// 		// validationErr := validate.Struct(entry)
// 		// if validationErr != nil {
// 		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
// 		// 	fmt.Println(validationErr)
// 		// 	return
// 		// }
// 		// entry.ID = primitive.NewObjectID()
// 		// result, insertErr := entryCollection.InsertOne(ctx, entry)
// 		// if insertErr != nil {
// 		// 	msg := fmt.Sprintf("order item was not created")
// 		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 		// 	fmt.Println(insertErr)
// 		// 	return
// 		// }
// 		defer cancel()
// 		c.JSON(http.StatusOK, result)
// 		// ctx.Done()
// 	}
// }
func GetEntries(c *gin.Context) {
	userID := c.Params.ByName("_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User

	err = userCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, user.Medications)

}
//** Latest
// func GetEntries(c *gin.Context) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 	// userID := c.MustGet("user_id").(string)
// 	// filter := bson.M{"medications": bson.M{}}
// 	// userID, _ := c.Get("user_id")

// 	var entries []bson.M                     // M is an unordered representation of a BSON document. This type should be used when the order of the elements does not matter. This type is handled as a regular map[string]interface{} when encoding and decoding. Elements will be serialized in an undefined, random order.


	
// 	cursor, err := entryCollection.Find(ctx, bson.M{}) // passing through empty object you get all values if you want specific you must declare/specify
	
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}
	
// 	if err = cursor.All(ctx, &entries); err != nil {
// 		// c.JSON serializes the given struct as JSON into the response body - it also sets the Content-Type as "application/json"
// 		// process of converting a data structure or object into a format that can be easily stored or transmitted
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		fmt.Println(err)
// 		return
// 	}
// 	defer cancel()
// 	fmt.Println(entries)
// 	// fmt.Println(userID)
// 	c.JSON(http.StatusOK, entries)
// }

// func GetUserMedications(c *gin.Context) {
//     // Get user ID from the authenticated user's token or session
//     userID, _ := c.Get("user_id")

//     var medications []models.Medication

//     // Assuming db is your MongoDB client

//     // Filter medications by user ID
//     cursor, err := userCollection.Find(context.Background(), bson.M{"user_id": userID})
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch medications"})
//         return
//     }
//     defer cursor.Close(context.Background())

//     for cursor.Next(context.Background()) {
//         var medication models.Medication
//         if err := cursor.Decode(&medication); err != nil {
//             c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode medication"})
//             return
//         }
//         medications = append(medications, medication)
//     }

//     if err := cursor.Err(); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate over medications"})
//         return
//     }

//     c.JSON(http.StatusOK, medications)
// }

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
			"id":          primitive.NewObjectID(),
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

func Logout(c *gin.Context) {
	// Access user ID from context (assuming stored by Authentication middleware)
	userID := c.MustGet("user_id").(string)
  
	// Invalidate user session on server-side (e.g., remove tokens from database/cache)
	err := helper.InvalidateUserSession(userID)
	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "error invalidating session"})
	  return
	}
  
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
  }
