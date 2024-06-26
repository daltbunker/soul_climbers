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

	newUser := types.User{
		Username: dbUser.Username,
		Email:    dbUser.Email,
	}

	return newUser, err
}

func GetUserByEmail(r *http.Request, email string) (types.User, error) {
	dbUser, err := DB.GetUserByEmail(r.Context(), email)

	user := types.User{
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Password: dbUser.Password,
	}

	return user, err
}

func GetUserByUsername(r *http.Request, username string) (types.User, error) {
	dbUser, err := DB.GetUserByUsername(r.Context(), username)

	user := types.User{
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Password: dbUser.Password,
	}

	return user, err
}

func NewResetToken(r *http.Request, email string) (types.ResetPassword, error) {
	dbResetToken, err := DB.CreateResetToken(r.Context(), database.CreateResetTokenParams{
		Token:      uuid.New(),
		Expiration: time.Now().Add(time.Minute * 10).UTC(),
		Email:      email,
	})

	resetToken := types.ResetPassword{
		Token:      dbResetToken.Token,
		Expiration: dbResetToken.Expiration,
		Email:      dbResetToken.Email,
	}

	return resetToken, err
}

func GetResetTokenByToken(r *http.Request, token uuid.UUID) (types.ResetPassword, error) {
	dbResetToken, err := DB.GetResetTokenByToken(r.Context(), token)

	resetToken := types.ResetPassword{
		Token:      dbResetToken.Token,
		Expiration: dbResetToken.Expiration,
		Email:      dbResetToken.Email,
	}

	return resetToken, err
}
