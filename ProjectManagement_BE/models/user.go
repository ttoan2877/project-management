package models

import (
	u "Projectmanagement_BE/utils"
	"os"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Token struct
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// UserID struct
type UserID struct {
	TaskID uint `json:"task_id" gorm:"-"`
	UserID uint `json:"user_id" gorm:"-"`
}

// User struct
type User struct {
	gorm.Model
	Name     *string   `json:"username"`
	Password *string   `json:"password"`
	Employee *Employee `json:"employee"`
	Projects []Project `gorm:"many2many:user_projects" json:"projects"`
	Tasks    []Task    `gorm:"many2many:user_tasks" json:"tasks"`
	Logs     []UserLog `json:"logs"`
	Token    string    `json:"token" gorm:"-"`
}

// Create - user model
func (user *User) Create() map[string]interface{} {

	if user.Name == nil || user.Password == nil {
		return u.Message(false, "Invalid request")
	}

	if user.Employee.Name == nil || user.Employee.PhoneNumber == nil || user.Employee.Mail == nil || user.Employee.Bio == nil {
		return u.Message(false, "Invalid request")
	}

	if msg, status := user.Validate(); !status {
		return msg
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	*user.Password = string(hashedPassword)

	GetDB().Create(user)
	GetDB().Create(user.Employee)

	if user.ID <= 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}

	// Create log
	go CreateUserLog(user.ID, 0, 0, "Register", 0)
	user.Password = nil //delete password

	response := u.Message(true, "User has been created")
	response["user"] = user
	return response
}

// Update - user model
func (user *User) Update(UserID uint) map[string]interface{} {

	// Get user by UserID
	updatedUser, ok := GetUserByID(UserID)
	if ok {
		if updatedUser == nil {
			return u.Message(false, "User not found")
		}
	}
	if !ok {
		return u.Message(false, "Error when query user")
	}

	// To update record, need to check all the valid request
	// username
	if user.Name != nil {
		temp := &User{}
		//check for errors and duplicate user name
		err := GetDB().Table("users").Where("name = ?", user.Name).First(temp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return u.Message(false, "Connection error. Please retry.")
		}
		if temp.Name != nil {
			return u.Message(false, "User name exists.")
		}
		updatedUser.Name = user.Name
	}

	// password
	if user.Password != nil {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		*updatedUser.Password = string(hashedPassword)
		updatedUser.Password = nil
	}
	GetDB().Save(updatedUser)

	if user.Employee != nil {
		if user.Employee.Name == nil && user.Employee.Mail == nil && user.Employee.Bio == nil && user.Employee.PhoneNumber == nil {
			return u.Message(false, "Invalid request")
		}
		employee, ok := GetEmployeeByUserID(updatedUser.ID)
		if ok {
			if employee == nil {
				return u.Message(false, "Employee not found")
			}
		}
		if !ok {
			return u.Message(false, "Error when query employee")
		}
		if user.Employee.Name != nil {
			employee.Name = user.Employee.Name
		}
		if user.Employee.Mail != nil {
			employee.Mail = user.Employee.Mail
		}
		if user.Employee.Bio != nil {
			employee.Bio = user.Employee.Bio
		}
		if user.Employee.PhoneNumber != nil {
			employee.PhoneNumber = user.Employee.PhoneNumber
		}
		GetDB().Save(employee)
		updatedUser.Employee = employee
	}

	// Create log
	go CreateUserLog(user.ID, 0, 0, "Updated", 0)
	// Respond
	response := u.Message(true, "")
	response["user"] = updatedUser
	return response
}

// UserAuthenticate - user model
func UserAuthenticate(name string, password string) map[string]interface{} {

	user := &User{}
	user, status := GetUserByName(name)
	if status {
		if user == nil {
			return u.Message(false, "User name not found")
		}
	} else {
		return u.Message(false, "Connection error")
	}

	err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Wrong password.")
	}

	//Worked! Logged In
	user.Password = nil

	//Create JWT token
	tk := &Token{UserID: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	// Create log
	go CreateUserLog(user.ID, 0, 0, "Logged in", 0)
	resp := u.Message(true, "Logged In")
	resp["user"] = user
	return resp
}

// GetUserByID - user model
func GetUserByID(id uint) (*User, bool) {
	user := &User{}
	err := GetDB().Table("users").Where("id = ?", id).Preload("Employee").First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}

	user.Password = nil
	return user, true
}

// GetUserByName - user model
func GetUserByName(name string) (*User, bool) {
	user := &User{}
	err := GetDB().Table("users").Where("name = ?", name).Preload("Employee").First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return user, true
}

// Validate - user model
func (user *User) Validate() (map[string]interface{}, bool) {

	if msg, status := user.Employee.Validate(); !status {
		return msg, status
	}

	temp := &User{}

	//check for errors and duplicate user name
	err := GetDB().Table("users").Where("name = ?", user.Name).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry."), false
	}
	if temp.Name != nil {
		return u.Message(false, "User name already in use by another user."), false
	}

	return u.Message(false, "Requirement passed."), true
}

// GetAllUserByRoleID - user model
// func GetAllUserByRoleID(RoleID uint) (*[]User, bool) {
// 	user := &[]User{}
// 	userID := &[]{"user_id"}
// 	err := GetDB().Table("user_projects").Where("role_id = ?", RoleID).Find().Error
// 	if err != nil {
// 		if len(*role) > 0 {
// 			return role, true
// 		}
// 		return nil, false
// 	}
// 	if len(*role) == 0 {
// 		return nil, true
// 	}
// 	return role, true
// }
