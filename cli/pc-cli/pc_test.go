package main

import (
	"fmt"
	"testing"

	"github.com/cncsmonster/mathelogic-check/pc"
	"github.com/stretchr/testify/assert"
)

// 验证A->A的证明
func TestA2A(t *testing.T) {
	checker := pc.New()
	conclusion := "A->A"
	proofs := []string{
		"A->(B->A)#1",
		"A->((B->A)->A)#1",
		"(A->((B->A)->A))->((A->(B->A))->(A->A))#2",
		"(A->(B->A))->(A->A)#0#1,2",
		"A->A#0#0,3",
	}
	for _, inference := range proofs {
		ok, err := checker.PushInference(inference)
		assert.True(t, err == nil, err)
		assert.True(t, ok, inference)
	}
	assert.True(t, checker.CheckConclusion(conclusion))
}

// 验证 A->(B->C) then B->(A->C)  (换前件规则证明)
func TestSwapPreditionRule(t *testing.T) {
	checker := pc.New()
	conclusion := "B->(A->C)"
	premises := []string{
		"A->(B->C)",
	}
	proofs := []string{
		"(A->(B->C))->((A->B)->(A->C))#2",
		"(A->B)->(A->C)#0#0,1",
		"((A->B)->(A->C))->(B->((A->B)->(A->C)))#1",
		"B->((A->B)->(A->C))#0#2,3",
		"(B->((A->B)->(A->C)))->((B->(A->B))->(B->(A->C)))#2",
		"((B->(A->B))->(B->(A->C)))#0#4,5",
		"B->(A->B)#1",
		"B->(A->C)#0#7,6",
	}
	for _, premise := range premises {
		checker.PushPremise(premise)
	}
	for _, inference := range proofs {
		ok, err := checker.PushInference(inference)
		assert.True(t, err == nil, err)
		assert.True(t, ok, inference)
	}
	assert.True(t, checker.CheckConclusion(conclusion))
}

// 测试第三公理
func TestThirdAxiom(t *testing.T) {
	checker := pc.New()
	conclusion := "B->A"
	premises := []string{
		"!A->!B",
	}
	proofs := []string{
		"(!A->!B)->(B->A)#3",
		"B->A#0#0,1",
	}
	for _, premise := range premises {
		checker.PushPremise(premise)
	}
	for _, inference := range proofs {
		ok, err := checker.PushInference(inference)
		assert.True(t, err == nil, err)
		assert.True(t, ok, inference)
	}
	assert.True(t, checker.CheckConclusion(conclusion))
}

// 测试规则的注册
func TestRegisterRule(t *testing.T) {
	// 测试注册加前件规则
	checker := pc.New()
	ruleId, err := checker.RegisterRule(pc.SWapPreRule)
	assert.Nil(t, err, err)
	conclusion := "B->(A->C)"
	premises := []string{
		"A->(B->C)",
	}
	proofs := []string{
		fmt.Sprintf("B->(A->C)#%d#0", ruleId),
	}
	for _, premise := range premises {
		checker.PushPremise(premise)
	}
	for _, inference := range proofs {
		ok, err := checker.PushInference(inference)
		assert.True(t, err == nil, err)
		assert.True(t, ok, inference)
	}
	assert.True(t, checker.CheckConclusion(conclusion))
}
