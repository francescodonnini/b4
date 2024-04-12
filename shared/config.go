package shared

import (
	"log"
	"os"
	"strconv"
)

type Settings interface {
	GetBool(key string) (bool, bool)
	GetFloat(key string) (float64, bool)
	GetInt(key string) (int, bool)
	GetString(key string) (string, bool)
}

type EnvSettings struct{}

func NewSettings() Settings {
	return &EnvSettings{}
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
