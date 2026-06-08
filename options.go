package pongo2

// Options allow you to change the behavior of template-engine.
// You can change the options before calling the Execute method.
type Options struct {
	// If this is set to true the first newline after a block is removed (block, not variable tag!). Defaults to false.
	TrimBlocks bool

	// If this is set to true leading spaces and tabs are stripped from the start of a line to a block. Defaults to false
	LStripBlocks bool

	// If this is set to true, the per-Execute scan that validates context keys
	// are valid identifiers (and don't clash with exported macro names) is
	// skipped. Enable only when context keys are trusted/known-valid. Keys that
	// aren't valid identifiers are unreferenceable by the template regardless,
	// so skipping the scan does not change render output. Defaults to false.
	SkipContextValidation bool
}

func newOptions() *Options {
	return &Options{
		TrimBlocks:            false,
		LStripBlocks:          false,
		SkipContextValidation: false,
	}
}

// Update updates this options from another options.
func (opt *Options) Update(other *Options) *Options {
	opt.TrimBlocks = other.TrimBlocks
	opt.LStripBlocks = other.LStripBlocks
	opt.SkipContextValidation = other.SkipContextValidation

	return opt
}
