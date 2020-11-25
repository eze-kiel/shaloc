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

func Test_stringInSlice(t *testing.T) {
	tests := []struct {
		name string
		s    string
		sl   []string
		want bool
	}{
		{name: "case1", s: "plop", sl: []string{"plip", "plop", "ploup"}, want: true},
		{name: "case2", s: "plop", sl: []string{"plip", "plap", "ploup"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringInSlice(tt.s, tt.sl); got != tt.want {
				t.Errorf("stringInSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
