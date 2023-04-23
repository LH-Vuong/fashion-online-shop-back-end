package middleware

import (
	"fmt"
	"net/http"
	"online_fashion_shop/api/common/errs"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/api/utils"
	"online_fashion_shop/initializers"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		fmt.Println(ctx.Cookie("access_token"))
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

		var user model.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			errs.HandleFailStatus(ctx, "The user belonging to this token no logger exists", http.StatusForbidden)
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
