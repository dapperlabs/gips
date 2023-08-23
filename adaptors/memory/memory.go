package memory

import "github.com/darron/gips/core"

type Memory struct {
	data map[string]*core.Project
}

func New() *Memory {
	return &Memory{
		data: make(map[string]*core.Project),
	}
}

func (m *Memory) Find(name string) (*core.Project, error) {
	p, ok := m.data[name]
	if ok {
		return p, nil
	}
	return nil, nil
}

func (m *Memory) FindIP(ip string) (*core.Project, error) {
	return nil, nil
}

func (m *Memory) GetAll() ([]*core.Project, error) {
	var projects []*core.Project
	for _, p := range m.data {
		projects = append(projects, p)
	}
	return projects, nil
}

func (m *Memory) Store(p *core.Project) (string, error) {
	m.data[p.Name] = p
	return "", nil
}
