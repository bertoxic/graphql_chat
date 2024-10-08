package handlers

import (
	"github.com/bertoxic/graphqlChat/internal/app"
	"github.com/bertoxic/graphqlChat/internal/database"
	"github.com/bertoxic/graphqlChat/internal/render"
	"net/http"
)

type Repository struct {
	a  *app.App
	db database.DatabaseRepo
}

func NewRepository(a *app.App, db database.DatabaseRepo) *Repository {
	return &Repository{a: a, db: db}
}

var Repo Repository

func NewRepo(repository *Repository) {
	Repo = *repository
}
func (au *Repository) HomePage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "home.page.gohtml")
}
