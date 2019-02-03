package babyenv

import (
	"os"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {
	type config struct {
		Debug  bool   `env:"DEBUG" envdefault:"false"`
		DBAddr string `env:"DB_ADDR" envdefault:"localhost:6379"`
		DBNum  int    `env:"DB_NUM"`
	}

	debugVal := true
	dbAddrVal := "example.com:6397"
	dbNumVal := 16

	os.Setenv("DEBUG", strconv.FormatBool(debugVal))
	os.Setenv("DB_ADDR", dbAddrVal)
	os.Setenv("DB_NUM", strconv.FormatInt(int64(dbNumVal), 10))

	var c config
	Parse(&c)

	if !c.Debug {
		t.Errorf("failed parsing bool; expected %v, got %v", debugVal, c.Debug)
	}
	if c.DBAddr != dbAddrVal {
		t.Errorf("failed parsing string; expected %v, got %v", dbAddrVal, c.DBAddr)
	}
	if c.DBNum != dbNumVal {
		t.Errorf("failed parsing int; expected %v, got %v", dbNumVal, c.DBNum)
	}
}

func TestParseWithDefaults(t *testing.T) {
	type config struct {
		Debug  bool   `env:"DEBUG" envdefault:"true"`
		DBAddr string `env:"DB_ADDR" envdefault:"localhost:6379"`
		DBNum  int    `env:"DB_NUM" envdefault:"16"`
	}

	debugVal := true
	dbAddrVal := "localhost:6379"
	dbNumVal := 16

	os.Unsetenv("DEBUG")
	os.Unsetenv("DB_ADDR")
	os.Unsetenv("DB_NUM")

	var c config
	Parse(&c)

	if c.Debug != debugVal {
		t.Errorf("failed parsing bool; expected %v, got %v", debugVal, c.Debug)
	}
	if c.DBAddr != dbAddrVal {
		t.Errorf("failed parsing string; expected %v, got %v", dbAddrVal, c.DBAddr)
	}
	if c.DBNum != dbNumVal {
		t.Errorf("failed parsing int; expected %v, got %v", dbNumVal, c.DBNum)
	}
}
