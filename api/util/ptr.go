package util

func Ptr[T comparable](v T) *T {
	return &v
}

func SafeValue[T comparable](v *T) T {
	if v == nil {
		temp := new(T)
		return *temp
	}
	return *v
}
