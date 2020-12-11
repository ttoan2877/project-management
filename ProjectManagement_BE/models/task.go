package models

import (
	u "Projectmanagement_BE/utils"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Task struct
type Task struct {
	gorm.Model
	Name        *string   `json:"name"`
	CreatorID   uint      `json:"creator_id"`
	ProjectID   *uint     `json:"project_id"`
	DateEnd     *string   `json:"date_end"`
	Deadline    *string   `json:"date_deadline"`
	Description *string   `json:"description"`
	Status      *uint     `json:"status"`
	Subtasks    []SubTask `json:"subtasks"`
	Logs        []TaskLog `json:"logs"`
}

// UserTask struct
type UserTask struct {
	UserID    uint           `json:"user_id" gorm:"primaryKey"`
	TaskID    uint           `json:"task_id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"date_created"`
	UpdateAt  time.Time      `json:"date_updated"`
	DeletedAt gorm.DeletedAt `json:"date_deleted"`
}

// SubTask struct
type SubTask struct {
	gorm.Model
	TaskID      uint    `json:"task_id"`
	Description *string `json:"description"`
	IsDone      *bool   `json:"is_done"`
}

// Create Task - model
func (task *Task) Create(UserID uint) map[string]interface{} {

	if task.ProjectID == nil || task.Description == nil || task.Name == nil {
		return u.Message(false, "Invalid request")
	}
	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserID, *task.ProjectID)
	if ok {
		if roleUserReq == nil {
			return u.Message(false, "No role between user request and project")
		}
	} else {
		return u.Message(false, "Connection error when query role between user request and project")
	}

	// check permissions
	if !strings.Contains(*roleUserReq.Permissions, "0") && !strings.Contains(*roleUserReq.Permissions, "2") {
		return u.Message(false, "Request user doesnt have permissions to create task")
	}

	// check valid task status
	statusType := [5]uint{1, 2, 3, 4, 5}
	errStatus := false
	for i := range statusType {
		if *task.Status == statusType[i] {
			errStatus = true
			break
		}
	}
	if !errStatus {
		return u.Message(false, "Invalid task status")
	}

	task.CreatorID = UserID
	GetDB().Create(task)

	if task.ID <= 0 {
		return u.Message(false, "Failed to create task, connection error??")
	}

	// Create log
	go CreateUserLog(UserID, 0, task.ID, "Created", 0)
	go CreateTaskLog(UserID, task.ID, "Created", 0)
	resp := u.Message(true, "Task has been created")
	resp["task"] = task
	return resp
}

