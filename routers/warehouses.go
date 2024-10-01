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
// @Summary         Get all warehouse
// @Description     Return list of warehouse
// @Param           query query requests.GetWarehouses true "Get warehouses"
// @Produce         application/json
// @Tags            Warehouses
// @Success         200 {object} utils.IPagination[[]models.Warehouses]
// @Router          /warehouses [get]
// @Security        BearerAuth
func GetWarehouses(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_WAREHOUSES, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &requests.GetWarehouses{
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

	warehouses := &utils.IPagination[[]models.Warehouses]{}
	warehouses, err := services.GetWarehouses(params)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, warehouses)
}

// FindByIdWarehouse  godoc
// @Summary         Get Single warehouse by id.
// @Description     Return the warehouse whoes WarehouseId value mathes id.
// @Param           id path string true "get warehouse by id"
// @Produce         application/json
// @Tags            Warehouses
// @Success         200 {object} models.Warehouses
// @Router          /warehouses/{id} [get]
// @Security        BearerAuth
func GetWarehouseById(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_WAREHOUSE_BY_ID, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	fmt.Println(id)

	warehouse, err := services.GetWarehouseById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, warehouse)
}

// CreateWarehouses    godoc
// @Summary            Create Warehouse
// @Description        Save warehouse data in Db.
// @Param              body body requests.CreateOrUpdateWarehouse true "Create warehouse"
// @Produce            application/json
// @Tags               Warehouses
// @Success            200 {object} models.Warehouses
// @Router             /warehouses [post]
// @Security           BearerAuth
func CreateWarehouse(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.CREATE_WAREHOUSE, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateWarehouse{}
	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	warehouse, err := services.CreateWarehouse(payload)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, warehouse)
}

// UpdateWarehouses    godoc
// @Summary            Update warehouse
// @Description        Update warehouse data.
// @Param              id path string true "update warehouse by id"
// @Param              body body requests.CreateOrUpdateWarehouse true "Update warehouse"
// @Tags               Warehouses
// @Produce            application/json
// @Success            200 {object} models.Warehouses
// @Router             /warehouses/{id} [put]
// @Security           BearerAuth
func UpdateWarehouse(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.UPDATE_WAREHOUSE, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateWarehouse{}
	id := ctx.Param("id")

	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	warehouse, httpStatus, err := services.UpdateWarehouse(id, payload)

	if err != nil {
		ctx.JSON(httpStatus, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, warehouse)
}

// DeleteWarehouses     godoc
// @Summary             Get Single warehouse by id.
// @Param               id path string true "delete warehouse"
// @Description         Return the warehouse whoes warehouseId valu mathes id.
// @Produce             application/json
// @Tags                Warehouses
// @Success             200 {object} models.Warehouses
// @Router              /warehouses/{id} [delete]
// @Security            BearerAuth
func DeleteWarehouse(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.DELETE_WAREHOUSE, actionerId.(string)); err != nil {
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
