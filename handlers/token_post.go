package handlers

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"passwords_api/mongo_db"
	"passwords_api/structures"
	"passwords_api/utils"
)

func TokenPost(c echo.Context) error {
	requestToken := new(structures.RequestToken)
	if err := c.Bind(requestToken); err != nil {
		return c.JSON(http.StatusBadRequest, "bad request - token")
	}

	token, err := utils.ValidateJWT(requestToken.TokenValue)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userData := structures.UserData{
		UserName:     claims["username"].(string),
		PasswordHash: claims["password_hash"].(string),
		DatabaseName: claims["database_name"].(string),
	}

	collection := mongo_db.MongoClient().Database(userData.DatabaseName).Collection("tokens")

	var result bson.M
	err = collection.FindOne(context.TODO(), bson.D{
		{"username", userData.UserName},
	}).Decode(&result)

	if result != nil {
		return c.JSON(http.StatusBadRequest, "bad request - already in use")
	}

	_, err = collection.InsertOne(context.TODO(), bson.D{
		{Key: "username", Value: userData.UserName},
		{Key: "password_hash", Value: userData.PasswordHash},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error - unable to save values")
	}

	return c.JSON(http.StatusCreated, "created - values were saved")
}