// Update - model
func (task *Task) Update(UserID uint) map[string]interface{} {

	if task.Name == nil && task.Description == nil && task.Deadline == nil {
		return u.Message(false, "Invalid request")
	}

	// get project
	project, ok := GetProjectByTaskID(task.ID)
	if ok {
		if project == nil {
			return u.Message(false, "Task is not in this project")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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
	if !strings.Contains(*roleUserReq.Permissions, "2") {
		return u.Message(false, "Request user doesnt have permissions to update task")
	}

	updatedTask, ok := GetTaskByID(task.ID)
	if ok {
		if updatedTask == nil {
			return u.Message(false, "Task not found")
		}
	}
	if !ok {
		return u.Message(false, "Error when query task")
	}

	if task.Name != nil {
		updatedTask.Name = task.Name
	}
	if task.Description != nil {
		updatedTask.Description = task.Description
	}
	if task.Deadline != nil {
		updatedTask.Deadline = task.Deadline
	}
	GetDB().Save(updatedTask)

	// Create log
	go CreateUserLog(UserID, 0, task.ID, "Updated", 0)
	go CreateTaskLog(UserID, task.ID, "Updated", 0)
	response := u.Message(true, "")
	response["task"] = updatedTask
	return response
}

// AddMember2Task - model
func AddMember2Task(UserRequestID uint, UserID uint, TaskID uint) map[string]interface{} {

	UserIDRespStr := strconv.FormatUint(uint64(UserID), 10)

	// get project
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, UserIDRespStr+": Task is not in this project")
		}
	} else {
		return u.Message(false, UserIDRespStr+": Connection error when query project")
	}

	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserRequestID, project.ID)
	if ok {
		if roleUserReq == nil {
			return u.Message(false, UserIDRespStr+": No role between user request and project")
		}
	} else {
		return u.Message(false, UserIDRespStr+": Connection error when query role between user request and project")
	}
	// check permissions
	if !strings.Contains(*roleUserReq.Permissions, "0") && !strings.Contains(*roleUserReq.Permissions, "3") {
		return u.Message(false, UserIDRespStr+": Request user doesnt have permissions to add member")
	}

	//--------------- User request passed ----------------

	// check relation of user added and project
	userProjectAdded, ok := GetUserProject(UserID, project.ID)
	if userProjectAdded == nil {
		if !ok {
			return u.Message(false, UserIDRespStr+": Connection error when query relation between user and project")
		}
		return u.Message(false, UserIDRespStr+": This user is not in project")
	}

	// check relation of user added and task
	userTaskAdded, ok := GetUserTask(UserID, TaskID)
	if userTaskAdded != nil {
		return u.Message(false, UserIDRespStr+": User already in task")
	}
	if userTaskAdded == nil {
		if !ok {
			return u.Message(false, UserIDRespStr+": Connection error when query relation between user and task")
		}
	}

	userTask := &UserTask{
		UserID: UserID,
		TaskID: TaskID,
	}

	if GetDB().Create(userTask).Error != nil {
		return u.Message(false, UserIDRespStr+": Error when add user to task")
	}
	// Create log
	go CreateTaskLog(UserRequestID, TaskID, "Assigned", UserID)
	go CreateUserLog(UserRequestID, 0, TaskID, "Assigned", UserID)
	go CreateUserLog(UserID, 0, TaskID, "Assigned by", UserRequestID)
	resp := u.Message(true, UserIDRespStr+": User has been added to task")
	return resp
}

// RemoveUserFromTask - model
func RemoveUserFromTask(UserRequestID uint, UserID uint, TaskID uint) map[string]interface{} {

	UserIDRespStr := strconv.FormatUint(uint64(UserID), 10)

	// check relation of user removed and task
	userTaskRemoved, ok := GetUserTask(UserID, TaskID)
	if userTaskRemoved == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user request and task")
		}
		if ok {
			return u.Message(false, "User not in task")
		}
	}

	// If user request is user removed, unassign from task
	if UserRequestID == UserID {
		Err := GetDB().Table("user_tasks").Where("user_id = ? AND task_id = ?", UserID, TaskID).Delete(userTaskRemoved).Error
		if Err != nil {
			return u.Message(false, "Error when unassign from task")
		}
		return u.Message(true, "")
	}

	// get project of task
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, UserIDRespStr+": Task is not in this project")
		}
	} else {
		return u.Message(false, UserIDRespStr+": Connection error when find project of this task")
	}

	// check role of user request and project
	roleUserReq, ok := GetRoleByUserProjectID(UserRequestID, project.ID)
	if ok {
		if roleUserReq == nil {
			return u.Message(false, UserIDRespStr+": No role between user request and project")
		}
	} else {
		return u.Message(false, UserIDRespStr+": Connection error when query role between user request and project")
	}

	// check permissions user request
	if !strings.Contains(*roleUserReq.Permissions, "0") && !strings.Contains(*roleUserReq.Permissions, "3") {
		return u.Message(false, UserIDRespStr+": Request user doesnt have permissions to remove others from task")
	}

	/*--------Remove user from task----------*/
	Err := GetDB().Table("user_tasks").Where("user_id = ? AND task_id = ?", UserID, TaskID).Delete(userTaskRemoved).Error
	if Err != nil {
		return u.Message(false, "Error when unassign from task")
	}

	// Create log
	go CreateTaskLog(UserRequestID, TaskID, "Removed", UserID)
	go CreateUserLog(UserRequestID, 0, TaskID, "Removed", UserID)
	go CreateUserLog(UserID, 0, TaskID, "Removed by", UserRequestID)
	return u.Message(true, "")

}

