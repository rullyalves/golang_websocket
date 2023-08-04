package websocket

type Expression interface {
	Match(context interface{}) bool
}

type AndExpression struct {
	expressions []Expression
}

func (r *AndExpression) Match(context interface{}) bool {
	for _, expression := range r.expressions {
		return expression.Match(context)
	}
	return true
}

type OrExpression struct {
	expressions []Expression
}

func (r *OrExpression) Match(context interface{}) bool {
	for _, expression := range r.expressions {
		return expression.Match(context)
	}
	return true
}

type GreaterThanExpression struct {
	value interface{}
}

func (r *GreaterThanExpression) Match(context interface{}) bool {
	return true
}

type LessThanExpression struct {
	value interface{}
}

func (r *LessThanExpression) Match(context interface{}) bool {
	return true
}

type StringExpression struct {
	Equals   string
	Contains string
	In       []string
}

type MessageFilter struct {
	name        *StringExpression
	description *StringExpression
	and         []*MessageFilter
	or          []*MessageFilter
	not         *MessageFilter
}
