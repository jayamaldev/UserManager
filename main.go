package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"user-manager/api"
	"user-manager/config"
	"user-manager/database"
	_ "user-manager/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error on configurations", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/doc/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/docs/swagger.json", cfg.APPPort)),
	))

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "docs"))
	fileServer(r, "/docs", filesDir)

	server, err := ConnectDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	r.Route("/users", server.UserRouter)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.APPPort),
		Handler: r,
	}
	go func() {
		fmt.Println("Server is Running on port 8080")
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Error on Listening server!")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Signal Received")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Forced Shutdown: ", err)
	}

	server.Pool.Close()
	log.Println("DB Connection Pools Closed")

	log.Println("Server Exited Gracefully")
}

func ConnectDatabase(cfg *config.Config) (*api.Server, error) {
	ctx := context.Background()

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	fmt.Println("DB connection string: ", connString)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	queries := database.New(pool)
	server := api.NewServer(queries, pool)

	return server, nil
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
