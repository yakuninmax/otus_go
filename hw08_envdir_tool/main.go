package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	envVars, err := ReadDir("testdata/env")
	if err != nil {
		log.Panic(err)
	}

	exitCode := RunCmd(os.Args[1:], envVars)

	fmt.Println(exitCode)

	for k, m := range envVars {
		fmt.Println(k, "value is", m.Value, "and it need to remove", m.NeedRemove)
	}

}
