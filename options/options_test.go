// package options_test
package options

import (
	"os"
	"testing"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
)

var filemap = map[string]string{
	// test proto
	"message.proto": `
syntax = "proto3";
package test;

import "ext.proto";

message Thing {
	option (test.msg_string) = "dogs";
	option (test.msg_bool) = true;

	string id = 1 [
		(test.field_string) = "cats",
		(test.field_bool) = true
	];
}
	`,

	// descriptor extensions
	"ext.proto": `
syntax = "proto3";
package test;

import "google/protobuf/descriptor.proto";

extend google.protobuf.MessageOptions {
	string msg_string = 5001;
	bool msg_bool = 5002;
}

extend google.protobuf.FieldOptions {
	string field_string = 5001;
	bool field_bool = 5002;
}
	`,
}

var p = protoparse.Parser{
	Accessor: protoparse.FileContentsFromMap(filemap),
}

var reader = NewOptionReader(nil)

type tablecase struct {
	code int
	name string
	node desc.Descriptor

	srcByID   interface{}
	srcByName interface{}

	expected interface{}
}

var tableCases map[string]tablecase

func TestMain(m *testing.M) {
	files, err := p.ParseFiles("message.proto")
	if err != nil {
		panic(err)
	}

	file := files[0]
	mdesc := file.FindSymbol("test.Thing").(*desc.MessageDescriptor)
	fdesc := file.FindSymbol("test.Thing.id").(*desc.FieldDescriptor)

	tableCases = map[string]tablecase{
		"test.Thing(test.msg_string)": {
			code:      5001,
			name:      "test.msg_string",
			node:      mdesc,
			srcByID:   new(string),
			srcByName: new(string),
			expected:  "dogs",
		},
		"test.Thing(test.msg_bool)": {
			code:      5002,
			name:      "test.msg_bool",
			node:      mdesc,
			srcByID:   new(bool),
			srcByName: new(bool),
			expected:  true,
		},
		"test.Thing.id(test.field_string)": {
			code:      5001,
			name:      "test.field_string",
			node:      fdesc,
			srcByID:   new(string),
			srcByName: new(string),
			expected:  "cats",
		},
		"test.Thing.id(test.field_bool)": {
			code:      5002,
			name:      "test.field_bool",
			node:      fdesc,
			srcByID:   new(bool),
			srcByName: new(bool),
			expected:  true,
		},
	}

	os.Exit(m.Run())
}

func TestAll(t *testing.T) {
	if tableCases == nil {
		t.Error("tableCases not initialized")
		t.FailNow()
		return
	}

	for name, args := range tableCases {
		name, args := name, args

		// GetByID
		t.Run(name+"....GetByID()", func(t *testing.T) {
			actual, err := reader.GetOptionByID(args.node, args.code)
			if err != nil {
				t.Error(err)
				return
			}

			if actual != args.expected {
				t.Errorf("actual: %v -- expected: %v", actual, args.expected)
			} else {
				t.Logf("%q == %v, which is expected", name, args.expected)
			}
		})

		// GetByName
		t.Run(name+"....GetByName()", func(t *testing.T) {
			actual, err := reader.GetOptionByName(args.node, args.name)
			if err != nil {
				t.Error(err)
				return
			}

			if actual != args.expected {
				t.Errorf("actual: %v -- expected: %v", actual, args.expected)
			} else {
				t.Logf("%q == %v, which is expected", name, args.expected)
			}
		})

		// ReadByID
		t.Run(name+"...ReadByID()", func(t *testing.T) {
			if err := reader.ReadOptionByID(args.node, args.code, args.srcByID); err != nil {
				t.Error(err)
				return
			}

			switch val := args.srcByID.(type) {
			case *string:
				if *val != args.expected.(string) {
					t.Errorf("actual: %v -- expected: %v", val, args.expected)
				} else {
					t.Logf("%q == %v, which is expected", name, *val)
				}
			case *bool:
				if *val != args.expected.(bool) {
					t.Errorf("actual: %v -- expected: %v", val, args.expected)
				} else {
					t.Logf("%q == %v, which is expected", name, *val)
				}
			}
		})

		// ReadByName
		t.Run(name+"...ReadByName()", func(t *testing.T) {
			if err := reader.ReadOptionByName(args.node, args.name, args.srcByName); err != nil {
				t.Error(err)
				return
			}

			switch val := args.srcByName.(type) {
			case *string:
				if *val != args.expected.(string) {
					t.Errorf("actual: %v -- expected: %v", val, args.expected)
				} else {
					t.Logf("%q == %v, which is expected", name, *val)
				}
			case *bool:
				if *val != args.expected.(bool) {
					t.Errorf("actual: %v -- expected: %v", val, args.expected)
				} else {
					t.Logf("%q == %v, which is expected", name, *val)
				}
			}
		})
	}
}
