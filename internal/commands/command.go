package commands

import "github.com/urfave/cli"

var Commands = []cli.Command{
	{
		Name:    "queue",
		Aliases: []string{"q"},
		Usage:   "队列job",
		Subcommands: []cli.Command{
			{
				Name:   "consumer",
				Usage:  "消费程序",
				Action:  Queue,
			},
		},
	},
	{
		Name:    "db",
		Aliases: []string{"d"},
		Usage:   "数据",
		Subcommands: []cli.Command{
			{
				Name:   "migrate",
				Usage:  "数据迁移",
				Action:  Migrate,
			},
		},
	},
}
