package configs

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	t.Setenv(EVNameConfigFile, "../.env.dev")

	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "success",
			want: NewConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringValue(t *testing.T) {
	t.Setenv("key", "value")

	type args struct {
		configFile map[string]string
		configKey  string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success_getStringValue",
			args: args{
				configFile: map[string]string{
					"key": "value",
				},
				configKey: "key",
			},
			want: "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStringValue(tt.args.configFile, tt.args.configKey); got != tt.want {
				t.Errorf("getStringValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getIntValue(t *testing.T) {
	t.Setenv("key", "1")

	type args struct {
		configFile map[string]string
		configKey  string
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "success_getIntValue",
			args: args{
				configFile: map[string]string{
					"key": "1",
				},
				configKey: "key",
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getIntValue(tt.args.configFile, tt.args.configKey); got != tt.want {
				t.Errorf("getIntValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBoolValue(t *testing.T) {
	t.Setenv("key", "true")

	type args struct {
		configFile map[string]string
		configKey  string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success_getBoolValue",
			args: args{
				configFile: map[string]string{
					"key": "true",
				},
				configKey: "key",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBoolValue(tt.args.configFile, tt.args.configKey); got != tt.want {
				t.Errorf("getBoolValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStringArrayValue(t *testing.T) {
	t.Setenv("key", "first,second,third")

	type args struct {
		configFile map[string]string
		configKey  string
	}

	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "success_getStringArrayValue",
			args: args{
				configFile: map[string]string{
					"key": "first,second,third",
				},
				configKey: "key",
			},
			want: []string{"first", "second", "third"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStringArrayValue(tt.args.configFile, tt.args.configKey); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getStringArrayValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
