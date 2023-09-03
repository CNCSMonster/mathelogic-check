package pc

import (
	"fmt"
	"strings"
)

// 第三规则  (!A->!B) -> (B -> A)
func check_for_third_axiom(inference PCInference, pcChecer *PCChecker) bool {
	//
	nA2nB, B2A, ok := split_A_B_from_A2B(inference.expr)
	nA2nB, B2A = simplify_expr(nA2nB), simplify_expr(B2A)
	if !ok {
		return false
	}
	B, A, ok3 := split_A_B_from_A2B(B2A)
	if !ok3 {
		return false
	}
	B, A = simplify_expr(B), simplify_expr(A)
	if strings.Contains(B, "->") {
		B = fmt.Sprintf("(%s)", B)
	}
	if strings.Contains(A, "->") {
		A = fmt.Sprintf("(%s)", A)
	}
	nA2nB_mode := fmt.Sprintf("!%s->!%s", A, B)
	return nA2nB_mode == nA2nB
}
