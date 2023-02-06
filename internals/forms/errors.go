package forms

type errors map[string][]string

// Add adds an error message for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	// in case there are multiple errors for a given key, it returns the first one (which should be the most important and be checked at first) and then, as conditions gradually fulfilled, it returns the next, still unhandled error, and so on.
	return es[0]
}
