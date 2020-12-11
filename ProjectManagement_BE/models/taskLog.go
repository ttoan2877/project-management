package models

import (
	"gorm.io/gorm"
)

// TaskLog - struct
type TaskLog struct {
	gorm.Model
	UserInfoID         uint   `json:"user_id"`
	UserInfoName       string `json:"user_name"`
	Activity           string `json:"activity"`
	TaskID             uint   `json:"task_id"`
	TargetUserInfoID   uint   `json:"target_user_id"`
	TargetUserInfoName string `json:"target_user_name"`
}

// CreateTaskLog - model
func CreateTaskLog(UserID uint, TaskID uint, Activity string, TargetUserID uint) {
	taskLog := &TaskLog{}

	taskLog.TaskID = TaskID

	taskLog.Activity = Activity

	if UserID != 0 {
		// get user by UserID
		user, ok := GetUserByID(UserID)
		if ok {
			taskLog.UserInfoID = user.ID
			taskLog.UserInfoName = *user.Name
		}
	}

	if TargetUserID != 0 {
		// Get targetUser by TargetUserID
		targetUser, ok := GetUserByID(TargetUserID)
		if ok {
			taskLog.TargetUserInfoID = targetUser.ID
			taskLog.TargetUserInfoName = *targetUser.Name
		}
	}

	GetDB().Create(taskLog)

}

// GetTaskInfoByTaskID - taskLog model
// func GetTaskInfoByTaskID(id uint) (*TaskInfoInLog, bool) {
// 	taskInfo := &TaskInfoInLog{}
// 	err := GetDB().Table("task_info_in_logs").Where("task_id = ?", id).First(taskInfo).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, true
// 		}
// 		return nil, false
// 	}

// 	return taskInfo, true
// }
