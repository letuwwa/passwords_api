package handlers

import (
	"context"
	"log"
	"net/http"
	"passwords_api/cmd/mongo_db"
	"passwords_api/cmd/structures"
	"passwords_api/cmd/utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func TokenPost(c echo.Context) error {
	requestToken := new(structures.RequestToken)
	if err := c.Bind(requestToken); err != nil {
		log.Printf("bind err: %s", err.Error())
		return c.JSON(http.StatusBadRequest, "bad request - invalid token request")
	}

	token, err := utils.ValidateJWT(requestToken.TokenValue)
	if err != nil {
		log.Printf("token err: %s", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userData := structures.UserData{}.FromToken(token)
	collection := mongo_db.MongoClient().Database(userData.Database).Collection("tokens")

	var result bson.M
	err = collection.FindOne(context.TODO(), bson.D{
		{"username", userData.Username},
	}).Decode(&result)

	if result != nil {
		log.Printf("in use err, username: %s", userData.Username)
		return c.JSON(http.StatusBadRequest, "bad request - already in use")
	}

	encryptedPass, encryptErr := utils.EncryptPassword(userData.Password)
	if encryptErr != nil {
		log.Printf("encrypt err: %s", encryptErr)
		return c.JSON(http.StatusInternalServerError, "internal server error - encrypt error")
	}

	_, err = collection.InsertOne(context.TODO(), bson.D{
		{Key: "password", Value: encryptedPass},
		{Key: "username", Value: userData.Username},
		{Key: "expired_at", Value: userData.ExpiredAt},
	})
	if err != nil {
		log.Printf("token save err: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error - unable to save values")
	}

	return c.JSON(http.StatusCreated, "created - values were saved")
}
