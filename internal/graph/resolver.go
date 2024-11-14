package graph

import (
	db "github.com/777Lava/ozonTest/internal/database"
	"github.com/777Lava/ozonTest/internal/entities"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	posts []*entities.Post
	comments []*entities.Comment
	PostRepo db.PostRepo
	CommentRepo db.CommentRepo
}
