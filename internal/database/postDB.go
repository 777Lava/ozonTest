package db

import (
	"github.com/777Lava/ozonTest/internal/entities"
	"gorm.io/gorm"
)

type PostRepo struct{ 
	DB *gorm.DB
}

// CreatePost inserts a new post into the database
func (r *PostRepo) CreatePost(input entities.NewPost) (*entities.Post, error){
	post := entities.Post{Author: input.Author, Title: input.Title, Content: input.Content}
	err := r.DB.Create(&post)
	if err.Error != nil { 
		return nil, err.Error
	}
	return &post, nil 
}

// GetPosts retrieves all posts from the database
func (r *PostRepo) GetPosts() ([]*entities.Post, error){
	var posts []*entities.Post
	err := r.DB.Find(&posts)
	if err.Error != nil { 
		return nil, err.Error
	}
	return posts, err.Error
}

// GetPost retrieves a specific post by its ID
func (r *PostRepo) GetPost(id int) (*entities.Post, error){
	var post *entities.Post
	err := r.DB.Where("id = ?", id).Find(&post)
	if err.Error != nil {
		return nil, err.Error
	}
	return post, err.Error
}