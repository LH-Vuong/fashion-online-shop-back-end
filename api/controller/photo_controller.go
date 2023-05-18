package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model/response"
	"online_fashion_shop/api/service"
)

type PhotoController struct {
	PhotoService service.PhotoService
}

// Upload Photo
//
//	@Summary		Upload image(s) to product photo to add photos
//	@Description	Create one if product has no photos yet, or add image(s) to photo list
//	@Tags			Photo
//	@Accept 		multipart/form-data
//	@Produce		json
//	@Param 			product_id 		path 		string 			true 	"ID of the product to add photos to"
//	@Param          photo_files[]   formData    []file			true    "The image(s) to upload"
//	@Success		200				{object}	response.BaseResponse[[]string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/product_photos/upload/{product_id} [post]
func (cl *PhotoController) Upload(c *gin.Context) {

	form, _ := c.MultipartForm()
	productId := c.Param("product_id")
	filesHeaders := form.File["photo_files[]"]
	reader := make([]io.Reader, len(filesHeaders))
	for index, fileHeader := range filesHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
			return
		}
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
			return
		}
		reader[index] = bytes.NewReader(fileBytes)
	}
	photos, err := cl.PhotoService.UploadPhoto(reader, productId)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse[[]string]{
		Data:    photos,
		Message: "list of photo url of product",
		Status:  "success",
	})

}

// DeleteOne Photo
//
//	@Summary		Delete one by proto_url
//	@Description	remove photo from photo store by url
//	@Tags			Photo
//	@Accept			json
//	@Produce		json
//	@Param          photo_url     		query       string   true    "photo_url"
//	@Param          product_id     		query       string   true    "photo_url"
//	@Success		200				{object}	response.BaseResponse[string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/product_photo [delete]
func (cl *PhotoController) DeleteOne(c *gin.Context) {
	query := c.Request.URL.Query()
	photoURL := query.Get("photo_url")
	productId := query.Get("product_id")
	err := cl.PhotoService.DeleteOne(photoURL, productId)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse[string]{
		Data:    photoURL,
		Message: fmt.Sprintf("remove photo'%s'", photoURL),
		Status:  "success",
	})
}

// DeleteAll Photo
//
//	@Summary		Delete all photos by product_id
//	@Description	remove all photos which belong product
//	@Tags			Photo
//	@Accept			json
//	@Produce		json
//	@Param          product_id     		path       string   true    "product_id"
//	@Success		200				{object}	response.BaseResponse[string]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/product_photos/{product_id} [delete]
func (cl *PhotoController) DeleteAll(c *gin.Context) {
	productId := c.Param("product_id")
	err := cl.PhotoService.DeletePhotoByProductId(productId)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse[string]{
		Data:    productId,
		Message: "remove all photo of product",
		Status:  "success",
	})

}
