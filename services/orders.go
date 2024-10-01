package services

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"vtcanteen/constants"
	error_constants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/utils"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

func GetOrders(params *requests.GetOrders) (data *utils.IPagination[[]models.Orders], err error) {
	var count int
	orders := []models.Orders{}

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

	_ = query.Model(&models.Orders{}).Count(&count).Error

	err = query.Find(&orders).Error
	data = utils.PaginateResult(orders, count, page, limit)
	return data, err
}

func CreateOrder(actionerId string, newOrder *requests.CreateOrUpdateOrder) (order *models.Orders, err error) {
	err = validateCreateOrUpdateOrder(newOrder)
	if err != nil {
		return nil, err
	}
	if actionerId != newOrder.AcceptedId {
		return nil, errors.New("Actioner is not match")
	}

	tx := utils.GetConnection().Begin()

	order = &models.Orders{}
	copier.Copy(order, newOrder)

	// save order initialize
	order.OrderStatus = constants.ORDER_STATUS_NEW
	err = tx.Create(&order).Error
	if err != nil {
		tx.Rollback()
		return nil, utils.GetError(err)
	}

	// save address initialize
	address := &models.Addresses{}
	copier.Copy(address, newOrder.Address)
	address.OrderId = order.Id
	err = tx.Create(&address).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	{
		// get fulfillmentStatus, paymentStatus of order and save order
		fulfillmentStatus, totalPayAmount, err := getFulfillmentStatus(tx, newOrder)
		order.FulfillmentStatus = fulfillmentStatus
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		paymentStatus, err := getPaymentStatus(tx, newOrder, totalPayAmount)
		order.PaymentStatus = paymentStatus
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		err = tx.Save(&order).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_ORDERS, order.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}
		_, err = CreateHistory(constants.ENTITY_CODE_ADDRESSES, address.Id, constants.ACTION_CREATE_ORDER)
		if err != nil {
			return nil, err
		}
	}

	return
}

