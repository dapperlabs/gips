package core

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/goombaio/namegenerator"
)

type ProjectService interface {
	Find(id string) (*Project, error)
	Store(p *Project) (string, error)
	GetAll() ([]*Project, error)
}

type Project struct {
	Name   string   `json:"name" faker:"gcpProject" db:"name"`
	Region string   `json:"region" faker:"gcpRegion" db:"region"`
	IPs    []string `json:"ips" faker:"ips" db:"ips"`
}

func UnmarshalJSONProject(j string) (Project, error) {
	var p Project
	err := json.Unmarshal([]byte(j), &p)
	return p, err
}

func FakeProject() Project {
	CustomFakerData()
	p := Project{}
	faker.FakeData(&p) //nolint
	return p
}

func FakeProjectJSON() string {
	p := FakeProject()
	j, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(j)
}

func CustomFakerData() {
	_ = faker.AddProvider("gcpProject", func(v reflect.Value) (interface{}, error) {
		nameGenerator := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())
		name := nameGenerator.Generate()
		return name, nil
	})
	_ = faker.AddProvider("gcpRegion", func(v reflect.Value) (interface{}, error) {
		return "us-west1", nil
	})
	_ = faker.AddProvider("ips", func(v reflect.Value) (interface{}, error) {
		return []string{faker.IPv4(), faker.IPv4(), faker.IPv4(), faker.IPv4()}, nil
	})
}
