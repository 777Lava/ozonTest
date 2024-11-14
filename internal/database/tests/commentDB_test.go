package tests

import (
	"fmt"
	"log"
	"os"

	"testing"

	db "github.com/777Lava/ozonTest/internal/database"
	"github.com/777Lava/ozonTest/internal/entities"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func TestCommentReop(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("can't find .env file. continue work using system env vars ")
	}

	line := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("TEST_DB_NAME"), os.Getenv("DB_PORT"))
	DB, err := gorm.Open(postgres.Open(line), &gorm.Config{})
	if err != nil {
		t.Errorf("faile to connect to DB: %v", err)
	}
	DB.Migrator().CreateTable(&entities.Post{})
	defer DB.Migrator().DropTable(&entities.Post{})

	p := entities.Post{
		ID:               1,
		Author:           "author 1",
		Title:            "title 1",
		Content:          "content 1",
		CommentsDisabled: false,
	}
	DB.Create(&p)

	cr := db.CommentRepo{DB: DB}

	t.Run("should create and return Comment", func(t *testing.T) {
		DB.Migrator().CreateTable(&entities.Comment{})
		defer DB.Migrator().DropTable(&entities.Comment{})

		want := entities.Comment{
			ID:       1,
			Author:   "test",
			PostID:   1,
			ParentID: nil,
			Content:  "test",
		}
		get, err := cr.CreateComment(entities.NewComment{
			Author:   "test",
			PostID:   1,
			ParentID: nil,
			Content:  "test",
		})

		require.Equal(t, nil, err)
		require.Equal(t, want.ID, get.ID)
		require.Equal(t, want.Author, get.Author)
		require.Equal(t, want.PostID, get.PostID)
		require.Equal(t, want.ParentID, get.ParentID)
		require.Equal(t, want.Content, get.Content)
	})

	t.Run("should return Replies", func(t *testing.T) {
		DB.Migrator().CreateTable(&entities.Comment{})
		defer DB.Migrator().DropTable(&entities.Comment{})

		c := entities.Comment{
			ID:       1,
			Author:   "test",
			PostID:   1,
			ParentID: nil,
			Content:  "test",
		}
		DB.Create(&c)

		ptr := func(i int) *int { return &i }(1)

		want := []entities.Comment{
			{
				ID:       1,
				Author:   "test",
				PostID:   1,
				ParentID: ptr,
				Content:  "test",
			},
			{
				ID:       2,
				Author:   "test",
				PostID:   1,
				ParentID: ptr,
				Content:  "test",
			},
		}
		DB.Create(&want)
		
		get, err := cr.GetReplies(p.ID)

		require.Equal(t, nil, err)
		for i := range get {
			require.Equal(t, want[i].ID, get[i].ID)
			require.Equal(t, want[i].Author, get[i].Author)
			require.Equal(t, want[i].PostID, get[i].PostID)
			require.Equal(t, want[i].ParentID, get[i].ParentID)
			require.Equal(t, want[i].Content, get[i].Content)
		}
	})

	t.Run("should return all Comments by postID", func(t *testing.T) {
		DB.Migrator().CreateTable(&entities.Comment{})
		defer DB.Migrator().DropTable(&entities.Comment{})

		want := []entities.Comment{
			{
				ID:       1,
				Author:   "test",
				PostID:   1,
				ParentID: nil,
				Content:  "test",
			},
			{
				ID:       2,
				Author:   "test",
				PostID:   1,
				ParentID: nil,
				Content:  "test",
			},
		}
		DB.Create(&want)
		get, err := cr.GetAllComments(p.ID)

		require.Equal(t, nil, err)
		for i := range get {
			require.Equal(t, want[i].ID, get[i].ID)
			require.Equal(t, want[i].Author, get[i].Author)
			require.Equal(t, want[i].PostID, get[i].PostID)
			require.Equal(t, want[i].ParentID, get[i].ParentID)
			require.Equal(t, want[i].Content, get[i].Content)
		}
	})

	t.Run("should return Comments", func(t *testing.T) {
		DB.Migrator().CreateTable(&entities.Comment{})
		defer DB.Migrator().DropTable(&entities.Comment{})

		want := []entities.Comment{
			{
				ID:       1,
				Author:   "test",
				PostID:   1,
				ParentID: nil,
				Content:  "test",
			},
			{
				ID:       2,
				Author:   "test",
				PostID:   1,
				ParentID: nil,
				Content:  "test",
			},
		}
		DB.Create(&want)
		get, err := cr.GetComments(p.ID, 1, 0)

		require.Equal(t, nil, err)
		require.Equal(t, want[0].ID, get[0].ID)
		require.Equal(t, want[0].Author, get[0].Author)
		require.Equal(t, want[0].PostID, get[0].PostID)
		require.Equal(t, want[0].ParentID, get[0].ParentID)
		require.Equal(t, want[0].Content, get[0].Content)
	})
	
}