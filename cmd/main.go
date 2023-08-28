package main

import (
	"fmt"

	mathelogic "github.com/cncsmonster/mathelogic-check"
	"github.com/cncsmonster/mathelogic-check/pc"
)

func main() {
	//获取一个pc定理证明序列验证器
	var m mathelogic.Interface = pc.New()
	m.PushPremise("A->B")
	m.PushPremise("A")
	ok, err := m.PushInference("B#0#1,0")
	if err != nil {
		fmt.Println(err)
	} else if ok {
		fmt.Println("推理成功")
	} else {
		fmt.Println("推理失败")
	}
}
