package ptegodactyl

import (
	"strconv"
	"time"
)

// Node contains node informations
type Node struct {
	Object     string `json:"object"`
	Attributes struct {
		ID                 int       `json:"id"`
		Public             bool      `json:"public"`
		Name               string    `json:"name"`
		Description        string    `json:"description"`
		LocationID         int       `json:"location_id"`
		Fqdn               string    `json:"fqdn"`
		Scheme             string    `json:"scheme"`
		BehindProxy        bool      `json:"behind_proxy"`
		MaintenanceMode    bool      `json:"maintenance_mode"`
		Memory             int       `json:"memory"`
		MemoryOverallocate int       `json:"memory_overallocate"`
		Disk               int       `json:"disk"`
		DiskOverallocate   int       `json:"disk_overallocate"`
		UploadSize         int       `json:"upload_size"`
		DaemonListen       int       `json:"daemon_listen"`
		DaemonSftp         int       `json:"daemon_sftp"`
		DaemonBase         string    `json:"daemon_base"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
	} `json:"attributes"`
	client *AppClient
}

// ListNodes list all servers the user has access to
func (c *AppClient) ListNodes() ([]Node, error) {
	nodes := []Node{}
	err := c.list("/application/nodes", &nodes)
	if err != nil {
		return nil, err
	}
	for k := range nodes {
		nodes[k].client = c
	}
	return nodes, nil
}

// GetNode return a node object based on nodeID
func (c *AppClient) GetNode(id int) (Node, error) {
	s := Node{client: c}
	err := c.get("/application/nodes/"+strconv.Itoa(id), &s)
	if err != nil {
		return s, err
	}
	return s, nil
}
