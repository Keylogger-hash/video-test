package main
import "github.com/gorilla/mux"
import "net/http"
import "os"
import "io"

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
	rw.Write([]byte("File upload success"))
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