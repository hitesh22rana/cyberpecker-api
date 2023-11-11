package main

import (
	"fmt"

	cybernews "github.com/hitesh22rana/cyberpecker-api/pkg"
)

func main() {
	data, err := cybernews.GetNews("general")

	if err != nil {
		panic(err)
	}

	fmt.Println(len(data))
}
