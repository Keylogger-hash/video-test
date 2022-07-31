package mongo


import (
	"context"
	"api/v1/video/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)



type File struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
	Path string `bson:"path"`
	Filename string `bson:"filename"`
	Processing bool `bson:"processing"`
	Progress int `bson:"progress"`
	ProcessingStatus string `bson:"processingstatus"`
}

type FileRepository struct {
	db *mongo.Collection
}

func NewFileRepository(db *mongo.Database,collection string) *FileRepository{
	return &FileRepository{
		db: db.Collection(collection),
	}
}

func (f FileRepository) CreateFile(ctx context.Context, file *models.File) error{
	model := toModel(file)
	_, err  := f.db.InsertOne(ctx,model)
	if err != nil {
		return err
	}
	return nil	
}

func (f FileRepository) DeleteFile(ctx context.Context, id string) error{
	objId, _ := primitive.ObjectIDFromHex(id)
	_, err := f.db.DeleteOne(ctx,bson.M{"Id":objId})
	return err
}


func toModel(f *models.File) *File {
	uid,_ := primitive.ObjectIDFromHex(f.Id)
	return &File{
		Id: uid,
		Path: f.Path,
		Filename: f.Filename,
		Processing: f.Processing,
		Progress: f.Progress,
		ProcessingStatus: f.ProcessingStatus,
	}
}