package graph

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/777Lava/ozonTest/internal/entities"
)

var commentPublishedChannel = make(map[int]chan *entities.Comment)

// Replies is the resolver for the replies field.
func (r *commentResolver) Replies(ctx context.Context, obj *entities.Comment) ([]*entities.Comment, error) {
	if os.Getenv("MODE") == "inmemory" {
		var replies []*entities.Comment

		for _, comment := range r.comments {
			if comment.ParentID != nil && *comment.ParentID == obj.ID {
				replies = append(replies, comment)
			}
		}

		return replies, nil
	}

	comments, err := r.CommentRepo.GetReplies(obj.ID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CreatePost is the resolver for the createPost field.
func (r *mutationResolver) CreatePost(ctx context.Context, input entities.NewPost) (*entities.Post, error) {
	if os.Getenv("MODE") == "inmemory" {
		post := &entities.Post{
			ID:               len(r.posts) + 1,
			Author:           input.Author,
			Title:            input.Title,
			Content:          input.Content,
			CommentsDisabled: input.CommentsDisabled,
			CreatedAt:        time.Now().Local(),
		}
		r.posts = append(r.posts, post)
		return post, nil
	}

	post, err := r.PostRepo.CreatePost(input)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input entities.NewComment) (*entities.Comment, error) {
	if len(input.Content) > 2000 {
		return nil, errors.New("Comment must be less than 2000 characters")
	}

	if os.Getenv("MODE") == "inmemory" {
		if input.PostID > len(r.posts) || input.PostID < 1 {
			return nil, errors.New(fmt.Sprintf("Post with id=%d does not exists", input.PostID))
		}

		if r.posts[input.PostID-1].CommentsDisabled {
			return nil, errors.New(fmt.Sprintf("Post with id=%d have comments disabled", input.PostID))
		}

		if input.ParentID != nil && (*input.ParentID > len(r.comments) || *input.ParentID < 1) {
			return nil, errors.New(fmt.Sprintf("Comment with id=%d does not exists", *input.ParentID))
		}

		if input.ParentID != nil && r.comments[*input.ParentID-1].PostID != input.PostID {
			return nil, errors.New("The comment and the reply to it must be in one post")
		}

		comment := &entities.Comment{
			ID:        len(r.comments) + 1,
			Author:    input.Author,
			PostID:    input.PostID,
			ParentID:  input.ParentID,
			Content:   input.Content,
			CreatedAt: time.Now().Local(),
		}
		r.comments = append(r.comments, comment)

		if _, ok := commentPublishedChannel[comment.PostID]; ok {
			commentPublishedChannel[comment.PostID] <- comment
		}
		return comment, nil
	}

	comment, err := r.CommentRepo.CreateComment(input)
	if err != nil {
		return nil, err
	}
	if _, ok := commentPublishedChannel[comment.PostID]; ok {
		commentPublishedChannel[comment.PostID] <- comment
	}

	return comment, nil
}

// Comments is the resolver for the comments field.
func (r *postResolver) Comments(ctx context.Context, obj *entities.Post) ([]*entities.Comment, error) {
	if obj.CommentsDisabled {
		return nil, nil
	}

	if os.Getenv("MODE") == "inmemory" {
		var comments []*entities.Comment

		for _, comment := range r.comments {
			if obj.ID == comment.PostID {
				comments = append(comments, comment)
			}
		}
		return comments, nil
	}

	comments, err := r.CommentRepo.GetAllComments(obj.ID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetPosts is the resolver for the getPosts field.
func (r *queryResolver) GetPosts(ctx context.Context) ([]*entities.Post, error) {
	if os.Getenv("MODE") == "inmemory" {
		return r.posts, nil
	}

	posts, err := r.PostRepo.GetPosts()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// GetPost is the resolver for the getPost field.
func (r *queryResolver) GetPost(ctx context.Context, id int) (*entities.Post, error) {
	if os.Getenv("MODE") == "inmemory" {
		if id < 1 || id > len(r.posts) {
			return nil, errors.New(fmt.Sprintf("Post with id=%d does not exists", id))
		}
		return r.posts[id-1], nil
	}

	post, err := r.PostRepo.GetPost(id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

// GetComments is the resolver for the getComments field.
func (r *queryResolver) GetComments(ctx context.Context, postID int, first *int, skip *int) ([]*entities.Comment, error) {
	if os.Getenv("MODE") == "inmemory" {
		if postID < 1 || postID > len(r.posts) {
			return nil, errors.New(fmt.Sprintf("Post with id=%d does not exists", postID))
		}

		if r.posts[postID-1].CommentsDisabled {
			return nil, errors.New(fmt.Sprintf("Post with id=%d have comments disabled", postID))
		}

		var comments []*entities.Comment

		for _, comment := range r.comments {
			if comment.PostID == postID {
				comments = append(comments, comment)
			}
		}

		if *skip > len(comments) {
			return nil, nil
		}

		if *first+*skip > len(comments) {
			return comments[*skip:], nil
		}

		return comments[*skip : *skip+*first], nil
	}

	comments, err := r.CommentRepo.GetComments(postID, *first, *skip)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

// CommentAdded is the resolver for the commentAdded field.
func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID int) (<-chan *entities.Comment, error) {
	if os.Getenv("MODE") == "inmemory" {
		if postID < 1 || postID > len(r.posts) {
			return nil, errors.New(fmt.Sprintf("Post with id=%d does not exists", postID))
		}

		if r.posts[postID-1].CommentsDisabled {
			return nil, errors.New(fmt.Sprintf("Post with id=%d have comments disabled", postID))
		}
	} else {
		_, err := r.PostRepo.GetPost(postID)
		if err != nil {
			return nil, err
		}
	}

	commentPublishedChannel[postID] = make(chan *entities.Comment)
	return commentPublishedChannel[postID], nil
}

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Post returns PostResolver implementation.
func (r *Resolver) Post() PostResolver { return &postResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
