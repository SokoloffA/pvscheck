package main

import "fmt"

type SuppressArgs struct {
}

func (args *SuppressArgs) Execute(_ []string) error {
	return fmt.Errorf("NOT IMPLEMENTED")
}
