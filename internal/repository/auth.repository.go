package repository

import (
	"context"
	"errors"
	"log"

	"github.com/federus1105/koda-b4-final-backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) Register(ctx context.Context, hashed string, user models.AuthRegister) (models.AuthRegister, error) {
	// --- START TRANSACTION ---
	tx, err := a.db.Begin(ctx)
	if err != nil {
		log.Println("Failed to begin transaction : ", err)
		return models.AuthRegister{}, err
	}
	defer tx.Rollback(ctx)
	var userID int

	// --- INSERT TABLE USER ---
	err = tx.QueryRow(ctx,
		`INSERT INTO users (email, password) VALUES ($1, $2) 
		RETURNING id`, user.Email, hashed).Scan(&userID)
	if err != nil {
		log.Println("Failed to insert user :", err)
		return models.AuthRegister{}, err
	}

	// --- INSERT TABLE ACCOUNT ---
	_, err = tx.Exec(ctx,
		`INSERT INTO account (id_users, fullname) VALUES ($1, $2)`,
		userID, user.Fullname)
	if err != nil {
		log.Println("Failed to insert account : ", err)
		return models.AuthRegister{}, err
	}

	// --- COMMIT TRANSACTION ---
	if err = tx.Commit(ctx); err != nil {
		log.Println("Failed to commit Transaction : ", err)
		return models.AuthRegister{}, err
	}

	return models.AuthRegister{
		Id:       userID,
		Email:    user.Email,
		Fullname: user.Fullname,
	}, nil
}

func (a *AuthRepository) Login(ctx context.Context, email string) (models.AuthLogin, error) {
	sql := `SELECT id, email, password, role FROM users WHERE email = $1`
	var user models.AuthLogin
	if err := a.db.QueryRow(ctx, sql, email).Scan(&user.Id, &user.Email, &user.Password, &user.Role); err != nil {
		if err == pgx.ErrNoRows {
			return models.AuthLogin{}, errors.New("user not found")
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		return models.AuthLogin{}, err
	}
	return user, nil
}
