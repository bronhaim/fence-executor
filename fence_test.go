package main

import (
	"errors"
	"testing"
	"reflect"
	"fence-executor/fence"
	"fence-executor/providers"
)

func TestVerifyAgentConfig(t *testing.T) {
	f := fence.New()
	provider := NewFakeProvider()
	err := provider.LoadAgents(0)
	if err != nil {
		t.Error("error:", err)
	}
	f.RegisterProvider("fakeprovider", provider)

	ac := fence.NewAgentConfig("fakeprovider", "missingagent01")
	err = f.VerifyAgentConfig(ac, false)
	if err == nil {
		t.Error(err)
	}

	ac = fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("missingparam", "bla")
	err = f.VerifyAgentConfig(ac, false)
	if err == nil {
		t.Error(err)
	}

	ac = fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param01", "bla")
	err = f.VerifyAgentConfig(ac, false)
	if err != nil {
		t.Error(err)
	}

	ac = fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param02", "bla")
	err = f.VerifyAgentConfig(ac, false)
	if err == nil {
		t.Error(err)
	}

	ac = fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param03", "option02")
	err = f.VerifyAgentConfig(ac, false)
	if err != nil {
		t.Error(err)
	}

	ac = fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param03", "bla")
	err = f.VerifyAgentConfig(ac, false)
	if err == nil {
		t.Error(err)
	}

	ac = fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param03", 1)
	err = f.VerifyAgentConfig(ac, false)
	want := errors.New("Parameter \"param03\" not of string type")
	if !ErrorEquals(err, want) {
		t.Errorf("Expecting \"%s\" error, found \"%s\"", err, want)
	}

	ac = fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param01", "bla")
	err = f.VerifyAgentConfig(ac, true)
	want = errors.New("Port name required")
	if !ErrorEquals(err, want) {
		t.Errorf("Expecting \"%s\" error, found \"%s\"", err, want)
	}

	ac = fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param01", "bla")
	ac.SetPort("port01")
	err = f.VerifyAgentConfig(ac, true)
	if err != nil {
		t.Error(err)
	}

}

func TestRunRHProvider(t *testing.T) {
	p := fence_providers.New(nil)
	a := &fence_providers.RHAgent{
		Command: "/bin/cat",
		Agent: &fence.Agent{
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
			Actions: []fence.Action{
				fence.Off,
			},
		},
	}
	p.agents[a.Name] = a

	ac := fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("missingparam", "bla")

	err := p.Run(ac, fence.None, 0)
	if err != nil {
		t.Error(err)
	}
}

