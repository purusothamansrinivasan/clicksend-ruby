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

func GetcountrypricingHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		countryVal, ok := args["country"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: country"), nil
		}
		country, ok := countryVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: country"), nil
		}
		currencyVal, ok := args["currency"]
		if !ok {
			return mcp.NewToolResultError("Missing required path parameter: currency"), nil
		}
		currency, ok := currencyVal.(string)
		if !ok {
			return mcp.NewToolResultError("Invalid path parameter: currency"), nil
		}
		url := fmt.Sprintf("%s/pricing/%s?currency=%s", cfg.BaseURL, country, currency)
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

func CreateGetcountrypricingTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_pricing_country?currency=currency",
		mcp.WithDescription("Get Country Pricing"),
		mcp.WithString("country", mcp.Required(), mcp.Description("Two-letter representation of the country.")),
		mcp.WithString("currency", mcp.Required(), mcp.Description("Three-letter representation of the currency.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    GetcountrypricingHandler(cfg),
	}
}
