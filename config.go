package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const defaultConfigYml = `PvsThreads: 4
BuildThreads: 4

# Specifies analyzer(s) and level(s) to be used for filtering, i.e.
# Set 0 for disable. Default is General = 1,2
Levels:
    # General Analysis
    General: 2

    # 64-bit errors;
    64-bit: 2

    # Micro-optimizations
    Optimizations: 2

    # Customers Specific Requests
    Customers: 0

    # MISRA
    MISRA: 0

Checks:
    # This file is marked with copyleft license
    - V::1042


    # Dangerous magic number
    - V::112

    # Qt warnings
    - V:qCDebug:1044
    - V:qCInfo:1044
    - V:qCWarning:1044
    - V:qCCritical:1044
`

type Config struct {
	PvsThreads   int           `yaml:"PvsThreads"`
	BuildThreads int           `yaml:"BuildThreads"`
	Levels       ConfigLevels  `yaml:"Levels"`
	Checks       []ConfigCheck `yaml:"Checks"`
}

type ConfigLevels struct {
	General       int `yaml:"General"`
	X64           int `yaml:"64-bit"`
	Optimizations int `yaml:"Optimizations"`
	Customers     int `yaml:"Customers"`
	MISRA         int `yaml:"MISRA"`
}

const DefaultConfigFile = ".pvscheck.yml"

type ConfigCheck string

func defaultConfig() (Config, error) {
	res := Config{
		PvsThreads:   0,
		BuildThreads: 0,
		Levels:       ConfigLevels{},
		Checks:       []ConfigCheck{},
	}
	err := yaml.Unmarshal([]byte(defaultConfigYml), &res)
	return res, err
}

func (c *Config) load(fileName string) error {

	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(f, c); err != nil {
		return err
	}

	return nil
}

func findConfigFiles(projectDir string) ([]string, error) {
	var err error
	dir, err := filepath.Abs(projectDir)
	if err != nil {
		return []string{}, err
	}

	res := []string{}
	path := strings.Split(dir, "/")

	for len(path) > 0 {
		f := strings.Join(path, "/") + DefaultConfigFile
		if fileExists(f) {
			res = append(res, f)
		}

		path = path[:len(path)-1]
	}

	return res, nil
}

func loadConfig(projectDir string) (Config, error) {
	files, err := findConfigFiles(projectDir)
	if err != nil {
		return Config{}, err
	}

	{
		file := os.Getenv("HOME") + "/.config/PVS-Studio/pvscheck.yml"
		if fileExists(file) {
			files = append(files, file)
		}
	}

	if len(files) == 0 {
		fmt.Println("The default config is used")
		return defaultConfig()
	}

	res := Config{}
	fmt.Println(files[0], " config file is used")
	res.load(files[0])
	//	cfg, err := ini.Load(fileName)
	//if err != nil {
	//return res, err
	//}

	//if err = cfg.MapTo(&res); err != nil {
	//return res, err
	//}

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
func (c Config) pvsMode() int {
	res := 0

	if c.Levels.X64 > 0 {
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
func (c Config) pvsLevels() string {
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

	addLevel(c.Levels.X64, "64")
	addLevel(c.Levels.General, "GA")
	addLevel(c.Levels.Optimizations, "OP")
	addLevel(c.Levels.Customers, "CS")
	addLevel(c.Levels.MISRA, "MISRA")

	return strings.Join(res, ";")
}
