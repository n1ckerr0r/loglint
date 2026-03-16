package detector

import (
	"go/ast"
	"go/types"
	"testing"
)

// Постарался написать максимально читаемыми
func TestDetectSlog_Simple(t *testing.T) {

	tests := []struct {
		name      string
		level     string
		identName string
		wantType  string
		wantOk    bool
	}{
		{
			name:      "slog debug",
			level:     "Debug",
			identName: "slog",
			wantType:  "slog",
			wantOk:    true,
		},
		{
			name:      "slog info",
			level:     "Info",
			identName: "slog",
			wantType:  "slog",
			wantOk:    true,
		},
		{
			name:      "не slog (другое имя)",
			level:     "Info",
			identName: "logger",
			wantOk:    false,
		},
		{
			name:      "неподдерживаемый уровень",
			level:     "Trace",
			identName: "slog",
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := &ast.SelectorExpr{
				X: &ast.Ident{Name: tt.identName},
				Sel: &ast.Ident{Name: tt.level},
			}

			got, ok := detectSlog(selector, tt.level)

			if ok != tt.wantOk {
				t.Errorf("detectSlog() ok = %v, want %v", ok, tt.wantOk)
				return
			}

			if ok && got.LoggerType != tt.wantType {
				t.Errorf("LoggerType = %v, want %v", got.LoggerType, tt.wantType)
			}
		})
	}
}

func TestDetect_Main(t *testing.T) {
	detector := New()

	tests := []struct {
		name     string
		setup    func() (*ast.CallExpr, *types.Info)
		wantType string
		wantOk   bool
	}{
		{
			name: "slog вызов",
			setup: func() (*ast.CallExpr, *types.Info) {
				call := &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "slog"},
						Sel: &ast.Ident{Name: "Info"},
					},
				}
				return call, &types.Info{}
			},
			wantType: "slog",
			wantOk:   true,
		},
		{
			name: "не валидный вызов (не selector)",
			setup: func() (*ast.CallExpr, *types.Info) {
				call := &ast.CallExpr{
					Fun: &ast.Ident{Name: "println"},
				}
				return call, &types.Info{}
			},
			wantOk: false,
		},
		{
			name: "nil аргументы",
			setup: func() (*ast.CallExpr, *types.Info) {
				return nil, nil
			},
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call, info := tt.setup()

			got, ok := detector.Detect(call, info)

			if ok != tt.wantOk {
				t.Errorf("Detect() ok = %v, want %v", ok, tt.wantOk)
			}

			if ok && got.LoggerType != tt.wantType {
				t.Errorf("LoggerType = %v, want %v", got.LoggerType, tt.wantType)
			}
		})
	}
}

func TestGetSelector(t *testing.T) {
	tests := []struct {
		name      string
		call      *ast.CallExpr
		wantOk    bool
		wantSelFn func(*ast.SelectorExpr) bool
	}{
		{
			name: "валидный selector",
			call: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "slog"},
					Sel: &ast.Ident{Name: "Info"},
				},
			},
			wantOk: true,
		},
		{
			name: "не selector (простой идентификатор)",
			call: &ast.CallExpr{
				Fun: &ast.Ident{Name: "println"},
			},
			wantOk: false,
		},
		{
			name:   "nil call",
			call:   nil,
			wantOk: false,
		},
		{
			name: "nil Fun",
			call: &ast.CallExpr{
				Fun: nil,
			},
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := getSelector(tt.call)

			if ok != tt.wantOk {
				t.Errorf("getSelector() ok = %v, want %v", ok, tt.wantOk)
			}

			if ok && got == nil {
				t.Error("getSelector() вернул nil при ok = true")
			}
		})
	}
}
