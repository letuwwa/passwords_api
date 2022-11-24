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
	"time"
)

func TokenCheck(c echo.Context) error {
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

	if err != nil {
		log.Printf("unable to find, username: %s", userData.Username)
		return c.JSON(http.StatusBadRequest, "bad request - unable  to find")
	}

	if !utils.IsPasswordHashValid(userData.Password, result["password"].(string)) {
		log.Print("invalid password")
		return c.JSON(http.StatusForbidden, "forbidden - invalid password")
	}

	if userData.ExpiredAt != "" {
		layout := "2006-01-02"
		date, err := time.Parse(layout, userData.ExpiredAt)
		if err != nil {
			log.Printf("error with parsing date: %s", err.Error())
			return c.JSON(http.StatusInternalServerError, "internal server error - date parsing")
		}

		if time.Now().Unix() > date.Unix() {
			log.Printf("password expired: %s", date)
			return c.JSON(http.StatusForbidden, "forbidden - password expired")
		}
	}

	return c.JSON(http.StatusOK, "ok - password is correct")
}
