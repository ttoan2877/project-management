package models

import (
	u "Projectmanagement_BE/utils"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

/*
	Permission rules:
	1 - Delete project/ create role / Update project - get all permission include delete project
	2 - Create/Delete task
	3 - Assign/Unassign member from task
	4 - Add/Remove member from project
*/

// Role struct - model
type Role struct {
	gorm.Model
	ProjectID   *uint   `json:"project_id"`
	Name        *string `json:"name"`
	Permissions *string `json:"permissions"`
	Description *string `json:"description"`
}

// Create - role models
func (role *Role) Create(UserID uint, ProjectID uint) map[string]interface{} {
	// check permissions in request
	permissions := strings.Split(*role.Permissions, ",")
	for i := range permissions {
		temp, err := strconv.Atoi(permissions[i])
		if err != nil {
			return u.Message(false, "Invalid permissions")
		}
		if temp < 1 || temp > 3 {
			return u.Message(false, "No permissions available")
		}
	}
	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserID, ProjectID)
	if ok {
		if roleUserReq == nil {
			return u.Message(false, "No role between user request and project")
		}
	} else {
		return u.Message(false, "Connection error when query role between user request and project")
	}

	// check permissions
	if !strings.Contains(*roleUserReq.Permissions, "0") && !strings.Contains(*roleUserReq.Permissions, "1") {
		return u.Message(false, "Request user doesnt have permissions to create role")
	}

	GetDB().Create(role)

	if role.ID <= 0 {
		return u.Message(false, "Failed to create role, connection error.")
	}

	// Create log
	go CreateProjectLog(UserID, ProjectID, 0, "Created a role", 0) // Log project
	go CreateUserLog(UserID, ProjectID, 0, "Created a role", 0)    // Log User

	response := u.Message(true, "Role has been created")
	response["role"] = role
	return response
}

// GetRoleByID - query role by role_id
func GetRoleByID(RoleID uint) (*Role, bool) {
	role := &Role{}
	err := GetDB().Table("roles").Where("id = ?", RoleID).First(role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return role, true
}

// GetRoleByUserProjectID - query role by user_id and project_id
func GetRoleByUserProjectID(UserID uint, ProjectID uint) (*Role, bool) {
	role := &Role{}
	err := GetDB().Table("roles").Joins("join user_projects on user_projects.role_id = roles.id").Where("user_projects.user_id = ? and user_projects.project_id = ? ", UserID, ProjectID).First(role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return role, true
}

// GetAllRoleByProjectID - query role by project_id
func GetAllRoleByProjectID(ProjectID uint) (*[]Role, bool) {
	role := &[]Role{}
	err := GetDB().Table("roles").Where("project_id = ?", ProjectID).Find(role).Error
	if err != nil {
		if len(*role) > 0 {
			return role, true
		}
		return nil, false
	}
	if len(*role) == 0 {
		return nil, true
	}
	return role, true
}
