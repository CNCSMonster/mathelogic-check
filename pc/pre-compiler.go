package pc

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// 把推理语句进行编译成基础语句单元
// 成分 字面表达式#使用的规则#依赖的表达式序号
func compileInference(inference string) (*PCInference, error) {
	inference = simplify_expr(inference)
	// 首先读取直接表达式
	strs := strings.Split(inference, "#")
	if len(strs) == 0 {
		return nil, errors.New("wrong syntax")
	}
	out := PCInference{}
	out.expr = simplify_expr(strs[0])
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

// 表达式化简, 去除两边多余括号,以及两边多余空格
// 表达式最简情况下,最外层最多只能有1个->,其他->必须在()中
func simplify_expr(base_expr string) string {
	outs := strings.Clone(base_expr)
	// 判断是否首尾存在成对的可以删除的括号
	for len(outs) > 2 && outs[0] == '(' && outs[len(outs)-1] == ')' {
		num_left := 1
		i := 1
		for i < len(outs)-1 {
			c := outs[i]
			if c == '(' {
				num_left += 1
			} else if c == ')' {
				num_left -= 1
			}
			if num_left == 0 {
				break
			}
			i += 1
		}
		if i != len(outs)-1 {
			break
		}
		outs = outs[1 : len(outs)-1]
	}

	inferences := strings.Split(outs, " ")
	outs = strings.Join(inferences, "")
	return outs
}

// 分离模式 A->(B->C)中的 A,B,C
func split_A_B_C(expr string) (A, B, C string, ok bool) {
	A, B2C, ok := split_A_B_from_A2B(expr)
	if !ok {
		return "", "", "", false
	}
	B2C = simplify_expr(B2C)
	B, C, ok2 := split_A_B_from_A2B(B2C)
	if !ok2 {
		return "", "", "", false
	}
	return A, B, C, true
}

// 从模式 A->B 中提取 A和B
func split_A_B_from_A2B(expr string) (A string, B string, ok bool) {
	A = ""
	B = ""
	num_left := 0
	state := 1
	ok = false
	sb := strings.Builder{}
	for i, c := range expr {
		switch state {
		// 处理左串
		// 处理left
		case 1:
			if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '!' || c == '$' {
				sb.WriteRune(c)
			} else if c == '(' {
				num_left += 1
				sb.WriteRune(c)
			} else if c == ')' && num_left > 0 {
				num_left -= 1
				sb.WriteRune(c)
			} else if c == '-' && num_left == 0 {
				A = sb.String()
				sb = strings.Builder{}
				state = 2
			} else if (c == '-' || c == '>') && num_left != 0 {
				sb.WriteRune(c)
			} else {
				break
			}
		// 中转2
		case 2:
			if c == '>' {
				state = 3
			} else {
				break
			}
		// 处理right
		case 3:
			if i == len(expr)-1 {
				if c == ')' && num_left == 1 {
					sb.WriteRune(c)
				} else if c != ')' && num_left == 0 {
					sb.WriteRune(c)
				} else {
					break
				}
				ok = true
				B = sb.String()
			} else if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '!' || c == '$' {
				sb.WriteRune(c)
			} else if c == '(' {
				num_left += 1
				sb.WriteRune(c)
			} else if c == ')' && num_left > 0 {
				sb.WriteRune(c)
				num_left -= 1
			} else if (c == '-' || c == '>') && num_left != 0 {
				sb.WriteRune(c)
			} else {
				break
			}
		}
	}
	return A, B, ok
}
