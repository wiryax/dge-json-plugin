package jp

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
	dge "github.com/wiryax/direct-graph-engine"
)

type JsonToTabular struct {
	blobStorageId,
	tabularStorageId string
}

func (jtt *JsonToTabular) Execute(gCtx *dge.GraphContext) error {
	rawJson, err := gCtx.GetTabularBlob(jtt.blobStorageId)
	if err != nil {
		return err
	}
	tabular, err := parseJsonToTabular(rawJson.GetRaw())
	if err != nil {
		return err
	}
	gCtx.SetTabularStorage(jtt.tabularStorageId, tabular)
	return nil
}

func parseJsonToTabular(raw []byte) (dge.Tabular, error) {
	if !gjson.Valid(string(raw)) {
		return dge.Tabular{}, errors.New("json is not valid")
	}
	result := gjson.Parse(string(raw))

	return parseValue(result)
}

func parseObject(result gjson.Result, parentCol string) (dge.Tabular, error) {
	var (
		err     error
		tabular dge.Tabular
		stack   []dge.Tabular
	)
	for key, val := range result.Map() {
		if val.IsObject() {
			newTabular, err := parseObject(val, parentCol+key)
			if err != nil {
				return tabular, err
			}

			stack = append(stack, newTabular)
			continue
		} else if val.IsArray() {
			newTabular, err := parseArray(val, key)
			if err != nil {
				return tabular, err
			}
			stack = append(stack, newTabular)
			continue
		} else {
			err = tabular.AddColumn(func(rows []dge.Variable) [][]dge.Variable {
				var temp [][]dge.Variable
				temp = append(temp, append(rows, dge.ParseVariable([]byte(val.String()))))
				return temp
			}, parentCol+key)
			if err != nil {
				return tabular, err
			}
		}
	}

	for _, s := range stack {
		tabular, err = tabular.Merge(s)
		if err != nil {
			return tabular, err
		}
	}

	return tabular, err
}

func parseArray(result gjson.Result, parentCol string) (dge.Tabular, error) {
	var (
		err     error
		tabular dge.Tabular
		queue   []dge.Tabular
		rows    []dge.Variable
	)

	for _, r := range result.Array() {
		if r.IsObject() {
			tempT, err := parseObject(r, parentCol)
			if err != nil {
				return tabular, err
			}
			queue = append(queue, tempT)
			continue
		} else if r.IsArray() {
			tempT, err := parseArray(r, parentCol)
			if err != nil {
				return tabular, err
			}
			queue = append(queue, tempT)
			continue
		} else {
			rows = append(rows, dge.ParseVariable([]byte(r.String())))
			continue
		}
	}

	if rows != nil {
		tabular = *dge.MakeTabular([]string{parentCol})
		for _, row := range rows {
			err = tabular.AddRow(row)
			if err != nil {
				return tabular, err
			}
		}
	}

	for _, s := range queue {
		tabular, err = tabular.Merge(s)
		if err != nil {
			return tabular, err
		}
	}

	return tabular, err
}

func parseValue(result gjson.Result) (dge.Tabular, error) {
	if result.IsObject() {
		return parseObject(result, "")
	} else if result.IsArray() {
		return parseArray(result, "")
	}
	return dge.Tabular{}, fmt.Errorf("cannot parse json with %s type", result.Type)
}
