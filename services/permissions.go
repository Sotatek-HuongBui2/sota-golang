package services

import (
	"encoding/json"
	"errors"
	"vtcanteen/constants"
	errorConstants "vtcanteen/constants/errors"
	"vtcanteen/utils"
)

func CheckPermission(action string, actionerId string, ids ...string) error {
	role, err := GetRoleByUserId(actionerId)
	if err != nil {
		return err
	}
	if !HasPermission(role.Permissions, action) {
		if ids != nil && actionerId == ids[0] {
			return nil
		}
		return errors.New(errorConstants.NOT_PERMISSION)
	}
	return nil
}

func HasPermission(permissions string, action string) bool {
	if permissions == "*" {
		return true
	}

	var permission utils.Permissions
	err := json.Unmarshal([]byte(permissions), &permission)
	if err != nil {
		return false
	}

	switch action {

	case constants.GET_USERS:
		return permission.GET_USERS
	case constants.GET_USER_BY_ID:
		return permission.GET_USER_BY_ID
	case constants.CREATE_USER:
		return permission.CREATE_USER
	case constants.UPDATE_USER:
		return permission.UPDATE_USER
	case constants.DELETE_USER:
		return permission.DELETE_USER

	case constants.GET_HISTORIES:
		return permission.GET_HISTORIES

	case constants.GET_ROLES:
		return permission.GET_ROLES
	case constants.GET_ROLE_BY_ID:
		return permission.GET_ROLE_BY_ID
	case constants.CREATE_ROLE:
		return permission.CREATE_ROLE
	case constants.UPDATE_ROLE:
		return permission.UPDATE_ROLE
	case constants.DELETE_ROLE:
		return permission.DELETE_ROLE

	case constants.GET_OUTLETS:
		return permission.GET_OUTLETS
	case constants.GET_OUTLET_BY_ID:
		return permission.GET_OUTLET_BY_ID
	case constants.CREATE_OUTLET:
		return permission.CREATE_OUTLET
	case constants.UPDATE_OUTLET:
		return permission.UPDATE_OUTLET
	case constants.DELETE_OUTLET:
		return permission.DELETE_OUTLET

	case constants.GET_WAREHOUSES:
		return permission.GET_WAREHOUSES
	case constants.GET_WAREHOUSE_BY_ID:
		return permission.GET_WAREHOUSE_BY_ID
	case constants.CREATE_WAREHOUSE:
		return permission.CREATE_WAREHOUSE
	case constants.UPDATE_WAREHOUSE:
		return permission.UPDATE_WAREHOUSE
	case constants.DELETE_WAREHOUSE:
		return permission.DELETE_WAREHOUSE

	case constants.GET_CUSTOMERS:
		return permission.GET_CUSTOMERS
	case constants.GET_CUSTOMER_BY_ID:
		return permission.GET_CUSTOMER_BY_ID
	case constants.CREATE_CUSTOMER:
		return permission.CREATE_CUSTOMER
	case constants.UPDATE_CUSTOMER:
		return permission.UPDATE_CUSTOMER
	case constants.DELETE_CUSTOMER:
		return permission.DELETE_CUSTOMER

	case constants.GET_WAREHOUSE_ITEMS:
		return permission.GET_WAREHOUSES
	case constants.GET_WAREHOUSE_ITEM_BY_ID:
		return permission.GET_WAREHOUSE_BY_ID
	case constants.CREATE_WAREHOUSE_ITEM:
		return permission.CREATE_WAREHOUSE_ITEM
	case constants.UPDATE_WAREHOUSE_ITEM:
		return permission.UPDATE_WAREHOUSE_ITEM
	case constants.DELETE_WAREHOUSE_ITEM:
		return permission.DELETE_WAREHOUSE_ITEM

	case constants.RECEIVE_LOWSTOCK_NOTIFICATION:
		return permission.RECEIVE_LOWSTOCK_NOTIFICATION

	case constants.GET_PRODUCTS:
		return permission.GET_PRODUCTS
	case constants.GET_PRODUCT_BY_ID:
		return permission.GET_PRODUCT_BY_ID
	case constants.CREATE_PRODUCT:
		return permission.CREATE_PRODUCT
	case constants.UPDATE_PRODUCT:
		return permission.UPDATE_PRODUCT
	case constants.DELETE_PRODUCT:
		return permission.DELETE_PRODUCT

	case constants.GET_CATEGORIES:
		return permission.GET_CATEGORIES
	case constants.GET_CATEGORY_BY_ID:
		return permission.GET_CATEGORY_BY_ID
	case constants.CREATE_CATEGORY:
		return permission.CREATE_CATEGORY
	case constants.UPDATE_CATEGORY:
		return permission.UPDATE_CATEGORY
	case constants.DELETE_CATEGORY:
		return permission.DELETE_CATEGORY

	case constants.GET_ORDERS:
		return permission.GET_ORDERS
	case constants.GET_ORDER_BY_ID:
		return permission.GET_ORDER_BY_ID
	case constants.CREATE_ORDER_BY_ADMIN:
		return permission.CREATE_ORDER_BY_ADMIN
	case constants.UPDATE_ORDER_BY_ADMIN:
		return permission.UPDATE_ORDER_BY_ADMIN
	case constants.DELETE_ORDER_BY_ADMIN:
		return permission.DELETE_ORDER_BY_ADMIN
	case constants.CANCEL_ORDER_BY_ADMIN:
		return permission.CANCEL_ORDER_BY_ADMIN
	default:
		return false
	}
}

func IsSuperAdmin(userId string) (bool, error) {
	role, err := GetRoleByUserId(userId)
	if err != nil {
		return false, err
	}
	if role.Permissions == "*" {
		return true, nil
	}
	return false, errors.New("User is not super admin")
}
