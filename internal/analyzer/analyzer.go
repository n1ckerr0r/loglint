package analyzer

import (
	"fmt"
	"go/ast"
	"log"
	"runtime/debug"

	"github.com/n1ckerr0r/loglint/internal/config"
	"github.com/n1ckerr0r/loglint/internal/detector"
	"github.com/n1ckerr0r/loglint/internal/extractor"
	"github.com/n1ckerr0r/loglint/internal/log_message"
	"github.com/n1ckerr0r/loglint/internal/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

func New(ruleSet []rules.Rule, cfg config.Config) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "loglint",
		Doc:  "checks log messages",
		Requires: []*analysis.Analyzer{
			inspect.Analyzer,
		},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return run(pass, ruleSet, cfg)
		},
	}
}

func run(pass *analysis.Pass, ruleSet []rules.Rule, cfg config.Config) (interface{}, error) {

	if pass.TypesInfo == nil {
		log.Println("no type info")
		return nil, nil
	}

	if len(ruleSet) == 0 {
		log.Println("no rules configured")
		return nil, nil
	}

	logDetector := detector.New()

	for _, file := range pass.Files {

		// Использует готовую реализацию ast дерева для обхода файла
		ast.Inspect(file, func(node ast.Node) bool {

			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			checkCall(pass, call, logDetector, &ruleSet, cfg)

			return true
		})
	}

	return nil, nil
}

func checkCall(
	pass *analysis.Pass,
	call *ast.CallExpr,
	logDetector *detector.Detector,
	ruleSet *[]rules.Rule,
	cfg config.Config,
) {

	// Лучше всего обрабатывать панику здесь, чтобы у пользователя было максимум информации и программа не падала
	defer func() {
		if r := recover(); r != nil {

			pass.Reportf(
				call.Pos(),
				"loglint internal error: %v\n%s",
				r,
				debug.Stack(),
			)
		}

	}()

	loggerCall, ok := logDetector.Detect(call, pass.TypesInfo)
	if !ok {
		return
	}

	text, ok := extractor.ExtractMessageFromArg(
		call,
		loggerCall.MsgIndex,
	)
	if !ok {
		return
	}

	msg := log_message.NewLogMessage(
		text,
		call.Pos(),
		loggerCall.LoggerType,
		loggerCall.Level,
	)

	for _, rule := range *ruleSet {

		err, fix := rule.Check(*msg)

		if err != nil {

			message := fmt.Sprintf("%s: %v", rule.Name(), err)

			if cfg.EnableSuggestedFix && fix != nil {
				message += fmt.Sprintf(" | fix: %s", fix.Message)

				if len(fix.TextEdits) > 0 {
					message += fmt.Sprintf(" : %s", string(fix.TextEdits[0].NewText))
				}
			}

			diagnostic := analysis.Diagnostic{
				Pos:     msg.Pos,
				Message: message,
			}

			if fix != nil {
				diagnostic.SuggestedFixes = []analysis.SuggestedFix{*fix}
			}

			pass.Report(diagnostic)
		}
	}
}
