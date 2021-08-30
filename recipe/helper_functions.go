package recipe

import (
	"errors"
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

func JoinFunc(p string) Processor {
	return func(s string) string {
		return p + s
	}
}

func Add(x string, y string) (string, error) {
	xnum, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return "", fmt.Errorf("first arg to Add was not numeric: %s", x)
	}
	ynum, err := strconv.ParseFloat(y, 64)
	if err != nil {
		return "", fmt.Errorf("second arg to Add was not numeric: %s", y)
	}

	sum := xnum + ynum

	return fmt.Sprintf("%f", sum), nil
}

func Subtract(x string, y string) (string, error) {
	xnum, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return "", fmt.Errorf("first arg to subtract was not numeric: %s", x)
	}
	ynum, err := strconv.ParseFloat(y, 64)
	if err != nil {
		return "", fmt.Errorf("second arg to subtract was not numeric: %s", y)
	}

	difference := xnum - ynum
	return fmt.Sprintf("%f", difference), nil
}

func Multiply(x string, y string) (string, error) {
	xnum, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return "", fmt.Errorf("error: first arg to multiply was not numeric, got '%s'", x)
	}
	ynum, err := strconv.ParseFloat(y, 64)
	if err != nil {
		return "", fmt.Errorf("error: second arg to multiply was not numeric, got '%s'", y)
	}

	product := xnum * ynum
	return fmt.Sprintf("%f", product), nil
}

func Divide(x string, y string) (string, error) {
	xnum, err := strconv.ParseFloat(x, 64)
	if err != nil {
		return "", fmt.Errorf("error: first arg to divide was not numeric, got '%s'", x)
	}
	ynum, err := strconv.ParseFloat(y, 64)
	if err != nil {
		return "", fmt.Errorf("error: second arg to divide was not numeric, got '%s'", y)
	}

	if ynum == 0.0 {
		return "", errors.New("error: attempt to divide by zero")
	}

	result := xnum / ynum
	return fmt.Sprintf("%f", result), nil
}

func NumberFormat(digits string, input string) (string, error) {
	digitsNum, err := strconv.Atoi(digits)
	if err != nil {
		return "", fmt.Errorf("error: digits must be an integer, got '%s'", digits)
	}

	inputNum, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return "", fmt.Errorf("error: input is not numeric: got '%s'", input)
	}

	format := fmt.Sprintf("%%.%df", digitsNum)

	return fmt.Sprintf(format, inputNum), nil
}

func OnlyDigits(input string) (string, error) {
	reg := regexp.MustCompile("[^0-9]+")
	return reg.ReplaceAllString(input, ""), nil
}

func Modulus(x string, y string) (string, error) {
	xNum, err := strconv.Atoi(x)
	if err != nil {
		return "", fmt.Errorf("first arg to mod was not an integer: '%s'", x)
	}
	yNum, err := strconv.Atoi(y)
	if err != nil {
		return "", fmt.Errorf("second arg to mod was not an integer: '%s'", y)
	}
	if yNum == 0 {
		return "", errors.New("attempt to divide by zero")
	}

	remainder := xNum % yNum

	return fmt.Sprintf("%d", remainder), nil
}

func Trim(input string) (string, error) {
	return strings.TrimSpace(input), nil
}

func FirstChars(count string, input string) (string, error) {
	num, err := strconv.Atoi(count)
	if err != nil {
		return "", fmt.Errorf("first arg is not an integer: got '%s'", count)
	}
	if num < 1 {
		return "", fmt.Errorf("first arg is not a positive integer: got '%s'", count)
	}

	if num > len(input) {
		return input, nil
	}

	return input[:num], nil
}

func RemoveDigits(input string) (string, error) {
	reg := regexp.MustCompile("[0-9]+")
	return reg.ReplaceAllString(input, ""), nil
}

func Change(from string, to string, input string) (string, error) {
	if input == from {
		return to, nil
	}
	return input, nil
}

func ChangeI(from string, to string, input string) (string, error) {
	if strings.ToLower(input) == strings.ToLower(from) {
		return to, nil
	}
	return input, nil
}

func IfEmpty(emptyVal string, notEmptyVal string, input string) (string, error) {
	if input == "" {
		return emptyVal, nil
	}
	return notEmptyVal, nil
}

func MassProcess(incoming []string, processor Processor) (out []string) {
	for _, s := range incoming {
		out = append(out, processor(s))
	}
	return
}
