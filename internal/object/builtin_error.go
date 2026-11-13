package object

type BuiltinTypeError struct {
	Message string
}

func (err BuiltinTypeError) Error() string {
	return err.Message
}
