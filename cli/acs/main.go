package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Cantabar/friend-quest-ark-controller/core"
	"github.com/urfave/cli"
)

func main() {
	instances := core.GetAcsInstances()
	app := &cli.App{
		Name:  "ACS Controller CLI",
		Usage: "Conrtol ACS instances",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "map",
				Usage:       "Specify which ACS map to target",
				Value:       "island",
				DefaultText: "island",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{""},
				Usage:   "List available hosts",
				Action: func(c *cli.Context) error {
					fmt.Println("Available hosts are: ", core.GetAcsHostNames(instances))
					return nil
				},
			},
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Start an instance",
				Action: func(c *cli.Context) error {
					instanceIndex := core.FindInstanceByAcsHost(instances, c.Args().First())
					if instanceIndex != -1 {
						core.StartAcsInstance(instances[instanceIndex])
					} else {
						fmt.Println("Unable to find instance: ", c.Args().First())
						fmt.Println("Available hosts are: ", core.GetAcsHostNames(instances))
					}
					return nil
				},
			},
			{
				Name:    "stop",
				Aliases: []string{""},
				Usage:   "Stop an instance",
				Action: func(c *cli.Context) error {
					instanceIndex := core.FindInstanceByAcsHost(instances, c.Args().First())
					if instanceIndex != -1 {
						core.StopAcsInstance(instances[instanceIndex])
					} else {
						fmt.Println("Unable to find instance: ", c.Args().First())
						fmt.Println("Available hosts are: ", core.GetAcsHostNames(instances))
					}
					return nil
				},
			},
			{
				Name:    "status",
				Aliases: []string{""},
				Usage:   "Stop an instance",
				Action: func(c *cli.Context) error {
					instanceIndex := core.FindInstanceByAcsHost(instances, c.Args().First())
					if instanceIndex != -1 {
						instances[instanceIndex].ActivePlayers = core.GetActivePlayers(instances[instanceIndex].PublicIPAddress)
						fmt.Println(instances[instanceIndex])
					} else {
						fmt.Println("Unable to find instance: ", c.Args().First())
						fmt.Println("Available hosts are: ", core.GetAcsHostNames(instances))
					}
					return nil
				},
			},
			{
				Name:    "ssh",
				Aliases: []string{""},
				Usage:   "SSH to an instance",
				Action: func(c *cli.Context) error {
					instanceIndex := core.FindInstanceByAcsHost(instances, c.Args().First())
					if instanceIndex != -1 {
						fmt.Println(instances[instanceIndex])
						core.OpenSSHToInstance(instances[instanceIndex])
					} else {
						fmt.Println("Unable to find instance: ", c.Args().First())
						fmt.Println("Available hosts are: ", core.GetAcsHostNames(instances))
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
