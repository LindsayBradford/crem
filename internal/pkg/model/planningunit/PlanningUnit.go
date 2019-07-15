package planningunit

import "fmt"

type Id uint64

func (i Id) String() string {
	return fmt.Sprintf("%d", i)
}

type Ids []Id
