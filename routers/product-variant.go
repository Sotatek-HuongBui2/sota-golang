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

// GetProductsVariant      godoc
// @Summary                Get all product variant
// @Description            Return list of product variant
// @Param                  query query requests.GetProductVariants true "Get products variant"
// @Param                  id path string true "get product variant by id"
// @Produce                application/json
// @Tags                   Product Variant
// @Success                200 {object} utils.IPagination[[]models.Products]
// @Router                 /products/{id}/variant [get]
// @Security               BearerAuth
func GetProductVariants(ctx *gin.Context) {
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

	parentId := ctx.Param("id")

	params := &requests.GetProductVariants{
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
	products, err := services.GetProductVariants(parentId, params)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, products)
}

// FindByIdProductVariant  godoc
// @Summary                Get Single product variant by id.
// @Description            Return the product variant whoes productId value mathes id.
// @Param                  id path string true "get product variant by id"
// @Param                  variant_id path string true "get product variant by id"
// @Produce                application/json
// @Tags                   Product Variant
// @Success                200 {object} models.Products
// @Router                 /products/{id}/variant/{variant_id} [get]
// @Security               BearerAuth
func GetProductVariantById(ctx *gin.Context) {
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

	parentId := ctx.Param("id")
	id := ctx.Param("variant_id")

	fmt.Println(id)

	product, err := services.GetProductVariantById(parentId, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// CreateProductVariant   godoc
// @Summary               Create Product Variant
// @Description           Save product variant data in Db.
// @Param                 body body requests.CreateOrUpdateProductVariant true "Create product variant"
// @Param                 id path string true "get product variant by id"
// @Produce               application/json
// @Tags                  Product Variant
// @Success               200 {object} models.Products
// @Router                /products/{id}/variant [post]
// @Security              BearerAuth
func CreateProductVariant(ctx *gin.Context) {
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

	parentId := ctx.Param("id")
	payload := &requests.CreateOrUpdateProduct{}
	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	product, err := services.CreateProductVariant(parentId, payload)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, product)
}

// UpdateProductVariant    godoc
// @Summary                Update product variant
// @Description            Update product variant data.
// @Param                  id path string true "parent id"
// @Param                  body body requests.CreateOrUpdateProductVariant true "Update product variant"
// @Param                  variant_id path string true "update product variant by variant_id"
// @Tags                   Product Variant
// @Produce                application/json
// @Success                200 {object} models.Products
// @Router                 /products/{id}/variant/{variant_id} [put]
// @Security               BearerAuth
func UpdateProductVariant(ctx *gin.Context) {
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

	payload := &requests.CreateOrUpdateProductVariant{}
	parentId := ctx.Param("id")
	id := ctx.Param("variant_id")

	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	product, httpStatus, err := services.UpdateProductVariant(parentId, id, payload)

	if err != nil {
		ctx.JSON(httpStatus, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}

// DeleteProductVariant   godoc
// @Summary               Get Single product variant by id.
// @Param                 id path string true "parent id"
// @Param                 variant_id path string true "variant id"
// @Description           Return the product variant whoes productId valu mathes id.
// @Produce               application/json
// @Tags                  Product Variant
// @Success               200 {object} models.Products
// @Router                /products/{id}/variant/{variant_id} [delete]
// @Security              BearerAuth
func DeleteProductVariant(ctx *gin.Context) {
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

	id := ctx.Param("variant_id")
	parentId := ctx.Param("id")

	product, err := services.DeleteProductVariant(parentId, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, product)
}
