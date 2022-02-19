package router

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/naylorpmax/homebrew-users-api/pkg/api"
	"github.com/naylorpmax/homebrew-users-api/pkg/middleware/apierror"
	"github.com/naylorpmax/homebrew-users-api/pkg/middleware/monitoring"
	"github.com/naylorpmax/homebrew-users-api/pkg/monster"
	"github.com/naylorpmax/homebrew-users-api/pkg/spell"
)

type (
	Config struct {
		OAuth2Config   *oauth2.Config
		Logger         zap.Logger
		MonsterService *monster.Service
		SpellService   *spell.Service
	}
)

func New(cfg Config) *mux.Router {
	routes := mux.NewRouter().StrictSlash(true)
	routes.Use(monitoring.Monitoring(cfg.Logger))

	home := &api.Home{}

	routes.Methods(api.HomeRoute.Method).
		Path(api.HomeRoute.Path).
		Handler(apierror.Middleware(home.Handler))

	login := &api.Login{
		OAuth2Config: cfg.OAuth2Config,
	}

	routes.Methods(api.LoginRoute.Method).
		Path(api.LoginRoute.Path).
		Handler(apierror.Middleware(login.Handler))

	welcome := &api.Welcome{
		OAuth2Config: cfg.OAuth2Config,
	}

	routes.Methods(api.WelcomeRoute.Method).
		Path(api.WelcomeRoute.Path).
		Handler(apierror.Middleware(welcome.Handler))

	monsterLookup := &api.MonsterLookup{
		MonsterService: cfg.MonsterService,
	}

	routes.Methods(api.MonsterLookupRoute.Method).
		Path(api.MonsterLookupRoute.Path).
		Handler(apierror.Middleware(monsterLookup.Handler))

	spellLookup := &api.SpellLookup{
		SpellService: cfg.SpellService,
	}

	routes.Methods(api.SpellLookupRoute.Method).
		Path(api.SpellLookupRoute.Path).
		Handler(apierror.Middleware(spellLookup.Handler))

	return routes
}
