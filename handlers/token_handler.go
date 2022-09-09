package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"passwords_api/db"
)

type RequestToken struct {
	TokenValue string `json:"token"`
}

func TokenPost(c echo.Context) error {
	coll := db.MongoClient().Database("mongo_app").Collection("test_tokens")
	requestToken := new(RequestToken)
	if err := c.Bind(requestToken); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	doc := bson.D{{"token", requestToken.TokenValue}}
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusAccepted, result)
}
