package recipe

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Processor func(string) string

func Lowercase(s string) string {
	return strings.ToLower(s)
}

func Uppercase(s string) string {
	return strings.ToUpper(s)
}

func NoDigits(s string) string {
	reg, err := regexp.Compile("[0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s, "")
}

func OnlyDigits(s string) string {
	reg := regexp.MustCompile("[^0-9]+")
	return reg.ReplaceAllString(s, "")
}

func JoinFunc(p string) Processor {
	return func(s string) string {
		return p + s
	}
}

func Add(x string, y string) (string, error) {
	xnum, err := strconv.Atoi(x)
	if err != nil {
		return "", fmt.Errorf("first arg to Add was not an integer: %s", x)
	}
	ynum, err := strconv.Atoi(y)
	if err != nil {
		return "", fmt.Errorf("second arg to Add was not an integer: %s", y)
	}

	sum := xnum + ynum

	return strconv.Itoa(sum), nil
}

func AddFloat(x string, y string, decimals string) (string, error) {
	xnum, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return "", fmt.Errorf("first arg to AddFloat was not numeric: %s", x)
	}
	ynum, err := strconv.ParseFloat(y, 64)
	if err != nil {
		return "", fmt.Errorf("second arg to AddFloat was not numeric: %s", y)
	}
	precision, err := strconv.Atoi(decimals)
	if err != nil {
		return "", fmt.Errorf("AddFloat precision should be an integer for number of decimals, or -1 for all, found %s", decimals)
	}

	format := "%f"
	if precision != -1 {
		format = fmt.Sprintf("%%.%df", precision)
	}

	sum := xnum + ynum

	return fmt.Sprintf(format, sum), nil
}

func Change(from string, to string, input string) (string, error) {
	if input == from {
		return to, nil
	}
	return input, nil
}

func MassProcess(incoming []string, processor Processor) (out []string) {
	for _, s := range incoming {
		out = append(out, processor(s))
	}
	return
}
