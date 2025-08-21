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

func UpdatepostreturnaddressHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		return_address_idVal, ok := args["return_address_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: return_address_id"), nil
		}
		return_address_id, ok := return_address_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: return_address_id"), nil
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
		url := fmt.Sprintf("%s/post/return-addresses/%s", cfg.BaseURL, return_address_id)
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

func CreateUpdatepostreturnaddressTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_post_return-addresses_return_address_id",
		mcp.WithDescription("Update Post Return Address"),
		mcp.WithString("return_address_id", mcp.Required(), mcp.Description("Your return address id.")),
		mcp.WithString("address_name", mcp.Required(), mcp.Description("Input parameter: Your address name.")),
		mcp.WithString("address_postal_code", mcp.Required(), mcp.Description("Input parameter: Your address postal code.")),
		mcp.WithString("address_state", mcp.Required(), mcp.Description("Input parameter: Your address state.")),
		mcp.WithString("address_city", mcp.Required(), mcp.Description("Input parameter: Your address city.")),
		mcp.WithString("address_country", mcp.Required(), mcp.Description("Input parameter: Two-letter country code defined in ISO 3166.")),
		mcp.WithString("address_line_1", mcp.Required(), mcp.Description("Input parameter: Your address line 1.")),
		mcp.WithString("address_line_2", mcp.Description("Input parameter: Your address line 2.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    UpdatepostreturnaddressHandler(cfg),
	}
}
