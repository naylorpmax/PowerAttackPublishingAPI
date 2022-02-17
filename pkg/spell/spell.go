package spell

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type (
	Service struct {
		DBConn *pgx.Conn
	}

	Object struct {
		Name        string `json:"name"`
		Level       string `json:"level"`
		School      string `json:"school"`
		CastingTime string `json:"castingTime"`
		Range       string `json:"range"`
		Components  string `json:"components"`
		Duration    string `json:"duration"`
		Classes     string `json:"classes"`
		Description string `json:"description"`
		Source      string `json:"source"`
	}
)

func (s *Service) Lookup(ctx context.Context, name string) ([]*Object, error) {
	results, err := s.DBConn.Query(
		ctx,
		"SELECT name, CAST(level AS TEXT), school, casting_time, range, components, duration, classes, description, source "+
			"FROM spells "+
			"WHERE name LIKE $1",
		[]interface{}{"%" + name + "%"}...)
	if err != nil {
		return nil, fmt.Errorf("error querying db: %v", err)
	}

	objects := make([]*Object, 0)
	for results.Next() {
		object := &Object{}
		err := results.Scan(&object.Name, &object.Level, &object.School, &object.CastingTime,
			&object.Range, &object.Components, &object.Duration, &object.Classes, &object.Description,
			&object.Source)
		if err != nil {
			return nil, fmt.Errorf("error scanning db result: %v", err.Error())
		}
		objects = append(objects, object)
	}
	return objects, nil
}
