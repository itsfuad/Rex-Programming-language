package typechecker

import (
	"fmt"
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/frontend/parser"
	"walrus/helpers"
)

func EvaluateIdenitifierExpr(expr ast.IdentifierExpr, env *Environment) RuntimeValue {
	if !env.HasVariable(expr.Identifier) {

		msg := fmt.Sprintf("variable %v is not declared in this scope\n", expr.Identifier)

		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, msg).Display()
	}

	runtimeVal, err := env.GetRuntimeValue(expr.Identifier)

	if err != nil {
		parser.MakeError(env.parser, expr.StartPos.Line, env.parser.FilePath, expr.StartPos, expr.EndPos, err.Error()).Display()
	}

	return runtimeVal
}

func EvaluateUnaryExpression(unary ast.UnaryExpr, env *Environment) RuntimeValue {

	expr := Evaluate(unary.Argument, env)

	switch unary.Operator.Value {
	case "-", "+":
		if !helpers.TypesMatchT[IntegerValue](expr) {
			panic(fmt.Sprintf("Invalid unary operation for type %v", expr))
		}

		var value int64

		if unary.Operator.Value == "-" {
			value = -expr.(IntegerValue).Value
		} else {
			value = expr.(IntegerValue).Value
		}

		return MAKE_INT(value, 32, true)

	case "!":
		if !helpers.TypesMatchT[BooleanValue](expr) {
			panic(fmt.Sprintf("Invalid unary operation for type %v", expr))
		}

		return BooleanValue{
			Type: ast.BoolType{
				Kind: ast.T_BOOLEAN,
			},
			Value: !expr.(BooleanValue).Value,
		}

	case "++", "--":
		if !helpers.TypesMatchT[IntegerValue](expr) {
			panic(fmt.Sprintf("Invalid unary operation for type %v", expr))
		}

		var value int64

		if unary.Operator.Value == "++" {
			value = expr.(IntegerValue).Value + 1
		} else {
			value = expr.(IntegerValue).Value - 1
		}

		runtimeVal := MAKE_INT(value, 32, true)

		if helpers.TypesMatchT[ast.IdentifierExpr](unary.Argument) {
			if env.HasVariable(unary.Argument.(ast.IdentifierExpr).Identifier) {
				env.AssignVariable(unary.Argument.(ast.IdentifierExpr).Identifier, runtimeVal)
			}
		}

		return runtimeVal

	default:
		return MAKE_NULL()
	}
}

func EvaluateBinaryExpr(binop ast.BinaryExpr, env *Environment) RuntimeValue {

	left := Evaluate(binop.Left, env)
	right := Evaluate(binop.Right, env)

	leftType := GetRuntimeType(left)
	rightType := GetRuntimeType(right)

	errMsg := fmt.Sprintf("Unsupported binary operation between %v and %v", leftType, rightType)

	switch binop.Operator.Value {
	case "+", "-", "*", "/", "^":
		if IsINT(left) || IsFLOAT(left) && IsINT(right) || IsFLOAT(right) {
			// Numeric expr
			val, err := evaluateNumericExpr(left, right, binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val

		} else if helpers.ContainsIn([]string{"==", "!="}, binop.Operator.Value) {
			// eval string expr
			val, err := evaluateStringExpr(left.(StringValue), right.(StringValue), binop.Operator)
			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		} else if binop.Operator.Value == "+" && helpers.TypesMatchT[StringValue](left) && helpers.TypesMatchT[StringValue](right) {
			// eval string concat
			val, err := evaluateStringConcat(left.(StringValue), right.(StringValue))
			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		} else {

			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()

			return nil
		}

	case "==", "!=", ">", "<", ">=", "<=":
		if helpers.TypesMatchT[IntegerValue](left) && helpers.TypesMatchT[IntegerValue](right) {
			// Logical expr
			val, err := evaluateLogicalExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		} else if binop.Operator.Value == "==" || binop.Operator.Value == "!=" {
			// bool expr
			val, err := evaluateBoolExpr(left.(BooleanValue), right.(BooleanValue), binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		} else {

			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()

			return nil
		}

	case "+=", "-=", "*=", "/=", "%=":
		if !helpers.TypesMatchT[ast.IdentifierExpr](binop.Left) || !helpers.TypesMatchT[IntegerValue](right) {

			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()

			return nil
		}

		if env.HasVariable((binop.Left).(ast.IdentifierExpr).Identifier) {

			exprVal, err := evaluateNumericExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			runtimeVal, err := env.AssignVariable((binop.Left).(ast.IdentifierExpr).Identifier, exprVal)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return runtimeVal
		} else {
			val, err := evaluateNumericExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator)

			if err != nil {
				parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
			}

			return val
		}

	case "&&", "||":
		val, err := evaluateLogicalExpr(left.(IntegerValue), right.(IntegerValue), binop.Operator)

		if err != nil {
			parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.StartPos, binop.EndPos, err.Error()).Display()
		}

		return val

	default:

		parser.MakeError(env.parser, binop.StartPos.Line, env.parser.FilePath, binop.Operator.StartPos, binop.Operator.EndPos, errMsg).Display()

		return nil
	}
}

