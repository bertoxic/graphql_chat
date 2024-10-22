package handlers

import (
	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/internal/render"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"net/http"
)

type Repository struct {
	a  *config.AppConfig
	db database.DatabaseRepo
}

func NewRepository(a *config.AppConfig, db database.DatabaseRepo) *Repository {
	return &Repository{a: a, db: db}
}

var Repo Repository

func NewRepo(repository *Repository) {
	Repo = *repository
}
func (au *Repository) HomePage(w http.ResponseWriter, r *http.Request) {
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
