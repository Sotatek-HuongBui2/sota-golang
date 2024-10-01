package requests

type Login struct {
	UserNameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type ChangePassword struct {
	OldPassword   string `json:"old_password"`
	NewPassword   string `json:"new_password"`
	RenewPassword string `json:"renew_password"`
}

type ResetPassword struct {
	NewPassword string `json:"new_password"`
}

type MailResetPassword struct {
	Email string `json:"email"`
}

type GetUsers struct {
	GetList
}

type Register struct {
	UserName   string `json:"user_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	RePassword string `json:"repassword"`
}

type VerifyRegister struct {
	VerificationCode string `json:"verification_token"`
}
