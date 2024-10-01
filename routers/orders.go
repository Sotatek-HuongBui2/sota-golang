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

// GetOrders        godoc
// @Summary         Get Orders
// @Description     Return list of order.
// @Param           query query requests.GetOrders true "Create order"
// @Produce         application/json
// @Tags            Orders
// @Success         200 {object} utils.IPagination[[]models.Orders]
// @Router          /orders [get]
// @Security        BearerAuth
func GetOrders(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_ORDERS, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &requests.GetOrders{
		GetList: requests.GetList{
			Limit:        ctx.Query("limit"),
			Sort:         ctx.Query("sort"),
			SortDir:      ctx.Query("sort_dir"),
			SearchFields: ctx.Query("search_fields"),
			Search:       ctx.Query("search"),
			Page:         ctx.Query("page"),
			Filter:       ctx.Query("filter")},
	}

	orders := &utils.IPagination[[]models.Orders]{}
	orders, err := services.GetOrders(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orders)
}

// GetOrderById     godoc
// @Summary         Get order by ID
// @Description     Return order whoes orderId value mathes id.
// @Param           id path string true "update order by id"
// @Produce         application/json
// @Tags            Orders
// @Success         200 {object} models.Orders
// @Router          /orders/{id} [get]
// @Security        BearerAuth
func GetOrderById(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	id := ctx.Param("id")

	// Check role permission
	if err := services.CheckPermission(constants.GET_ORDER_BY_ID, actionerId.(string), id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := services.GetOrderById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

// CreateOrderByAdmin      godoc
// @Summary                Create order by admin
// @Description            Save order data in Db.
// @Param                  body body requests.CreateOrUpdateOrder true "Create order"
// @Produce                application/json
// @Tags                   Orders
// @Success                200 {object} models.Orders
// @Router                 /orders/admin [post]
// @Security               BearerAuth
func CreateOrderByAdmin(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.CREATE_ORDER_BY_ADMIN, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestOrder := &requests.CreateOrUpdateOrder{}
	if err := ctx.BindJSON(requestOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
	}

	order, err := services.CreateOrder(actionerId.(string), requestOrder)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

// UpdateOrderByAdmin      godoc
// @Summary                Update order by admin
// @Description            Save order data in Db.
// @Param                  body body requests.CreateOrUpdateOrder true "Update order"
// @Produce                application/json
// @Tags                   Orders
// @Success                200 {object} models.Orders
// @Router                 /orders/admin/{id} [post]
// @Security               BearerAuth
func UpdateOrderByAdmin(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	requestOrder := &requests.CreateOrUpdateOrder{}
	id := ctx.Param("id")

	orderUpdate, err := services.GetOrderById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customerUpdate, err := services.GetCustomerById(orderUpdate.CustomerId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.UPDATE_ORDER_BY_ADMIN, actionerId.(string), customerUpdate.UserId); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.BindJSON(requestOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
	}

	orderUpdated, err := services.UpdateOrder(actionerId.(string), id, requestOrder)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, orderUpdated)
}

// DeleteOrderByAdmin      godoc
// @Summary                Delete order by admin
// @Description            Delete order, update deleteTime in Db.
// @Param                  id path string true "delete order"
// @Produce                application/json
// @Tags                   Orders
// @Success                200 {object} models.Orders
// @Router                 /orders/admin/{id} [delete]
// @Security               BearerAuth
func DeleteOrderByAdmin(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	id := ctx.Param("id")

	// Check role permission
	if err := services.CheckPermission(constants.DELETE_ORDER_BY_ADMIN, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := services.DeleteOrder(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}

// CancelOrderByAdmin      godoc
// @Summary                Cancel order by admin
// @Description            Cancel order
// @Param                  id path string true "cancel order"
// @Produce                application/json
// @Tags                   Orders
// @Success                200 {object} models.Orders
// @Router                 /orders/admin/cancel/{id} [post]
// @Security               BearerAuth
func CancelOrder(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	id := ctx.Param("id")

	// Check role permission
	if err := services.CheckPermission(constants.CANCEL_ORDER_BY_ADMIN, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := services.CancelOrder(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, order)
}
