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

func GetProducts(params *requests.GetProducts) (data *utils.IPagination[[]models.Products], err error) {
	var count int
	products := []models.Products{}

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

	_ = query.Model(&models.Products{}).Count(&count).Error

	err = query.Preload("Options").Preload("CustomOptions").Find(&products).Error
	data = utils.PaginateResult(products, count, page, limit)
	return data, err
}

func GetProductById(id string) (product *models.Products, err error) {
	product = &models.Products{Id: id}
	err = utils.GetConnection().Preload("Options").Preload("CustomOptions").Find(product, "id = ?", id).Error

	if err != nil {
		return product, errors.New(errorConstants.PRODUCT_NOT_FOUND)
	}

	return product, err
}

func GetProductByName(productName string) (product *models.Products, err error) {
	product = &models.Products{}
	err = utils.GetConnection().First(product, models.Products{ProductName: productName}).Error
	return product, err
}

func CreateProduct(newProduct *requests.CreateOrUpdateProduct) (product *models.Products, err error) {
	tx := utils.GetConnection().Begin()

	product = &models.Products{
		ParentId:        newProduct.ParentId,
		ProductName:     newProduct.ProductName,
		SKU:             newProduct.SKU,
		Barcode:         newProduct.Barcode,
		Type:            newProduct.Type,
		Price:           newProduct.Price,
		SpecialPrice:    newProduct.SpecialPrice,
		ManageStock:     newProduct.ManageStock,
		SafetyThreshold: newProduct.SafetyThreshold,
		Taxable:         newProduct.Taxable,
		ImageURL:        newProduct.ImageURL,
		IsVariant:       newProduct.IsVariant,
		IsActive:        newProduct.IsActive,
	}

	err = tx.Create(&product).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if newProduct.IsActive == true {
		if newProduct.Type == constants.CONFIG {
			for _, v := range newProduct.Options {
				v.ProductId = product.Id
				err = tx.Create(&v).Error
			}

			if err != nil {
				tx.Rollback()
				return nil, err
			}
		} else if newProduct.Type == constants.BUNDLE {
			newCustomOption := newProduct.CustomOptions
			newCustomOption.ProductId = product.Id
			err = tx.Create(&newCustomOption).Error

			if err != nil {
				tx.Rollback()
				return nil, utils.GetError(err)
			}
		}
	}

	tx.Commit()
	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_PRODUCTS, product.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}
	}

	return product, err
}

func UpdateProduct(id string, newProduct *requests.CreateOrUpdateProduct) (product *models.Products, httpStatus int, err error) {
	product = &models.Products{Id: id}
	err = utils.GetConnection().First(product).Error

	if newProduct.Id != id {
		return product, http.StatusBadRequest, errors.New(errorConstants.ID_NOT_MATCH)
	}

	if err != nil {
		fmt.Println(err)
		return product, http.StatusNotFound, errors.New(errorConstants.PRODUCT_NOT_FOUND)
	}

	tx := utils.GetConnection().Begin()

	if newProduct.IsActive == true {
		if newProduct.Type == constants.CONFIG {
			arrayIdsIsExist := []string{}
			for _, v := range newProduct.Options {
				arrayIdsIsExist = append(arrayIdsIsExist, v.Id)
			}

			deletedTime := time.Now()
			tx.Model(models.Options{}).Where("id NOT IN (?) AND product_id = ?", arrayIdsIsExist, product.Id).Update(models.Options{DeletedAt: &deletedTime})

			for _, v := range newProduct.Options {
				v.ProductId = product.Id
				if len(strings.TrimSpace(v.Id)) == 0 {
					v.Id = uuid.New().String()
				}
				err = tx.Model(models.Options{}).Save(&v).Error
			}

			if err != nil {
				tx.Rollback()
				return nil, http.StatusInternalServerError, err
			}
		} else if newProduct.Type == constants.BUNDLE {

			newCustomOption := newProduct.CustomOptions
			if newCustomOption != nil {
				arrayIdsIsExist := []string{}
				for _, v := range newCustomOption.OptionItems {
					arrayIdsIsExist = append(arrayIdsIsExist, v.Id)
				}

				newCustomOption.ProductId = id
				err = tx.Save(&newCustomOption).Error

				deletedTime := time.Now()
				tx.Model(models.OptionItems{}).Where("id NOT IN (?) AND option_id = ?", arrayIdsIsExist, newCustomOption.Id).Update(models.Options{DeletedAt: &deletedTime})

				for _, v := range newCustomOption.OptionItems {
					v.OptionId = newCustomOption.Id
					if len(strings.TrimSpace(v.Id)) == 0 {
						v.Id = uuid.New().String()
					}
					err = tx.Save(&v).Error
				}

				if err != nil {
					tx.Rollback()
					return nil, http.StatusInternalServerError, utils.GetError(err)
				}
			}

		}
	}

	err = tx.Model(product).Update(newProduct).Error
	tx.Commit()
	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_PRODUCTS, product.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	return product, http.StatusBadRequest, err
}

func DeleteProduct(id string) (product *models.Products, err error) {
	product = &models.Products{Id: id}
	err = utils.GetConnection().First(product).Error

	if err != nil {
		fmt.Println(err)
		return product, errors.New(errorConstants.PRODUCT_NOT_FOUND)
	}
	deletedTime := time.Now()
	product.DeletedAt = &deletedTime
	err = utils.GetConnection().Model(product).Update(&product).Error

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_PRODUCTS, product.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return product, err
}
