package db

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/daltbunker/soul_climbers/internal/database"
	"github.com/daltbunker/soul_climbers/types"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

var DB *database.Queries

func InitDatabase(dbURL string) {
	conn, err := sql.Open("postgres", dbURL)
	DB = database.New(conn)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
}

func NewUser(r *http.Request, user types.User) (types.User, error) {
	dbUser, err := DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return types.User{}, err 
	}

	newUser := types.User{
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Role: 	  dbUser.Role,
		SoulScore: dbUser.SoulScore,
	}

	return newUser, nil
}

func GetUserByEmail(r *http.Request, email string) (types.User, error) {
	dbUser, err := DB.GetUserByEmail(r.Context(), email)
	if err != nil {
		return types.User{}, err
	}

	user := types.User{
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Password: dbUser.Password,
		Role:	  dbUser.Role,
		SoulScore: dbUser.SoulScore,
	}

	return user, nil
}

func GetUserByUsername(r *http.Request, username string) (types.User, error) {
	dbUser, err := DB.GetUserByUsername(r.Context(), username)
	if err != nil {
		return types.User{}, err 
	}

	user := types.User{
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Password: dbUser.Password,
		Role: 	  dbUser.Role,
		SoulScore: dbUser.SoulScore,
	}

	return user, nil
}

func NewResetToken(r *http.Request, email string) (types.ResetPassword, error) {
	dbResetToken, err := DB.CreateResetToken(r.Context(), database.CreateResetTokenParams{
		Token:      uuid.New(),
		Expiration: time.Now().Add(time.Minute * 10).UTC(),
		Email:      email,
	})
	if err != nil {
		return types.ResetPassword{}, err
	}

	resetToken := types.ResetPassword{
		Token:      dbResetToken.Token,
		Expiration: dbResetToken.Expiration,
		Email:      dbResetToken.Email,
	}

	return resetToken, nil
}

func GetResetTokenByToken(r *http.Request, token uuid.UUID) (types.ResetPassword, error) {
	dbResetToken, err := DB.GetResetTokenByToken(r.Context(), token)
	if err != nil {
		return types.ResetPassword{}, err 
	}

	resetToken := types.ResetPassword{
		Token:      dbResetToken.Token,
		Expiration: dbResetToken.Expiration,
		Email:      dbResetToken.Email,
	}

	return resetToken, nil
}

