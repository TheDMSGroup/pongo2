package pongo2

// Version string
const Version = "4.0.2"

// Must panics, if a Template couldn't successfully parsed. This is how you
// would use it:
//
//	var baseTemplate = pongo2.Must(pongo2.FromFile("templates/base.html"))
func Must(tpl *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return tpl
}

// ExecuteTemplateValue is a convenience function to execute a template string and return the raw value
// instead of a string. This is useful when you want to evaluate expressions that return non-string types.
func ExecuteTemplateValue(templateStr string, context Context) (interface{}, error) {
	tpl, err := FromString(templateStr)
	if err != nil {
		return nil, err
	}
	return tpl.ExecuteValue(context)
}

// ExecuteTemplateRaw is a convenience function to execute a template string and return all raw values
// without string conversion. This is useful for templates with multiple expressions.
func ExecuteTemplateRaw(templateStr string, context Context) ([]interface{}, error) {
	tpl, err := FromString(templateStr)
	if err != nil {
		return nil, err
	}
	return tpl.ExecuteRaw(context)
}
