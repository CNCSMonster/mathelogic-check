package pc

import (
	"fmt"
	"strings"
)

// 第二公理, (A->(B->C)) -> ((A->B)->(A->C))
func check_for_second_axiom(inference PCInference, pcChecer *PCChecker) bool {
	// fmt.Println("check for axiom two")
	// fmt.Println("process", inference.expr)
	// 首先提取出左侧,然后提取出右侧
	// 然后对两侧
	left, right, ok := split_A_B_from_A2B(inference.expr)
	if !ok {
		return false
	}
	// fmt.Println("left:", left, "right:", right)
	left, right = simplify_expr(left), simplify_expr(right)

	A, B2C, ok := split_A_B_from_A2B(left)
	if !ok {
		return false
	}
	// fmt.Println("A:", A, "B2C:", B2C)
	B, C, ok2 := split_A_B_from_A2B(simplify_expr(B2C))
	B, C = simplify_expr(B), simplify_expr(C)
	if !ok2 {
		return false
	}
	if strings.Contains(B, "->") {
		B = fmt.Sprintf("(%s)", B)
	}
	if strings.Contains(A, "->") {
		A = fmt.Sprintf("(%s)", A)
	}
	if strings.Contains(C, "->") {
		C = fmt.Sprintf("(%s)", C)
	}
	// 对于右部的模式,使用左部组装进行匹配
	right_mod := fmt.Sprintf("(%s->%s)->(%s->%s)", A, B, A, C)
	fmt.Println(right, right_mod)
	return right_mod == right
}
