package error_constants

const (
	UNAUTHORIZED               = "Unauthorized"
	ID_NOT_MATCH               = "Id not match"
	INVALID_JSON_INPUT         = "Invalid json input"
	NOT_PERMISSION             = "Not permission"
	ERR_ENCODING_JSON_METADATA = "Error encoding json metadata"
	ERR_SEND_EMAIL             = "Error when send email"
)

const (
	DB_ERROR_CODE_DUPLICATE_ENTRY             = 1062
	DB_ERROR_CODE_INSERT_MISSING_FOREGION_KEY = 1452
)
