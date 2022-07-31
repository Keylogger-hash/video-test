package main

import (
	"api/v1/video/structs"
	"io"
	"net/http"
	"os"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"api/v1/video/files/repository/mongo"
)

const MAX_UPLOAD_SIZE int64=1024*1024*1024



func HandleFile(rw http.ResponseWriter,r *http.Request) {
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)	
	if err != nil {
		http.Error(rw,"Too much file",400)
	}
	file,header, err := r.FormFile("file")
	if err != nil {
		http.Error(rw,"File not found in post request",400)
	}
	defer file.Close()
	dst,err := os.Create(header.Filename)
	if err != nil {
		http.Error(rw,"Can't create filename",500)
	}
	defer dst.Close()
	if _,err := io.Copy(dst,file); err != nil{
		http.Error(rw,"Can't upload file", 500)
	}
	id := uuid.New().String()
    
	resp := &structs.Response{Id:id}
	rawBytes, err := easyjson.Marshal(resp)
	rw.Write(rawBytes)
}
func HandleDeleteFile(rw http.ResponseWriter,r *http.Request) {
	rw.Write([]byte("DELETE"))
}
func HandleGetInfoFile(rw http.ResponseWriter,r *http.Request) {
	rw.Write([]byte("GET INFO"))

}
func HandlePatchFile(rw http.ResponseWriter,r *http.Request) {
	rw.Write([]byte("PATCH"))
}

func main(){
	r := mux.NewRouter()
	r.HandleFunc("/file",HandleFile).Methods("POST")
	r.HandleFunc("/file/{id}",HandleDeleteFile).Methods("DELETE")
	r.HandleFunc("/file/{id}",HandleGetInfoFile).Methods("GET")
	r.HandleFunc("/file/{id}",HandlePatchFile).Methods("PATCH")
	http.ListenAndServe(":8080",r)
	
}