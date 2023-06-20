package service

import (
	"log"
	"net/http"
	"online_fashion_shop/api/common/errs"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/api/repository"
	"online_fashion_shop/api/utils"
	"online_fashion_shop/initializers"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	SignUp(*gin.Context, model.SignUpModel)
	VerifyAccount(*gin.Context, string)
	SignIn(*gin.Context, model.SignInModel)
	ResendVerifyEmail(*gin.Context, string)
	ChangeUserPassword(*gin.Context, model.ChangePasswordModel)

	CreateUserAddress(*gin.Context, model.CreateUserAddressModel)
	DeleteUserAddress(*gin.Context, string)
	UpdateUserAddress(*gin.Context, model.UpdateUserAddressModel)
	UpdateUserInfo(*gin.Context, model.UpdateUserInfoModel)
	GetUserAddressList(*gin.Context, model.GetUserAddressListModel)

	AddUserWishlistItem(*gin.Context, model.AddWishListModel)
	DeleteUserWishlistItems(*gin.Context, model.DeleteWishListModel)
	GetUserWishlist(*gin.Context, model.GetUserWishlistModel)
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
	existEmail, err := service.userRepo.GetUserByEmail(ctx, payload.Email)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateUser")
		return
	}

	if existEmail != nil {
		errs.HandleFailStatus(ctx, "User with that email has already existed!", http.StatusBadRequest)
		return
	}

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

	// Check if the file exists.
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateUser")

		return
	} else {
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

	if existVerify == nil {
		errs.HandleFailStatus(ctx, "Token is expired or not exist!", http.StatusBadRequest)

		return
	}

	if time.Until(existVerify.CreatedAt) > 6 {
		errs.HandleFailStatus(ctx, "Token expired", http.StatusBadRequest)
		return
	}

	existUser, err := service.userRepo.GetUserById(ctx, existVerify.UserId)

	if err != nil {
		errs.HandleFailStatus(ctx, "User's not found!", http.StatusNotFound)

		return
	}

	existUser.Verified = true
	existUser.UpdatedAt = time.Now()

	_, err = service.userRepo.UpdateUser(ctx, existUser)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "VerifyAccount")
		return
	}

	if err = service.userRepo.DeleteUserVerify(ctx, existVerify); err != nil {
		errs.HandleErrorStatus(ctx, err, "VerifyAccount")
		return
	}

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

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn*time.Minute, existUser.Id, config.AccessTokenPrivateKey)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateToken")
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn*time.Minute, existUser.Id, config.RefreshTokenPrivateKey)
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateToken")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success",
		"data": model.LoginResponseModel{
			AccessToken:  access_token,
			RefreshToken: refresh_token,
			User:         *existUser,
		},
	})
}

func (service *UserServiceImpl) ResendVerifyEmail(ctx *gin.Context, email string) {
	existUser, err := service.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "ResendVerifyEmail")
		return
	}

	if existUser == nil {
		errs.HandleFailStatus(ctx, "User with that email does not exist!", http.StatusBadRequest)
		return
	}

	if existUser.Verified {
		errs.HandleFailStatus(ctx, "User with that email is already verified!", http.StatusBadRequest)
		return
	}

	code := randstr.String(20)

	verificationCode := utils.Encode(code)
	newUserVerify := model.UserVerify{
		UniqueToken: verificationCode,
		UserId:      existUser.Id,
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
		Fullname: existUser.Fullname,
	}

	filePath := "api/templates/account-verification.html"

	// Check if the file exists.
	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateUser")
		return
	} else {
		utils.SendVerifyEmail(existUser, &emailData, filePath)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "New email has been seed to your email!"})
}

func (service *UserServiceImpl) ChangeUserPassword(ctx *gin.Context, payload model.ChangePasswordModel) {
	if payload.Password != payload.PasswordConfirm {
		errs.HandleFailStatus(ctx, "Password do not match!", http.StatusBadRequest)
		return
	}

	currentUser := ctx.MustGet("currentUser").(model.User)

	existUser, err := service.userRepo.GetUserById(ctx, currentUser.Id)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "ChangePassword")
		return
	}

	if existUser == nil {
		errs.HandleFailStatus(ctx, "User does not exist!", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "ChangePassword")
	}

	existUser.Password = hashedPassword
	existUser.UpdatedAt = time.Now()

	existUser, err = service.userRepo.UpdateUser(ctx, existUser)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "ChangePassword")
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": existUser})
}

func (service *UserServiceImpl) AddUserWishlistItem(ctx *gin.Context, payload model.AddWishListModel) {
	//Todo: check product exist
	currentUser := ctx.MustGet("currentUser").(model.User)
	existWishlistItem, err := service.userRepo.GetUserWishlistItemByProductId(ctx, payload.ProductId)

	if err != nil {
		if err != mongo.ErrNoDocuments {
			errs.HandleErrorStatus(ctx, err, "AddUserWishlistItem")
			return
		}
	}

	if existWishlistItem != nil {
		errs.HandleFailStatus(ctx, "Wishlist item has already existed!", http.StatusBadRequest)
		return
	}

	newWishlistItem := model.UserWishlist{
		ProductId:    payload.ProductId,
		UserId:       currentUser.Id,
		ProductImage: payload.ProductImage,
		Status:       "active",
		CreatedAt:    time.Now(),
	}

	rs, err := service.userRepo.CreateWishlistItem(ctx, &newWishlistItem)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "AddUserWishlistItem")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": rs})
}

