package plugins

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

var pluginName = "PostgreSqlQueryPlugin"

type PostgreSqlQueryPlugin struct {
	ConnectionString string
	Query            string
	InstanceName     string
	IsSingleValue    bool
}

func (plugin PostgreSqlQueryPlugin) Name() string {
	return pluginName
}

func (plugin PostgreSqlQueryPlugin) Collect() (map[string]interface{}, error) {
	pluginName = plugin.InstanceName

	db, err := sql.Open("postgres", plugin.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query(plugin.Query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

	results, err := plugin.ProcessData(rows, columns)
	return results, err
}

func (plugin PostgreSqlQueryPlugin) ProcessData(rows *sql.Rows, columns []string) (map[string]interface{}, error) {
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

		if plugin.IsSingleValue {
			return row, nil
		}

		results[strconv.Itoa(rowNumber)] = row
		rowNumber++
	}
	return results, nil
}