func EvaluateAssignmentExpr(assignNode ast.AssignmentExpr, env *Environment) RuntimeValue {

	var err error

	if assignNode.Assigne.Kind != ast.IDENTIFIER {
		err = fmt.Errorf("invalid left-hand side in assignment expression %v", assignNode.Assigne)
		parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, assignNode.Assigne.StartPos, assignNode.Assigne.EndPos, err.Error()).Display()
	}

	//if assigne is any of "false", "true", "null";
	if helpers.ContainsIn([]string{"false", "true", "null"}, assignNode.Assigne.Identifier) {
		err = fmt.Errorf("cannot assign to built-in constant %v", assignNode.Assigne.Identifier)
		parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, assignNode.Assigne.StartPos, assignNode.Assigne.EndPos, err.Error()).Display()
	}

	assigneValue, err := env.GetRuntimeValue(assignNode.Assigne.Identifier)

	if err != nil {
		valStart, valEnd := assignNode.Value.GetPos()
		parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, valStart, valEnd, err.Error()).Display()
	}

	value := Evaluate(assignNode.Value, env)

	switch assignNode.Operator.Kind {
	case lexer.PLUS_EQUALS_TOKEN, lexer.MINUS_EQUALS_TOKEN, lexer.TIMES_EQUALS_TOKEN, lexer.DIVIDE_EQUALS_TOKEN, lexer.MODULO_EQUALS_TOKEN:
		if !IsINT(assigneValue) || !IsINT(value) {

			err = fmt.Errorf("invalid operation between %v and %v", GetRuntimeType(assigneValue), GetRuntimeType(value))

			parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, assignNode.Operator.StartPos, assignNode.Operator.EndPos, err.Error()).Display()
		}

		value, err = evaluateNumericExpr(assigneValue.(IntegerValue), value.(IntegerValue), assignNode.Operator)

		if err != nil {

			valStart, valEnd := assignNode.Value.GetPos()

			parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, valStart, valEnd, err.Error()).Display()
		}
	}

	runtimeVal, err := env.AssignVariable(assignNode.Assigne.Identifier, value)

	if err != nil {
		start, end := assignNode.Value.GetPos()
		parser.MakeError(env.parser, assignNode.StartPos.Line, env.parser.FilePath, start, end, err.Error()).Display()
	}

	return runtimeVal
}

func evaluateNumericExpr(left RuntimeValue, right RuntimeValue, operator lexer.Token) (RuntimeValue, error) {

	if left != nil && right != nil {

		// evaluate both left, right as a, b where a and b can be int or float
		// if a is int return int, if a is float return float
		if IsINT(left) {
			if IsINT(right) {
				return evaluateIntInt(left.(IntegerValue), right.(IntegerValue), operator)
			} else {
				return evaluateIntFloat(left.(IntegerValue), right.(FloatValue), operator)
			}
		} else if IsFLOAT(left) {
			if IsINT(right) {
				return evaluateFloatInt(left.(FloatValue), right.(IntegerValue), operator)
			} else {
				return evaluateFloatFloat(left.(FloatValue), right.(FloatValue), operator)
			}
		}
	}
	return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
}

func evaluateIntInt(left IntegerValue, right IntegerValue, operator lexer.Token) (RuntimeValue, error) {
	switch operator.Value {
	case "+", "+=":
		return MAKE_INT(left.Value+right.Value, 32, true), nil
	case "-", "-=":
		return MAKE_INT(left.Value-right.Value, 32, true), nil
	case "*", "*=":
		return MAKE_INT(left.Value*right.Value, 32, true), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_INT(left.Value/right.Value, 32, true), nil
	case "%", "%=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_INT(left.Value%right.Value, 32, true), nil
	case "^":
		//power operation
		number := left.Value
		power := right.Value

		//use bit shifting to calculate power
		result := int64(1)
		for power > 0 {
			if power&1 == 1 {
				result *= number
			}
			number *= number
			power >>= 1
		}

		return MAKE_INT(result, 32, true), nil
	default:
		return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
	}
}

