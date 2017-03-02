package looli

type Error struct {
	Err  error
	Code int
	Meta interface{}
}

func (err *Error) Error() string {
	return err.Err.Error()
}
