package config

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"
	"unicode"
)

func TestNewConfigFromFile(t *testing.T) {
	var err error
	configPath := "../testdata/conf"

	cfg := &Config{"./conf", "./out/", 1, 1, 10, ".*", 4, nil}
	cfg.TargetUrlRE, err = regexp.Compile(cfg.TargetUrl)
	// write Config to file
	file, err := os.Create(configPath)
	if err != nil {
		t.Errorf("fail to open test configure file: %s\n", err)
		return
	}
	defer file.Close()

	val := reflect.ValueOf(*cfg)
	typ := reflect.TypeOf(*cfg)

	file.WriteString("[spider]")
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		// lower first letter of fieldName
		runes := []rune(fieldName)
		if len(runes) > 0 && unicode.IsUpper(runes[0]) {
			runes[0] = unicode.ToLower(runes[0])
		}
		fieldName = string(runes)
		fieldValue := field.Interface()
		_, err := file.WriteString(fmt.Sprintf("%s = %v\n", fieldName, fieldValue))
		if err != nil {
			t.Errorf("fail to write test Config to configure file: %s\n", err)
			return
		}
	}
	// build a new Config Based on file
	cfgFromFile, err := NewConfigFromFile(configPath)
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}
	// compare two Config
	if !reflect.DeepEqual(cfg, cfgFromFile) {
		t.Errorf("%v != %v", cfg, cfgFromFile)
	}
}

func TestValidateConfig(t *testing.T) {
	type args struct {
		cfg *Config
	}
	cfg1 := Config{"./conf", "./out/", 1, 1, 10, ".*", 4, nil}

	cfg2 := cfg1
	cfg2.UrlListFile = ""

	cfg3 := cfg1
	cfg3.ThreadCount = 0

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{&cfg1}, false},
		{"empty filepath string", args{&cfg2}, true},
		{"negative int", args{&cfg3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfig(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
