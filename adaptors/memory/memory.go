package memory

import "github.com/darron/gips/core"

type Memory struct{}

func (m *Memory) Find(name string) (*core.Project, error) {
	return nil, nil
}

func (m *Memory) FindIP(ip string) (*core.Project, error) {
	return nil, nil
}

func (m *Memory) GetAll() ([]*core.Project, error) {
	return nil, nil
}

func (m *Memory) Store(p *core.Project) (string, error) {
	return "", nil
}
