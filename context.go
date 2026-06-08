package pongo2

import (
	"errors"
)

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

func isValidIdentifier(s string) bool {
	for i := range s {
		if !isValidIdentifierChar(s[i]) {
			return false
		}
	}
	return true
}

func isValidIdentifierChar(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c >= '0' && c <= '9') ||
		c == '_'
}

// Update updates this context with the key/value-pairs from another context.
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
// You can safely use the Private context (a Scope) to provide data to the
// user's template (like a 'forloop'-information). The Shared-context is used
// to share data between tags. All ExecutionContexts share this context.
//
// Public is a read-only view over the user-supplied context and the set's
// globals; the caller's map is used directly (never copied) and must not be
// mutated during rendering. Private is a non-copying scope chain: a child
// context layers its own variables over its parent.
//
// To create your own execution context within tags, use the
// NewChildExecutionContext(parent) function.
type ExecutionContext struct {
	template *Template

	Autoescape bool

	// Public is the read-only view of user data overlaying the set's globals.
	// Access via Public.Get/Has/Range, never by indexing.
	Public PublicContext

	// parent links a child context to the one it was derived from, forming the
	// scope chain walked by Private. nil for a root context.
	parent *ExecutionContext

	// privateVars holds this context's own private layer, allocated lazily on
	// first write. Parent layers are reached through parent, never copied.
	privateVars Context

	// Private is engine-managed scoped data (e.g. 'forloop', {% set %} vars,
	// macros). Access via Private.Get/Set/Delete/Range, never by indexing.
	Private Scope

	Shared Context
}

var pongo2MetaContext = Context{
	"version": Version,
}

func newExecutionContext(tpl *Template, ctx Context) *ExecutionContext {
	newctx := &ExecutionContext{
		template: tpl,

		// Make the pongo2-related funcs/vars available to the context.
		privateVars: Context{"pongo2": pongo2MetaContext},
		Public:      PublicContext{vars: ctx, globals: tpl.set.Globals},
		Shared:      make(Context),
		Autoescape:  autoescape,
	}
	newctx.Private = Scope{ec: newctx}
	return newctx
}

func NewChildExecutionContext(parent *ExecutionContext) *ExecutionContext {
	newctx := &ExecutionContext{
		template: parent.template,
		parent:   parent,

		Public:     parent.Public,
		Shared:     parent.Shared,
		Autoescape: parent.Autoescape,
	}
	newctx.Private = Scope{ec: newctx}
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

func (ctx *ExecutionContext) Logf(format string, args ...interface{}) {
	ctx.template.set.logf(format, args...)
}
