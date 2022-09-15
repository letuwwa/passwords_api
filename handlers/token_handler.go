package handlers

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"passwords_api/db"
	"passwords_api/utils"
)

type RequestToken struct {
	TokenValue string `json:"token"`
}

type UserData struct {
	UserName     string `json:"user_name"`
	PasswordHash string `json:"password_hash"`
}

func TokenPost(c echo.Context) error {
	coll := db.MongoClient().Database("mongo_app").Collection("tokens")
	requestToken := new(RequestToken)
	if err := c.Bind(requestToken); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	token, _ := utils.ValidateJWT(requestToken.TokenValue)
	claims := token.Claims.(jwt.MapClaims)
	userData := UserData{
		UserName:     claims["user_name"].(string),
		PasswordHash: claims["password_hash"].(string),
	}
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"user_name", userData.UserName}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			doc := bson.D{{"user_name", userData.UserName}, {"password_hash", userData.PasswordHash}}
			_, err := coll.InsertOne(context.TODO(), doc)
			if err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}
			return c.JSON(http.StatusCreated, "created")
		}
	}
	if userData.PasswordHash != result["password_hash"] {
		return c.JSON(http.StatusForbidden, "forbidden")
	}
	return c.JSON(http.StatusOK, "ok")
}
