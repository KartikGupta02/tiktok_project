package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/kartik/tiktok_project/internal/controller"

	route "github.com/kartik/tiktok_project/internal/handlers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:K@rtik3275@(127.0.0.1:3306)/mydb")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully Connected with database")

	defer db.Close()

	controller.SetDB(db)

	route.RegisterUserRoutes()
	route.RegisterVideoRoutes()
	http.HandleFunc("/delete_user", controller.DeleteUser)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
