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

func NewJsonToTabular(blobId, tabularId string) *JsonToTabular {
	return &JsonToTabular{
		blobStorageId:    blobId,
		tabularStorageId: tabularId,
	}
}

func (jtt *JsonToTabular) Execute(gCtx *dge.GraphContext) error {
	rawJson, err := gCtx.GetBlob(jtt.blobStorageId)
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
	s := string(raw)
	if !gjson.Valid(s) {
		return dge.Tabular{}, errors.New("json is not valid")
	}

	result := gjson.Parse(s)
	return parseValue(result)
}

func parseObject(result gjson.Result, tabular *dge.Tabular, parentCol string) error {
	var (
		err error
	)
	for key, val := range result.Map() {
		if val.IsObject() {
			err := parseObject(val, tabular, parentCol+key)
			if err != nil {
				return err
			}
			continue
		} else if val.IsArray() {
			err := parseArray(val, tabular, key)
			if err != nil {
				return err
			}
			continue
		} else {
			tabular.AddOrSetColumn(parentCol+key, dge.ParseVariable([]byte(val.String())))
		}
	}

	return err
}

func parseArray(result gjson.Result, tabular *dge.Tabular, parentCol string) error {
	for _, r := range result.Array() {
		if r.IsObject() {
			err := parseObject(r, tabular, parentCol)
			if err != nil {
				return err
			}
		} else if r.IsArray() {
			err := parseArray(r, tabular, parentCol)
			if err != nil {
				return err
			}
		} else {
			tabular.AddOrSetColumn(parentCol, dge.ParseVariable([]byte(r.String())))
		}
	}

	return nil
}

func parseValue(result gjson.Result) (dge.Tabular, error) {
	var (
		tabular = dge.MakeTabular()
		err     error
	)
	if result.IsObject() {
		err = parseObject(result, tabular, "")
	} else if result.IsArray() {
		err = parseArray(result, tabular, "")
	} else {
		err = fmt.Errorf("cannot parse json with %s type", result.Type)
	}
	return *tabular, err
}
