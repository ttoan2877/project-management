package main

import (
	app "Projectmanagement_BE/app"
	controller "Projectmanagement_BE/controller"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func handleRequests() {
	router := mux.NewRouter()

	// default api
	router.HandleFunc("/", nil)

	// users api
	router.Handle("/user", notImplement)
	router.HandleFunc("/api/user/login", controller.AuthenticateUser).Methods("POST") // added log
	router.HandleFunc("/api/user/register", controller.RegisterUser).Methods("POST")  // added log
	router.HandleFunc("/api/user/update-info", controller.UpdateUser).Methods("POST") // added log
	router.HandleFunc("/api/user/get-by-id", controller.GetUserByID).Methods("POST")
	router.HandleFunc("/api/user/search-project", controller.SearchProject).Methods("POST")
	router.HandleFunc("/api/user/search-task", controller.SearchTask).Methods("POST")
	router.HandleFunc("/api/user/search-user", controller.SearchUser).Methods("POST")

	// projects api
	router.HandleFunc("/api/project/create", controller.CreateProject).Methods("POST")                // added log
	router.HandleFunc("/api/project/add-list-user", controller.AddListMember2Project).Methods("POST") // added log
	router.HandleFunc("/api/project/create-role", controller.CreateRole).Methods("POST")              // added log
	router.HandleFunc("/api/project/search-user", controller.SearchUserInProject).Methods("POST")
	router.HandleFunc("/api/project/search-task", controller.SearchTaskInProject).Methods("POST")
	router.HandleFunc("/api/project/update-info", controller.UpdateProject).Methods("POST")         // added log
	router.HandleFunc("/api/project/remove-user", controller.RemoveUserFromProject).Methods("POST") // added log
	router.HandleFunc("/api/project/get-by-id", controller.GetProjectByID).Methods("POST")
	router.HandleFunc("/api/project/search-user-task", controller.SearchUserTaskInProject).Methods("POST")

	// tasks api
	router.HandleFunc("/api/task/create", controller.CreateTask).Methods("POST")          // need add log
	router.HandleFunc("/api/task/assign", controller.AssignTask).Methods("POST")          // need add log
	router.HandleFunc("/api/task/set-todo", controller.SetTODOTask).Methods("POST")       // need add log
	router.HandleFunc("/api/task/set-doing", controller.SetDOINGTask).Methods("POST")     // need add log
	router.HandleFunc("/api/task/set-done", controller.SetDONETask).Methods("POST")       // need add log
	router.HandleFunc("/api/task/set-waiting", controller.SetWAITINGTask).Methods("POST") // need add log
	router.HandleFunc("/api/task/set-delete", controller.SetDELETETask).Methods("POST")   // need add log
	router.HandleFunc("/api/task/update-info", controller.UpdateTask).Methods("POST")     // need add log
	router.HandleFunc("/api/task/unassign-user", controller.UnassignTask).Methods("POST") // need add log
	router.HandleFunc("/api/task/get-by-id", controller.GetTaskByID).Methods("POST")
	router.HandleFunc("/api/task/create-subtask", controller.CreateSubtask).Methods("POST") // need add log
	router.HandleFunc("/api/task/search-user", controller.SearchUserInTask).Methods("POST") // testing
	router.HandleFunc("/api/task/update-subtask", controller.UpdateSubTask).Methods("POST")

	// logs api
	router.Handle("/api/log", notImplement)

	router.Use(app.JwtAuthentication)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(":"+port, router)
	fmt.Println(err)
	if err == nil {
	}
}

func main() {
	handleRequests()
}

// in case api is not implemented yet
var notImplement = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not implemented"))
})
