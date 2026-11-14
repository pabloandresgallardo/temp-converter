package temperature_test

import (
	"testing"

	"github.com/gallardo/temp-converter/internal/domain/temperature"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	cases := []struct {
		name  string
		value float64
		from  temperature.Unit
		to    temperature.Unit
		want  float64
	}{
		{"0C→F", 0, temperature.Celsius, temperature.Fahrenheit, 32},
		{"100C→F", 100, temperature.Celsius, temperature.Fahrenheit, 212},
		{"212F→C", 212, temperature.Fahrenheit, temperature.Celsius, 100},
		{"0C→K", 0, temperature.Celsius, temperature.Kelvin, 273.15},
		{"273.15K→C", 273.15, temperature.Kelvin, temperature.Celsius, 0},
		{"32F→K", 32, temperature.Fahrenheit, temperature.Kelvin, 273.15},
		{"273.15K→F", 273.15, temperature.Kelvin, temperature.Fahrenheit, 32},
		{"SameUnit", 100, temperature.Celsius, temperature.Celsius, 100},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := temperature.Convert(tc.value, tc.from, tc.to)
			require.NoError(t, err)
			assert.InDelta(t, tc.want, got, 0.001)
		})
	}
}

func TestConvert_InvalidUnit(t *testing.T) {
	_, err := temperature.Convert(100, "rankine", temperature.Celsius)
	assert.ErrorContains(t, err, "invalid unit")
}
