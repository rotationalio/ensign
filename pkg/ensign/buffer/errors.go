package buffer

import "errors"

var (
	ErrBufferFull  = errors.New("cannot write event: buffer is full")
	ErrBufferEmpty = errors.New("cannot read event: buffer is empty")
)
