package main

import "testing"

func Test_validate200Record(t *testing.T) {
	type args struct {
		record []string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy path",
			args: args{
				record: []string{"200", "NEM1201009", "E1E2", "1", "E1", "N1", "01009", "kWh", "30", "20050610"},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "not enough elements",
			args: args{
				record: []string{"1", "3"},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate200Record(tt.args.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate200Record() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validate200Record() = %v, want %v", got, tt.want)
			}
		})
	}
}
