package main

import (
	"testing"
)

func Test_ensureFileNameConvention(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
	}{
		{
			"test cldt",
			args{"CLDT_CSomeType.codelist"},
			true,
			"CLDT_CSomeType.codelist",
		},
		{
			"test cldt2",
			args{"CL_CSomeType.codelist"},
			false,
			"CLDT_CSomeType.codelist",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ensureFileNameConvention(tt.args.filePath)
			if got != tt.want {
				t.Errorf("ensureFileNameConvention() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ensureFileNameConvention() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_shortenPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"simple test",
			args{"Folder1\\Hello\\Test.XX"},
			"F\\H\\Test.XX",
		},
		{
			"empty test",
			args{"Test.XX"},
			"Test.XX",
		},
		{
			"relative test",
			args{"..\\Empty\\Test\\Test.XX"},
			"..\\E\\T\\Test.XX",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shortenPath(tt.args.path); got != tt.want {
				t.Errorf("shortenPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
