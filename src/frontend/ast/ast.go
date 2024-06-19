package ast

import (
	"walrus/frontend/lexer"
)

type NODE_TYPE string

const (
	// Statements
	PROGRAM NODE_TYPE = "PROGRAM"

	MODULE_STATEMENT NODE_TYPE = "MODULE_STATEMENT"
	IMPORT_STATEMENT NODE_TYPE = "IMPORT_STATEMENT"

	BLOCK_STATEMENT                NODE_TYPE = "BLOCK_STATEMENT"
	VARIABLE_DECLARATION_STATEMENT NODE_TYPE = "VARIABLE_DECLARATION_STATEMENT"
	CONTROL_FLOW_STATEMENT         NODE_TYPE = "CONTROL_FLOW_STATEMENT"
	WHILE_STATEMENT                NODE_TYPE = "WHILE_STATEMENT"
	SWITCH_STATEMENT               NODE_TYPE = "SWITCH_STATEMENT"
	SWITCH_CASE_STATEMENT          NODE_TYPE = "SWITCH_CASE_STATEMENT"
	DEFAULT_CASE_STATEMENT         NODE_TYPE = "DEFAULT_CASE_STATEMENT"
	IF_STATEMENT                   NODE_TYPE = "IF_STATEMENT"
	ELSE_STATEMENT                 NODE_TYPE = "ELSE_STATEMENT"
	FOR_LOOP_STATEMENT             NODE_TYPE = "FOR_LOOP_STATEMENT"
	FOREACH_LOOP_STATEMENT         NODE_TYPE = "FOREACH_LOOP_STATEMENT"
	FN_DECLARATION_STATEMENT       NODE_TYPE = "FN_DECLARATION_STATEMENT"
	FN_PROTOTYPE_STATEMENT         NODE_TYPE = "FN_PROTOTYPE_STATEMENT"
	RETURN_STATEMENT               NODE_TYPE = "RETURN_STATEMENT"
	BREAK_STATEMENT                NODE_TYPE = "BREAK_STATEMENT"
	CONTINUE_STATEMENT             NODE_TYPE = "CONTINUE_STATEMENT"
	TRAIT_STATEMENT                NODE_TYPE = "TRAIT_STATEMENT"
	STRUCT_STATEMENT               NODE_TYPE = "STRUCT_STATEMENT"
	IMPLEMENTS_STATEMENT           NODE_TYPE = "IMPLEMENTS_STATEMENT"

	// Literals
	NUMERIC_LITERAL   NODE_TYPE = "NUMERIC_LITERAL"
	STRING_LITERAL    NODE_TYPE = "STRING_LITERAL"
	CHARACTER_LITERAL NODE_TYPE = "CHARACTER_LITERAL"
	BOOLEAN_LITERAL   NODE_TYPE = "BOOLEAN_LITERAL"
	NULL_LITERAL      NODE_TYPE = "NULL_LITERAL"
	VOID_LITERAL      NODE_TYPE = "VOID_LITERAL"
	ARRAY_LITERALS    NODE_TYPE = "ARRAY_LITERALS"
	STRUCT_LITERAL    NODE_TYPE = "STRUCT_LITERAL"

	STRUCT_PROPERTY NODE_TYPE = "STRUCT_PROPERTY"

	// Expressions
	ASSIGNMENT_EXPRESSION NODE_TYPE = "ASSIGNMENT_EXPRESSION"
	IDENTIFIER            NODE_TYPE = "IDENTIFIER"
	BINARY_EXPRESSION     NODE_TYPE = "BINARY_EXPRESSION"
	LOGICAL_EXPRESSION    NODE_TYPE = "LOGICAL_EXPRESSION"

	// Functions
	FUNCTION_PARAMETER NODE_TYPE = "FUNCTION_PARAMETER"

	FUNCTION_CALL_EXPRESSION NODE_TYPE = "FUNCTION_CALL_EXPRESSION"

	// Unary Operations
	UNARY_EXPRESSION NODE_TYPE = "UNARY_EXPRESSION"
)

type Node interface {
	iNode()
	GetPos() (lexer.Position, lexer.Position)
}

type Statement interface {
	Node
	iStatement()
}

type Expression interface {
	Node
	iExpression()
}

type Type interface {
	IType() DATA_TYPE
}

type BaseStmt struct {
	Kind     NODE_TYPE
	StartPos lexer.Position
	EndPos   lexer.Position
}
