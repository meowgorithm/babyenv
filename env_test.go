package babyenv

import (
	"os"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {
	type config struct {
		Debug  bool   `env:"DEBUG" default:"false"`
		DBAddr string `env:"DB_ADDR" default:"localhost:6379"`
		DBNum  int    `env:"DB_NUM"`
	}

	debugVal := true
	dbAddrVal := "example.com:6397"
	dbNumVal := 16

	os.Setenv("DEBUG", strconv.FormatBool(debugVal))
	os.Setenv("DB_ADDR", dbAddrVal)
	os.Setenv("DB_NUM", strconv.FormatInt(int64(dbNumVal), 10))

	var c config
	if err := Parse(&c); err != nil {
		t.Errorf("error while parsing: %v", err)
		return
	}

	if !c.Debug {
		t.Errorf("failed parsing bool; expected %#v, got %#v", debugVal, c.Debug)
	}
	if c.DBAddr != dbAddrVal {
		t.Errorf("failed parsing string; expected %#v, got %#v", dbAddrVal, c.DBAddr)
	}
	if c.DBNum != dbNumVal {
		t.Errorf("failed parsing int; expected %#v, got %#v", dbNumVal, c.DBNum)
	}
}

func TestParseWithDefaults(t *testing.T) {
	type config struct {
		Debug  bool   `env:"DEBUG" default:"true"`
		DBAddr string `env:"DB_ADDR" default:"localhost:6379"`
		DBNum  int    `env:"DB_NUM" default:"16"`
	}

	debugVal := true
	dbAddrVal := "localhost:6379"
	dbNumVal := 16

	os.Unsetenv("DEBUG")
	os.Unsetenv("DB_ADDR")
	os.Unsetenv("DB_NUM")

	var c config
	if err := Parse(&c); err != nil {
		t.Errorf("error while parsing: %v", err)
		return
	}

	if c.Debug != debugVal {
		t.Errorf("failed parsing bool; expected %#v, got %#v", debugVal, c.Debug)
	}
	if c.DBAddr != dbAddrVal {
		t.Errorf("failed parsing string; expected %#v, got %#v", dbAddrVal, c.DBAddr)
	}
	if c.DBNum != dbNumVal {
		t.Errorf("failed parsing int; expected %#v, got %#v", dbNumVal, c.DBNum)
	}
}

func TestParsePointers(t *testing.T) {
	type config struct {
		A *string `env:"A"`
		B *bool   `env:"B"`
		C *int    `env:"C"`
	}

	a := "xxx"
	b := true
	c := 1
	os.Setenv("A", a)
	os.Setenv("B", strconv.FormatBool(b))
	os.Setenv("C", strconv.FormatInt(int64(c), 10))

	var cfg config
	if err := Parse(&cfg); err != nil {
		t.Errorf("error while parsing: %v", err)
		return
	}

	if cfg.A == nil {
		t.Errorf("failed parsing *string; expected %#v, got nil", a)
	} else if *cfg.A != a {
		t.Errorf("failed parsing *string; expected %#v, got %#v", a, *cfg.A)
	}

	if cfg.B == nil {
		t.Errorf("failed parsing *bool; expected %#v, got nil", b)
	} else if *cfg.B != b {
		t.Errorf("failed parsing *bool; expected %#v, got %#v", b, *cfg.B)
	}

	if cfg.C == nil {
		t.Errorf("failed parsing *int; expected %#v, got nil", c)
	} else if *cfg.C != c {
		t.Errorf("failed parsing *int; expected %#v, got %#v", c, *cfg.C)
	}
}
