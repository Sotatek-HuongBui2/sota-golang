package routers

import (
	"os"
	"vtcanteen/middlewares"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Create() (g *gin.Engine) {

	g = gin.Default()
	if os.Getenv("GO_ENV") == "development" {
		g.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	v1 := g.Group("api/v1")
	{
		v1.POST("/login", Login)
		v1.POST("/reset-password", SendMailResetPassword)
		v1.POST("/register", Register)

		// Use authentication middlewares
		v1.Use(middlewares.RequireAuthentication())

		users := v1.Group("/users")
		{
			// Add user handlers
			users.GET("", GetUsers)
			users.GET("/:id", GetUserById)
			users.POST("", CreateUser)
			users.PUT("/:id", UpdateUser)
			users.DELETE("/:id", DeleteUser)
			users.PUT("/change-password", ChangePassword)
			users.POST("/reset-password", ResetPassword)
			users.POST("/resend-verification", ResendMailVerificationRegister)
			users.POST("/verify-register", VerifyRegister)
		}

		roles := v1.Group("/roles")
		{
			roles.GET("", GetRoles)
			roles.GET("/:id", GetRoleById)
			roles.POST("", CreateRole)
			roles.PUT("/:id", UpdateRole)
			roles.DELETE("/:id", DeleteRole)
		}

		outlets := v1.Group("/outlets")
		{
			outlets.GET("", GetOutlets)
			outlets.GET("/:id", GetOutletById)
			outlets.POST("", CreateOutlet)
			outlets.PUT("/:id", UpdateOutlet)
			outlets.DELETE("/:id", DeleteOutlet)
		}

		warehouses := v1.Group("/warehouses")
		{
			warehouses.GET("", GetWarehouses)
			warehouses.GET("/:id", GetWarehouseById)
			warehouses.POST("", CreateWarehouse)
			warehouses.PUT("/:id", UpdateWarehouse)
			warehouses.DELETE("/:id", DeleteWarehouse)
		}

		histories := v1.Group("/histories")
		{
			histories.GET("", GetHistories)
		}

		categories := v1.Group("/categories")
		{
			categories.GET("", GetCategories)
			categories.GET("/:id", GetCategoryById)
			categories.POST("", CreateCategory)
			categories.PUT("/:id", UpdateCategory)
			categories.DELETE("/:id", DeleteCategory)
		}

		products := v1.Group("/products")
		{
			products.GET("", GetProducts)
			products.GET("/:id", GetProductById)
			products.POST("", CreateProduct)
			products.PUT("/:id", UpdateProduct)
			products.DELETE("/:id", DeleteProduct)
		}

		productVariant := v1.Group("/products/:id/variant")
		{
			productVariant.GET("", GetProductVariants)
			productVariant.GET("/:variant_id", GetProductVariantById)
			productVariant.POST("", CreateProductVariant)
			productVariant.PUT("/:variant_id", UpdateProductVariant)
			productVariant.DELETE("/:variant_id", DeleteProductVariant)
		}

		customers := v1.Group("/customers")
		{
			customers.GET("", GetCustomers)
			customers.GET("/:id", GetCustomerById)
			customers.POST("", CreateCustomer)
			customers.PUT("/:id", UpdateCustomer)
			customers.DELETE("/:id", DeleteCustomer)
		}

		warehouseItems := v1.Group("/warehouse-items")
		{
			warehouseItems.GET("", GetWarehouseItems)
			warehouseItems.GET("/:id", GetWarehouseItemById)
			warehouseItems.POST("", CreateWarehouseItem)
			warehouseItems.PUT("/:id", UpdateWarehouseItem)
			warehouseItems.DELETE("/:id", DeleteWarehouseItem)
		}

		orders := v1.Group("/orders")
		{
			orders.GET("", GetOrders)
			orders.GET("/:id", GetOrderById)
			orders.POST("/admin", CreateOrderByAdmin)
			orders.PUT("/admin/:id", UpdateOrderByAdmin)
			orders.DELETE("/admin/:id", DeleteOrderByAdmin)
			orders.POST("/admin/cancel/:id", CancelOrder)
		}
	}

	return
}