// Create SubTask - model
func (subTask *SubTask) Create(UserID uint, TaskID uint) map[string]interface{} {
	// get project
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task is not available")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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
	if !strings.Contains(*roleUserReq.Permissions, "0") && !strings.Contains(*roleUserReq.Permissions, "3") {
		// if user doesnt have permission, check relation of user and task
		// check relation of user request and task
		userTask, ok := GetUserTask(UserID, TaskID)
		if userTask == nil {
			if !ok {
				return u.Message(false, "Connection error when query relation between user and task")
			}
			return u.Message(false, "User is not in this task")
		}
	}
	GetDB().Create(subTask)

	if subTask.ID <= 0 {
		return u.Message(false, "Failed to create subtask, connection error.")
	}

	resp := u.Message(true, "Sub task has been created")
	// Create log
	go CreateTaskLog(UserID, TaskID, "Created subtask", 0)
	go CreateUserLog(UserID, 0, TaskID, "Created subtask", 0)
	resp["subtask"] = subTask
	return resp
}

// Update SubTask - model
func (subtask *SubTask) Update(UserID uint) map[string]interface{} {

	updatedSubTask, ok := GetSubtaskByID(subtask.ID)
	if ok {
		if updatedSubTask == nil {
			return u.Message(false, "Subtask not found")
		}
	}
	if !ok {
		return u.Message(false, "Error when query subtask")
	}
	TaskID := updatedSubTask.TaskID
	// get project
	project, ok := GetProjectByTaskID(TaskID)
	if ok {
		if project == nil {
			return u.Message(false, "Task is not available")
		}
	} else {
		return u.Message(false, "Connection error when query project")
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
	if !strings.Contains(*roleUserReq.Permissions, "0") && !strings.Contains(*roleUserReq.Permissions, "3") {
		// if user doesnt have permission, check relation of user and task
		// check relation of user request and task
		userTask, ok := GetUserTask(UserID, TaskID)
		if userTask == nil {
			if !ok {
				return u.Message(false, "Connection error when query relation between user and task")
			}
			return u.Message(false, "User is not in this task")
		}
	}

	if subtask.Description != nil {
		updatedSubTask.Description = subtask.Description
	}
	if subtask.IsDone != nil {
		updatedSubTask.IsDone = subtask.IsDone
	}
	GetDB().Save(updatedSubTask)

	// Create log
	go CreateUserLog(UserID, 0, TaskID, "Updated a subtask", 0)
	go CreateTaskLog(UserID, TaskID, "Update a subtask", 0)
	response := u.Message(true, "")
	response["subtask"] = updatedSubTask
	return response
}

// GetUserTask - get relation of user_id and task_id
func GetUserTask(UserID uint, TaskID uint) (*UserTask, bool) {
	userTask := &UserTask{}
	err := GetDB().Table("user_tasks").Where("user_id = ? AND task_id = ?", UserID, TaskID).First(userTask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return userTask, true
}

// GetSubTaskByTaskID - get sub tasks of task by task_id
func GetSubTaskByTaskID(TaskID uint) (*[]SubTask, bool) {
	subtask := &[]SubTask{}
	err := GetDB().Table("sub_tasks").Where("task_id = ?", TaskID).Find(subtask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return subtask, true

}

// GetNotDoneSubtask - get sub tasks which are not done
func GetNotDoneSubtask(TaskID uint) (*SubTask, bool) {
	subTask := &SubTask{}
	err := GetDB().Table("sub_tasks").Where("task_id = ? and is_done = ?", TaskID, false).First(subTask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}
	return subTask, true
}

// SetStatusTODO - set status to to do
func SetStatusTODO(UserID uint, TaskID uint) map[string]interface{} {

	// check relation of user request and task
	userTaskAdded, ok := GetUserTask(UserID, TaskID)
	if userTaskAdded == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user and task")
		}
		return u.Message(false, "User is not in this task")
	}

	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 1).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}
	// Create log
	go CreateTaskLog(UserID, TaskID, "Set status to todo", 0)
	go CreateUserLog(UserID, 0, TaskID, "Set status to todo", 0)
	return u.Message(true, "Set status success")
}

