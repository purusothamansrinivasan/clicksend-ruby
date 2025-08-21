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

func TransferacontactHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		from_list_idVal, ok := args["from_list_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: from_list_id"), nil
		}
		from_list_id, ok := from_list_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: from_list_id"), nil
		}
		contact_idVal, ok := args["contact_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: contact_id"), nil
		}
		contact_id, ok := contact_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: contact_id"), nil
		}
		to_list_idVal, ok := args["to_list_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: to_list_id"), nil
		}
		to_list_id, ok := to_list_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: to_list_id"), nil
		}
		url := fmt.Sprintf("%s/lists/%s/contacts/%s/%s", cfg.BaseURL, from_list_id, contact_id, to_list_id)
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

func CreateTransferacontactTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_lists_from_list_id_contacts_contact_id_to_list_id",
		mcp.WithDescription("Transfer a Contact"),
		mcp.WithString("from_list_id", mcp.Required(), mcp.Description("From list id.")),
		mcp.WithString("contact_id", mcp.Required(), mcp.Description("Contact ID.")),
		mcp.WithString("to_list_id", mcp.Required(), mcp.Description("To list id.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    TransferacontactHandler(cfg),
	}
}
