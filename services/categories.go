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

func GetCategories(params *requests.GetCategories) (data *utils.IPagination[[]models.Categories], err error) {
	var count int
	categories := []models.Categories{}

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

	_ = query.Model(&models.Categories{}).Count(&count).Error

	err = query.Find(&categories).Error
	data = utils.PaginateResult(categories, count, page, limit)

	return data, err
}

func GetCategoryById(id string) (category *models.Categories, err error) {
	category = &models.Categories{Id: id}
	err = utils.GetConnection().First(category).Error

	if err != nil {
		return category, errors.New(errorConstants.CATEGORY_NOT_FOUND)
	}

	return category, err
}

func GetCategoryByName(categoryName string) (category *models.Categories, err error) {
	category = &models.Categories{}
	err = utils.GetConnection().First(category, models.Categories{CategoryName: categoryName}).Error
	return category, err
}

func CreateCategory(newCategory *requests.CreateOrUpdateCategory) (category *models.Categories, err error) {
	id := uuid.New()

	category = &models.Categories{Id: id.String(), ParentId: newCategory.ParentId, CategoryName: newCategory.CategoryName, IsActive: true}

	fmt.Println("==============", &category)
	err = utils.GetConnection().Create(&category).Error

	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_CATEGORIES, category.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}
	}
	return category, err
}

func UpdateCategory(id string, newCategory *requests.CreateOrUpdateCategory) (category *models.Categories, httpStatus int, err error) {
	category = &models.Categories{Id: id}
	err = utils.GetConnection().First(category).Error

	if newCategory.Id != id {
		return category, http.StatusBadRequest, errors.New(errorConstants.ID_NOT_MATCH)
	}

	if err != nil {
		fmt.Println(err)
		return category, http.StatusNotFound, errors.New(errorConstants.CATEGORY_NOT_FOUND)
	}

	err = utils.GetConnection().Model(category).Update(newCategory).Error

	if err != nil {
		return nil, http.StatusBadRequest, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_CATEGORIES, category.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	return category, http.StatusBadRequest, err
}

func DeleteCategory(id string) (category *models.Categories, err error) {
	category = &models.Categories{Id: id}
	err = utils.GetConnection().First(category).Error

	if err != nil {
		fmt.Println(err)
		return category, errors.New(errorConstants.CATEGORY_NOT_FOUND)
	}
	deletedTime := time.Now()
	category.DeletedAt = &deletedTime
	err = utils.GetConnection().Model(category).Update(&category).Error
	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_CATEGORIES, category.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}
	return category, err
}
