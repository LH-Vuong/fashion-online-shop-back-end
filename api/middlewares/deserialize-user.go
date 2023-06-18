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
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)
		queryToken := ctx.Param("token")

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		}

		if queryToken != "" {
			access_token = queryToken
		}

		if access_token == "" {
			errs.HandleFailStatus(ctx, "You are not logged in!", http.StatusUnauthorized)
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
			errs.HandleErrorStatus(ctx, err, "VerifyToken")
		}

		ctx.Set("currentUser", *user)
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
	err = rs.Decode(&user)

	if err != nil || user == nil {
		return nil, err
	}

	return
}

func DeserializeAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)
		queryToken := ctx.Param("token")

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		}

		if queryToken != "" {
			access_token = queryToken
		}

		if access_token == "" {
			errs.HandleFailStatus(ctx, "You are not logged in!", http.StatusUnauthorized)
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)

		if err != nil {
			errs.HandleFailStatus(ctx, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := getAdminById(ctx, fmt.Sprint(sub), config)

		if err != nil || user == nil {
			errs.HandleFailStatus(ctx, "You are not log in!", http.StatusUnauthorized)
			return
		}

		ctx.Set("currentUser", *user)
	}
}

func getAdminById(ctx context.Context, userId string, config initializers.Config) (user *model.User, err error) {
	var cl initializers.Client
	err = fmt.Errorf("wrong credentials")
	var userRoleMapping []*model.UserRoleMapping
	cl, err = initializers.NewClient(config.MongoUrl)

	if err != nil {
		return nil, err
	}

	userCollection := cl.Database("fashion_shop").Collection("user")
	userRoleMappingCollection := cl.Database("fashion_shop").Collection("user_role_mapping")
	userRoleCollection := cl.Database("fashion_shop").Collection("user_role")

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

	if err := rs.Decode(&user); err != nil || user == nil {
		return nil, err
	}

	query := bson.M{"user_id": userId}

	rs1, err := userRoleMappingCollection.Find(ctx, query)

	if err != nil {
		return nil, err
	}

	if err = rs1.All(ctx, &userRoleMapping); err != nil || len(userRoleMapping) == 0 {
		return nil, err
	}

	for _, v := range userRoleMapping {
		objId, err := primitive.ObjectIDFromHex(v.RoleId)

		if err != nil {
			return nil, err
		}
		rs1 := userRoleCollection.FindOne(ctx, bson.M{"_id": objId})

		if err != nil {
			return nil, err
		}

		var role model.UserRole

		if err = rs1.Decode(&role); err != nil {
			return nil, err
		}

		if role.Role != "admin" {
			return nil, err
		}
	}

	err = nil
	return
}
