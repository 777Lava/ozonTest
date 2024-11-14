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

func TestPostRep(t *testing.T) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("can't find .env file. continue work using system env vars")
	}


	line := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("TEST_DB_NAME"), os.Getenv("DB_PORT"))

	DB, err := gorm.Open(postgres.Open(line), &gorm.Config{})
	if err != nil {
		t.Errorf("faile to connect to DB: %v", err)
	}

	r := db.PostRepo{DB: DB}
	
	t.Run("create and return post", func(t *testing.T) {
		DB.Migrator().CreateTable(&entities.Post{})
		defer DB.Migrator().DropTable(&entities.Post{})

		want := entities.Post{
			ID: 1,
			Title: "title",
			Author: "author",
			Content: "content",
			CommentsDisabled: false,
		}
		have, err := r.CreatePost(entities.NewPost{
			Title: want.Title,
			Author: want.Author,
			Content: want.Content,
			CommentsDisabled: want.CommentsDisabled,
		})
		require.Equal(t,err, nil)
		require.Equal(t, want.ID, have.ID)
		require.Equal(t, want.Author, have.Author)
		require.Equal(t, want.Title, have.Title)
		require.Equal(t, want.Content, have.Content)
		require.Equal(t, want.CommentsDisabled, have.CommentsDisabled)
	})
	t.Run("check current post", func(t *testing.T) {
		DB.Migrator().CreateTable(&entities.Post{})
		defer DB.Migrator().DropTable(&entities.Post{})

		want := entities.Post{
			ID:               1,
			Author:           "author 1",
			Title:            "title 1",
			Content:          "content 1",
			CommentsDisabled: false,
		}
		DB.Create(&want)
		get, err := r.GetPost(want.ID)

		require.Equal(t, nil, err)
		require.Equal(t, want.ID, get.ID)
		require.Equal(t, want.Author, get.Author)
		require.Equal(t, want.Title, get.Title)
		require.Equal(t, want.Content, get.Content)
		require.Equal(t, want.CommentsDisabled, get.CommentsDisabled)
	})

	t.Run("check return all posts", func(t *testing.T) {
		DB.Migrator().CreateTable(&entities.Post{})
		defer DB.Migrator().DropTable(&entities.Post{})

		want := []entities.Post{
			{
				ID:               1,
				Author:           "author 1",
				Title:            "title 1",
				Content:          "content 1",
				CommentsDisabled: false,
			},
			{
				ID:               2,
				Author:           "author 2",
				Title:            "titile 2",
				Content:          "content 2",
				CommentsDisabled: true,
			},
		}
		DB.Create(&want)
		get, err := r.GetPosts()

		for i := range want {
			require.Equal(t, nil, err)
			require.Equal(t, want[i].ID, get[i].ID)
			require.Equal(t, want[i].Author, get[i].Author)
			require.Equal(t, want[i].Title, get[i].Title)
			require.Equal(t, want[i].Content, get[i].Content)
			require.Equal(t, want[i].CommentsDisabled, get[i].CommentsDisabled)
		}
	})
}