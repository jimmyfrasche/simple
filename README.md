#simple [![GoDoc](https://godoc.org/github.com/jimmyfrasche/simple?status.png)](https://godoc.org/github.com/jimmyfrasche/simple)
Package simple simplifies direct use of io.Readers.

Download:
```shell
go get github.com/jimmyfrasche/simple
```

* * *
Package simple simplifies direct use of io.Readers.

Generally, you should use *bufio.Reader or the helper functions in io/ioutil
instead of this package.
However, low level io code and library code may mean using an io.Reader directly.

This package provides two helpers.

The first is a struct that wraps an arbitrary io.Reader
and guarantees a stronger contract

The second a function that simplifies calling an arbitrary io.Reader.

They play well together.



* * *
Automatically generated by [autoreadme](https://github.com/jimmyfrasche/autoreadme) on 2015.08.23