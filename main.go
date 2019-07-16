package main

import (
	"flag"
	"io/ioutil"
	"log"

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
		executer := new(execute.Executer)
		err = executer.RunFuncCase(v)
		if err != nil {
			panic(err)
		}
		log.Printf("\n FuncName: %s\n Total: %d\n Success: %d\n Failed: %d\n Passed: %.2f%s\n", v.Name, len(v.Cases), executer.Success, executer.Failed, float64(executer.Success)/float64(len(v.Cases))*100, "%")
	}

	// jsonBytes, err := json.Marshal(funcCases)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(string(jsonBytes))

}
