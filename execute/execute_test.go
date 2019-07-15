package execute

import (
	"PrincessMononoke/models"
	"fmt"
	"strings"
	"testing"
)

func TestVerifyJSON(t *testing.T) {
	s, e := VerifyJSON(`$.data==2 && ($.items[0].name=="test" || $.items[1].value==true)`, `{"data":2,"items":[{"name":"test"},{"value":true}]}`)
	if e != nil {
		t.Fatal(e)
	}
	if s != true {
		t.Fatal(s)
	}
}

func TestVerifyTEXT(t *testing.T) {
	s, e := VerifyTEXT(`$.indexOf('"b"')>-1`, strings.Replace(`a"b"cd`, `"`, `\"`, -1))
	if e != nil {
		t.Fatal(e)
	}
	if s != true {
		t.Fatal(s)
	}
}

func TestGetValuesByBody(t *testing.T) {
	keys := GetKeysByString("http://{{$.items[0].name}}/{{$.items[1].value}}/?id={{$.data}}")
	if len(keys) < 3 {
		t.Fatalf("Number of errors %d, %v", len(keys), keys)
	}
	vals := GetValuesByBody(keys, `{"data":2,"items":[{"name":"test"},{"value":true}]}`)
	if len(vals) != len(keys) {
		t.Fatalf("Number of errors %d", len(vals))
	}
	t.Log(vals)
}

func TestMakeACase(t *testing.T) {
	lastStep := &models.CaseStep{
		Name:        "step 1",
		Level:       "P0",
		Method:      "GET",
		URL:         "http://www.liveramp.com",
		Data:        "",
		ResType:     "JSON",
		Verfication: "",

		Result: &models.CaseOutput{
			Status:  200,
			Body:    []byte(`{"data":2,"items":[{"name":"test"},{"value":true}]}`),
			Success: true,
			Error:   nil,
		},
	}

	nextSetp := &models.CaseStep{
		Name:        "step 2",
		Level:       "P0",
		Method:      "POST",
		URL:         "http://www.liveramp.com/{{$.items[0].name}}/{{$.items[1].value}}",
		Data:        `{"data":{{$.data}}}`,
		ResType:     "JSON",
		Verfication: "",

		Result: nil,
	}
	caseStep := MakeACase(lastStep, nextSetp)

	url := nextSetp.URL
	data := nextSetp.Data
	keys := make([]string, 0)
	keys = append(keys, GetKeysByString(url)...)
	keys = append(keys, GetKeysByString(data)...)

	vals := GetValuesByBody(keys, string(lastStep.Result.Body))
	for k, v := range vals {
		url = strings.Replace(url, fmt.Sprintf("{{%s}}", k), v, -1)
		data = strings.Replace(data, fmt.Sprintf("{{%s}}", k), v, -1)
	}

	if caseStep.URL != url {
		t.Fatalf("%s != %s", caseStep.URL, url)
	}
	if caseStep.Data != data {
		t.Fatalf("%s != %s", caseStep.Data, data)
	}
}
