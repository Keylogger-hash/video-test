package http

import (
	"api/v1/video/files/delivery/http/structs"
	fusecase "api/v1/video/files/usecase"
	"context"
	"io"
	"net/http"
	"os"
	"github.com/mailru/easyjson"

	"github.com/spf13/viper"
)

const MAX_UPLOAD_SIZE int64=1024*1024*1024


type Handler struct {
	fC fusecase.FileUseCase
}


func NewHandler(f fusecase.FileUseCase) *Handler{
	return &Handler{
		fC: f,
	}
}


func (h *Handler) Create(rw http.ResponseWriter,r *http.Request){
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)	
	if err != nil {
		http.Error(rw,"Too much file",400)
	}
	file,header, err := r.FormFile("file")
	if err != nil {
		http.Error(rw,"File not found in post request",400)
	}
	defer file.Close()
	mediaFolder := viper.GetString("media")
	filePath := mediaFolder+"/"+header.Filename
	dst,err := os.Create(filePath)
	if err != nil {
		http.Error(rw,"Can't create filename",500)
	}
	defer dst.Close()
	if _,err := io.Copy(dst,file); err != nil{
		http.Error(rw,"Can't upload file", 500)
	}
	ctx := context.Background()
	h.fC.CreateFile(ctx,filePath,header.Filename)
    
	resp := &structs.Response{Id:id}
	rawBytes, err := easyjson.Marshal(resp)
	rw.Write(rawBytes)
}

