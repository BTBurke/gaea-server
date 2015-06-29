package error

type APIError struct {
    Code int
    Developer string
    User string
}

func (a APIError) Error() string { return a.User }