package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"golang.org/x/oauth2"

	"github.com/naylorpmax/homebrew-users-api/pkg/monster"
	"github.com/naylorpmax/homebrew-users-api/pkg/router"
	"github.com/naylorpmax/homebrew-users-api/pkg/spell"
)

func main() {
	patreonClientID := os.Getenv("PATREON_CLIENT_ID")
	if patreonClientID == "" {
		fmt.Println("missing required environment variable: $PATREON_CLIENT_ID")
	}

	patreonClientSecret := os.Getenv("PATREON_CLIENT_SECRET")
	if patreonClientID == "" {
		fmt.Println("missing required environment variable: $PATREON_CLIENT_SECRET")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		fmt.Println("missing required environment variable: $DB_URL")
		os.Exit(1)
	}
	fmt.Println(dbURL)

	oauth2Config := &oauth2.Config{
		ClientID:     patreonClientID,
		ClientSecret: patreonClientSecret,
		RedirectURL:  "http://localhost:8080/welcome",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://patreon.com/api/oauth2/token",
			AuthURL:  "https://patreon.com/oauth2/authorize",
		},
	}

	// TODO: initialize logger

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		fmt.Println("unable to connect to database: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close(ctx)

	select {
	case <-ctx.Done():
		fmt.Println("context cancelled before connecting to the database: ", ctx.Err().Error())
	default:
	}

	routerCfg := router.Config{
		MonsterService: &monster.Service{
			DBConn: conn,
		},
		SpellService: &spell.Service{
			DBConn: conn,
		},
		OAuth2Config: oauth2Config,
	}

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router.New(routerCfg),
	}

	fmt.Println("listening on :8080")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("shutting down server: ", err.Error())
		os.Exit(1)
	}
}
