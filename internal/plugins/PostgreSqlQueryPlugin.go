package plugins

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

var pluginName = "PostgreSqlQueryPlugin"

type Query struct {
	Name          string
	QueryText     string
	QueryTimeout  float64
	IsSingleValue bool
}

type PostgreSqlQueryPlugin struct {
	InstanceName     string
	ConnectionString string
	Queries          []interface{}
}

func (plugin PostgreSqlQueryPlugin) Name() string {
	return pluginName
}

func (plugin PostgreSqlQueryPlugin) Collect() (map[string]interface{}, error) {
	pluginName = plugin.InstanceName

	results := make(map[string]interface{})

	db, err := sql.Open("postgres", plugin.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	for _, query := range plugin.Queries {
		data, ok := query.(map[string]interface{})
		if !ok {
			fmt.Println("Invalid data type")
			return results, nil
		}

		queryContext, cancel := context.WithTimeout(context.Background(), time.Duration(data["QueryTimeout"].(float64))*time.Second)
		defer cancel()

		rows, err := db.QueryContext(queryContext, data["QueryText"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %v", err)
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return nil, fmt.Errorf("failed to get columns: %v", err)
		}

		queryResult, err := plugin.ProcessData(rows, columns, data["IsSingleValue"].(bool))
		if err != nil {
			results[data["Name"].(string)] = nil
		}

		results[data["Name"].(string)] = queryResult

	}
	return results, err
}

func (plugin PostgreSqlQueryPlugin) ProcessData(rows *sql.Rows, columns []string, IsSingleValue bool) (map[string]interface{}, error) {
	results := make(map[string]interface{}, 0)
	var rowNumber = 0

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}

		convertByteArraysInMapToStrings(row)

		if IsSingleValue {
			return row, nil
		}

		results[strconv.Itoa(rowNumber)] = row
		rowNumber++
	}
	return results, nil
}

func convertByteArraysInMapToStrings(row map[string]interface{}) {
	for key, value := range row {
		if v, ok := value.([]uint8); ok {
			row[key] = string(v)
		}
	}
}
