package pc

// 分离规则 deps:[前提,蕴含命题],如 [A,A->B]
func check_for_modus_penus(inference PCInference, pcChecker *PCChecker) bool {
	// fmt.Println("check for modus penus")
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
	premise = simplify_expr(premise)
	// TODO
	front, conclusion, ok := split_A_B_from_A2B(depstr)
	front = simplify_expr(front)
	conclusion = simplify_expr(conclusion)
	return ok && front == premise && inference.expr == conclusion
}
