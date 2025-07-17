package main

// markdown specification files was stolen from https://github.com/markedjs/marked project.
// Including spec/commonmark/commonmark.0.31.2.json file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
)

type MDTest struct {
	md, html, section string
}

type TestSection struct {
	tests []MDTest
	name  string
}

type TestSections []TestSection

type TestSuite struct {
	name      string
	path      string
	format    FileFormat
	specTests func(s *TestSuite, unmarshaled interface{}) (TestSections, error)
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
		specTests: func(s *TestSuite, unmarshaled interface{}) (TestSections, error) {
			switch tests := unmarshaled.(type) {
			case []interface{}:
				var testSections TestSections
				for _, test := range tests {
					if t, ok := test.(map[string]interface{}); ok {
						sectionName := t["section"].(string)
						if sectionName == "" {
							// TODO(kra53n): put it to TestSections
							sectionName = "general"
						}
						testSections.add(sectionName, MDTest{
							md:   t["markdown"].(string),
							html: t["html"].(string),
						})
					}
				}
				return testSections, nil
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

func (ts *TestSections) add(name string, test MDTest) {
	for i := range *ts {
		if (*ts)[i].name == name {
			(*ts)[i].tests = append((*ts)[i].tests, test)
			return
		}
	}
	*ts = append((*ts), TestSection{
		name:  name,
		tests: []MDTest{test},
	})
}

func (s *TestSuite) MDTests() (TestSections, error) {
	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return nil, err
	}

	unmarshaled, err := s.unmarshal(data)
	if err != nil {
		return nil, err
	}

	if s.specTests == nil {
		return nil, errors.New(fmt.Sprint("run ", s.path, ": parse function was not attached in Spec struct value for specTests slice"))
	}
	return s.specTests(s, unmarshaled)
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

func renderHTMLFromMD(md string) string {
	mdBytes := []byte(md)
	tks := Lex(mdBytes)
	ast := Parse(mdBytes, tks)
	res := Render(mdBytes, ast)
	return res
}

func TestSpecs(t *testing.T) {
	for _, testSuite := range testSuites {
		t.Run(testSuite.name, func(subtest *testing.T) {
			sections, err := testSuite.MDTests()
			if err != nil {
				subtest.Error(err)
			}

			for _, section := range sections {
				if section.name != "List items" {
					continue
				}
				subtest.Run(section.name, func(sectionTest *testing.T) {
					for _, mdTest := range section.tests {
						src := renderHTMLFromMD(mdTest.md)
						if src != mdTest.html {
							sectionTest.Errorf("src != dst\n md: %q\nsrc: %q\ndst: %q\n", mdTest.md, src, mdTest.html)
						}
					}
				})
			}
		})
	}
}
