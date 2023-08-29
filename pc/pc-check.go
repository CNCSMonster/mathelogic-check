package pc

import (
	"errors"
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
