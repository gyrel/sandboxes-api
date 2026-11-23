package version

import (
	"testing"

	localv1 "github.com/rheic/sandboxes-api/gen/go/sandboxes/local/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The custom-secret surface reuses SecretSource for resolver sources and the
// existing List/Remove RPCs. Responses expose placement metadata only: the
// placeholder is returned masked and the value never leaves the host.
func TestSecretSetCustomContractIsAdditive(t *testing.T) {
	service := localv1.File_sandboxes_local_v1_sandbox_proto.Services().ByName("SecretService")
	if service == nil || service.Methods().ByName("SetCustom") == nil {
		t.Fatal("SecretService.SetCustom RPC is missing from the generated descriptor")
	}

	request := (&localv1.SecretServiceSetCustomRequest{}).ProtoReflect().Descriptor()
	source := request.Fields().ByName("source")
	if source == nil || source.Number() != 5 || source.Message().FullName() != (&localv1.SecretSource{}).ProtoReflect().Descriptor().FullName() {
		t.Fatal("SetCustom must reuse SecretSource for reference, command, refresh and verify options")
	}
	for name, number := range map[string]protoreflect.FieldNumber{
		"service": 1, "scope": 2, "sandbox_name": 3, "value": 4, "env": 6,
		"placeholder": 7, "hosts": 8,
	} {
		field := request.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("SecretServiceSetCustomRequest.%s field %d is missing or renumbered", name, number)
		}
	}

	summary := (&localv1.SecretSummary{}).ProtoReflect().Descriptor()
	for name, number := range map[string]protoreflect.FieldNumber{
		"env_name": 11, "hosts": 12, "placeholder_masked": 13,
	} {
		field := summary.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("SecretSummary.%s field %d is missing or renumbered", name, number)
		}
	}
	for _, descriptor := range []protoreflect.MessageDescriptor{summary, (&localv1.SecretServiceSetCustomResponse{}).ProtoReflect().Descriptor()} {
		for _, forbidden := range []protoreflect.Name{"value", "placeholder", "secret_value", "password", "token"} {
			if descriptor.Fields().ByName(forbidden) != nil {
				t.Fatalf("%s returns forbidden field %q", descriptor.FullName(), forbidden)
			}
		}
	}
}

func TestSecretSetCustomContractRoundTrip(t *testing.T) {
	want := &localv1.SecretServiceSetCustomRequest{
		Service:     "example",
		SandboxName: "sandbox-test",
		Env:         "API_KEY",
		Placeholder: "sk-{rand}",
		Hosts:       []string{"api.example.test", "*.example.test"},
		Source: &localv1.SecretSource{
			Reference: "op://vault/item/field",
			Refresh:   "1h",
			NoVerify:  true,
		},
	}
	encoded, err := proto.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got localv1.SecretServiceSetCustomRequest
	if err := proto.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(want, &got) {
		t.Fatalf("set-custom round-trip mismatch: got %v, want %v", &got, want)
	}

	response := &localv1.SecretServiceSetCustomResponse{ApprovalRequired: true, ApprovalId: "approval-test"}
	encoded, err = proto.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}
	var decoded localv1.SecretServiceSetCustomResponse
	if err := proto.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(response, &decoded) {
		t.Fatalf("set-custom response round-trip mismatch: got %v, want %v", &decoded, response)
	}
}
