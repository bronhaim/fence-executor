package tests

import (
	"errors"
	"testing"
	"reflect"
	"fence-executor/utils"
	"fence-executor/providers"
)

func TestVerifyAgentConfig(t *testing.T) {
	f := utils.CreateNewFence()
	provider := NewFakeProvider()
	err := provider.LoadAgents(0)
	if err != nil {
		t.Error("error:", err)
	}
	f.RegisterProvider("fakeprovider", provider)

	ac := utils.NewAgentConfig("fakeprovider", "missingagent01")
	err = f.VerifyAgentConfig(ac, false)
	if err == nil {
		t.Error(err)
	}

	ac = utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("missingparam", "bla")
	err = f.VerifyAgentConfig(ac, false)
	if err == nil {
		t.Error(err)
	}

	ac = utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param01", "bla")
	err = f.VerifyAgentConfig(ac, false)
	if err != nil {
		t.Error(err)
	}

	ac = utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param02", "bla")
	err = f.VerifyAgentConfig(ac, false)
	if err == nil {
		t.Error(err)
	}

	ac = utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param03", "option02")
	err = f.VerifyAgentConfig(ac, false)
	if err != nil {
		t.Error(err)
	}

	ac = utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param03", "bla")
	err = f.VerifyAgentConfig(ac, false)
	if err == nil {
		t.Error(err)
	}

	ac = utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param03", 1)
	err = f.VerifyAgentConfig(ac, false)
	want := errors.New("Parameter \"param03\" not of string type")
	if !ErrorEquals(err, want) {
		t.Errorf("Expecting \"%s\" error, found \"%s\"", err, want)
	}

	ac = utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param01", "bla")
	err = f.VerifyAgentConfig(ac, true)
	want = errors.New("Port name required")
	if !ErrorEquals(err, want) {
		t.Errorf("Expecting \"%s\" error, found \"%s\"", err, want)
	}

	ac = utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("param01", "bla")
	ac.SetPort("port01")
	err = f.VerifyAgentConfig(ac, true)
	if err != nil {
		t.Error(err)
	}

}

func TestRunRHProvider(t *testing.T) {
	p := providers.CreateRHProvider(nil)
	a := &providers.RHAgent{
		Command: "/bin/cat",
		Agent: &utils.Agent{
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
			Actions: []utils.Action{
				utils.Off,
			},
		},
	}

	p.Agents[a.Name] = a

	ac := utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetParameter("missingparam", "bla")

	err := p.Run(ac, utils.None, 0)
	if err != nil {
		t.Error(err)
	}
}

func TestToResourceAgent(t *testing.T) {
	tests := []struct {
		xml      string
		expected *providers.RHAgent
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

			&providers.RHAgent{
				Agent: &utils.Agent{Name: "agent01",
					ShortDesc: "Fence agent01",
					LongDesc:  "Fence Agent 01.",
					Parameters: map[string]*utils.Parameter{
						"param01": &utils.Parameter{Name: "param01", Desc: "param01", Unique: false, Required: false, ContentType: utils.String, Default: "123"},
						"param02": &utils.Parameter{Name: "param02", Desc: "param02", Unique: true, Required: true, ContentType: utils.String, Default: nil},
						"param03": &utils.Parameter{Name: "param03", Desc: "param03", Unique: false, Required: true, ContentType: utils.Boolean, Default: false},
						"param04": &utils.Parameter{Name: "param04", Desc: "param04", Unique: true, Required: false, ContentType: utils.Boolean, Default: true},
						"param05": &utils.Parameter{Name: "param05", Desc: "param05",
							ContentType: utils.String,
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
					Actions: []utils.Action{
						utils.On,
						utils.Off,
						utils.Status,
						utils.List,
						utils.Monitor,
					},
					UnfenceAction:   utils.On,
					UnfenceOnTarget: true,
				},
			},
		},
	}
	for _, test := range tests {
		rha, err := providers.ParseMetadata([]byte(test.xml))
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
		rha, err := providers.ParseMetadata([]byte(test.xml))
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
		out utils.Action
		err error
	}{
		{"on", utils.On, nil},
		{"off", utils.Off, nil},
		{"reboot", utils.Reboot, nil},
		{"status", utils.Status, nil},
		{"list", utils.List, nil},
		{"monitor", utils.Monitor, nil},
		{"badaction", 0, errors.New("Unknown fence action: badaction")},
	}

	for _, test := range tests {
		action, err := utils.StringToAction(test.in)

		if !ErrorEquals(err, test.err) {
			t.Errorf("Expecting \"%s\" error, found \"%s\"", err, test.err)
		}
		if action != test.out {
			t.Errorf("Wrong fence action %s for action %s", utils.ActionMap[action], test.in)
		}
	}
}

func TestActionToString(t *testing.T) {
	tests := []struct {
		out string
		in  utils.Action
		err error
	}{
		{"on", utils.On, nil},
		{"off", utils.Off, nil},
		{"reboot", utils.Reboot, nil},
		{"status", utils.Status, nil},
		{"list", utils.List, nil},
		{"monitor", utils.Monitor, nil},
		{"", utils.None, errors.New("Unknown fence action: badaction")},
	}

	for _, test := range tests {
		action := utils.ActionToString(test.in)

		if action != test.out {
			t.Errorf("Wrong fence action %s for action %s", action, utils.ActionMap[test.in])
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
	f := utils.CreateNewFence()
	provider := NewFakeProvider()
	err := provider.LoadAgents(0)
	if err != nil {
		t.Error("error:", err)
	}
	f.RegisterProvider("fakeprovider", provider)

	ac := utils.NewAgentConfig("fakeprovider", "agent01")
	ac.SetPort("port01")

	err = f.Run(ac, utils.On, 0)
	if err != nil {
		t.Error(err)
	}

	err = f.Run(ac, utils.Off, 0)
	if err != nil {
		t.Error(err)
	}
}
