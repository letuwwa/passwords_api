package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"passwords_api/mongo_db"
	"passwords_api/structures"
	"passwords_api/utils"
)

func TokenReset(c echo.Context) error {
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
		filter := bson.D{{"name", userData.Username}}
		_, err := collection.ReplaceOne(context.TODO(), filter, userData)
		if err != nil {
			log.Printf("updating password err: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, "internal server error - error during updating password")
		}
		return c.JSON(http.StatusOK, "ok - user password was updated")
	}

	return c.JSON(http.StatusInternalServerError, "internal server error - user not found")
}
