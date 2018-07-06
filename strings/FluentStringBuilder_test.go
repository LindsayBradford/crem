// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package strings

import "fmt"

func ExampleFluentBuilder_Add() {
	PetsAndFoods := new(FluentBuilder).Add("cat", "-dog-", "canary, ").Add("pie", "-cake-", "tart").String()
	fmt.Print(PetsAndFoods)
	// Output: cat-dog-canary, pie-cake-tart
}
