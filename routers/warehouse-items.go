package routers

import (
	"errors"
	"fmt"
	"net/http"
	"vtcanteen/constants"
	errorConstants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/services"
	"vtcanteen/utils"

	"github.com/gin-gonic/gin"
)

// GetWarehouse     godoc
// @Summary         Get all warehouse item
// @Description     Return list of warehouse item
// @Param           query query requests.GetWarehouseItems true "Get WarehouseItems"
// @Produce         application/json
// @Tags            WarehouseItems
// @Success         200 {object} utils.IPagination[[]models.WarehouseItems]
// @Router          /warehouse-items [get]
// @Security        BearerAuth
func GetWarehouseItems(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_WAREHOUSE_ITEMS, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &requests.GetWarehouseItems{
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

	warehouseItems := &utils.IPagination[[]models.WarehouseItems]{}

	warehouseItems, err := services.GetWarehouseItems(params)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, warehouseItems)
}

// FindByIdWarehouse  godoc
// @Summary         Get Single warehouse item by id.
// @Description     Return the warehouse item whoes WarehouseItemId value mathes id.
// @Param           id path string true "get warehouse item by id"
// @Produce         application/json
// @Tags            WarehouseItems
// @Success         200 {object} models.WarehouseItems
// @Router          /warehouse-items/{id} [get]
// @Security        BearerAuth
func GetWarehouseItemById(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_WAREHOUSE_ITEM_BY_ID, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	fmt.Println(id)

	warehouse, err := services.GetWarehouseItemById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, warehouse)
}

// CreateWarehouseItem    godoc
// @Summary               Create Warehouse item
// @Description           Save warehouse item data in Db.
// @Param                 body body models.WarehouseItems true "Create warehouse item"
// @Produce               application/json
// @Tags                  WarehouseItems
// @Success               200 {object} models.WarehouseItems
// @Router                /warehouse-items [post]
// @Security              BearerAuth
func CreateWarehouseItem(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.CREATE_WAREHOUSE_ITEM, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &models.WarehouseItems{}
	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	warehouse, err := services.CreateWarehouseItem(payload)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, warehouse)
}

// UpdateWarehouseItem    godoc
// @Summary               Update warehouse item
// @Description           Update warehouse item data.
// @Param                 id path string true "update warehouse item by id"
// @Param                 body body models.WarehouseItems true "Update warehouse"
// @Tags                  WarehouseItems
// @Produce               application/json
// @Success               200 {object} models.WarehouseItems
// @Router                /warehouse-items/{id} [put]
// @Security              BearerAuth
func UpdateWarehouseItem(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.UPDATE_WAREHOUSE_ITEM, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &models.WarehouseItems{}
	id := ctx.Param("id")

	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	warehouse, err := services.UpdateWarehouseItem(id, payload)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, warehouse)
}

// DeleteWarehouseItem  godoc
// @Summary             Get Single warehouse item by id.
// @Param               id path string true "delete warehouse item"
// @Description         Return the warehouse item whoes warehouseItemId value mathes id.
// @Produce             application/json
// @Tags                WarehouseItems
// @Success             200 {object} models.WarehouseItems
// @Router              /warehouse-items/{id} [delete]
// @Security            BearerAuth
func DeleteWarehouseItem(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.DELETE_WAREHOUSE_ITEM, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	fmt.Println(id)
	warehouse, err := services.DeleteWarehouse(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, warehouse)
}
