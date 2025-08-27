package model

type Filterer interface {
	Filter()
}

// Filter is a leave node in a filter tree.
// It represents a single condition that can be applied to a query.
type Filter struct {
	Column         string
	Value          any
	ComparisonType Comparison
}

func (f *Filter) Filter() {}

// FilterNode is a node in a filter tree.
// It can contain multiple Filter or FilterNode instances and represents a logical connection (And/Or)
type FilterNode struct {
	Nodes      []Filterer
	Connection ConnectionType
}

func (f *FilterNode) Filter() {}

type Comparison int64

const (
	Equal Comparison = iota
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
	NotEqual
	Like
	ILike
)

type ConnectionType int64

const (
	And ConnectionType = 0
	Or  ConnectionType = 1
)
