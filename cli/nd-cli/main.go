package main

import (
	"fmt"

	"github.com/cncsmonster/mathelogic-check/pc"
)

func main() {
	var checker *pc.PCChecker = pc.New()
	checker.PushPremise("A")
	checker.PushPremise("A->B")
	checker.PushInference("B#0#0,1")
	fmt.Println(checker.Len())
	ok := checker.CheckConclusion("B#0#0,1")
	fmt.Println(ok)
	pc_checker_cli()
}
