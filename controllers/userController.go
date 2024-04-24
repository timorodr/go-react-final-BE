// package controllers

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"net/http"
// 	"time"
	
// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator/v10"
	
// 	"github.com/timorodr/go-react-final/server/routes"
// 	helper "github.com/timorodr/go-react-final/server/helpers"
	
// 	"github.com/timorodr/go-react-final/server/models"
	
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"golang.org/x/crypto/bcrypt"
// )

// var userCollection *mongo.Collection = routes.OpenCollection(routes.Client, "user")
// var validate = validator.New()

// // HashPassword is used to encrypt the password before it is stored in the DB
// func HashPassword(password string) string {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	return string(bytes)
// }

// // VerifyPassword checks the input password while verifying it with the passward in the DB.
// func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
// 	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
// 	check := true
// 	msg := ""

// 	if err != nil {
// 		msg = fmt.Sprintf("login or password is incorrect")
// 		check = false
// 	}

// 	return check, msg
// }

// // CreateUser is the api used to tget a single user
// func SignUp() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		var user models.User

// 		if err := c.BindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		validationErr := validate.Struct(user)
// 		if validationErr != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
// 			return
// 		}

// 		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
// 		defer cancel()
// 		if err != nil {
// 			log.Panic(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
// 			return
// 		}

// 		password := HashPassword(*user.Password)
// 		user.Password = &password

// 		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
// 		defer cancel()
// 		if err != nil {
// 			log.Panic(err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
// 			return
// 		}

// 		if count > 0 {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
// 			return
// 		}

// 		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// 		user.ID = primitive.NewObjectID()
// 		user.User_id = user.ID.Hex()
// 		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, user.User_id)
// 		user.Token = &token
// 		user.Refresh_token = &refreshToken

// 		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
// 		if insertErr != nil {
// 			msg := fmt.Sprintf("User item was not created")
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 			return
// 		}
// 		defer cancel()

// 		c.JSON(http.StatusOK, resultInsertionNumber)

// 	}
// }

// // Login is the api used to tget a single user
// func Login() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
// 		var user models.User
// 		var foundUser models.User

// 		if err := c.BindJSON(&user); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or passowrd is incorrect"})
// 			return
// 		}
// 		defer cancel()

// 		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
// 		if passwordIsValid != true {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 			return
// 		}
// 		defer cancel()

// 		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, foundUser.User_id)

// 		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

// 		c.JSON(http.StatusOK, foundUser)

// 	}
// }

// // AddMedicationForUser adds a medication to a user's medications list
// // func AddMedicationForUser(ctx context.Context, client *mongo.Client, userID string, medication Medication) (*User, error) {
// //     // Find the user by ID
// //     // userCollection := client.Database("your_database_name").Collection("users")
// //     filter := bson.M{"user_id": userID}
// //     var user models.Medication
// //     err := userCollection.FindOne(ctx, filter).Decode(&user)
// //     if err != nil {
// //         if err == mongo.ErrNoDocuments {
// //             return nil, fmt.Errorf("user with ID %s not found", userID)
// //         }
// //         return nil, fmt.Errorf("error finding user: %w", err)
// //     }

// //     // Append the medication to the user's Medications slice
// //     user.Medications = append(user.Medications, medication)

// //     // Update the user document in the database
// //     update := bson.M{"$set": bson.M{"medications": user.Medications}}
// //     _, err = userCollection.UpdateByID(ctx, userID, update)
// //     if err != nil {
// //         return nil, fmt.Errorf("error updating user medications: %w", err)
// //     }

// //     return &user, nil
// // }