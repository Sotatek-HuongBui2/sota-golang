package services

import (
	"errors"
	"strconv"
	"strings"
	"time"
	error_constants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/utils"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func GetHistories(params *requests.GetHistories) (data *utils.IPagination[[]models.Histories], err error) {
	var count int
	historys := []models.Histories{}

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

	_ = query.Model(&models.Histories{}).Count(&count).Error

	err = query.Find(&historys).Error
	data = utils.PaginateResult(historys, count, page, limit)
	return data, err
}

func CreateHistory(entityCode string, entityId string, actionName string, metadata ...string) (history *models.Histories, err error) {
	return createHistory(utils.GetConnection(), entityCode, entityId, actionName, metadata...)
}

func CreateHistoryTx(tx *gorm.DB, entityCode string, entityId string, actionName string, metadata ...string) (history *models.Histories, err error) {
	return createHistory(tx, entityCode, entityId, actionName, metadata...)
}

func createHistory(tx *gorm.DB, entityCode string, entityId string, actionName string, metadata ...string) (history *models.Histories, err error) {
	history = &models.Histories{}
	history.Id = uuid.New().String()
	history.EntityCode = entityCode
	history.EntityId = entityId
	history.ActionName = actionName
	if metadata != nil {
		history.Metadata = metadata[0]
	}
	history.ProcessedAt = time.Now().String()

	err = tx.Create(&history).Error
	if err != nil {
		return nil, errors.New(error_constants.HISTORY_CREATE_FAILED)
	}
	return history, nil
}
