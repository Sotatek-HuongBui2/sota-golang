package services

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"vtcanteen/constants"
	error_constants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/utils"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

func GetCustomers(params *requests.GetCustomers) (data *utils.IPagination[[]models.Customers], err error) {
	var count int
	customers := []models.Customers{}

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

	_ = query.Model(&models.Customers{}).Count(&count).Error

	err = query.Find(&customers).Error
	data = utils.PaginateResult(customers, count, page, limit)
	return data, err
}

func CreateCustomer(newCustomer *models.Customers) (customer *models.Customers, err error) {
	//validation
	{
		if !isValidPassword(newCustomer.Password) {
			return nil, errors.New(error_constants.PASSWORD_NOT_SECURE)
		}
	}

	// create user
	user := &models.Users{}
	copier.Copy(&user, &newCustomer)
	fmt.Println(user)
	role, err := GetRoleByName(os.Getenv("ROLE_USER_NAME"))
	if err != nil {
		return nil, errors.New(error_constants.ROLE_NOT_FOUND)
	}
	user.Id = uuid.New().String()
	user.RoleId = role.Id
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	user.Password = string(hashedPassword)
	user.IsActive = true
	verificationCode, userMetadataStr, err := GenerateUserMetadataForVerifyRegister()
	if err != nil {
		return nil, errors.New(error_constants.ERR_ENCODING_JSON_METADATA)
	}
	user.Metadata = string(userMetadataStr)

	err = utils.GetConnection().Create(&user).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	//create customer
	customer = newCustomer
	customer.Id = uuid.New().String()
	customer.UserId = user.Id
	customer.Password = user.Password
	customer.Metadata = user.Metadata
	err = utils.GetConnection().Create(&customer).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err = CreateHistory(constants.ENTITY_CODE_USERS, user.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}

		_, err = CreateHistory(constants.ENTITY_CODE_CUSTOMERS, customer.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}
	}

	//Send mail for verification code
	_, err = SendMailVerificationRegister(user.Email, verificationCode)
	if err != nil {
		return customer, errors.New(error_constants.ERR_SEND_EMAIL)
	}

	return customer, err
}

func UpdateCustomer(actionerId string, id string, customerUpdate *models.Customers) (customer *models.Customers, err error) {
	//validation
	{
		if utils.IsStringEmpty(id) {
			return nil, errors.New("Customer id is empty")
		}
		if id != customerUpdate.Id {
			return nil, errors.New("Customer id not match")
		}
		customer, err = GetCustomerById(id)
		if err != nil {
			return nil, err
		}
		if customer.UserId != customerUpdate.UserId {
			return nil, errors.New("Cannot change user id")
		}

		if customer.UserName != customerUpdate.UserName {
			return nil, errors.New("Cannot change user name")
		}

		roleActioner, err := GetRoleByUserId(actionerId)
		if err != nil {
			return nil, err
		}
		if HasPermission(roleActioner.Permissions, constants.UPDATE_CUSTOMER) {
			if customer.Password != customerUpdate.Password {
				if isValidPassword(customerUpdate.Password) {
					hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(customerUpdate.Password), bcrypt.DefaultCost)
					customerUpdate.Password = string(hashedPassword)
				} else {
					return nil, errors.New(error_constants.PASSWORD_NOT_SECURE)
				}
			}
		} else {
			if customer.Password != customerUpdate.Password {
				return nil, errors.New("Cannot change password in this case")
			}
			if customer.Metadata != "" {
				return nil, errors.New("Cannot change metadata")
			}
		}
	}

	customer = customerUpdate
	err = utils.GetConnection().Save(customer).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	userOld, err := GetUserById(customer.UserId)
	if err != nil {
		return nil, utils.GetError(err)
	}
	user := &models.Users{}
	copier.Copy(&user, &customer)
	user.Id = userOld.Id
	user.RoleId = userOld.RoleId
	err = utils.GetConnection().Save(user).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_CUSTOMERS, customer.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, err
		}

		_, err = CreateHistory(constants.ENTITY_CODE_USERS, user.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, err
		}
	}

	return customer, err
}

func DeleteCustomer(id string) (customer *models.Customers, err error) {
	customer, err = GetCustomerById(id)
	if err != nil {
		return nil, err
	}

	deletedTime := time.Now()
	customer.DeletedAt = &deletedTime
	err = utils.GetConnection().Save(customer).Error

	DeleteUser(customer.UserId)

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_CUSTOMERS, customer.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return customer, err
}

func GetCustomerById(id string) (customer *models.Customers, err error) {
	customer = &models.Customers{}
	err = utils.GetConnection().First(customer, models.Customers{Id: id}).Error
	if err != nil {
		return nil, errors.New(error_constants.CUSTOMER_NOT_FOUND)
	}
	return customer, err
}

func GetCustomerByUserId(id string) (customer *models.Customers, err error) {
	customer = &models.Customers{}
	err = utils.GetConnection().First(customer, models.Customers{UserId: id}).Error
	if err != nil {
		return nil, errors.New(error_constants.CUSTOMER_NOT_FOUND)
	}
	return customer, err
}
