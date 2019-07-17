package planningunit

import "fmt"

type Id uint64

func Float64ToId(value float64) Id {
	return Id(value)
}

func (i Id) String() string {
	return fmt.Sprintf("%d", i)
}

type Ids []Id
