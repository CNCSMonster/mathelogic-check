package pc

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 使用自定义规则进行验证
// 自定义规则的语法如下描述
/*
使用使用冒号(:)来分割语句(忽略所有空格以及换行等字符)
规则由若干语句来定义
第一行语句为规则原型,记录为 `表达式#依赖`,如:$2->($1->$3)#$4
之后若干行为规则语义
比如:$4=$1->($2->$3)

又或者只有规则原型(类似公理)
$1->$1
表示A->A等形式
*/
// 定义pc规则命题逻辑表达式的抽象语法树,以进行匹配
type pcRuleAst struct {
	left  *pcRuleAst
	right *pcRuleAst
	item  int
	kind  int //0:!,1:terminal,2:->
}
type pcRule struct {
	expr pcRuleAst
	deps []int
	eqs  [][]*pcRuleAst
}

// 注册规则,如果注册成功返回该规则所处的编号,注册失败返回失败原因
func (pcChecker *PCChecker) RegisterRule(rule string) (int, error) {
	cp := regexp.MustCompile("[\n\t ]")
	rule = cp.ReplaceAllString(rule, "")
	sentences := strings.Split(rule, ":")
	if len(sentences) == 0 {
		return -1, errors.New("syntax error!")
	}
	// 首先检查原型,取出依赖
	model := sentences[0]
	// 依赖
	ruleDeps := []int{}
	if strings.Contains(model, "#") {
		splits := strings.Split(model, "#")
		if len(splits) > 2 {
			return -1, errors.New("syntax error!")
		}
		model = splits[0]
		depStrs := strings.Split(splits[1], ",")
		for _, ds := range depStrs {
			if ds[0] != '$' || len(ds) < 2 {
				return -1, errors.New("syntax error!")
			} else {
				dep, err := strconv.Atoi(ds[1:])
				if err != nil {
					return -1, err
				}
				ruleDeps = append(ruleDeps, dep)
			}
		}
	}
	//
	pcast, err := compilePcRuleModel2PcRuleAst(model)
	if err != nil {
		return -1, err
	}
	// 对于等价规则进行提取
	eqs := [][]pcRuleAst{}
	for _, sentence := range sentences[1:] {
		split := strings.Split(sentence, "=")
		if len(split) != 2 {
			return -1, errors.New("syntax error!")
		}
		front, back := split[0], split[1]
		past1, err := compilePcRuleModel2PcRuleAst(front)
		if err != nil {
			return -1, fmt.Errorf(front, err)
		}
		past2, err2 := compilePcRuleModel2PcRuleAst(back)
		if err2 != nil {
			return -1, fmt.Errorf(back, err2)
		}
		eqs = append(eqs, []pcRuleAst{past1, past2})
	}
	// 取完依赖和模板,处理剩下的定义,根据剩下的定义生成结构
	newCheckFunc := func(pcInference PCInference, pcCheckerForNewRule *PCChecker) bool {
		// 首先提取依赖,对于依赖进行记录表
		deps := pcInference.depends
		item2Str := map[int]string{}
		if len(deps) != len(ruleDeps) {
			return false
		}
		// 根据依赖先初始化字符串表
		for index, dep := range deps {
			depInferenceStr, _ := pcCheckerForNewRule.get_expr(dep)
			depInference, _ := compileInference(depInferenceStr)
			depStr := depInference.expr
			item2Str[ruleDeps[index]] = depStr
		}
		// 然后 结合字符串表以及依赖比较ruleAst和 要检查的语句的expr
		match := compareRuleAstWithExprWithStrMap(&pcast, pcInference.expr, item2Str)
		if !match {
			return false
		}
		// 如果仍然匹配,则接下来基于规则进行匹配
		for _, eq := range eqs {
			front, back := eq[0], eq[1]
			fs := compileAst2Str(&front, item2Str)
			bs := compileAst2Str(&back, item2Str)
			fs, bs = simplify_expr(fs), simplify_expr(bs)
			if fs != bs {
				return false
			}
		}
		return true
	}
	// 注册新的检查函数
	pcChecker.rules[len(pcChecker.rules)] = newCheckFunc
	return len(pcChecker.rules) - 1, nil
}

// 把规则model表示式编译成pcRule对应的抽象语法树,编译失败返回失败error原因
func compilePcRuleModel2PcRuleAst(ruleModel string) (pcRuleAst, error) {
	out := pcRuleAst{left: nil, right: nil, item: 0, kind: -1}
	ruleModel = simplify_expr(ruleModel)
	if len(ruleModel) < 2 {
		return out, errors.New("syntax error!")
	}
	// 否定规则
	if ruleModel[0] == '!' {
		item, err := compilePcRuleModel2PcRuleAst(ruleModel[1:])
		if err != nil {
			return out, err
		}
		out.kind = 0
		out.left = &item
		return out, nil
	}
	// 蕴含规则
	if A, B, ok := split_A_B_from_A2B(ruleModel); ok {
		it, err := compilePcRuleModel2PcRuleAst(A)
		if err != nil {
			return out, err
		}
		it2, err2 := compilePcRuleModel2PcRuleAst(B)
		if err2 != nil {
			return out, err2
		}
		out.kind = 2
		out.left = &it
		out.right = &it2
		return out, nil
	}
	// 终点规则
	if ruleModel[0] == '$' {
		item, err := strconv.Atoi(ruleModel[1:])
		if err != nil {
			return out, err
		}
		out.kind = 1
		out.item = item
		return out, nil
	} else {
		return out, errors.New("syntax error!")
	}
}

// 检查expr是否符合pcModel,要求,输入的expre必须是经过化简去除至少最外层多余的括号的
func compareRuleAstWithExprWithStrMap(past *pcRuleAst, expr string, item2Strs map[int]string) bool {
	if past == nil {
		return false
	}
	// 比较终止情况,到终点了
	if past.kind == 1 {
		if target, ok := item2Strs[past.item]; !ok {
			item2Strs[past.item] = expr
			return true
		} else {
			return target == expr
		}
	} else if past.kind == 0 {
		// 否定情况
		if expr[0] != '!' {
			return false
		} else {
			expr = simplify_expr(expr[1:])
			return compareRuleAstWithExprWithStrMap(past.left, expr, item2Strs)
		}
	} else if past.kind == 2 {
		A, B, ok := split_A_B_from_A2B(expr)
		if !ok {
			return false
		}
		A, B = simplify_expr(A), simplify_expr(B)
		if !compareRuleAstWithExprWithStrMap(past.left, A, item2Strs) {
			return false
		}
		if !compareRuleAstWithExprWithStrMap(past.right, B, item2Strs) {
			return false
		}
		return true
	} else {
		return false
	}
}

// 使用量表由ast到字符串,如果得到的命题最外层为->类型,则会带上一对括号,比如(A->B)
func compileAst2Str(ast *pcRuleAst, item2Strs map[int]string) string {
	if ast == nil {
		return ""
	}
	if ast.kind == 1 {
		return item2Strs[ast.item]
	} else if ast.kind == 0 {
		return strings.Join([]string{"!", compileAst2Str(ast.left, item2Strs)}, "")
	} else if ast.kind == 2 {
		sb := strings.Builder{}
		sb.WriteRune('(')
		sb.WriteString(compileAst2Str(ast.left, item2Strs))
		sb.WriteString("->")
		sb.WriteString(compileAst2Str(ast.right, item2Strs))
		sb.WriteString(")")
		return sb.String()
	} else {
		panic(ast.kind)
	}
}
