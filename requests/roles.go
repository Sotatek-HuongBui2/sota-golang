package requests

type CreateOrUpdateRole struct {
	Id          string `json:"id"`
	RoleName    string `json:"name"`
	Permissions string `json:"permissions"`
	Level       int    `json:"level"`
}

type GetRoles struct {
	GetList
}
