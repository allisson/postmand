package repository

import (
	"github.com/allisson/postmand"
	"github.com/huandu/go-sqlbuilder"
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
	if listOptions.OrderBy != "" {
		sb.OrderBy(listOptions.OrderBy)
	}
	return sb.Build()
}

func insertQuery(tableName string, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.PostgreSQL)
	ib := theStruct.InsertInto(tableName, structValue)
	return ib.Build()
}

func updateQuery(tableName string, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.PostgreSQL)
	ib := theStruct.Update(tableName, structValue)
	return ib.Build()
}
