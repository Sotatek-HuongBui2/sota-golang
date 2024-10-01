package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vtcanteen/constants"
	errorConstants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/utils"

	"github.com/google/uuid"
)

func GetWarehouses(params *requests.GetWarehouses) (data *utils.IPagination[[]models.Warehouses], err error) {
	var count int
	warehouse := []models.Warehouses{}

	query := utils.GetConnection()

	var limit int64
	var page int64

	if !(utils.IsStringEmpty(params.Limit) && utils.IsStringEmpty(params.Page)) {
		limit, _ = strconv.ParseInt(params.Limit, 10, 0)
		page, _ = strconv.ParseInt(params.Page, 10, 0)
		query = query.Limit(limit).Offset(limit * (page - 1))
	}

	if !(utils.IsStringEmpty(params.Search) || utils.IsStringEmpty(params.SearchFields)) {
		searchQuery := ""
		arraySearchField := strings.Split(params.SearchFields, ",")
		for i, v := range arraySearchField {
			if i > 0 {
				searchQuery += " OR "
			}
			searchQuery += v + " LIKE " + "\"%" + params.Search + "%\""
		}

		query = query.Where(searchQuery)
	}

	if !(utils.IsStringEmpty(params.Sort) && utils.IsStringEmpty(params.SortDir)) {
		query = query.Order(params.Sort + " " + params.SortDir)
	}

	_ = query.Model(&models.Warehouses{}).Count(&count).Error

	err = query.Find(&warehouse).Error
	data = utils.PaginateResult(warehouse, count, page, limit)
	return data, err
}

func GetWarehouseById(id string) (warehouse *models.Warehouses, err error) {
	warehouse = &models.Warehouses{Id: id}
	err = utils.GetConnection().First(warehouse).Error

	if err != nil {
		return warehouse, errors.New(errorConstants.WAREHOUSE_NOT_FOUND)
	}

	return warehouse, err
}

func GetWarehouseByName(warehouseName string) (warehouse *models.Warehouses, err error) {
	warehouse = &models.Warehouses{}
	err = utils.GetConnection().First(warehouse, models.Warehouses{WarehouseName: warehouseName}).Error
	return warehouse, err
}

func CreateWarehouse(newWarehouse *requests.CreateOrUpdateWarehouse) (warehouse *models.Warehouses, err error) {
	id := uuid.New()

	warehouse = &models.Warehouses{
		Id:            id.String(),
		WarehouseName: newWarehouse.WarehouseName,
		Country:       newWarehouse.Country,
		CountryCode:   newWarehouse.CountryCode,
		Region:        newWarehouse.Region,
		RegionCode:    newWarehouse.RegionCode,
		IsActive:      true,
	}
	err = utils.GetConnection().Create(&warehouse).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_WAREHOUSES, warehouse.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}
	}

	return warehouse, err
}
func UpdateWarehouse(id string, newWarehouse *requests.CreateOrUpdateWarehouse) (warehouse *models.Warehouses, httpStatus int, err error) {
	warehouse = &models.Warehouses{Id: id}
	err = utils.GetConnection().First(warehouse).Error

	if newWarehouse.Id != id {
		return warehouse, http.StatusBadRequest, errors.New(errorConstants.ID_NOT_MATCH)
	}

	if err != nil {
		fmt.Println(err)
		return warehouse, http.StatusNotFound, errors.New(errorConstants.WAREHOUSE_NOT_FOUND)
	}

	err = utils.GetConnection().Model(warehouse).Update(newWarehouse).Error
	if err != nil {
		return nil, http.StatusBadRequest, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_WAREHOUSES, warehouse.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	return warehouse, http.StatusBadRequest, err
}
func DeleteWarehouse(id string) (warehouse *models.Warehouses, err error) {
	warehouse = &models.Warehouses{Id: id}
	err = utils.GetConnection().First(warehouse).Error

	if err != nil {
		fmt.Println(err)
		return warehouse, errors.New(errorConstants.WAREHOUSE_NOT_FOUND)
	}
	deletedTime := time.Now()
	warehouse.DeletedAt = &deletedTime
	err = utils.GetConnection().Model(warehouse).Update(&warehouse).Error

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_WAREHOUSES, warehouse.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return warehouse, err
}
