package router

import (
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"

	"github.com/naylorpmax/homebrew-users-api/pkg/api"
	"github.com/naylorpmax/homebrew-users-api/pkg/middleware"
	"github.com/naylorpmax/homebrew-users-api/pkg/monster"
	"github.com/naylorpmax/homebrew-users-api/pkg/spell"
)

type (
	Config struct {
		// dbClient *db.Client
		OAuth2Config *oauth2.Config
		// Logger logger.Logger
		MonsterService *monster.Service
		SpellService   *spell.Service
	}
)

func New(cfg Config) *mux.Router {
	routes := mux.NewRouter().StrictSlash(true)

	// routes.Use(/*middleware*/)
	// health := &api.Health{DBClient: cfg.DBClient}

	// routes.Methods(api.Health.Method).
	// 	Path(api.Health.Path).
	// 	Handler(health.Handler)

	home := &api.Home{}

	routes.Methods(api.HomeRoute.Method).
		Path(api.HomeRoute.Path).
		Handler(middleware.Middleware(home.Handler))

	login := &api.Login{
		OAuth2Config: cfg.OAuth2Config,
	}

	routes.Methods(api.LoginRoute.Method).
		Path(api.LoginRoute.Path).
		Handler(middleware.Middleware(login.Handler))

	welcome := &api.Welcome{
		OAuth2Config: cfg.OAuth2Config,
	}

	routes.Methods(api.WelcomeRoute.Method).
		Path(api.WelcomeRoute.Path).
		Handler(middleware.Middleware(welcome.Handler))

	monsterLookup := &api.MonsterLookup{
		MonsterService: cfg.MonsterService,
	}

	routes.Methods(api.MonsterLookupRoute.Method).
		Path(api.MonsterLookupRoute.Path).
		Handler(middleware.Middleware(monsterLookup.Handler))

	spellLookup := &api.SpellLookup{
		SpellService: cfg.SpellService,
	}

	routes.Methods(api.SpellLookupRoute.Method).
		Path(api.SpellLookupRoute.Path).
		Handler(middleware.Middleware(spellLookup.Handler))

	return routes
}