func (service *UserServiceImpl) DeleteUserWishlistItems(ctx *gin.Context, payload model.DeleteWishListModel) {
	var deletedWishlistItem []*model.UserWishlist
	var rs []string
	for _, val := range payload.DeleteIds {
		existWishlistItem, err := service.userRepo.GetUserWishlistItemById(ctx, val)

		if err != nil {
			errs.HandleErrorStatus(ctx, err, "DeleteUserWishlistItems")
			return
		}

		if existWishlistItem == nil {
			errs.HandleFailStatus(ctx, "Wishlist item is not existed!", http.StatusBadRequest)
		}

		deletedWishlistItem = append(deletedWishlistItem, existWishlistItem)
	}

	for _, val := range deletedWishlistItem {
		deletedId, err := service.userRepo.DeleteWishlistItem(ctx, val)

		if err != nil {
			errs.HandleErrorStatus(ctx, err, "DeleteUserWishlistItems")
			return
		}
		rs = append(rs, deletedId)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": rs})
}

func (service *UserServiceImpl) GetUserWishlist(ctx *gin.Context, payload model.GetUserWishlistModel) {
	currentUser := ctx.MustGet("currentUser").(model.User)

	rs, total, err := service.userRepo.GetUserWishlist(ctx, currentUser.Id, payload.Page, payload.PageSize)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetUserWishlist")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": rs, "total": total})
}

func (service *UserServiceImpl) CreateUserAddress(ctx *gin.Context, payload model.CreateUserAddressModel) {
	currentUser := ctx.MustGet("currentUser").(model.User)

	newAddress := model.UserAddress{
		UserId:     currentUser.Id,
		ProvinceId: payload.ProvinceId,
		DistrictId: payload.DistrictId,
		WardCode:   payload.WardCode,
		Address:    payload.Address,
		Name:       payload.Name,
		IsDefault:  payload.IsDefault,
		CreatedAt:  time.Now(),
	}

	count, err := service.userRepo.GetUserAddressCount(ctx, currentUser.Id)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateUserAddress")
		return
	}

	rs, err := service.userRepo.CreateUserAddress(ctx, &newAddress)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateUserAddress")
		return
	}

	if count == 0 || newAddress.IsDefault {
		if err := service.userRepo.SetDefaultAddress(ctx, currentUser.Id, rs.Id); err != nil {
			errs.HandleErrorStatus(ctx, err, "CreateUserAddress")
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": rs})
}

func (service *UserServiceImpl) UpdateUserInfo(ctx *gin.Context, payload model.UpdateUserInfoModel) {
	currentUser := ctx.MustGet("currentUser").(model.User)

	existUser, err := service.userRepo.GetUserById(ctx, currentUser.Id)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "UpdateUserInfo")
		return
	}

	if existUser == nil {
		errs.HandleFailStatus(ctx, "User does not exist!", http.StatusBadRequest)
		return
	}

	existUser.Fullname = payload.Fullname
	existUser.PhoneNumber = payload.PhoneNumber
	existUser.Photo = payload.Photo
	existUser.UpdatedAt = time.Now()

	rs, err := service.userRepo.UpdateUser(ctx, existUser)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "UpdateUserInfo")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": rs})
}

func (service *UserServiceImpl) UpdateUserAddress(ctx *gin.Context, payload model.UpdateUserAddressModel) {
	existAddress, err := service.userRepo.GetUserAddressById(ctx, payload.Id)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "UpdateUserAddress")
		return
	}

	if existAddress == nil {
		errs.HandleFailStatus(ctx, "User's Address is not existed!", http.StatusNotFound)
		return
	}

	existAddress.Address = payload.Address
	existAddress.ProvinceId = payload.ProvinceId
	existAddress.DistrictId = payload.DistrictId
	existAddress.WardCode = payload.WardCode
	existAddress.Name = payload.Name
	existAddress.IsDefault = payload.IsDefault
	existAddress.UpdatedAt = time.Now()

	rs, err := service.userRepo.UpdateUserAddress(ctx, existAddress)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateUserAddress")
		return
	}

	count, err := service.userRepo.GetUserAddressCount(ctx, existAddress.UserId)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "CreateUserAddress")
		return
	}

	if count == 0 || existAddress.IsDefault {
		if err := service.userRepo.SetDefaultAddress(ctx, existAddress.UserId, existAddress.Id); err != nil {
			errs.HandleErrorStatus(ctx, err, "CreateUserAddress")
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": rs})
}

func (service *UserServiceImpl) DeleteUserAddress(ctx *gin.Context, addressId string) {
	existAddress, err := service.userRepo.GetUserAddressById(ctx, addressId)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "DeleteUserAddress")
		return
	}

	if existAddress == nil {
		errs.HandleFailStatus(ctx, "User's address is not existed!", http.StatusNotFound)
		return
	}

	rs, err := service.userRepo.DeleteUserAddress(ctx, existAddress)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "DeleteUserAddress")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": rs})
}

func (service *UserServiceImpl) GetUserAddressList(ctx *gin.Context, payload model.GetUserAddressListModel) {
	currentUser := ctx.MustGet("currentUser").(model.User)

	rs, total, err := service.userRepo.GetUserAddressList(ctx, currentUser.Id, payload.Page, payload.PageSize)

	if err != nil {
		errs.HandleErrorStatus(ctx, err, "GetUserAddressList")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": rs, "total": total})
}
