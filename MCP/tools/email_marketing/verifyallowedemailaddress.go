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

func VerifyallowedemailaddressHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		email_address_idVal, ok := args["email_address_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: email_address_id"), nil
		}
		email_address_id, ok := email_address_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: email_address_id"), nil
		}
		activation_tokenVal, ok := args["activation_token"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: activation_token"), nil
		}
		activation_token, ok := activation_tokenVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: activation_token"), nil
		}
		url := fmt.Sprintf("%s/email/address-verify/%s/verify/%s", cfg.BaseURL, email_address_id, activation_token)
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

func CreateVerifyallowedemailaddressTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_email_address-verify_email_address_id_verify_activation_token",
		mcp.WithDescription("Verify Allowed Email Address"),
		mcp.WithString("email_address_id", mcp.Required(), mcp.Description("The email address id you want to access.")),
		mcp.WithString("activation_token", mcp.Required(), mcp.Description("6E8B-4FDB-99A7-7ED08DF97BCC (required, string) - Your activation token.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    VerifyallowedemailaddressHandler(cfg),
	}
}
