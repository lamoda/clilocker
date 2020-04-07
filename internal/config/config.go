package config

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/spf13/viper"
	"strings"
)

const (
	LocalFileScheme     = "local"
	RedisScheme         = "redis"
	RedisSentinelScheme = "redis-sentinel"

	NoLimit = 0
)

type Config struct {
	ConfigFile string `short:"f" long:"config"`
	LockerDsn  string `short:"d" long:"locker_dsn" mapstructure:"locker_dsn"`
	Limit      int    `short:"l" long:"limit" mapstructure:"limit"`
	CommandId  string `short:"i" long:"id" mapstructure:"command_id"`
	Command    string
	Args       []string
}

func New(args []string) (*Config, error) {
	options := defaultOptions()

	err := options.fillFromArgs(args)
	if err != nil {
		return nil, err
	}

	err = options.fillFromConfig()
	if err != nil {
		return nil, err
	}

	err = options.fillFromArgs(args)
	if err != nil {
		return nil, err
	}

	return options, nil
}

func defaultOptions() *Config {
	return &Config{
		ConfigFile: "",
		LockerDsn:  LocalFileScheme + "://",
		Limit:      NoLimit,
	}
}

func (c *Config) fillFromArgs(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("at least one command argument must be defined")
	}

	commandStart := -1
	for i, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			commandStart = i
			break
		}
	}

	if commandStart < 0 {
		return fmt.Errorf("at least one command argument must be defined")
	}

	_, err := flags.ParseArgs(c, args[:commandStart])
	if err != nil {
		return err
	}
	commandArgs := args[commandStart:]

	if c.CommandId == "" {
		c.CommandId = strings.Join(commandArgs, " ")
	}

	c.Command = commandArgs[0]
	c.Args = commandArgs[1:]

	return nil
}

func (c *Config) fillFromConfig() error {
	v := viper.New()

	if c.ConfigFile != "" {
		v.SetConfigFile(c.ConfigFile)
	} else {
		v.AddConfigPath(".")
		v.SetConfigName(".clilocker")
	}
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil
	}

	err = v.Unmarshal(c)

	if err != nil {
		return err
	}

	return nil
}
