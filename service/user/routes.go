package user

import (
	"fmt"
	"go-mysql-restapi/constants/errors"
	route "go-mysql-restapi/constants/router"
	"go-mysql-restapi/service/auth"
	"go-mysql-restapi/types"
	"go-mysql-restapi/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	store types.UserStore
}

type userService struct {
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// test router
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	}).Methods("GET")

	router.Handle(route.ROUTE_USER_REGISTER, auth.RoleMiddleware()(http.HandlerFunc(h.handleRegisterUser))).Methods("POST")

	router.HandleFunc(route.ROUTE_USER_LOGIN, h.handleLoginUser).Methods("POST")

	router.Handle(route.ROUTE_USER_DELETE, auth.RoleMiddleware()(http.HandlerFunc(h.handleDeleteUser))).Methods("DELETE")
}

func (h *Handler) handleRegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start handleRegisterUser")
	var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Check user if exist
	userByEmail, err := h.store.GetUserByEmail(user.Email)

	if err != nil && err != errors.ErrUserNotFound {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if userByEmail != nil {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrUserExist)
		return
	}

	// Check balance must be more than 0
	if user.Balance < 0 {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidInputBalance)
		return
	}

	// Hash user password
	hashedPassword := auth.HashPassword(user.Password)

	if hashedPassword == "" {
		utils.WriteError(w, http.StatusInternalServerError, errors.ErrInternalServer)
		return
	}

	// Create user
	err = h.store.CreateUser(types.User{
		Name:     user.Name,
		Age:      user.Age,
		Email:    user.Email,
		Password: hashedPassword,
		Balance:  user.Balance,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, errors.ResUserCreated)

}

func (h *Handler) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	var user types.RegisterUserPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userId, err := h.store.LoginUser(user.Email, user.Password)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if userId == "" {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrInvalidEmailOrPasswordLogin)
		return
	}

	isAdmin, _ := h.store.GetUserRoleByID(userId)

	if userId != "" {
		token, err := auth.GenerateJWT(user.Email, isAdmin)

		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Login success", "token": token})
	}
}

func (h *Handler) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	if userID == "1" {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrCannotDeleteRootAdmin)
		return
	}

	// Check user if exist
	existUser, err := h.store.GetUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if existUser == nil {
		utils.WriteError(w, http.StatusBadRequest, errors.ErrUserNotFound)
		return
	}

	err = h.store.DeleteUser(userID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "Delete success"})

}
