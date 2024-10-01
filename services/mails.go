package services

import (
	"fmt"
	"net/smtp"
	"os"
	"vtcanteen/constants"
)

func SendMail(toList []string, msg string) {
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	username := os.Getenv("MAILDEV_USER")
	password := os.Getenv("MAILDEV_PASS")

	body := []byte(msg)

	auth := smtp.PlainAuth("", username, password, host)

	err := smtp.SendMail(host+":"+port, auth, username, toList, body)

	// handling the errors
	if err != nil {
		fmt.Println("[Send Mail Error]: ", err)
	} else {
		fmt.Println("Successfully sent mail to all user in toList")
	}

}

func SendMailVerificationRegister(mail string, verificationCode string) (resp interface{}, err error) {
	subject := "Subject: VTCC Verify Account!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body> <h3>Your verification code is: %s</h3> </body></html>", verificationCode)

	template := subject + mime + body
	SendMail([]string{mail}, template)
	resp = interface{}(map[string]string{"message": "Email is sent"})
	return
}

func SendMailNotifyLowstock(warehouseItemId string) (resp interface{}, err error) {
	warehouseItem, err := GetWarehouseItemById(warehouseItemId)
	if err != nil {
		return nil, err
	}

	subject := "Subject: VTCC Verify Account!\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body> <h3 style='color:#ff9966'>Lowstock Warning</h3><br/> <p>Warehouse item <span style='color: #339900;'>{%s}</span> with product <span style='color: #339900;'>{%s}</span> is running out of stock.<br/> Please update the item or change the status of the item! </p> </body></html>", warehouseItemId, warehouseItem.ProductId)

	template := subject + mime + body
	users, err := GetUserByPermission(constants.RECEIVE_LOWSTOCK_NOTIFICATION)
	if err != nil {
		return nil, err
	}
	for _, user := range *users {
		SendMail([]string{user.UserName}, template)
	}
	resp = interface{}(map[string]string{"message": "Email is sent"})
	return
}
