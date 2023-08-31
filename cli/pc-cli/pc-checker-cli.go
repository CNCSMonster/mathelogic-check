package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cncsmonster/mathelogic-check/pc"
	"github.com/urfave/cli/v2"
)

// 检查pc推理机器
func pc_checker_cli() {
	app := &cli.App{
		Name:   "pc checker",
		Usage:  "make an pc proof checker",
		Action: PcCheckAction,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func PcCheckAction(ctx *cli.Context) error {
	// 读取结论
	checker := pc.New()
	fmt.Println("input conclusion:")
	var conclusion string
	fmt.Scan(&conclusion)
	// 输入前提
	fmt.Println("input premises(use \\q to quit):")
	lines := []string{}
	for {
		var premise string
		fmt.Scan(&premise)
		if premise == "\\q" {
			break
		}
		checker.PushPremise(premise)
		fmt.Println(len(lines), premise)
		lines = append(lines, premise)
	}
	fmt.Println("input your inferences(\\q for quit):")
	for {
		var inference string = ""
		fmt.Scan(&inference)
		if inference == "\\q" {
			break
		}
		ok, err := checker.PushInference(inference)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !ok {
			panic(inference)
		}
		fmt.Println(len(lines), inference)
		lines = append(lines, inference)
		if checker.CheckConclusion(conclusion) {
			fmt.Println("prove successfully!")
			break
		}
	}

	return nil
}
