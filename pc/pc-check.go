package pc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
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

// 分离规则 deps:[前提,蕴含命题],如 [A,A->B]
func check_for_modus_penus(inference PCInference, pcChecker *PCChecker) bool {
	fmt.Println("check for modus penus")
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
	strs := strings.Split(depstr, "->")
	if len(strs) < 2 {
		return false
	}
	if strs[0] != premise {
		return false
	}
	strs = strs[1:]
	outcome := strings.Join(strs, "")
	if outcome == inference.expr {
		return true
	}
	return false
}

// 分离模式 A->(B->C)中的 A,B,C
func split_A_B_C(expr string) (A, B, C string, ok bool) {
	A, B, C = "", "", ""
	sb := strings.Builder{}
	state := 0
	numLeftParenthesis := 0 //自动机使用的参数
	ok = true
	fmt.Println("process", expr)
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

func check_for_third_axiom(inference PCInference, pcChecer *PCChecker) bool {
	fmt.Println("check for axiom three")
	return false
}
