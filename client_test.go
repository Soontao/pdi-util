package pdiutil

import "testing"

func TestGetReleaseVersionForTenant(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"byd test",
			args{"my600101.sapbyd.cn"},
			false,
		},
		{
			"c4c test",
			args{"my312006.crm.ondemand.com"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRt, err := GetReleaseVersionForTenant(tt.args.host)
			t.Logf("Get release: %v for %v", gotRt, tt.args.host)
			if gotRt == "" {
				t.Errorf("GetReleaseVersionForTenant() gotRt is empty")
				return
			}
			if len(gotRt) != 4 {
				t.Errorf("GetReleaseVersionForTenant() gotRt = %v, is not valid", gotRt)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReleaseVersionForTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
