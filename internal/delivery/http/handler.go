package httpdelivery

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Nurdaulet-no/auth-svc/internal/domain"
	"github.com/Nurdaulet-no/auth-svc/internal/usecase"
	"github.com/Nurdaulet-no/auth-svc/pkg/jwt"
)

type Handler struct {
	auth *usecase.AuthService
	jwt *jwt.Manager
}

func NewHandler(authService *usecase.AuthService, jwtManager *jwt.Manager) *Handler{
	return  &Handler{
		auth: authService,
		jwt: jwtManager,
	}
}

type registerRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Login string `json:"login"`
	Password string `json:"password"`
}

func writeJson(w http.ResponseWriter, code int, v any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request){
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error" : "bad json"})
		return
	}

	u, err := h.auth.Register(req.Email, req.Password)
	if err != nil {
		if err == domain.ErrEmailTaken {
			writeJson(w, http.StatusConflict, map[string]string{"error": "email taken"})
			return
		}
		writeJson (w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJson (w, http.StatusCreated, map[string]any{
		"id": u.ID,
		"email": u.Email,
	})	
}

func (h*Handler) Login (w http.ResponseWriter, r *http.Request){
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson (w, http.StatusBadRequest, map[string]string{"error": "bad json"})
		return
	}
	 
	token, err := h.auth.Login(req.Login, req.Password)
	if err != nil {
		if err == domain.ErrInvalidCreadiantals {
			writeJson (w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}
		writeJson (w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	
	writeJson (w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc (func (w http.ResponseWriter, r *http.Request){
		authH := r.Header.Get("Authorization")
		if !strings.HasPrefix(authH, "Bearer "){
			writeJson (w, http.StatusUnauthorized, map[string]string{"error": "missing or invalid authorization header"})
			return
		}
		tokenStr := strings.TrimPrefix (authH, "Bearer ")
		userID, err := h.jwt.Parse(tokenStr)
		if err != nil {
			writeJson (w, http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			return
		}

		ctx := WithUserID(r.Context(), userID)
		next.ServeHTTP (w, r.WithContext(ctx))
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeJson(w, http.StatusUnauthorized, map[string]string{"error": "no auth"})
		return
	}

	u, err := h.auth.Me(userID)
	if err != nil {
		writeJson(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	writeJson(w, http.StatusOK, map[string]any{
		"id":    u.ID,
		"email": u.Email,
	})
}