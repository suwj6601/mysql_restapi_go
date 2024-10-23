package api

import (
	"database/sql"
	"fmt"
	"go-mysql-restapi/service/user"

	"github.com/gorilla/mux"
)

type Api struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *Api) SetUpRoutes() {
	userStore := user.NewStore(a.DB)
	userHandler := user.NewHandler(userStore)

	userHandler.RegisterRoutes(a.Router)
}

func Test() {
	fmt.Print('a')
}
