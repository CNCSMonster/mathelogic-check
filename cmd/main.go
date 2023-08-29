package main

import (
	"fmt"

	mathelogic "github.com/cncsmonster/mathelogic-check"
	"github.com/cncsmonster/mathelogic-check/pc"
)

// 检查pc推理机器
func pc_checker() {

}

func main() {
	//获取一个pc定理证明序列验证器
	var m mathelogic.Interface = pc.New()
	// 尝试证明自身推出
	s1 := "A->((B->A)->A)"
	s2 := "(A->((B->A)->A))->((A->(B->A))->(A->A))"
	// s3 := "A->(B->A)"
	// s4 := "((A->(B->A))->(A->A))"
	ok, err := m.PushInference(fmt.Sprintf("%s#1", s1))
	fmt.Println(ok, err)
	ok, err = m.PushInference(fmt.Sprintf("%s#2", s2))
	fmt.Println(ok, err)
	// ok, err = m.PushInference(fmt.Sprintf("%s#1", s3))
	// fmt.Println(ok, err)
	// ok, err = m.PushInference(fmt.Sprintf("%s#0#0,1", s4))
	// fmt.Println(ok, err)
	ok, err = m.PushInference("(!A->!B)->(B->A)#3")
	fmt.Println(ok, err)
}
