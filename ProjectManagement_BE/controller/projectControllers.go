package controller

import (
	m "Projectmanagement_BE/models"
	u "Projectmanagement_BE/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// RequestProjectID - form to get request user id
type RequestProjectID struct {
	ProjectID *uint `json:"project_id" sql:"-"`
}

// RequestAddUserProject - struct to get form request user + project
type RequestAddUserProject struct {
	ProjectID    *uint                 `json:"project_id"  sql:"-"`
	ListUserRole []RequestListUserRole `json:"user_role_ids" sql:"-"`
}

// RequestListUserRole - struct to get list user + role
type RequestListUserRole struct {
	UserID *uint `json:"user_id" sql:"-"`
	RoleID *uint `json:"role_id" sql:"-"`
}

// RequestListUserProject - struct to get list user + project
type RequestListUserProject struct {
	ProjectID *uint           `json:"project_id" sql:"-"`
	ListUser  []RequestUserID `json:"user_ids" sql:"-"`
}

// CreateProject - controller
var CreateProject = func(w http.ResponseWriter, r *http.Request) {

	project := &m.Project{}
	err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		fmt.Println(err)
		return
	}
	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	resp := project.Create(user) //Create project with user id
	u.Respond(w, resp)
}

// UpdateProject - controller
var UpdateProject = func(w http.ResponseWriter, r *http.Request) {

	project := &m.Project{}
	err := json.NewDecoder(r.Body).Decode(project) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		fmt.Println(err)
		return
	}
	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	resp := project.Update(UserID)
	u.Respond(w, resp)
}

// AddListMember2Project - controller
var AddListMember2Project = func(w http.ResponseWriter, r *http.Request) {

	request := RequestAddUserProject{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint)
	var temp []map[string]interface{}
	var wg sync.WaitGroup
	wg.Add(len(request.ListUserRole))
	for i := range request.ListUserRole {
		go func(i int) {
			defer wg.Done()
			resp := m.AddMember2Project(UserID, *request.ListUserRole[i].UserID, *request.ProjectID, *request.ListUserRole[i].RoleID)
			temp = append(temp, resp)
		}(i)
	}
	wg.Wait()
	u.MultipleRespond(w, temp)
}

// CreateRole - controller
var CreateRole = func(w http.ResponseWriter, r *http.Request) {
	role := &m.Role{}
	err := json.NewDecoder(r.Body).Decode(role) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		fmt.Println(err)
		return
	}
	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	resp := role.Create(UserID, *role.ProjectID) //Create role with user id
	u.Respond(w, resp)
}

// GetProjectByID - controller
var GetProjectByID = func(w http.ResponseWriter, r *http.Request) {
	request := &RequestProjectID{}
	err := json.NewDecoder(r.Body).Decode(request) // decode request body
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint) //Grab the id of the user that send the request

	if request.ProjectID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	project, ok := m.GetProjectByID(*request.ProjectID)
	if !ok {
		u.Respond(w, u.Message(false, "Error when find project"))
		return
	}
	if ok {
		if project == nil {
			u.Respond(w, u.Message(false, "Project not found"))
			return
		}
	}

	if request.ProjectID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	// check relation between user and project
	userProject, ok := m.GetUserProject(UserID, *request.ProjectID)
	if !ok {
		u.Respond(w, u.Message(false, "Error when find project and user relation"))
		return
	}
	if ok {
		if userProject == nil {
			u.Respond(w, u.Message(false, "No relation between user and project"))
			return
		}
	}

	// if project found and user in project
	// get role of user and project
	role, ok := m.GetRoleByUserProjectID(UserID, *request.ProjectID)
	if !ok {
		u.Respond(w, u.Message(false, "Error when find role of user"))
		return
	}

	roles, ok := m.GetAllRoleByProjectID(*request.ProjectID)
	if !ok {
		u.Respond(w, u.Message(false, "Error when find role of project"))
		return
	}
	resp := u.Message(true, "")
	resp["project"] = project
	resp["role"] = role
	resp["roles"] = roles

	u.Respond(w, resp)

}

// SearchUserInProject - controller
var SearchUserInProject = func(w http.ResponseWriter, r *http.Request) {
	request := &m.RequestSearchUser{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.PageIndex == nil || request.PageSize == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint)
	if request.ProjectID != nil {
		result, ok := m.SearchUserInProject(UserID, request.Query, request.ProjectID, request.PageSize, request.PageIndex)
		if ok {
			if result != nil {
				resp := u.Message(true, "")
				resp["result"] = result
				u.Respond(w, resp)
				return
			}
		}
		if !ok {
			u.Respond(w, u.Message(false, "Error when connect to database"))
			return
		}
	}
	u.Respond(w, u.Message(true, ""))
	return
}

// SearchTaskInProject - controller
var SearchTaskInProject = func(w http.ResponseWriter, r *http.Request) {
	request := &m.RequestSearchProjectTask{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.PageIndex == nil || request.PageSize == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserID := r.Context().Value("user").(uint)
	if request.ProjectID != nil {
		result, ok := m.SearchTaskInProject(UserID, request.Query, request.ProjectID, request.Status, request.PageSize, request.PageIndex)
		if ok {
			if result != nil {
				resp := u.Message(true, "")
				resp["result"] = result
				u.Respond(w, resp)
				return
			}
		}
		if !ok {
			u.Respond(w, u.Message(false, "Error when connect to database"))
			return
		}
	}
	u.Respond(w, u.Message(true, ""))
	return
}

// SearchUserTaskInProject - controller
var SearchUserTaskInProject = func(w http.ResponseWriter, r *http.Request) {
	request := &m.RequestSearchUserTaskInProject{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.PageIndex == nil || request.PageSize == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	if request.UserID != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	UserRequestID := r.Context().Value("user").(uint)

	if request.ProjectID != nil {
		result, ok := m.SearchUserTaskInProject(UserRequestID, request.UserID, request.ProjectID, request.Query, request.Status, request.PageSize, request.PageIndex)
		if ok {
			if result != nil {
				resp := u.Message(true, "")
				resp["result"] = result
				u.Respond(w, resp)
				return
			}
		}
		if !ok {
			u.Respond(w, u.Message(false, "Error when connect to database"))
			return
		}
	}
	u.Respond(w, u.Message(true, ""))
	return
}

// RemoveUserFromProject - controller
var RemoveUserFromProject = func(w http.ResponseWriter, r *http.Request) {
	request := &RequestListUserProject{}
	err := json.NewDecoder(r.Body).Decode(request) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		fmt.Println(err)
		return
	}
	UserID := r.Context().Value("user").(uint)
	var temp []map[string]interface{}

	if request.ProjectID == nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	for i := range request.ListUser {
		resp := m.DeleteUserProject(UserID, *request.ListUser[i].UserID, *request.ProjectID)
		temp = append(temp, resp)
	}
	u.MultipleRespond(w, temp)
}
