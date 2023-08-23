package core

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarkshalJSONProject(t *testing.T) {
	p := FakeProject()
	j, err := json.Marshal(p)
	if err != nil {
		t.Error(err)
	}
	p2, err := UnmarshalJSONProject(string(j))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(p, p2) {
		t.Error("Those should match")
	}
}

func TestProjectFakerJSON(t *testing.T) {
	j := FakeProjectJSON()
	if j == "" {
		t.Error("JSON was blank")
	}
}
