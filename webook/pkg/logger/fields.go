package logger

func Strings(key, val string) Field {
	return Field{
		Key:   key,
		Value: val,
	}
}
