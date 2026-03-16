package detector

import (
	"go/ast"
	"go/types"
)

type LoggerCall struct {
	LoggerType string
	Level      string
	MsgIndex   int
}

// Сущность для определения является ли вызов логгером
type Detector struct {
	loggers map[string]map[string]int
}

func New() *Detector {
	return &Detector{
		loggers: make(map[string]map[string]int),
	}
}

func (detector *Detector) Detect(call *ast.CallExpr, info *types.Info) (*LoggerCall, bool) {

	if call == nil || info == nil {
		return nil, false
	}

	selector, ok := getSelector(call)
	if !ok {
		return nil, false
	}

	level := selector.Sel.Name

	if logger, ok := detectSlog(selector, level); ok {
		return logger, true
	}

	if lc, ok := detectZap(selector, level, info); ok {
		return lc, true
	}

	return nil, false
}

// Вспомогательные функции для анализа узла дерево и конкретного лога
func getSelector(call *ast.CallExpr) (*ast.SelectorExpr, bool) {

	if call == nil || call.Fun == nil {
		return nil, false
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	return selector, true
}

func detectSlog(selector *ast.SelectorExpr, level string) (*LoggerCall, bool) {

	if selector == nil || selector.X == nil {
		return nil, false
	}

	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return nil, false
	}

	if ident.Name != "slog" {
		return nil, false
	}

	switch level {
	case "Debug", "Info", "Warn", "Error":
		return &LoggerCall{
			LoggerType: "slog",
			Level:      level,
			MsgIndex:   0,
		}, true
	}

	return nil, false
}

func detectZap(selector *ast.SelectorExpr, level string, info *types.Info) (*LoggerCall, bool) {

	if selector == nil || selector.X == nil || info == nil {
		return nil, false
	}

	if !isZapLogger(selector.X, info) {
		return nil, false
	}

	switch level {
	case "Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal":
		return &LoggerCall{
			LoggerType: "zap",
			Level:      level,
			MsgIndex:   0,
		}, true
	}

	return nil, false
}

func isZapLogger(expr ast.Expr, info *types.Info) bool {

	if expr == nil || info == nil {
		return false
	}

	typ := info.TypeOf(expr)
	if typ == nil {
		return false
	}

	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	return obj.Pkg().Path() == "go.uber.org/zap" &&
		obj.Name() == "Logger"
}
