package pc

import "fmt"

// 第三规则  (!A->!B) -> (B -> A)
func check_for_third_axiom(inference PCInference, pcChecer *PCChecker) bool {
	//
	nA2nB, B2A, ok := split_A_B_from_A2B(inference.expr)
	nA2nB, B2A = simplify_expr(nA2nB), simplify_expr(B2A)
	if !ok {
		return false
	}
	fmt.Println(nA2nB, B2A)
	B, A, ok3 := split_A_B_from_A2B(B2A)
	if !ok3 {
		return false
	}
	nA2nB_mode := fmt.Sprintf("!%s->!%s", A, B)
	fmt.Println(nA2nB_mode, nA2nB)
	return nA2nB_mode == nA2nB
}