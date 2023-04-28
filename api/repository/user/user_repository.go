package repository

import (
	"context"
	model "online_fashion_shop/api/model/user"
	"online_fashion_shop/initializers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository interface {
	CreateUser(context.Context, *model.User) (*model.User, error)
	CreateUserVerify(context.Context, *model.UserVerify) (*model.UserVerify, error)
	GetUserByEmail(context.Context, string) (*model.User, error)
	GetUserById(context.Context, string) (*model.User, error)
	DeleteUser(context.Context, *model.User) (string, error)
	UpdateUser(context.Context, *model.User) (*model.User, error)

	GetVerifyByUniqueToken(context.Context, string) (*model.UserVerify, error)
	UpdateUserVerify(context.Context, *model.UserVerify) error
	DeleteUserVerify(context.Context, *model.UserVerify) error

	CreateWishlistItem(context.Context, *model.UserWishlist) (*model.UserWishlist, error)
	DeleteWishlistItem(context.Context, *model.UserWishlist) (string, error)
	GetUserWishlist(context.Context, string, int64, int64) (*[]model.UserWishlist, int64, error)
	GetUserWishlistItemByProductId(context.Context, string) (*model.UserWishlist, error)
	GetUserWishlistItemById(context.Context, string) (*model.UserWishlist, error)

	CreateUserAddress(context.Context, *model.UserAddress) (*model.UserAddress, error)
	DeleteUserAddress(context.Context, *model.UserAddress) (string, error)
	GetUserAddressList(context.Context, string, int64, int64) (*[]model.UserAddress, int64, error)
	UpdateUserAddress(context.Context, *model.UserAddress) (*model.UserAddress, error)
	GetUserAddressById(context.Context, string) (*model.UserAddress, error)
}

type userRepository struct {
	userCollection         initializers.Collection
	userVerifyCollection   initializers.Collection
	userWishlistCollection initializers.Collection
	userAddressCollection  initializers.Collection
}

