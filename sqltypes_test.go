package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestStructEmbedding(t *testing.T) {
	tim := time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := []byte(`{"a":123,"b":true,"c":123.123,"d":"string","e":"2017-01-01T00:00:00Z","f":[1,2,3]}`)
	type embed struct {
		A Int64   `json:"a,omitempty"`
		B Bool    `json:"b,omitempty"`
		C Float64 `json:"c,omitempty"`
		D String  `json:"d,omitempty"`
		E Time    `json:"e,omitempty"`
		F RawJSON `json:"f,omitempty"`
	}
	em := embed{
		A: Int64{Valid: true, Int64: 123},
		B: Bool{Valid: true, Bool: true},
		C: Float64{Valid: true, Float64: 123.123},
		D: String{Valid: true, String: "string"},
		E: Time{Valid: true, Time: tim},
		F: RawJSON(`[1,2,3]`),
	}
	b, err := json.Marshal(em)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, b) {
		t.Fatal("not the same JSON!")
	}
	if !(string(b) == string(expected)) {
		t.Fatal("not the same!")
	}

	var em2 embed
	if err = json.Unmarshal(expected, &em2); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(em2, em) {
		t.Fatal("not correct")
	}
}

func TestNullString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       String
		source  []byte
		wantErr bool
	}{
		{
			name:    "explicit null",
			source:  []byte(`null`),
			wantErr: false,
		},
		{
			name:    "string null",
			source:  []byte(`"null"`),
			wantErr: false, // this one SHOULD be valid
		},
		{
			name:    "valid",
			source:  []byte(`"hello"`),
			wantErr: false,
		},
		{
			name:    "invalid",
			source:  []byte(`{"key":"value"}`),
			wantErr: true,
		},
		{
			name:    "empty",
			source:  []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnmarshalJSON(tt.source); (err != nil) != tt.wantErr {
				t.Errorf("String.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNullString_Value(t *testing.T) {
	tests := []struct {
		name    string
		n       String
		want    driver.Value
		wantErr bool
	}{
		{
			name: "valid",
			n: String{
				Valid:  true,
				String: "hello",
			},
			want:    driver.Value("hello"),
			wantErr: false,
		},
		{
			name: "invalid",
			n: String{
				Valid: false,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullString_Scan(t *testing.T) {
	tests := []struct {
		name    string
		n       *String
		wantErr bool
		src     any
	}{
		{
			name: "valid",
			n: &String{
				String: "hello",
				Valid:  true,
			},
			src:     "",
			wantErr: false,
		},
		{
			name: "nil value",
			n: &String{
				String: "hello",
				Valid:  false,
			},
			src:     nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Scan(tt.src); (err != nil) != tt.wantErr {
				t.Errorf("String.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.n.Valid && tt.src != nil {
				t.Errorf("should return null")
			}
			if tt.n.Valid && tt.src != tt.n.String {
				t.Errorf("invalid value")
			}
		})
	}
}

func TestNullString_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       *String
		want    []byte
		wantErr bool
	}{
		{
			name: "valid",
			n: &String{
				String: "hello",
				Valid:  true,
			},
			want:    []byte(`"hello"`),
			wantErr: false,
		},
		{
			name: "valid null",
			n: &String{
				String: "",
				Valid:  false,
			},
			want:    []byte(`null`),
			wantErr: false,
		},
		{
			name: "invalid",
			n: &String{
				Valid: true,
			},
			want:    []byte(`""`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("String.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullBool_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name         string
		n            Bool
		source       []byte
		wantErr      bool
		wantValidity bool
	}{
		{
			name:         "explicit null",
			source:       []byte(`null`),
			wantErr:      false,
			wantValidity: false,
		},
		{
			name:         "valid",
			source:       []byte(`false`),
			wantErr:      false,
			wantValidity: true,
		},
		{
			name:         "invalid",
			source:       []byte(`{"key":"value"}`),
			wantErr:      true,
			wantValidity: false,
		},
		{
			name:         "empty",
			source:       []byte{},
			wantErr:      true,
			wantValidity: false,
		},
		{
			name:         "explicit null",
			source:       []byte("null"),
			wantErr:      false,
			wantValidity: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnmarshalJSON(tt.source); (err != nil) != tt.wantErr && tt.n.Valid == tt.wantValidity {
				t.Errorf("Bool.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNullBool_Value(t *testing.T) {
	tests := []struct {
		name    string
		n       Bool
		want    driver.Value
		wantErr bool
	}{
		{
			name: "valid",
			n: Bool{
				Valid: true,
				Bool:  true,
			},
			want:    driver.Value(true),
			wantErr: false,
		},
		{
			name: "invalid",
			n: Bool{
				Valid: false,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullBool_Scan(t *testing.T) {
	tests := []struct {
		name    string
		n       *Bool
		wantErr bool
		src     any
	}{
		{
			name: "valid",
			n: &Bool{
				Bool:  true,
				Valid: true,
			},
			src:     true,
			wantErr: false,
		},
		{
			name: "nil value",
			n: &Bool{
				Bool:  true,
				Valid: false,
			},
			src:     false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Scan(tt.src); (err != nil) != tt.wantErr {
				t.Errorf("Bool.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.n.Valid && tt.src != nil {
				t.Errorf("should return null")
			}
			if tt.n.Valid && tt.src != tt.n.Bool {
				t.Errorf("invalid value")
			}
		})
	}
}

func TestNullBool_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       *Bool
		want    []byte
		wantErr bool
	}{
		{
			name: "valid",
			n: &Bool{
				Valid: true,
			},
			want:    []byte(`false`),
			wantErr: false,
		},
		{
			name: "valid null",
			n: &Bool{
				Valid: false,
			},
			want:    []byte(`null`),
			wantErr: false,
		},
		{
			name: "invalid",
			n: &Bool{
				Valid: true,
			},
			want:    []byte(`false`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       Time
		source  []byte
		wantErr bool
	}{
		{
			name:    "valid",
			source:  []byte(`"2017-11-24T00:00:00Z"`),
			wantErr: false,
		},
		{
			name:    "invalid",
			source:  []byte(`{"key":"value"}`),
			wantErr: true,
		},
		{
			name:    "empty",
			source:  []byte{},
			wantErr: true,
		},
		{
			name:    "explicit null",
			source:  []byte(`null`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnmarshalJSON(tt.source); (err != nil) != tt.wantErr {
				t.Errorf("Bool.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTime_Value(t *testing.T) {
	tim := time.Now()
	tests := []struct {
		name    string
		n       Time
		want    driver.Value
		wantErr bool
	}{
		{
			name: "valid",
			n: Time{
				Valid: true,
				Time:  tim,
			},
			want:    driver.Value(tim),
			wantErr: false,
		},
		{
			name: "invalid",
			n: Time{
				Valid: false,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Time.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime_Scan(t *testing.T) {
	tim := time.Now()
	tests := []struct {
		name    string
		n       *Time
		wantErr bool
		src     any
	}{
		{
			name: "valid",
			n: &Time{
				Time:  tim,
				Valid: true,
			},
			src:     tim,
			wantErr: false,
		},
		{
			name: "nil value",
			n: &Time{
				Time:  tim,
				Valid: false,
			},
			src:     time.Now(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Scan(tt.src); (err != nil) != tt.wantErr {
				t.Errorf("Time.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.n.Valid && tt.src != nil {
				t.Errorf("should return null")
			}
			if tt.n.Valid && tt.src != tt.n.Time {
				t.Errorf("invalid value")
			}
		})
	}
}

func TestTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       *Time
		want    []byte
		wantErr bool
	}{
		{
			name: "valid",
			n: &Time{
				Time:  time.Date(2017, 11, 24, 0, 0, 0, 0, time.UTC),
				Valid: true,
			},
			want:    []byte(`"2017-11-24T00:00:00Z"`),
			wantErr: false,
		},
		{
			name: "valid null",
			n: &Time{
				Valid: false,
			},
			want:    []byte(`null`),
			wantErr: false,
		},
		{
			name: "invalid",
			n: &Time{
				Valid: true,
			},
			want:    []byte(`"0001-01-01T00:00:00Z"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Time.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullInt64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       Int64
		source  []byte
		wantErr bool
	}{
		{
			name:    "valid",
			source:  []byte(`123`),
			wantErr: false,
		},
		{
			name:    "invalid",
			source:  []byte(`{"key":"value"}`),
			wantErr: true,
		},
		{
			name:    "empty",
			source:  []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnmarshalJSON(tt.source); (err != nil) != tt.wantErr {
				t.Errorf("Int64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNullInt64_Value(t *testing.T) {
	tests := []struct {
		name    string
		n       Int64
		want    driver.Value
		wantErr bool
	}{
		{
			name: "valid",
			n: Int64{
				Valid: true,
				Int64: 123,
			},
			want:    driver.Value(int64(123)),
			wantErr: false,
		},
		{
			name: "invalid",
			n: Int64{
				Valid: false,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullInt64_Scan(t *testing.T) {
	tests := []struct {
		name    string
		n       *Int64
		wantErr bool
		src     any
	}{
		{
			name: "valid",
			n: &Int64{
				Int64: 123,
				Valid: true,
			},
			src:     int64(123),
			wantErr: false,
		},
		{
			name: "nil value",
			n: &Int64{
				Valid: false,
			},
			src:     int64(123),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Scan(tt.src); (err != nil) != tt.wantErr {
				t.Errorf("Int64.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.n.Valid && tt.src != nil {
				t.Errorf("should return null")
			}
			if tt.n.Valid && tt.src != tt.n.Int64 {
				t.Errorf("invalid value")
			}
		})
	}
}

func TestNullInt64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       *Int64
		want    []byte
		wantErr bool
	}{
		{
			name: "valid",
			n: &Int64{
				Int64: 123,
				Valid: true,
			},
			want:    []byte(`123`),
			wantErr: false,
		},
		{
			name: "valid null",
			n: &Int64{
				Valid: false,
			},
			want:    []byte(`null`),
			wantErr: false,
		},
		{
			name: "invalid",
			n: &Int64{
				Valid: true,
			},
			want:    []byte(`0`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullFloat64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       Float64
		source  []byte
		wantErr bool
	}{
		{
			name:    "explicit null",
			source:  []byte(`null`),
			wantErr: false,
		},
		{
			name:    "valid",
			source:  []byte(`123.123`),
			wantErr: false,
		},
		{
			name:    "invalid",
			source:  []byte(`{"key":"value"}`),
			wantErr: true,
		},
		{
			name:    "empty",
			source:  []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UnmarshalJSON(tt.source); (err != nil) != tt.wantErr {
				t.Errorf("Float64.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNullFloat64_Value(t *testing.T) {
	tests := []struct {
		name    string
		n       Float64
		want    driver.Value
		wantErr bool
	}{
		{
			name: "valid",
			n: Float64{
				Valid:   true,
				Float64: 123.123,
			},
			want:    driver.Value(123.123),
			wantErr: false,
		},
		{
			name: "invalid",
			n: Float64{
				Valid: false,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float64.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNullFloat64_Scan(t *testing.T) {
	tests := []struct {
		name    string
		n       *Float64
		wantErr bool
		src     any
	}{
		{
			name: "valid",
			n: &Float64{
				Float64: 123.123,
				Valid:   true,
			},
			src:     float64(123),
			wantErr: false,
		},
		{
			name: "nil value",
			n: &Float64{
				Valid: false,
			},
			src:     123.123,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Scan(tt.src); (err != nil) != tt.wantErr {
				t.Errorf("Float64.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.n.Valid && tt.src != nil {
				t.Errorf("should return null")
			}
			if tt.n.Valid && tt.src != tt.n.Float64 {
				t.Errorf("invalid value")
			}
		})
	}
}

func TestNullFloat64_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		n       *Float64
		want    []byte
		wantErr bool
	}{
		{
			name: "valid",
			n: &Float64{
				Float64: 123.123,
				Valid:   true,
			},
			want:    []byte(`123.123`),
			wantErr: false,
		},
		{
			name: "valid null",
			n: &Float64{
				Valid: false,
			},
			want:    []byte(`null`),
			wantErr: false,
		},
		{
			name: "invalid",
			n: &Float64{
				Valid: true,
			},
			want:    []byte(`0`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Float64.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullBool(t *testing.T) {
	b := true
	bb := MakeBool(&b)
	if !bb.Valid {
		t.Errorf("expected valid, got %v", bb.Valid)
	}
	if !bb.Bool {
		t.Errorf("expected true, got %v", bb.Bool)
	}

	var b2 *bool
	bb2 := MakeBool(b2)
	if bb2.Valid {
		t.Errorf("expected not valid, got %v", bb2.Valid)
	}
	if bb2.Bool {
		t.Errorf("expected false, got %v", bb2.Bool)
	}
}

func TestToNullInt64(t *testing.T) {
	b := int64(123)
	bb := MakeInt64(&b)
	if !bb.Valid {
		t.Errorf("expected valid, got %v", bb.Valid)
	}
	if bb.Int64 != 123 {
		t.Errorf("expected 123, got %v", bb.Int64)
	}

	var b2 *int64
	bb2 := MakeInt64(b2)
	if bb2.Valid {
		t.Errorf("expected not valid, got %v", bb2.Valid)
	}
	if bb2.Int64 != 0 {
		t.Errorf("expected 0, got %v", bb2.Int64)
	}
}

func TestToNullFloat64(t *testing.T) {
	b := 123.123
	bb := MakeFloat64(&b)
	if !bb.Valid {
		t.Errorf("expected valid, got %v", bb.Valid)
	}
	if bb.Float64 != 123.123 {
		t.Errorf("expected 123.123, got %v", bb.Float64)
	}

	var b2 *float64
	bb2 := MakeFloat64(b2)
	if bb2.Valid {
		t.Errorf("expected not valid, got %v", bb2.Valid)
	}
	if bb2.Float64 != 0 {
		t.Errorf("expected 0, got %v", bb2.Float64)
	}
}

func TestToNullString(t *testing.T) {
	b := "qwe"
	bb := MakeString(&b)
	if !bb.Valid {
		t.Errorf("expected valid, got %v", bb.Valid)
	}
	if bb.String != "qwe" {
		t.Errorf("expected qwe, got %v", bb.String)
	}

	var b2 *string
	bb2 := MakeString(b2)
	if bb2.Valid {
		t.Errorf("expected not valid, got %v", bb2.Valid)
	}
	if bb2.String != "" {
		t.Errorf("expected <empty string>, got %v", bb2.String)
	}
}

func TestToTime(t *testing.T) {
	tim := time.Now()
	bb := MakeTime(tim)
	if !bb.Valid {
		t.Errorf("expected valid, got %v", bb.Valid)
	}
	if bb.Time != tim {
		t.Errorf("expected %v, got %v", tim, bb.Time)
	}

	tim = time.Time{}
	bb = MakeTime(tim)
	if bb.Valid {
		t.Errorf("expected invalid, got %v", bb.Valid)
	}
	if bb.Time != tim {
		t.Errorf("expected %v, got %v", tim, bb.Time)
	}
}

func TestRawJSON_MarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		data []byte
		exp  string
	}{
		{
			name: "empty data",
			data: []byte{},
			exp:  "null",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rj := RawJSON(c.data)
			b, err := rj.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(b) != c.exp {
				t.Fatalf("\nexp: %q\ngot: %q", c.exp, string(b))
			}
		})
	}
}

type Person struct {
	Name string
	Age  int
}

func (p Person) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Person) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, p)
	default:
		return fmt.Errorf("unexpected type: %T", src)
	}
}

func TestNullBoxing(t *testing.T) {
	str := Null[string]{
		V:     "hello",
		Valid: true,
	}
	num := Null[int]{
		V:     123,
		Valid: true,
	}

	p := Person{
		Name: "John",
		Age:  30,
	}
	comp := Null[Person]{
		V:     p,
		Valid: true,
	}

	t.Run("strings", func(t *testing.T) {
		t.Run("MarshalJSON", func(t *testing.T) {
			b, err := str.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(b) != `"hello"` {
				t.Fatalf("unexpected value: %q", string(b))
			}
		})

		t.Run("UnmarshalJSON", func(t *testing.T) {
			var s2 Null[string]
			err := s2.UnmarshalJSON([]byte(`"hello"`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if s2.V != "hello" {
				t.Fatalf("unexpected value: %q", s2.V)
			}

			if !s2.Valid {
				t.Fatalf("expected valid")
			}
		})

		t.Run("Scan", func(t *testing.T) {
			var s3 Null[string]
			err := s3.Scan("hello")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if s3.V != "hello" {
				t.Fatalf("unexpected value: %q", s3.V)
			}

			if !s3.Valid {
				t.Fatalf("expected valid")
			}
		})

		t.Run("Scan null", func(t *testing.T) {
			var s4 Null[string]
			err := s4.Scan(nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if s4.V != "" {
				t.Fatalf("unexpected value: %q", s4.V)
			}

			if s4.Valid {
				t.Fatalf("expected not valid")
			}
		})

		t.Run("Value", func(t *testing.T) {
			v, err := str.Value()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if v != "hello" {
				t.Fatalf("unexpected value: %q", v)
			}
		})

		t.Run("Value null", func(t *testing.T) {
			var s5 Null[string]
			v, err := s5.Value()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if v != nil {
				t.Fatalf("unexpected value: %v", v)
			}

			if s5.Valid {
				t.Fatalf("expected not valid")
			}
		})
	})

	t.Run("ints", func(t *testing.T) {
		t.Run("MarshalJSON", func(t *testing.T) {
			b, err := num.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(b) != "123" {
				t.Fatalf("unexpected value: %q", string(b))
			}
		})

		t.Run("UnmarshalJSON", func(t *testing.T) {
			var num2 Null[int]
			err := num2.UnmarshalJSON([]byte("123"))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if num2.V != 123 {
				t.Fatalf("unexpected value: %d", num2.V)
			}
		})

		t.Run("Scan", func(t *testing.T) {
			var num3 Null[int]
			err := num3.Scan(123)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if num3.V != 123 {
				t.Fatalf("unexpected value: %d", num3.V)
			}

			if !num3.Valid {
				t.Fatalf("expected valid")
			}
		})

		t.Run("Scan null", func(t *testing.T) {
			var num4 Null[int]
			err := num4.Scan(nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if num4.V != 0 {
				t.Fatalf("unexpected value: %d", num4.V)
			}

			if num4.Valid {
				t.Fatalf("expected not valid")
			}
		})

		t.Run("Value", func(t *testing.T) {
			v, err := num.Value()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if v != 123 {
				t.Fatalf("unexpected value: %d", v)
			}
		})

		t.Run("Value null", func(t *testing.T) {
			var num5 Null[int]

			v, err := num5.Value()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if v != nil {
				t.Fatalf("unexpected value: %v", v)
			}

			if num5.Valid {
				t.Fatalf("expected not valid")
			}
		})
	})

	t.Run("complex", func(t *testing.T) {
		t.Run("MarshalJSON", func(t *testing.T) {
			b, err := comp.MarshalJSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			exp := `{"Name":"John","Age":30}`
			if string(b) != exp {
				t.Fatalf("unexpected value: %q", string(b))
			}
		})

		t.Run("UnmarshalJSON", func(t *testing.T) {
			var p2 Null[Person]
			err := p2.UnmarshalJSON([]byte(`{"Name":"John","Age":30}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if comp.V != p2.V {
				t.Fatalf("unexpected value: %+v", comp.V)
			}
		})

		t.Run("Scan", func(t *testing.T) {
			var p3 Null[Person]
			err := p3.Scan([]byte(`{"Name":"John","Age":30}`))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if comp.V != p3.V {
				t.Fatalf("unexpected value: %+v", comp.V)
			}

			if !comp.Valid {
				t.Fatalf("expected valid")
			}
		})

		t.Run("Scan null", func(t *testing.T) {
			p := Null[Person]{}
			err := p.Scan(nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if p.V != (Person{}) {
				t.Fatalf("unexpected value: %+v", p.V)
			}

			if p.Valid {
				t.Fatalf("expected not valid")
			}
		})

		t.Run("Value", func(t *testing.T) {
			v, err := comp.Value()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			p := Person{}
			if err := json.Unmarshal(v.([]byte), &p); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if comp.V != p {
				t.Fatalf("unexpected value: %+v", comp.V)
			}
		})

		t.Run("Value null", func(t *testing.T) {
			var complex2 Null[Person]

			v, err := complex2.Value()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			raw, _ := json.Marshal(Person{})

			if bytes.Compare(v.([]byte), raw) != 0 {
				t.Fatalf("unexpected value: %+s", v)
			}

			if complex2.V != (Person{}) {
				t.Fatalf("expected not valid")
			}

			if complex2.Valid {
				t.Fatalf("expected not valid")
			}
		})
	})
}
