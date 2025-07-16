package main

// markdown specification files was stolen from https://github.com/markedjs/marked project.
// Including spec/commonmark/commonmark.0.31.2.json file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type MDTest struct {
	md, html, section string
}

type TestSuite struct {
	name      string
	path      string
	format    FileFormat
	specTests func(s *TestSuite, unmarshaled interface{}) ([]MDTest, error)
	tests     []MDTest
}

type FileFormat int

const (
	Json FileFormat = iota
)

var testSuites []TestSuite = []TestSuite{
	TestSuite{
		name:   "commonmmark",
		path:   "spec/commonmark/commonmark.0.31.2.json",
		format: Json,
		specTests: func(s *TestSuite, unmarshaled interface{}) ([]MDTest, error) {
			switch tests := unmarshaled.(type) {
			case []interface{}:
				specTests := make([]MDTest, 0, len(tests))
				for _, test := range tests {
					if t, ok := test.(map[string]interface{}); ok {
						specTests = append(specTests, MDTest{
							md:      t["markdown"].(string),
							html:    t["html"].(string),
							section: t["section"].(string),
						})
					}
				}
				return specTests, nil
			}
			return nil, errors.New(fmt.Sprint("extractFields ", s.path, ": due to some reasons"))
		},
	},
	TestSuite{
		name:   "other md standard",
		path:   "spec/other/other.0.0.1.json",
		format: Json,
	},
}

func runSpecTests() {
	for _, s := range testSuites {
		s.run()
	}
}

func (s *TestSuite) run() {
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		printTestErr(err)
	}

	unmarshaled, err := s.unmarshal(data)
	if err != nil {
		printTestErr(err)
	}

	if s.specTests == nil {
		printTestErr(errors.New(fmt.Sprint("run ", s.path, ": parse function was not attached in Spec struct value for specTests slice")))
	}

	specTests, err := s.specTests(s, unmarshaled)
	for _, specTest := range specTests {
		err := specTest.test(s)
		if err != nil {
			printTestErr(err)
		}
	}
}

func (s *TestSuite) unmarshal(data []byte) (interface{}, error) {
	var err error
	var unmarshaled interface{}

	switch s.format {
	case Json:
		err = json.Unmarshal(data, &unmarshaled)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprint("unmarshal ", s.path, ": SpecTestFormat(", s.format, ") was not implemented"))
	}

	return unmarshaled, nil
}

func printTestErr(e error) {
	fmt.Println("[ERROR] ", e)
	os.Exit(69)
}

func (t *MDTest) test(s *TestSuite) error {
	fmt.Println(s.name, s.path, t.section)
	return nil
}

func TestSpecs(t *testing.T) {
	for _, testSuite := range testSuites {
		_ = testSuite
	}
}
