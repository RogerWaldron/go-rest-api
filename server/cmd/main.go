package main

import "fmt"
func Run() error {
	return nil
}
func main(){
	err := Run()
	if err != nil {
		fmt.Println(err)
	}
}