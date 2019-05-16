package brew

import (
	"testing"
)

func TestTempDummy_Get(t *testing.T) {

	type fields struct {
		fn   callback
		temp float64
	}
	type test struct {
	}
	tests := []struct {
		name    string
		fields  fields
		want    float64
		wantErr bool
	}{
		{"Callback increase", fields{func(x float64) float64 { return x + 4.0 }, 22.0}, 26.0, false},
		{"Callback decrease", fields{func(x float64) float64 { return x - 4.0 }, 22.0}, 18.0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := &TempDummy{
				Fn:   tt.fields.fn,
				Temp: tt.fields.temp,
			}
			got, err := td.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("TempDummy.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TempDummy.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
