package main;

import (
	"math"
	"strings"
	"errors"
	"fmt"
)

type Unit struct {
	// The scale factor from the base metric unit
	// to this unit
	scaleFactor float64
	// The offset from metric to this unit
	// mostly for temperature units
	unitOffset float64
}

var Distances = map[string]Unit{
	"metre": {1.0, 0.0},
	"meter": {1.0, 0.0},
	"foot":  {0.3048, 0.0},
	"feet":  {0.3048, 0.0},
	"inch":  {0.0254, 0.0},
	"inche": {0.0254, 0.0},
	"yard":  {0.9144, 0.0},
	"mile":  {1609.344, 0.0},
}

var Weights = map[string]Unit{
	"gram": {1.0, 0.0},
	"ounce": {28.34952, 0.0},
	"pound": {453.5924, 0.0},
}

var Temps = map[string]Unit {
	"celsiu": {1.0, 0.0},
	"fahrenheit": {0.5555555555555, 32.0},
	"kelvin": {1.0, -273.15},
}

var SiPrefixes = map[string]int{
	"exa": 18,
	"peta": 15,
	"tera": 12,
	"giga": 9,
	"mega": 6,
	"kilo": 3,
	"hecto": 2,
	"deca": 1,
	"deci": -1,
	"centi": -2,
	"milli": -3,
	"micro": -6,
	"nano": -9,
	"pico": -12,
	"femto": -15,
	"atto": -18,
}

var ShortForms = map[string]string {
	"cm": "centimeter",
	"m": "meter",
	"ft": "foot",
	"yd": "yard",
	"mi": "mile",
	"km": "kilometer",
	"c": "celsiu",
	"f": "fahrenheit",
	"k": "kelvin",
	"mg": "milligram",
	"g": "gram",
	"kg": "kilogram",
	"oz": "ounce",
	"lb": "pound",
}

var Categories = []map[string]Unit {Distances, Weights, Temps}

func ConvertVal(val float64, unit1 string, unit2 string) (float64, error) {
	unit1Ex, unit2Ex := unit1, unit2
	
	if ShortForms[unit1] != "" {
		unit1Ex = ShortForms[unit1]
	}

	if ShortForms[unit2] != "" {
		unit2Ex = ShortForms[unit2]
	}
	
	base1Str, scale1 := parsePrefixes(unit1Ex)
	base2Str, scale2 := parsePrefixes(unit2Ex)

	var base1, base2 Unit;
	
	catloop:
	for _, cat := range Categories {
		isCat := false
		for unitStr := range cat {
			if base1Str == unitStr {
				base1 = cat[base1Str]
				isCat = true
			}
			if base2Str == unitStr {
				base2 = cat[base2Str]
				isCat = true
			}
		}
		if isCat {
			break catloop
		}
	}

	if(base1.scaleFactor == 0.0) {
		return 0.0, errors.New(fmt.Sprintf("invalid unit '%v'", unit1))
	}

	if(base2.scaleFactor == 0.0) {
		return 0.0, errors.New(fmt.Sprintf("invalid unit '%v'", unit2))
	}

	return (val - base1.unitOffset) * ((base1.scaleFactor * scale1) / (base2.scaleFactor * scale2)) + base2.unitOffset, nil
}

// takes in a prefixed unit, outputs the name without prefixes
// and the scale factor specified by the prefixes
func parsePrefixes(unit string) (string, float64) {
	for k := range SiPrefixes {
		if strings.HasPrefix(unit, k) {
			return strings.TrimPrefix(unit, k), math.Pow10(SiPrefixes[k])
		}
	}
	return unit, 1.0
}