func TestToResourceAgent(t *testing.T) {
	tests := []struct {
		xml      string
		expected *fence_providers.RHAgent
	}{
		{`
			<?xml version="1.0" ?>
			<resource-agent name="agent01" shortdesc="Fence agent01" >
			<symlink name="agentalias01" shortdesc="Fence agent alias01"/>
			<symlink name="agentalias02" shortdesc="Fence agent alias02"/>
			<longdesc>Fence Agent 01.</longdesc>
			<vendor-url>vendor01</vendor-url>
			<parameters>
			        <parameter name="param01" unique="0" required="0">
			                <content type="string" default="123"  />
			                <shortdesc lang="en">param01</shortdesc>
			        </parameter>
			        <parameter name="param02" unique="1" required="1">
			                <content type="string"/>
			                <shortdesc lang="en">param02</shortdesc>
			        </parameter>
			        <parameter name="param03" unique="0" required="1">
			                <content type="boolean" default="0" />
			                <shortdesc lang="en">param03</shortdesc>
			        </parameter>
			        <parameter name="param04" unique="1" required="0">
			                <content type="boolean" default="1" />
			                <shortdesc lang="en">param04</shortdesc>
			        </parameter>
			        <parameter name="param05" unique="0" required="0">
			                <content type="select" default="option01"  >
			                        <option value="option01" />
			                        <option value="option02" />
			                </content>
			                <shortdesc lang="en">param05</shortdesc>
			        </parameter>
			</parameters>
			<actions>
			        <action name="on" on_target="1" automatic="1"/>
			        <action name="off" />
			        <action name="status" />
			        <action name="list" />
			        <action name="monitor" />
			        <action name="metadata" />
			</actions>
			</resource-agent>
			`,

			&fence_providers.RHAgent{
				Agent: &fence.Agent{Name: "agent01",
					ShortDesc: "Fence agent01",
					LongDesc:  "Fence Agent 01.",
					Parameters: map[string]*fence.Parameter{
						"param01": &fence.Parameter{Name: "param01", Desc: "param01", Unique: false, Required: false, ContentType: fence.String, Default: "123"},
						"param02": &fence.Parameter{Name: "param02", Desc: "param02", Unique: true, Required: true, ContentType: fence.String, Default: nil},
						"param03": &fence.Parameter{Name: "param03", Desc: "param03", Unique: false, Required: true, ContentType: fence.Boolean, Default: false},
						"param04": &fence.Parameter{Name: "param04", Desc: "param04", Unique: true, Required: false, ContentType: fence.Boolean, Default: true},
						"param05": &fence.Parameter{Name: "param05", Desc: "param05",
							ContentType: fence.String,
							HasOptions:  true,
							Unique:      false,
							Required:    false,
							Default:     "option01",
							Options: []interface{}{
								"option01",
								"option02",
							},
						},
					},
					Actions: []fence.Action{
						fence.On,
						fence.Off,
						fence.Status,
						fence.List,
						fence.Monitor,
					},
					UnfenceAction:   fence.On,
					UnfenceOnTarget: true,
				},
			},
		},
	}
	for _, test := range tests {
		rha, err := fence_providers.ParseMetadata([]byte(test.xml))
		if err != nil {
			t.Error(err)
		}
		agent, err := rha.ToResourceAgent()
		if err != nil {
			t.Error(err)
		}
		if !reflect.DeepEqual(agent, test.expected) {
			t.Error("Agent definition different from the expected one")
		}
	}

	errortests := []struct {
		xml string
		err error
	}{
		{`
			<?xml version="1.0" ?>
			<resource-agent name="agent01" shortdesc="Fence agent01" >
			<parameters>
			        <parameter name="param01" unique="0" required="0">
			                <content type="boolean" default="badbool" />
			                <shortdesc lang="en">param01</shortdesc>
			        </parameter>
			</parameters>
			</resource-agent>
			`,

			errors.New("strconv.ParseBool: parsing \"badbool\": invalid syntax"),
		},
		{`
			<?xml version="1.0" ?>
			<resource-agent name="agent01" shortdesc="Fence agent01" >
			<parameters>
			        <parameter name="param01" unique="0" required="0">
			                <content type="wrongtype" />
			                <shortdesc lang="en">param01</shortdesc>
			        </parameter>
			</parameters>
			</resource-agent>
			`,

			errors.New("Agent: agent01, parameter: param01. Wrong content type: wrongtype"),
		},
	}

	for _, test := range errortests {
		rha, err := fence_providers.ParseMetadata([]byte(test.xml))
		if err != nil {
			t.Error(err)
		}
		_, err = rha.ToResourceAgent()
		if !ErrorEquals(err, test.err) {
			t.Errorf("Expecting \"%s\" error, found \"%s\"", err, test.err)
		}
	}

}

func TestStringToAction(t *testing.T) {
	tests := []struct {
		in  string
		out fence.Action
		err error
	}{
		{"on", fence.On, nil},
		{"off", fence.Off, nil},
		{"reboot", fence.Reboot, nil},
		{"status", fence.Status, nil},
		{"list", fence.List, nil},
		{"monitor", fence.Monitor, nil},
		{"badaction", 0, errors.New("Unknown fence action: badaction")},
	}

	for _, test := range tests {
		action, err := fence.StringToAction(test.in)

		if !ErrorEquals(err, test.err) {
			t.Errorf("Expecting \"%s\" error, found \"%s\"", err, test.err)
		}
		if action != test.out {
			t.Errorf("Wrong fence action %s for action %s", fence.ActionMap[action], test.in)
		}
	}
}

func TestActionToString(t *testing.T) {
	tests := []struct {
		out string
		in  fence.Action
		err error
	}{
		{"on", fence.On, nil},
		{"off", fence.Off, nil},
		{"reboot", fence.Reboot, nil},
		{"status", fence.Status, nil},
		{"list", fence.List, nil},
		{"monitor", fence.Monitor, nil},
		{"", fence.None, errors.New("Unknown fence action: badaction")},
	}

	for _, test := range tests {
		action := fence.ActionToString(test.in)

		if action != test.out {
			t.Errorf("Wrong fence action %s for action %s", action, fence.ActionMap[test.in])
		}
	}
}

func ErrorEquals(err1 error, err2 error) bool {
	if err1 == nil || err2 == nil {
		return err1 == err2
	}
	return err1.Error() == err2.Error()
}


func TestRunFence(t *testing.T) {
	f := fence.New()
	provider := NewFakeProvider()
	err := provider.LoadAgents(0)
	if err != nil {
		t.Error("error:", err)
	}
	f.RegisterProvider("fakeprovider", provider)

	ac := fence.NewAgentConfig("fakeprovider", "agent01")
	ac.SetPort("port01")

	err = f.Run(ac, fence.On, 0)
	if err != nil {
		t.Error(err)
	}

	err = f.Run(ac, fence.Off, 0)
	if err != nil {
		t.Error(err)
	}
}
