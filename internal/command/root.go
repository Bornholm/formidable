package command

import "github.com/urfave/cli/v2"

func Root() []*cli.Command {
	return []*cli.Command{
		Edit(),
		Set(),
		Get(),
		Delete(),
		Check(),
	}
}
