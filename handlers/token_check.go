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

func TokenCheck(c echo.Context) error {
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
		Username: claims["username"].(string),
		Password: claims["password"].(string),
		Database: claims["database"].(string),
	}

	collection := mongo_db.MongoClient().Database(userData.Database).Collection("tokens")

	var result bson.M
	err = collection.FindOne(context.TODO(), bson.D{
		{"username", userData.Username},
	}).Decode(&result)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "bad request - unable  to find")
	}

	if userData.Password != result["password"] {
		return c.JSON(http.StatusForbidden, "forbidden - not valid hash")
	}

	return c.JSON(http.StatusOK, "ok")
}
