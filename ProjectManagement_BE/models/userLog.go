package models

import (
	"gorm.io/gorm"
)

// UserLog struct
type UserLog struct {
	gorm.Model
	UserID             uint   `json:"user_id"`
	ProjectInfoID      uint   `json:"project_info"`
	ProjectInfoName    string `json:"project_name"`
	TaskInfoID         uint   `json:"task_info"`
	TaskInfoName       string `json:"task_name"`
	TargetUserInfoID   uint   `json:"target_user_info"`
	TargetUserInfoName string `json:"target_user_name"`
	Activity           string `json:"activity"`
}

// CreateUserLog - model
func CreateUserLog(UserID uint, ProjectID uint, TaskID uint, Activity string, TargetUserID uint) {
	userLog := &UserLog{}

	userLog.UserID = UserID

	if ProjectID != 0 {
		// get project by ProjectID
		project, ok := GetProjectByID(ProjectID)
		if ok {
			userLog.ProjectInfoID = project.ID
			userLog.ProjectInfoName = *project.Name
		}
	}
	if TaskID != 0 {
		// get task by TaskID
		task, ok := GetTaskByID(TaskID)
		if ok {
			userLog.TaskInfoID = task.ID
			userLog.TaskInfoName = *task.Name
		}
	}

	userLog.Activity = Activity

	if TargetUserID != 0 {
		// Get user by UserID
		targetUser, ok := GetUserByID(TargetUserID)
		if ok {
			userLog.TargetUserInfoID = targetUser.ID
			userLog.TargetUserInfoName = *targetUser.Name
		}
	}

	GetDB().Create(userLog)

}

// GetUserInfoByUserID - user model
// func GetUserInfoByUserID(id uint) (*UserInfoInLog, bool) {
// 	userInfo := &UserInfoInLog{}
// 	err := GetDB().Table("user_info_in_logs").Where("user_id = ?", id).First(userInfo).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, true
// 		}
// 		return nil, false
// 	}

// 	return userInfo, true
// }
