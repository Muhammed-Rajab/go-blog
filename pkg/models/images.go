package models

import (
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImageModel struct {
	ID       primitive.ObjectID `bson:"_id" form:"_id"`
	Caption  string             `bson:"caption" form:"caption"`
	Location string             `bson:"location" form:"location"`
	Slug     string             `bson:"slug" form:"slug"`
	AddedAt  time.Time          `bson:"added_at" form:"added_at"`
}

type Images struct {
	collection *mongo.Collection
}

func NewImages(collection *mongo.Collection) Images {
	return Images{collection: collection}
}

func (i *Images) CreateSlug(caption string) string {
	base := slug.Make(caption)
	unique := base
	counter := 1

	for {
		if !i.CheckImageExistsBySlug(unique) {
			break
		}

		counter++
		unique = fmt.Sprintf("%s-%d", unique, counter)
	}
	return unique
}

func (i *Images) AddImage(image ImageModel) (primitive.ObjectID, error) {

	image.AddedAt = time.Now()
	image.ID = primitive.NewObjectID()

	res, err := i.collection.InsertOne(context.TODO(), image)
	if err != nil {
		return primitive.NilObjectID, errors.Join(errors.New("failed to add image to the database"), err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		return oid, nil
	}
	return primitive.NilObjectID, nil
}

func (i *Images) FindImageByID(id string) (*ImageModel, error) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Join(errors.New("invalid object id provided"), err)
	}

	var image ImageModel
	if err := i.collection.FindOne(context.TODO(), bson.M{
		"_id": objectId,
	}).Decode(&image); err != nil {
		return nil, errors.Join(errors.New("failed to find image by id"), err)
	}
	return &image, nil
}

func (i *Images) FindImageBySlug(slug string) (*ImageModel, error) {
	var image ImageModel
	if err := i.collection.FindOne(context.TODO(), bson.M{
		"slug": slug,
	}).Decode(&image); err != nil {
		return nil, errors.Join(errors.New("failed to find image by slug"), err)
	}
	return &image, nil
}

func (i *Images) CheckImageExistsBySlug(slug string) bool {
	if found, err := i.FindImageBySlug(slug); found == nil || err != nil {
		return false
	}
	return true
}

func (i *Images) FindImages(filter bson.M, pageNo, imagesPerPage int) ([]*ImageModel, error) {
	limit := imagesPerPage
	skip := (pageNo - 1) * limit

	var images []*ImageModel
	options := options.Find().SetSort(bson.M{"created_at": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := i.collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, errors.Join(errors.New("failed to find images"), err)
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &images); err != nil {
		return nil, errors.Join(errors.New("failed to find images"), err)
	}

	return images, nil
}

func (i *Images) DeleteImageByID(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Join(errors.New("invalid object id provided"), err)
	}

	// Remove the object from the server
	img, err := i.FindImageByID(id)
	if err != nil {
		return err
	}

	// Remove the file if it exists
	if _, err := os.Stat(img.Location); err == nil {
		err := os.Remove(img.Location)
		if err != nil {
			return errors.New("failed to remove error: " + err.Error())
		}
	}
	// ! Commenting this for now cause im not sure
	// ! If I want this feature now.
	// } else if os.IsNotExist(err) {
	// 	return errors.New("file doesn't exist")
	// } else {
	// 	return err
	// }

	if _, err := i.collection.DeleteOne(context.TODO(), bson.M{
		"_id": objectId,
	}); err != nil {
		return errors.Join(errors.New("failed to delete image"), err)
	}
	return nil
}

func (i *Images) ValidateImage(file multipart.File) bool {
	reader := io.TeeReader(file, os.Stdout)
	_, _, err := image.Decode(reader)
	log.Print(err)
	return err == nil
}
