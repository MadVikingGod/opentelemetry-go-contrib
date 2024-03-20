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

package semconv

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
)

func TestDupTraceRequest(t *testing.T) {
	t.Setenv("OTEL_HTTP_CLIENT_COMPATIBILITY_MODE", "http/dup")
	serv := NewHTTPServer()
	want := func(req testServerReq) []attribute.KeyValue {
		return []attribute.KeyValue{
			attribute.String("http.method", "GET"),
			attribute.String("http.request.method", "GET"),
			attribute.String("http.scheme", "http"),
			attribute.String("url.scheme", "http"),
			attribute.String("net.host.name", req.hostname),
			attribute.String("server.address", req.hostname),
			attribute.Int("net.host.port", req.serverPort),
			attribute.Int("server.port", req.serverPort),
			attribute.String("net.sock.peer.addr", req.peerAddr),
			attribute.String("network.peer.address", req.peerAddr),
			attribute.Int("net.sock.peer.port", req.peerPort),
			attribute.Int("network.peer.port", req.peerPort),
			attribute.String("user_agent.original", "Go-http-client/1.1"),
			attribute.String("http.client_ip", req.clientIP),
			attribute.String("client.address", req.clientIP),
			attribute.String("net.protocol.version", "1.1"),
			attribute.String("network.protocol.version", "1.1"),
			attribute.String("http.target", "/"),
			attribute.String("url.path", "/"),
		}
	}
	testTraceRequest(t, serv, want)
}

func TestDupMethod(t *testing.T) {
	testCases := []struct {
		method string
		n      int
		want   []attribute.KeyValue
	}{
		{
			method: http.MethodPost,
			n:      2,
			want: []attribute.KeyValue{
				attribute.String("http.method", "POST"),
				attribute.String("http.request.method", "POST"),
			},
		},
		{
			method: "Put",
			n:      3,
			want: []attribute.KeyValue{
				attribute.String("http.method", "Put"),
				attribute.String("http.request.method", "PUT"),
				attribute.String("http.request.method_original", "Put"),
			},
		},
		{
			method: "Unknown",
			n:      3,
			want: []attribute.KeyValue{
				attribute.String("http.method", "Unknown"),
				attribute.String("http.request.method", "GET"),
				attribute.String("http.request.method_original", "Unknown"),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.method, func(t *testing.T) {
			attrs := make([]attribute.KeyValue, 5)
			n := dupHTTPServer{}.method(tt.method, attrs[1:])
			require.Equal(t, tt.n, n, "Length doesn't match")
			require.ElementsMatch(t, tt.want, attrs[1:n+1])
		})
	}
}

func TestDupTraceResponse(t *testing.T) {
	t.Setenv("OTEL_HTTP_CLIENT_COMPATIBILITY_MODE", "http/dup")
	serv := NewHTTPServer()
	want := []attribute.KeyValue{
		attribute.Int("http.request_content_length", 701),
		attribute.Int("http.request.body.size", 701),
		attribute.String("http.read_error", "read error"),
		attribute.Int("http.response_content_length", 802),
		attribute.Int("http.response.body.size", 802),
		attribute.String("http.write_error", "write error"),
		attribute.Int("http.status_code", 200),
		attribute.Int("http.response.status_code", 200),
	}
	testTraceResponse(t, serv, want)
}

func TestDupClient(t *testing.T) {
	t.Setenv("OTEL_HTTP_CLIENT_COMPATIBILITY_MODE", "http/dup")
	client := NewHTTPClient()
	want := []attribute.KeyValue{
		// Old Attributes
		attribute.String("http.method", "PoST"),
		attribute.String("http.url", "https://fake.url.local:8080/path"),
		attribute.String("net.peer.name", "fake.url.local"),
		attribute.Int("net.peer.port", 8080),
		attribute.Int64("http.request_content_length", 4),
		// New Attributes
		attribute.String("http.request.method", "POST"),
		attribute.String("http.request.method_original", "PoST"),
		attribute.String("network.peer.address", "fake.url.local"),
		attribute.String("server.address", "fake.url.local"),
		attribute.Int("network.peer.port", 8080),
		attribute.Int("server.port", 8080),
		attribute.String("url.full", "https://fake.url.local:8080/path"),
		attribute.Int("http.request.body.size", 4),
		// Common Attributes
		attribute.String("user_agent.original", "http-test-client"),
	}
	testClientTraceRequest(t, client, want)

	// want = []attribute.KeyValue{
	// 	attribute.Int64("http.status_code", 201),
	// 	attribute.Int64("http.response_content_length", 397),
	// }
	// testClientTraceResponse(t, client, want)
}