func NewUserRepositoryImpl(userCollection initializers.Collection,
	userVerifyCollection initializers.Collection,
	userWishlistCollection initializers.Collection,
	userAddressCollection initializers.Collection) UserRepository {
	return &userRepository{
		userCollection:         userCollection,
		userVerifyCollection:   userVerifyCollection,
		userWishlistCollection: userWishlistCollection,
		userAddressCollection:  userAddressCollection,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	res, err := r.userCollection.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.Id = res.(primitive.ObjectID).Hex()
	return user, err
}

func (r *userRepository) CreateUserVerify(ctx context.Context, userVerify *model.UserVerify) (*model.UserVerify, error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	res, err := r.userVerifyCollection.InsertOne(ctx, userVerify)

	if err != nil {
		return nil, err
	}

	userVerify.Id = res.(primitive.ObjectID).Hex()
	return userVerify, err
}

func (r *userRepository) GetVerifyByUniqueToken(ctx context.Context, uniqueToken string) (userVerify *model.UserVerify, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	rs := r.userVerifyCollection.FindOne(ctx, bson.M{"unique_token": uniqueToken})

	err = rs.Decode(&userVerify)

	return userVerify, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (user *model.User, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	rs := r.userCollection.FindOne(ctx, bson.M{"email": email})

	err = rs.Decode(&user)

	return user, nil
}

func (r *userRepository) GetUserById(ctx context.Context, id string) (user *model.User, err error) {
	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	ctx, cancel := initializers.InitContext()

	defer cancel()

	rs := r.userCollection.FindOne(ctx, bson.M{"_id": objId})

	err = rs.Decode(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUserVerify(ctx context.Context, userVerify *model.UserVerify) error {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	query := bson.M{"_id": userVerify.Id}
	update := bson.M{"$set": userVerify}
	_, err := r.userVerifyCollection.UpdateOne(ctx, query, update)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteUserVerify(ctx context.Context, userVerify *model.UserVerify) error {
	objId, err := primitive.ObjectIDFromHex(userVerify.Id)

	if err != nil {
		return err
	}

	ctx, cancel := initializers.InitContext()

	defer cancel()

	query := bson.M{"_id": objId}
	_, err = r.userVerifyCollection.DeleteOne(ctx, query)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, user *model.User) (string, error) {
	objId, err := primitive.ObjectIDFromHex(user.Id)

	if err != nil {
		return "", err
	}

	ctx, cancel := initializers.InitContext()

	defer cancel()

	query := bson.M{"_id": objId}
	_, err = r.userVerifyCollection.DeleteOne(ctx, query)

	if err != nil {
		return "", err
	}

	return user.Id, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	objId, err := primitive.ObjectIDFromHex(user.Id)

	if err != nil {
		return nil, err
	}

	ctx, cancel := initializers.InitContext()

	defer cancel()

	updateData, err := bson.Marshal(user)
	if err != nil {
		return nil, err
	}

	updatedDataMap := bson.M{}
	err = bson.Unmarshal(updateData, &updatedDataMap)

	if err != nil {
		return nil, err
	}

	delete(updatedDataMap, "_id")

	query := bson.M{"_id": objId}
	update := bson.M{"$set": updatedDataMap}

	_, err = r.userCollection.UpdateOne(ctx, query, update)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) CreateWishlistItem(ctx context.Context, userWislist *model.UserWishlist) (*model.UserWishlist, error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	res, err := r.userWishlistCollection.InsertOne(ctx, userWislist)

	if err != nil {
		return nil, err
	}

	userWislist.Id = res.(primitive.ObjectID).Hex()
	return userWislist, err
}

func (r *userRepository) DeleteWishlistItem(ctx context.Context, userWislist *model.UserWishlist) (string, error) {
	objId, err := primitive.ObjectIDFromHex(userWislist.Id)

	if err != nil {
		return "", err
	}

	ctx, cancel := initializers.InitContext()

	defer cancel()

	query := bson.M{"_id": objId}
	_, err = r.userVerifyCollection.DeleteOne(ctx, query)

	if err != nil {
		return "", err
	}

	return userWislist.Id, nil
}

func (r *userRepository) GetUserWishlist(ctx context.Context, userId string, page int64, pageSize int64) (wishlistItems *[]model.UserWishlist, total int64, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	opts := options.Find().
		SetSkip(page - 1).
		SetLimit(pageSize)

	query := bson.M{"user_id": userId}

	rs, err := r.userWishlistCollection.Find(ctx, query, opts)

	if err != nil {
		return nil, 0, err
	}

	total, err = r.userWishlistCollection.CountDocuments(ctx, query)

	if err != nil {
		return nil, 0, err
	}

	err = rs.All(ctx, &wishlistItems)

	if err != nil {
		return nil, 0, err
	}
	return
}

func (r *userRepository) GetUserWishlistItemByProductId(ctx context.Context, productId string) (wishlistItem *model.UserWishlist, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	query := bson.M{"product_id": productId}

	rs := r.userWishlistCollection.FindOne(ctx, query)

	if err != nil {
		return nil, err
	}

	if err = rs.Decode(&wishlistItem); err != nil {
		return nil, err
	}

	return wishlistItem, nil
}

func (r *userRepository) GetUserWishlistItemById(ctx context.Context, id string) (wishlistItem *model.UserWishlist, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	query := bson.M{"_id": id}

	rs := r.userWishlistCollection.FindOne(ctx, query)

	if err != nil {
		return nil, err
	}

	if err = rs.Decode(&wishlistItem); err != nil {
		return nil, err
	}

	return wishlistItem, nil
}

func (r *userRepository) CreateUserAddress(ctx context.Context, userAddress *model.UserAddress) (*model.UserAddress, error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	res, err := r.userAddressCollection.InsertOne(ctx, userAddress)

	if err != nil {
		return nil, err
	}

	userAddress.Id = res.(primitive.ObjectID).Hex()
	return userAddress, err
}

func (r *userRepository) DeleteUserAddress(ctx context.Context, userAddress *model.UserAddress) (string, error) {
	objId, err := primitive.ObjectIDFromHex(userAddress.Id)

	if err != nil {
		return "", err
	}

	ctx, cancel := initializers.InitContext()

	defer cancel()

	query := bson.M{"_id": objId}
	_, err = r.userAddressCollection.DeleteOne(ctx, query)

	if err != nil {
		return "", err
	}

	return userAddress.Id, nil
}

func (r *userRepository) GetUserAddressList(ctx context.Context, userId string, page int64, pageSize int64) (userAddressList *[]model.UserAddress, total int64, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	opts := options.Find().
		SetSkip(page - 1).
		SetLimit(pageSize)

	query := bson.M{"user_id": userId}

	rs, err := r.userAddressCollection.Find(ctx, query, opts)

	if err != nil {
		return nil, 0, err
	}

	total, err = r.userCollection.CountDocuments(ctx, query)

	if err != nil {
		return nil, 0, err
	}

	err = rs.All(ctx, &userAddressList)

	if err != nil {
		return nil, 0, err
	}
	return
}

func (r *userRepository) UpdateUserAddress(ctx context.Context, userAddress *model.UserAddress) (*model.UserAddress, error) {
	objId, err := primitive.ObjectIDFromHex(userAddress.Id)

	if err != nil {
		return nil, err
	}

	ctx, cancel := initializers.InitContext()

	defer cancel()

	updateData, err := bson.Marshal(userAddress)
	if err != nil {
		return nil, err
	}

	updatedDataMap := bson.M{}
	err = bson.Unmarshal(updateData, &updatedDataMap)

	if err != nil {
		return nil, err
	}

	delete(updatedDataMap, "_id")

	query := bson.M{"_id": objId}
	update := bson.M{"$set": updatedDataMap}

	_, err = r.userAddressCollection.UpdateOne(ctx, query, update)

	if err != nil {
		return nil, err
	}

	return userAddress, nil
}

func (r *userRepository) GetUserAddressById(ctx context.Context, id string) (userAddress *model.UserAddress, err error) {
	ctx, cancel := initializers.InitContext()

	defer cancel()

	query := bson.M{"_id": id}

	rs := r.userAddressCollection.FindOne(ctx, query)

	if err != nil {
		return nil, err
	}

	if err = rs.Decode(&userAddress); err != nil {
		return nil, err
	}

	return userAddress, nil
}
