package pc

import (
	"fmt"
	"regexp"
)

// 第一公理: A->(B->A)
// 规定文法 A,B 属于 [A-Ba-b]*
func check_for_first_axiom(inference PCInference, pcChecer *PCChecker) bool {
	// fmt.Println("check for atom one")
	A, B, C, ok := split_A_B_C(inference.expr)
	if !ok {
		return false
	}
	fmt.Println(A, B, C)
	// 对分离出来的三个成分进行检查,
	if A != C {
		return false
	}
	re := regexp.MustCompile("^[A-Za-z()->!]+$")
	return re.MatchString(A) && re.MatchString(B) && re.MatchString(C)
}
