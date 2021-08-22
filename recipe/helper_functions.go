package recipe

import (
	"log"
	"regexp"
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

func MassProcess(incoming []string, processor Processor) (out []string) {
	for _, s := range incoming {
		out = append(out, processor(s))
	}
	return
}
