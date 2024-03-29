package main

import (
	"MovieService/internal/pkg/utils/jwt"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log/slog"
	"net/http"
	"os"

	authHandler "MovieService/internal/pkg/auth/http"
	authRepo "MovieService/internal/pkg/auth/repo"
	authUsecase "MovieService/internal/pkg/auth/usecase"

	actorsHandler "MovieService/internal/pkg/actors/http"
	actorsRepo "MovieService/internal/pkg/actors/repo"
	actorsUsecase "MovieService/internal/pkg/actors/usecase"

	moviesHandler "MovieService/internal/pkg/movies/http"
	moviesRepo "MovieService/internal/pkg/movies/repo"
	moviesUsecase "MovieService/internal/pkg/movies/usecase"
)

// Логгер
// Swagger

// @title MovieService
// @version 1.0
// @description Api for movie db

// @host localhost:8080
// @schemes http
// @BasePath /
func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() (err error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Не удалось загрузить файл .env")
		return err
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	secretKey := os.Getenv("SECRET_KEY")

	db, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName))
	if err != nil {
		err = fmt.Errorf("error happened in sql.Open: %w", err)

		return err
	}
	defer db.Close()

	if err = db.Ping(context.Background()); err != nil {
		err = fmt.Errorf("error happened in db.Ping: %w", err)

		return err
	}

	err = jwt.LoadSecret(secretKey)
	if err != nil {
		fmt.Println("загрузка секрета")
		return err
	}

	log := &slog.Logger{}

	authRepo := authRepo.NewAuthRepo(db)
	authUsecase := authUsecase.NewAuthUsecase(authRepo)
	authHandler := authHandler.NewAuthHandler(log, authUsecase)

	actorRepo := actorsRepo.NewActorsRepo(db)
	actorUsecase := actorsUsecase.NewActorsUsecase(actorRepo)
	actorHandler := actorsHandler.NewActorsHandler(log, actorUsecase)

	movieRepo := moviesRepo.NewMoviesRepo(db)
	movieUsecase := moviesUsecase.NewMoviesUsecase(movieRepo)
	movieHandler := moviesHandler.NewMoviesHandler(log, movieUsecase)

	mux := http.NewServeMux()

	mux.Handle("/api/actors/", &actorHandler)
	mux.Handle("/api/movies/", &movieHandler)
	mux.Handle("/api/auth/", &authHandler)
	mux.Handle("/api/actors", &actorHandler)
	mux.Handle("/api/movies", &movieHandler)
	mux.Handle("/api/auth", &authHandler)

	return http.ListenAndServe(":8080", mux)
}
