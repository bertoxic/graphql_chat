package jwt

import (
	"context"
	"fmt"
	"github.com/bertoxic/graphqlChat/internal/models"
	"github.com/bertoxic/graphqlChat/pkg/config"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var signatureType jwa.SignatureAlgorithm = jwa.HS256
var (
	Conf         *config.AppConfig
	tokenService *TokenService
)

func TestMain(m *testing.M) {
	var err interface{}
	Conf, err = config.NewConfig(".env.test", "80")
	if err != nil {
		fmt.Printf("%v", err)
	}
	// Conf.JWT = config.JWT{Secret: []byte(os.Getenv(
	//	"JWT_SECRET")), Issuer: os.Getenv("ISSUER")}
	tokenService = NewTokenService(Conf)
	defer os.Exit(m.Run())
	// print(conf.DataBaseINFO.URL)
}

func TestTokenService_CreateAccessToken(t *testing.T) {
	t.Run("testing create Access token", func(t *testing.T) {
		ctx := context.Background()
		user := models.UserDetails{ID: "342"}

		token, err := tokenService.CreateAccessToken(ctx, user)
		require.NoError(t, err)
		tok, err := jwt.Parse(
			[]byte(token),
			jwt.WithValidate(true),
			jwt.WithVerify(signatureType, Conf.JWT.Secret),
			jwt.WithIssuer(Conf.JWT.Issuer),
		)
		require.NoError(t, err)

		t.Logf("Parsed Token: %+v, Generated Token: %+v\n", tok, token)
	})
}
