package mathelogiccheck

// ps,约定推理序列，包括前提,第一个语句的下标为0
type Interface interface {
	PushPremise(premise string)                       //无条件推入前提
	PushInference(inference string) (bool, error)     //如果推理语句能够加入,推理成功,如果不能够加入,推理失败
	PopInference()                                    //移除最后一条推理，如果推理序列为空,则不移除
	Len() int                                         //获取推理序列(包括前提的长度)
	Get(index int) string                             //获取某个位置的语句
	Equal(inferenceOne, inferenceAnother string) bool //判断两个推理语句是否是等价的,如果是,说明
}
