package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	_ "modernc.org/sqlite"
)

func errorResult(msg string, err error) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg + ": " + err.Error()},
		},
	}, nil, nil
}

func ConnectDB() (*sql.DB, error) {
	path := os.Getenv("STORE_DB_PATH")
	if path == "" {
		return nil, fmt.Errorf("STORE_DB_PATH is not set")
	}
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

type NoInput struct{}

func list_tables(ctx context.Context, req *mcp.CallToolRequest, _ NoInput) (*mcp.CallToolResult, any, error) {
	db, err := ConnectDB()
	if err != nil {
		return errorResult("Could not connect to the database", err)
	}
	defer db.Close()

	query := `
	SELECT name
	FROM sqlite_master
	WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
	ORDER BY name
	`
	rows, err := db.Query(query)
	if err != nil {
		return errorResult("Could not list tables", err)
	}
	defer rows.Close()

	result := []string{}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return errorResult("Could not read a table name", err)
		}
		result = append(result, name)
	}
	if err := rows.Err(); err != nil {
		return errorResult("Could not finish reading tables", err)
	}

	if len(result) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No tables found in the database."},
			},
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: strings.Join(result, "\n")},
		},
	}, nil, nil
}
