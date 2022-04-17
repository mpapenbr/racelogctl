package util

import (
	"fmt"
	"racelogctl/internal"
	"reflect"
	"testing"
)

func TestDuplicateArray(t *testing.T) {
	type args struct {
		src []interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		// TODO: Add test cases.
		{
			name: "standard",
			args: args{src: []interface{}{1, 2}},
			want: []interface{}{1, 2},
		},
		{
			name: "nil input",
			args: args{src: nil},
			want: []interface{}{},
		},
		{
			name: "empty input",
			args: args{src: []interface{}{}},
			want: []interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DuplicateArray(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DuplicateArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatchSession(t *testing.T) {
	type args struct {
		src       []interface{}
		patchData []interface{}
	}

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "standard",
			args: args{src: []interface{}{1, 2, 3}, patchData: []interface{}{[]interface{}{0, 10}, []interface{}{1, 20}}},
			want: []interface{}{10, 20, 3},
		},
		{
			name: "new data for empty list",
			args: args{src: []interface{}{}, patchData: []interface{}{[]interface{}{0, 10}, []interface{}{2, 20}}},
			want: []interface{}{10, nil, 20},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PatchSession(tt.args.src, tt.args.patchData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PatchArray() = %v, want %v", got, tt.want)
			}
		})
	}
	fmt.Printf("tests: %v\n", tests)
}

func TestProcessDeltaStates(t *testing.T) {
	// data := []byte(`
	// [{"type": 1, "payload": {"cars": [["RUN",15,1.2], ["RUN",16,1.4]]}}],
	// [{"type": 8, "payload": {"cars": [[0,0,"OUT"], [0,2, 99.99], [1,1,27]]}}],
	// `)
	type1 := internal.State{
		Type: 1,
		Payload: internal.Payload{
			Cars:     [][]interface{}{{"RUN", 15, 1.2}, {"RUN", 16, 1.4}},
			Session:  []interface{}{"A", "B"},
			Messages: [][]interface{}{{"M1", "M2"}},
		},
	}
	type1Empty := internal.State{
		Type: 1,
		Payload: internal.Payload{
			Cars:     [][]interface{}{},
			Session:  []interface{}{},
			Messages: [][]interface{}{{"M1", "M2"}},
		},
	}
	type8 := internal.State{
		Type: 8,
		Payload: internal.Payload{
			Cars:     [][]interface{}{{0, 0, "OUT"}, {0, 1, 20}, {1, 2, 99.99}},
			Session:  []interface{}{[]interface{}{0, "C"}, []interface{}{1, "D"}},
			Messages: [][]interface{}{{"M3"}},
		},
	}
	result := internal.State{
		Type: 1,
		Payload: internal.Payload{
			Cars:     [][]interface{}{{"OUT", 20, 1.2}, {"RUN", 16, 99.99}},
			Session:  []interface{}{"C", "D"},
			Messages: [][]interface{}{{"M3"}},
		},
	}
	result2 := internal.State{
		Type: 1,
		Payload: internal.Payload{
			Cars:     [][]interface{}{{"OUT", 20, nil}, {nil, nil, 99.99}},
			Session:  []interface{}{"C", "D"},
			Messages: [][]interface{}{{"M3"}},
		},
	}

	type args struct {
		current  internal.State
		incoming internal.State
	}

	tests := []struct {
		name string
		args args
		want internal.State
	}{
		{
			name: "both empty",
			args: args{internal.State{}, internal.State{}},
			want: internal.State{},
		},
		{
			name: "incoming type 1",
			args: args{internal.State{}, type1},
			want: type1,
		},
		{
			name: "incoming type 8",
			args: args{type1, type8},
			want: result,
		},
		{
			name: "incoming type 8 on empty type 1",
			args: args{type1Empty, type8},
			want: result2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessDeltaStates(tt.args.current, tt.args.incoming); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessDeltaStates() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestConvertJsonToGo(t *testing.T) {
	type args struct {
		jsonData []byte
	}
	tests := []struct {
		name string
		args args
		want internal.State
	}{
		{
			// note: when converting json to go struct, the "int" values become float64
			name: "float64 check",
			args: args{[]byte(`{"type":1, "payload": {"cars": [["RUN",15,1.2]]}}`)},
			want: internal.State{
				Type:      1,
				Payload:   internal.Payload{Cars: [][]interface{}{{"RUN", float64(15), 1.2}}},
				Timestamp: 0,
			},
		},
		{
			name: "standard",
			args: args{[]byte(`{"type":1, "payload": {"cars": [["RUN",15,1.2], ["RUN",16,1.4]]}}`)},
			want: internal.State{
				Type:      1,
				Payload:   internal.Payload{Cars: [][]interface{}{{"RUN", float64(15), 1.2}, {"RUN", float64(16), 1.4}}},
				Timestamp: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertJsonToGo(tt.args.jsonData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertJsonToGo() = %v, want %v", got, tt.want)
			}
		})
	}
}
