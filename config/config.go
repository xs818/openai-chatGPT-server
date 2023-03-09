package config

import (
	"bytes"
	_ "embed"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type Configuration struct {
	Server  Server  `toml:"server"`
	ChatGPT ChatGPT `toml:"chatGPT"`
}

type Server struct {
	Name string `toml:"name"`
	Port int    `toml:"port"`
	Mode string `toml:"mode"`
}

type ChatGPT struct {
	// gpt apikey
	ApiKey string `toml:"apiKey" json:"apiKey"`
	// AI特征
	BotDesc string `toml:"botDesc" json:"botDesc"`
	// 代理
	Proxy string `toml:"proxy" json:"proxy"`
	// GPT请求最大字符数
	MaxTokens int `toml:"maxTokens" json:"maxTokens"`
	// GPT模型
	Model string `toml:"model" json:"model"`
	// 热度
	Temperature      float64 `toml:"temperature" json:"temperature"`
	TopP             float32 `toml:"topP" json:"topP"`
	PresencePenalty  float32 `toml:"presencePenalty" json:"presencePenalty"`
	FrequencyPenalty float32 `toml:"frequencyPenalty" json:"frequencyPenalty"`
}

var (
	//go:embed dev_config.toml
	devConfig []byte
)

func NewConfig(env Environment)  *Configuration{
	var config Configuration
	var r io.Reader
	switch env.Value() {
	case "dev":
		r = bytes.NewReader(devConfig)
	default:
		r = bytes.NewReader(devConfig)
	}

	// 优先使用环境变量的配置
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.SetEnvPrefix("CHAT")
	BindEnvs(config, "CHAT")

	viper.SetConfigType("toml")

	if err := viper.ReadConfig(r); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	viper.SetConfigName(Active().Value() + "_config")
	viper.AddConfigPath(currentAbPath())

	configFile := filepath.Join(currentAbPath(), Active().Value()+"_config.toml")
	_, ok := IsExists(configFile)
	if !ok {
		if err := os.MkdirAll(filepath.Dir(configFile), 0766); err != nil {
			panic(err)
		}

		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(&config); err != nil {
			panic(err)
		}
	})

	os.Setenv("name", config.Server.Name)

	return &config
}

func BindEnvs(config interface{}, prefix string) {
	value := reflect.ValueOf(config)
	valueType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := field.Type()
		if fieldType.Kind() == reflect.Struct {
			BindEnvs(field.Interface(), prefix)
		} else {
			tag := valueType.Field(i).Name
			if tag != "" {
				fieldName := valueType.Name() + "." + strings.ToLower(tag)
				envName := prefix +"_" + strings.ToUpper(tag)
				viper.BindEnv(fieldName, envName)
			}
		}
	}
}

func currentAbPath() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func IsExists(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
	return f, err == nil || os.IsExist(err)
}