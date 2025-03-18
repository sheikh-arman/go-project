package router

import "github.com/sheikh-arman/go-project/pkg/database"

func StartMainRoute() {
	database.DbConnect()
}
