package models

import (
	u "Projectmanagement_BE/utils"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Project struct
type Project struct {
	gorm.Model
	Name        *string      `json:"name"`
	Tasks       []Task       `json:"tasks"`
	CreatorID   uint         `json:"creator_id"`
	Logs        []ProjectLog `json:"logs"`
	Roles       []Role       `json:"roles"` // role in project
	Status      *uint        `json:"status"`
	Description *string      `json:"description"`
}

// UserProject struct - project user relation
type UserProject struct {
	UserID        uint           `json:"user_id" gorm:"primaryKey"`
	ProjectID     uint           `json:"project_id" gorm:"primaryKey"`
	RoleID        uint           `json:"role_id"`
	AddedByUserID uint           `json:"added_by_user_id"`
	CreatedAt     time.Time      `json:"date_created"`
	UpdateAt      time.Time      `json:"date_updated"`
	DeletedAt     gorm.DeletedAt `json:"date_deleted"`
}

// Create project
func (project *Project) Create(UserID uint) map[string]interface{} {

	if project.Name == nil {
		return u.Message(false, "Invalid request")
	}
	project.CreatorID = UserID
	GetDB().Create(project)

	if project.ID <= 0 {
		return u.Message(false, "Failed to create project, connection error.")
	}

	// create role creator and set permimsison
	des := "the author of this project"
	name := "Creator"
	per := "1,2,3,4"
	role := &Role{
		ProjectID:   &project.ID,
		Description: &des,
		Name:        &name,
		Permissions: &per,
	}

	GetDB().Create(role)
	if role.ID <= 0 {
		return u.Message(false, "Failed to create role, conenction error.")
	}

	// create relation for creator to project
	userproject := &UserProject{
		ProjectID:     project.ID,
		UserID:        UserID,
		RoleID:        role.ID,
		AddedByUserID: UserID,
	}
	// userproject.ProjectID = project.ID
	// userproject.UserID = UserID
	// userproject.RoleID = role.ID
	// userproject.AddedByUserID = UserID

	GetDB().Create(userproject)

	// Create log
	go CreateProjectLog(UserID, project.ID, 0, "Created", 0) // Log project
	go CreateUserLog(UserID, project.ID, 0, "Created", 0)    // Log User

	resp := u.Message(true, "")
	resp["project"] = project
	resp["user_project"] = userproject
	resp["role"] = role
	return resp
}

// Update - Project model
func (project *Project) Update(UserID uint) map[string]interface{} {

	if project.Name == nil && project.Description == nil {
		return u.Message(false, "Invalid request")
	}
	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserID, project.ID)
	if ok {
		if roleUserReq == nil {
			return u.Message(false, "No role between user request and project")
		}
	} else {
		return u.Message(false, "Connection error when query role between user request and project")
	}
	// check permissions
	if !strings.Contains(*roleUserReq.Permissions, "1") {
		return u.Message(false, "Request user doesnt have permissions to update project")
	}

	updatedProject, ok := GetProjectByID(project.ID)
	if ok {
		if updatedProject == nil {
			return u.Message(false, "Project not found")
		}
	}
	if !ok {
		return u.Message(false, "Error when query project")
	}

	if project.Name != nil {
		updatedProject.Name = project.Name
	}
	if project.Description != nil {
		updatedProject.Description = project.Description
	}
	GetDB().Save(updatedProject)

	// Create log
	// Get user by UserID
	user, ok := GetUserByID(UserID)
	if ok {
		if user == nil {
			return u.Message(false, "User not found")
		}
	}
	if !ok {
		return u.Message(false, "Error when query user")
	}

	// Create log
	go CreateProjectLog(UserID, project.ID, 0, "Updated", 0) // Log project
	go CreateUserLog(UserID, project.ID, 0, "Updated", 0)    // Log User

	response := u.Message(true, "")
	response["project"] = updatedProject
	return response
}

// AddMember2Project - Add member to project
func AddMember2Project(UserRequestID uint, UserID uint, ProjectID uint, RoleID uint) map[string]interface{} {

	UserIDRespStr := strconv.FormatUint(uint64(UserID), 10)

	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserRequestID, ProjectID)
	if ok {
		if roleUserReq == nil {
			return u.Message(false, UserIDRespStr+": No role between user request and project")
		}
	} else {
		return u.Message(false, UserIDRespStr+": Connection error when query role between user request and project")
	}
	// check permissions
	if !strings.Contains(*roleUserReq.Permissions, "4") {
		return u.Message(false, UserIDRespStr+": Request user doesnt have permissions to add member")
	}

	/*--------------- User request passed ----------------*/

	// check relation of user added and project
	userProjectAdded, ok := GetUserProject(UserID, ProjectID)
	if userProjectAdded != nil {
		return u.Message(false, UserIDRespStr+": User already in project")
	}
	if userProjectAdded == nil {
		if !ok {
			return u.Message(false, UserIDRespStr+": Connection error when query relation between user and project")
		}
	}

	// Check whether project has role or not
	roleUserAdded, ok := GetRoleByID(RoleID)
	if ok {
		if roleUserAdded != nil && *roleUserAdded.ProjectID != ProjectID {
			return u.Message(false, UserIDRespStr+": Project doesnt have this role")
		}
	} else {
		return u.Message(false, UserIDRespStr+": Connection error when query role in project")
	}

	// Add user to project
	userProject := &UserProject{
		UserID:        UserID,
		ProjectID:     ProjectID,
		RoleID:        RoleID,
		AddedByUserID: UserRequestID,
	}
	// userProject.UserID = UserID
	// userProject.ProjectID = ProjectID
	// userProject.RoleID = RoleID
	// userProject.AddedByUserID = UserRequestID
	if GetDB().Create(userProject).Error != nil {
		return u.Message(false, UserIDRespStr+": Error when add user to project")
	}

	// Create log
	go CreateProjectLog(UserRequestID, ProjectID, 0, "Assigned", UserID)
	go CreateUserLog(UserRequestID, ProjectID, 0, "Assigned", UserID)
	go CreateUserLog(UserID, ProjectID, 0, "Assigned by", UserRequestID)
	resp := u.Message(true, UserIDRespStr)
	return resp
}

