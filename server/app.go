package server

import (
	"api/v1/video/files/repository"
	fmongo "api/v1/video/files/repository/mongo"
	fusecase "api/v1/video/files/repository/usecase"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
type App struct {
	httpServer *http.Server
	fileR repository.FileRepository
}


func NewApp() *App{
	db := InitDB()
	fileRepo := fmongo.NewFileRepository(db,viper.GetString("mongo.files_collection"))
	return &App{
		fileR: fusecase.NewFileUseCase(fileRepo),
	}
}

func (a *App) Run(port string) error {
	router := mux.NewRouter()
	a.httpServer = &http.Server{
		Addr: ":"+port,
		Handler: router,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 10,

	}
	go func(){
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatal("Failed to listen and serve: %v", err)
		}
	}()
	quit := make(chan os.Signal,1)
	signal.Notify(quit,os.Interrupt,os.Interrupt)
	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(),5 * time.Second)
	defer shutdown()
	return a.httpServer.Shutdown(ctx)
}

func InitDB() *mongo.Database{
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatal("Error occured while establishing mongo connect")
	}
	ctx, cancel := context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(),nil)
	if err != nil {
		log.Fatal(err)
	}
	return client.Database(viper.GetString("mongo.name"))
}