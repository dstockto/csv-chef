package recipe

import (
	"errors"
	"fmt"
	"github.com/carmo-evan/strtotime"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Processor func(string) string

var Now = time.Now

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

func TrimZeros(input string) (string, error) {
	f, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return "", fmt.Errorf("trimZeros: input is not numeric: got '%s'", input)
	}
	return strconv.FormatFloat(f, 'f', -1, 64), nil
}

func Repeat(count string, input string) (string, error) {
	num, err := strconv.Atoi(count)
	if err != nil {
		return "", fmt.Errorf("first arg is not an integer: got '%s'", count)
	}
	if num < 0 {
		return "", fmt.Errorf("first arg is negative: got '%d'", num)
	}

	return strings.Repeat(input, num), nil
}

func ReplaceString(search string, replace string, input string) (string, error) {
	return strings.ReplaceAll(input, search, replace), nil
}

func Today(now func() time.Time) (string, error) {
	return now().Format("2006-01-02"), nil
}

func FormatDate(format string, normalDate string) (string, error) {
	timestamp, err := time.Parse(time.RFC3339, normalDate)
	if err != nil {
		return normalDate, nil
	}
	return timestamp.Format(format), nil
}

func FormatDateF(format string, normalDate string) (string, error) {
	timestamp, err := time.Parse(time.RFC3339, normalDate)
	if err != nil {
		return "", fmt.Errorf("expected RFC3339 format for input date: '%s'", normalDate)
	}
	return timestamp.Format(format), nil
}

func ReadDate(format string, input string) (string, error) {
	timestamp, err := time.Parse(format, input)
	if err != nil {
		return input, nil
	}
	return timestamp.Format(time.RFC3339), nil
}

func ReadDateF(format string, input string) (string, error) {
	timestamp, err := time.Parse(format, input)
	if err != nil {
		return "", fmt.Errorf("unrecognized date '%s' for format: '%s'", input, format)
	}
	return timestamp.Format(time.RFC3339), nil
}

func SmartDate(date string) (string, error) {
	if _, err := time.Parse(time.RFC3339, date); err == nil {
		return date, nil
	}

	tz := time.UTC
	d, err := strtotime.Parse(date, 0)
	if err != nil {
		return "", err
	}

	pt := time.Unix(d, 0).In(tz)

	return pt.Format(time.RFC3339), nil
}

func IsPast(past string, future string, date string) (string, error) {
	if date == "" {
		return "", nil
	}
	normalizedDate, err := SmartDate(date)
	if err != nil {
		return "", fmt.Errorf("unable to recognize date: %v", err)
	}
	actualDate, err := time.Parse(time.RFC3339, normalizedDate)
	if err != nil {
		return "", fmt.Errorf("unable to parse date: %v", err)
	}
	now := Now()
	if now.After(actualDate) {
		return past, nil
	}
	return future, nil
}

func IsFuture(future string, past string, date string) (string, error) {
	normalizedDate, err := SmartDate(date)
	if err != nil {
		return "", fmt.Errorf("unable to recognize date: %v", err)
	}
	actualDate, err := time.Parse(time.RFC3339, normalizedDate)
	if err != nil {
		return "", fmt.Errorf("unable to parse date: %v", err)
	}
	now := Now()
	if now.Before(actualDate) {
		return future, nil
	}
	return past, nil
}

func Power(number string, power string) (string, error) {
	num, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return "", fmt.Errorf("unrecognized number '%s' for num parameter", number)
	}
	pow, err := strconv.ParseFloat(power, 64)
	if err != nil {
		return "", fmt.Errorf("unrecognized number '%s' for power parameter", power)
	}
	result := math.Pow(num, pow)
	return fmt.Sprintf("%f", result), nil
}

func Age(dob string) (string, error) {
	normalizedDate, err := SmartDate(dob)
	if err != nil {
		return "", err
	}
	// Ignoring the error because SmartDate would have already failed on a bad date
	birthdate, _ := time.Parse(time.RFC3339, normalizedDate)
	now := Now()
	years := now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		years--
	}

	return fmt.Sprintf("%d", years), nil
}

func NowTime(now func() time.Time) (string, error) {
	return now().Format(time.RFC3339), nil
}

func FirstChars(count string, input string) (string, error) {
	num, err := strconv.Atoi(count)
	if err != nil {
		return "", fmt.Errorf("first arg is not an integer: got '%s'", count)
	}
	if num < 0 {
		return "", fmt.Errorf("first arg is negative: got '%s'", count)
	}

	r := []rune(input)

	if num > len(r) {
		return input, nil
	}

	return string(r[:num]), nil
}

