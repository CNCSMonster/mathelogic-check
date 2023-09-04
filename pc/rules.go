package pc

// 可选规则列表
const (
	A2ARule string = `
	$1->$1
	`
	SWapPreRule string = `
	$1->($2->$3)#$4:
	$4=$2->($1->$3)
	`
)
