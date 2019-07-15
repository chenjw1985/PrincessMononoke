package parser

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestConvertToFuncCases(t *testing.T) {
	bs, err := ioutil.ReadFile("./usecase.xlsx")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	funcCases, err := ConvertToFuncCases(bs)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	jsonBytes, err := json.Marshal(funcCases)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Log(string(jsonBytes))
}
