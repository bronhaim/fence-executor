package tests

import (
	"fmt"
	"log"
	"time"
	"fence-executor/utils"
)

type FakeProvider struct {
	agents utils.Agents
}

func NewFakeProvider() utils.FenceProvider {
	return &FakeProvider{agents: make(utils.Agents)}
}

func (p *FakeProvider) LoadAgents(timeout time.Duration) error {
	a := &utils.Agent{
		Name: "agent01",
		Parameters: map[string]*utils.Parameter{
			"param01": &utils.Parameter{Name: "param01", ContentType: utils.String},
			"param02": &utils.Parameter{Name: "param03", ContentType: utils.Boolean},
			"param03": &utils.Parameter{Name: "param03",
				ContentType: utils.String,
				HasOptions:  true,
				Options: []interface{}{
					"option01",
					"option02",
				},
			},
		},
		MultiplePorts:   true,
		DefaultAction:   utils.Reboot,
		UnfenceAction:   utils.On,
		UnfenceOnTarget: false,
		Actions: []utils.Action{
			utils.On,
			utils.Off,
			utils.Reboot,
		},
	}
	p.agents[a.Name] = a

	return nil
}

func (p *FakeProvider) GetAgents() (utils.Agents, error) {
	return p.agents, nil
}

func (p *FakeProvider) GetAgent(name string) (*utils.Agent, error) {
	a, ok := p.agents[name]
	if !ok {
		return nil, fmt.Errorf("Unknown agent: %s", name)
	}
	return a, nil
}

func (p *FakeProvider) Status(ac *utils.AgentConfig, timeout time.Duration) (utils.DeviceStatus, error) {
	return utils.Ok, nil
}
func (p *FakeProvider) Monitor(ac *utils.AgentConfig, timeout time.Duration) (utils.DeviceStatus, error) {
	return utils.Ok, nil
}
func (p *FakeProvider) List(ac *utils.AgentConfig, timeout time.Duration) (utils.PortList, error) {
	return utils.PortList{
		utils.PortName{Name: "device01", Alias: "alias01"},
		utils.PortName{Name: "device02", Alias: "alias02"},
	}, nil
}
func (p *FakeProvider) Run(ac *utils.AgentConfig, action utils.Action, timeout time.Duration) error {
	_, err := p.GetAgent(ac.Name)
	if err != nil {
		log.Print("error: ", err)
		return err
	}
	return nil
}
