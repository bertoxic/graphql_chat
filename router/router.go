package router

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bertoxic/graphqlChat/graph"
	"github.com/bertoxic/graphqlChat/graph/resolvers"
	"github.com/bertoxic/graphqlChat/internal/app"
	"github.com/bertoxic/graphqlChat/internal/handlers"
	"github.com/bertoxic/graphqlChat/internal/jwt"
	"github.com/bertoxic/graphqlChat/internal/middlewares"
	"github.com/bertoxic/graphqlChat/internal/posts"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

var mux = chi.NewRouter()

func Routes(app *app.App) http.Handler {
	mux.Use(middleware.Logger)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Recoverer)
	tokenService := jwt.NewTokenService(app.Config)
	postRepo := posts.NewPostRepo(app.DB)
	postService := posts.NewPostServiceImpl(postRepo)

	mux.Use(middlewares.AuthMiddleWare(tokenService))
	mux.Use(middleware.Timeout(time.Second * 45))
	mux.Get("/", handlers.Repo.HomePage)
	mux.Get("/googleLogin", handlers.Repo.HandleGoogleLogin)
	mux.Get("/googleCallback", handlers.Repo.HandleGoogleCallback)
	mux.Handle("/play", playground.Handler("Graphql-chat", "/query"))
	mux.Handle("/query", handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &resolvers.Resolver{
					AuthService:     app.Services.AuthService,
					AuthUserService: app.Services.UserAuthService,
					PostService:     postService,
					UserService:     *app.Services.UserService,
				},
			},
		),
	))
	return mux
}
