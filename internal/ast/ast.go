package ast

type Program struct {
	Name  string
	Setup *Setup
	Run   *Run
}

type Setup struct {
	Declarations []Declaration
}

type Run struct {
	Instructions []Instruction
}

type Declaration interface {
	declarationNode()
}

type VarDeclaration struct {
	Names     []string
	Type      DataType
	InitValue string
	IsArray   bool
	ArraySize int
	Line      int
	Column    int
}

func (*VarDeclaration) declarationNode() {}

type ConstDeclaration struct {
	Name   string
	Type   DataType
	Value  string
	Line   int
	Column int
}

func (*ConstDeclaration) declarationNode() {}

type Instruction interface {
	instructionNode()
}

type Assignment struct {
	Target        string
	IsArrayAccess bool
	Index         Expression
	Value         Expression
	Line          int
	Column        int
}

func (*Assignment) instructionNode() {}

type IfStatement struct {
	Condition Condition
	ThenBlock []Instruction
	ElseBlock []Instruction
	Line      int
	Column    int
}

func (*IfStatement) instructionNode() {}

type WhileLoop struct {
	Condition Condition
	Body      []Instruction
	Line      int
	Column    int
}

func (*WhileLoop) instructionNode() {}

type ForLoop struct {
	Variable string
	From     Expression
	To       Expression
	Body     []Instruction
	Line     int
	Column   int
}

func (*ForLoop) instructionNode() {}

type InFunction struct {
	Variable string
	Line     int
	Column   int
}

func (*InFunction) instructionNode() {}

type OutFunction struct {
	Arguments []interface{} // can be string or identifier
	Line      int
	Column    int
}

func (*OutFunction) instructionNode() {}

type Expression interface {
	expressionNode()
	getType() DataType
}

type IntConstant struct {
	Value int
}

func (*IntConstant) expressionNode()      {}
func (ic *IntConstant) getType() DataType { return TYPE_INTEGER }

type FloatConstant struct {
	Value float64
}

func (*FloatConstant) expressionNode()      {}
func (fc *FloatConstant) getType() DataType { return TYPE_FLOAT }

type Identifier struct {
	Name          string
	IsArrayAccess bool
	Index         Expression
}

func (*Identifier) expressionNode() {}
func (id *Identifier) getType() DataType {
	return TYPE_UNKNOWN
}

type BinaryOp struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (*BinaryOp) expressionNode() {}
func (bo *BinaryOp) getType() DataType {
	leftType := bo.Left.getType()
	rightType := bo.Right.getType()
	return MergeTypes(leftType, rightType)
}

type UnaryOp struct {
	Operator string
	Operand  Expression
}

func (*UnaryOp) expressionNode() {}
func (uo *UnaryOp) getType() DataType {
	return uo.Operand.getType()
}

type Condition interface {
	conditionNode()
}

type ComparisonCondition struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (*ComparisonCondition) conditionNode() {}

type LogicalCondition struct {
	Left     Condition
	Operator string
	Right    Condition
}

func (*LogicalCondition) conditionNode() {}

type NotCondition struct {
	Operand Condition
}

func (*NotCondition) conditionNode() {}
