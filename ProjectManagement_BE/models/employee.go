package models

import (
	u "Projectmanagement_BE/utils"

	"gorm.io/gorm"
)

// Employee - struct
type Employee struct {
	gorm.Model
	UserID      *uint   `json:"user_id"`
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phone_number"`
	Mail        *string `json:"mail"`
	Bio         *string `json:"bio"`
}

// Validate - models
func (employee *Employee) Validate() (map[string]interface{}, bool) {

	if status, msg := u.CheckValidMail(*employee.Mail); !status {
		return u.Message(status, msg), false
	}

	if status, msg := u.CheckValidPhone(*employee.PhoneNumber); !status {
		return u.Message(status, msg), false
	}
	temp := &Employee{}

	//check for errors and duplicate emails
	err := GetDB().Table("employees").Where("mail = ?", employee.Mail).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Mail != nil {
		return u.Message(false, "Email address already in use by another employee."), false
	}

	// check for errors and duplicate phone nummbers
	err = GetDB().Table("employees").Where("phone_number = ?", employee.PhoneNumber).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.PhoneNumber != nil {
		return u.Message(false, "Phone number already in use by another employee."), false
	}

	return u.Message(false, "Employee requirement passed"), true
}

// GetEmployeeByUserID - model
func GetEmployeeByUserID(UserID uint) (*Employee, bool) {
	employee := &Employee{}
	err := GetDB().Table("employees").Where("user_id = ?", UserID).First(employee).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return employee, true
}
