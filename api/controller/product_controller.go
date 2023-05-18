package controller

import (
	"errors"
	"math"
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model/product"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/model/response"
	"online_fashion_shop/api/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	Service service.ProductService
}

func parseListProductsRequest(c *gin.Context) (*request.ListProductsRequest, error) {
	var req request.ListProductsRequest

	// Parse the query parameters
	queryValues := c.Request.URL.Query()

	// Parse the brands' parameter (comma-separated list)
	if brandsStr := queryValues.Get("brands"); brandsStr != "" {
		req.Brands = strings.Split(brandsStr, ",")
	}

	// Parse the colors' parameter (comma-separated list)
	if colorsStr := queryValues.Get("colors"); colorsStr != "" {
		req.Colors = strings.Split(colorsStr, ",")
	}

	// Parse the rate parameter
	if rateStr := queryValues.Get("rate"); rateStr != "" {
		rate, err := strconv.Atoi(rateStr)
		if err != nil {
			return nil, err
		}
		req.Rate = rate
	}

	// Parse the tags' parameter (comma-separated list)
	if tagsStr := queryValues.Get("tags"); tagsStr != "" {
		req.Tags = strings.Split(tagsStr, ",")
	}

	// Parse the genders' parameter (comma-separated list)
	if genderStr := queryValues.Get("genders"); genderStr != "" {
		req.Genders = strings.Split(genderStr, ",")
	}

	// Parse the type parameter (comma-separated list)
	if typeStr := queryValues.Get("types"); typeStr != "" {
		req.Types = strings.Split(typeStr, ",")
	}

	// Parse the price parameters

	req.MaxPrice = math.MaxInt64
	req.MinPrice = 0

	if priceStr := queryValues.Get("price"); priceStr != "" {
		priceRange := strings.Split(priceStr, ",")
		if len(priceRange) == 2 {
			minPrice, err := strconv.ParseInt(priceRange[0], 10, 64)
			if err == nil {
				req.MinPrice = minPrice
			}
			maxPrice, err := strconv.ParseInt(priceRange[1], 10, 64)
			if err == nil {
				req.MaxPrice = maxPrice
			}

		}
	}

	// Parse the name parameter
	req.KeyWord = queryValues.Get("name")

	// Parse the page and page_size parameters
	if pageStr := queryValues.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}
		req.Page = page - 1
	} else {
		req.Page = 0
	}
	if pageSizeStr := queryValues.Get("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)

		if err != nil {
			pageSize = 10
		}

		if pageSize > request.PageMaximum {
			return nil, errors.New("reach limit, 10_000 record per request")
		}

		if pageSize < 1 {
			pageSize = request.PageMinimum
		}

		req.PageSize = pageSize
	} else {
		req.PageSize = 10
	}

	return &req, nil
}

// Get Product
//
//	@Summary		get product info
//	@Description	get the product info
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Param          id     	path       string    true    "product's id"
//	@Success		200				{object}	response.BaseResponse[product.Product]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/product/{id} [get]
func (cl ProductController) Get(c *gin.Context) {
	productId := c.Param("id")
	productDetail, err := cl.Service.Get(productId)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
	}
	c.JSON(200, response.BaseResponse[product.Product]{
		Data:    *productDetail,
		Message: "",
		Status:  "success",
	})
}

// List Product
//
//	@Summary		list of product
//	@Description	list of the product
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Param          brands    	query       string    false    "a list of brand name separated by commas"
//	@Param          colors    	query       string    false    "a list of color name separated by commas (FULL UPPERCASE format)"
//	@Param          tags      	query       string    false    "a list of tag name ['HOT','NEW','SALE'] separated by commas"
//	@Param          genders   	query       string    false    "a list of gender type ['KID','WOMEN','MEN'] separated by commas"
//	@Param          types     	query       string    false    "a list of type name separated by commas"
//	@Param          rate	  	query       int    	  false    "Minimum of avg rate of product"
//	@Param          price	  	query       string	  false    "Range of values in format 'min_value,max_value' "
//	@Param          name	  	query       string	  false    "Key work relate to products' name "
//	@Param          page	  	query       int	   	  false    "current page's number ,start at 1"
//	@Param          page_size	query       int		  false    "Length per page from '1' to '10000'"
//	@Success		200				{object}	response.PagingResponse[product.Product]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/products/ [get]
func (cl ProductController) List(c *gin.Context) {

	req, err := parseListProductsRequest(c)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	products, total, err := cl.Service.List(*req)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(200, response.PagingResponse[*product.Product]{
		Total:  total,
		Length: len(products),
		Status: "success",
		Data:   products,
	})
}

// Create
//
//	@Summary		create a new product
//	@Description	create a new product and return created result
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Param          request     		body       product.Product   true    "info of created product"
//	@Success		200				{object}	response.BaseResponse[product.Product]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/product [put]
func (cl ProductController) Create(c *gin.Context) {
	var createdProduct product.Product
	c.BindJSON(&createdProduct)
	err := cl.Service.Create(&createdProduct)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse[*product.Product]{
		Data:    &createdProduct,
		Message: "created one",
		Status:  "success",
	})
}

// Update
//
//	@Summary		update info of product
//	@Description	provide fields will be updated, fields don't provide will be omitted , and updated result
//	@Tags			product
//	@Accept			json
//	@Produce		json
//	@Param          request     		body       product.Product   true    "fields want to update with provided value"
//	@Param          request     		body       product.Product   true    "fields want to update with provided value"
//	@Success		200				{object}	response.BaseResponse[product.Product]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/product [post]
func (cl ProductController) Update(c *gin.Context) {
	var updateProduct product.Product
	err := c.BindJSON(&updateProduct)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	id := c.Param("product_id")
	updateProduct.Id = id
	err = cl.Service.Update(&updateProduct)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, response.BaseResponse[*product.Product]{
		Data:    &updateProduct,
		Message: "updated one",
		Status:  "success",
	})

}
