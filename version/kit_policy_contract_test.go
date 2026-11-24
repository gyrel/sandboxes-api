package version

import (
	"testing"

	localv1 "github.com/rheic/sandboxes-api/gen/go/sandboxes/local/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// KitPolicyService is an additive contract: it introduces a new service and new
// messages without changing any existing descriptor. These tests pin the field
// numbers, enum values and redaction invariants so a future edit cannot silently
// renumber or smuggle secret material into the surface.
func TestKitPolicyServiceDescriptorPresence(t *testing.T) {
	service := localv1.File_sandboxes_local_v1_sandbox_proto.Services().ByName("KitPolicyService")
	if service == nil {
		t.Fatal("KitPolicyService is missing from the generated descriptor")
	}
	for _, method := range []protoreflect.Name{"CheckKit", "GetKitPolicyConfig"} {
		if service.Methods().ByName(method) == nil {
			t.Fatalf("KitPolicyService.%s RPC is missing from the generated descriptor", method)
		}
	}
	for _, messageName := range []string{
		"CheckKitRequest", "CheckKitResponse",
		"GetKitPolicyConfigRequest", "GetKitPolicyConfigResponse",
	} {
		if localv1.File_sandboxes_local_v1_sandbox_proto.Messages().ByName(protoreflect.Name(messageName)) == nil {
			t.Fatalf("%s is missing from the generated descriptor", messageName)
		}
	}
	for _, enumName := range []string{"KitSourceKind", "CheckKitReason"} {
		if localv1.File_sandboxes_local_v1_sandbox_proto.Enums().ByName(protoreflect.Name(enumName)) == nil {
			t.Fatalf("%s is missing from the generated descriptor", enumName)
		}
	}
}

func TestKitPolicyContractIsAdditive(t *testing.T) {
	checkRequest := (&localv1.CheckKitRequest{}).ProtoReflect().Descriptor()
	for name, number := range map[string]protoreflect.FieldNumber{
		"source_kind": 1, "source_ref": 2, "sandbox_id": 3, "signature_present": 4,
	} {
		field := checkRequest.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("CheckKitRequest.%s field %d is missing or renumbered", name, number)
		}
	}

	checkResponse := (&localv1.CheckKitResponse{}).ProtoReflect().Descriptor()
	for name, number := range map[string]protoreflect.FieldNumber{
		"allowed": 1, "reason": 2, "message": 3, "checked_settings": 4,
	} {
		field := checkResponse.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("CheckKitResponse.%s field %d is missing or renumbered", name, number)
		}
	}

	configResponse := (&localv1.GetKitPolicyConfigResponse{}).ProtoReflect().Descriptor()
	for name, number := range map[string]protoreflect.FieldNumber{
		"allowed_sources": 1, "allow_local_kits": 2, "require_signature": 3, "managed": 4,
	} {
		field := configResponse.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("GetKitPolicyConfigResponse.%s field %d is missing or renumbered", name, number)
		}
	}

	// Forbidden-field scan: no secret material may be carried by the new
	// messages. checked_settings carries names only.
	for _, descriptor := range []protoreflect.MessageDescriptor{
		checkRequest, checkResponse,
		(&localv1.GetKitPolicyConfigRequest{}).ProtoReflect().Descriptor(),
		configResponse,
	} {
		for _, forbidden := range []protoreflect.Name{"value", "token", "content"} {
			if descriptor.Fields().ByName(forbidden) != nil {
				t.Fatalf("%s carries forbidden field %q", descriptor.FullName(), forbidden)
			}
		}
	}
}

func TestKitPolicyEnumStability(t *testing.T) {
	sourceKind := localv1.File_sandboxes_local_v1_sandbox_proto.Enums().ByName("KitSourceKind")
	if sourceKind == nil {
		t.Fatal("KitSourceKind enum is missing")
	}
	for name, number := range map[string]protoreflect.EnumNumber{
		"KIT_SOURCE_KIND_UNSPECIFIED": 0,
		"KIT_SOURCE_KIND_DIRECTORY":   1,
		"KIT_SOURCE_KIND_ZIP":         2,
		"KIT_SOURCE_KIND_GIT":         3,
		"KIT_SOURCE_KIND_OCI":         4,
	} {
		value := sourceKind.Values().ByName(protoreflect.Name(name))
		if value == nil || value.Number() != number {
			t.Fatalf("KitSourceKind.%s = %v, want %d", name, value, number)
		}
	}

	reason := localv1.File_sandboxes_local_v1_sandbox_proto.Enums().ByName("CheckKitReason")
	if reason == nil {
		t.Fatal("CheckKitReason enum is missing")
	}
	for name, number := range map[string]protoreflect.EnumNumber{
		"CHECK_KIT_REASON_UNSPECIFIED":        0,
		"CHECK_KIT_REASON_OK":                 1,
		"CHECK_KIT_REASON_SOURCE_NOT_ALLOWED": 2,
		"CHECK_KIT_REASON_LOCAL_DISABLED":     3,
		"CHECK_KIT_REASON_SIGNATURE_REQUIRED": 4,
		"CHECK_KIT_REASON_POLICY_UNAVAILABLE": 5,
		"CHECK_KIT_REASON_INVALID_REFERENCE":  6,
	} {
		value := reason.Values().ByName(protoreflect.Name(name))
		if value == nil || value.Number() != number {
			t.Fatalf("CheckKitReason.%s = %v, want %d", name, value, number)
		}
	}
}

func TestKitPolicyContractRoundTrip(t *testing.T) {
	want := &localv1.CheckKitRequest{
		SourceKind:       localv1.KitSourceKind_KIT_SOURCE_KIND_DIRECTORY,
		SourceRef:        "/kits/example",
		SandboxId:        "sandbox-test",
		SignaturePresent: true,
	}
	encoded, err := proto.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got localv1.CheckKitRequest
	if err := proto.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(want, &got) {
		t.Fatalf("check-kit request round-trip mismatch: got %v, want %v", &got, want)
	}

	response := &localv1.CheckKitResponse{
		Allowed:         false,
		Reason:          localv1.CheckKitReason_CHECK_KIT_REASON_POLICY_UNAVAILABLE,
		Message:         "policy unavailable",
		CheckedSettings: []string{"kit.allowedSources", "kit.allowLocalKits", "kit.requireSignature"},
	}
	encoded, err = proto.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	var decoded localv1.CheckKitResponse
	if err := proto.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(response, &decoded) {
		t.Fatalf("check-kit response round-trip mismatch: got %v, want %v", &decoded, response)
	}

	config := &localv1.GetKitPolicyConfigResponse{
		AllowedSources:   []localv1.KitSourceKind{localv1.KitSourceKind_KIT_SOURCE_KIND_OCI, localv1.KitSourceKind_KIT_SOURCE_KIND_GIT},
		AllowLocalKits:   true,
		RequireSignature: false,
		Managed:          true,
	}
	encoded, err = proto.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	var decodedConfig localv1.GetKitPolicyConfigResponse
	if err := proto.Unmarshal(encoded, &decodedConfig); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(config, &decodedConfig) {
		t.Fatalf("kit policy config round-trip mismatch: got %v, want %v", &decodedConfig, config)
	}
}
