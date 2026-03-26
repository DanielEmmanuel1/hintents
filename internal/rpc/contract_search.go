// Copyright 2026 Erst Users
// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	defaultContractSearchLimit = 20
	maxContractSearchPages     = 5
)

type SearchContractsOptions struct {
	Query      string
	HorizonURL string
	Limit      int
	Timeout    time.Duration
}

type ContractSummary struct {
	ID                 string
	Symbol             string
	Creator            string
	LastModifiedLedger int64
	LastModifiedTime   string
}

type horizonContractsPage struct {
	Embedded struct {
		Records []map[string]any `json:"records"`
	} `json:"_embedded"`
	Links struct {
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
	} `json:"_links"`
}

func SearchContracts(ctx context.Context, options SearchContractsOptions) ([]ContractSummary, error) {
	query := strings.TrimSpace(options.Query)
	if query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	horizonURL := strings.TrimRight(strings.TrimSpace(options.HorizonURL), "/")
	if horizonURL == "" {
		horizonURL = strings.TrimRight(TestnetHorizonURL, "/")
	}

	limit := options.Limit
	if limit <= 0 {
		limit = defaultContractSearchLimit
	}

	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}

	client := &http.Client{Timeout: timeout}
	nextURL := buildContractsURL(horizonURL)
	results := make([]ContractSummary, 0, limit)

	for page := 0; page < maxContractSearchPages && nextURL != ""; page++ {
		pageResults, next, err := fetchContractPage(ctx, client, nextURL, query, limit-len(results))
		if err != nil {
			return nil, err
		}
		results = append(results, pageResults...)
		if len(results) >= limit {
			return results[:limit], nil
		}
		nextURL = next
	}

	return results, nil
}

func buildContractsURL(horizonURL string) string {
	u, _ := url.Parse(horizonURL)
	u.Path = path.Join(u.Path, "contracts")
	values := u.Query()
	values.Set("limit", strconv.Itoa(200))
	values.Set("order", "desc")
	u.RawQuery = values.Encode()
	return u.String()
}

func fetchContractPage(
	ctx context.Context,
	client *http.Client,
	pageURL string,
	query string,
	remaining int,
) ([]ContractSummary, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("fetch contracts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("contracts endpoint returned HTTP %d", resp.StatusCode)
	}

	var page horizonContractsPage
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, "", fmt.Errorf("decode contracts response: %w", err)
	}

	matches := make([]ContractSummary, 0, remaining)
	for _, raw := range page.Embedded.Records {
		summary := normalizeContractSummary(raw)
		if matchesQuery(summary, query) {
			matches = append(matches, summary)
			if len(matches) >= remaining {
				break
			}
		}
	}

	return matches, page.Links.Next.Href, nil
}

func normalizeContractSummary(raw map[string]any) ContractSummary {
	summary := ContractSummary{
		ID:                 firstNonEmpty(stringValue(raw["id"]), stringValue(raw["contract_id"])),
		Symbol:             firstNonEmpty(stringValue(raw["symbol"]), stringValue(raw["asset_code"])),
		Creator:            firstNonEmpty(stringValue(raw["creator"]), stringValue(raw["creator_account"]), stringValue(raw["sponsor"])),
		LastModifiedLedger: intValue(raw["last_modified_ledger"], raw["last_modified_ledger_seq"]),
		LastModifiedTime:   stringValue(raw["last_modified_time"]),
	}
	return summary
}

func matchesQuery(summary ContractSummary, query string) bool {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return false
	}
	fields := []string{summary.ID, summary.Symbol, summary.Creator}
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), q) {
			return true
		}
	}
	return false
}

func stringValue(v any) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(s)
}

func intValue(values ...any) int64 {
	for _, v := range values {
		switch n := v.(type) {
		case float64:
			return int64(n)
		case int64:
			return n
		case int:
			return int64(n)
		case json.Number:
			if parsed, err := n.Int64(); err == nil {
				return parsed
			}
		}
	}
	return 0
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
