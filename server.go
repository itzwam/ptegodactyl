package ptegodactyl

import (
	"fmt"
	"time"
)

// Server is all infos on a server
type Server struct {
	Attributes struct {
		Description   string `json:"description"`
		FeatureLimits struct {
			Allocations int64 `json:"allocations"`
			Databases   int64 `json:"databases"`
		} `json:"feature_limits"`
		Identifier string `json:"identifier"`
		Limits     struct {
			CPU    int64 `json:"cpu"`
			Disk   int64 `json:"disk"`
			Io     int64 `json:"io"`
			Memory int64 `json:"memory"`
			Swap   int64 `json:"swap"`
		} `json:"limits"`
		Name        string `json:"name"`
		ServerOwner bool   `json:"server_owner"`
		UUID        string `json:"uuid"`
	} `json:"attributes"`
	Object string `json:"object"`
	client *Client
}

// AppServer is all infos on server (with admin key)
type AppServer struct {
	Object     string `json:"object"`
	Attributes struct {
		ID          int         `json:"id"`
		ExternalID  interface{} `json:"external_id"`
		UUID        string      `json:"uuid"`
		Identifier  string      `json:"identifier"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Suspended   bool        `json:"suspended"`
		Limits      struct {
			Memory int `json:"memory"`
			Swap   int `json:"swap"`
			Disk   int `json:"disk"`
			Io     int `json:"io"`
			CPU    int `json:"cpu"`
		} `json:"limits"`
		FeatureLimits struct {
			Databases   int `json:"databases"`
			Allocations int `json:"allocations"`
		} `json:"feature_limits"`
		User       int         `json:"user"`
		Node       int         `json:"node"`
		Allocation int         `json:"allocation"`
		Nest       int         `json:"nest"`
		Egg        int         `json:"egg"`
		Pack       interface{} `json:"pack"`
		Container  struct {
			StartupCommand string `json:"startup_command"`
			Image          string `json:"image"`
			Installed      bool   `json:"installed"`
			Environment    struct {
				SERVERJARFILE   string `json:"SERVER_JARFILE"`
				VANILLAVERSION  string `json:"VANILLA_VERSION"`
				STARTUP         string `json:"STARTUP"`
				PSERVERLOCATION string `json:"P_SERVER_LOCATION"`
				PSERVERUUID     string `json:"P_SERVER_UUID"`
			} `json:"environment"`
		} `json:"container"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"attributes"`
	client *AppClient
}

// ServerStats is all stats on a MC server
type ServerStats struct {
	Object     string `json:"object"`
	Attributes struct {
		State  string `json:"state"`
		Memory struct {
			Current int `json:"current"`
			Limit   int `json:"limit"`
		} `json:"memory"`
		CPU struct {
			Current float64   `json:"current"`
			Cores   []float64 `json:"cores"`
			Limit   int       `json:"limit"`
		} `json:"cpu"`
		Disk struct {
			Current int `json:"current"`
			Limit   int `json:"limit"`
		} `json:"disk"`
	} `json:"attributes"`
}

// ListServers list all servers the user has access to
func (c *Client) ListServers() ([]Server, error) {
	servers := []Server{}
	err := c.list("/client", &servers)
	if err != nil {
		return nil, err
	}
	for k := range servers {
		servers[k].client = c
	}
	return servers, nil
}

// GetServer return a server object based on serverID
func (c *Client) GetServer(id string) (Server, error) {
	s := Server{client: c}
	err := c.get("/client/servers/"+id, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// ListServers list all servers the user has access to
func (c *AppClient) ListServers() ([]AppServer, error) {
	servers := []AppServer{}
	err := c.list("/application/servers", &servers)
	if err != nil {
		return nil, err
	}
	for k := range servers {
		servers[k].client = c
	}
	return servers, nil
}

// GetServer return a server object based on serverID
func (c *AppClient) GetServer(id string) (AppServer, error) {
	s := AppServer{client: c}
	err := c.get("/application/servers/"+id, &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

// GetStats return server stats based on serverID
func (s *Server) GetStats() (ServerStats, error) {
	var stats ServerStats
	err := s.client.get("/client/servers/"+s.Attributes.Identifier+"/utilization", &stats)
	if err != nil {
		return stats, err
	}
	return stats, nil
}

// ConsolePayload is a payload containing command
type ConsolePayload struct {
	Command string `json:"command"`
}

// SendConsoleCommand sends a command to the server
func (s *Server) SendConsoleCommand(command string) error {
	cmd := ConsolePayload{Command: command}
	err := s.client.send("/client/servers/"+s.Attributes.Identifier+"/command", &cmd, nil)
	if err != nil {
		return err
	}
	return nil
}

// PowerAction is a power action
type PowerAction string

const (
	PowerStart   PowerAction = "start"
	PowerStop    PowerAction = "stop"
	PowerKill    PowerAction = "kill"
	PowerRestart PowerAction = "restart"
)

// PowerPayload is a payload containing powerAction
type PowerPayload struct {
	State string `json:"signal"`
}

// SendPowerAction sends a command to the server
func (s *Server) SendPowerAction(state PowerAction) error {
	cmd := PowerPayload{State: fmt.Sprint(state)}
	err := s.client.send("/client/servers/"+s.Attributes.Identifier+"/power", &cmd, nil)
	if err != nil {
		return err
	}
	return nil
}
