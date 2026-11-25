package version

import (
	"testing"

	localv1 "github.com/gyrel/sandboxes-api/gen/go/sandboxes/local/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Approval disclosure messages must never carry secret material. These field
// names are the ones a careless extension would most likely add.
var forbiddenDisclosureFields = []protoreflect.Name{"value", "secret", "token", "password", "header_value", "placeholder"}

func TestApprovalServiceContractIsAdditive(t *testing.T) {
	service := localv1.File_sandboxes_local_v1_sandbox_proto.Services().ByName("ApprovalService")
	if service == nil {
		t.Fatal("ApprovalService is missing from the generated descriptor")
	}
	for _, method := range []protoreflect.Name{"List", "Inspect", "Respond"} {
		if service.Methods().ByName(method) == nil {
			t.Fatalf("ApprovalService.%s RPC is missing from the generated descriptor", method)
		}
	}

	approval := (&localv1.Approval{}).ProtoReflect().Descriptor()
	for name, number := range map[string]protoreflect.FieldNumber{
		"id": 1, "kind": 2, "state": 3, "sandbox_id": 4, "action_summary": 5,
		"destination_host": 6, "destination_path": 7, "headers": 8,
		"env_names": 9, "env_names_truncated": 10, "cwd_summary": 11,
		"created_at_unix": 12,
	} {
		field := approval.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("Approval.%s field %d is missing or renumbered", name, number)
		}
	}
	for _, descriptor := range []protoreflect.MessageDescriptor{approval, (&localv1.ApprovalHeader{}).ProtoReflect().Descriptor()} {
		for _, forbidden := range forbiddenDisclosureFields {
			if descriptor.Fields().ByName(forbidden) != nil {
				t.Fatalf("%s carries forbidden secret-bearing field %q", descriptor.FullName(), forbidden)
			}
		}
	}

	respond := (&localv1.ApprovalServiceRespondRequest{}).ProtoReflect().Descriptor()
	for name, number := range map[string]protoreflect.FieldNumber{
		"approval_id": 1, "decision": 2, "persist_scope": 3,
	} {
		field := respond.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("ApprovalServiceRespondRequest.%s field %d is missing or renumbered", name, number)
		}
	}
	if respond.Fields().ByName("approved") != nil {
		t.Fatal("ApprovalServiceRespondRequest accepts a caller-provided approval boolean")
	}
}

func TestApprovalContractRoundTrip(t *testing.T) {
	want := &localv1.ApprovalServiceListResponse{Approvals: []*localv1.Approval{{
		Id:                "approval-test",
		Kind:              localv1.ApprovalKind_APPROVAL_KIND_MCP_CALL,
		State:             localv1.ApprovalState_APPROVAL_STATE_PENDING,
		SandboxId:         "sandbox-test",
		ActionSummary:     "call fixture-pong on fixture",
		DestinationHost:   "mcp.example.test",
		DestinationPath:   "/v1/mcp",
		Headers:           []*localv1.ApprovalHeader{{Name: "X-API-Key", ValueClass: localv1.HeaderValueClass_HEADER_VALUE_CLASS_VAULT, VaultReferenceIds: []string{"vault-ref-1"}}},
		EnvNames:          []string{"PATH"},
		EnvNamesTruncated: false,
		CwdSummary:        "<workspace>/app",
		CreatedAtUnix:     1,
	}}}
	encoded, err := proto.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got localv1.ApprovalServiceListResponse
	if err := proto.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(want, &got) {
		t.Fatalf("approval round-trip mismatch: got %v, want %v", &got, want)
	}

	request := &localv1.ApprovalServiceRespondRequest{
		ApprovalId:   "approval-test",
		Decision:     localv1.ApprovalDecision_APPROVAL_DECISION_APPROVE,
		PersistScope: localv1.ApprovalPersistScope_APPROVAL_PERSIST_SCOPE_EXACT,
	}
	encoded, err = proto.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	var decoded localv1.ApprovalServiceRespondRequest
	if err := proto.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(request, &decoded) {
		t.Fatalf("respond round-trip mismatch: got %v, want %v", &decoded, request)
	}
}
