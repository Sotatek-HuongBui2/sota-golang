package routers

import (
	"errors"
	"fmt"
	"net/http"
	"vtcanteen/constants"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/services"

	errorConstants "vtcanteen/constants/errors"
	"vtcanteen/utils"

	"github.com/gin-gonic/gin"
)

// GetProducts      godoc
// @Summary         Get all product
// @Description     Return list of product
// @Param           query query requests.GetProducts true "Get products"
// @Produce         application/json
// @Tags            Product
// @Success         200 {object} utils.IPagination[[]models.Products]
// @Router          /products [get]
// @Security        BearerAuth
func GetProducts(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_PRODUCTS, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &requests.GetProducts{
		GetList: requests.GetList{
			Limit:        ctx.Query("limit"),
			Sort:         ctx.Query("sort"),
			SortDir:      ctx.Query("sort_dir"),
			SearchFields: ctx.Query("search_fields"),
			Search:       ctx.Query("search"),
			Page:         ctx.Query("page"),
			Filter:       ctx.Query("filter"),
		},
	}

	products := &utils.IPagination[[]models.Products]{}
	products, err := services.GetProducts(params)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

// FindByIdProduct  godoc
// @Summary         Get Single product by id.
// @Description     Return the product whoes productId value mathes id.
// @Param           id path string true "get product by id"
// @Produce         application/json
// @Tags            Product
// @Success         200 {object} models.Products
// @Router          /products/{id} [get]
// @Security        BearerAuth
func GetProductById(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_PRODUCT_BY_ID, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	fmt.Println(id)

	product, err := services.GetProductById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// CreateProduct    godoc
// @Summary         Create Product
// @Description     Save product data in Db.
// @Param           body body requests.CreateOrUpdateProduct true "Create product"
// @Produce         application/json
// @Tags            Product
// @Success         200 {object} models.Products
// @Router          /products [post]
// @Security        BearerAuth
func CreateProduct(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.CREATE_PRODUCT, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateProduct{}
	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	product, err := services.CreateProduct(payload)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, product)
}

// UpdateProduct   godoc
// @Summary         Update product
// @Description     Update product data.
// @Param           id path string true "update product by id"
// @Param           body body requests.CreateOrUpdateProduct true "Update product"
// @Tags            Product
// @Produce         application/json
// @Success         200 {object} models.Products
// @Router          /products/{id} [put]
// @Security        BearerAuth
func UpdateProduct(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.UPDATE_PRODUCT, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateProduct{}
	id := ctx.Param("id")

	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	product, httpStatus, err := services.UpdateProduct(id, payload)

	if err != nil {
		ctx.JSON(httpStatus, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// DeleteProduct   godoc
// @Summary         Get Single product by id.
// @Param           id path string true "delete product"
// @Description     Return the product whoes productId valu mathes id.
// @Produce         application/json
// @Tags            Product
// @Success         200 {object} models.Products
// @Router          /products/{id} [delete]
// @Security        BearerAuth
func DeleteProduct(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.DELETE_PRODUCT, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	fmt.Println(id)
	product, err := services.DeleteProduct(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}
