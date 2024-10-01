package services

import (
	"encoding/json"
	"errors"
	"os"
	"time"
	"vtcanteen/constants"
	error_constants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/utils"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(userNameOrEmail string, password string) (tokenStr string, err error) {
	user, err := GetUserByUserNameOrEmail(userNameOrEmail)
	if err != nil {
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("Invalid password")
	}

	tokenStr, err = CreateToken(user.Id)
	return
}

func Register(registerInput *requests.Register) (newUser *models.Users, resp interface{}, err error) {
	if !isValidPassword(registerInput.Password) {
		return nil, nil, errors.New(error_constants.PASSWORD_NOT_SECURE)
	}
	if registerInput.Password != registerInput.RePassword {
		return nil, nil, errors.New("Password confirm incorrect")
	}

	//Create user
	role, err := GetRoleByName(os.Getenv("ROLE_USER_NAME"))
	if err != nil {
		return nil, nil, errors.New(error_constants.ROLE_NOT_FOUND)
	}
	newUser = &models.Users{}
	newUser.Id = uuid.New().String()
	newUser.RoleId = role.Id
	newUser.UserName = registerInput.UserName
	newUser.Email = registerInput.Email
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerInput.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	newUser.Password = string(hashedPassword)
	newUser.IsActive = true

	verificationCode, userMetadataStr, err := GenerateUserMetadataForVerifyRegister()
	if err != nil {
		return nil, nil, errors.New(error_constants.ERR_ENCODING_JSON_METADATA)
	}
	newUser.Metadata = string(userMetadataStr)

	err = utils.GetConnection().Create(&newUser).Error
	if err != nil {
		return nil, nil, utils.GetError(err)
	}
	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_USERS, newUser.Id, constants.ACTION_REGISTER)
		if err != nil {
			return nil, nil, err
		}
	}

	//create customer
	customer := &models.Customers{}
	copier.Copy(&customer, &newUser)
	customer.Id = uuid.New().String()
	customer.UserId = newUser.Id
	err = utils.GetConnection().Create(&customer).Error
	if err != nil {
		return nil, nil, utils.GetError(err)
	}
	{
		// create history
		_, err := CreateHistory(constants.ENTITY_CODE_CUSTOMERS, customer.Id, constants.ACTION_REGISTER)
		if err != nil {
			return nil, nil, err
		}
	}

	//Send mail for verification code
	resp, err = SendMailVerificationRegister(newUser.Email, verificationCode)
	if err != nil {
		return newUser, nil, errors.New(error_constants.ERR_SEND_EMAIL)
	}

	return newUser, resp, err
}

func ResendMailVerificationRegister(userId string) (resp interface{}, err error) {
	user, err := GetUserById(userId)
	if err != nil {
		return nil, err
	}
	var userMetadata utils.UserMetadatas
	err = json.Unmarshal([]byte(user.Metadata), &userMetadata)
	if err != nil {
		return nil, errors.New(error_constants.USER_METADATA_WRONG_FORMAT)
	}
	if userMetadata.IsVerified {
		return nil, errors.New(error_constants.USER_VERIED)
	}

	verificationCode, userMetadataStr, err := GenerateUserMetadataForVerifyRegister()
	if err != nil {
		return nil, errors.New(error_constants.ERR_ENCODING_JSON_METADATA)
	}
	user.Metadata = string(userMetadataStr)

	err = utils.GetConnection().Save(&user).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	//Send mail for verification code
	resp, err = SendMailVerificationRegister(user.Email, verificationCode)
	if err != nil {
		return nil, errors.New(error_constants.ERR_SEND_EMAIL)
	}

	return resp, err
}

func VerifyRegister(userId string, verificationCode string) (user *models.Users, err error) {
	user, err = GetUserById(userId)
	if err != nil {
		return nil, err
	}
	var userMetadata utils.UserMetadatas
	err = json.Unmarshal([]byte(user.Metadata), &userMetadata)
	if err != nil {
		return nil, errors.New(error_constants.USER_METADATA_WRONG_FORMAT)
	}
	if userMetadata.IsVerified {
		return nil, errors.New(error_constants.USER_VERIED)
	}
	if verificationCode != userMetadata.VerificationCode {
		return nil, errors.New(error_constants.VERIFICATION_CODE_INCORRECT)
	}
	if time.Now().Compare(*userMetadata.ExpireAt) == 1 {
		return nil, errors.New(error_constants.VERIFICATION_CODE_IS_EXPIRED)
	}

	userMetadata = utils.UserMetadatas{
		IsVerified:       true,
		VerificationCode: "",
		ExpireAt:         nil,
	}
	userMetadataStr, err := json.Marshal(userMetadata)
	if err != nil {
		return nil, errors.New(error_constants.ERR_ENCODING_JSON_METADATA)
	}
	user.Metadata = string(userMetadataStr)
	err = utils.GetConnection().Save(&user).Error
	if err != nil {
		return nil, utils.GetError(err)
	}

	return user, err
}

func GenerateUserMetadataForVerifyRegister() (verificationCode string, userMetadataStr []byte, err error) {
	verificationCode = randstr.String(10)
	expireAt := time.Unix(time.Now().Unix(), 0).Add(time.Minute * 2)
	userMetadataObj := utils.UserMetadatas{
		IsVerified:       false,
		VerificationCode: verificationCode,
		ExpireAt:         &expireAt,
	}
	userMetadataStr, err = json.Marshal(userMetadataObj)
	if err != nil {
		return "", nil, errors.New(error_constants.ERR_ENCODING_JSON_METADATA)
	}
	return verificationCode, userMetadataStr, err
}
