package myerrors

type CurError struct {
	ErrMsg string
}

// Текущая ошибка, которая выводится на страницу клиента
var Cur_error = CurError{ErrMsg: ""}
