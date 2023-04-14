package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"online_fashion_shop/api/model/request"
	"online_fashion_shop/api/service"
	"strconv"
	"strings"
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

	req.MaxPrice = 1000000000000
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

func (cl ProductController) Get(c *gin.Context) {
	productId := c.Param("id")
	product := cl.Service.Get(productId)
	c.JSON(200, product)
}

func (cl ProductController) List(c *gin.Context) {

	req, err := parseListProductsRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	products := cl.Service.List(*req)
	c.JSON(200, products)

}
