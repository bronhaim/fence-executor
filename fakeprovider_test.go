package main

import (
	"fmt"
	"log"
	"time"
	"fence-executor/fence-providers"
	"fence-executor/fence"
)

type FakeProvider struct {
	agents fence.Agents
}

func NewFakeProvider() fence_providers.FenceProvider {
	return &FakeProvider{agents: make(fence.Agents)}
}

func (p *FakeProvider) LoadAgents(timeout time.Duration) error {
	a := &fence.Agent{
		Name: "agent01",
		Parameters: map[string]*fence.Parameter{
			"param01": &fence.Parameter{Name: "param01", ContentType: fence.String},
			"param02": &fence.Parameter{Name: "param03", ContentType: fence.Boolean},
			"param03": &fence.Parameter{Name: "param03",
				ContentType: fence.String,
				HasOptions:  true,
				Options: []interface{}{
					"option01",
					"option02",
				},
			},
		},
		MultiplePorts:   true,
		DefaultAction:   fence.Reboot,
		UnfenceAction:   fence.On,
		UnfenceOnTarget: false,
		Actions: []fence.Action{
			fence.On,
			fence.Off,
			fence.Reboot,
		},
	}
	p.agents[a.Name] = a

	return nil
}

func (p *FakeProvider) GetAgents() (fence.Agents, error) {
	return p.agents, nil
}

func (p *FakeProvider) GetAgent(name string) (*fence.Agent, error) {
	a, ok := p.agents[name]
	if !ok {
		return nil, fmt.Errorf("Unknown agent: %s", name)
	}
	return a, nil
}

func (p *FakeProvider) Status(ac *fence.AgentConfig, timeout time.Duration) (fence.DeviceStatus, error) {
	return fence.Ok, nil
}
func (p *FakeProvider) Monitor(ac *fence.AgentConfig, timeout time.Duration) (fence.DeviceStatus, error) {
	return fence.Ok, nil
}
func (p *FakeProvider) List(ac *fence.AgentConfig, timeout time.Duration) (fence.PortList, error) {
	return fence.PortList{
		fence.PortName{Name: "device01", Alias: "alias01"},
		fence.PortName{Name: "device02", Alias: "alias02"},
	}, nil
}
func (p *FakeProvider) Run(ac *fence.AgentConfig, action fence.Action, timeout time.Duration) error {
	_, err := p.GetAgent(ac.Name)
	if err != nil {
		log.Print("error: ", err)
		return err
	}
	return nil
}
