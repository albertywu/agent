package buildkite

import (
	"fmt"
	"github.com/buildkite/agent/buildkite/logger"
)

type AgentRegistration struct {
	// The access token for the agent
	AccessToken string `json:"access_token"`

	// Hostname of the machine
	Hostname string `json:"hostname"`

	// Operating system for this machine
	OS string `json:"os"`

	// If this agent is allowed to perform script evaluation
	ScriptEval bool `json:"script_eval_enabled"`

	// The priority of the agent
	Priority string `json:"priority,omitempty"`

	// The name of the new agent
	Name string `json:"name"`

	// Meta data for the agent
	MetaData []string `json:"meta_data"`
}

func (c *Client) AgentRegister(name string, priority string, metaData []string, scriptEval bool) (string, error) {
	os, err := MachineOSDump()
	hostname, err := MachineHostname()

	// Create the agent registration
	var registration AgentRegistration
	registration.Name = name
	registration.Priority = priority
	registration.Hostname = hostname
	registration.OS = os
	registration.MetaData = metaData
	registration.ScriptEval = scriptEval

	logger.Info("Registering Agent with Buildkite...")

	// Register and return the agent
	err = c.Post(&registration, "/register", registration)
	if err != nil {
		return "", err
	}

	logger.Info("Agent successfully registered")

	return registration.AccessToken, nil
}

func (a *AgentRegistration) String() string {
	return fmt.Sprintf("AgentRegistration{AccessToken: %s}", a.AccessToken)
}
