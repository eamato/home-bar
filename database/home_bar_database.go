package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"home-bar/configs"
	"home-bar/domain"
	"home-bar/internal"
	"strings"
	"time"
)

func NewHomeBarDatabase(config *configs.Config) Client {
	return &homeBarClient{
		config: config,
	}
}

type homeBarClient struct {
	config *configs.Config
	db     *sqlx.DB
}

func (c *homeBarClient) Database(databaseName string) Database {
	if c.config == nil {
		internal.PrintFatal("config file nil in client", nil)
	}
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		c.config.DatabaseConfig.User,
		c.config.DatabaseConfig.Password,
		c.config.DatabaseConfig.Host,
		c.config.DatabaseConfig.Port,
		databaseName)

	db, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		internal.PrintFatal("", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	c.db = db

	return &homeBarDatabase{
		db: db,
	}
}

func (c *homeBarClient) Connect(ctx context.Context) error {
	return nil
}

func (c *homeBarClient) Disconnect(ctx context.Context) error {
	if c.db == nil {
		internal.PrintFatal("db is nil in client", nil)
	}

	return c.db.Close()
}

func (c *homeBarClient) Ping(ctx context.Context) error {
	if c.db == nil {
		internal.PrintFatal("db is nil in client", nil)
	}

	return c.db.PingContext(ctx)
}

type homeBarDatabase struct {
	db *sqlx.DB
}

func (d *homeBarDatabase) Collection(name string) Collection {
	return &homeBarCollection{
		name: name,
		db:   d.db,
	}
}

type homeBarCollection struct {
	name string
	db   *sqlx.DB
}

func (col *homeBarCollection) UpsertOne(opts ...UpsertionOption) (interface{}, error) {
	if col.db == nil {
		return 0, domain.NewCustomError("db is nil in collection")
	}

	var options UpsertionOptions
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			internal.PrintWarning("UpsertOne options error: %s", err.Error())
		}
	}

	if options.FieldsValues == nil {
		return 0, domain.NewCustomError("UpsertOne FieldsValues nil")
	}

	keys := internal.GetKeys(options.FieldsValues)

	arguments := make([]any, 0, len(keys))

	fields := strings.Join(keys, ", ")
	values := strings.Repeat("?, ", len(options.FieldsValues)-1) + "?"

	updateFields := make([]string, 0, len(keys))
	for i := range keys {
		updateFields = append(updateFields, fmt.Sprintf("%s = VALUES(%s)", keys[i], keys[i]))
		arguments = append(arguments, internal.GetActualInterfaceValue(options.FieldsValues[keys[i]]))
	}

	update := strings.Join(updateFields, ", ")

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s",
		col.name, fields, values, update)

	result, err := col.db.Exec(query, arguments...)
	if err != nil {
		return 0, domain.NewCustomError(fmt.Sprintf("UpsertOne query exec error: %s", err.Error()))
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, domain.NewCustomError(fmt.Sprintf("UpsertOne LastInsertId error: %s", err.Error()))
	}

	return id, nil
}

func (col *homeBarCollection) FindMany(opts ...SelectionOption) (MultiResult, error) {
	if col.db == nil {
		return nil, domain.NewCustomError("DB is nil in collection")
	}

	var options SelectionOptions
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			internal.PrintWarning("FindMany options error: %s", err.Error())
		}
	}

	selectionFields := "*"
	if options.Fields != nil {
		selectionFields = strings.Join(options.Fields, ", ")
	}

	whereClause := ""
	if options.WhereClause != nil {
		whereClause = fmt.Sprintf(" WHERE %s", *options.WhereClause)
	}

	groupByClause := ""
	if options.GroupByFields != nil {
		selectionFields = fmt.Sprintf(" GROUP BY (%s)", strings.Join(options.Fields, ", "))
	}

	joinClause := ""
	if options.JoinClause != nil {
		for pair := options.JoinClause.Oldest(); pair != nil; pair = pair.Next() {
			joinClause += fmt.Sprintf(" JOIN %s ON %s", pair.Key, pair.Value)
		}
	}

	take := ""
	skip := ""
	if options.Pagination != nil {
		take = fmt.Sprintf(" LIMIT %d", options.Pagination.Take)
		skip = fmt.Sprintf(" OFFSET %d", options.Pagination.Skip)
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s%s%s%s%s%s",
		selectionFields,
		col.name,
		joinClause,
		whereClause,
		groupByClause,
		take,
		skip)

	rows, err := col.db.Queryx(query)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("FindMany error: %s", err.Error()))
	}

	return &homeBarMultiResult{Rows: rows}, nil
}

