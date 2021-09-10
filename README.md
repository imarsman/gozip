# gozip
An implementation of the zip utility in Go as an exercise.

Zip was first implemented in 1989 by the PKZIP corporation. Its specification
was put in the public domain in the same year (see
[Wikipedia](https://en.wikipedia.org/wiki/ZIP_(file_format))) for further information.

Go has a good library for the zip protocol. This project provides a way to
explore the protocol, find any issues, and overall learn more about the things
that are involved, namely i/o and file handling. I also find it generally useful
to learn how to pay close attention to exactly how things work.

There are a huge number of possible options for the zip utility. I will likely
not implement most of them. As far as I can tell the Go zip library does not
work with updating an existing archive. In testing extra bytes ended up being
added to the archive with each update. For this reason the code first gets an
inventory of all files in an archive along with things like modified time and
then uses that to create a new archive after deleting the old one. This is not
ideal and if possible a switch will be made to updating rather than replacing.

Zip has enough arguments without also trying to implement unzip. I'll save that
for another project.

## Arguments

* `gozip -h` - print usage information
* `gozip -l <zipfile>` - list contents of zipfile
* `gozip <zipfile> <file>...` - default to freshen
* `gozip -a <zipfile> <file>...` - add and update files
* `gozip -u <zipfile> <file>...` - update if newer and add
* `gozip -f <zipfile> <file>...` - only update newer files already in archive

## Usage

To build

`go build .`

To run tests

`go test -v .`

## Simple tests

It's not as fast as the native zip.

```
time for i in {1..1000}; do ./gozip test/archive2.zip ./sample/; done

real	0m6.936s
user	0m3.693s
sys	0m2.772s

time for i in {1..1000}; do zip add -q test/archive2.zip ./sample/; done

real	0m4.237s
user	0m2.072s
sys	0m1.760s
```

## Notes

Creating new archives and listing them works

```
  compressed uncompressed      date       time        name
---------------------------------------------------------------------------
        3495        7210   2021-09-08  19:02:28  sample/1.txt
        2330        4621   2021-09-08  19:02:28  sample/2.txt
        1174        2178   2021-09-08  19:02:28  sample/3.txt
        1021        1827   2021-09-08  19:02:28  sample/4.txt
         497         918   2021-09-08  19:02:28  sample/5.txt
        3495        7210   2021-09-07  23:26:00  sample/orig/1.txt
        2330        4621   2021-09-07  23:26:00  sample/orig/2.txt
        1174        2178   2021-09-07  23:26:00  sample/orig/3.txt
        1021        1827   2021-09-07  23:26:00  sample/orig/4.txt
         497         918   2021-09-07  23:26:00  sample/orig/5.txt
---------------------------------------------------------------------------
       17034       33508                         10
```

Wherever a feature of the zip utility is implemented an attempt will be made to
impliment it to behave identically.
