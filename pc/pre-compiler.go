package pc

import (
	"errors"
	"strconv"
	"strings"
)

// 把推理语句进行编译成基础语句单元
// 成分 字面表达式#使用的规则#依赖的表达式序号
func compileInference(inference string) (*PCInference, error) {
	//去掉所有空格
	inferences := strings.Split(inference, " ")
	inference = strings.Join(inferences, "")
	// 首先读取直接表达式
	strs := strings.Split(inference, "#")
	if len(strs) == 0 {
		return nil, errors.New("wrong syntax")
	}
	out := PCInference{}
	out.expr = strs[0]
	// 解析规则
	if len(strs) >= 2 {
		rule, err := strconv.Atoi(strs[1])
		if err != nil {
			return nil, errors.New("expect int for rule mark")
		}
		out.rule = rule
	}
	// 解析依赖的语句
	if len(strs) >= 3 {
		deps := strings.Split(strs[2], ",")
		out.depends = []int{}
		for _, v := range deps {
			val, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.New("expect int order for depends")
			}
			out.depends = append(out.depends, val)
		}
	}
	return &out, nil
}
