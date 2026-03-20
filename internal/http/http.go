package http

import (
	"encoding/json"
	"image/png"
	"log"
	"net/http"
	"test/internal/domain"
	presenters "test/internal/http/presenters"
	"test/internal/usecases"
)

type Object struct {
	taskService usecases.TaskService
	authService usecases.AuthService
}

func NewHandler(taskService usecases.TaskService, authService usecases.AuthService) *Object {
	return &Object{
		taskService: taskService,
		authService: authService,
	}
}

// @Summary      Get Task Status
// @Description  Retrieve the status of a task by its ID.
// @Tags         Task
// @Param        taskID    path      string  true  "Task ID"
// @Param        Authorization header    string  true  "Authorization token in the format 'Bearer <token>'"
// @Produce      json
// @Success      200  {object}  presenters.HandlerResponse
// @Failure      400  {string}  string  "Bad Request"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /task/status/{taskID} [get]
func (s *Object) GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	req, err := presenters.CreateGetHandlerRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	taskStatus, err := s.taskService.GetTaskStatus(req.TaskId)
	presenters.ProcessErrorAndResponse(w, &presenters.HandlerResponse{Value: taskStatus}, err)
}

// @Summary      Get Task Result
// @Description  Retrieve the result of a task by its ID.
// @Tags         Task
// @Param        taskID    path      string  true  "Task ID"
// @Param        Authorization header    string  true  "Authorization token in the format 'Bearer <token>'"
// @Produce      json
// @Success      200  {object}  presenters.HandlerResponse
// @Failure      400  {string}  string  "Bad Request"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /task/result/{taskID} [get]
func (s *Object) GetResultHandler(w http.ResponseWriter, r *http.Request) {
	req, err := presenters.CreateGetHandlerRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := s.taskService.GetTaskResult(req.TaskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if res == nil {
		http.Error(w, "No result yet, please check the status", http.StatusAccepted)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	err = png.Encode(w, res)
	if err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}

}

func (s *Object) CommitHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.Task
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding request body in commitHandler: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.taskService.PutTask(req)
	if err != nil {
		log.Printf("Error saving task in commitHandler: %v", err)
		http.Error(w, "Failed to save task", http.StatusInternalServerError)
		return
	}

	log.Printf("Task %s successfully committed", req.Id)
	presenters.ProcessErrorAndResponse(w, &presenters.HandlerResponse{Value: "task was done, you can check"}, nil)
}

// @Summary      Create Task with Filters
// @Description  Create a new task with specified filters and parameters.
// @Tags         Task
// @Accept       json
// @Produce      json
// @Param        body  body      presenters.FilterRequest  true  "Filter details"
// @Param        Authorization header    string  true  "Authorization token in the format 'Bearer <token>'"
// @Success      200  {object}  presenters.HandlerResponse
// @Failure      400  {string}  string  "Bad Request"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /task [post]
func (s *Object) PostTaskHandlerWithFilters(w http.ResponseWriter, r *http.Request) {
	req, err := presenters.ExtractFiltersFromBody(r)
	if err != nil {
		log.Printf("Error decoding request body in PostTaskHandlerWithFilters: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("decoding request body in PostTaskHandlerWithFilters: %v %v", req.Filter.Name, req.Filter.Parameters)
	taskId, err := s.taskService.NewTask(req.Filter.Name, req.Filter.Parameters.Value)
	presenters.ProcessErrorAndResponse(w, &presenters.HandlerResponse{Value: taskId}, err)
}

// @Summary      Register User
// @Description  Register a new user with a username and password.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      presenters.AuthRequest  true  "User credentials"
// @Success      200  {object}  presenters.HandlerResponse
// @Failure      400  {string}  string  "Bad Request"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /auth/register [post]
func (s *Object) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req, err := presenters.CreateAuthRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.authService.Register(req.Username, req.Password)
	presenters.ProcessErrorAndResponse(w, &presenters.HandlerResponse{Value: "Registration successful"}, err)
}

// @Summary      Login User
// @Description  Login a user and return a session token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      presenters.AuthRequest  true  "User credentials"
// @Success      200  {object}  presenters.HandlerResponse
// @Failure      400  {string}  string  "Bad Request"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /auth/login [post]
func (s *Object) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req, err := presenters.CreateAuthRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sessionToken, err := s.authService.Login(req.Username, req.Password)
	presenters.ProcessErrorAndResponse(w, &presenters.AuthenticatedResponse{Value: "log in was succesful", Token: sessionToken}, err)
}

// @Summary      Middleware for Authentication
// @Description  Middleware to validate user authentication token.
// @Tags         Middleware
// @Failure      400  {string}  string  "Bad Request"
// @Failure      401  {string}  string  "Unauthorized"
func (s *Object) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := presenters.ExtractTokenFromHeader(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = s.authService.Auth(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
