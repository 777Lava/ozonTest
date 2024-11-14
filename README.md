# ozonTest

## Start
1. Clone repository `git@github.com:777Lava/ozonTest.git`
2. Run `docker compose up -d` in root directory 
3. Go to http://localhost:8000
4. For stop docker run `docker compose stop` 

## P.S. Комментарии по проекту
Тесты запускаются сразу после сборки проекта.
Всю документацию можно посмотреть на http://localhost:8000

## P.S.S. Примеры запросов 
Просмотр всех постов с комментариями 
```
query {
  getPosts{
    id
    author
    title
    content
    comments {
      id
      author
      content
      createdAt
      replies
    }
    createdAt
  }
}
```
Просмотр поста по id

```
query {
  getPost(id: 1) {  # Подставьте нужный ID поста
    id
    title
    content
    author
    createdAt
    comments{
      content
      replies{
        content
      }
    }
  }
}
```

Создание поста

```
mutation { 
  createPost(
    input: {
      author: "Andrew",
      title: "My First Post",
      content: "This is the content of the post." 
      commentsDisabled: false
    }) { id 
    	author 
    	title 
    	content 
    	createdAt }
}
```
Создание комментария (parentId - комментарий на который хотим ответить)

```
mutation {
  createComment(input: {
    author: "alexey",
    postId: 1,
    content: "cool"
  }) {
    id
    author
    postID
    parentId
    content
    replies {
      id
      author
      content
    }
    createdAt
  }
}
```

