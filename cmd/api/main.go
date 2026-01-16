package main

import (
	"database/sql"
	"net/http"
	"os"

	"notesApp/internal/auth"
	"notesApp/internal/http/handler"
	"notesApp/internal/user"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)


func main(){
var port string = os.Getenv("PORT")
var dbUrl string = os.Getenv("DATABASE_URL")
var jwt_Secret string = os.Getenv("JWT_SECRET")

var pool *sql.DB //database connection pool
var dsn string = dbUrl 
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	if dsn == ""{
		logger.Fatal("DATABASE URL EMPTY")
	}
	if jwt_Secret == ""{
		logger.Fatal("JWT SECRET CANNOT BE EMPTY")
	}
	if port == ""{
		port = "8080"
	}
	pool, err = sql.Open("pgx", dsn)
	if err != nil {
		logger.Fatal("UNABLE TO CREATE CONNECTION POOL! CHECK DB URL")
	}
	err = pool.Ping()
	if err != nil {
		logger.Fatal("Error", zap.String("reason", "Unable to connect to database"))
	}
	userRepo := user.NewPostgresRepository(pool)
	jwtSecret := []byte(jwt_Secret)
	authService := *auth.NewService(userRepo, jwtSecret)
	authHandler := handler.NewAuthHandler(userRepo, authService, logger)
	r := chi.NewRouter()
	r.Use(auth.LoggingMiddleware(logger))
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTMiddleware(jwtSecret))
		r.Post("/notes", authHandler.Notes)
		r.Get("/notes", authHandler.GetNotes)
		r.Delete("/notes", authHandler.DeleteNotes)
		r.Patch("/notes", authHandler.UpdateNotes)
		r.Delete("/user", authHandler.DeleteUser)
	})
	//r.Get("/test", testHandler)

	port = ":" + port
	http.ListenAndServe(port, r)
}


