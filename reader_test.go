package simple

import (
	"fmt"
	"io"
	"strings"
)

//basic is a simple io.Reader.
//It's just a byte slice that consumes itself as it's read.
//This means it can't store an error for the next read and has
//to return data and an error.
type basic []byte

//NewBasic creates a basic reader from the string s.
func NewBasic(s string) basic {
	return basic([]byte(s))
}

//Read implements the contract for io.Reader exactly,
//but it takes the option to return an error even if n > 0.
func (b *basic) Read(p []byte) (n int, err error) {
	//handle exceptional case
	if b == nil {
		return 0, io.ErrUnexpectedEOF
	}

	//once EOF, always EOF
	if len(*b) == 0 {
		return 0, io.EOF
	}

	//return 0, nil if nothing to read into
	if len(p) == 0 {
		return 0, nil
	}

	n = copy(p, *b)
	if len(p) <= len(*b) {
		//we fit in p, just slice away what was read
		*b = (*b)[len(p):]
	} else {
		//we have exhausted the input without filling p, so signal EOF as well
		*b = (*b)[:0]
		err = io.EOF
	}

	return n, err
}

//Using our simplistic io.Reader above, we see the care required to deal
//with a general io.Reader.
func ExampleReader() {
	//Basic is a simple reader that can return n > 0 and err != nil
	b := NewBasic("Hello, World!")
	p := make([]byte, 10) // note that len(p) < len(b)

	//It's easy to lose data if you expect a better behaved io.Reader
	for {
		//This is how we WANT to code our read loop
		n, err := b.Read(p)
		if err != nil {
			//but b returned an error with n > 0, data loss
			fmt.Printf("whoops, lost: %s\n", p[:n])
			break
		}
		fmt.Printf("%s\n", p[:n])
	}

	//Let's try that again, and handle all the cases this time
	b = NewBasic("Hello, World!")
	for {
		//This is how we HAVE to code our read loop
		//unless we happen to know we're using an io.Reader
		//that saves errors for subsequent reads, like *bufio.Reader.
		n, err := b.Read(p)
		//since we have to consider both cases simultaneously, we can't
		//be sure anything was read
		if n > 0 {
			fmt.Printf("%s\n", p[:n])
		}
		//need to handle the error after the input
		if err != nil {
			break
		}
	}

	//By wrapping b, we can write the loop we want without having to
	//have any special knowledge of the io.Reader being used.

	b = NewBasic("Hello, World!")
	r := NewReader(&b)
	for {
		n, err := r.Read(p)
		if err != nil {
			break
		}
		fmt.Printf("%s\n", p[:n])
	}

	// Output:
	// Hello, Wor
	// whoops, lost: ld!
	// Hello, Wor
	// ld!
	// Hello, Wor
	// ld!
}

func ExampleRead() {
	//Basic is a simple reader that can return n > 0 and err != nil
	b := NewBasic("Hello, World!")
	p := make([]byte, 10)
	r := NewReader(&b) //we use Reader here so we can handle data/errors disjointly
	for {
		//The Read helper takes care of indices, no matter how little is read.
		p, err := Read(r, p)
		if err != nil {
			break
		}
		fmt.Printf("%s\n", p)
	}
	// Output:
	// Hello, Wor
	// ld!
}

func ExampleReader_Err() {
	sr := strings.NewReader("Hello, World!")
	r := NewReader(sr)
	p := make([]byte, 10)

	p, err := Read(r, p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", p)

	if err = r.Err(); err != nil {
		fmt.Println(err)
		return
	}

	//we had to call Err above, before invoking a method on the underlying
	//io.Reader, sr, in case the previous read
	//resulted an error or it will, at best, be lost
	//and, at worst, be reported at the wrong time.
	if _, err = sr.Seek(-3, 1); err != nil {
		fmt.Println(err)
		return
	}

	//nothing special has
	p, err = Read(r, p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", p)

	// Output:
	// Hello, Wor
	// World!
}
