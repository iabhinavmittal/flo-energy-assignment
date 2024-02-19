package main

import "testing"

func Test_validate200Record(t *testing.T) {
	type args struct {
		record []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				record: []string{"200", "NEM1201009", "E1E2", "1", "E1", "N1", "01009", "kWh", "30", "20050610"},
			},
			wantErr: false,
		},
		{
			name: "not enough elements",
			args: args{
				record: []string{"1", "3"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate200Record(tt.args.record); (err != nil) != tt.wantErr {
				t.Errorf("validate200Record() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO:
// - empty file
// - invalid file, assuming only csv files for now
// - missing NEM value
// - missing interval value
// - wrong interval value?
// - missing date value
// - wrong date format
// - insufficient number of consumption values according to interval?
// - large number of records: x ?
// - no 200 record found
// - no 300 record found
// - wrong order, 300 records are above 200
