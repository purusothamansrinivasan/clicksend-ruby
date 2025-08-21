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

func SendfaxHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
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
		url := fmt.Sprintf("%s/fax/send", cfg.BaseURL)
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

func CreateSendfaxTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("post_fax_send",
		mcp.WithDescription("Send Fax"),
		mcp.WithArray("messages", mcp.Required(), mcp.Description("Input parameter: Your messages.")),
		mcp.WithString("source", mcp.Description("Input parameter: Your method of sending e.g. 'wordpress', 'php', 'c#'.")),
		mcp.WithString("from_email", mcp.Description("Input parameter: An email address where the reply should be emailed to.")),
		mcp.WithString("to", mcp.Required(), mcp.Description("Input parameter: Recipient number in E.164 format or local format ([more info](https://help.clicksend.com/SMS/what-format-does-the-recipient-phone-number-need-to-be-in)).")),
		mcp.WithString("country", mcp.Description("Input parameter: Recipient country.")),
		mcp.WithString("file_url", mcp.Required(), mcp.Description("Input parameter: Your URL to your PDF file.")),
		mcp.WithString("list_id", mcp.Description("Input parameter: Your list ID if sending to a whole list. Can be used instead of 'to'.")),
		mcp.WithString("custom_string", mcp.Description("Input parameter: Your reference. Will be passed back with all replies and delivery reports.")),
		mcp.WithString("from", mcp.Description("Input parameter: Your sender id. Must be a valid fax number.")),
		mcp.WithString("schedule", mcp.Description("Input parameter: Leave blank for immediate delivery. Your schedule time as a [unix timestamp](http://help.clicksend.com/what-is-a-unix-timestamp).")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    SendfaxHandler(cfg),
	}
}
