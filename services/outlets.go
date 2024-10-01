package services

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vtcanteen/constants"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/utils"

	errorConstants "vtcanteen/constants/errors"

	"github.com/google/uuid"
)

func GetOutlets(params *requests.GetOutlets) (data *utils.IPagination[[]models.Outlets], err error) {
	var count int
	outlets := []models.Outlets{}

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

	_ = query.Model(&models.Outlets{}).Count(&count).Error

	err = query.Find(&outlets).Error
	data = utils.PaginateResult(outlets, count, page, limit)
	return data, err
}

func GetOutletById(id string) (outlet *models.Outlets, err error) {
	outlet = &models.Outlets{Id: id}
	err = utils.GetConnection().First(outlet).Error

	if err != nil {
		return outlet, errors.New(errorConstants.OUTLET_NOT_FOUND)
	}

	return outlet, err
}

func GetOutletByName(outletName string) (outlet *models.Outlets, err error) {
	outlet = &models.Outlets{}
	err = utils.GetConnection().First(outlet, models.Outlets{OutletName: outletName}).Error
	return outlet, err
}

func CreateOutlet(newOutlet *requests.CreateOrUpdateOutlet) (outlet *models.Outlets, err error) {
	id := uuid.New()

	outlet = &models.Outlets{Id: id.String(), WarehouseId: newOutlet.WarehouseId, OutletName: newOutlet.OutletName, IsActive: true}
	err = utils.GetConnection().Create(&outlet).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_OUTLETS, outlet.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}
	}

	return outlet, err
}

func UpdateOutlet(id string, newOutlet *requests.CreateOrUpdateOutlet) (outlet *models.Outlets, httpStatus int, err error) {
	outlet = &models.Outlets{Id: id}
	err = utils.GetConnection().First(outlet).Error

	if newOutlet.Id != id {
		return outlet, http.StatusBadRequest, errors.New(errorConstants.ID_NOT_MATCH)
	}

	if err != nil {
		fmt.Println(err)
		return outlet, http.StatusNotFound, errors.New(errorConstants.OUTLET_NOT_FOUND)
	}

	err = utils.GetConnection().Model(outlet).Update(newOutlet).Error
	if err != nil {
		return nil, http.StatusBadRequest, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_OUTLETS, outlet.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	return outlet, http.StatusBadRequest, err
}

func DeleteOutlet(id string) (outlet *models.Outlets, err error) {
	outlet = &models.Outlets{Id: id}
	err = utils.GetConnection().First(outlet).Error

	if err != nil {
		fmt.Println(err)
		return outlet, errors.New(errorConstants.OUTLET_NOT_FOUND)
	}
	deletedTime := time.Now()
	outlet.DeletedAt = &deletedTime
	err = utils.GetConnection().Model(outlet).Update(&outlet).Error

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_OUTLETS, outlet.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return outlet, err
}
