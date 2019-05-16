package brew

import (
	"math"
	"testing"
	"time"
)

func TestKettle_GoToTemp(t *testing.T) {
	type fields struct {
		Temp     TempSensor
		Heater   Control
		Agitator Control
	}
	type args struct {
		tempTo float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Dummy", fields{
			&TempDummy{"TempDummy", func(x float64) float64 { return x + 4.0 }, 22.0},
			&SSRDummy{},
			&SSRDummy{},
		}, args{30.0}, false},
		{"Dummy", fields{
			&TempDummy{"TempDummy", func(x float64) float64 { return x - 4.0 }, 22.0},
			&SSRDummy{},
			&SSRDummy{},
		}, args{10.0}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kettle{
				Temp:     tt.fields.Temp,
				Heater:   tt.fields.Heater,
				Agitator: tt.fields.Agitator,
			}
			if err := k.GoToTemp(tt.args.tempTo); (err != nil) != tt.wantErr {
				t.Errorf("Kettle.GoToTemp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKettle_HoldTemp(t *testing.T) {
	type fields struct {
		Temp     TempSensor
		Heater   Control
		Agitator Control
	}
	type args struct {
		temp     float64
		holdTime time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Dummy", fields{
			&TempDummy{"TempDummy", func(x float64) float64 {
				if math.Mod(x, 2.0) > 0 {
					return x - 1.0
				}
				return x + 1.0
			}, 30.0},
			&SSRDummy{},
			&SSRDummy{},
		}, args{30.0, time.Second * 4}, false},
	}
	//add timeout
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kettle{
				Temp:     tt.fields.Temp,
				Heater:   tt.fields.Heater,
				Agitator: tt.fields.Agitator,
			}
			if err := k.HoldTemp(tt.args.temp, tt.args.holdTime); (err != nil) != tt.wantErr {
				t.Errorf("Kettle.HoldTemp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
