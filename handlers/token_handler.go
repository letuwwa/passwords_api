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

const dbName = "mongo_app"
const collectionName = "test_tokens"

func TokenPost(c echo.Context) error {
	coll := db.MongoClient().Database(dbName).Collection(collectionName)
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
