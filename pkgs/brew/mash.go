package brew

type Masher struct {
	Temp     TempSensor
	Heater   Control
	Agitator Control
}

func (m *Masher) Mash() error {
	return nil
}
