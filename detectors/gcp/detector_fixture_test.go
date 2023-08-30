// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gcp // import "go.opentelemetry.io/contrib/detectors/gcp"

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type testCase struct {
	Description        string            `json:"description"`
	EnvVars            map[string]string `json:"envVars",omitempty`
	Metadata           map[string]string `json:"metadata",omitempty`
	ExpectedAttributes map[string]string `json:"expectedAttributes",omitempty`
}

type testCases struct {
	TestCases []testCase `json:"testCases",omitempty`
}

func TestFixtureDetection(t *testing.T) {
	// Load TestCases
	bytes, err := os.ReadFile("test_cases.json")
	assert.NoError(t, err)
	var tcs testCases
	err = json.Unmarshal(bytes, &tcs)
	assert.NoError(t, err)
	// Test fixtures
	for _, tc := range tcs.TestCases {
		t.Run(tc.Description, func(t *testing.T) {
			// Set environment variables from the fixtures
			for k, v := range tc.EnvVars {
				t.Setenv(k, v)
			}
			// Handle metadata server requests
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				meta, ok := tc.Metadata[r.URL.Path]
				if !ok {
					t.Errorf("metadata path %v not found", r.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(meta))
			}))
			defer server.Close()

			trimmedURL, ok := strings.CutPrefix(server.URL, "http://")
			assert.True(t, ok)
			// Point detector at our server
			t.Setenv("GCE_METADATA_HOST", trimmedURL)
			// Check detected resources
			res, err := NewDetector().Detect(context.Background())
			assert.NoError(t, err)
			var attrs []attribute.KeyValue
			for expectedKey, expectedVal := range tc.ExpectedAttributes {
				attrs = append(attrs, attribute.String(expectedKey, expectedVal))
			}
			expectedResource := resource.NewWithAttributes(semconv.SchemaURL, attrs...)
			assert.Equal(t, expectedResource, res, "Resource object returned is incorrect")
		})
	}
}
