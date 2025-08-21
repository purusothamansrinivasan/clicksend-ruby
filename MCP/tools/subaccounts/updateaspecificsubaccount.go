package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"bytes"

	"github.com/clicksend-rest-api-v3/mcp-server/config"
	"github.com/clicksend-rest-api-v3/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func UpdateaspecificsubaccountHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		subaccount_idVal, ok := args["subaccount_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: subaccount_id"), nil
		}
		subaccount_id, ok := subaccount_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: subaccount_id"), nil
		}
		// Create properly typed request body using the generated schema
		var requestBody map[string]interface{}
		
		// Optimized: Single marshal/unmarshal with JSON tags handling field mapping
		if argsJSON, err := json.Marshal(args); err == nil {
			if err := json.Unmarshal(argsJSON, &requestBody); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to convert arguments to request type: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal arguments: %v", err)), nil
		}
		
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to encode request body", err), nil
		}
		url := fmt.Sprintf("%s/subaccounts/%s", cfg.BaseURL, subaccount_id)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
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

func CreateUpdateaspecificsubaccountTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_subaccounts_subaccount_id",
		mcp.WithDescription("Update a specific subaccount"),
		mcp.WithString("subaccount_id", mcp.Required(), mcp.Description("The subaccount ID you want to access.")),
		mcp.WithString("first_name", mcp.Description("Input parameter: Your firstname.")),
		mcp.WithString("phone_number", mcp.Description("Input parameter: Your phone number in E.164 format.")),
		mcp.WithString("access_settings", mcp.Description("Input parameter: Your access settings flag value, must be 1 or 0.")),
		mcp.WithString("email", mcp.Description("Input parameter: Your new email.")),
		mcp.WithString("access_contacts", mcp.Description("Input parameter: Your access contacts flag value, must be 1 or 0.")),
		mcp.WithString("access_users", mcp.Description("Input parameter: Your access users flag value, must be 1 or 0.")),
		mcp.WithString("access_billing", mcp.Description("Input parameter: Your access billing flag value, must be 1 or 0.")),
		mcp.WithString("last_name", mcp.Description("Input parameter: Your lastname.")),
		mcp.WithString("password", mcp.Description("Input parameter: Your new password.")),
		mcp.WithString("share_campaigns", mcp.Description("Input parameter: Your share campaigns flag value, must be 1 or 0.")),
		mcp.WithString("access_reporting", mcp.Description("Input parameter: Your access reporting flag value, must be 1 or 0.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    UpdateaspecificsubaccountHandler(cfg),
	}
}
