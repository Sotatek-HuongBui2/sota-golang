package services

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
	"vtcanteen/constants"
	errorConstants "vtcanteen/constants/errors"
	error_constants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckExistRoleAndUser() {
	var countRole int
	utils.GetConnection().Model(&models.Roles{}).Count(&countRole)
	var countUser int
	utils.GetConnection().Model(&models.Users{}).Count(&countUser)

	newRole := &requests.CreateOrUpdateRole{}
	newRole.RoleName = "Super Admin"
	newRole.Permissions = "*"
	newRole.Level = constants.LEVEL_1

	user := &models.Users{}
	user.UserName = os.Getenv("ADMIN_USERNAME")
	user.Password = os.Getenv("ADMIN_PASSWORD")

	if countRole == 0 && countUser == 0 {
		role, err := CreateRole(newRole)
		if err != nil {
			fmt.Println(err)
			panic("Create role failed")
		}
		user.RoleId = role.Id
		CreateUser("", user)
	}
	if countRole == 0 && countUser != 0 {
		role, err := CreateRole(newRole)
		if err != nil {
			panic("Create role failed")
		}

		_, err = GetUserByUserName(os.Getenv("ADMIN_USERNAME"))
		if err != nil {
			user.RoleId = role.Id
			CreateUser("", user)
		} else {
			user.RoleId = role.Id
			err = utils.GetConnection().Save(user).Error
			if err != nil {
				panic("Update user failed")
			}
		}

	}
	if countUser == 0 && countRole != 0 {
		role := &models.Roles{}
		err := utils.GetConnection().First(role, &models.Roles{Permissions: "*"}).Error
		if err != nil {
			role, err = CreateRole(newRole)
			if err != nil {
				panic("Create role failed")
			}
		}
		user.RoleId = role.Id
		CreateUser("", user)
	}
}

func GetUsers(params *requests.GetUsers) (data *utils.IPagination[[]models.Users], err error) {
	var count int
	users := []models.Users{}

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

	_ = query.Model(&models.Users{}).Count(&count).Error

	err = query.Find(&users).Error
	data = utils.PaginateResult(users, count, page, limit)
	return data, err
}

func CreateUser(actionerId string, newUser *models.Users) (user *models.Users, err error) {
	//validation
	{
		if !isValidPassword(newUser.Password) {
			return nil, errors.New(errorConstants.PASSWORD_NOT_SECURE)
		}

		if !utils.IsStringEmpty(actionerId) {
			roleActioner, err := GetRoleByUserId(actionerId)
			if err != nil {
				return nil, err
			}
			roleNewUser, err := GetRoleById(newUser.RoleId)
			if err != nil {
				return nil, err
			}
			if roleNewUser.Level < roleActioner.Level {
				return nil, errors.New(errorConstants.ERROR_SET_ROLE_LEVEL)
			}
		}
	}

	user = newUser
	user.Id = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	user.Password = string(hashedPassword)
	user.IsActive = true

	err = utils.GetConnection().Create(&user).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err = CreateHistory(constants.ENTITY_CODE_USERS, user.Id, constants.ACTION_CREATE)
		if err != nil {
			return nil, err
		}
	}

	return user, err
}

func UpdateUser(actionerId string, id string, userUpdate *models.Users) (user *models.Users, err error) {
	//validation
	{
		if utils.IsStringEmpty(id) {
			return nil, errors.New("User id is empty")
		}
		if id != userUpdate.Id {
			return nil, errors.New("User id not match")
		}
		user, err = GetUserById(id)
		if err != nil {
			return nil, err
		}
		if user.UserName != userUpdate.UserName {
			return nil, errors.New("Cannot change user name")
		}

		roleActioner, err := GetRoleByUserId(actionerId)
		if err != nil {
			return nil, err
		}
		roleUserUpdate, err := GetRoleById(userUpdate.RoleId)
		if err != nil {
			return nil, err
		}
		if roleUserUpdate.Level < roleActioner.Level {
			fmt.Println(roleUserUpdate.Level)
			fmt.Println(roleActioner.Level)
			return nil, errors.New(errorConstants.ERROR_SET_ROLE_LEVEL)
		}
		if user.Password != userUpdate.Password {

			if HasPermission(roleActioner.Permissions, constants.UPDATE_USER) {
				if isValidPassword(userUpdate.Password) {
					hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.DefaultCost)
					userUpdate.Password = string(hashedPassword)
				} else {
					return nil, errors.New(errorConstants.PASSWORD_NOT_SECURE)
				}
			} else {
				return nil, errors.New("Cannot change password in this case")
			}
		}
	}

	user = userUpdate
	err = utils.GetConnection().Save(user).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_USERS, user.Id, constants.ACTION_UPDATE)
		if err != nil {
			return nil, err
		}
	}

	return user, err
}

func DeleteUser(id string) (user *models.Users, err error) {
	user, err = GetUserById(id)
	if err != nil {
		return nil, err
	}

	deletedTime := time.Now()
	user.DeletedAt = &deletedTime
	err = utils.GetConnection().Save(user).Error

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_USERS, user.Id, constants.ACTION_DELETE)
		if err != nil {
			return nil, err
		}
	}

	return user, err
}

