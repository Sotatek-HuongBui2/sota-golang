package services

import (
	"errors"
	errorConstants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/utils"
)

func GetAddressByOrderId(orderId string) (address *models.Addresses, err error) {
	address = &models.Addresses{}
	err = utils.GetConnection().Find(&address).Where("orderId = ?", orderId).Error
	if err != nil {
		return nil, errors.New(errorConstants.ADDRESS_NOT_FOUND)
	}

	return address, nil
}
