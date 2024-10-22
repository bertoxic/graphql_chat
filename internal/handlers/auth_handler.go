package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/internal/render"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
)

type Repository struct {
	a  *config.AppConfig
	db database.DatabaseRepo
}

var (
	oauth2Config *oauth2.Config
	// Generate a random state string to prevent CSRF attacks
	oauthStateString, _ = generateStateOauthCookie(32)
)

func initOauth() {
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
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func NewRepository(a *config.AppConfig, db database.DatabaseRepo) *Repository {
	initOauth()
	return &Repository{a: a, db: db}
}

var Repo Repository

func NewRepo(repository *Repository) {

	Repo = *repository
}
func (au *Repository) HomePage(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()
	//ctx := context.Context(context.Background())
	//user, err := au.db.GetUserByEmail(ctx, "henry@mail.com")
	//if err != nil {
	//	fmt.Printf("", errorx.New(errorx.ErrCodeDatabase, "", err))
	//	return
	//}
	render.Template(w, "home.page.gohtml", &models.TemplateData{
		StringMap: map[string]string{"name": "user.UserName"},
	})

}

func (au *Repository) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (au *Repository) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	//state := r.FormValue("state")
	//if state != oauthStateString {
	//	log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
	//	http.Error(w, "Invalid state", http.StatusBadRequest)
	//	return
	//}

	code := r.FormValue("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("oauth2Config.Exchange() failed with '%s'\n", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := oauth2Config.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("Failed to get user info: %s", err)
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		log.Printf("Failed to decode user info: %s", err)
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Now you have userInfo, you can save it to your database
	saveUserToDB(userInfo)
}

func saveUserToDB(userInfo map[string]interface{}) {
	// Logic to save user info to your database
	// Example: userInfo["id"], userInfo["email"], etc.
}
