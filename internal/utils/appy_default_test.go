package utils

import (
	"errors"
	"reflect"
	"testing"
)

// Mock struct
type TestStruct struct {
	BoolField   *bool   `default:"true"`
	IntField    *int    `default:"42"`
	Int8Field   *int8   `default:"127"`
	Uint32Field *uint32 `default:"1000"`
	StringField *string `default:"hello"`
	NoTagField  *int
	NonNilField *int
	NonPtrField int `default:"99"`
}

func TestApplyDefaults(t *testing.T) {
	// List of func, return a pointer to its value
	newInt := func(i int) *int { return &i }
	newString := func(s string) *string { return &s }
	newBool := func(b bool) *bool { return &b }
	newInt8 := func(i int8) *int8 { return &i }
	newUint32 := func(u uint32) *uint32 { return &u }

	tests := []struct {
		name    string
		input   any
		want    any
		wantErr error
	}{
		{
			name:    "Non-pointer input",
			input:   TestStruct{},
			wantErr: errors.New("obj should be a pointer and not a nil pointer"),
		},
		{
			name:    "Nil pointer input",
			input:   (*TestStruct)(nil),
			wantErr: errors.New("obj should be a pointer and not a nil pointer"),
		},
		{
			name: "Valid input with defaults",
			input: &TestStruct{
				NonNilField: newInt(50),
				NonPtrField: 99,
			},
			want: &TestStruct{
				BoolField:   newBool(true),
				IntField:    newInt(42),
				Int8Field:   newInt8(127),
				Uint32Field: newUint32(1000),
				StringField: newString("hello"),
				NoTagField:  nil,
				NonNilField: newInt(50),
				NonPtrField: 99,
			},
			wantErr: nil,
		},
		{
			name: "Invalid default tag for int",
			input: &struct {
				IntField *int `default:"invalid"`
			}{},
			want: &struct {
				IntField *int `default:"invalid"`
			}{IntField: nil},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApplyDefaults(tt.input)
			if !errors.Is(err, tt.wantErr) &&
				(err == nil || tt.wantErr == nil || (err.Error() != tt.wantErr.Error())) {
				t.Errorf("ApplyDefaults() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if !reflect.DeepEqual(tt.input, tt.want) {
					t.Errorf("ApplyDefaults() result = %+v, want %+v", tt.input, tt.want)
				}
			}
		})
	}
}
