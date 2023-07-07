package database

import (
	"github.com/wk8/go-ordered-map"
)

type Pagination struct {
	Take int
	Skip int
}

type SelectionOptions struct {
	WhereClause   *string
	Fields        []string
	GroupByFields []string
	JoinClause    *orderedmap.OrderedMap
	Pagination    *Pagination
}

type SelectionOption func(*SelectionOptions) error

func WithWhere(whereClause string) SelectionOption {
	return func(so *SelectionOptions) error {
		so.WhereClause = &whereClause
		return nil
	}
}

func WithFields(fields []string) SelectionOption {
	return func(so *SelectionOptions) error {
		so.Fields = fields
		return nil
	}
}

func WithGroupBy(fields []string) SelectionOption {
	return func(so *SelectionOptions) error {
		so.GroupByFields = fields
		return nil
	}
}

func WithJoin(joinClause *orderedmap.OrderedMap) SelectionOption {
	return func(so *SelectionOptions) error {
		so.JoinClause = joinClause
		return nil
	}
}

func WithPagination(take, skip int) SelectionOption {
	return func(so *SelectionOptions) error {
		so.Pagination = &Pagination{
			Take: take,
			Skip: skip,
		}

		return nil
	}
}

type UpsertionOptions struct {
	FieldsValues map[string]interface{}
}

type UpsertionOption func(*UpsertionOptions) error

func WithFieldsValues(fieldsValues map[string]interface{}) UpsertionOption {
	return func(uo *UpsertionOptions) error {
		uo.FieldsValues = fieldsValues
		return nil
	}
}

type DeletionOptions struct {
	WhereClause *string
}

type DeletionOption func(*DeletionOptions) error

func WithDeleteWhere(whereClause string) DeletionOption {
	return func(do *DeletionOptions) error {
		do.WhereClause = &whereClause
		return nil
	}
}
