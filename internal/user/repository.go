package user

import "github.com/google/uuid"

type User struct{
	Id uuid.UUID
	Email string
	Password string
}

type Notes struct {
	Title string
	MySpace string
}

type Repository interface {
	Create(email, password string) error
	GetByEmail(email string) (*User, error)
	DeleteUser(Id string) error

	SaveRefreshToken(email, refreshToken string) error
	GetRefreshToken(email string) (string, error)

	GetNotes(Id string) ([]Notes, error)
	MyNotes(Id string, notes Notes) error
	DeleteNotes(Id string) error
	UpdateNotes(updatedNotes Notes, Id string) error
}