// GetUserProject - get relation of user_id and project_id
func GetUserProject(UserID uint, ProjectID uint) (*UserProject, bool) {
	userProject := &UserProject{}
	err := GetDB().Table("user_projects").Where("user_id = ? AND project_id = ?", UserID, ProjectID).First(userProject).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return userProject, true
}

// GetProjectByTaskID - get project_id by task_id
func GetProjectByTaskID(TaskID uint) (*Project, bool) {
	project := &Project{}
	err := GetDB().Table("projects").Joins("join tasks on tasks.project_id = projects.id").Where("tasks.id = ? ", TaskID).First(project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return project, true
}

// GetProjectByID - project model
func GetProjectByID(id uint) (*Project, bool) {
	project := &Project{}
	err := GetDB().Table("projects").Where("id = ?", id).First(project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}

	return project, true
}

// DeleteUserProject - UserProject model
func DeleteUserProject(UserRequestID uint, UserID uint, ProjectID uint) map[string]interface{} {

	UserIDRespStr := strconv.FormatUint(uint64(UserID), 10)
	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserRequestID, ProjectID)
	if ok {
		if roleUserReq == nil {
			return u.Message(false, UserIDRespStr+": No role between user request and project")
		}
	} else {
		return u.Message(false, UserIDRespStr+": Connection error when query role between user request and project")
	}

	userProjectRemoved := &UserProject{}
	// if user removed is not user request, check permissions of user request
	if UserRequestID != UserID {
		// check permissions
		if !strings.Contains(*roleUserReq.Permissions, "0") && !strings.Contains(*roleUserReq.Permissions, "2") {
			return u.Message(false, UserIDRespStr+": User request doesnt have permissions to delete member")
		}

		// check relation between user removed and project
		userProjectRemoved, ok := GetUserProject(UserID, ProjectID)

		if userProjectRemoved == nil {
			if !ok {
				return u.Message(false, UserIDRespStr+": Connection error when query relation between user and project")
			}
			if ok {
				return u.Message(false, UserIDRespStr+": User is not in project")
			}
		}

		// check role of user removed
		roleUserRemoved, ok := GetRoleByUserProjectID(UserID, ProjectID)
		if ok {
			if roleUserRemoved == nil {
				return u.Message(false, UserIDRespStr+": No role between user removed and project")
			}
		} else {
			return u.Message(false, UserIDRespStr+": Connection error when query role between user removed and project")
		}

		// check permissions of user removed
		if strings.Contains(*roleUserRemoved.Permissions, "0") && strings.Contains(*roleUserReq.Permissions, "2") {
			return u.Message(false, UserIDRespStr+": Can not remove creator from project")
		}
		if strings.Contains(*roleUserRemoved.Permissions, "2") && strings.Contains(*roleUserReq.Permissions, "2") {
			return u.Message(false, UserIDRespStr+": Can not remove user with the same permission from project")
		}
	}

	/*--------------- User request passed ----------------*/
	// delete user - project relation
	Err := GetDB().Where("user_id = ? AND project_id = ?", UserID, ProjectID).Delete(userProjectRemoved).Error
	if Err != nil {
		return u.Message(false, UserIDRespStr+": Error when delete user from project")
	}

	// delete user from tasks assigned
	// Create log
	go CreateProjectLog(UserID, ProjectID, 0, "Deleted", 0) // Log project
	go CreateUserLog(UserID, ProjectID, 0, "Deleteed", 0)   // Log User
	return u.Message(true, "")
}

// GetProjectByUserID - project model
func GetProjectByUserID(UserID uint, Status *uint, PageSize *uint, PageIndex *uint) (*[]Project, bool) {
	project := &[]Project{}
	if Status == nil {
		if PageSize != nil && PageIndex != nil {
			pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
			err := GetDB().Table("projects").Joins("join user_projects on projects.id = user_projects.project_id").
				Where("user_projects.user_id = ?", UserID).
				Offset(offset).Limit(pageSize).Find(project).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, true
				}
				return nil, false
			}
		}
		return project, true
	}
	if PageSize != nil && PageIndex != nil {
		pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
		err := GetDB().Table("projects").Joins("join user_projects on projects.id = user_projects.project_id").
			Where("user_projects.user_id = ? AND projects.status = ?", UserID, Status).
			Offset(offset).Limit(pageSize).Find(project).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, true
			}
			return nil, false
		}
	}
	return project, true
}
