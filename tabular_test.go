package jp

import (
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
	dge "github.com/wiryax/direct-graph-engine"
)

func TestJsonToTabular(t *testing.T) {
	tests := []struct {
		title,
		payload string
		wantErr  bool
		expected func() dge.Tabular
	}{
		{
			title:   "single object",
			payload: `{"a":"a","b":"b","c":"c"}`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular([]string{"a", "b", "c"})
				e.AddRow(dge.ParseVariable([]byte("a")), dge.ParseVariable([]byte("b")), dge.ParseVariable([]byte("c")))
				return *e
			},
		}, {
			title:   "invalid JSON",
			payload: `{"a"}`,
			wantErr: true,
			expected: func() dge.Tabular {
				return dge.Tabular{}
			},
		}, {
			title:   "object array",
			payload: `{"a": ["1","2"]}`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular([]string{"a"})
				e.AddRow(dge.ParseVariable([]byte("1")))
				e.AddRow(dge.ParseVariable([]byte("2")))
				return *e
			},
		}, {
			title:   "object multi dimensional array with object",
			payload: `{"a": [{"b": "1"},{"b": "2"}]}`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular([]string{"ab"})
				e.AddRow(dge.ParseVariable([]byte("1")))
				e.AddRow(dge.ParseVariable([]byte("2")))
				return *e
			},
		}, {
			title:   "nested object",
			payload: `{"a": {"a":"1", "b": "2"}, "b" : "3"}`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular([]string{"b", "aa", "ab"})
				e.AddRow(dge.ParseVariable([]byte("3")), dge.ParseVariable([]byte("1")), dge.ParseVariable([]byte("2")))
				return *e
			},
		}, {
			title:   "inconsistency structure",
			payload: `[{"a":"1","b":"2"},"3", 4]`,
			wantErr: false,
			expected: func() dge.Tabular {
				e := dge.MakeTabular([]string{"", "a", "b"})
				e.AddRow(dge.ParseVariable([]byte("3")), dge.ParseVariable([]byte("1")), dge.ParseVariable([]byte("2")))
				e.AddRow(dge.ParseVariable([]byte("4")), dge.ParseVariable([]byte("1")), dge.ParseVariable([]byte("2")))
				return *e
			},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			e := test.expected()
			result, err := parseJsonToTabular([]byte(test.payload))
			if test.wantErr != (err != nil) {
				t.Fatalf("unexpected err result. want %v, got %v. err: %v", test.wantErr, (err != nil), err)
			}

			if !reflect.DeepEqual(e, result) {
				t.Errorf("unexpected result. want %v, got %v", e.String(), result.String())
			}
		})
	}
}

func TestJsonToTabular_NestedObject(t *testing.T) {
	b := []byte(`{"a":"a","b":{"c":"c","d":"d"}}`)

	expected := *dge.MakeTabular([]string{"a", "bc", "bd"})

	expected.AddRow(dge.ParseVariable([]byte("a")), dge.ParseVariable([]byte("c")), dge.ParseVariable([]byte("d")))

	result, err := parseObject(gjson.Parse(string(b)), "")
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("unexpected result. want %v, got %v", expected, result)
	}
}

func TestJsonToTabular_Array(t *testing.T) {
	b := []byte(`["a1", "a2"]`)
	expected := *dge.MakeTabular([]string{"a"})

	expected.AddRow(dge.ParseVariable([]byte("a1")))
	expected.AddRow(dge.ParseVariable([]byte("a2")))

	result, err := parseArray(gjson.Parse(string(b)), "a")
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("unexpected result. want %v got %v", expected.String(), result.String())
	}
}
func TestJsonToTabular_ArrayWithObject(t *testing.T) {
	b := []byte(`[{"b": "b1", "c": "c1"}, "a1"]`)
	expected := *dge.MakeTabular([]string{"a", "ab", "ac"})

	expected.AddRow(dge.ParseVariable([]byte("a1")), dge.ParseVariable([]byte("b1")), dge.ParseVariable([]byte("c1")))

	result, err := parseArray(gjson.Parse(string(b)), "a")
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("unexpected result. want %v got %v", expected.String(), result.String())
	}
}

func TestJsonToTabular_MultiDimensionalArray(t *testing.T) {
	b := []byte(`["a3",["a1", "a2"]]`)
	expected := *dge.MakeTabular([]string{"a"})

	expected.AddRow(dge.ParseVariable([]byte("a3")))
	expected.AddRow(dge.ParseVariable([]byte("a1")))
	expected.AddRow(dge.ParseVariable([]byte("a2")))

	result, err := parseArray(gjson.Parse(string(b)), "a")
	if err != nil {
		t.Fatalf("unexpected err %v", err)
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("unexpected result. want %v got %v", expected.String(), result.String())
	}
}
