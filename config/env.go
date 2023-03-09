package config

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	active Environment
	dev    Environment = &environment{value: "dev"}
	test   Environment = &environment{value: "test"}
	pro    Environment = &environment{value: "pro"}
	local  Environment = &environment{value: "local"}
)

var _ Environment = (*environment)(nil)

// Environment 环境配置
type Environment interface {
	Value() string
	IsDev() bool
	IsTest() bool
	IsPro() bool
	IsLocal() bool
}

type environment struct {
	value string
}

func (e *environment) Value() string {
	return e.value
}

func (e *environment) IsDev() bool {
	return e.value == "dev"
}

func (e *environment) IsTest() bool {
	return e.value == "test"
}

func (e *environment) IsPro() bool {
	return e.value == "pro"
}

func (e *environment) IsLocal() bool {
	return e.value == "local"
}

func LoadEnv(c *cli.Context) Environment {
	env := c.String("env")
	switch env {
	case "dev":
		os.Setenv("env", "dev")
		active = dev
	case "test":
		os.Setenv("env", "test")
		active = test
	case "pro":
		os.Setenv("env", "pro")
		active = pro
	case "local":
		os.Setenv("env", "local")
		active = local
	default:
		os.Setenv("env", "dev")
		active = dev
		fmt.Println("Warning: '-env' cannot be found, or it is illegal. The default 'dev' will be used.")
	}

	return active
}

// Active 当前配置的env
func Active() Environment {
	return active
}
