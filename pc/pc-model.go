package pc

import (
	"fmt"
)

// pc只有一种推理规则,蕴含规则 A,A->B
// pc有三大公理:1.A->(B->A) 2.(A->(B->C)) -> (A->B)->(A->C) 3.
// pc只有两大联结词:1.-> 2.!(ps,为了方便打,使用!作为否定联结词)
// T
const (
	// pc唯一推理规则,分离规则 若A,A->B,则B
	ONLY_ONE_MODUS_PONENS = iota
	// 公理1  (  A->(B->A))
	AXIOM_ONE
	// 公理2 （蕴含展开 (A->(B->C)) -> (A->B)->(A->C))
	AXIOM_TWO
	// 公理3  (双重否定 (!B->!A)->(A->B))
	AXIOM_TRHEE
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
	rules      map[int]func(PCInference, *PCChecker) bool
}

func New() *PCChecker {
	out := &PCChecker{premise: []string{}, inferences: []string{}, rules: make(map[int]func(PCInference, *PCChecker) bool)}
	out.rules[ONLY_ONE_MODUS_PONENS] = check_for_modus_penus
	out.rules[AXIOM_ONE] = check_for_first_axiom
	out.rules[AXIOM_TWO] = check_for_second_axiom
	out.rules[AXIOM_TRHEE] = check_for_third_axiom
	return out
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
	inferenceStructp, err := compileInference(inference)
	if err != nil {
		return false, err
	}
	inferenceStruct := *inferenceStructp
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

// 弹出语句
func (pcChecker *PCChecker) PopInference() {
	pcChecker.inferences = pcChecker.inferences[:len(pcChecker.inferences)-1]
}
