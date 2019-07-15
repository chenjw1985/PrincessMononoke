package parser

import (
	"bytes"
	"fmt"

	"PrincessMononoke/models"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func ConvertToFuncCases(byteData []byte) ([]*models.FuncCase, error) {
	r := bytes.NewReader(byteData)
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	sheets := f.GetSheetMap()
	funcCases := make([]*models.FuncCase, len(sheets))
	for k, v := range sheets {
		rows, err := f.GetRows(v)
		if err != nil {
			return nil, err
		}
		funcCase, err := parseSheet(v, rows[1:])
		if err != nil {
			return nil, err
		}
		funcCases[k-1] = funcCase
	}

	return funcCases, nil
}

func parseSheet(name string, rows [][]string) (*models.FuncCase, error) {
	funcCase := &models.FuncCase{
		Name:  name,
		Cases: make([]*models.TestCase, 0, len(rows)),
	}

	var last *models.TestCase
	for _, v := range rows {
		if v[0] != "" {
			if last != nil {
				funcCase.Cases = append(funcCase.Cases, last)
			}
			last = &models.TestCase{
				Name:  v[0],
				Steps: make([]*models.CaseStep, 0),
			}
		}
		caseStep, err := parseRow(v)
		if err != nil {
			return nil, err
		}
		last.Steps = append(last.Steps, caseStep)
	}
	funcCase.Cases = append(funcCase.Cases, last)

	return funcCase, nil
}

func parseRow(cols []string) (*models.CaseStep, error) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover a panic error %v", r)
		}
	}()
	caseStep := &models.CaseStep{
		Level:       cols[1],
		Method:      cols[2],
		Name:        cols[3],
		URL:         cols[4],
		Data:        cols[5],
		ResType:     cols[6],
		Verfication: cols[7],
		Result:      &models.CaseOutput{},
	}
	return caseStep, err
}
