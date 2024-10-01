package services

import (
	"errors"
	errorConstants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/utils"
)

func GetItemsByOrderId(orderId string) (items []*models.Items, err error) {
	items = []*models.Items{}
	err = utils.GetConnection().Find(&items).Where("orderId = ?", orderId).Error
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, errors.New(errorConstants.ITEM_NOT_FOUND)
	}

	return items, nil
}
