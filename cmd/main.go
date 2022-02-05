package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/naylorpmax/homebrew-users-api/pkg/monster"
	"github.com/naylorpmax/homebrew-users-api/pkg/spell"
	"github.com/naylorpmax/homebrew-users-api/pkg/router"
)

func main() {
	oauth2Config := &oauth2.Config{
		ClientID:     os.Getenv("PATREON_CLIENT_ID"),
		ClientSecret: os.Getenv("PATREON_CLIENT_SECRET"),
        RedirectURL:  "http://localhost:8080/welcome",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://patreon.com/api/oauth2/token",
			AuthURL: "https://patreon.com/oauth2/authorize",
		},
	}
	
	// TODO: initialize logger

	// TODO: initialize DB client

	spellSvc := &spell.Service{}
	monsterSvc := &monster.Service{}

	routerCfg := router.Config{
		MonsterService: monsterSvc,
		SpellService: spellSvc,
		OAuth2Config: oauth2Config,
	}

	server := http.Server{
		Addr: "localhost:8080",
		Handler: router.New(routerCfg),
	}

	fmt.Println("listening on :8080")
	log.Fatal(server.ListenAndServe())
}
