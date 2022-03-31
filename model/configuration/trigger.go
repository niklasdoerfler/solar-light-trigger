package configuration

type TriggerConfiguration struct {
	Enabled    bool
	Name       string
	Threshold  float64
	Hysteresis float64
}
