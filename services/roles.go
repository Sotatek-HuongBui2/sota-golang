package services

import (
	"encoding/json"
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

func GetRoles(params *requests.GetRoles) (data *utils.IPagination[[]models.Roles], err error) {
	var count int
	roles := []models.Roles{}

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

	_ = query.Model(&models.Roles{}).Count(&count).Error

	err = query.Find(&roles).Error
	data = utils.PaginateResult(roles, count, page, limit)
	return data, err
}

func GetRoleById(id string) (role *models.Roles, err error) {
	role = &models.Roles{Id: id}
	err = utils.GetConnection().First(role, models.Roles{Id: id}).Error

	if err != nil {
		return role, errors.New(errorConstants.ROLE_NOT_FOUND)
	}

	return role, err
}

func GetRoleByName(roleName string) (role *models.Roles, err error) {
	role = &models.Roles{RoleName: roleName}
	err = utils.GetConnection().First(role, models.Roles{RoleName: roleName}).Error
	return role, err
}

func CreateRole(newRole *requests.CreateOrUpdateRole) (role *models.Roles, err error) {
	id := uuid.New()

	newRole.Permissions, err = FormatPermission(newRole.Permissions)
	if err != nil {
		return nil, err
	}
	role = &models.Roles{Id: id.String(), RoleName: newRole.RoleName, Permissions: newRole.Permissions}
	err = utils.GetConnection().Create(&role).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_ROLES, role.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}
	}

	return role, err
}

func GetRoleByUserId(userId string) (role *models.Roles, err error) {
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if utils.IsStringEmpty(user.RoleId) {
		return nil, errors.New("User has no role")
	}
	role, err = GetRoleById(user.RoleId)
	if err != nil {
		return nil, err
	}
	return role, err
}

func UpdateRole(id string, newRole *requests.CreateOrUpdateRole) (role *models.Roles, httpStatus int, err error) {

	role = &models.Roles{Id: id}
	err = utils.GetConnection().First(role).Error

	if newRole.Id != id {
		return role, http.StatusBadRequest, errors.New(errorConstants.ID_NOT_MATCH)
	}

	if err != nil {
		fmt.Println(err)
		return role, http.StatusNotFound, errors.New(errorConstants.ROLE_NOT_FOUND)
	}

	newRole.Permissions, err = FormatPermission(newRole.Permissions)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	err = utils.GetConnection().Model(role).Update(newRole).Error
	if err != nil {
		return nil, http.StatusBadRequest, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_ROLES, role.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, http.StatusBadRequest, err
		}
	}

	return role, http.StatusBadRequest, err
}

func DeleteRole(id string) (role *models.Roles, err error) {
	role = &models.Roles{Id: id}
	err = utils.GetConnection().First(role).Error

	if err != nil {
		fmt.Println(err)
		return role, errors.New(errorConstants.ROLE_NOT_FOUND)
	}
	deletedTime := time.Now()
	role.DeletedAt = &deletedTime
	err = utils.GetConnection().Model(role).Update(&role).Error

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_ROLES, role.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return role, err
}

func FormatPermission(permissions string) (permissionsOutput string, err error) {
	if permissions == "*" {
		return permissions, nil
	}
	var permissionsJson map[string]bool
	err = json.NewDecoder(strings.NewReader(permissions)).Decode(&permissionsJson)
	if err != nil {
		return "", errors.New("Error decode permission json ")
	}
	permissionsEncode, err := json.Marshal(permissionsJson)

	if err != nil {
		return "", errors.New("Error encode permission json ")
	}
	permissionsOutput = string(permissionsEncode)
	return permissionsOutput, nil
}
