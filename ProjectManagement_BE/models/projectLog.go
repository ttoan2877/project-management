package models

import (
	"gorm.io/gorm"
)

// ProjectLog struct
type ProjectLog struct {
	gorm.Model
	UserInfoID         uint   `json:"user_id"`
	UserInfoName       string `json:"user_name"`
	ProjectID          uint   `json:"project_id"`
	TaskInfoID         uint   `json:"task_id"`
	TaskInfoName       string `json:"task_name"`
	Activity           string `json:"activity"`
	TargetUserInfoID   uint   `json:"target_user_id"`
	TargetUserInfoName string `json:"target_user_name"`
}

// CreateProjectLog - model
func CreateProjectLog(UserID uint, ProjectID uint, TaskID uint, Activity string, TargetUserID uint) {
	projectLog := &ProjectLog{}

	if UserID != 0 {
		// get user by UserID
		user, ok := GetUserByID(UserID)
		if ok {
			projectLog.UserInfoID = user.ID
			projectLog.UserInfoName = *user.Name
		}
	}

	projectLog.ProjectID = ProjectID

	if TaskID != 0 {
		// get task by TaskID
		task, ok := GetTaskByID(TaskID)
		if ok {
			projectLog.TaskInfoID = task.ID
			projectLog.TaskInfoName = *task.Name
		}
	}

	projectLog.Activity = Activity

	if TargetUserID != 0 {
		// Get targetUser by UserID
		targetUser, ok := GetUserByID(TargetUserID)
		if ok {
			projectLog.UserInfoID = targetUser.ID
			projectLog.UserInfoName = *targetUser.Name
		}
	}

	GetDB().Create(projectLog)

}

// GetProjectInfoByProjectID - projectLog model
// func GetProjectInfoByProjectID(id uint) (*ProjectInfoInLog, bool) {
// 	projectInfo := &ProjectInfoInLog{}
// 	err := GetDB().Table("project_info_in_logs").Where("project_id = ?", id).First(projectInfo).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, true
// 		}
// 		return nil, false
// 	}

// 	return projectInfo, true
// }
