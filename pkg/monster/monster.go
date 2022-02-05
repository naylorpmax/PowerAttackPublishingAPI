package monster

type (
	Service struct {
		// DBClient *db.Client
	}

	Object struct {
		Name                string            `json:"name"`
		Type                string            `json:"type"`
		Alignment           string            `json:"alignment"`
		ArmorClass          string            `json:"armorClass,omitempty"`
		HitPoints           string            `json:"hitPoints"`
		Speed               string            `json:"speed"`
		AbilityScores       *Scores           `json:"abilityScores"`
		Skills              string            `json:"skills,omitempty"`
		DamageResistances   string            `json:"damageResistances,omitempty"`
		DamageImmunities    string            `json:"damageImmunities,omitempty"`
		ConditionImmunities string            `json:"conditionImmunities,omitempty"`
		Senses              string            `json:"senses,omitempty"`
		Languages           string            `json:"languages,omitempty"`
		Challenge           string            `json:"challenge,omitempty"`
		Traits              map[string]string `json:"traits,omitempty"`
		Actions             map[string]string `json:"actions,omitempty"`
		Source              string            `json:"source,omitempty"`
	}

	Scores struct {
		Strength     string
		Dexterity    string
		Constitution string
		Intelligence string
		Wisdom       string
		Charisma     string
	}
)