func (col *homeBarCollection) FindOne(opts ...SelectionOption) (SingleResult, error) {
	if col.db == nil {
		return nil, domain.NewCustomError("DB is nil in collection")
	}

	var options SelectionOptions
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			internal.PrintWarning("FindOne options error: %s", err.Error())
		}
	}

	selectionFields := "*"
	if options.Fields != nil {
		selectionFields = strings.Join(options.Fields, ", ")
	}

	whereClause := ""
	if options.WhereClause != nil {
		whereClause = fmt.Sprintf(" WHERE %s", *options.WhereClause)
	}

	groupByClause := ""
	if options.GroupByFields != nil {
		selectionFields = fmt.Sprintf(" GROUP BY (%s)", strings.Join(options.Fields, ", "))
	}

	joinClause := ""
	if options.JoinClause != nil {
		for pair := options.JoinClause.Oldest(); pair != nil; pair = pair.Next() {
			joinClause += fmt.Sprintf(" JOIN %s ON %s", pair.Key, pair.Value)
		}
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s%s%s%s",
		selectionFields,
		col.name,
		joinClause,
		whereClause,
		groupByClause)

	return &homeBarSingleResult{
		Row: col.db.QueryRowx(query),
	}, nil
}

func (col *homeBarCollection) DeleteOne(opts ...DeletionOption) (interface{}, error) {
	if col.db == nil {
		return nil, domain.NewCustomError("DB is nil in collection")
	}

	var options DeletionOptions
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			internal.PrintWarning("DeleteOne options error: %s", err.Error())
		}
	}

	whereClause := ""
	if options.WhereClause != nil {
		whereClause = fmt.Sprintf(" WHERE %s", *options.WhereClause)
	}

	query := fmt.Sprintf("DELETE FROM %s%s", col.name, whereClause)

	result, err := col.db.Exec(query)
	if err != nil {
		return 0, domain.NewCustomError(fmt.Sprintf("DeleteOne query exec error: %s", err.Error()))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, domain.NewCustomError(fmt.Sprintf("UpsertOne LastInsertId error: %s", err.Error()))
	}

	return rowsAffected, nil
}

type homeBarSingleResult struct {
	Row *sqlx.Row
}

func (h *homeBarSingleResult) Decode(target interface{}) error {
	if h.Row == nil {
		return domain.NewCustomError("FindOne row is nil")
	}

	err := h.Row.StructScan(target)
	if err != nil {
		if err == sql.ErrNoRows {
			internal.PrintWarning("%s returned no rows", "Query")
		} else {
			return domain.NewCustomError(fmt.Sprintf("Decode StructScan error: %s", err.Error()))
		}
	}

	return nil
}

type homeBarMultiResult struct {
	Rows *sqlx.Rows
}

func (h *homeBarMultiResult) Decode(target interface{}) ([]interface{}, error) {
	if h.Rows == nil {
		return nil, domain.NewCustomError("FindMany rows is nil")
	}

	result := make([]interface{}, 0)
	for h.Rows.Next() {
		err := h.Rows.StructScan(target)
		if err != nil {
			if err == sql.ErrNoRows {
				internal.PrintWarning("%s returned no rows", "Query")
			} else {
				return nil, domain.NewCustomError(fmt.Sprintf("Decode StructScan error: %s", err.Error()))
			}
		}
		result = append(result, internal.Clone(target))
	}

	defer h.Rows.Close()

	return result, nil
}
