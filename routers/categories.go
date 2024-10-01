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

// GetCategories    godoc
// @Summary         Get all category
// @Description     Return list of category
// @Param           query query requests.GetCategories true "Get categories"
// @Produce         application/json
// @Tags            Category
// @Success         200 {object} utils.IPagination[[]models.Categories]
// @Router          /categories [get]
// @Security        BearerAuth
func GetCategories(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_HISTORIES, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &requests.GetCategories{
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

	categories := &utils.IPagination[[]models.Categories]{}
	categories, err := services.GetCategories(params)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, categories)
}

// FindByIdCategory godoc
// @Summary         Get Single category by id.
// @Description     Return the category whoes categoryId value mathes id.
// @Param           id path string true "get category by id"
// @Produce         application/json
// @Tags            Category
// @Success         200 {object} models.Categories
// @Router          /categories/{id} [get]
// @Security        BearerAuth
func GetCategoryById(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_CATEGORY_BY_ID, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	fmt.Println(id)

	category, err := services.GetCategoryById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

// CreateCategory   godoc
// @Summary         Create Category
// @Description     Save category data in Db.
// @Param           body body requests.CreateOrUpdateCategory true "Create category"
// @Produce         application/json
// @Tags            Category
// @Success         200 {object} models.Categories
// @Router          /categories [post]
// @Security        BearerAuth
func CreateCategory(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.CREATE_CATEGORY, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateCategory{}
	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	category, err := services.CreateCategory(payload)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, category)
}

// UpdateCategory   godoc
// @Summary         Update category
// @Description     Update category data.
// @Param           id path string true "update category by id"
// @Param           body body requests.CreateOrUpdateCategory true "Update category"
// @Tags            Category
// @Produce         application/json
// @Success         200 {object} models.Categories
// @Router          /categories/{id} [put]
// @Security        BearerAuth
func UpdateCategory(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.UPDATE_CATEGORY, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateCategory{}
	id := ctx.Param("id")

	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	category, httpStatus, err := services.UpdateCategory(id, payload)

	if err != nil {
		ctx.JSON(httpStatus, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, category)
}

// DeleteCategory   godoc
// @Summary         Get Single category by id.
// @Param           id path string true "delete category"
// @Description     Return the category whoes categoryId valu mathes id.
// @Produce         application/json
// @Tags            Category
// @Success         200 {object} models.Categories
// @Router          /categories/{id} [delete]
// @Security        BearerAuth
func DeleteCategory(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.DELETE_CATEGORY, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	fmt.Println(id)
	category, err := services.DeleteCategory(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, category)
}
