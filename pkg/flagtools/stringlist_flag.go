package flagtools

import (
	"fmt"
	"strings"
)

type StringList []string

func (sl *StringList) String() string {
	return fmt.Sprint(*sl)
}

func (sl *StringList) Set(value string) error {
	for _, s := range strings.Split(value, ";") {
		*sl = append(*sl, s)
	}
	return nil
}

func (*StringList) Type() string {
	return "stringList"
}
