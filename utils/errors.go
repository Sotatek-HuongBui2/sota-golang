package utils

import (
	"errors"
	"strings"
	"vtcanteen/constants"
	errorConsants "vtcanteen/constants/errors"

	"github.com/go-sql-driver/mysql"
)

func GetError(dbErr error) (err error) {
	ne, ok := dbErr.(*mysql.MySQLError)
	if !ok {
		return dbErr
	}
	if ne.Number == errorConsants.DB_ERROR_CODE_DUPLICATE_ENTRY {
		switch true {
		case strings.Contains(ne.Message, constants.UNIQUE_OUTLET_NAME):
			return errors.New(errorConsants.OUTLET_NAME_IS_EXISTED)

		case strings.Contains(ne.Message, constants.UNIQUE_ROLE_NAME):
			return errors.New(errorConsants.ROLE_NAME_IS_EXISTED)

		case strings.Contains(ne.Message, constants.UNIQUE_USER_USER_NAME):
			return errors.New(errorConsants.USER_NAME_IS_EXISTED)

		case strings.Contains(ne.Message, constants.UNIQUE_USER_EMAIL):
			return errors.New(errorConsants.EMAIL_IS_EXISTED)

		case strings.Contains(ne.Message, constants.UNIQUE_WAREHOUSE_NAME):
			return errors.New(errorConsants.WAREHOUSE_NAME_IS_EXISTED)

		case strings.Contains(ne.Message, constants.UNIQUE_ORDER_NUMBER):
			return errors.New(errorConsants.ORDER_NUMBER_IS_EXISTED)

		default:
			return dbErr
		}
	}

	if ne.Number == errorConsants.DB_ERROR_CODE_INSERT_MISSING_FOREGION_KEY {
		return errors.New("Missing foreign key when insert or update")
	}

	return nil
}
