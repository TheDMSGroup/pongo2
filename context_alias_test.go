package pongo2_test

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/flosch/pongo2/v4"
)

func copyContext(c pongo2.Context) pongo2.Context {
	out := make(pongo2.Context, len(c))
	for k, v := range c {
		out[k] = v
	}
	return out
}

// Execute must leave the caller's context untouched, including when tags
// write variables of their own (set, for, with, cycle, firstof, macros).
func TestExecuteDoesNotMutateContext(t *testing.T) {
	set := pongo2.NewSet("no globals", pongo2.MustNewLocalFileSystemLoader(""))
	ctx := pongo2.Context{
		"name":  "jane",
		"items": []string{"a", "b"},
		"lead":  map[string]interface{}{"state": "CA"},
	}
	before := copyContext(ctx)

	tpl, err := set.FromString(`{% set greeting = "hi" %}{{ greeting }} {{ name }}` +
		`{% for i in items %}{{ forloop.Counter }}{{ i }}{% cycle "x" "y" %}{% endfor %}` +
		`{% with name="bob" %}{{ name }}{% endwith %}{{ name }}` +
		`{% firstof missing lead.state %}` +
		`{% macro m(v) %}<{{ v }}>{% endmacro %}{{ m(name) }}`)
	if err != nil {
		t.Fatal(err)
	}
	out, err := tpl.Execute(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expected := "hi jane1ax2bybobjaneCA<jane>"
	if out != expected {
		t.Errorf("expected \"%s\", got \"%s\"", expected, out)
	}
	if !reflect.DeepEqual(ctx, before) {
		t.Errorf("context was mutated: expected %v, got %v", before, ctx)
	}
}

// User context keys win over globals, and globals never leak into the
// caller's context.
func TestExecuteGlobalsAndContext(t *testing.T) {
	set := pongo2.NewSet("globals", pongo2.MustNewLocalFileSystemLoader(""))
	set.Globals["a"] = "global-a"
	set.Globals["b"] = "global-b"
	ctx := pongo2.Context{"a": "ctx-a"}

	tpl, err := set.FromString("{{ a }} {{ b }}")
	if err != nil {
		t.Fatal(err)
	}
	out, err := tpl.Execute(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if out != "ctx-a global-b" {
		t.Errorf("expected \"ctx-a global-b\", got \"%s\"", out)
	}
	if len(ctx) != 1 {
		t.Errorf("expected the caller's context to keep 1 key, got %v", ctx)
	}
}

// mastodon-dynamic-integrations matches on this message to drop invalid keys
// and render again, so the check and its wording must stay.
func TestExecuteRejectsInvalidIdentifier(t *testing.T) {
	set := pongo2.NewSet("invalid key", pongo2.MustNewLocalFileSystemLoader(""))
	tpl, err := set.FromString("{{ name }}")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tpl.Execute(pongo2.Context{"name": "jane", "first-name": "jane"})
	if err == nil || !strings.Contains(err.Error(), "is not a valid identifier") {
		t.Errorf("expected an invalid identifier error, got %v", err)
	}
}

func TestExecuteRejectsMacroClash(t *testing.T) {
	set := pongo2.NewSet("macro clash", pongo2.MustNewLocalFileSystemLoader(""))
	tpl, err := set.FromString("{% macro m() export %}x{% endmacro %}{{ m() }}")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tpl.Execute(pongo2.Context{"m": "value"})
	if err == nil || !strings.Contains(err.Error(), "clashes with macro") {
		t.Errorf("expected a macro clash error, got %v", err)
	}
}

// Many goroutines render against one shared context map. Run with -race:
// any write to the caller's map during Execute shows up as a data race.
func TestExecuteConcurrentSharedContext(t *testing.T) {
	set := pongo2.NewSet("concurrent", pongo2.MustNewLocalFileSystemLoader(""))
	ctx := pongo2.Context{"name": "jane", "items": []int{1, 2, 3}}
	tpl, err := set.FromString(`{% set x = name %}{% for i in items %}{{ x }}{{ i }}{% endfor %}`)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	errs := make(chan error, 16)
	for g := 0; g < 16; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				out, err := tpl.Execute(ctx)
				if err != nil {
					errs <- err
					return
				}
				if out != "jane1jane2jane3" {
					errs <- fmt.Errorf("expected \"jane1jane2jane3\", got \"%s\"", out)
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}

func largeContext(keys int) pongo2.Context {
	ctx := make(pongo2.Context, keys)
	for i := 0; i < keys; i++ {
		ctx[fmt.Sprintf("k%d", i)] = fmt.Sprintf("value %d", i)
	}
	ctx["items"] = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	return ctx
}

// One field of a large context, the shape of an integration rendering one
// mapped field from a whole lead.
func BenchmarkExecuteLargeContext(b *testing.B) {
	set := pongo2.NewSet("bench", pongo2.MustNewLocalFileSystemLoader(""))
	tpl, err := set.FromString("{{ k1 }}")
	if err != nil {
		b.Fatal(err)
	}
	for _, keys := range []int{10, 100, 1000} {
		ctx := largeContext(keys)
		b.Run(fmt.Sprintf("keys=%d", keys), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := tpl.Execute(ctx); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// A loop creates a child context per iteration, which this change leaves as is.
func BenchmarkExecuteLoop(b *testing.B) {
	set := pongo2.NewSet("bench", pongo2.MustNewLocalFileSystemLoader(""))
	tpl, err := set.FromString("{% for i in items %}{{ i }}{{ k1 }}{% endfor %}")
	if err != nil {
		b.Fatal(err)
	}
	ctx := largeContext(100)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := tpl.Execute(ctx); err != nil {
			b.Fatal(err)
		}
	}
}
