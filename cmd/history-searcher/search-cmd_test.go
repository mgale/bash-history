package main

import (
	"reflect"
	"testing"
)

func Test_convertArgsToQuery(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want []rune
	}{
		{"emptyArgs", args{[]string{}}, []rune{'*'}},
		{"simleArgs", args{[]string{"ls", " ", "-l"}}, []rune{'l', 's', ' ', '-', 'l'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertArgsToQuery(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertArgsToQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
