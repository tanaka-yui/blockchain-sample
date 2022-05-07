package command

import (
	"blockchain/pkg/ossignal"
	"github.com/spf13/cobra"
	"os"
)

type Command struct {
	*cobra.Command
}

func NewCommand() *Command {
	return &Command{
		Command: &cobra.Command{},
	}
}

func (c *Command) WaitSignal(runner func()) os.Signal {
	c.Exec(runner)
	return <-ossignal.Quit()
}

func (c *Command) Exec(runner func()) {
	c.Run = func(cmd *cobra.Command, args []string) {
		runner()
	}
	if err := c.Execute(); err != nil {
		return
	}
	helpFlag := c.Flag("help")
	if helpFlag != nil && helpFlag.Changed {
		return
	}
}
