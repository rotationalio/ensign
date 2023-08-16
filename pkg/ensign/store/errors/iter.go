package errors

func NewIter(err error) ErrorIterator {
	return ErrorIterator{err: err}
}

// ErrorIterator implements the iterator.Iterator interface and can be used to return
// errors from methods that need to return iterators.
type ErrorIterator struct {
	err error
}

func (e ErrorIterator) Key() []byte   { return nil }
func (e ErrorIterator) Value() []byte { return nil }
func (e ErrorIterator) Next() bool    { return false }
func (e ErrorIterator) Prev() bool    { return false }
func (e ErrorIterator) Error() error  { return e.err }
func (e ErrorIterator) Release()      {}
