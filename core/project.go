package core

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/goombaio/namegenerator"
)

type ProjectService interface {
	Find(name string) (*Project, error)
	FindIP(ip string) (*Project, error)
	GetAll() ([]*Project, error)
	Store(p *Project) (string, error)
}

type Project struct {
	Name    string             `json:"name" faker:"gcpProject" db:"name"`
	Regions []ProjectRegionIPs `json:"regions" db:"regions"`
}

type ProjectRegionIPs struct {
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
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		regions := []string{"us-west1", "us-east4", "europe-west2", "us-west2", "asia-south1"}
		n := r.Intn(len(regions))
		return regions[n], nil
	})
	_ = faker.AddProvider("ips", func(v reflect.Value) (interface{}, error) {
		return []string{faker.IPv4(), faker.IPv4(), faker.IPv4(), faker.IPv4()}, nil
	})
}
