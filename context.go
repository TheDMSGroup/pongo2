package pongo2

import (
	"errors"
	"fmt"
	"regexp"
)

var reIdentifiers = regexp.MustCompile("^[a-zA-Z0-9_]+$")

var autoescape = true

func SetAutoescape(newValue bool) {
	autoescape = newValue
}

// A Context type provides constants, variables, instances or functions to a template.
//
// pongo2 automatically provides meta-information or functions through the "pongo2"-key.
// Currently, context["pongo2"] contains the following keys:
//  1. version: returns the version string
//
// Template examples for accessing items from your context:
//
//	{{ myconstant }}
//	{{ myfunc("test", 42) }}
//	{{ user.name }}
//	{{ pongo2.version }}
type Context map[string]interface{}

type context interface {
	GetValue(key string) (interface{}, bool)
	SetValue(string, interface{})
	GetIdentifiers() []string
}

type mapContext map[string]interface{}

func (m mapContext) GetValue(key string) (interface{}, bool) {
	v, ok := m[key]
	return v, ok
}

func (m mapContext) SetValue(key string, val interface{}) {
	m[key] = val
}

func (m mapContext) GetIdentifiers() []string {
	identifiers := make([]string, 0, len(m))
	for i := range m {
		identifiers = append(identifiers, i)
	}
	return identifiers
}

func checkForValidIdentifiers(identifiers []string) *Error {
	for _, v := range identifiers {
		if !reIdentifiers.MatchString(v) {
			return &Error{
				Sender:    "checkForValidIdentifiers",
				OrigError: fmt.Errorf("context-key '%s' is not a valid identifier", v),
			}
		}
	}
	return nil
}

type mergedContext struct {
	a, b context
}

func (m mergedContext) GetValue(key string) (interface{}, bool) {
	val, had := m.b.GetValue(key)
	if !had {
		val, had = m.a.GetValue(key)
	}
	return val, had
}

func (m mergedContext) SetValue(key string, val interface{}) {
	m.b.SetValue(key, val)
}

func (m mergedContext) GetIdentifiers() []string {
	identifiers := make([]string, 0)
	if len(m.a.GetIdentifiers()) > 0 {
		identifiers = append(identifiers, m.a.GetIdentifiers()...)
	}
	if len(m.b.GetIdentifiers()) > 0 {
		identifiers = append(identifiers, m.b.GetIdentifiers()...)
	}
	return identifiers
}

func mergeContexts(a, b context) context {
	return mergedContext{a: a, b: b}
}

// Update updates this context with the key/value-pairs from another context.
// / todo solve by wrapping
func (c Context) Update(other Context) Context {
	for k, v := range other {
		c[k] = v
	}
	return c
}

// ExecutionContext contains all data important for the current rendering state.
//
// If you're writing a custom tag, your tag's Execute()-function will
// have access to the ExecutionContext. This struct stores anything
// about the current rendering process's Context including
// the Context provided by the user (field Public).
// You can safely use the Private context to provide data to the user's
// template (like a 'forloop'-information). The Shared-context is used
// to share data between tags. All ExecutionContexts share this context.
//
// Please be careful when accessing the Public data.
// PLEASE DO NOT MODIFY THE PUBLIC CONTEXT (read-only).
//
// To create your own execution context within tags, use the
// NewChildExecutionContext(parent) function.
type ExecutionContext struct {
	template   *Template
	macroDepth int

	Autoescape bool
	Public     context
	Private    context
	Shared     context
}

var pongo2MetaContext = mapContext{
	"version": Version,
}

func newExecutionContext(tpl *Template, ctx context) *ExecutionContext {
	// Make the pongo2-related funcs/vars available to the context
	privateCtx := mapContext{
		"pongo2": pongo2MetaContext,
	}

	return &ExecutionContext{
		template: tpl,

		Public:     mergeContexts(ctx, mapContext{}),
		Private:    privateCtx,
		Autoescape: autoescape,
	}
}

func NewChildExecutionContext(parent *ExecutionContext) *ExecutionContext {
	newctx := &ExecutionContext{
		template: parent.template,

		Public:     parent.Public,
		Private:    mergeContexts(parent.Private, mapContext{}),
		Autoescape: parent.Autoescape,
	}
	newctx.Shared = parent.Shared

	return newctx
}

func (ctx *ExecutionContext) Error(msg string, token *Token) *Error {
	return ctx.OrigError(errors.New(msg), token)
}

func (ctx *ExecutionContext) OrigError(err error, token *Token) *Error {
	filename := ctx.template.name
	var line, col int
	if token != nil {
		// No tokens available
		// TODO: Add location (from where?)
		filename = token.Filename
		line = token.Line
		col = token.Col
	}
	return &Error{
		Template:  ctx.template,
		Filename:  filename,
		Line:      line,
		Column:    col,
		Token:     token,
		Sender:    "execution",
		OrigError: err,
	}
}

func (ctx *ExecutionContext) Logf(format string, args ...any) {
	ctx.template.set.logf(format, args...)
}
