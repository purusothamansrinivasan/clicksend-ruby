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

func GetspecificvoicereceiptHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		message_idVal, ok := args["message_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: message_id"), nil
		}
		message_id, ok := message_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: message_id"), nil
		}
		url := fmt.Sprintf("%s/voice/receipts/%s", cfg.BaseURL, message_id)
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

func CreateGetspecificvoicereceiptTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_voice_receipts_message_id",
		mcp.WithDescription("Get Specific Voice Receipt"),
		mcp.WithString("message_id", mcp.Required(), mcp.Description("3055-45F1-9B79-F2C43509FD16 (string, required) - The voice receipt message id.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetspecificvoicereceiptHandler(cfg),
	}
}
