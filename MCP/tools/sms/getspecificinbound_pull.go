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

func Getspecificinbound_pullHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		outbound_message_idVal, ok := args["outbound_message_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: outbound_message_id"), nil
		}
		outbound_message_id, ok := outbound_message_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: outbound_message_id"), nil
		}
		url := fmt.Sprintf("%s/sms/inbound/%s", cfg.BaseURL, outbound_message_id)
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

func CreateGetspecificinbound_pullTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_sms_inbound_outbound_message_id",
		mcp.WithDescription("Get Specific Inbound - Pull"),
		mcp.WithString("outbound_message_id", mcp.Required(), mcp.Description("Message ID of the original outbound message, to which the inbound message is a reply. Must be a valid GUID.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Getspecificinbound_pullHandler(cfg),
	}
}
