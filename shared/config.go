package shared

import (
	"log"
	"os"
	"strconv"
)

type Settings interface {
	GetBool(key string) (bool, bool)
	GetBoolOrDefault(key string, b bool) bool
	GetFloat(key string) (float64, bool)
	GetFloatOrDefault(key string, f float64) float64
	GetInt(key string) (int, bool)
	GetIntOrDefault(key string, i int) int
	GetString(key string) (string, bool)
	GetStringOrDefault(key string, s string) string
}

type EnvSettings struct{}

func NewSettings() Settings {
	return &EnvSettings{}
}

func (e *EnvSettings) GetBoolOrDefault(key string, b bool) bool {
	v, ok := e.GetBool(key)
	if !ok {
		v = b
	}
	return v
}

func (e *EnvSettings) GetFloatOrDefault(key string, f float64) float64 {
	v, ok := e.GetFloat(key)
	if !ok {
		v = f
	}
	return v
}

func (e *EnvSettings) GetIntOrDefault(key string, i int) int {
	v, ok := e.GetInt(key)
	if !ok {
		v = i
	}
	return v
}

func (e *EnvSettings) GetStringOrDefault(key string, s string) string {
	v, ok := e.GetString(key)
	if !ok {
		v = s
	}
	return v
}

func (e *EnvSettings) GetBool(key string) (bool, bool) {
	val, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		log.Printf("cannot read key %s, error: %s\n", key, err)
		return val, false
	}
	return val, true
}

func (e *EnvSettings) GetFloat(key string) (float64, bool) {
	val, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil {
		log.Printf("cannot read key %s, error: %s\n", key, err)
		return val, false
	}
	return val, true
}

func (e *EnvSettings) GetInt(key string) (int, bool) {
	val, err := strconv.ParseInt(os.Getenv(key), 10, 32)
	if err != nil {
		log.Printf("cannot read key %s, error: %s\n", key, err)
		return int(val), false
	}
	return int(val), true
}

func (e *EnvSettings) GetString(key string) (string, bool) {
	val := os.Getenv(key)
	return val, val != ""
}
