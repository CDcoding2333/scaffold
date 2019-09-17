package flagtools

import "fmt"
import "strconv"

type IntList []int32

func (sl *IntList) String() string {
	return fmt.Sprint(*sl)
}

func (sl *IntList) Set(value string) error {
	i, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}
	*sl = append(*sl, int32(i))
	return nil
}

func (*IntList) Type() string {
	return "intList"
}
