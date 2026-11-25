package version

import (
	"testing"

	localv1 "github.com/gyrel/sandboxes-api/gen/go/sandboxes/local/v1"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestSkillsContractFieldNumbersAndRoundTrip(t *testing.T) {
	for _, field := range []struct {
		message protoreflect.ProtoMessage
		number  protoreflect.FieldNumber
	}{{&localv1.CreateRequest{}, 18}, {&localv1.Sandbox{}, 23}} {
		descriptor := field.message.ProtoReflect().Descriptor().Fields().ByName("skills")
		if descriptor == nil || descriptor.Number() != field.number || descriptor.Kind() != protoreflect.MessageKind {
			t.Fatalf("invalid skills field: %v", descriptor)
		}
	}
	for _, mode := range []localv1.SkillsMode{localv1.SkillsMode_SKILLS_MODE_OFF, localv1.SkillsMode_SKILLS_MODE_READONLY, localv1.SkillsMode_SKILLS_MODE_READWRITE} {
		original := &localv1.CreateRequest{Skills: &localv1.SkillsOptions{Mode: mode, StorePath: "/synthetic/store"}}
		wire, err := proto.Marshal(original)
		if err != nil {
			t.Fatal(err)
		}
		var result localv1.CreateRequest
		if err := proto.Unmarshal(wire, &result); err != nil {
			t.Fatal(err)
		}
		if !proto.Equal(original, &result) {
			t.Fatalf("round trip: %v", &result)
		}
	}
}
