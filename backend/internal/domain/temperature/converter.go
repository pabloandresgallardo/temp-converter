package temperature

import "fmt"

// Unit representa una escala de temperatura.
type Unit string

const (
	Celsius    Unit = "celsius"
	Fahrenheit Unit = "fahrenheit"
	Kelvin     Unit = "kelvin"
)

// Convert convierte value de la unidad from a la unidad to.
// Retorna error si alguna unidad es inválida.
func Convert(value float64, from, to Unit) (float64, error) {
	toCelsius := map[Unit]func(float64) float64{
		Celsius:    func(v float64) float64 { return v },
		Fahrenheit: func(v float64) float64 { return (v - 32) * 5 / 9 },
		Kelvin:     func(v float64) float64 { return v - 273.15 },
	}
	fromCelsius := map[Unit]func(float64) float64{
		Celsius:    func(v float64) float64 { return v },
		Fahrenheit: func(v float64) float64 { return v*9/5 + 32 },
		Kelvin:     func(v float64) float64 { return v + 273.15 },
	}

	toC, ok := toCelsius[from]
	if !ok {
		return 0, fmt.Errorf("invalid unit: %q", from)
	}
	fromC, ok := fromCelsius[to]
	if !ok {
		return 0, fmt.Errorf("invalid unit: %q", to)
	}
	return fromC(toC(value)), nil
}
