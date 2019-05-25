package main

import "github.com/boothgames/nightfury/cmd"

//go:generate ./scripts/mocks

func main() {
	cmd.Execute()
}
