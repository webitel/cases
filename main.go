package main

import "github.com/webitel/cases/cmd/main"

//go:generate go run github.com/bufbuild/buf/cmd/buf@v1.42.0 generate --template buf.gen.yaml

func main() {
	cmd.Run()
}
