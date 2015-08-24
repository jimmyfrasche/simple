//Package simple simplifies direct use of io.Readers.
//
//Generally, you should use *bufio.Reader or the helper functions in io/ioutil
//instead of this package.
//However, low level io code and library code may mean using an io.Reader directly.
//
//This package provides two helpers.
//
//The first is a struct that wraps an arbitrary io.Reader
//and guarantees a stronger contract
//
//The second a function that simplifies calling an arbitrary io.Reader.
//
//They play well together.
package simple

import "io"

//Reader wraps any io.Reader and strengthens the io.Reader contract
//by never returning an error when n > 0.
//
//The io.Reader contract allows reads to either return both data and an error
//or to return the error on the next read request in lieu of an actual read.
//This means to avoid data loss it must always be assumed that the reader
//can return an error on a read with n > 0.
//But it is much easier to code a read loop if you can assume
//that either data was read or an error was returned but not both.
//The *simple.Reader stores the error if n > 0 and returns it on the next call,
//allowing simpler code.
//
//Stored errors are discarded after being returned.
//
//Some care must be taken if used in conjunction with method calls
//on the wrapped io.Reader.
//See the Err method for more details.
type Reader struct {
	err error
	r   io.Reader
}

//NewReader wraps an io.Reader in a simple.Reader.
func NewReader(r io.Reader) *Reader {
	if r == nil {
		panic("cannot wrap nil io.Reader")
	}

	return &Reader{
		r: r,
	}
}

//Read wraps the underlying Read to ensure err == nil if n > 0.
//
//If the wrapped io.Reader returns an error when n > 0, it is stored until
//the next call to Read (or Err) which will then return and discard the error
//without making a Read.
//
//If you need to access a different method of the wrapped io.Reader
//after a successful read, then it is your responsibility to first call Err.
func (r *Reader) Read(p []byte) (n int, err error) {
	//if we had a previous error stored, return it and clear the store
	if r.err != nil {
		return 0, r.Err()
	}

	n, err = r.r.Read(p)

	//error and data returned, store error for next call
	if n != 0 && err != nil {
		r.err = err
		return n, nil
	}

	//otherwise just return
	return n, err
}

//Err returns, then discards, any error stored from the last Read.
//
//It is only necessary to check this if you make a successful read,
//then attempt to call a different method (for example, Seek)
//on the wrapped io.Reader.
//
//If Err is called twice in a row, with no intervening reads, the second call
//will always return nil.
func (r *Reader) Err() error {
	//we need to unstick the error in case it was transitory
	//or only applied to a particular read.
	var err error
	err, r.err = r.err, nil
	return err
}

//Read grows p to its capacity, calls r.Read with p,
//and slices p to contain only the returned data before returning it.
//
//This allows the buffer to be reused between iterations of a read loop,
//without having to fuss over indices.
//
//Combine with *Reader for best experience.
func Read(r io.Reader, p []byte) ([]byte, error) {
	p = p[:cap(p)]
	n, err := r.Read(p)
	if p != nil {
		p = p[:n]
	}
	return p, err
}
