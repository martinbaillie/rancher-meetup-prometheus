//go:generate go-bindata -o cow.go cow.ascii
package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	cowBytes, _ = Asset("cow.ascii")
	cow         = string(cowBytes)
	proverbs    = []string{"",
		"Don't communicate by sharing memory, share memory by communicating.",
		"Concurrency is not parallelism.",
		"Channels orchestrate; mutexes serialize.",
		"The bigger the interface, the weaker the abstraction.",
		"Make the zero value useful.",
		"interface{} says nothing.",
		"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
		"A little copying is better than a little dependency.",
		"Syscall must always be guarded with build tags.",
		"Cgo must always be guarded with build tags.",
		"Cgo is not Go.",
		"With the unsafe package there are no guarantees.",
		"Clear is better than clever.",
		"Reflection is never clear.",
		"Errors are values.",
		"Don't just check errors, handle them gracefully.",
		"Design the architecture, name the components, document the details.",
		"Documentation is for users.",
		"Don't panic.",
	}
)

// Service provides Go proverb related operations.
type Service interface {
	Textsay(int) (int, string)
	Cowsay(int) (int, string)
}

type service struct{}

func (service) Textsay(index int) (int, string) {
	return proverb(index)
}

func (service) Cowsay(index int) (int, string) {
	retindex, say := proverb(index)
	line := strings.Repeat("-", utf8.RuneCountInString(say)+4)
	return retindex, fmt.Sprintf("%s\n| %s |\n%s\n%s", line, say, line, cow)
}

func proverb(index int) (int, string) {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(2600)) * time.Millisecond)

	for index == 0 {
		index = rand.Intn(len(proverbs))
	}
	return index, proverbs[index]
}

// ServiceMiddleware is a chainable behavior modifier for Service.
type ServiceMiddleware func(Service) Service
