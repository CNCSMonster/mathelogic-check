package pc

import (
	"fmt"
	"strings"
)

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
	strs[0] = simplify_expr(strs[0])
	premise = simplify_expr(premise)
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
