package engine

import (
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	PvsThreads   int          `ini:"PVS threads"`
	BuildThreads int          `ini:"Build threads"`
	Levels       ConfigLevels `ini:"LEVELS"`
}

type ConfigLevels struct {
	General       int
	x64           int `ini:"64-bit"`
	Optimizations int
	Customers     int
	MISRA         int
}

func NewConfig() (Config, error) {
	res := Config{
		PvsThreads:   4,
		BuildThreads: 4,
		Levels: ConfigLevels{
			General:       2,
			x64:           2,
			Optimizations: 2,
			Customers:     0,
			MISRA:         0,
		},
	}

	fileName := os.Getenv("HOME") + "/.config/PVS-Studio/pvscheck.conf"
	if !fileExists(fileName) {
		return res, nil
	}

	cfg, err := ini.Load(fileName)
	if err != nil {
		return res, err
	}

	if err = cfg.MapTo(&res); err != nil {
		return res, err
	}

	return res, nil
}

/* **********************************************
 	MODE defines the type of warnings:
	    1 - 64-bit errors;
        2 - reserved;
        4 - General Analysis;
        8 - Micro-optimizations;
        16 - Customers Specific Requests;
        32 - MISRA.
        Modes can be combined by adding the values
 ********************************************** */
func (c Config) PvsMode() int {
	res := 0

	if c.Levels.x64 > 0 {
		res += 1
	}

	if c.Levels.General > 0 {
		res += 4
	}

	if c.Levels.Optimizations > 0 {
		res += 8
	}

	if c.Levels.Customers > 0 {
		res += 16
	}

	if c.Levels.MISRA > 0 {
		res += 32
	}

	return res
}

/* **********************************************
	Specifies analyzer(s) and level(s)
	to be used for filtering, i.e.
		'GA:1,2;64:1;OP:1,2,3;CS:1;MISRA:1,2'
	Default: GA:1,2
********************************************** */
func (c Config) PvsLevels() string {
	res := []string{}

	addLevel := func(v int, key string) {
		if v == 1 {
			res = append(res, key+":1")
		}

		if v == 2 {
			res = append(res, key+":1,2")
		}

		if v == 3 {
			res = append(res, key+":1,2,3")
		}
	}

	addLevel(c.Levels.x64, "64")
	addLevel(c.Levels.General, "GA")
	addLevel(c.Levels.Optimizations, "OP")
	addLevel(c.Levels.Customers, "CS")
	addLevel(c.Levels.MISRA, "MISRA")

	return strings.Join(res, ";")
}
