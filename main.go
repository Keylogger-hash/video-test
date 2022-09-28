package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var schema string = `CREATE TABLE IF NOT EXISTS video_test(id integer primary key autoincrement,
	uuid varchar(255) not null unique,
	filepath varchar(255),
	filename varchar(255),
	processing bool default false,
	processingSuccess bool default false);
`

// Video
type Video struct {
	Id                int    `db:"id" json:"id" xml:"id"`
	UUID              string `db:"uuid" json:"uuid" xml:"uuid"`
	Filepath          string `db:"filepath" json:"filepath" xml:"filepath"`
	Filename          string `db:"filename" json:"filename" xml:"filename"`
	Processing        bool   `db:"processing" json:"processing" xml:"processing"`
	ProcessingSuccess bool   `db:"processingSuccess" json:"processingSuccess" xml:"processingSuccess"`
}


type VideoWidthHeight struct {
	Width int `db:"Width" json:"Width" xml:"Width`
	Height int  `db:"Height" json:"Height" xml:"Height`
}
type Error struct {
	Error string `json:"Error" xml:"Error"`
}

type Response struct {
	Success bool     `json:"Status" xml:"Status"`
	Data    []*Video `json:"Data,omitempty" xml:"Data,omitempty"`
	Errors  []*Error `json:"Errors,omitempty" xml:"Errors,omitempty"`
}

type ResponseDetail struct {
	Success bool     `json:"Status" xml:"Status"`
	Video   *Video   `json:"Data,omitempty" xml:"Data,omitempty"`
	Errors  []*Error `json:"Errors,omitempty" xml:"Errors,omitempty"`
}

func resizeVideo(db *sqlx.DB,wg *sync.WaitGroup,videofilepath string, filename string,uuid string,width int, height int) {
	defer wg.Done()
	out_filepath := filepath.Join(videofilepath, "out.mp4")

	input_filepath := filepath.Join(videofilepath, filename)
	strWidth := strconv.Itoa(width)
	strHeight := strconv.Itoa(height)
	cmd := exec.Command("ffmpeg", "-i", input_filepath, "-y", "-vf",fmt.Sprintf("scale=%s:%s",strWidth,strHeight), out_filepath)
	fmt.Println(cmd)
	fmt.Println(input_filepath)
	fmt.Println(out_filepath)
	bytes, err := cmd.CombinedOutput()
	fmt.Println("Running ffmpeg task", string(bytes))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	wg.Add(1)
	go UpdateProcessingQuery(db,wg,uuid,false,true)

}

func GetVideoQuery(db *sqlx.DB, uuid string, videoChan chan *Video, wg *sync.WaitGroup) {
	defer wg.Done()
	video := &Video{}
	err := db.Get(video, "SELECT * FROM video_test WHERE uuid=$1", uuid)
	if err == sql.ErrNoRows {
		videoChan <- nil
		return
	}
	videoChan <- video
	return
}
func GetAllVideoQuery(db *sqlx.DB, videoChan chan []*Video, wg *sync.WaitGroup) {
	defer wg.Done()
	videos := []*Video{}
	err := db.Select(&videos, "SELECT * FROM video_test;")
	if err != nil {
		panic(err)
	}
	videoChan <- videos
}
func UpdateProcessingQuery(db *sqlx.DB, wg *sync.WaitGroup,uuid string, processing bool, processingSuccess bool) {
	defer wg.Done()
	tx := db.MustBegin()

	_ = tx.MustExec("UPDATE video_test SET processing=$1,processingSuccess=$2 WHERE uuid=$3 ", processing,processingSuccess, uuid)
	tx.Commit()
}

func DeleteProcessingQuery(db *sqlx.DB, wg *sync.WaitGroup,uuid string ) {
	defer wg.Done()
	tx := db.MustBegin()

	_ = tx.MustExec("DELETE FROM video_test WHERE uuid=$1 ", uuid)
	tx.Commit()

}

