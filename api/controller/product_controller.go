package controller

import (
	"errors"
	"math"
	"net/http"
	"online_fashion_shop/api/common/errs"
	"online_fashion_shop/api/model"
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
			return nil, err
		}
		req.Page = page
	} else {
		req.Page = 0
	}
	if pageSizeStr := queryValues.Get("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)

		if err != nil {
			return nil, err
		}

		if pageSize > request.PageMaximum {
			return nil, errors.New("reach limit, 10_000 record per request")
		}

		if pageSize < 1 {
			pageSize = request.PageMinimum
		}

		req.PageSize = pageSize
	} else {
		req.PageSize = request.PageMinimum
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
//	@Success		200				{object}	model.Product
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/product/{id} [get]
func (cl ProductController) Get(c *gin.Context) {
	productId := c.Param("id")
	product, err := cl.Service.Get(productId)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
	}
	c.JSON(200, gin.H{"status": "success", "data": product})
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
//	@Param          page	  	query       int	   	  false    "current page's number"
//	@Param          page_size	query       int		  false    "Length per page from '1' to '10000'"
//	@Success		200				{object}	response.PagingResponse[model.Product]
//	@Failure		400				{object}	string
//	@Failure		401				{object}	string
//	@Router			/products/ [get]
func (cl ProductController) List(c *gin.Context) {

	req, err := parseListProductsRequest(c)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}
	products, totalPages, err := cl.Service.List(*req)
	if err != nil {
		errs.HandleFailStatus(c, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(200, response.PagingResponse[*model.Product]{
		TotalPage: totalPages,
		Page:      req.Page,
		Status:    "success",
		Data:      products,
	})
}
