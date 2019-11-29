package accounts

import (
	"strings"

	"github.com/bradfitz/slice"
	"github.com/jmoiron/sqlx"
)

// Predefined query section type
const (
	QsWhere QuerySection = iota + 1
	QsOrderBy
	QsLimit
	QsOffset
)

// Predefined users list ordering
const (
	// ORDER BY created_at ASC
	OrderByCreatedAtAsc Order = iota + 1
	// ORDER BY created_at DESC
	OrderByCreatedAtDesc
	// ORDER BY updated_at ASC
	OrderByUpdatedAtAsc
	// ORDER BY updated_at DESC
	OrderByUpdatedAtDesc
)

var orderQueryMap = map[Order]string{
	OrderByCreatedAtAsc:  "created_at ASC",
	OrderByCreatedAtDesc: "created_at DESC",
	OrderByUpdatedAtAsc:  "updated_at ASC",
	OrderByUpdatedAtDesc: "updated_at DESC",
}

type (
	// Condition struct
	Condition interface {
		Query() string
		Params() []interface{}
		Type() QuerySection
	}

	// Order type
	Order int

	// QuerySection type
	QuerySection int

	condition struct {
		query  string
		params []interface{}
		t      QuerySection
	}
)

func (o Order) String() string {
	if v, ok := orderQueryMap[o]; ok {
		return v
	}
	return ""
}

// Query string
func (c condition) Query() string {
	return c.query
}

// Params for query string
func (c condition) Params() []interface{} {
	return c.params
}

// Type of query section string
func (c condition) Type() QuerySection {
	return c.t
}

// Disabled func add condition to select items only with disabled=v
func Disabled(v bool) Condition {
	return condition{
		query:  "disabled=?",
		params: []interface{}{v},
		t:      QsWhere,
	}
}

// Role func add condition to select query
func Role(v string) Condition {
	return condition{
		query:  "role=?",
		params: []interface{}{v},
		t:      QsWhere,
	}
}

// UserID func add condition to select query
func UserID(id string) Condition {
	return condition{
		query:  "user_id=?",
		params: []interface{}{id},
		t:      QsWhere,
	}
}

// UserIDs func add condition to select query
func UserIDs(ids []string) Condition {
	query, args, _ := sqlx.In("user_id IN (?)", ids)
	return condition{
		query:  query,
		params: args,
		t:      QsWhere,
	}
}

// AccountID func add condition to select query
func AccountID(id string) Condition {
	return condition{
		query:  "account_id=?",
		params: []interface{}{id},
		t:      QsWhere,
	}
}

// AccountIDs func add condition to select query
func AccountIDs(ids []string) Condition {
	query, args, _ := sqlx.In("account_id IN (?)", ids)
	return condition{
		query:  query,
		params: args,
		t:      QsWhere,
	}
}

// IDs func add condition to select query
func IDs(ids []string) Condition {
	query, args, _ := sqlx.In("id IN (?)", ids)
	return condition{
		query:  query,
		params: args,
		t:      QsWhere,
	}
}

// OrderBy func add ordering to selected list
func OrderBy(order ...Order) Condition {
	q := "ORDER BY"
	for k, o := range order {
		if k > 0 {
			q = q + ", " + o.String()
		} else {
			q = q + " " + o.String()
		}
	}
	return condition{
		query: q,
		t:     QsOrderBy,
	}
}

// Limit func add limit to select query
func Limit(v int) Condition {
	if v <= 0 {
		v = 100
	}
	return condition{
		query:  "LIMIT ?",
		params: []interface{}{v},
		t:      QsLimit,
	}
}

// Offset func add offset to select query
func Offset(v int) Condition {
	if v < 0 {
		v = 0
	}
	return condition{
		query:  "OFFSET ?",
		params: []interface{}{v},
		t:      QsOffset,
	}
}

// ConditionsToQuery represents conditions slice to query string and parameters slice
func conditionsToQuery(cs ...Condition) (q string, params []interface{}) {
	if len(cs) == 0 {
		return
	}

	slice.Sort(cs[:], func(i, j int) bool {
		return cs[i].Type() < cs[j].Type()
	})

	var where, limit, offset, orderBy string
	for _, c := range cs {
		switch c.Type() {
		case QsWhere:
			where += " AND " + c.Query()
		case QsLimit:
			limit = c.Query()
		case QsOffset:
			offset = c.Query()
		case QsOrderBy:
			orderBy = c.Query()
		}
		params = append(params, c.Params()...)
	}

	if where != "" {
		q = " WHERE " + strings.TrimPrefix(where, " AND ")
	}
	if orderBy != "" {
		q += " " + orderBy
	}
	if limit != "" {
		q += " " + limit
	}
	if offset != "" {
		q += " " + offset
	}

	return q, params
}
