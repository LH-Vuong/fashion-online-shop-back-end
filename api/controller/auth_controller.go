package controller

import (
	"online_fashion_shop/api/common/errs"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/api/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Service service.UserService
}

func NewAuthController(service service.UserService) AuthController {
	return AuthController{service}
}

// SignUp User godoc
//	@Summary		sign up user
//	@Description	Endpoint to allow a user to sign up with their details
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			SignUpModel	body		model.SignUpModel	true	"User's profile"
//	@Success		201			{object}	string				"success"
//	@Failure		500			{object}	string				"error"
//	@Failure		400			{object}	string				"fail"
//	@Router			/auth/sign-up [post]
func (c AuthController) SignUp(ctx *gin.Context) {
	var payload model.SignUpModel

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	c.Service.SignUp(ctx, payload)
}

// Verify Email godoc
//	@Summary		Verify user email
//	@Description	Verify user email using verification code
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			uniqueToken			query		int		false	"unique token"
//	@Success		200					{object}	string
//	@Failure		400					{object}	string
//	@Failure		409					{object}	string
//	@Router			/auth/verify [get]
func (c AuthController) VerifyAccount(ctx *gin.Context) {
	uniqueToken := ctx.DefaultQuery("uniqueToken", "")

	c.Service.VerifyAccount(ctx, uniqueToken)
}

// Verify Email godoc
//	@Summary		Login
//	@Description	Login and set access token to header
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			SignInModel		body		model.SignInModel	true	"User's credentials"
//	@Success		200				{object}	string
//	@Failure		400				{object}	string
//	@Failure		409				{object}	string
//	@Router			/auth/sign-in [post]
func (c AuthController) SignIn(ctx *gin.Context) {
	var payload model.SignInModel

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errs.HandleErrorStatus(ctx, err, "ShouldBindJSON")
		return
	}

	c.Service.SignIn(ctx, payload)
}
