package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" validate:"omitempty"`
	Title     string             `bson:"title" validate:"required"`
	Desc      string             `bson:"description" validate:""`
	Content   string             `bson:"content" validate:""`
	Tags      []string           `bson:"tags" validate:""`
	CreatedAt time.Time          `bson:"created_at" validate:""`
	Slug      string             `bson:"slug" validate:""`
	Published bool               `bson:"published" validate:""`
}

func (b *BlogModel) String() string {
	return fmt.Sprintf(`
Title: %s
Desc: %s
Created At: %s
Published: %t
Tags: %v
----------------

%s

================


`, b.Title, b.Desc, b.CreatedAt, b.Published, b.Tags, b.Content)
}

type Blogs struct {
	collection *mongo.Collection
}

func NewBlogs(collection *mongo.Collection) *Blogs {
	return &Blogs{
		collection: collection,
	}
}

func (b *Blogs) CreateSlug(title string) string {

	base := slug.Make(title)
	unique := base
	counter := 1

	for {
		if !b.CheckBlogExistsBySlug(unique) {
			break
		}

		counter++
		unique = fmt.Sprintf("%s-%d", unique, counter)
	}
	return unique
}

// CREATE
func (b *Blogs) AddBlog(blog BlogModel) error {
	validate := validator.New()

	if err := validate.Struct(blog); err != nil {
		return errors.Join(errors.New("failed to validate blog"), err)
	}

	blog.ID = primitive.NewObjectID()
	blog.CreatedAt = time.Now()
	blog.Published = false
	blog.Slug = b.CreateSlug(blog.Title)

	if _, err := b.collection.InsertOne(context.TODO(), blog); err != nil {
		return errors.Join(errors.New("failed to add blog to the database"), err)
	}

	return nil
}

// DELETE
func (b *Blogs) DeleteBlogByID(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Join(errors.New("invalid object id provided"), err)
	}
	if _, err := b.collection.DeleteOne(context.TODO(), bson.M{
		"_id": objectId,
	}); err != nil {
		return errors.Join(errors.New("failed to delete blog"), err)
	}
	return nil
}

// READ
func (b *Blogs) FindBlogs(filter bson.M, pageNo, blogsPerPage int) ([]*BlogModel, error) {
	limit := blogsPerPage
	skip := (pageNo - 1) * limit

	var blogs []*BlogModel
	options := options.Find().SetSort(bson.M{"created_at": -1}).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := b.collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, errors.Join(errors.New("failed to find blogs"), err)
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &blogs); err != nil {
		return nil, errors.Join(errors.New("failed to find blogs"), err)
	}

	return blogs, nil
}

func (b *Blogs) FindBlogByID(id string) (*BlogModel, error) {

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.Join(errors.New("invalid object id provided"), err)
	}

	var blog BlogModel
	if err := b.collection.FindOne(context.TODO(), bson.M{
		"_id": objectId,
	}).Decode(&blog); err != nil {
		return nil, errors.Join(errors.New("failed to find blog by id"), err)
	}
	return &blog, nil
}

func (b *Blogs) FindBlogBySlug(slug string) (*BlogModel, error) {
	var blog BlogModel
	if err := b.collection.FindOne(context.TODO(), bson.M{
		"slug": slug,
	}).Decode(&blog); err != nil {
		return nil, errors.Join(errors.New("failed to find blog by slug"), err)
	}
	return &blog, nil
}

func (b *Blogs) FindBlogsByTitle(title string, pageNo, blogsPerPage int) ([](*BlogModel), error) {
	blogs, err := b.FindBlogs(bson.M{
		"title": bson.M{
			"$regex":   title,
			"$options": "i",
		},
	}, pageNo, blogsPerPage)

	return blogs, err
}

func (b *Blogs) CheckBlogExistsByID(id string) bool {
	if found, err := b.FindBlogByID(id); found == nil || err != nil {
		return false
	}
	return true
}

func (b *Blogs) CheckBlogExistsBySlug(slug string) bool {
	if found, err := b.FindBlogBySlug(slug); found == nil || err != nil {
		return false
	}
	return true
}

// UPDATE
func (b *Blogs) PublishDraftBlogByID(id string) error {
	if !b.CheckBlogExistsByID(id) {
		return errors.New("blog doesn't exists")
	}

	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err := b.collection.UpdateByID(context.TODO(), objectId, bson.M{
		"$set": bson.M{
			"published": bson.M{
				"$ne": true,
			},
		},
	})
	if err != nil {
		return errors.Join(errors.New("failed to publish blog"), err)
	}
	return nil
}

func (b *Blogs) UpdateBlogByID(id string, body BlogModel) error {

	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return errors.Join(errors.New("failed to validate blog"), err)
	}

	old, err := b.FindBlogByID(id)
	if err != nil {
		return errors.Join(errors.New("failed to update blog"), err)
	}

	updatedBlog := BlogModel{}
	updatedBlog.ID = old.ID
	updatedBlog.Published = old.Published
	updatedBlog.CreatedAt = old.CreatedAt
	updatedBlog.Slug = old.Slug

	updatedBlog.Content = body.Content
	updatedBlog.Tags = body.Tags
	updatedBlog.Title = body.Title
	updatedBlog.Desc = body.Desc
	if updatedBlog.Title != old.Title {
		updatedBlog.Slug = b.CreateSlug(updatedBlog.Title)
	}

	if _, err := b.collection.UpdateByID(context.TODO(), updatedBlog.ID, bson.M{
		"$set": updatedBlog,
	}); err != nil {
		return errors.Join(errors.New("failed to update blog"), err)
	}

	return nil
}