func main() {
	e := echo.New()
	db, err := sqlx.Connect("sqlite3", "db.sqlite3")
	db.Exec(schema)
	var wg sync.WaitGroup
	if err != nil {
		e.Logger.Fatal(err)
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/files", func(c echo.Context) error {
		videoChan := make(chan ([]*Video), 1)
		wg.Add(1)
		go GetAllVideoQuery(db, videoChan, &wg)
		videos := <-videoChan
		resp := &Response{
			Success: true,
			Data:    videos,
		}
		return c.JSON(http.StatusOK, resp)
	})
	e.PATCH("/file/:id/patch", func(c echo.Context) error {
		uuid := c.Param("id")
		videowidthheight := &VideoWidthHeight{}
		if err:=c.Bind(videowidthheight); err != nil {
			resp := &ResponseDetail{
				Success: false,
				Errors: []*Error{
					{Error: fmt.Sprintf("%v",err)},
				},
			}
			return c.JSON(http.StatusNotFound, resp)
		}
		if videowidthheight.Height<=20 || videowidthheight.Width<=20{
			resp := &ResponseDetail{
				Success: false,
				Errors: []*Error{
					{Error: "Height and width can be less than 20"},
				},
			}
			return c.JSON(http.StatusNotFound, resp)

		}
		videoChan := make(chan *Video, 1)
		wg.Add(1)
		go GetVideoQuery(db, uuid, videoChan, &wg)
		video := <-videoChan
		videopath := filepath.Join(video.Filepath, video.Filename)
		_, err := os.Stat(videopath)
		if video == nil || err != nil {
			resp := &ResponseDetail{
				Success: false,
				Errors: []*Error{
					{Error: "Video not found"},
				},
			}
			return c.JSON(http.StatusNotFound, resp)
		} else {
			wg.Add(1)
			go UpdateProcessingQuery(db, &wg,uuid,true,false)
			wg.Add(1)
			go resizeVideo(db,&wg,video.Filepath, video.Filename,uuid,videowidthheight.Width,videowidthheight.Height)
			resp := &ResponseDetail{
				Success: true,
			}
			return c.JSON(http.StatusOK, resp)
		}
	})
	e.DELETE("/file/:id/delete", func(c echo.Context) error {
		uuid := c.Param("id")
		videoChan := make(chan *Video, 1)
		wg.Add(1)
		go GetVideoQuery(db, uuid, videoChan, &wg)
		video := <-videoChan
		videofilepath := filepath.Join(video.Filepath)
		_,err := os.Stat(videofilepath)
		if video == nil || err != nil{
			resp := &ResponseDetail{
				Success: true,
				Errors: *&[]*Error{
					{Error: "Video not found"},
				},
				
			}
			return c.JSON(http.StatusNotFound,resp)
		} else {
			os.RemoveAll(videofilepath)
			wg.Add(1)
			go DeleteProcessingQuery(db,&wg,uuid)
			resp := &ResponseDetail{
				Success: true,
			}
			return c.JSON(http.StatusOK, resp)
		}
		

	})
	e.GET("/file/:id", func(c echo.Context) error {
		if err != nil {
			e.Logger.Fatal(err)
		}
		// video := &Video{}
		uuid := c.Param("id")
		videoChan := make(chan *Video, 1)
		wg.Add(1)
		go GetVideoQuery(db, uuid, videoChan, &wg)
		video := <-videoChan
		if video == nil {
			errorVideoNotFound := &Error{Error: "Video not found"}
			resp := &ResponseDetail{
				Success: false,
				Errors:  []*Error{errorVideoNotFound},
			}
			return c.JSON(http.StatusNotFound, resp)
		}

		resp := &ResponseDetail{
			Success: true,
			Video:   video,
		}
		return c.JSON(http.StatusOK, resp)

	})
	e.POST("/file/add", func(c echo.Context) error {
		// name := c.FormValue("name")
		file, err := c.FormFile("file")
		if err != nil {
			e.Logger.Fatal(err)
		}
		if file.Header.Get("Content-Type") != "video/mp4" {
			e := &Error{
				Error: "Content-type should be mp4",
			}
			resp := &Response{
				Success: false,
				Errors:  []*Error{e},
			}
			return c.JSON(http.StatusOK, resp)
		}
		src, err := file.Open()

		defer src.Close()
		if err != nil {
			e.Logger.Fatal(err)
		}
		pathUUID := uuid.New().String()
		uuidFilepath := filepath.Join("upload", pathUUID)
		os.Mkdir(uuidFilepath, 0777)
		path := filepath.Join(uuidFilepath, file.Filename)
		dst, err := os.Create(path)
		defer dst.Close()
		if err != nil {
			e.Logger.Fatal(err)
		}
		if _, err := io.Copy(dst, src); err != nil {
			e.Logger.Fatal(err)
		}
		db, err := sqlx.Connect("sqlite3", "db.sqlite3")
		if err != nil {
			e.Logger.Fatal(err)
		}
		tx := db.MustBegin()
		if err != nil {
			e.Logger.Fatal(err)
		}
		dbUUID := uuid.New().String()
		result := tx.MustExec("INSERT INTO video_test(uuid,filepath,filename) VALUES($1,$2,$3)", dbUUID, uuidFilepath, file.Filename)
		tx.Commit()
		lastId, err := result.LastInsertId()
		if err != nil {
			e.Logger.Fatal(err)
		}
		video := &Video{Id: int(lastId), UUID: dbUUID, Filepath: uuidFilepath, Filename: file.Filename}

		resp := &ResponseDetail{
			Success: true,
			Video:   video,
		}
		return c.JSON(200, resp)

	})
	runtime.GOMAXPROCS(8)
	e.Logger.Fatal(e.Start(":1323"))
	wg.Wait()

}
