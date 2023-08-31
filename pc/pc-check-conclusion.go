package pc

// 结论
func (pcChecker *PCChecker) CheckConclusion(conclusion string) bool {
	conclusion_struct, err := compileInference(conclusion)
	if err != nil {
		return false
	}
	expr := conclusion_struct.expr
	for _, inference := range pcChecker.inferences {
		// 检查是否在
		inf, err2 := compileInference(inference)
		inf_expr := inf.expr
		if err2 != nil {
			panic(pcChecker)
		}
		if inf_expr != expr {
			continue
		}
		// 判断该推理是否符合要求
		ok := true
		for _, dep := range inf.depends {
			if dep >= pcChecker.Len() {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
		return false
	}
	return false
}