func evaluateIntFloat(left IntegerValue, right FloatValue, operator lexer.Token) (RuntimeValue, error) {
	switch operator.Value {
	case "+", "+=":
		return MAKE_INT(int64(float64(left.Value)+right.Value), 64, true), nil
	case "-", "-=":
		return MAKE_INT(int64(float64(left.Value)-right.Value), 64, true), nil
	case "*", "*=":
		return MAKE_INT(int64(float64(left.Value)*right.Value), 64, true), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_INT(int64(float64(left.Value)/right.Value), 64, true), nil
	case "^":
		//power operation
		number := float64(left.Value)
		power := right.Value

		//use bit shifting to calculate power
		result := 1.0
		for power > 0 {
			if int64(power)&1 == 1 {
				result *= number
			}
			number *= number
			power /= 2
		}

		return MAKE_INT(int64(result), 64, true), nil
	default:
		return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
	}
}

func evaluateFloatInt(left FloatValue, right IntegerValue, operator lexer.Token) (RuntimeValue, error) {
	switch operator.Value {
	case "+", "+=":
		return MAKE_FLOAT(left.Value+float64(right.Value), 64), nil
	case "-", "-=":
		return MAKE_FLOAT(left.Value-float64(right.Value), 64), nil
	case "*", "*=":
		return MAKE_FLOAT(left.Value*float64(right.Value), 64), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_FLOAT(left.Value/float64(right.Value), 64), nil
	case "^":
		//power operation
		number := left.Value
		power := float64(right.Value)

		//use bit shifting to calculate power
		result := 1.0
		for power > 0 {
			if int64(power)&1 == 1 {
				result *= number
			}
			number *= number
			power /= 2
		}

		return MAKE_FLOAT(result, 64), nil
	default:
		return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
	}
}

func evaluateFloatFloat(left FloatValue, right FloatValue, operator lexer.Token) (RuntimeValue, error) {

	switch operator.Value {
	case "+", "+=":
		return MAKE_FLOAT(left.Value+right.Value, 64), nil
	case "-", "-=":
		return MAKE_FLOAT(left.Value-right.Value, 64), nil
	case "*", "*=":
		return MAKE_FLOAT(left.Value*right.Value, 64), nil
	case "/", "/=":
		if right.Value == 0 {
			return nil, fmt.Errorf("division by zero is forbidden")
		}
		return MAKE_FLOAT(left.Value/right.Value, 64), nil
	case "^":
		//power operation
		number := left.Value
		power := right.Value

		//use bit shifting to calculate power
		result := 1.0
		for power > 0 {
			if int64(power)&1 == 1 {
				result *= number
			}
			number *= number
			power /= 2
		}

		return MAKE_FLOAT(result, 64), nil
	default:
		return nil, fmt.Errorf("cannot evaluate numeric operation. unsupported operator %v", operator.Value)
	}
}

func evaluateLogicalExpr(left IntegerValue, right IntegerValue, operator lexer.Token) (RuntimeValue, error) {

	result := false

	switch operator.Value {
	case "==":
		result = left.Value == right.Value
	case "!=":
		result = left.Value != right.Value
	case ">":
		result = left.Value > right.Value
	case "<":
		result = left.Value < right.Value
	case ">=":
		result = left.Value >= right.Value
	case "<=":
		result = left.Value <= right.Value
	case "&&":
		result = left.Value != 0 && right.Value != 0
	case "||":
		result = left.Value != 0 || right.Value != 0
	default:
		return nil, fmt.Errorf("cannot evaluate logical expression. unsupported operator %v", operator.Value)
	}

	return MAKE_BOOL(result), nil
}

func evaluateBoolExpr(left RuntimeValue, right RuntimeValue, operator lexer.Token) (RuntimeValue, error) {
	result := false

	sameType := helpers.TypesMatch(left, right)

	switch operator.Value {
	case "==":
		if sameType {
			result = left == right
		} else {
			result = false
		}
	case "!=":
		if sameType {
			result = left != right
		} else {
			result = true
		}
	case "&&":
		if sameType && helpers.TypesMatchT[BooleanValue](left) {
			result = left.(BooleanValue).Value && right.(BooleanValue).Value
		} else {
			result = IsTruthy(left) && IsTruthy(right)
		}
	case "||":
		if sameType && helpers.TypesMatchT[BooleanValue](left) {
			result = left.(BooleanValue).Value || right.(BooleanValue).Value
		} else {
			result = IsTruthy(left) || IsTruthy(right)
		}
	default:
		return nil, fmt.Errorf("cannot evaluate boolean expression. unsupported operator %v", operator)
	}

	return MAKE_BOOL(result), nil
}

func evaluateStringExpr(left StringValue, right StringValue, operator lexer.Token) (RuntimeValue, error) {
	result := false

	switch operator.Value {
	case "==":
		result = left.Value == right.Value
	case "!=":
		result = left.Value != right.Value
	default:
		return nil, fmt.Errorf("cannot evaluate string operation. unsupported operator %v", operator)
	}

	return MAKE_BOOL(result), nil
}

func evaluateStringConcat(left StringValue, right StringValue) (RuntimeValue, error) {
	return MAKE_STRING(left.Value + right.Value), nil
}
