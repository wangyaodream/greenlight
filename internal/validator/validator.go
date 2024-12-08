package validator

import "regexp"

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(field, message string) {
	if _, ok := v.Errors[field]; !ok {
		v.Errors[field] = message
	}
}

func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

func PermiteedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func Matchs(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// 判断切片中是否有重复的元素
func Unique[T comparable](values []T) bool {
	encountered := make(map[T]bool)

	for _, value := range values {
		encountered[value] = true
	}

	return len(values) == len(encountered)
}
