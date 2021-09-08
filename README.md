# gozip
An implementation of the zip utility in Go as an exercise.

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
-2, etc. It may be that an argument will need to have a different identifier to
work around this.