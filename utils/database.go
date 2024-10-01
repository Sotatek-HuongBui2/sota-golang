package utils

import (
	"fmt"
	"os"
	"vtcanteen/constants"
	"vtcanteen/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func GetConnection() *gorm.DB {
	db.LogMode(true)
	return db
}

func ConnectDB() (err error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASS"), os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_SCHEMA"))
	db, err = gorm.Open("mysql", connection)
	if err != nil {
		panic(err)
	}

	// Migration database
	err = db.AutoMigrate(
		&models.Addresses{},
		&models.Categories{},
		&models.CustomOptions{},
		&models.Customers{},
		&models.Histories{},
		&models.Items{},
		&models.Medias{},
		&models.OptionItems{},
		&models.Options{},
		&models.Orders{},
		&models.Outlets{},
		&models.Payments{},
		&models.ProductCategories{},
		&models.Products{},
		&models.Roles{},
		&models.Transactions{},
		&models.Users{},
		&models.WarehouseItems{},
		&models.Warehouses{},
	).Error
	if err != nil {
		return
	}

	addForeignKey()

	addUnique()

	return
}

func addForeignKey() {
	db.Model(&models.Addresses{}).AddForeignKey(
		"order_id", "orders(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.CustomOptions{}).AddForeignKey(
		"product_id", "products(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Customers{}).AddForeignKey(
		"user_id", "users(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Items{}).AddForeignKey(
		"order_id", "orders(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Items{}).AddForeignKey(
		"product_id", "products(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.OptionItems{}).AddForeignKey(
		"option_id", "custom_options(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.OptionItems{}).AddForeignKey(
		"product_id", "products(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Options{}).AddForeignKey(
		"product_id", "products(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Orders{}).AddForeignKey(
		"customer_id", "customers(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Orders{}).AddForeignKey(
		"accepted_id", "users(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Outlets{}).AddForeignKey(
		"warehouse_id", "warehouses(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Payments{}).AddForeignKey(
		"transaction_id", "transactions(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.ProductCategories{}).AddForeignKey(
		"product_id", "products(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.ProductCategories{}).AddForeignKey(
		"category_id", "categories(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Transactions{}).AddForeignKey(
		"order_id", "orders(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.Users{}).AddForeignKey(
		"role_id", "roles(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.WarehouseItems{}).AddForeignKey(
		"warehouse_id", "warehouses(id)", "CASCADE", "CASCADE",
	)

	db.Model(&models.WarehouseItems{}).AddForeignKey(
		"product_id", "products(id)", "CASCADE", "CASCADE",
	)
}

func addUnique() {
	db.Model(&models.Categories{}).AddUniqueIndex(constants.UNIQUE_CATEGORY_NAME, "category_name")

	db.Model(&models.CustomOptions{}).AddUniqueIndex(constants.UNIQUE_CUSTOM_OPTION_NAME, "option_name")

	db.Model(&models.Customers{}).AddUniqueIndex(constants.UNIQUE_CUSTOMER_USER_NAME, "user_name")

	db.Model(&models.Customers{}).AddUniqueIndex(constants.UNIQUE_CUSTOMER_EMAIL, "email")

	db.Model(&models.Items{}).AddUniqueIndex(constants.UNIQUE_ITEM_BARCODE, "barcode")

	db.Model(&models.Options{}).AddUniqueIndex(constants.UNIQUE_OPTION_NAME, "option_name")

	db.Model(&models.Orders{}).AddUniqueIndex(constants.UNIQUE_ORDER_NUMBER, "order_number")

	db.Model(&models.Outlets{}).AddUniqueIndex(constants.UNIQUE_OUTLET_NAME, "outlet_name")

	db.Model(&models.Roles{}).AddUniqueIndex(constants.UNIQUE_ROLE_NAME, "role_name")

	db.Model(&models.Users{}).AddUniqueIndex(constants.UNIQUE_USER_USER_NAME, "user_name")

	db.Model(&models.Users{}).AddUniqueIndex(constants.UNIQUE_USER_EMAIL, "email")

	db.Model(&models.Warehouses{}).AddUniqueIndex(constants.UNIQUE_WAREHOUSE_NAME, "warehouse_name")
}
