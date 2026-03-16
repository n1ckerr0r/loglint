package plugin

import (
	"testing"
)

func TestNew(t *testing.T) {
	p, err := New(nil)
	if err != nil {
		t.Errorf("New(nil) ошибка: %v", err)
	}
	if p == nil {
		t.Error("New(nil) вернул nil")
	}

	p, err = New(map[string]any{})
	if err != nil {
		t.Errorf("New(empty) ошибка: %v", err)
	}
	if p == nil {
		t.Error("New(empty) вернул nil")
	}

	cfg := map[string]any{
		"check-lowercase": true,
		"check-english":   true,
		"sensitive-keywords": []string{"pass"},
	}

	p, err = New(cfg)
	if err != nil {
		t.Errorf("New(with config) ошибка: %v", err)
	}
}

func TestBuildAnalyzers(t *testing.T) {
	p, _ := New(nil)

	a, err := p.BuildAnalyzers()

	if err != nil {
		t.Errorf("BuildAnalyzers() ошибка: %v", err)
	}
	if len(a) != 1 {
		t.Errorf("хотим 1 анализатор, получили %d", len(a))
	}
}
