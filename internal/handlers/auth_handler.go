package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/bertoxic/graphqlChat/internal/auth"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/internal/render"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Repository struct {
	app         *config.AppConfig
	db          database.DatabaseRepo
	logger      *log.Logger
	authService AuthService
}

type AuthService interface {
	Register(ctx context.Context, input models.RegistrationInput) (*auth.AuthResponse, error)
	Login(ctx context.Context, input models.LoginInput) (*auth.AuthResponse, error)
}

type jsonResponse struct {
	Status      string      `json:"status"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data,omitempty"`
	RedirectURL string      `json:"redirect_url,omitempty"`
}

var (
	oauth2Config     *oauth2.Config
	oauthStateString string
)

func NewRepository(app *config.AppConfig, db database.DatabaseRepo, authService AuthService) *Repository {
	logger := log.New(os.Stdout, "[AUTH] ", log.LstdFlags|log.Lshortfile)
	repo := &Repository{
		app:         app,
		db:          db,
		logger:      logger,
		authService: authService,
	}
	initOauth()
	return repo
}

var Repo *Repository

func NewRepo(repo *Repository) {
	Repo = repo
}

func initOauth() {
	var err error
	oauthStateString, err = generateStateOauthCookie(32)
	if err != nil {
		log.Fatalf("Failed to generate OAuth state: %v", err)
	}

	//oauth2Config = &oauth2.Config{
	//	ClientID:     os.Getenv("CLIENT_ID"),
	//	ClientSecret: os.Getenv("CLIENT_SECRET"),
	//	Scopes: []string{
	//		"https://www.googleapis.com/auth/userinfo.email",
	//		"https://www.googleapis.com/auth/userinfo.profile",
	//	},
	//	Endpoint:     google.Endpoint,
	//	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	//}
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/googleCallback",
	}
}

func generateStateOauthCookie(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (r *Repository) HomePage(w http.ResponseWriter, req *http.Request) {
	render.Template(w, "home.page.gohtml", &models.TemplateData{
		StringMap: map[string]string{"name": "Welcome"},
	})
}

// HandleRegister handles user registration
func (r *Repository) HandleRegister(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.RegistrationInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		r.jsonResponse(w, jsonResponse{
			Status:  "error",
			Message: "Invalid request body",
		}, http.StatusBadRequest)
		return
	}

	authResp, err := r.authService.Register(req.Context(), input)
	if err != nil {
		r.logger.Printf("Registration failed: %v", err)
		r.jsonResponse(w, jsonResponse{
			Status:  "error",
			Message: "Registration failed",
		}, http.StatusInternalServerError)
		return
	}

	r.jsonResponse(w, jsonResponse{
		Status:  "success",
		Message: "Registration successful",
		Data:    authResp,
	}, http.StatusOK)
}

// HandleLogin handles user login
func (r *Repository) HandleLogin(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input models.LoginInput
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		r.jsonResponse(w, jsonResponse{
			Status:  "error",
			Message: "Invalid request body",
		}, http.StatusBadRequest)
		return
	}

	authResp, err := r.authService.Login(req.Context(), input)
	if err != nil {
		r.logger.Printf("Login failed: %v", err)
		r.jsonResponse(w, jsonResponse{
			Status:  "error",
			Message: "Invalid credentials",
		}, http.StatusUnauthorized)
		return
	}

	r.jsonResponse(w, jsonResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    authResp,
	}, http.StatusOK)
}

// HandleGoogleLogin initiates Google OAuth flow
func (r *Repository) HandleGoogleLogin(w http.ResponseWriter, req *http.Request) {
	state, err := generateStateOauthCookie(32)
	if err != nil {
		r.logger.Printf("Failed to generate state: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauthState",
		Value:    state,
		Expires:  time.Now().Add(10 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	url := oauth2Config.AuthCodeURL(state)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

// HandleGoogleCallback handles the Google OAuth callback
func (r *Repository) HandleGoogleCallback(w http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")
	stateCookie, err := req.Cookie("oauthState")
	if err != nil || state != stateCookie.Value {
		r.logger.Printf("Invalid OAuth state")
		http.Error(w, "Invalid OAuth state", http.StatusBadRequest)
		return
	}

	code := req.FormValue("code")
	token, err := oauth2Config.Exchange(req.Context(), code)
	if err != nil {
		r.logger.Printf("Token exchange failed: %v", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	userInfo, err := r.fetchGoogleUserInfo(req.Context(), token)
	if err != nil {
		r.logger.Printf("Failed to get user info: %v", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}

	// Create registration input from Google user info
	input := models.RegistrationInput{
		Email:    userInfo["email"].(string),
		Username: userInfo["name"].(string),
		Password: generateSecurePassword(), // Generate a secure random password for OAuth users
	}

	// Use the existing Register method from AuthService
	authResp, err := r.authService.Register(req.Context(), input)
	if err != nil {
		r.logger.Printf("Failed to register Google user: %v", err)
		http.Error(w, "Failed to create user account", http.StatusInternalServerError)
		return
	}

	r.jsonResponse(w, jsonResponse{
		Status:      "success",
		Message:     "Google authentication successful",
		Data:        authResp,
		RedirectURL: "/dashboard",
	}, http.StatusOK)
}

// Helper functions

func (r *Repository) fetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (map[string]interface{}, error) {
	client := oauth2Config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (r *Repository) jsonResponse(w http.ResponseWriter, response jsonResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func generateSecurePassword() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
