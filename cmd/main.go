package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/naylorpmax/homebrew-users-api/pkg/monster"
	"github.com/naylorpmax/homebrew-users-api/pkg/router"
	"github.com/naylorpmax/homebrew-users-api/pkg/spell"
)

func main() {
	mainLogger, _ := zap.NewProduction()
	defer func() {
		mainLogger.Sync()
		fmt.Println("syncing")
	}()

	mainLogger.Info("starting application")

	patreonClientID := os.Getenv("PATREON_CLIENT_ID")
	if patreonClientID == "" {
		mainLogger.Fatal("missing required environment variable",
			zap.String("environment_variable", "PATREON_CLIENT_ID"),
		)
	}

	patreonClientSecret := os.Getenv("PATREON_CLIENT_SECRET")
	if patreonClientID == "" {
		mainLogger.Fatal("missing required environment variable",
			zap.String("environment_variable", "PATREON_CLIENT_SECRET"),
		)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		mainLogger.Fatal("missing required environment variable",
			zap.String("environment_variable", "DB_URL"),
		)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     patreonClientID,
		ClientSecret: patreonClientSecret,
		RedirectURL:  "http://localhost:8080/welcome",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://patreon.com/api/oauth2/token",
			AuthURL:  "https://patreon.com/oauth2/authorize",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		mainLogger.Fatal("unable to connect to database",
			zap.String("error", err.Error()),
		)
	}
	defer conn.Close(ctx)

	select {
	case <-ctx.Done():
		mainLogger.Fatal("context cancelled before connecting to the database",
			zap.String("error", ctx.Err().Error()),
		)
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
		Logger:       *mainLogger,
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
