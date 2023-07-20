package pc

import (
	"fmt"

	mathelogic "github.com/cncsmonster/mathelogic-check"
)

// pc只有一种推理规则,蕴含规则 A,A->B
// pc有三大公理:1.A->(B->A) 2.(A->(B->C)) -> (A->B)->(A->C) 3.
// pc只有两大联结词:1.-> 2.!(ps,为了方便打,使用!作为否定联结词)
// T
const (
	ONLY_ONE_MODUS_PONENS = iota
	AXIOM_ONE             = iota
	AXIOM_TWO             = iota
	AXIOM_TRHEE           = iota
)

// pc推理语句
type PCInference struct {
	// 序号
	order int
	// 字符串表达
	expr string
	// 使用到的规则(包括公理,推理规则,定理)
	rule int
	// 依赖的序号
	depends []int
}

// pc推理检查机
type PCChecker struct {
	premise    []string //前提
	inferences []string //非前提部分推理语句
}

func (pcChecker *PCChecker) Len() int {
	return len(pcChecker.premise) + len(pcChecker.inferences)
}

// 获取第i条语句,如果获取失败,
func (pcChecker *PCChecker) Get(index int) string {
	if index > pcChecker.Len() {
		return ""
	} else if index >= len(pcChecker.premise) {
		return pcChecker.inferences[index]
	} else {
		return pcChecker.premise[index]
	}
}

// 加入前提,加入前提不用检查,可以直接加入
func (pcChecker *PCChecker) PushPremise(premise string) {
	//
	pcChecker.premise = append(pcChecker.premise, premise)

}

// 加入推理语句 (加入推理语句之前要先进行检查)
func (pcChecker *PCChecker) PushInference(inference string) (bool, error) {
	inferenceStruct := compileInference(inference)
	//然后检查inference
	// TODO
	ok, err := pcChecker.checkInference(inferenceStruct)
	if err != nil {
		fmt.Println(err)
	}
	if ok {
		pcChecker.inferences = append(pcChecker.inferences, inferenceStruct.expr)
	}
	return ok, err
}

func (pcChecker *PCChecker) checkInference(inference PCInference) (bool, error) {
	return false, nil
}

// /首先对推理语句进行编译成基础语句单元
func compileInference(inference string) PCInference {
	return PCInference{}
}

// 弹出语句
func (pcChecker *PCChecker) PopInference() {
	pcChecker.inferences = pcChecker.inferences[:len(pcChecker.inferences)-1]
}

// 判断两个pc语句是否一致
func (pcChecker *PCChecker) Equal(inferenceOne, inferenceAnother string) bool {
	pcinfs1 := compileInference(inferenceOne)
	pcinfs2 := compileInference(inferenceAnother)
	if pcinfs1.expr != pcinfs2.expr {
		return false
	}
	if pcinfs1.rule != pcinfs2.rule {
		return false
	}
	if len(pcinfs1.depends) != len(pcinfs2.depends) {
		return false
	}
	for i, v := range pcinfs1.depends {
		if pcinfs2.depends[i] != v {
			return false
		}
	}

	return true
}

func init() {
	// TODO 初始化 操作
	var m mathelogic.Interface = &PCChecker{}
	m.PushInference("gg")
}
