package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Nurdaulet-no/auth-svc/internal/delivery/http"
	"github.com/Nurdaulet-no/auth-svc/internal/repository/memory"
	"github.com/Nurdaulet-no/auth-svc/internal/usecase"
	"github.com/Nurdaulet-no/auth-svc/pkg/jwt"
)


func main() {
	userRepo := memory.NewUserRepo()
	jwtm := jwt.NewManager("dev-secret-change-me", 24*time.Hour)

	idGen := func() string { return time.Now().Format("20060102150405.000000000") }

	authUC := usecase.NewAuthService(userRepo, jwtm, idGen)
	h := httpdelivery.NewHandler(authUC, jwtm)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/register", h.Register)
	mux.HandleFunc("/login", h.Login)

	me := http.HandlerFunc(h.Me)
	mux.Handle("/me", h.AuthMiddleware(me))

	log.Println("HTTP server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
