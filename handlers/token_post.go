package handlers

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"passwords_api/mongo_db"
	"passwords_api/structures"
	"passwords_api/utils"
)

func TokenPost(c echo.Context) error {
	requestToken := new(structures.RequestToken)
	if err := c.Bind(requestToken); err != nil {
		log.Printf("bind err: %s", err.Error())
		return c.JSON(http.StatusBadRequest, "bad request - token")
	}

	token, err := utils.ValidateJWT(requestToken.TokenValue)
	if err != nil {
		log.Printf("token err: %s", err.Error())
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userData := structures.UserData{
		Username: claims["username"].(string),
		Password: claims["password"].(string),
		Database: claims["database"].(string),
	}

	collection := mongo_db.MongoClient().Database(userData.Database).Collection("tokens")

	var result bson.M
	err = collection.FindOne(context.TODO(), bson.D{
		{"username", userData.Username},
	}).Decode(&result)

	if result != nil {
		log.Printf("in use err, username: %s", userData.Username)
		return c.JSON(http.StatusBadRequest, "bad request - already in use")
	}

	_, err = collection.InsertOne(context.TODO(), bson.D{
		{Key: "username", Value: userData.Username},
		{Key: "password", Value: userData.Password},
	})
	if err != nil {
		log.Printf("token save err: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, "internal server error - unable to save values")
	}

	return c.JSON(http.StatusCreated, "created - values were saved")
}
