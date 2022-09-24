package dectalkdapi

type MMError struct {
	Code MMResult
}

func (err *MMError) Error() string {
	return mmErrorMessages[err.Code]
}
