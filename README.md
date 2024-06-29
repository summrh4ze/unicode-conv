# unicode-conv

Command line tool for changing unicode encoding of files. Usage similar to iconv:

```go run . -f <from encoding> -t <to encoding> input_file > output_file```

Supported conversions between UTF-8, UTF-16 (LE/BE) and UTF-32 (LE/BE)

Tested using this [utf-8 demo file](https://www.w3.org/2001/06/utf-8-test/UTF-8-demo.html)