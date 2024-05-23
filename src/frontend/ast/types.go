package ast

type DATA_TYPE string

const (
	// Primitive Types
	VOID DATA_TYPE 	= "VOID"

	INTEGER   		= "INTEGER"
	FLOATING  		= "FLOATING"
	BOOLEAN   		= "BOOLEAN"
	STRING    		= "STRING"
	CHARACTER 		= "CHARECTER"
	NULL      		= "NULL"

	// Derived Types
	ARRAY 			= "ARRAY"
)

type IntegerType struct {
	Kind     DATA_TYPE
	BitSize  uint8
	IsSigned bool
}

func (i IntegerType) _type() {}

type FloatingType struct {
	Kind     DATA_TYPE
	BitSize  uint8
}

func (f FloatingType) _type() {}

type BooleanType struct {
	Kind     DATA_TYPE
}

func (b BooleanType) _type() {}

type StringType struct {
	Kind     DATA_TYPE
}

func (s StringType) _type() {}

type CharecterType struct {
	Kind     DATA_TYPE
}

func (c CharecterType) _type() {}

type NullType struct {
	Kind     DATA_TYPE
}

func (n NullType) _type() {}

type VoidType struct {
	Kind     DATA_TYPE
}

func (v VoidType) _type() {}

type ArrayType struct {
	Kind        DATA_TYPE
	ElementType Type
	Size        int
}

func (a ArrayType) _type() {}
