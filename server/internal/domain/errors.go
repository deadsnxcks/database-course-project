package domain

import (
	"errors"
)

type Kind int

const (
	KindUnknown Kind = iota
	KindNotFound
	KindConflict
	KindInUse
	KindUnprocessable
)

type Error struct {
	Kind Kind
	Msg  string
	err  error
}

func (k Kind) String() string {
	switch k {
	case KindNotFound:
		return "not_found"
	case KindConflict:
		return "conflict"
	case KindInUse:
		return "in_use"
	case KindUnprocessable:
		return "unprocessable"
	default:
		return "unknown"
	}
}

func (e *Error) Error() string { return e.Msg }

func (e *Error) Unwrap() error { return e.err }

func (e *Error) Wrap(err error) *Error {
	return &Error{
		Kind: e.Kind,
		Msg:  e.Msg,
		err:  err,
	}
}

func NotFound(msg string) *Error      { return &Error{Kind: KindNotFound, Msg: msg} }
func Conflict(msg string) *Error      { return &Error{Kind: KindConflict, Msg: msg} }
func InUse(msg string) *Error         { return &Error{Kind: KindInUse, Msg: msg} }
func Unprocessable(msg string) *Error { return &Error{Kind: KindUnprocessable, Msg: msg} }

func KindOf(err error) Kind {
	var de *Error
	if errors.As(err, &de) {
		return de.Kind
	}

	return KindUnknown
}
