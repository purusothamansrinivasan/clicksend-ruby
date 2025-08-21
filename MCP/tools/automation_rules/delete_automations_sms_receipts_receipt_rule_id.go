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

func Delete_automations_sms_receipts_receipt_rule_idHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		receipt_rule_idVal, ok := args["receipt_rule_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: receipt_rule_id"), nil
		}
		receipt_rule_id, ok := receipt_rule_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: receipt_rule_id"), nil
		}
		url := fmt.Sprintf("%s/automations/sms/receipts/%s", cfg.BaseURL, receipt_rule_id)
		req, err := http.NewRequest("DELETE", url, nil)
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

func CreateDelete_automations_sms_receipts_receipt_rule_idTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("delete_automations_sms_receipts_receipt_rule_id",
		mcp.WithDescription("Delete a rule"),
		mcp.WithString("receipt_rule_id", mcp.Required(), mcp.Description("Receipt Rule ID.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Delete_automations_sms_receipts_receipt_rule_idHandler(cfg),
	}
}
