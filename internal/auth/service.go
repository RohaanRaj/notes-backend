package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"time"

	"notesApp/domain"
	"notesApp/internal/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct{
	db user.Repository
	jwtSecret []byte
}

func NewService(db user.Repository, jwtSecret []byte) *Service{
	return &Service{
		db: db,
		jwtSecret: jwtSecret,
	}
}

func(s *Service) DeleteUser(Id string) error{
	err := s.db.DeleteUser(Id)
	if err != nil{
		return err
	}
	return nil
}

func(s *Service) UpdateNotes(updatedNotes user.Notes, Id string) error{

	err := s.db.UpdateNotes(updatedNotes, Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return domain.ErrNotFound

		default: return err
		}
	}
	return nil

}

func(s *Service) DeleteNotes(Id string) error{
		
	err := s.db.DeleteNotes(Id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return domain.ErrNotFound

		default: return err
		}
	}
	return nil
}

func(s *Service) GetNotes(Id string) ([]user.Notes, error){
	notes, err := s.db.GetNotes(Id)
	if err != nil {
		switch err{
			case sql.ErrNoRows:
				return nil, domain.ErrNotFound
			default : return nil, err
	}
}
	return notes, nil

}

func(s *Service) Notes(Id string, notes user.Notes) error{

	if err := s.db.MyNotes(Id, notes); err != nil{
		switch err {
		case sql.ErrNoRows:
			return domain.ErrNotFound

		default: return err
		}
	}
	return nil
}

func (s *Service) Login(email, password string) (AuthTokens, error){

	u, err := s.db.GetByEmail(email)
	if err != nil{ 
		switch err {
		case sql.ErrNoRows:
			return AuthTokens{}, domain.ErrNotRegistered

		default: return AuthTokens{}, err
		}

	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return AuthTokens{}, domain.ErrInternalServer
	}


	claims := Claims{
		UserId: u.Id.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15*time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return AuthTokens{}, domain.ErrInternalServer
	}


	if refreshToken, _ := s.db.GetRefreshToken(u.Email); refreshToken != ""{
		result := AuthTokens{
			AccessToken: tokenStr,
			RefreshToken: refreshToken,
		}
	return result, nil
	}

	refreshToken := generateRefreshToken()
	if err := s.db.SaveRefreshToken(u.Email, refreshToken); err != nil {
		return AuthTokens{}, domain.ErrDb
	}
	result := AuthTokens{
			AccessToken: tokenStr,
			RefreshToken: refreshToken,
	}
	return result, nil

}

func generateRefreshToken() string{
	temp := make([]byte, 32)
	rand.Read(temp)
	return base64.URLEncoding.EncodeToString(temp)
}


func(s *Service) Register(user user.User ) error{

	hashedpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.ErrInternalServer
	}
	strPass := string(hashedpass)
	if err := s.db.Create(user.Email, strPass); err != nil {
		return domain.ErrUserAlreadyExists
	}
	return nil 
}

