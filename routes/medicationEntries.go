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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "login or passowrd is incorrect"})
			return
		}

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


func AddEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Params.ByName("id")
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


func GetEntries(c *gin.Context) {
	userID := c.Params.ByName("id")
	fmt.Println(userID)
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


func UpdateEntry(c *gin.Context) {
    userID := c.Params.ByName("id")
    medicationID := c.Params.ByName("medication_id")
    
    if userID == "" || medicationID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user or medication ID"})
        return
    }

    userObjectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

	medicationIDObj, err := primitive.ObjectIDFromHex(medicationID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid medication ID"})
        return
    }

    var updateMedication models.Medication
    if err := c.BindJSON(&updateMedication); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        fmt.Println(err)
        return
    }

	filter := bson.M{
        "_id": userObjectID,
        "medications._id": medicationIDObj,
    }

    update := bson.D{
        {Key: "$set", Value: bson.M{
            "medications.$.name":        updateMedication.Name,
            "medications.$.dosage":      updateMedication.Dosage,
            "medications.$.description": updateMedication.Description,
        }},
    }

	ctx := context.TODO()
    result, err := userCollection.UpdateOne(ctx, filter, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update medication"})
        return
    }

	if result.ModifiedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "medication not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "medication updated successfully", "result": result})
}


func DeleteEntry(c *gin.Context) {
    userID := c.Params.ByName("id")
    medicationID := c.Params.ByName("medication_id")
	fmt.Println(medicationID, "user:", userID)
    
    if userID == "" || medicationID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user or medication ID"})
        return
    }

    userObjectID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }
    medicationObjID, err := primitive.ObjectIDFromHex(medicationID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
        return
    }

	filter := bson.M{
        "_id": userObjectID,
        "medications._id": medicationObjID,
    }

	update := bson.M{
        "$pull": bson.M{"medications": bson.M{"_id": medicationObjID}},
    }

	ctx := context.TODO()
    result, err := userCollection.UpdateOne(ctx, filter, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete medication"})
        return
    }

	if result.ModifiedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "medication not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "medication deleted successfully", "result": result})
}


func Logout(c *gin.Context) {
	// Access user ID from context
	userID := c.MustGet("user_id").(string)
  
	// Invalidate user session on server-side
	err := helper.InvalidateUserSession(userID)
	if err != nil {
	  c.JSON(http.StatusInternalServerError, gin.H{"error": "error invalidating session"})
	  return
	}
  
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
  }
