package extractor

import (
	"go/ast"
	"go/token"
	"strconv"
)

// Извлекает строковое сообщение из вызова логгера.
// Поддерживаются только строковые аргументы
func ExtractMessageFromArg(call *ast.CallExpr, argIndex int) (string, bool) {

	if call == nil {
		return "", false
	}

	if len(call.Args) <= argIndex {
		return "", false
	}

	// Сначала попробуем получить простую строку
	if lit, ok := call.Args[argIndex].(*ast.BasicLit); ok && lit.Kind == token.STRING {
		msg, err := strconv.Unquote(lit.Value)
		if err == nil {
			return msg, true
		}
		return lit.Value, false
	}

	// Теперь попробуем получить конкатенацию строк
	if binary, ok := call.Args[argIndex].(*ast.BinaryExpr); ok && binary.Op == token.ADD {
		if msg, ok := extractFromBinary(binary); ok {
			return msg, true
		}
	}

	// Пробуем получить из вызова функции
	if callExpr, ok := call.Args[argIndex].(*ast.CallExpr); ok {
		if msg, ok := extractFromCall(callExpr); ok {
			return msg, true
		}
	}

	return "", false
}

// Для тех логов, которые есть условие сообщение хранится в первом аргументе
func ExtractMessage(call *ast.CallExpr) (string, bool) {
	return ExtractMessageFromArg(call, 0)
}

// Вспомогательные функции
func extractFromBinary(expr *ast.BinaryExpr) (string, bool) {

	var result string

	// Левая часть
	if lit, ok := expr.X.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		left, _ := strconv.Unquote(lit.Value)
		result += left
	} else {
		return "", false
	}

	// Правая часть
	switch y := expr.Y.(type) {
	case *ast.BasicLit:
		if y.Kind == token.STRING {
			right, _ := strconv.Unquote(y.Value)
			result += right
		} else {
			return "", false
		}
	case *ast.BinaryExpr:
		if right, ok := extractFromBinary(y); ok {
			result += right
		} else {
			return "", false
		}
	default:
		return "", false
	}

	return result, true
}

func extractFromCall(call *ast.CallExpr) (string, bool) {

	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "fmt" && sel.Sel.Name == "Sprintf" {
			if len(call.Args) > 0 {
				if format, ok := call.Args[0].(*ast.BasicLit); ok && format.Kind == token.STRING {
					msg, _ := strconv.Unquote(format.Value)
					return msg + " (format string)", true
				}
			}
		}
	}
	return "", false
}
