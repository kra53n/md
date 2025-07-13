package main

// markdown specification files was stolen from https://github.com/markedjs/marked project.
// Including spec/commonmark/commonmark.0.31.2.json file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type SpecTest struct {
	md, html, section string
}

type SpecTestInfo struct {
	name   string
	path   string
	format SpecTestFormat
	parse  func(unmarshaled interface{}) []SpecTest
}

type SpecTestFormat int

const (
	Json SpecTestFormat = iota
)

var specTestInfos []SpecTestInfo = []SpecTestInfo{
	SpecTestInfo{
		name:   "commonmmark",
		path:   "spec/commonmark/commonmark.0.31.2.json",
		format: Json,
		parse: func(unmarshaled interface{}) []SpecTest {
			switch tests := unmarshaled.(type) {
			case []interface{}:
				specTests := make([]SpecTest, 0, len(tests))
				for _, test := range tests {
					if t, ok := test.(map[string]interface{}); ok {
						specTests = append(specTests, SpecTest{
							md: t["markdown"].(string),
							html: t["html"].(string),
							section: t["section"].(string),
						})
					}
				}
				return specTests
			default:
				fmt.Println("watafak")
			}
			return nil
		},
	},
	SpecTestInfo{
		name:   "other md standard",
		path:   "spec/other/other.0.0.1.json",
		format: Json,
	},
}

func runSpecTests() {
	for _, spec := range specTestInfos {
		spec.run()
	}
}

func (s *SpecTestInfo) run() {
	specTests, err := s.read()
	if err != nil {
		fmt.Println("[ERROR] ", err)
		return
	}
	for _, specTest := range specTests {
		err := specTest.test(s)
		if err != nil {
			fmt.Println("[ERROR] ", err)
			return
		}
	}
}

func (s *SpecTestInfo) read() ([]SpecTest, error) {
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	var unmarshaled interface{}

	switch s.format {
	case Json:
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprint("parse ", s.path, ": SpecTestFormat(", s.format, ") was not implemented"))
	}

	if s.parse == nil {
		return nil, errors.New(fmt.Sprint("parse ", s.path, ": parse function was not attached in Spec struct value for specTests slice"))
	}

	return s.parse(unmarshaled), nil
}

func (st *SpecTest) test(i *SpecTestInfo) error {
	fmt.Println(i.name, i.path, st.section)
	return nil
}
