package pdiutil

import (
	"reflect"
	"testing"
	"time"
)

func TestParseXrepDateString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			"basic test",
			args{"20181231092019.5268080 "},
			time.Date(int(2018), time.December, int(31), int(9), int(20), int(19), int(0), time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseXrepDateString(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseXrepDateString() = %v, want %v", got, tt.want)
			}
		})
	}
}
