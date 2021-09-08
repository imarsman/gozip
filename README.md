# gozip
An implementation of the zip utility in Go as an exercise.

Zip was first implemented in 1989 by the PKZIP corporation. Its specification
was put in the public domain in the same year (see
[Wikipedia](https://en.wikipedia.org/wiki/ZIP_(file_format))) for further information.

Go has a good library for the zip protocol. This project provides a way to
explore the protocol, find any issues, and overall learn more about the things
that are involved, namely i/o and file handling. I also find it generally useful
to learn how to pay close attention to exactly how things work.

## Usage

To build

`go build .`

To run tests

`go test -v .`

## Notes

Creating new archives and listing them works

```
compressed   uncompressed     date       time      name
----------------------------------------------------------------------
    7210         7210      2021-09-08  17:06:39  sample/1.txt
    4621         4621      2021-09-08  17:06:39  sample/2.txt
    2178         2178      2021-09-08  17:06:39  sample/3.txt
    1827         1827      2021-09-08  17:06:39  sample/4.txt
     918          918      2021-09-08  17:06:39  sample/5.txt
    7210         7210      2021-09-07  23:26:00  sample/orig/1.txt
    4621         4621      2021-09-07  23:26:00  sample/orig/2.txt
    2178         2178      2021-09-07  23:26:00  sample/orig/3.txt
    1827         1827      2021-09-07  23:26:00  sample/orig/4.txt
     918          918      2021-09-07  23:26:00  sample/orig/5.txt
----------------------------------------------------------------------
   33508        33508                            10
```

Wherever a feature of the zip utility is implemented an attempt will be made to
impliment it to behave identically.

The argument parsing library used here does not deal with arguments such as -1,
-2, -, etc. It may be that an argument will need to have a different identifier to
work around this.
