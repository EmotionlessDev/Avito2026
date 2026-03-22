package validator

import "regexp"

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) In(value string, list ...string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}

func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (v *Validator) Unique(values []string) bool {
	uniqueValues := make(map[string]struct{})
	for _, value := range values {
		uniqueValues[value] = struct{}{}
	}
	return len(uniqueValues) == len(values)
}
