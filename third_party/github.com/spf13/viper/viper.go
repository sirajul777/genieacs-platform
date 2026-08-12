package viper

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Viper struct {
	values   map[string]any
	prefix   string
	replacer *strings.Replacer
	file     string
}

func New() *Viper                                      { return &Viper{values: map[string]any{}} }
func (v *Viper) SetEnvPrefix(prefix string)            { v.prefix = strings.ToUpper(prefix) }
func (v *Viper) SetEnvKeyReplacer(r *strings.Replacer) { v.replacer = r }
func (v *Viper) AutomaticEnv()                         {}
func (v *Viper) SetDefault(key string, value any)      { v.values[key] = value }
func (v *Viper) SetConfigFile(file string)             { v.file = file }
func (v *Viper) ReadInConfig() error {
	f, err := os.Open(v.file)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	section := ""
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, ":") {
			section = strings.TrimSuffix(line, ":")
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if section != "" {
			key = section + "." + key
		}
		v.values[key] = val
	}
	return s.Err()
}
func (v *Viper) Unmarshal(target any) error {
	rv := reflect.ValueOf(target).Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		prefix := f.Tag.Get("mapstructure")
		fv := rv.Field(i)
		if fv.Kind() == reflect.Struct {
			fillStruct(fv, prefix, v)
		}
	}
	return nil
}
func fillStruct(rv reflect.Value, prefix string, v *Viper) {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		key := prefix + "." + f.Tag.Get("mapstructure")
		fv := rv.Field(i)
		if fv.Kind() == reflect.Struct {
			fillStruct(fv, key, v)
			continue
		}
		switch fv.Kind() {
		case reflect.String:
			if val, ok := lookup(v, key); ok {
				fv.SetString(val)
			}
		case reflect.Int:
			if val, ok := lookup(v, key); ok {
				parsed, _ := strconv.Atoi(val)
				fv.SetInt(int64(parsed))
			}
		}
	}
}
func lookup(v *Viper, key string) (string, bool) {
	envKey := strings.ToUpper(key)
	if v.replacer != nil {
		envKey = v.replacer.Replace(envKey)
	}
	if v.prefix != "" {
		envKey = v.prefix + "_" + envKey
	}
	if val, ok := os.LookupEnv(envKey); ok {
		return val, true
	}
	if val, ok := v.values[key]; ok {
		switch typed := val.(type) {
		case string:
			return typed, true
		case int:
			return strconv.Itoa(typed), true
		}
	}
	return "", false
}
