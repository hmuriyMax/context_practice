package main

import (
	"context"
	"fmt"
)

func main() {
	var ctx context.Context
	fmt.Printf("%s", ctx.Value("key"))
}
