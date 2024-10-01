package routers

import (
	"errors"
	"net/http"
	"vtcanteen/constants"
	errorConstants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/services"
	"vtcanteen/utils"

	"github.com/gin-gonic/gin"
)

// GetCustomers     godoc
// @Summary         Get Customers
// @Description     Return list of customer.
// @Param           query query requests.GetCustomers true "Create customer"
// @Produce         application/json
// @Tags            Customers
// @Success         200 {object} utils.IPagination[[]models.Customers]
// @Router          /Customers [get]
// @Security        BearerAuth
func GetCustomers(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_CUSTOMERS, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &requests.GetCustomers{
		GetList: requests.GetList{
			Limit:        ctx.Query("limit"),
			Sort:         ctx.Query("sort"),
			SortDir:      ctx.Query("sort_dir"),
			SearchFields: ctx.Query("search_fields"),
			Search:       ctx.Query("search"),
			Page:         ctx.Query("page"),
			Filter:       ctx.Query("filter")},
	}

	Customers := &utils.IPagination[[]models.Customers]{}
	Customers, err := services.GetCustomers(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, Customers)
}

// GetCustomerById  godoc
// @Summary         Get customer by ID
// @Description     Return customer whoes customerId value mathes id.
// @Param           id path string true "update customer by id"
// @Produce         application/json
// @Tags            Customers
// @Success         200 {object} models.Customers
// @Router          /Customers/{id} [get]
// @Security        BearerAuth
func GetCustomerById(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	id := ctx.Param("id")

	// Check role permission
	if err := services.CheckPermission(constants.GET_CUSTOMER_BY_ID, actionerId.(string), id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := services.GetCustomerById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

// CreateCustomer   godoc
// @Summary         Create customer
// @Description     Save customer data in Db.
// @Param           body body models.Customers true "Create customer"
// @Produce         application/json
// @Tags            Customers
// @Success         200 {object} models.Customers
// @Router          /Customers [post]
// @Security        BearerAuth
func CreateCustomer(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.CREATE_CUSTOMER, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := &models.Customers{}
	if err := ctx.BindJSON(customer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
	}

	customer, err := services.CreateCustomer(customer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

// UpdateCustomer   godoc
// @Summary         Update customer
// @Description     Update customer data in Db.
// @Param           body body models.Customers true "Update customer"
// @Produce         application/json
// @Tags            Customers
// @Success         200 {object} models.Customers
// @Router          /Customers/{id} [put]
// @Security        BearerAuth
func UpdateCustomer(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	id := ctx.Param("id")
	customerUpdate, err := services.GetCustomerById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.UPDATE_CUSTOMER, actionerId.(string), customerUpdate.UserId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := &models.Customers{}
	if err := ctx.BindJSON(customer); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
	}

	customer, err = services.UpdateCustomer(actionerId.(string), id, customer)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}

// DeleteCustomer   godoc
// @Summary         Delete customer
// @Description     Delete customer, update deleteTime in Db.
// @Param           id path string true "delete customer"
// @Produce         application/json
// @Tags            Customers
// @Success         200 {object} models.Customers
// @Router          /Customers/{id} [delete]
// @Security        BearerAuth
func DeleteCustomer(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	id := ctx.Param("id")

	// Check role permission
	if err := services.CheckPermission(constants.DELETE_CUSTOMER, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := services.DeleteCustomer(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, customer)
}
