package main

import (
	"flag"
	"io/ioutil"

	"PrincessMononoke/execute"
	"PrincessMononoke/parser"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "", "dir for input file")
	flag.Parse()

	byteData, err := ioutil.ReadFile(dir)
	if err != nil {
		panic(err)
	}
	funcCases, err := parser.ConvertToFuncCases(byteData)
	if err != nil {
		panic(err)
	}
	for _, v := range funcCases {
		err = execute.RunFuncCase(v)
		if err != nil {
			panic(err)
		}
	}

	// jsonBytes, err := json.Marshal(funcCases)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(jsonBytes))

}
