package zap

import (
	"fmt"
	"log"
)

type Field struct {
	key   string
	value any
}

type Logger struct{}

func (l *Logger) Info(msg string, fields ...Field)  { log.Println(format("INFO", msg, fields...)) }
func (l *Logger) Error(msg string, fields ...Field) { log.Println(format("ERROR", msg, fields...)) }
func (l *Logger) Sync() error                       { return nil }

func String(key, value string) Field { return Field{key: key, value: value} }
func Error(err error) Field          { return Field{key: "error", value: err} }

type AtomicLevel struct{ level string }

func NewAtomicLevel() AtomicLevel                      { return AtomicLevel{level: "info"} }
func (a *AtomicLevel) UnmarshalText(text []byte) error { a.level = string(text); return nil }

type Config struct{ Level AtomicLevel }

func NewProductionConfig() Config        { return Config{Level: NewAtomicLevel()} }
func NewDevelopmentConfig() Config       { return Config{Level: NewAtomicLevel()} }
func (c Config) Build() (*Logger, error) { return &Logger{}, nil }

func format(level, msg string, fields ...Field) string {
	out := level + " " + msg
	for _, field := range fields {
		out += fmt.Sprintf(" %s=%v", field.key, field.value)
	}
	return out
}
