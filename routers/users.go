package routers

import (
	"errors"
	"fmt"
	"net/http"
	"vtcanteen/constants"
	errorConstants "vtcanteen/constants/errors"
	"vtcanteen/models"
	"vtcanteen/requests"
	"vtcanteen/services"
	"vtcanteen/utils"

	"github.com/gin-gonic/gin"
)

// Login            godoc
// @Summary         Login
// @Description     Return string token jwt
// @Param           body body requests.Login true "Login"
// @Produce         application/json
// @Tags            Authentication
// @Success         200 {string} token
// @Router          /login [post]
func Login(ctx *gin.Context) {
	loginRequest := &requests.Login{}
	if err := ctx.BindJSON(loginRequest); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	token, err := services.Authenticate(loginRequest.UserNameOrEmail, loginRequest.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("Invalid username or password").Error()})
		return
	}
	ctx.JSON(200, gin.H{"token": token})
}

// GetUsers         godoc
// @Summary         Get Users
// @Description     Return list of user.
// @Param           query query requests.GetUsers true "Create user"
// @Produce         application/json
// @Tags            Users
// @Success         200 {object} utils.IPagination[[]models.Users]
// @Router          /users [get]
// @Security        BearerAuth
func GetUsers(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.GET_USERS, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &requests.GetUsers{
		GetList: requests.GetList{
			Limit:        ctx.Query("limit"),
			Sort:         ctx.Query("sort"),
			SortDir:      ctx.Query("sort_dir"),
			SearchFields: ctx.Query("search_fields"),
			Search:       ctx.Query("search"),
			Page:         ctx.Query("page"),
			Filter:       ctx.Query("filter")},
	}

	users := &utils.IPagination[[]models.Users]{}
	users, err := services.GetUsers(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// GetUserById      godoc
// @Summary         Get User by ID
// @Description     Return user whoes userId valu mathes id.
// @Param           id path string true "update user by id"
// @Produce         application/json
// @Tags            Users
// @Success         200 {object} models.Users
// @Router          /users/{id} [get]
// @Security        BearerAuth
func GetUserById(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	id := ctx.Param("id")

	// Check role permission
	if err := services.CheckPermission(constants.GET_USER_BY_ID, actionerId.(string), id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.GetUserById(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// CreateUser       godoc
// @Summary         Create User
// @Description     Save user data in Db.
// @Param           body body models.Users true "Create user"
// @Produce         application/json
// @Tags            Users
// @Success         200 {object} models.Users
// @Router          /users [post]
// @Security        BearerAuth
func CreateUser(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	// Check role permission
	if err := services.CheckPermission(constants.CREATE_USER, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.Users{}
	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
	}

	user, err := services.CreateUser(actionerId.(string), user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// UpdateUser       godoc
// @Summary         Update User
// @Description     Update user data in Db.
// @Param           body body models.Users true "Update user"
// @Produce         application/json
// @Tags            Users
// @Success         200 {object} models.Users
// @Router          /users/{id} [put]
// @Security        BearerAuth
func UpdateUser(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	id := ctx.Param("id")

	// Check role permission
	if err := services.CheckPermission(constants.UPDATE_USER, actionerId.(string), id); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.Users{}
	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
	}

	user, err := services.UpdateUser(actionerId.(string), id, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// DeleteUser       godoc
// @Summary         Delete User
// @Description     Delete user, update deleteTime in Db.
// @Param           id path string true "delete user"
// @Produce         application/json
// @Tags            Users
// @Success         200 {object} models.Users
// @Router          /users/{id} [delete]
// @Security        BearerAuth
func DeleteUser(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	id := ctx.Param("id")

	// Check role permission
	if err := services.CheckPermission(constants.DELETE_USER, actionerId.(string)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.DeleteUser(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// ChangePassword   godoc
// @Summary         Change password
// @Description     Change password.
// @Param           body body requests.ChangePassword true "Change password"
// @Produce         application/json
// @Tags            Users
// @Success         200 {object} models.Users
// @Router          /users/change-password [put]
// @Security        BearerAuth
func ChangePassword(ctx *gin.Context) {
	changePasswordObj := &requests.ChangePassword{}

	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	if err := ctx.BindJSON(changePasswordObj); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	user, err := services.ChangePassword(userId.(string), changePasswordObj)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// ResetPassword    godoc
// @Summary         Reset Password
// @Description     Send mail to reset password.
// @Param           body body requests.MailResetPassword true  "email send to reset password"
// @Produce         application/json
// @Tags            default
// @Router          /reset-password [post]
func SendMailResetPassword(ctx *gin.Context) {
	payload := &requests.MailResetPassword{}
	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
	}

	mess, err := services.SendMailResetPassword(payload.Email)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mess)
}

// ResetPassword    godoc
// @Summary         Reset Password
// @Description     reset password.
// @Param           body body requests.ResetPassword true  "new password to reset"
// @Produce         application/json
// @Tags            Users
// @Router          /users/reset-password [post]
// @Security        BearerAuth
func ResetPassword(ctx *gin.Context) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	payload := &requests.ResetPassword{}
	if err := ctx.BindJSON(payload); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
	}
	user, _ := services.ResetPassword(userId.(string), payload.NewPassword)

	ctx.JSON(http.StatusOK, user)
}

// Register         godoc
// @Summary         Register
// @Description     Return register new account
// @Param           body body requests.Register true "Register"
// @Produce         application/json
// @Tags            Authentication
// @Router          /register [post]
func Register(ctx *gin.Context) {
	registerRequest := &requests.Register{}
	if err := ctx.BindJSON(registerRequest); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	user, mess, err := services.Register(registerRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registerResponse := utils.RegisterResponse{
		User:        user,
		MessageMail: mess,
	}

	ctx.JSON(http.StatusAccepted, registerResponse)
}

// ResendMailVerificationRegister         godoc
// @Summary                               Register
// @Description                           Return register new account
// @Produce                               application/json
// @Tags                                  Users
// @Router                                /users/resend-verification [post]
// @Security                              BearerAuth
func ResendMailVerificationRegister(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}

	mess, err := services.ResendMailVerificationRegister(actionerId.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mess)
}

// VerifyRegister         godoc
// @Summary               VerifyRegister
// @Description           Return verify register
// @Produce               application/json
// @Tags                  Users
// @Router                /users/verify-register [post]
// @Security              BearerAuth
func VerifyRegister(ctx *gin.Context) {
	// Authentication
	actionerId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New(errorConstants.UNAUTHORIZED).Error()})
		return
	}
	verifyRegisterRequest := &requests.VerifyRegister{}
	if err := ctx.BindJSON(verifyRegisterRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New(errorConstants.INVALID_JSON_INPUT).Error()})
		return
	}

	mess, err := services.VerifyRegister(actionerId.(string), verifyRegisterRequest.VerificationCode)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, mess)
}
