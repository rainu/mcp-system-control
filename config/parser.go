package config

import (
	"mcp-system-control/config/model"
	"os"

	"github.com/rainu/go-yacl"
)

const EnvironmentPrefix = "MCP_SYSTEM_CONTROL_"

func Parse(arguments []string, env []string) *model.Config {
	// for possible config file path
	cf := &model.ConfigFile{}
	config := yacl.NewConfig(cf, yacl.WithPrefixEnv(EnvironmentPrefix))
	handleErr(func() error { return config.ParseEnvironment(env...) })
	handleErr(func() error { return config.ParseArguments(arguments...) })

	c := &model.Config{}
	config = yacl.NewConfig(c, yacl.WithPrefixEnv(EnvironmentPrefix))
	config.ApplyDefaults()

	processYamlFiles(config, cf.Path)
	handleErr(func() error { return config.ParseEnvironment(env...) })
	handleErr(func() error { return config.ParseArguments(arguments...) })
	checkHelp(c, config)

	return c
}

func handleErr(f func() error) {
	if err := f(); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
