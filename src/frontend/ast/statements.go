package ast

import "rexlang/frontend/lexer"

type ModuleStmt struct {
	Kind       NODE_TYPE
	ModuleName string
	StartPos   lexer.Position
	EndPos     lexer.Position
}
func (m ModuleStmt) node() {} // implements the Statement interface
func (m ModuleStmt) GetPos() (lexer.Position, lexer.Position) {
	return m.StartPos, m.EndPos
}
func (m ModuleStmt) _statement() {}

type ImportStmt struct {
	Kind       	NODE_TYPE
	ModuleName	string
	Identifiers []string
	StartPos   	lexer.Position
	EndPos     	lexer.Position
}

func (i ImportStmt) node() {} // implements the Statement interface
func (i ImportStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i ImportStmt) _statement() {}

type ProgramStmt struct {
	FileName   string
	ModuleName string
	Imports    []ImportStmt
	Contents   []Node
	StartPos   lexer.Position
	EndPos     lexer.Position
}

func (p ProgramStmt) node() {} // implements the Statement interface
func (p ProgramStmt) GetPos() (lexer.Position, lexer.Position) {
	return p.StartPos, p.EndPos
}
func (p ProgramStmt) _statement() {}

type BlockStmt struct {
	Kind     NODE_TYPE
	Body     []Node
	StartPos lexer.Position
	EndPos   lexer.Position
}

func (b BlockStmt) node() {}
func (b BlockStmt) GetPos() (lexer.Position, lexer.Position) {
	return b.StartPos, b.EndPos
}
func (b BlockStmt) _statement() {}

type VariableDclStml struct {
	Kind         NODE_TYPE
	IsConstant   bool
	Identifier   string
	Value        Expression
	ExplicitType Type
	StartPos     lexer.Position
	EndPos       lexer.Position
}

func (v VariableDclStml) node() {}
func (v VariableDclStml) GetPos() (lexer.Position, lexer.Position) {
	return v.StartPos, v.EndPos
}
func (v VariableDclStml) _statement() {}

type FunctionDeclStmt struct {
	Kind         NODE_TYPE
	FunctionName string
	Parameters   map[string]Type
	ReturnType   Type
	Block        BlockStmt
	StartPos     lexer.Position
	EndPos       lexer.Position
}

func (f FunctionDeclStmt) node() {}
func (f FunctionDeclStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f FunctionDeclStmt) _statement() {}

type ReturnStmt struct {
	Kind       NODE_TYPE
	Expression Expression
	StartPos   lexer.Position
	EndPos     lexer.Position
}

func (r ReturnStmt) node() {}
func (r ReturnStmt) GetPos() (lexer.Position, lexer.Position) {
	return r.StartPos, r.EndPos
}
func (r ReturnStmt) _statement() {}









type StructProperty struct {
	IsStatic bool
	IsPublic bool
	ReadOnly bool
	Type     Type
	StartPos lexer.Position
	EndPos   lexer.Position
}

type StructMethod struct {
	IsStatic   bool
	IsPublic   bool
	Parameters map[string]Type
	ReturnType Type
	StartPos   lexer.Position
	EndPos     lexer.Position
}

type StructDeclStatement struct {
	Kind       NODE_TYPE
	StructName string
	Properties map[string]StructProperty
	Methods    map[string]StructMethod
	StartPos   lexer.Position
	EndPos     lexer.Position
}
func (s StructDeclStatement) node() {}
func (s StructDeclStatement) GetPos() (lexer.Position, lexer.Position) {
	return s.StartPos, s.EndPos
}
func (s StructDeclStatement) _statement() {}



type IfStmt struct {
	Kind      NODE_TYPE
	Condition Expression
	Block     BlockStmt
	Alternate interface{}
	StartPos  lexer.Position
	EndPos    lexer.Position
}

func (i IfStmt) node() {}
func (i IfStmt) GetPos() (lexer.Position, lexer.Position) {
	return i.StartPos, i.EndPos
}
func (i IfStmt) _statement() {}

type ForStmt struct {
	Kind       NODE_TYPE
	Variable   string
	Init       Expression
	Condition  Expression
	Post	   Expression
	Block      BlockStmt
	StartPos   lexer.Position
	EndPos     lexer.Position
}
func (f ForStmt) node() {}
func (f ForStmt) GetPos() (lexer.Position, lexer.Position) {
	return f.StartPos, f.EndPos
}
func (f ForStmt) _statement() {}