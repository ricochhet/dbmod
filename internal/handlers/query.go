package handlers

import "github.com/ricochhet/dbmod/pkg/logx"

type Query struct {
	Name  string
	Query func(*Database) (*Database, error)
}

type Queries []Query

func (q *Queries) Run(db *Database, skip map[string]struct{}) *Database {
	cur := db

	for _, query := range *q {
		if _, ok := skip[query.Name]; ok {
			continue
		}

		result, err := query.Query(cur)
		if err != nil {
			logx.Infof("[%s] error: %v\n", query.Name, err)
			continue
		}

		if result != nil {
			cur = result
		}
	}

	return cur
}

func (q *Queries) skip(names []string) map[string]struct{} {
	if len(names) == 1 && names[0] == "all" {
		return nil
	}

	want := make(map[string]struct{}, len(names))
	for _, n := range names {
		want[n] = struct{}{}
	}

	skip := make(map[string]struct{})

	for _, q := range *q {
		if _, ok := want[q.Name]; !ok {
			skip[q.Name] = struct{}{}
		}
	}

	return skip
}
