package db

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/daltbunker/soul_climbers/internal/database"
	"github.com/daltbunker/soul_climbers/types"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

var DB *database.Queries

func InitDatabase(dbURL string) error {
	conn, err := sql.Open("postgres", dbURL)
	DB = database.New(conn)
	return err
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
		CreatedBy: user.UsersID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return types.Blog{}, err
	}

	newBlog := types.Blog{
		Body:      string(dbBlog.Body[:]),
		Title:     dbBlog.Title,
		CreatedBy: user.Username,
		CreatedAt: dbBlog.CreatedAt.Format("02 Jan 2006"),
	}

	return newBlog, nil
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
		blog.Body = string(b.Body)
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
	blog.Body = string(dbBlog.Body)
	blog.CreatedBy = dbBlog.Username
	blog.CreatedAt = dbBlog.CreatedAt.Format("02 Jan 2006")
	return blog, nil
}
