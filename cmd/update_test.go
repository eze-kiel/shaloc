package cmd

import (
	"reflect"
	"testing"
)

func Test_getVersionsList(t *testing.T) {
	tests := []struct {
		name string
		r    releases
		want []string
	}{
		{
			name: "typical case",
			r: releases{
				{TagName: "v1.2.3"},
				{TagName: "v1.2.2"},
				{TagName: "v1.2.1"},
			},
			want: []string{"v1.2.3", "v1.2.2", "v1.2.1"},
		},
		{
			name: "empty release",
			r:    releases{{}},
			want: []string{""},
		},
		{
			name: "one release",
			r:    releases{{TagName: "v1.2.3"}},
			want: []string{"v1.2.3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getVersionsList(tt.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getVersionsList() = %v, want %v", got, tt.want)
			}
		})
	}
}
