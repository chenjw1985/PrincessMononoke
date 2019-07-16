package execute

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"PrincessMononoke/models"

	"github.com/ddliu/go-httpclient"
	"github.com/robertkrimen/otto"
)

type Executer struct {
	Total   int
	Success int
	Failed  int
}

func (e *Executer) RunFuncCase(funcCase *models.FuncCase) error {
	for _, v := range funcCase.Cases {
		err := e.RunTestCase(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Executer) RunTestCase(testCase *models.TestCase) error {
	var caseOutput *models.CaseOutput
	e.Total++
	var acase *models.CaseStep
	for k, v := range testCase.Steps {
		if k == 0 {
			acase = v
		} else {
			acase = e.MakeACase(testCase.Steps[k-1], v)
		}
		caseOutput = e.RunSteps(acase)
		log.Printf("Run case [%s] step [%d:%s] ==> %t, Verfication: %s\n", testCase.Name, k, acase.Name, caseOutput.Success, acase.Verfication)
		if caseOutput.Success == false {
			e.Failed++
			return caseOutput.Error
		}
		v.Result = caseOutput
	}
	e.Success++
	return nil
}

func (e *Executer) MakeACase(lastStep, nextSetp *models.CaseStep) *models.CaseStep {
	caseStep := &models.CaseStep{
		Name:        nextSetp.Name,
		Level:       nextSetp.Level,
		Method:      nextSetp.Method,
		URL:         nextSetp.URL,
		Data:        nextSetp.Data,
		ResType:     nextSetp.ResType,
		Verfication: nextSetp.Verfication,

		Result: nil,
	}
	addr := nextSetp.URL
	data := nextSetp.Data
	vfdata := nextSetp.Verfication
	keys := make([]string, 0)
	keys = append(keys, GetKeysByString(addr)...)
	keys = append(keys, GetKeysByString(data)...)
	keys = append(keys, GetKeysByString(vfdata)...)

	vals := GetValuesByBody(keys, string(lastStep.Result.Body))
	for k, v := range vals {
		addr = strings.Replace(addr, fmt.Sprintf("{{%s}}", k), v, -1)
		data = strings.Replace(data, fmt.Sprintf("{{%s}}", k), v, -1)
		vfdata = strings.Replace(vfdata, fmt.Sprintf("{{%s}}", k), v, -1)
	}
	caseStep.URL = addr
	caseStep.Data = data
	caseStep.Verfication = vfdata

	return caseStep
}

func (e *Executer) RunSteps(caseStep *models.CaseStep) (caseOutput *models.CaseOutput) {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "unit test httpclient",
		"Accept-Language":        "en-us",
	})
	var res *httpclient.Response
	var err error

	addr := caseStep.URL
	data := caseStep.Data

	caseOutput = new(models.CaseOutput)

	switch strings.ToUpper(caseStep.Method) {
	case "GET":
		res, err = httpclient.Get(addr)
	case "POST":
		values, err := url.ParseQuery(data)
		if err != nil {
			caseOutput.Success = false
			caseOutput.Error = err
			return
		}
		res, err = httpclient.Post(addr, values)
	case "POSTJSON":
		res, err = httpclient.PostJson(addr, data)
	case "PUTJSON":
		res, err = httpclient.PutJson(addr, data)
	case "DELETE":
		res, err = httpclient.Delete(addr)
	default:
		caseOutput.Success = false
		caseOutput.Error = fmt.Errorf("Can not support this method [%s]", caseStep.Method)
	}

	if err != nil {
		caseOutput.Success = false
		caseOutput.Error = err
		return
	}
	body, err := res.ReadAll()
	if err != nil {
		caseOutput.Success = false
		caseOutput.Error = err
		return
	}
	caseOutput.Body = body
	caseOutput.Status = res.StatusCode

	if caseOutput.Status < 200 || caseOutput.Status >= 400 {
		caseOutput.Success = false
		caseOutput.Error = fmt.Errorf("Incorrect HTTP status code %d", caseOutput.Status)
		return
	}
	var success bool
	if caseStep.ResType == "JSON" {
		success, err = VerifyJSON(caseStep.Verfication, string(body))
	} else if caseStep.ResType == "TEXT" {
		success, err = VerifyTEXT(caseStep.Verfication, strings.Replace(string(body), `"`, `\"`, -1))
	}
	caseOutput.Success = success
	caseOutput.Error = err
	return
}

func VerifyJSON(verifier, data string) (bool, error) {
	vm := otto.New()
	tpl := `
		var $=%s;
		var success=false;
		if(%s){
			success=true;
		}
	`
	vmdata := fmt.Sprintf(tpl, data, verifier)
	// log.Println(vmdata)
	_, err := vm.Run(vmdata)
	if err != nil {
		return false, err
	}
	if value, err := vm.Get("success"); err == nil {
		if valueBool, err := value.ToBoolean(); err != nil || valueBool == false {
			return false, nil
		}
	}
	return true, nil
}

func VerifyTEXT(verifier, data string) (bool, error) {
	vm := otto.New()
	tpl := `
		var $="%s";
		var success=false;
		if(%s){
			success=true;
		}
	`
	vmdata := fmt.Sprintf(tpl, data, verifier)
	// log.Println(vmdata)
	_, err := vm.Run(vmdata)
	if err != nil {
		return false, err
	}
	if value, err := vm.Get("success"); err == nil {
		if valueBool, err := value.ToBoolean(); err != nil || valueBool == false {
			return false, nil
		}
	}
	return true, nil
}

func GetKeysByString(data string) []string {
	keys := regexp.MustCompile(`\{\{(.*?)\}\}`).FindAllString(data, -1)
	for k, v := range keys {
		keys[k] = strings.Replace(strings.Replace(v, "}}", "", -1), "{{", "", -1)
	}
	return keys
}

func GetValuesByBody(keys []string, data string) map[string]string {
	vals := make(map[string]string)
	for _, v := range keys {
		vals[v] = parseValue(v, data)
	}
	return vals
}

func parseValue(k, data string) string {
	vm := otto.New()
	tpl := `
		var $=%s;
		var flag=%s.toString();
	`
	vmdata := fmt.Sprintf(tpl, data, k)
	// log.Println(vmdata)
	_, err := vm.Run(vmdata)
	if err != nil {
		return ""
	}
	if value, err := vm.Get("flag"); err == nil {
		if valueStr, err := value.ToString(); err == nil {
			return valueStr
		}
	}
	return ""
}