func GetUserByUserNameOrEmail(userNameOrEmail string) (user *models.Users, err error) {
	user = &models.Users{}
	err = utils.GetConnection().Where("user_name = ? OR email = ?", userNameOrEmail, userNameOrEmail).First(user).Error
	return user, err
}

func GetUserById(id string) (user *models.Users, err error) {
	user = &models.Users{}
	err = utils.GetConnection().First(user, models.Users{Id: id}).Error
	if err != nil {
		return nil, errors.New(errorConstants.USER_NOT_FOUND)
	}
	return user, err
}

func GetUserByUserName(userName string) (user *models.Users, err error) {
	user = &models.Users{}
	err = utils.GetConnection().First(user, models.Users{UserName: userName}).Error
	return user, err
}

func GetUserByPermission(permission string) (users *[]models.Users, err error) {
	users = &[]models.Users{}
	query := utils.GetConnection().
		Model(&models.Users{}).
		Joins("LEFT JOIN `roles` ON `users`.`role_id` = `roles`.`id`")
	if permission == "*" {
		query = query.Where("`roles`.`permissions` = '*'")
		err = query.Find(&users).Error
	} else {
		permission = "%" + permission + ":true%"
		query = query.Where("`roles`.`permissions` LIKE ?", permission)
		err = query.Find(&users).Error
	}

	return users, err
}

func GetUserByEmail(email string) (user *models.Users, err error) {
	user = &models.Users{}
	err = utils.GetConnection().First(user, models.Users{Email: email}).Error
	return user, err
}

func ChangePassword(id string, changePasswordObj *requests.ChangePassword) (user *models.Users, err error) {
	user, err = validateChangePassword(id, changePasswordObj)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordObj.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	// Saving data user
	user.Password = string(hashedPassword)
	err = utils.GetConnection().Save(user).Error

	customer, err := GetCustomerByUserId(user.Id)
	if customer != nil {
		customer.Password = user.Password

		err = utils.GetConnection().Save(customer).Error
		_, err := CreateHistory(constants.ENTITY_CODE_CUSTOMERS, customer.Id, constants.ACTION_CHANGE_PASSWORD)
		if err != nil {
			return nil, err
		}
	}

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_USERS, user.Id, constants.ACTION_CHANGE_PASSWORD)
		if err != nil {
			return nil, err
		}
	}

	return user, err
}

func validateChangePassword(id string, changePasswordObj *requests.ChangePassword) (user *models.Users, err error) {
	user = &models.Users{}
	err = utils.GetConnection().Where("id = ?", id).First(user).Error

	if err != nil {
		return nil, errors.New("User is not exists")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePasswordObj.OldPassword)); err != nil {
		fmt.Println(err)
		return user, errors.New("Password incorrect")
	}
	if !isValidPassword(changePasswordObj.NewPassword) {
		return user, errors.New("Password is not secure")
	}
	if changePasswordObj.NewPassword == changePasswordObj.OldPassword {
		fmt.Println(err)
		return user, errors.New("Same old password")
	}
	if changePasswordObj.NewPassword != changePasswordObj.RenewPassword {
		return user, errors.New("Renew password incorrect")
	}
	return user, err
}

func isValidPassword(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func SendMailResetPassword(email string) (resp interface{}, err error) {
	user := &models.Users{}
	err = utils.GetConnection().First(user, models.Users{Email: email}).Error
	if user == nil {
		resp = interface{}(map[string]string{"message": "Fail"})
		return resp, errors.New(errorConstants.USER_NOT_FOUND)
	}
	token, err := CreateToken(user.Id)
	toEmail := []string{email}

	subject := "Subject: VTCC Reset password!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body> <a href=\"http://localhost:3000/reset-password?token=%s\">Click here to reset your password!</a> </body></html>", token)

	template := subject + mime + body
	SendMail(toEmail, template)
	resp = interface{}(map[string]string{"message": "Email is sent"})
	return
}

func ResetPassword(userId string, newPassword string) (user *models.Users, err error) {
	user = &models.Users{Id: userId}
	err = utils.GetConnection().First(user).Error

	if err != nil {
		fmt.Println(err)
		return user, err
	}

	if !isValidPassword(newPassword) {
		return user, errors.New("Password is not secure")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword)); err == nil {
		fmt.Println(err)
		return user, errors.New("Same old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)

	user.Password = string(hashedPassword)

	if err != nil {
		return
	}

	err = utils.GetConnection().Save(user).Error

	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_USERS, user.Id, constants.ACTION_RESET_PASSWORD)
		if err != nil {
			return nil, err
		}
	}

	return user, err
}

func GetUserByCustomerId(id string) (user *models.Users, err error) {
	customer, err := GetCustomerById(id)
	if err != nil {
		return nil, err
	}
	user = &models.Users{}
	err = utils.GetConnection().First(customer, models.Customers{UserId: customer.UserId}).Error
	if err != nil {
		return nil, errors.New(error_constants.USER_NOT_FOUND)
	}
	return user, err
}
