package version

import (
	"testing"

	localv1 "github.com/rheic/sandboxes-api/gen/go/sandboxes/local/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// The organization policy view is additive: existing fields and messages keep
// their numbers, so older clients can ignore the new RPC and the new origin
// field without changing behavior.
func TestOrgPolicyContractIsAdditive(t *testing.T) {
	ruleDescriptor := (&localv1.NetworkRule{}).ProtoReflect().Descriptor()
	for name, number := range map[string]protoreflect.FieldNumber{
		"protocol": 1, "cidr": 2, "port": 3, "domain": 4, "rule_id": 5, "origin": 6,
	} {
		field := ruleDescriptor.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("NetworkRule.%s field %d is missing or renumbered", name, number)
		}
	}

	responseDescriptor := (&localv1.GetOrgPolicyResponse{}).ProtoReflect().Descriptor()
	for name, number := range map[string]protoreflect.FieldNumber{
		"org": 1, "schema_version": 2, "revision": 3, "bundle_id": 4,
		"enforcement_enabled": 5, "loaded_at_unix": 6, "issued_at_unix": 7,
		"expires_at_unix": 8, "rules": 9, "ignored_rules": 10,
	} {
		field := responseDescriptor.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("GetOrgPolicyResponse.%s field %d is missing or renumbered", name, number)
		}
	}

	ruleView := (&localv1.OrgPolicyRule{}).ProtoReflect().Descriptor()
	for name, number := range map[string]protoreflect.FieldNumber{
		"rule_id": 1, "action": 2, "source": 3, "resource_type": 4,
		"protocol": 5, "cidr": 6, "domain": 7, "port": 8, "enforced": 9,
		"advisory_reason": 10,
	} {
		field := ruleView.Fields().ByName(protoreflect.Name(name))
		if field == nil || field.Number() != number {
			t.Fatalf("OrgPolicyRule.%s field %d is missing or renumbered", name, number)
		}
	}

	want := &localv1.GetOrgPolicyResponse{
		Org: "acme.example", SchemaVersion: 1, Revision: 7, BundleId: "sha256:deadbeef",
		EnforcementEnabled: true, LoadedAtUnix: 1, IssuedAtUnix: 2, ExpiresAtUnix: 3,
		IgnoredRules: 1,
		Rules: []*localv1.OrgPolicyRule{{
			RuleId: "rule-001", Action: "deny", Source: "org", ResourceType: "network",
			Protocol: "tcp", Cidr: "203.0.113.7/32", Port: 443, Enforced: true,
		}},
	}
	encoded, err := proto.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got localv1.GetOrgPolicyResponse
	if err := proto.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	if !proto.Equal(want, &got) {
		t.Fatalf("org policy view round-trip mismatch: got %v, want %v", &got, want)
	}
}