func NewBlog(r *http.Request, blog types.Blog) (types.Blog, error) {
	user, err := DB.GetUserByUsername(r.Context(), blog.CreatedBy)
	if err != nil {
		return types.Blog{}, err
	}

	dbBlog, err := DB.CreateBlog(r.Context(), database.CreateBlogParams{
		Body:      []byte(blog.Body),
		Title:     blog.Title,
		Excerpt:   blog.Excerpt,
		IsPublished: blog.IsPublished,
		CreatedBy: user.UsersID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return types.Blog{}, err
	}

	newBlog := types.Blog{
		Id: dbBlog.BlogID,
		Body:      string(dbBlog.Body[:]),
		Title:     dbBlog.Title,
		Excerpt:   dbBlog.Excerpt,
		IsPublished: dbBlog.IsPublished,
		CreatedBy: user.Username,
		CreatedAt: dbBlog.CreatedAt.Format("02 Jan 2006"),
	}

	return newBlog, nil
}

func UpdateBlog(r *http.Request, blog types.Blog) (types.Blog, error) {

	dbBlog, err := DB.UpdateBlog(r.Context(), database.UpdateBlogParams{
		Title:     blog.Title,
		Body:      []byte(blog.Body),
		Excerpt:   blog.Excerpt,
		IsPublished: blog.IsPublished,
		UpdatedAt: time.Now().UTC(),
		BlogID: blog.Id,
	})
	if err != nil {
		return types.Blog{}, err
	}

	updatedBlog := types.Blog{
		Id: dbBlog.BlogID,
		Body:      string(dbBlog.Body[:]),
		Title:     dbBlog.Title,
		Excerpt:   dbBlog.Excerpt,
		IsPublished: dbBlog.IsPublished,
		CreatedBy: blog.CreatedBy,
		CreatedAt: dbBlog.CreatedAt.Format("02 Jan 2006"),
	}

	return updatedBlog, nil
}

func GetAllBlogs(r *http.Request) ([]types.Blog, error) {
	dbBlogs, err := DB.GetAllBlogs(r.Context())
	if err != nil {
		return []types.Blog{}, err
	}
	blogs := []types.Blog{}
	for _, b := range dbBlogs {
		blog := types.Blog{}
		blog.Id = b.BlogID
		blog.Title = b.Title
		blog.Excerpt = b.Excerpt
		blog.IsPublished = b.IsPublished
		blog.ImgName = b.ImgName.String
		blog.CreatedBy = b.Username
		blog.CreatedAt = b.CreatedAt.Format("02 Jan 2006")
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func GetBlogById(r *http.Request, id int32) (types.Blog, error) {
	dbBlog, err := DB.GetBlogById(r.Context(), id)
	if err != nil {
		return types.Blog{}, err
	}

	blog := types.Blog{}
	blog.Title = dbBlog.Title
	blog.Id = dbBlog.BlogID
	blog.Excerpt = dbBlog.Excerpt
	blog.Body = string(dbBlog.Body)
	blog.IsPublished = dbBlog.IsPublished
	blog.CreatedBy = dbBlog.Username
	blog.CreatedAt = dbBlog.CreatedAt.Format("02 Jan 2006")
	blog.ImgName = dbBlog.ImgName.String

	return blog, nil
}

func GetBlogByTitle(r *http.Request, title string) (types.Blog, error) {
	dbBlog, err := DB.GetBlogByTitle(r.Context(), title)
	if err != nil {
		return types.Blog{}, err
	}

	blog := types.Blog{}
	blog.Title = dbBlog.Title
	blog.Id = dbBlog.BlogID
	blog.Excerpt = dbBlog.Excerpt
	blog.Body = string(dbBlog.Body)
	blog.IsPublished = dbBlog.IsPublished

	return blog, nil
}

func NewBlogImg(r *http.Request, blogImg types.BlogImg) (types.BlogImg, error) {
	dbBlogImg, err := DB.CreateBlogImg(r.Context(), database.CreateBlogImgParams{
		ImgName: blogImg.ImgName,
		Img: blogImg.Img,
		BlogID: blogImg.BlogId,
		CreatedAt: time.Now().UTC(),	
		UpdatedAt: time.Now().UTC(),	
	})

	if err != nil {
		return types.BlogImg{}, err 
	}

	newBlogImg := types.BlogImg{
		ImgName: dbBlogImg.ImgName,
		Img: dbBlogImg.Img,
		BlogId: dbBlogImg.BlogID,
		CreatedAt: dbBlogImg.CreatedAt.Format("02 Jan 2006"),
		UpdatedAt: dbBlogImg.UpdatedAt.Format("02 Jan 2006"),
	}

	return newBlogImg, nil
}

func GetBlogImg(r *http.Request, id int32) (types.BlogImg, error) {
	dbBlogImg, err := DB.GetBlogImg(r.Context(), id)
	if err != nil {
		return types.BlogImg{}, err
	}

	blogImg := types.BlogImg{
		ImgName: dbBlogImg.ImgName,
		Img: dbBlogImg.Img,
		BlogId: dbBlogImg.BlogID,
	}

	return blogImg, nil
}

func DeleteBlogImg(r *http.Request, id int32) (types.BlogImg, error) {
	dbBlogImg, err := DB.DeleteBlogImg(r.Context(), id)
	if err != nil {
		return types.BlogImg{}, err
	}

	blogImg := types.BlogImg{
		ImgName: dbBlogImg.ImgName,
		Img: dbBlogImg.Img,
		BlogId: dbBlogImg.BlogID,
	}

	return blogImg, nil
}

func GetBlogsByCreator(r *http.Request, username string) ([]types.Blog, error) {
	user, err := DB.GetUserByUsername(r.Context(), username)
	if err != nil {
		return []types.Blog{}, err
	}

	dbBlogs, err := DB.GetBlogsByCreator(r.Context(), user.UsersID)
	if err != nil {
		return []types.Blog{}, err
	}

	blogs := []types.Blog{}
	for _, b := range dbBlogs {
		blog := types.Blog{}
		blog.Id = b.BlogID
		blog.Title = b.Title
		blog.Excerpt = b.Excerpt
		blog.IsPublished = b.IsPublished
		blog.CreatedBy = b.Username
		blog.CreatedAt = b.CreatedAt.Format("02 Jan 2006")
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func GetPlacementTestQuestions(r *http.Request) ([]types.Question, error) {
    dbQuestions, err := DB.GetTestQuestions(r.Context())
    if err != nil {
        return []types.Question{}, err
    }

    questions := []types.Question{}
    for _, q := range dbQuestions {
        question := types.Question{}
        question.Id = strconv.Itoa(int(q.TestQuestionID))
        question.InputType = q.InputType
        question.Text = q.QuestionText
        question.Answers = strings.Split(q.Answers, ",")
        question.AnswerPoints = strings.Split(q.AnswerPoints, ",")
		question.PossiblePoints = int(q.Points)
        questions = append(questions, question)
    }

    return questions, nil
}

func InsertPlacementTest(r *http.Request, username string, score int32) error {
	user, err := DB.GetUserByUsername(r.Context(), username)
	if err != nil {
		return err
	}

	_, err = DB.InsertPlacmentTest(r.Context(), database.InsertPlacmentTestParams{
		UserID: user.UsersID,
		Score: score,
		CreatedAt: time.Now().UTC(),	
		UpdatedAt: time.Now().UTC(),	
	})

	return err
}

func GetPlacementTest(r *http.Request, username string) (int32, error) {
	user, err := DB.GetUserByUsername(r.Context(), username)
	if err != nil {
		return 0, err
	}

	dbPlacementTest, err := DB.GetPlacementTest(r.Context(), user.UsersID)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return dbPlacementTest, err
}

func SetUserSoulScore(r *http.Request, username string, soulScore int32) (int32, error) {
	dbUserSoulScore, err := DB.SetUserSoulScore(r.Context(), database.SetUserSoulScoreParams{
		SoulScore: soulScore,
		Username: username,
	})
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return dbUserSoulScore.SoulScore, err
}

func DeleteBlog(r *http.Request, id int32) error {
	_, err := DB.DeleteBlog(r.Context(), id)
	return err
}