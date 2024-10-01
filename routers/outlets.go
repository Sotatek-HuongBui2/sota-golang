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

// GetOutlet        godoc
// @Summary         Get all outlet
// @Description     Return list of outlet
// @Param           query query requests.GetOutlets true "Get outlets"
// @Produce         application/json
// @Tags            Outlets
// @Success         200 {object} utils.IPagination[[]models.Outlets]
// @Router          /outlets [get]
// @Security        BearerAuth
func GetOutlets(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_OUTLETS, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &requests.GetOutlets{
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

	outlets := &utils.IPagination[[]models.Outlets]{}
	outlets, err := services.GetOutlets(params)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, outlets)
}

// FindByIdOutlets  godoc
// @Summary         Get Single outlet by id.
// @Description     Return the outlet whoes OutletId value mathes id.
// @Param           id path string true "get outlet by id"
// @Produce         application/json
// @Tags            Outlets
// @Success         200 {object} models.Outlets
// @Router          /outlets/{id} [get]
// @Security        BearerAuth
func GetOutletById(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_OUTLET_BY_ID, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	fmt.Println(id)

	outlet, err := services.GetOutletById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, outlet)
}

// CreateOutlets    godoc
// @Summary         Create Outlet
// @Description     Save outlet data in Db.
// @Param           body body requests.CreateOrUpdateOutlet true "Create outlet"
// @Produce         application/json
// @Tags            Outlets
// @Success         200 {object} models.Outlets
// @Router          /outlets [post]
// @Security        BearerAuth
func CreateOutlet(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.CREATE_OUTLET, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateOutlet{}
	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT)})
		return
	}

	outlet, err := services.CreateOutlet(payload)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, outlet)
}

// UpdateOutlets    godoc
// @Summary         Update outlet
// @Description     Update outlet data.
// @Param           id path string true "update outlet by id"
// @Param           body body requests.CreateOrUpdateOutlet true "Update outlet"
// @Tags            Outlets
// @Produce         application/json
// @Success         200 {object} models.Outlets
// @Router          /outlets/{id} [put]
// @Security        BearerAuth
func UpdateOutlet(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.UPDATE_OUTLET, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateOutlet{}
	id := ctx.Param("id")

	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	outlet, httpStatus, err := services.UpdateOutlet(id, payload)

	if err != nil {
		ctx.JSON(httpStatus, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, outlet)
}

// DeleteOutlets    godoc
// @Summary         Get Single outlet by id.
// @Param           id path string true "delete outlet"
// @Description     Return the outlet whoes outletId valu mathes id.
// @Produce         application/json
// @Tags            Outlets
// @Success         200 {object} models.Outlets
// @Router          /outlets/{id} [delete]
// @Security        BearerAuth
func DeleteOutlet(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.DELETE_OUTLET, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	fmt.Println(id)
	outlet, err := services.DeleteOutlet(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, outlet)
}
