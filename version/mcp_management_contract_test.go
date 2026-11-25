package version

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	localv1 "github.com/gyrel/sandboxes-api/gen/go/sandboxes/local/v1"
)

// Deferred: registry/manifest URLs, static MCP configuration, discovery/DCR
// configuration and OAuth scope overrides need a separate registration contract.
func TestMCPManagementContract(t *testing.T) {
	file := localv1.File_sandboxes_local_v1_sandbox_proto
	service := file.Services().ByName("SandboxService")
	for _, name := range []protoreflect.Name{"MCPLoadServer", "MCPUnloadServer", "InspectMCPServer", "GetMCPAuthStatus", "RemoveMCPAuth", "StartMCPAuth", "MCPRegisterServer", "MCPRemoveServer", "MCPCall"} {
		if service.Methods().ByName(name) == nil {
			t.Errorf("missing RPC %s", name)
		}
	}
	// NB: MCPRegisterServerRequest is deliberately not scanned: it predates this
	// contract and legitimately carries secret_headers (vault placeholder
	// references, values never cross). New messages must not name any field
	// with value/token/secret material.
	for _, name := range []protoreflect.Name{"MCPLoadServerRequest", "MCPLoadServerResponse", "MCPUnloadServerRequest", "MCPUnloadServerResponse", "InspectMCPServerRequest", "InspectMCPServerResponse", "MCPServerDetail", "GetMCPAuthStatusRequest", "GetMCPAuthStatusResponse", "MCPServerAuthStatus", "RemoveMCPAuthRequest", "RemoveMCPAuthResponse", "StartMCPAuthRequest", "StartMCPAuthResponse"} {
		msg := file.Messages().ByName(name)
		if msg == nil {
			t.Errorf("missing message %s", name)
			continue
		}
		for i := 0; i < msg.Fields().Len(); i++ {
			field := msg.Fields().Get(i)
			for _, forbidden := range []string{"value", "token", "secret", "header_value", "client_secret"} {
				if strings.Contains(string(field.Name()), forbidden) {
					t.Errorf("%s.%s contains %s", name, field.Name(), forbidden)
				}
			}
		}
	}
	for _, tc := range []struct {
		name   protoreflect.Name
		values []protoreflect.Name
	}{
		{"MCPAuthStatus", []protoreflect.Name{"MCP_AUTH_STATUS_UNSPECIFIED", "MCP_AUTH_STATUS_NOT_CONFIGURED", "MCP_AUTH_STATUS_READY", "MCP_AUTH_STATUS_PENDING_APPROVAL", "MCP_AUTH_STATUS_EXPIRED", "MCP_AUTH_STATUS_ERROR"}},
		{"MCPGatewayState", []protoreflect.Name{"MCP_GATEWAY_STATE_UNSPECIFIED", "MCP_GATEWAY_STATE_UNLOADED", "MCP_GATEWAY_STATE_LOADED"}},
	} {
		e := file.Enums().ByName(tc.name)
		for i, n := range tc.values {
			v := e.Values().ByName(n)
			if v == nil || v.Number() != protoreflect.EnumNumber(i) {
				t.Errorf("%s.%s must be %d", tc.name, n, i)
			}
		}
	}
}
