package utils

type Permissions struct {
	GET_USERS      bool `json:"get_users"`
	GET_USER_BY_ID bool `json:"get_user_by_id"`
	CREATE_USER    bool `json:"create_user"`
	UPDATE_USER    bool `json:"update_user"`
	DELETE_USER    bool `json:"delete_user"`

	GET_HISTORIES bool `json:"get_histories"`

	GET_ROLES      bool `json:"get_roles"`
	GET_ROLE_BY_ID bool `json:"get_role_by_id"`
	CREATE_ROLE    bool `json:"create_role"`
	UPDATE_ROLE    bool `json:"update_role"`
	DELETE_ROLE    bool `json:"delete_role"`

	GET_OUTLETS      bool `json:"get_outlets"`
	GET_OUTLET_BY_ID bool `json:"get_outlet_by_id"`
	CREATE_OUTLET    bool `json:"create_outlet"`
	UPDATE_OUTLET    bool `json:"update_outlet"`
	DELETE_OUTLET    bool `json:"delete_outlet"`

	GET_WAREHOUSES      bool `json:"get_warehouses"`
	GET_WAREHOUSE_BY_ID bool `json:"get_warehouse_by_id"`
	CREATE_WAREHOUSE    bool `json:"create_warehouse"`
	UPDATE_WAREHOUSE    bool `json:"update_warehouse"`
	DELETE_WAREHOUSE    bool `json:"delete_warehouse"`

	GET_CATEGORIES     bool `json:"get_categories"`
	GET_CATEGORY_BY_ID bool `json:"get_category_by_id"`
	CREATE_CATEGORY    bool `json:"create_category"`
	UPDATE_CATEGORY    bool `json:"update_category"`
	DELETE_CATEGORY    bool `json:"delete_category"`

	GET_PRODUCTS      bool `json:"get_products"`
	GET_PRODUCT_BY_ID bool `json:"get_product_by_id"`
	CREATE_PRODUCT    bool `json:"create_product"`
	UPDATE_PRODUCT    bool `json:"update_product"`
	DELETE_PRODUCT    bool `json:"delete_product"`

	GET_CUSTOMERS      bool `json:"get_customers"`
	GET_CUSTOMER_BY_ID bool `json:"get_customer_by_id"`
	CREATE_CUSTOMER    bool `json:"create_customer"`
	UPDATE_CUSTOMER    bool `json:"update_customer"`
	DELETE_CUSTOMER    bool `json:"delete_customer"`

	GET_WAREHOUSE_ITEMS      bool `json:"get_warehouse_items"`
	GET_WAREHOUSE_ITEM_BY_ID bool `json:"get_warehouse_item_by_id"`
	CREATE_WAREHOUSE_ITEM    bool `json:"create_warehouse_item"`
	UPDATE_WAREHOUSE_ITEM    bool `json:"update_warehouse_item"`
	DELETE_WAREHOUSE_ITEM    bool `json:"delete_warehouse_item"`

	RECEIVE_LOWSTOCK_NOTIFICATION bool `json:"RECEIVE_LOWSTOCK_NOTIFICATION"`

	GET_ORDERS            bool `json:"get_orders"`
	GET_ORDER_BY_ID       bool `json:"get_order_by_id"`
	CREATE_ORDER_BY_ADMIN bool `json:"create_order_by_admin"`
	UPDATE_ORDER_BY_ADMIN bool `json:"update_order_by_admin"`
	DELETE_ORDER_BY_ADMIN bool `json:"delete_order_by_admin"`
	CANCEL_ORDER_BY_ADMIN bool `json:"cancel_order_by_admin"`
}
