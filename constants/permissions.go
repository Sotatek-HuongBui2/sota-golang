package constants

const (
	GET_USERS      = "get_users"
	GET_USER_BY_ID = "get_user_by_id"
	CREATE_USER    = "create_user"
	UPDATE_USER    = "update_user"
	DELETE_USER    = "delete_user"

	GET_HISTORIES = "get_histories"

	GET_ROLES      = "get_roles"
	GET_ROLE_BY_ID = "get_role_by_id"
	CREATE_ROLE    = "create_role"
	UPDATE_ROLE    = "update_role"
	DELETE_ROLE    = "delete_role"

	GET_OUTLETS      = "get_outlets"
	GET_OUTLET_BY_ID = "get_outlet_by_id"
	CREATE_OUTLET    = "create_outlet"
	UPDATE_OUTLET    = "update_outlet"
	DELETE_OUTLET    = "delete_outlet"

	GET_WAREHOUSES      = "get_warehouses"
	GET_WAREHOUSE_BY_ID = "get_warehouse_by_id"
	CREATE_WAREHOUSE    = "create_warehouse"
	UPDATE_WAREHOUSE    = "update_warehouse"
	DELETE_WAREHOUSE    = "delete_warehouse"

	GET_CATEGORIES     = "get_categories"
	GET_CATEGORY_BY_ID = "get_category_by_id"
	CREATE_CATEGORY    = "create_category"
	UPDATE_CATEGORY    = "update_category"
	DELETE_CATEGORY    = "delete_category"

	GET_PRODUCTS      = "get_products"
	GET_PRODUCT_BY_ID = "get_product_by_id"
	CREATE_PRODUCT    = "create_product"
	UPDATE_PRODUCT    = "update_product"
	DELETE_PRODUCT    = "delete_product"

	GET_CUSTOMERS      = "get_customers"
	GET_CUSTOMER_BY_ID = "get_customer_by_id"
	CREATE_CUSTOMER    = "create_customer"
	UPDATE_CUSTOMER    = "update_customer"
	DELETE_CUSTOMER    = "delete_customer"

	GET_WAREHOUSE_ITEMS      = "get_warehouse_items"
	GET_WAREHOUSE_ITEM_BY_ID = "get_warehouse_item_by_id"
	CREATE_WAREHOUSE_ITEM    = "create_warehouse_item"
	UPDATE_WAREHOUSE_ITEM    = "update_warehouse_item"
	DELETE_WAREHOUSE_ITEM    = "delete_warehouse_item"

	RECEIVE_LOWSTOCK_NOTIFICATION = "RECEIVE_LOWSTOCK_NOTIFICATION"

	GET_ORDERS            = "get_orders"
	GET_ORDER_BY_ID       = "get_order_by_id"
	CREATE_ORDER_BY_ADMIN = "create_order_by_admin"
	UPDATE_ORDER_BY_ADMIN = "update_order_by_admin"
	DELETE_ORDER_BY_ADMIN = "delete_order_by_admin"
	CANCEL_ORDER_BY_ADMIN = "cancel_order_by_admin"
)

// level
const (
	LEVEL_1 int = iota + 1
	LEVEL_2
	LEVEL_3
)
