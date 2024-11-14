package db

import (
	"fmt"

	"github.com/777Lava/ozonTest/internal/entities"
	"gorm.io/gorm"
)

type CommentRepo struct {
	DB *gorm.DB
}

// CreateComment inserts a new comment into the database
func (r *CommentRepo) CreateComment(input entities.NewComment) (*entities.Comment, error) {
	postRepo := PostRepo{DB: r.DB}
	post, err := postRepo.GetPost(input.PostID)
	if err != nil {
		return nil, err
	}
	if post.CommentsDisabled {
		return nil, fmt.Errorf("comments to post with id = %d are disabled ", *input.ParentID)
	}
	if input.ParentID != nil {
		var parent *entities.Comment
		err := r.DB.Where("id = ?", *input.ParentID).Find(&parent).Error
		if err != nil {
			return nil, fmt.Errorf("comment with id = %d does not exist", *input.ParentID)
		}
		if input.PostID != parent.PostID {
			return nil, fmt.Errorf("comment with id = %d not in this post", *input.ParentID)
		}
	}
	comment := entities.Comment{Author: input.Author, ParentID: input.ParentID, PostID: input.PostID, Content: input.Content}
	r.DB.Create(&comment)
	return &comment, nil
}

// GetComments retrieves a paginated list of comments for a given post.
func (r *CommentRepo) GetComments(postID, first, skip int) ([]*entities.Comment, error) {
	postRepo := PostRepo{DB: r.DB}
	post, err := postRepo.GetPost(postID)
	if err != nil {
		return nil, err
	}
	if post.CommentsDisabled {
		return nil, fmt.Errorf("post with id = %d has comments disabled", postID)
	}

	var comments []*entities.Comment
	err = r.DB.Limit(first).Offset(skip).Where("post_id = ?", postID).Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("cannot get comments for post with id = %d", postID)
	}
	return comments, nil
}

// GetAllComments retrieves all comments for a specific post
func (r *CommentRepo) GetAllComments(postId int) ([]*entities.Comment, error) {
	var comments []*entities.Comment
	err := r.DB.Where("post_id = ?", postId).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetReplies retrieves all replies to a specific comment (identified by `parentId`)
func (r *CommentRepo) GetReplies(parentId int) ([]*entities.Comment, error) {
	var comments []*entities.Comment
	err := r.DB.Where("parent_id = ?", parentId).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}
