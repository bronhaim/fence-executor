package main

import (
	"log"
	"time"
	"os"
	"fmt"
	"fence-executor/providers"
	"fence-executor/utils"
)

func executeFence(parameters map[string]string) error {
	f := utils.CreateNewFence()
	provider := providers.CreateRHProvider(nil)
	f.RegisterProvider("redhat", provider)
	err := f.LoadAgents(10 * time.Second)
	if err != nil {
		log.Print("error loading agents:", err)
		return err
	}

	ac := utils.NewAgentConfig(parameters["provider"], parameters["agent"])

	ac.SetParameter("--ip", parameters["address"])
	ac.SetParameter("--username", parameters["username"])
	ac.SetParameter("--password", parameters["password"])
	ac.SetParameter("--plug", parameters["port"])

	err = f.Run(ac, utils.Status, 10*time.Second)
	if err != nil {
		log.Print("error: ", err)
		return err
	}
	log.Print("Fenced was executed!")
	return nil
}


func main() {
	var args = os.Args[1:]
	var parameters map[string]string
	parameters = make(map[string]string)
	fmt.Println(args)

	parameters["address"] = args[0]
	parameters["username"] = args[1]
	parameters["password"] = args[2]
	parameters["agent"] = args[3]
	parameters["provider"] = args[4]
	parameters["port"] = args[5]

	fmt.Println(parameters)
	executeFence(parameters)

	fmt.Println("Done")
}