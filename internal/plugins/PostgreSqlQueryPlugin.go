package plugins

import (
    "fmt"
    "database/sql"
    _ "github.com/lib/pq"
)

var pluginName = "PostgreSqlQueryPlugin"


type PostgreSqlQueryPlugin struct {
    ConnectionString string 
    Query            string
	PluginName       string
}

func (p PostgreSqlQueryPlugin) Name() string {
    return pluginName
}

func (p PostgreSqlQueryPlugin) Collect() (map[string]interface{}, error) {
    pluginName = p.PluginName
	
    db, err := sql.Open("postgres", p.ConnectionString)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }
    defer db.Close()

    rows, err := db.Query(p.Query)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %v", err)
    }
    defer rows.Close()

    columns, err := rows.Columns()
    if err != nil {
        return nil, fmt.Errorf("failed to get columns: %v", err)
    }

    results := make([]map[string]interface{}, 0)
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
        results = append(results, row)
    }

    return map[string]interface{}{
        "rows": results,
    }, nil
}
