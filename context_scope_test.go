package pongo2_test

import (
	"strings"
	"testing"

	"github.com/flosch/pongo2/v4"
)

// Context-key validation is on by default and rejects invalid identifiers.
func TestContextValidationRejectsInvalidKeyByDefault(t *testing.T) {
	tpl, err := pongo2.FromString(`{{ x }}`)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tpl.Execute(pongo2.Context{"x": 2, "foo-bar": 1})
	if err == nil {
		t.Fatal("expected error for invalid context key 'foo-bar', got nil")
	}
	if !strings.Contains(err.Error(), "not a valid identifier") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// With SkipContextValidation, the scan is skipped: invalid keys no longer error
// and valid keys still render.
func TestSkipContextValidation(t *testing.T) {
	tpl, err := pongo2.FromString(`{{ x }}`)
	if err != nil {
		t.Fatal(err)
	}
	tpl.Options.SkipContextValidation = true

	out, err := tpl.Execute(pongo2.Context{"x": 2, "foo-bar": 1})
	if err != nil {
		t.Fatalf("expected no error with validation skipped, got: %v", err)
	}
	if out != "2" {
		t.Fatalf("expected '2', got %q", out)
	}
}

// The scope chain layers child scopes over parents without copying: nested for
// loops expose forloop.parentloop, and outer loop variables remain visible.
func TestScopeChainNestedForLoop(t *testing.T) {
	tpl, err := pongo2.FromString(
		`{% for a in outer %}{% for b in inner %}{{ a }}{{ b }}({{ forloop.Parentloop.Counter }}) {% endfor %}{% endfor %}`)
	if err != nil {
		t.Fatal(err)
	}
	out, err := tpl.Execute(pongo2.Context{
		"outer": []string{"x", "y"},
		"inner": []string{"1", "2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	// outer a stays visible inside the inner loop; parentloop counter tracks outer.
	expected := "x1(1) x2(1) y1(2) y2(2) "
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}

// A {% with %} block shadows an outer variable within its scope; the outer
// binding is unchanged outside (parent scope never mutated by a child).
func TestScopeWithShadowing(t *testing.T) {
	tpl, err := pongo2.FromString(`{{ v }}{% with v="inner" %}{{ v }}{% endwith %}{{ v }}`)
	if err != nil {
		t.Fatal(err)
	}
	out, err := tpl.Execute(pongo2.Context{"v": "outer"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "outerinnerouter" {
		t.Fatalf("expected 'outerinnerouter', got %q", out)
	}
}

// The user context is read-only and reachable via the Public view; globals fall
// through. A {% set %} writes to the private scope, shadowing public.
func TestScopePrivateShadowsPublic(t *testing.T) {
	tpl, err := pongo2.FromString(`{{ v }}{% set v="private" %}{{ v }}`)
	if err != nil {
		t.Fatal(err)
	}
	out, err := tpl.Execute(pongo2.Context{"v": "public"})
	if err != nil {
		t.Fatal(err)
	}
	if out != "publicprivate" {
		t.Fatalf("expected 'publicprivate', got %q", out)
	}
}
