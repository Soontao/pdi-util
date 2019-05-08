package pdiutil

import (
	"io/ioutil"
	"testing"
)

func TestCountElementForBODL(t *testing.T) {
	bytes, _ := ioutil.ReadFile("./ast/test_data/sample_bo.bo")

	type args struct {
		source []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Count Sample BO Elements",
			args{bytes},
			8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountElementForBODL(tt.args.source); got != tt.want {
				t.Errorf("CountElementForBODL() = %v, want %v", got, tt.want)
			}
		})
	}
}
