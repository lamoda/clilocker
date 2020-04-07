package command

import (
	"context"
	"github.com/lamoda/clilocker/internal/config"
	"github.com/lamoda/clilocker/internal/lock"
	services2 "github.com/lamoda/clilocker/internal/services"
	"os"
	"os/exec"
)

type Command struct {
	cmd      *exec.Cmd
	config   *config.Config
	services *services2.Services
}

func New(config *config.Config, services *services2.Services) (*Command, error) {
	return &Command{
		cmd:      exec.CommandContext(context.Background(), config.Command, config.Args...),
		config:   config,
		services: services,
	}, nil
}

func (c *Command) Run() error {
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr

	unlock, err := c.tryToTakeALock(c.services.Locker, c.config.Limit)
	if err != nil {
		return err
	}
	if unlock == nil {
		return nil
	}
	defer unlock()

	return c.cmd.Run()
}

func (c *Command) tryToTakeALock(locker lock.Locker, limit int) (func(), error) {
	if limit < 1 {
		return func() {}, nil
	}

	lockInstance, err := locker.Lock(c.config.CommandId, limit)
	if err != nil {
		return nil, err
	}
	if !lockInstance.Taken() {
		return nil, nil
	}
	return func() {
		_ = lockInstance.Release()
	}, nil
}
