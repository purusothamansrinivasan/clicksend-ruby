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

func MarkedvoicereceiptsasreadHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		date_beforeVal, ok := args["date_before"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: date_before"), nil
		}
		date_before, ok := date_beforeVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: date_before"), nil
		}
		url := fmt.Sprintf("%s/voice/receipts-read?date_before=%s", cfg.BaseURL, date_before)
		req, err := http.NewRequest("PUT", url, nil)
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

func CreateMarkedvoicereceiptsasreadTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_voice_receipts-read?date_before=date_before",
		mcp.WithDescription("Marked Voice Receipts as Read"),
		mcp.WithString("date_before", mcp.Required(), mcp.Description("An optional [unix timestamp](http://help.clicksend.com/what-is-a-unix-timestamp) - mark all as read before this timestamp. If not given, all receipts will be marked as read.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    MarkedvoicereceiptsasreadHandler(cfg),
	}
}
