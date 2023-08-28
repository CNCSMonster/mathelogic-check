package pc

import (
	"errors"
	"fmt"
	"strings"
)

func (pcChecker *PCChecker) checkInference(inference PCInference) (bool, error) {
	// 检查使用的 规则/公理/定理 是否存在
	rule, ok := pcChecker.rules[inference.rule]
	if !ok {
		return false, errors.New("使用的推理规则/公理/定理不存在")
	}
	ok = rule(inference, pcChecker)
	return ok, nil
}

// 根据序号获取语句
func (pcChecker *PCChecker) get_expr(index int) (string, error) {
	if index < len(pcChecker.premise) {
		return pcChecker.premise[index], nil
	} else if index < pcChecker.Len() {
		return pcChecker.inferences[index-len(pcChecker.premise)], nil
	} else {
		return "", errors.New("index out of bound")
	}
}

// 分离规则 deps:[前提,蕴含命题],如 [A,A->B]
func check_for_modus_penus(inference PCInference, pcChecker *PCChecker) bool {
	fmt.Println("check for modus penus")
	if inference.rule != ONLY_ONE_MODUS_PONENS {
		return false
	}
	if len(inference.depends) != 2 {
		return false
	}
	premise_str, err := pcChecker.get_expr(inference.depends[0])
	if err != nil {
		return false
	}
	premise_struct, err2 := compileInference(premise_str)
	if err2 != nil {
		panic("")
	}
	premise := premise_struct.expr
	//
	dep_struct, err3 := pcChecker.get_expr(inference.depends[1])
	if err3 != nil {
		return false
	}
	dep_expr, err4 := compileInference(dep_struct)
	if err4 != nil {
		panic("")
	}
	depstr := dep_expr.expr
	strs := strings.Split(depstr, "->")
	if len(strs) < 2 {
		return false
	}
	if strs[0] != premise {
		return false
	}
	strs = strs[1:]
	outcome := strings.Join(strs, "")
	if outcome == inference.expr {
		return true
	}
	return false
}

func check_for_first_axiom(inference PCInference, pcChecer *PCChecker) bool {
	fmt.Println("check for atom one")
	return false
}

func check_for_second_axiom(inference PCInference, pcChecer *PCChecker) bool {
	fmt.Println("check for axiom two")
	return false
}

func check_for_third_axiom(inference PCInference, pcChecer *PCChecker) bool {
	fmt.Println("check for axiom three")
	return false
}
