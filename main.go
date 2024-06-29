package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
)

var availableEncodings = []string{"UTF-8", "UTF-16", "UTF-32", "UTF-16LE", "UTF-32LE", "UTF-32BE"}

func checkEncodings(encodings ...string) error {
	for _, e := range encodings {
		if !slices.Contains(availableEncodings, e) {
			return fmt.Errorf("unknown encoding %s", e)
		}
	}
	return nil
}

func parseArguments(args []string) (*Converter, error) {
	if len(args) != 6 {
		return nil, errors.New("wrong number of arguments")
	}

	if args[1] == "-f" && args[3] == "-t" {
		if err := checkEncodings(args[2], args[4]); err != nil {
			return nil, err
		}
		return &Converter{
			FromEncoding:  args[2],
			ToEncoding:    args[4],
			InputFilePath: args[5],
		}, nil
	} else if args[1] == "-t" && args[3] == "-f" {
		if err := checkEncodings(args[2], args[4]); err != nil {
			return nil, err
		}
		return &Converter{
			FromEncoding:  args[4],
			ToEncoding:    args[2],
			InputFilePath: args[5],
		}, nil
	} else {
		return nil, errors.New("wrong arguments")
	}
}

func main() {
	converter, err := parseArguments(os.Args)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("usage: unicode-conv -f <input_encoding> -t <output_encoding> <input_file>\n")
		fmt.Printf("accepted encodings: %v\n", availableEncodings)
		os.Exit(1)
	}
	bytes, convErr := converter.Convert()
	if convErr != nil {
		fmt.Printf("conversion error %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s", string(bytes))
}
