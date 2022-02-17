package spell

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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

const (
	selectClause = "SELECT name, CAST(level AS TEXT), school, casting_time, range, " +
		"components, duration, classes, description, source " +
		"FROM spells "
)

func New(dbConn *pgx.Conn) (*Service, error) {
	if dbConn == nil {
		return nil, errors.New("unable to initialize servicer: missing db connection")
	}
	return &Service{
		DBConn: dbConn,
	}, nil
}

func (s *Service) Lookup(ctx context.Context, name *string, level *string) ([]*Object, error) {
	args := make(map[string]string)
	if name != nil && *name != "" {
		args["name"] = *name
	}
	if level != nil && *level != "" {
		args["level"] = *level
	}
	if len(args) == 0 {
		return nil, errors.New("unable to build query: no arguments provided to query")
	}
	sql, queryArgs := buildQuery(args)

	results, err := s.DBConn.Query(ctx, sql, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %v", err)
	}

	objects, err := scanObjects(results)
	if err != nil {
		return nil, fmt.Errorf("unable to scan objects: %v", err)
	}
	return objects, nil
}

func buildQuery(args map[string]string) (string, []interface{}) {
	whereClause := make([]string, 0)
	queryArgs := make([]interface{}, 0)
	i := 1
	for name, value := range args {
		argIndex := strconv.Itoa(i)
		whereClause = append(whereClause, "CAST("+name+" AS TEXT) LIKE $"+argIndex+" ")
		queryArgs = append(queryArgs, "%"+value+"%")
		i += 1
	}
	return selectClause + " WHERE " + strings.Join(whereClause, " AND "), queryArgs
}

func scanObjects(results pgx.Rows) ([]*Object, error) {
	objects := make([]*Object, 0)
	for results.Next() {
		object := &Object{}
		err := results.Scan(&object.Name, &object.Level, &object.School, &object.CastingTime,
			&object.Range, &object.Components, &object.Duration, &object.Classes, &object.Description,
			&object.Source)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}
