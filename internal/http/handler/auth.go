package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"notesApp/domain"
	"notesApp/internal/auth"
	"notesApp/internal/user"

	"go.uber.org/zap"
)


type AuthHandler struct {
	db user.Repository
	s auth.Service
	logger *zap.Logger 
}

func NewAuthHandler (db user.Repository, s auth.Service, logger *zap.Logger) *AuthHandler{
	return &AuthHandler{
		db: db,
		s: s,
		logger: logger,
	}
}


func(a *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request){
	Id := r.Context().Value(auth.IdContextKey).(string)
	err := a.s.DeleteUser(Id)
	if err != nil{
		a.logger.Error("Unable to delete", zap.String("user_id", Id), zap.Error(err))
		http.Error(w, "unable to delete user", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User Deleted"))
}

func(a *AuthHandler) UpdateNotes(w http.ResponseWriter, r *http.Request){
	var updatedNotes user.Notes
	if err := json.NewDecoder(r.Body).Decode(&updatedNotes); err != nil{
		a.logger.Warn("Malformed input")
		http.Error(w, "Error while decoding json", http.StatusBadRequest)
		return
	}

	Id := r.Context().Value(auth.IdContextKey).(string)

	err := a.s.UpdateNotes(updatedNotes, Id)
	if errors.Is(err, domain.ErrNotFound){
		a.logger.Info("Requested resource not found", zap.String("user_id", Id))
		http.Error(w, "Requested resource is not found", http.StatusNotFound)
		return
	}
	if err != nil{
		a.logger.Error("Database Error", zap.String("user_id", Id), zap.Error(err))
		http.Error(w, "Unable to update", http.StatusServiceUnavailable)
		return
	}
	w.Write([]byte("Updated successfully\n"))
}

func(a *AuthHandler) DeleteNotes(w http.ResponseWriter, r *http.Request){
	Id := r.Context().Value(auth.IdContextKey).(string)

	err := a.s.DeleteNotes(Id)
	if errors.Is(err, domain.ErrNotFound){
		a.logger.Info("Requested resource not found", zap.String("user_id", Id))
		http.Error(w, "Requested resource is not found", http.StatusNotFound)
		return
	}
	if err != nil{
		http.Error(w, "Unable to delete", http.StatusServiceUnavailable)
		a.logger.Error("Database Error", zap.String("user_id", Id))
		return
	}
	w.Write([]byte("Deleted successfully\n"))

}

func(a *AuthHandler) GetNotes(w http.ResponseWriter, r *http.Request){
	Id := r.Context().Value(auth.IdContextKey).(string)

	notes, err := a.s.GetNotes(Id)
	if errors.Is(err, domain.ErrInvalidCredentials){
		a.logger.Warn("Access Denied", zap.String("user_id", Id))
		http.Error(w, "Details not found", http.StatusUnauthorized)
		return
	}
	if errors.Is(err, domain.ErrNotFound){
		a.logger.Info("Requested resource not found", zap.String("user_id", Id))
		http.Error(w, "Requested resource is not found", http.StatusNotFound)
		return
	}
	if err != nil{
		a.logger.Error("Database Error", zap.String("user_id", Id))
		http.Error(w, "Unable to fetch notes right now", http.StatusServiceUnavailable)
		return
	}

	json.NewEncoder(w).Encode(notes)

}

func (a *AuthHandler) Notes(w http.ResponseWriter, r *http.Request){
	var notes user.Notes
	if err := json.NewDecoder(r.Body).Decode(&notes); err != nil{
		a.logger.Warn("Malformed input")
		http.Error(w, "Error while decoding error", http.StatusBadRequest)
		return
	}

	Id := r.Context().Value(auth.IdContextKey).(string)

	err := a.s.Notes(Id, notes)

	if errors.Is(err, domain.ErrNotFound){
		a.logger.Info("Not found", zap.String("user_id", Id))
		http.Error(w, "Requested resource is not found", http.StatusNotFound)
		return
	}
	if errors.Is(err, domain.ErrInvalidCredentials){
		a.logger.Warn("Access Denied", zap.String("user_id", Id))
		http.Error(w, "Details not found", http.StatusUnauthorized)
		return
	}
	if err != nil {
		a.logger.Error("Unexpected Error while taking notes",zap.String("user_id", Id), zap.Error(err))
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Notes taken"))
}

func(a *AuthHandler) Register(w http.ResponseWriter, r *http.Request){

	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil{
		a.logger.Warn("Malformed Input")
		http.Error(w, "Error while decoding error", http.StatusBadRequest)
		return
	}
	err := a.s.Register(u)
	if errors.Is(err, domain.ErrInternalServer){
		a.logger.Error("Internal Server Error!!", zap.Error(err))
		http.Error(w,"Internal server error", http.StatusInternalServerError)
		return
	}
	if errors.Is(err, domain.ErrDb){
		a.logger.Error("Registration error", zap.String("reason", "invalid_credentials"), zap.Error(err))
		http.Error(w,"Database error", http.StatusServiceUnavailable)
		return
	}
	if errors.Is(err, domain.ErrUserAlreadyExists){
		a.logger.Info("Registration failed", zap.String("email", u.Email), zap.String("reason","User already exists"), zap.Error(err))
		http.Error(w,"User exists", http.StatusConflict)
		return
	}
	if err != nil {
		a.logger.Error("Registration failed", zap.String("email", u.Email), zap.Error(err))
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User Registered Successfully\n"))
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request){
	var cred auth.LoginReq

	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		a.logger.Warn("Malformed input")
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	tokens, err := a.s.Login(cred.Email, cred.Password)
	if errors.Is(domain.ErrInternalServer, err){
		a.logger.Error("Internal Server Error", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if errors.Is(domain.ErrNotRegistered, err){
		a.logger.Info("User Not Registered", zap.String("email", cred.Email))
		http.Error(w, "User not registered", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}
