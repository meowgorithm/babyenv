package babyenv

import (
	"os"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {
	type config struct {
		A bool   `env:"A"`
		B string `env:"B"`
		C int    `env:"C"`
	}

	a := true
	b := "xxx"
	c := 16

	os.Setenv("A", strconv.FormatBool(a))
	os.Setenv("B", b)
	os.Setenv("C", strconv.FormatInt(int64(c), 10))

	var cfg config
	if err := Parse(&cfg); err != nil {
		t.Errorf("error while parsing: %v", err)
		return
	}

	if !cfg.A {
		t.Errorf("failed parsing bool; expected %#v, got %#v", a, cfg.A)
	}
	if cfg.B != b {
		t.Errorf("failed parsing string; expected %#v, got %#v", b, cfg.B)
	}
	if cfg.C != c {
		t.Errorf("failed parsing int; expected %#v, got %#v", c, cfg.C)
	}
}

func TestParseWithDefaults(t *testing.T) {
	type config struct {
		A bool   `env:"A" default:"true"`
		B string `env:"B" default:"xxx"`
		C int    `env:"C" default:"16"`
	}

	a := true
	b := "xxx"
	c := 16

	os.Unsetenv("A")
	os.Unsetenv("B")
	os.Unsetenv("C")

	var cfg config
	if err := Parse(&cfg); err != nil {
		t.Errorf("error while parsing: %v", err)
		return
	}

	if cfg.A != a {
		t.Errorf("failed parsing bool; expected %#v, got %#v", a, cfg.A)
	}
	if cfg.B != b {
		t.Errorf("failed parsing string; expected %#v, got %#v", b, cfg.B)
	}
	if cfg.C != c {
		t.Errorf("failed parsing int; expected %#v, got %#v", c, cfg.C)
	}
}

func TestParsePointers(t *testing.T) {
	type config struct {
		A *bool   `env:"A"`
		B *string `env:"B"`
		C *int    `env:"C"`
	}

	a := true
	b := "xxx"
	c := 16
	os.Setenv("A", strconv.FormatBool(a))
	os.Setenv("B", b)
	os.Setenv("C", strconv.FormatInt(int64(c), 10))

	var cfg config
	if err := Parse(&cfg); err != nil {
		t.Errorf("error while parsing: %v", err)
		return
	}

	if cfg.A == nil {
		t.Errorf("failed parsing *bool; expected %#v, got nil", a)
	} else if *cfg.A != a {
		t.Errorf("failed parsing *bool; expected %#v, got %#v", a, *cfg.A)
	}

	if cfg.B == nil {
		t.Errorf("failed parsing *string; expected %#v, got nil", b)
	} else if *cfg.B != b {
		t.Errorf("failed parsing *string; expected %#v, got %#v", b, *cfg.B)
	}

	if cfg.C == nil {
		t.Errorf("failed parsing *int; expected %#v, got nil", c)
	} else if *cfg.C != c {
		t.Errorf("failed parsing *int; expected %#v, got %#v", c, *cfg.C)
	}
}
