package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/naylorpmax/homebrew-users-api/pkg/middleware"
	"github.com/naylorpmax/homebrew-users-api/pkg/monster"
)

type (
	MonsterRequest struct {
		Name *string
	}

	MonsterLookup struct {
		MonsterService *monster.Service
	}
)

func (m *MonsterLookup) Handler(w http.ResponseWriter, r *http.Request) error {
	if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
		return &middleware.Error{
			Message:    "unsupported media type",
			Details:    fmt.Sprintf("expected 'application/json', got '%v'", contentType),
			StatusCode: http.StatusBadRequest,
		}
	}

	monsterReq := &MonsterRequest{}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(monsterReq); err != nil {
		return &middleware.Error{
			Message:    "unable to unmarshal request",
			Details:    err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	// temporary static data
	object := &monster.Object{
		Name:       "Swarm of Red Pikmin",
		Type:       "Small swarm of Tiny plants",
		Alignment:  "Neutral",
		ArmorClass: "11",
		HitPoints:  "15 (6d4)",
		Speed:      "20 ft, burrow 35 ft",
		AbilityScores: &monster.Scores{
			Strength:     "14 (+2)",
			Dexterity:    "12 (+1)",
			Constitution: "10 (+0)",
			Intelligence: "6 (-2)",
			Wisdom:       "10 (+0)",
			Charisma:     "12 (+1)",
		},
		Skills:              "Athletics +4",
		DamageImmunities:    "fire",
		ConditionImmunities: "exhaustion",
		Senses:              "darkvision 60 ft, passive Perception",
		Languages:           "Pikish; understands Common and Sylvan but cannot speak them",
		Challenge:           "1/4",
		Traits: map[string]string{
			"Natural Followers": "Any creature within 30 feet of the swarm that it can see and hear can attempt a DC 11 Wisdom (Animal Handling) or Charisma (Persuasion) check to charm the swarm. Once the swarm has been charmed by a creature, the swarm becomes immune to being charmed by any other creature as long as long as the creature it is charmed by is within 60 feet.",
			"Natural Growth":    "For every 10 minutes the swarm remains in direct sunlight, they automatically restore 1d4 hit points. If the swarm has full hit points and would regain hit points in this way, they instead gain a permanent +1 bonus to AC and their maximum number of hits points doubles. This can happen up to two times.",
			"Swarm":             "The swarm can occupy another creature’s space and vice-versa. The swarm can move through any opening large enough for a Tiny insect. The swarm can’t regain hit points except through their Natural Growth feature, and can’t gain temporary hit points.",
		},
		Actions: map[string]string{
			"Slams (swarm has more than half HP)": "Melee Weapon Attack: +4 to hit, reach 0 ft., one target in the swarm's space. Hit: 7 (2d4+2) bludgeoning damage.",
			"Slams (swarm has less than half HP)": "Melee Weapon Attack: +4 to hit, reach 0 ft., one target in the swarm’s space. Hit: 4 (1d4+2) bludgeoning damage.",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(object)
	if err != nil {
		return err
	}
	return nil
}
