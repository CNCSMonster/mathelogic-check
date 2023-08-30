package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// 检查pc推理机器
func pc_checker_cli() {
	app := &cli.App{
		Name:  "pc checker",
		Usage: "make an pc proof checker",
		Action: func(ctx *cli.Context) error {
			fmt.Println("start!")

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
