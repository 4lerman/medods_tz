package routes

import (
	"database/sql"
	"github.com/4lerman/medods_tz/internal/service/user"
	"github.com/gorilla/mux"
)

func Routes(router *mux.Router, db *sql.DB) {
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(db)
	userService := user.NewHandler(userStore)
	userService.RegisterRoutes(subrouter)
}