func UpdateOrder(actionerId string, id string, orderUpdate *requests.CreateOrUpdateOrder) (order *models.Orders, err error) {
	err = validateCreateOrUpdateOrder(orderUpdate)
	if err != nil {
		return nil, err
	}

	if actionerId != orderUpdate.AcceptedId {
		return nil, errors.New("Actioner is not match")
	}

	if orderUpdate.Id != id {
		return nil, errors.New(error_constants.ID_NOT_MATCH)
	}
	tx := utils.GetConnection().Begin()

	order, err = GetOrderById(orderUpdate.Id)
	if err != nil {
		return nil, err
	}

	if order.OrderStatus != constants.ORDER_STATUS_NEW {
		return nil, errors.New("Can not update status")
	}

	err = tx.Save(&order).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	address, err := GetAddressByOrderId(order.Id)
	if err != nil {
		return nil, err
	}
	orderUpdate.Address.Id = address.Id
	err = tx.Save(&orderUpdate.Address).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	itemsExisted, err := GetItemsByOrderId(order.Id)
	itemsUpdated := []*models.Items{}

	if order.PaymentStatus == constants.PAYMENT_STATUS_UNPAID {
		for i := 0; i < len(orderUpdate.Items); i++ {
			itemUpdate := orderUpdate.Items[i]
			// Update items have id
			if itemUpdate.Id != "" {
				for j := 0; j < len(itemsExisted); j++ {
					if itemUpdate.Id == itemsExisted[j].Id {
						itemsUpdated = append(itemsUpdated, itemUpdate)
						err = tx.Save(&itemUpdate).Error
						if err != nil {
							tx.Rollback()
							return nil, err
						}
						continue
					}
				}

			} else {
				// create items donot have id
				err = tx.Create(&itemUpdate).Error
				if err != nil {
					tx.Rollback()
					return nil, err
				}
			}

			_, err = CreateHistoryTx(tx, constants.ENTITY_CODE_ITEMS, itemUpdate.Id, constants.ACTION_UPDATE_ORDER)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	// delete items remaining
	if len(itemsUpdated) != 0 {
		itemsRemainingNotUpdated := []*models.Items{}
		for _, itemExisted := range itemsExisted {
			found := false
			for _, itemUpdated := range itemsUpdated {
				if itemUpdated.Id == itemExisted.Id {
					found = true
					break
				}
			}
			if !found {
				itemsRemainingNotUpdated = append(itemsRemainingNotUpdated, itemExisted)
			}
		}

		for _, item := range itemsRemainingNotUpdated {
			deleteTime := time.Now()
			item.DeletedAt = &deleteTime

			err = tx.Update(&item).Error
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			_, err = CreateHistoryTx(tx, constants.ENTITY_CODE_ITEMS, item.Id, constants.ACTION_UPDATE_ORDER)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	tx.Commit()

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_ORDERS, order.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, err
		}
		_, err = CreateHistory(constants.ENTITY_CODE_ADDRESSES, address.Id, constants.ACTION_UPDATE_ORDER)
		if err != nil {
			return nil, err
		}
	}
	return
}

func DeleteOrder(id string) (order *models.Orders, err error) {
	order, err = GetOrderById(id)
	if err != nil {
		return nil, err
	}

	deletedTime := time.Now()
	order.DeletedAt = &deletedTime
	err = utils.GetConnection().Save(order).Error

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_ORDERS, order.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return order, err
}

func GetOrderById(id string) (order *models.Orders, err error) {
	order = &models.Orders{}
	err = utils.GetConnection().First(order, models.Orders{Id: id}).Error
	if err != nil {
		return nil, errors.New(error_constants.ORDER_NOT_FOUND)
	}
	return order, err
}

func validateCreateOrUpdateOrder(order *requests.CreateOrUpdateOrder) (err error) {
	if _, err := GetCustomerById(order.CustomerId); err != nil {
		return err
	}
	if _, err := GetUserById(order.AcceptedId); err != nil {
		return err
	}

	if order.Address == nil {
		return errors.New("Address can not empty")
	}

	if order.Items == nil {
		return errors.New("Items can not empty")
	}
	return
}

func getFulfillmentStatus(tx *gorm.DB, order *requests.CreateOrUpdateOrder) (fulfillmentStatus string, totalPayAmount float64, err error) {
	totalOrderedQty := 0
	totalFulfillmentQty := 0
	totalRefundQty := 0
	totalPayAmount = 0

	for _, item := range order.Items {
		item.OrderId = order.Id

		{
			// convert orderedQty, fulfillmentQty, refundQty to int; price to float
			var orderedQty int
			var refundQty int
			if i, err := strconv.Atoi(item.OrderedQty); err == nil {
				orderedQty = i
				totalOrderedQty += i
			} else {
				return constants.EMPTY_STRING, 0, err
			}
			if i, err := strconv.Atoi(item.FulfilledQty); err == nil {
				totalFulfillmentQty += i
			} else {
				return constants.EMPTY_STRING, 0, err
			}
			if i, err := strconv.Atoi(item.RefundQty); err == nil {
				refundQty = i
				totalRefundQty += i
			} else {
				return constants.EMPTY_STRING, 0, err
			}
			var priceReality float64
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				specialPrice, err := strconv.ParseFloat(item.SpecialPrice, 64)
				if err != nil || specialPrice == 0 {
					priceReality = price
				} else {
					priceReality = specialPrice
				}
			} else {
				return constants.EMPTY_STRING, 0, err
			}
			// caculate total need pay for order
			totalPayAmount += priceReality * float64(orderedQty-refundQty)
		}

		err = tx.Create(&item).Error
		if err != nil {
			tx.Rollback()
			return constants.EMPTY_STRING, 0, err
		}
		_, err := CreateHistoryTx(tx, constants.ENTITY_CODE_ITEMS, item.Id, constants.ACTION_CREATE_ORDER)
		if err != nil {
			tx.Rollback()
			return constants.EMPTY_STRING, 0, err
		}
	}

	{
		// check fulfillmentStatus
		fulfillmentStatus = constants.FULFILLMENT_STATUS_UNFULFILLED
		if totalOrderedQty-totalRefundQty == totalFulfillmentQty {
			fulfillmentStatus = constants.FULFILLMENT_STATUS_FULFILLED
		}
		if totalFulfillmentQty > 0 {
			if totalOrderedQty-totalRefundQty > totalFulfillmentQty {
				fulfillmentStatus = constants.FULFILLMENT_STATUS_PATRIAL_FULFILLED
			} else {
				fulfillmentStatus = constants.FULFILLMENT_STATUS_REDUNDANT_FULFILLED
			}
		}
	}

	return fulfillmentStatus, totalPayAmount, err
}

func getPaymentStatus(tx *gorm.DB, order *requests.CreateOrUpdateOrder, totalPayAmount float64) (paymentStatus string, err error) {
	var totalPaidAmount float64 = 0

	for _, transaction := range order.Transactions {
		transaction.OrderId = order.Id
		newTransaction := models.Transactions{}
		copier.Copy(newTransaction, transaction)

		err = tx.Create(&newTransaction).Error
		if err != nil {
			tx.Rollback()
			return constants.EMPTY_STRING, err
		}
		_, err := CreateHistoryTx(tx, constants.ENTITY_CODE_TRANSACTIONS, transaction.Id, constants.ACTION_CREATE_ORDER)
		if err != nil {
			tx.Rollback()
			return constants.EMPTY_STRING, err
		}
		for _, payment := range transaction.Payments {
			payment.TransactionId = transaction.Id
			if f, err := strconv.ParseFloat(payment.PaidAmount, 64); err == nil {
				totalPaidAmount += f
			} else {
				return constants.EMPTY_STRING, err
			}
			err = tx.Create(&payment).Error

			if err != nil {
				tx.Rollback()
				return constants.EMPTY_STRING, err
			}
			_, err := CreateHistoryTx(tx, constants.ENTITY_CODE_PAYMENTS, transaction.Id, constants.ACTION_CREATE_ORDER)
			if err != nil {
				tx.Rollback()
				return constants.EMPTY_STRING, err
			}
		}
	}

	{
		// check paymentStatus
		paymentStatus = constants.PAYMENT_STATUS_UNPAID
		if totalPayAmount == totalPaidAmount {
			paymentStatus = constants.PAYMENT_STATUS_PAID
		}
		if totalPaidAmount > 0 {
			if totalPayAmount > totalPaidAmount {
				paymentStatus = constants.PAYMENT_STATUS_PATRIAL_PAID
			} else {
				paymentStatus = constants.PAYMENT_STATUS_REDUNDANT_PAID
			}
		}
	}

	return paymentStatus, err
}

func CancelOrder(id string) (order *models.Orders, err error) {
	order, err = GetOrderById(id)
	if err != nil {
		return nil, err
	}
	order.OrderStatus = constants.ORDER_STATUS_CANCELLED

	err = utils.GetConnection().Save(order).Error

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_ORDERS, order.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return order, err
}
