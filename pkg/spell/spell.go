package spell

type (
	Service struct{
		// DBClient *db.Client
	}

	Object struct {
		Name        string   `json:"name"`
		Level       string   `json:"level"`
		School      string   `json:"school"`
		CastingTime string   `json:"castingTime"`
		Range       string   `json:"range"`
		Components  string   `json:"components"`
		Duration    string   `json:"duration"`
		Classes     []string `json:"classes"`
		Description string   `json:"description"`
		Source      string   `json:"source"`
	}
)