// SetStatusDOING - set status to doing
func SetStatusDOING(UserID uint, TaskID uint) map[string]interface{} {

	// check relation of user request and task
	userTaskAdded, ok := GetUserTask(UserID, TaskID)
	if userTaskAdded == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user and task")
		}
		return u.Message(false, "User is not in this task")
	}

	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 2).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}

	// Create log
	go CreateTaskLog(UserID, TaskID, "Set status to doing", 0)
	go CreateUserLog(UserID, 0, TaskID, "Set status to doing", 0)
	return u.Message(true, "")
}

// SetStatusWAITING - set status to waiting
func SetStatusWAITING(UserID uint, TaskID uint) map[string]interface{} {

	// check relation of user request and task
	userTaskAdded, ok := GetUserTask(UserID, TaskID)
	if userTaskAdded == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user and task")
		}
		return u.Message(false, "User is not in this task")
	}

	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 4).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}
	// Create log
	go CreateTaskLog(UserID, TaskID, "Set status to waiting", 0)
	go CreateUserLog(UserID, 0, TaskID, "Set status to waiting", 0)
	return u.Message(true, "")
}

// SetStatusDELETE - set status to delete
func SetStatusDELETE(UserID uint, TaskID uint) map[string]interface{} {

	// check relation of user request and task
	userTaskAdded, ok := GetUserTask(UserID, TaskID)
	if userTaskAdded == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user and task")
		}
		return u.Message(false, "User is not in this task")
	}

	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 5).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}
	// Create log
	go CreateTaskLog(UserID, TaskID, "Set status to delete", 0)
	go CreateUserLog(UserID, 0, TaskID, "Set status to delete", 0)
	return u.Message(true, "Set status success")
}

// SetStatusDONE - set status to done
func SetStatusDONE(UserID uint, TaskID uint) map[string]interface{} {

	// check relation of user request and task
	userTaskAdded, ok := GetUserTask(UserID, TaskID)
	if userTaskAdded == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between user and task")
		}
		return u.Message(false, "User is not in this task")
	}

	// check subtasks done
	subTask, ok := GetNotDoneSubtask(TaskID)
	if subTask == nil {
		if !ok {
			return u.Message(false, "Connection error when query relation between subtask")
		}
	}
	if subTask != nil { //  not done subtask exists
		return u.Message(false, "Still remains subtask")
	}

	// update task by id
	err := GetDB().Table("tasks").Where("ID = ?", TaskID).Update("status", 3).Error
	if err != nil {
		return u.Message(false, "Error when set status")
	}

	// Create log
	go CreateTaskLog(UserID, TaskID, "Set status to done", 0)
	go CreateUserLog(UserID, 0, TaskID, "Set status to done", 0)
	return u.Message(true, "Set status success")
}

// GetTaskByID - task model
func GetTaskByID(id uint) (*Task, bool) {
	task := &Task{}
	err := GetDB().Table("tasks").Where("id = ?", id).Preload("Subtasks").First(task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}

	return task, true
}

// GetSubtaskByID - task model
func GetSubtaskByID(id uint) (*SubTask, bool) {
	subtask := &SubTask{}
	err := GetDB().Table("sub_tasks").Where("id = ?", id).First(subtask).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, true
		}
		return nil, false
	}

	return subtask, true
}

// GetTaskByUserID - task model
func GetTaskByUserID(UserID uint, Status *uint, PageSize *uint, PageIndex *uint) (*[]Task, bool) {
	task := &[]Task{}

	if Status == nil {
		if PageSize != nil && PageIndex != nil {
			pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
			err := GetDB().Table("tasks").Joins("join user_tasks on tasks.id = user_tasks.task_id").
				Where("user_tasks.user_id = ?", UserID).
				Offset(offset).Limit(pageSize).Preload("Subtasks").Find(task).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return nil, true
				}
				return nil, false
			}
		}
		return task, true
	}
	if PageSize != nil && PageIndex != nil {
		pageSize, offset := CalculatePaginate(*PageSize, *PageIndex)
		err := GetDB().Table("tasks").Joins("join user_tasks on tasks.id = user_tasks.task_id").
			Where("user_tasks.user_id = ? AND tasks.status = ?", UserID, Status).
			Offset(offset).Limit(pageSize).Preload("Subtasks").Find(task).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, true
			}
			return nil, false
		}
	}
	return task, true
}