func LastChars(count string, input string) (string, error) {
	num, err := strconv.Atoi(count)
	if err != nil {
		return "", fmt.Errorf("first arg is not an integer: got '%s'", count)
	}
	if num < 0 {
		return "", fmt.Errorf("first arg is negative: got '%s'", count)
	}

	r := []rune(input)

	runeCount := len(r)
	if num > runeCount {
		return input, nil
	}

	return string(r[runeCount-num:]), nil
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
	if strings.EqualFold(input, from) {
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

func Coalesce(a string, b string) (string, error) {
	if a != "" {
		return a, nil
	}
	return b, nil
}

func Nth(delimiter string, index string, input string) (string, error) {
	if delimiter == "" {
		return "", errors.New("delimiter must not be empty")
	}
	idx, err := strconv.Atoi(index)
	if err != nil {
		return "", fmt.Errorf("index is not an integer: got '%s'", index)
	}
	if idx < 1 {
		return "", fmt.Errorf("index must be >= 1: got '%d'", idx)
	}

	fields := strings.Split(input, delimiter)
	if idx > len(fields) {
		return "", nil
	}
	return fields[idx-1], nil
}

func PadLeft(width string, pad string, input string) (string, error) {
	w, err := strconv.Atoi(width)
	if err != nil {
		return "", fmt.Errorf("width is not an integer: got '%s'", width)
	}
	if w < 0 {
		return "", fmt.Errorf("width must be non-negative: got '%d'", w)
	}
	if pad == "" {
		return "", errors.New("pad must not be empty")
	}

	r := []rune(input)
	if len(r) >= w {
		return input, nil
	}

	padRunes := []rune(pad)
	var prefix []rune
	for len(prefix)+len(r) < w {
		prefix = append(prefix, padRunes...)
	}
	combined := append(prefix, r...)
	// Trim from the left so the result is exactly w runes
	return string(combined[len(combined)-w:]), nil
}

func PadRight(width string, pad string, input string) (string, error) {
	w, err := strconv.Atoi(width)
	if err != nil {
		return "", fmt.Errorf("width is not an integer: got '%s'", width)
	}
	if w < 0 {
		return "", fmt.Errorf("width must be non-negative: got '%d'", w)
	}
	if pad == "" {
		return "", errors.New("pad must not be empty")
	}

	r := []rune(input)
	if len(r) >= w {
		return input, nil
	}

	padRunes := []rune(pad)
	combined := append([]rune{}, r...)
	for len(combined) < w {
		combined = append(combined, padRunes...)
	}
	// Trim the overflow from the right so the result is exactly w runes
	return string(combined[:w]), nil
}

func TitleCase(input string) (string, error) {
	words := strings.Fields(input)
	for i, word := range words {
		r := []rune(word)
		first := strings.ToUpper(string(r[0]))
		rest := strings.ToLower(string(r[1:]))
		words[i] = first + rest
	}
	return strings.Join(words, " "), nil
}

func RegexReplace(pattern string, replacement string, input string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regular expression '%s': %v", pattern, err)
	}
	return re.ReplaceAllString(input, replacement), nil
}

func Substring(start string, length string, input string) (string, error) {
	startNum, err := strconv.Atoi(start)
	if err != nil {
		return "", fmt.Errorf("start is not an integer: got '%s'", start)
	}
	if startNum < 1 {
		return "", fmt.Errorf("start must be >= 1: got '%d'", startNum)
	}
	lengthNum, err := strconv.Atoi(length)
	if err != nil {
		return "", fmt.Errorf("length is not an integer: got '%s'", length)
	}
	if lengthNum < 0 {
		return "", fmt.Errorf("length must be non-negative: got '%d'", lengthNum)
	}

	r := []rune(input)
	startIdx := startNum - 1
	if startIdx >= len(r) {
		return "", nil
	}
	end := startIdx + lengthNum
	if end > len(r) {
		end = len(r)
	}
	return string(r[startIdx:end]), nil
}

// SanitizeField guards against CSV formula injection by prefixing a
// single quote to any value that begins with a character a spreadsheet
// may interpret as a formula.
func SanitizeField(s string) string {
	if s == "" {
		return s
	}
	switch s[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + s
	}
	return s
}
