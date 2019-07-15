package parser

import (
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
	t.Log(funcCases)
}
