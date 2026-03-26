// Copyright 2026 Erst Users
// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchContracts_MatchesBySymbolCreatorAndPartialID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/contracts" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"_embedded": map[string]any{
				"records": []map[string]any{
					{
						"id":                   "CAAA111",
						"symbol":               "USDC",
						"creator":              "GCREATORA",
						"last_modified_ledger": 123,
						"last_modified_time":   "2026-03-20T01:02:03Z",
					},
					{
						"id":                   "CBBB222",
						"symbol":               "MEME",
						"creator":              "GCREATORB",
						"last_modified_ledger": 456,
						"last_modified_time":   "2026-03-21T01:02:03Z",
					},
				},
			},
			"_links": map[string]any{
				"next": map[string]any{"href": ""},
			},
		})
	}))
	defer server.Close()

	ctx := context.Background()

	bySymbol, err := SearchContracts(ctx, SearchContractsOptions{
		Query:      "usdc",
		HorizonURL: server.URL,
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("search by symbol failed: %v", err)
	}
	if len(bySymbol) != 1 || bySymbol[0].ID != "CAAA111" {
		t.Fatalf("unexpected symbol results: %+v", bySymbol)
	}

	byCreator, err := SearchContracts(ctx, SearchContractsOptions{
		Query:      "GCREATORB",
		HorizonURL: server.URL,
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("search by creator failed: %v", err)
	}
	if len(byCreator) != 1 || byCreator[0].ID != "CBBB222" {
		t.Fatalf("unexpected creator results: %+v", byCreator)
	}

	byPartialID, err := SearchContracts(ctx, SearchContractsOptions{
		Query:      "caaa",
		HorizonURL: server.URL,
		Limit:      10,
	})
	if err != nil {
		t.Fatalf("search by id failed: %v", err)
	}
	if len(byPartialID) != 1 || byPartialID[0].ID != "CAAA111" {
		t.Fatalf("unexpected partial id results: %+v", byPartialID)
	}
}

func TestSearchContracts_RespectsLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"_embedded": map[string]any{
				"records": []map[string]any{
					{"id": "CA1", "symbol": "ABC", "creator": "G1", "last_modified_ledger": 10},
					{"id": "CA2", "symbol": "ABC", "creator": "G2", "last_modified_ledger": 20},
					{"id": "CA3", "symbol": "ABC", "creator": "G3", "last_modified_ledger": 30},
				},
			},
			"_links": map[string]any{
				"next": map[string]any{"href": ""},
			},
		})
	}))
	defer server.Close()

	results, err := SearchContracts(context.Background(), SearchContractsOptions{
		Query:      "abc",
		HorizonURL: server.URL,
		Limit:      2,
	})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}
