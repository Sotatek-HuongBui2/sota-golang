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

// GetRoles         godoc
// @Summary         Get all role
// @Description     Return list of role
// @Param           query query requests.GetRoles true "Create role"
// @Produce         application/json
// @Tags            Roles
// @Success         200 {object} utils.IPagination[[]models.Roles]
// @Router          /roles [get]
// @Security        BearerAuth
func GetRoles(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_ROLES, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &requests.GetRoles{
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
	roles := &utils.IPagination[[]models.Roles]{}
	roles, err := services.GetRoles(params)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, roles)
}

// FindByIdRoles    godoc
// @Summary         Get Single role by id.
// @Description     Return the role whoes roleId valu mathes id.
// @Param           id path string true "update roles by id"
// @Produce         application/json
// @Tags            Roles
// @Success         200 {object} models.Roles
// @Router          /roles/{id} [get]
// @Security        BearerAuth
func GetRoleById(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_ROLE_BY_ID, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")

	fmt.Println(id)

	role, err := services.GetRoleById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, role)
}

// CreateRoles      godoc
// @Summary         Create Role
// @Description     Save role data in Db.
// @Param           body body requests.CreateOrUpdateRole true "Create role"
// @Produce         application/json
// @Tags            Roles
// @Success         200 {object} models.Roles
// @Router          /roles [post]
// @Security        BearerAuth
func CreateRole(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.CREATE_ROLE, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateRole{}
	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	role, err := services.CreateRole(payload)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, role)
}

// UpdateRoles      godoc
// @Summary         Update role
// @Description     Update role data.
// @Param           id path string true "update role by id"
// @Param           body body requests.CreateOrUpdateRole true  "Update role"
// @Tags            Roles
// @Produce         application/json
// @Success         200 {object} models.Roles
// @Router          /roles/{id} [put]
// @Security        BearerAuth
func UpdateRole(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.UPDATE_ROLE, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payload := &requests.CreateOrUpdateRole{}
	id := ctx.Param("id")

	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	role, httpStatus, err := services.UpdateRole(id, payload)

	if err != nil {
		ctx.JSON(httpStatus, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, role)
}

// DeleteRoles      godoc
// @Summary         Get Single role by id.
// @Param           id path string true "delete role"
// @Description     Return the role whoes roleId valu mathes id.
// @Produce         application/json
// @Tags            Roles
// @Success         200 {object} models.Roles
// @Router          /roles/{id} [delete]
// @Security        BearerAuth
func DeleteRole(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.DELETE_ROLE, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := ctx.Param("id")
	fmt.Println(id)
	role, err := services.DeleteRole(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, role)
}
