package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"
	"time"
	"fence-executor/fence"
	"fence-executor/fence-providers"
)


type FenceParameters struct {
	Address string
	Username string
	Password string
	options map[string]string
}


func Index(w http.ResponseWriter, r *http.Request) {
	var parameters FenceParameters
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &parameters); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := executeFence(parameters)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}


func executeFence(parameters FenceParameters) error {
	f := fence.New()
	provider := fence_providers.New(nil)
	f.RegisterProvider("redhat", provider)
	err := f.LoadAgents(10 * time.Second)
	if err != nil {
		log.Print("error loading agents:", err)
		return err
	}

	ac := fence.NewAgentConfig("redhat", "fence_apc_snmp")

	ac.SetParameter("--ip", parameters.Address)
	ac.SetParameter("--username", parameters.Username)
	ac.SetParameter("--password", parameters.Password)
	ac.SetParameter("--plug", parameters.options["plug"])

	err = f.Run(ac, fence.Status, 10*time.Second)
	if err != nil {
		log.Print("error: ", err)
		return err
	}
	log.Print("Fenced was executed!")
	return nil
}


func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":7777", router))
}