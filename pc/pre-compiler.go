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
	A, B, C = "", "", ""
	sb := strings.Builder{}
	state := 0
	numLeftParenthesis := 0 //自动机使用的参数
	ok = true
	// fmt.Println("process", expr)
	for _, c := range expr {
		switch state {
		// 获取A
		case 0:
			if unicode.IsLetter(c) || c == '!' {
				sb.WriteRune(c)
			} else if c == '(' {
				sb.WriteRune(c)
				numLeftParenthesis += 1
			} else if c == ')' && numLeftParenthesis > 0 {
				sb.WriteRune(c)
				numLeftParenthesis -= 1
			} else if c == '-' && numLeftParenthesis == 0 {
				A = sb.String()
				sb = strings.Builder{}
				state = 1
			} else if c == '-' || c == '>' {
				sb.WriteRune(c)
			} else {
				ok = false
				break
			}
		// A,B中转
		case 1:
			if c == '>' {
				state = 2
			} else {
				ok = false
				break
			}
		// 进入(B->A)
		case 2:
			if c == '(' {
				state = 3
			} else {
				ok = false
				break
			}
		//准备B
		case 3:
			if unicode.IsLetter(c) || c == '!' {
				sb.WriteRune(c)
			} else if c == '(' {
				sb.WriteRune(c)
				numLeftParenthesis += 1
			} else if c == ')' && numLeftParenthesis > 0 {
				sb.WriteRune(c)
				numLeftParenthesis -= 1
			} else if c == '-' && numLeftParenthesis == 0 {
				B = sb.String()
				sb = strings.Builder{}
				state = 4
			} else if c == '-' || c == '>' {
				sb.WriteRune(c)
			} else {
				ok = false
				break
			}
		//B和第二个A的中转
		case 4:
			if c == '>' {
				state = 5
			} else {
				ok = false
				break
			}
		case 5:
			if unicode.IsLetter(c) || c == '!' {
				sb.WriteRune(c)
			} else if c == '(' {
				sb.WriteRune(c)
				numLeftParenthesis += 1
			} else if c == ')' && numLeftParenthesis > 0 {
				sb.WriteRune(c)
				numLeftParenthesis -= 1
			} else if c == ')' && numLeftParenthesis == 0 {
				C = sb.String()
				state = 6
			} else if c == '-' && numLeftParenthesis == 0 {
				ok = false
				break
			} else if c == '-' || c == '>' {
				sb.WriteRune(c)
			} else {
				ok = false
				break
			}
		case 6:
			// 到达终结状态后后面不应该存在字符
			ok = false
			break
		}
	}
	return A, B, C, ok
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
			if unicode.IsLetter(c) || c == '!' {
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
			} else if unicode.IsLetter(c) || c == '!' {
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
