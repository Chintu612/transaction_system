package routes

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"transaction_system/app/controllers"
)

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Welcome to transaction system")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "Hi, I am transaction system. I am healthy")
}

func InitRoutes(router *httprouter.Router) {

	router.GET("/", HomeHandler)
	router.GET("/health-check", HealthCheckHandler)

	transactionController := controllers.NewTransactionController()
	router.PUT("/transactionservice/transaction/:transaction_id", transactionController.CreateTransaction)
	router.GET("/transactionservice/types/:type", transactionController.GetTransactionsByType)
	router.GET("/transactionservice/sum/:transaction_id", transactionController.GetTransitiveSum)
}
