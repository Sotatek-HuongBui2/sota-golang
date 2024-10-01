package services

import (
	"errors"
	"fmt"
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

func GetWarehouseItems(params *requests.GetWarehouseItems) (data *utils.IPagination[[]models.WarehouseItems], err error) {
	var count int
	warehouse := []models.WarehouseItems{}

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

	_ = query.Model(&models.WarehouseItems{}).Count(&count).Error

	err = query.Find(&warehouse).Error
	data = utils.PaginateResult(warehouse, count, page, limit)
	return data, err
}

func GetWarehouseItemById(id string) (warehouseItem *models.WarehouseItems, err error) {
	warehouseItem = &models.WarehouseItems{Id: id}
	err = utils.GetConnection().First(warehouseItem).Error

	if err != nil {
		return warehouseItem, errors.New(errorConstants.WAREHOUSE_NOT_FOUND)
	}

	return warehouseItem, err
}

func CreateWarehouseItem(newWarehouseItem *models.WarehouseItems) (warehouseItem *models.WarehouseItems, err error) {
	err = validateWarehouseItem(newWarehouseItem, false)
	if err != nil {
		return nil, err
	}
	warehouseItem = newWarehouseItem
	warehouseItem.Id = uuid.New().String()
	err = utils.GetConnection().Create(&warehouseItem).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_WAREHOUSE_ITEMS, warehouseItem.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}
	}

	return warehouseItem, err
}
func UpdateWarehouseItem(id string, newWarehouseItem *models.WarehouseItems) (warehouseItem *models.WarehouseItems, err error) {
	err = validateWarehouseItem(newWarehouseItem, true)
	if err != nil {
		return nil, err
	}
	warehouseItem = &models.WarehouseItems{Id: id}
	err = utils.GetConnection().First(warehouseItem).Error

	if newWarehouseItem.Id != id {
		return warehouseItem, errors.New(errorConstants.ID_NOT_MATCH)
	}

	if err != nil {
		return warehouseItem, errors.New(errorConstants.WAREHOUSE_NOT_FOUND)
	}

	err = utils.GetConnection().Model(warehouseItem).Update(newWarehouseItem).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_WAREHOUSE_ITEMS, warehouseItem.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, err
		}
	}

	return warehouseItem, err
}

func DeleteWarehouseItem(id string) (warehouseItem *models.WarehouseItems, err error) {
	warehouseItem = &models.WarehouseItems{Id: id}
	err = utils.GetConnection().First(warehouseItem).Error

	if err != nil {
		fmt.Println(err)
		return warehouseItem, errors.New(errorConstants.WAREHOUSE_NOT_FOUND)
	}
	deletedTime := time.Now()
	warehouseItem.DeletedAt = &deletedTime
	err = utils.GetConnection().Model(warehouseItem).Update(&warehouseItem).Error

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_WAREHOUSE_ITEMS, warehouseItem.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return warehouseItem, err
}

func validateWarehouseItem(warehouseItem *models.WarehouseItems, isUpdate bool) (err error) {
	if _, err = GetWarehouseById(warehouseItem.WarehouseId); err != nil {
		return err
	}
	if _, err = GetProductById(warehouseItem.ProductId); err != nil {
		return err
	}
	var availableQty int
	var safetyThreshold int
	if i, err := strconv.Atoi(warehouseItem.AvailaibleQty); err == nil {
		if i < 0 {
			return errors.New(errorConstants.INVALID_WAREHOUSE_ITEM_AVAILAIBLE_QTY)
		}
		availableQty = i
	} else {
		return err
	}
	if i, err := strconv.Atoi(warehouseItem.SafetyThreshold); err == nil {
		if i < 0 {
			return errors.New(errorConstants.INVALID_WAREHOUSE_ITEM_SAFETY_THRESHOLD)
		}
		safetyThreshold = i
	} else {
		return err
	}
	if isUpdate && availableQty <= safetyThreshold {
		return errors.New("Available quantity must be greater than safety threshold")
	}
	return nil
}

func ValidateLowStock(warehouseItemId string) (err error) {
	warehouseItem, err := GetWarehouseItemById(warehouseItemId)
	if err != nil {
		return err
	}

	availableQty, err := strconv.Atoi(warehouseItem.AvailaibleQty)
	if err != nil {
		return errors.New(errorConstants.INVALID_WAREHOUSE_ITEM_AVAILAIBLE_QTY)
	}
	safetyThreshold, err := strconv.Atoi(warehouseItem.SafetyThreshold)
	if err != nil {
		return errors.New(errorConstants.INVALID_WAREHOUSE_ITEM_SAFETY_THRESHOLD)
	}

	if availableQty <= safetyThreshold {
		_, err = SendMailNotifyLowstock(warehouseItemId)
	}

	return err
}
