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

func Put_automations_fax_inbound_inbound_rule_idHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		inbound_rule_idVal, ok := args["inbound_rule_id"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: inbound_rule_id"), nil
		}
		inbound_rule_id, ok := inbound_rule_idVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: inbound_rule_id"), nil
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
		url := fmt.Sprintf("%s/automations/fax/inbound/%s", cfg.BaseURL, inbound_rule_id)
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

func CreatePut_automations_fax_inbound_inbound_rule_idTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_automations_fax_inbound_inbound_rule_id",
		mcp.WithDescription("Update a rule"),
		mcp.WithString("inbound_rule_id", mcp.Required(), mcp.Description("Fax inbound rule id")),
		mcp.WithString("rule_name", mcp.Required(), mcp.Description("Input parameter: Rule Name")),
		mcp.WithString("action", mcp.Required(), mcp.Description("Input parameter: Action")),
		mcp.WithString("action_address", mcp.Required(), mcp.Description("Input parameter: Action Address")),
		mcp.WithString("dedicated_number", mcp.Required(), mcp.Description("Input parameter: Decicated Number")),
		mcp.WithString("enabled", mcp.Required(), mcp.Description("Input parameter: Enable")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Put_automations_fax_inbound_inbound_rule_idHandler(cfg),
	}
}
