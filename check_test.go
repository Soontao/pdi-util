package main

import "testing"

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
