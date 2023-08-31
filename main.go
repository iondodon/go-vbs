package main

import (
	"github.com/gorilla/mux"
	"go-vbs/controller"
	"go-vbs/repository"
	"go-vbs/usecase"
	"net/http"
)

// a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

func main() {
	vrp := repository.NewVehicleRepository()
	gvuc := usecase.NewGetVehicleUseCase(vrp)
	vc := controller.NewVehicleController(gvuc)

	r := mux.NewRouter()

	r.HandleFunc("/vehicles/{uuid}", vc.HandleGetVehicleByUUID)

	http.Handle("/", r)

	http.ListenAndServe(":8000", nil)
}
