package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
