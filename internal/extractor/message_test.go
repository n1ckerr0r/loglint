package extractor

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestExtractMessageFromArg(t *testing.T) {
	tests := []struct {
		name   string
		arg    ast.Expr
		want   string
		wantOk bool
	}{
		{
			name:   "простая строка",
			arg:    &ast.BasicLit{Kind: token.STRING, Value: `"привет"`},
			want:   "привет",
			wantOk: true,
		},
		{
			name:   "не строка",
			arg:    &ast.BasicLit{Kind: token.INT, Value: "123"},
			want:   "",
			wantOk: false,
		},
		{
			name:   "пустой вызов",
			arg:    nil,
			want:   "",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call := &ast.CallExpr{Args: []ast.Expr{tt.arg}}
			got, ok := ExtractMessageFromArg(call, 0)

			if ok != tt.wantOk || got != tt.want {
				t.Errorf("got %q, %v; want %q, %v", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestExtractMessage(t *testing.T) {
	call := &ast.CallExpr{
		Args: []ast.Expr{
			&ast.BasicLit{Kind: token.STRING, Value: `"тест"`},
		},
	}

	got, ok := ExtractMessage(call)
	if !ok || got != "тест" {
		t.Errorf("got %q, %v; want %q, true", got, ok, "тест")
	}
}

func TestExtractMessageFromArg_BadIndex(t *testing.T) {
	call := &ast.CallExpr{Args: []ast.Expr{}}

	got, ok := ExtractMessageFromArg(call, 0)
	if ok || got != "" {
		t.Errorf("got %q, %v; want empty, false", got, ok)
	}
}

func TestExtractMessageFromArg_NilCall(t *testing.T) {
	got, ok := ExtractMessageFromArg(nil, 0)
	if ok || got != "" {
		t.Errorf("got %q, %v; want empty, false", got, ok)
	}
}
