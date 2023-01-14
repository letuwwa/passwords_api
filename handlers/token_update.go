package handlers

import (
	"context"

	"log"
	"net/http"
	"passwords_api/mongo_db"
	"passwords_api/structures"
	"passwords_api/utils"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func TokenUpdate(c echo.Context) error {
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
		encryptedPass, encryptErr := utils.EncryptPassword(userData.Password)
		if encryptErr != nil {
			log.Printf("encrypt err: %s", encryptErr)
			return c.JSON(http.StatusInternalServerError, "internal server error - encrypt error")
		}

		filter := bson.D{{"username", userData.Username}}

		_, err := collection.ReplaceOne(context.TODO(), filter, bson.D{
			{Key: "password", Value: encryptedPass},
			{Key: "username", Value: userData.Username},
			{Key: "expired_at", Value: userData.ExpiredAt},
		})
		if err != nil {
			log.Printf("updating password err: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, "internal server error - error during updating password")
		}

		return c.JSON(http.StatusOK, "ok - user password was updated")
	}

	return c.JSON(http.StatusInternalServerError, "internal server error - user not found")
}
