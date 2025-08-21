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

func CreateanewcontactHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		list_idVal, ok := args["list_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: list_id"), nil
		}
		list_id, ok := list_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: list_id"), nil
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
		url := fmt.Sprintf("%s/lists/%s/contacts", cfg.BaseURL, list_id)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
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

func CreateCreateanewcontactTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_lists_list_id_contacts",
		mcp.WithDescription("Create a new contact"),
		mcp.WithString("list_id", mcp.Required(), mcp.Description("Your contact list id where your contact be associated.")),
		mcp.WithString("last_name", mcp.Description("Input parameter: Contact lastname.")),
		mcp.WithString("custom_4", mcp.Description("Input parameter: Contact custom 4 text.")),
		mcp.WithString("email", mcp.Description("Input parameter: Contact email.")),
		mcp.WithString("address_country", mcp.Description("Input parameter: Contact two-letter country code defined in ISO 3166.")),
		mcp.WithString("address_postal_code", mcp.Description("Input parameter: Contact postal code.")),
		mcp.WithString("custom_1", mcp.Description("Input parameter: Contact custom 1 text.")),
		mcp.WithString("custom_2", mcp.Description("Input parameter: Contact custom 2 text.")),
		mcp.WithString("organization_name", mcp.Description("Input parameter: Your organization name.")),
		mcp.WithString("address_line_2", mcp.Description("Input parameter: Contact address line 2.")),
		mcp.WithString("address_state", mcp.Description("Input parameter: Contact state.")),
		mcp.WithString("first_name", mcp.Description("Input parameter: Contact firstname.")),
		mcp.WithString("custom_3", mcp.Description("Input parameter: Contact custom 3 text.")),
		mcp.WithString("fax_number", mcp.Description("Input parameter: Contact fax number.")),
		mcp.WithString("phone_number", mcp.Required(), mcp.Description("Input parameter: Contact phone number in E.164 format.")),
		mcp.WithString("address_city", mcp.Description("Input parameter: Contact city.")),
		mcp.WithString("address_line_1", mcp.Description("Input parameter: Contact address line 1.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    CreateanewcontactHandler(cfg),
	}
}
