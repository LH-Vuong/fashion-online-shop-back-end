package middleware

import (
	"context"
	"fmt"
	"net/http"
	"online_fashion_shop/api/common/errs"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/api/utils"
	"online_fashion_shop/initializers"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if err == nil {
			access_token = cookie
		}

		if access_token == "" {
			errs.HandleFailStatus(ctx, "You are not logged in", http.StatusUnauthorized)
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			errs.HandleFailStatus(ctx, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := getUserById(ctx, fmt.Sprint(sub), config)

		if err != nil {
			errs.HandleFailStatus(ctx, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx.Set("currentUser", *user)
		ctx.Next()
	}
}

func getUserById(ctx context.Context, userId string, config initializers.Config) (user *model.User, err error) {
	var cl initializers.Client

	cl, err = initializers.NewClient(config.MongoUrl)

	if err != nil {
		return nil, err
	}

	userCollection := cl.Database("fashion_shop").Collection("user")

	objId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return nil, err
	}

	ct, cancel := initializers.InitContext()

	defer cancel()

	rs := userCollection.FindOne(ct, bson.M{"_id": objId})

	if err != nil {
		return nil, err
	}

	err = rs.Decode(&user)

	return
}
