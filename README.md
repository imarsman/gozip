# gozip
An implementation of the zip utility in Go as an exercise.

Zip was first implemented in 1989 by the PKZIP corporation. Its specification
was put in the public domain in the same year (see
[Wikipedia](https://en.wikipedia.org/wiki/ZIP_(file_format))) for further information.

Go has a good library for the zip protocol. This project provides a way to
explore the protocol, find any issues, and overall learn more about the things
that are involved, namely i/o and file handling.

## Usage

To build

`go build .`

To run tests

`go test -v .`

## Notes

Wherever a feature of the zip utility is implemented an attempt will be made to
impliment it to behave identically.

The argument parsing library used here does not deal with arguments such as -1,
-2, -, etc. It may be that an argument will need to have a different identifier to
work around this.
