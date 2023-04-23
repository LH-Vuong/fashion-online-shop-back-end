package service

import (
	"fmt"
	"log"
	"net/http"
	"online_fashion_shop/api/common/errs"
	model "online_fashion_shop/api/model/user"
	repository "online_fashion_shop/api/repository/user"
	"online_fashion_shop/api/utils"
	"online_fashion_shop/initializers"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
)

type UserService interface {
	SignUp(*gin.Context, model.SignUpModel)
	VerifyAccount(*gin.Context, string)
	SignIn(*gin.Context, model.SignInModel)
}

func NewUserServiceImpl(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func (service *UserServiceImpl) SignUp(ctx *gin.Context, payload model.SignUpModel) {
	if payload.Password != payload.PasswordConfirm {
		errs.HandleFailStatus(ctx, "Password do not match!", http.StatusBadRequest)
		return
	}

	hashedPasswed, err := utils.HashPassword(payload.Password)

	if err != nil {
		errs.HandleFailStatus(ctx, "Hashed password failed!", http.StatusInternalServerError)
		return
	}

	newUser := model.User{
		Id:        uuid.New(),
		Fullname:  payload.Fullname,
		Password:  hashedPasswed,
		Email:     payload.Email,
		Verified:  false,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	rs, err := service.userRepo.CreateUser(ctx, &newUser)

	if err != nil {
		if strings.Contains(err.Error(), "uplicate key value violates unique") {
			errs.HandleFailStatus(ctx, "User with that email already exists", http.StatusBadRequest)
			return
		} else {
			errs.HandleErrorStatus(ctx, err, "CreateUser")
			return
		}
	}

	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	newUserVerify := model.UserVerify{
		Id:          uuid.New(),
		UniqueToken: verificationCode,
		UserId:      newUser.Id,
		Status:      "active",
		CreatedAt:   time.Now(),
	}

	_, err = service.userRepo.CreateUserVerify(ctx, &newUserVerify)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateUser")
		return
	}

	emailData := utils.VerifyEmailData{
		Token:    verificationCode,
		Fullname: newUser.Fullname,
	}

	filePath := "api/templates/account-verification.html"
	fileInfo, err := os.Stat(filePath)

	// Check if the file exists.
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist")
		} else {
			fmt.Println("Error checking file existence:", err)
		}
	} else {
		fmt.Println("File exists with size", fileInfo.Size(), "bytes")
		utils.SendVerifyEmail(&newUser, &emailData, filePath)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": rs})
}

func (service *UserServiceImpl) VerifyAccount(ctx *gin.Context, uniqueToken string) {
	existVerify, err := service.userRepo.GetVerifyByUniqueToken(ctx, uniqueToken)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "VerifyAccount")

		return
	}

	if existVerify.CreatedAt.Sub(time.Now()).Hours() > 6 {
		errs.HandleFailStatus(ctx, "Token expired", http.StatusBadRequest)

		return
	}

	existUser, err := service.userRepo.GetUserById(ctx, existVerify.UserId.String())

	if err != nil {
		errs.HandleFailStatus(ctx, "User not found!", http.StatusNotFound)
	}

	existUser.Verified = true
	existUser.UpdatedAt = time.Now()

	_, err = service.userRepo.UpdateUser(ctx, existUser)

	err = service.userRepo.DeleteUserVerify(ctx, existVerify)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (service *UserServiceImpl) SignIn(ctx *gin.Context, payload model.SignInModel) {
	existUser, err := service.userRepo.GetUserByEmail(ctx, payload.Email)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "SignIn")

		return
	}

	if err := utils.VerifyPassword(existUser.Password, payload.Password); err != nil {
		errs.HandleFailStatus(ctx, "Email or password is invalid!", http.StatusBadRequest)
		return
	}

	if !existUser.Verified {
		errs.HandleFailStatus(ctx, "User not verified", http.StatusBadRequest)

		return
	}

	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, existUser.Id, config.AccessTokenPrivateKey)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateToken")
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, existUser.Id, config.RefreshTokenPrivateKey)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateToken")
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge, "/", "", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge, "/", "", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge, "/", "", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": existUser})
}
