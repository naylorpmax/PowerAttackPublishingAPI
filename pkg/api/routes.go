package api

import(
	"net/http"
)

const (
	HealthEndpoint = "/health"

	HomeEndpoint = "/home"

	LoginEndpoint = "/login"

	WelcomeEndpoint = "/welcome"

	MonsterLookupEndpoint = "/monster/lookup"

	SpellLookupEndpoint = "/spell/lookup"
)

type (
	Route struct {
		Method string
		Path string
	}
)

var (
	HealthRoute = Route{
		Method: http.MethodGet,
		Path: HealthEndpoint,
	}

	HomeRoute = Route{
		Method: http.MethodGet,
		Path: HomeEndpoint,
	}

	LoginRoute = Route{
		Method: http.MethodGet,
		Path: LoginEndpoint,
	}

	WelcomeRoute = Route{
		Method: http.MethodGet,
		Path: WelcomeEndpoint,
	}

	MonsterLookupRoute = Route{
		Method: http.MethodPost,
		Path: MonsterLookupEndpoint,
	}

	SpellLookupRoute = Route{
		Method: http.MethodPost,
		Path: SpellLookupEndpoint,
	}
)
