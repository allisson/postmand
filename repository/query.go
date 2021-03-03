package repository

import (
	"log"

	"github.com/allisson/postmand"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

func getQuery(tableName string, getOptions postmand.RepositoryGetOptions) (string, []interface{}) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From(tableName)
	for key, value := range getOptions.Filters {
		sb.Where(sb.Equal(key, value))
	}
	return sb.Build()
}

func listQuery(tableName string, listOptions postmand.RepositoryListOptions) (string, []interface{}) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From(tableName).Limit(listOptions.Limit).Offset(listOptions.Offset)
	for key, value := range listOptions.Filters {
		sb.Where(sb.Equal(key, value))
	}
	if listOptions.OrderBy != "" && listOptions.Order != "" {
		sb.OrderBy(listOptions.OrderBy)
		switch listOptions.Order {
		case "asc", "ASC":
			sb.Asc()
		case "desc", "DESC":
			sb.Desc()
		}
	}
	return sb.Build()
}

func insertQuery(tableName string, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.PostgreSQL)
	ib := theStruct.InsertInto(tableName, structValue)
	return ib.Build()
}

func updateQuery(tableName string, id postmand.ID, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.PostgreSQL)
	ib := theStruct.Update(tableName, structValue)
	ib.Where(ib.Equal("id", id))
	return ib.Build()
}

func rollback(msg string, tx *sqlx.Tx) {
	if err := tx.Rollback(); err != nil {
		log.Printf("%s: unable to rollback: %v\n", msg, err)
	}
}

func commit(msg string, tx *sqlx.Tx) {
	if err := tx.Commit(); err != nil {
		log.Printf("%s: unable to commit: %v\n", msg, err)
		rollback(msg, tx)
	}
}
