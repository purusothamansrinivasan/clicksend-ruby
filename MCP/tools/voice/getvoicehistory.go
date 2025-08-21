package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/clicksend-rest-api-v3/mcp-server/config"
	"github.com/clicksend-rest-api-v3/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetvoicehistoryHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		date_fromVal, ok := args["date_from"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: date_from"), nil
		}
		date_from, ok := date_fromVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: date_from"), nil
		}
		date_toVal, ok := args["date_to"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: date_to"), nil
		}
		date_to, ok := date_toVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: date_to"), nil
		}
		url := fmt.Sprintf("%s/voice/history?date_from=%s&date_to=%s", cfg.BaseURL, date_from, date_to)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// Set authentication based on auth type
		if cfg.BasicAuth != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Basic %s", cfg.BasicAuth))
		}
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateGetvoicehistoryTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_voice_history?date_from=date_from&date_to=date_to",
		mcp.WithDescription("Get Voice History"),
		mcp.WithString("date_from", mcp.Required(), mcp.Description("Timestamp (from) used to show records by date.")),
		mcp.WithString("date_to", mcp.Required(), mcp.Description("Timestamp (to) used to show recrods by date.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetvoicehistoryHandler(cfg),
	}
}
