package http

import "github.com/gorilla/mux"
import 	fusecase "api/v1/video/files/usecase"
import "api/v1/video/files/delivery/http"


func RegisterHTTPEndpoints(router *mux.Router,fc fusecase.FileUseCase){
	h := NewHandler(fc)

	r.HandleFunc("/file",h.Create).Methods("POST")
	r.HandleFunc("/file/{id}",HandleDeleteFile).Methods("DELETE")
	r.HandleFunc("/file/{id}",HandleGetInfoFile).Methods("GET")
	r.HandleFunc("/file/{id}",HandlePatchFile).Methods("PATCH")

}