package hetznerrobot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HetznerRobotServerResponse struct {
	Server HetznerRobotServer `json:"server"`
}

type HetznerRobotServerSubnet struct {
	IP   string `json:"ip"`
	Mask string `json:"mask"`
}

type HetznerRobotServer struct {
	ServerIP         string                     `json:"server_ip"`
	ServerIPv6       string                     `json:"server_ipv6_net"`
	ServerNumber     int                        `json:"server_number"`
	ServerName       string                     `json:"server_name"`
	Product          string                     `json:"product"`
	DataCenter       string                     `json:"dc"`
	Traffic          string                     `json:"traffic"`
	Status           string                     `json:"status"`
	Cancelled        bool                       `json:"cancelled"`
	PaidUntil        string                     `json:"paid_until"`
	IPs              []string                   `json:"ip"`
	Subnets          []HetznerRobotServerSubnet `json:"subnet"`
	LinkedStoragebox int                        `json:"linked_storagebox"`

	Reset   bool `json:"reset"`
	Rescue  bool `json:"rescue"`
	VNC     bool `json:"vnc"`
	Windows bool `json:"windows"`
	Plesk   bool `json:"plesk"`
	CPanel  bool `json:"cpanel"`
	Wol     bool `json:"wol"`
	HotSwap bool `json:"hot_swap"`
}

type HetznerRobotServerRenameRequestBody struct {
	Name string `json:"server_name"`
}

type HetznerRobotServersResponse struct {
	Server []HetznerRobotServer `json:"server"`
}

func (c *HetznerRobotClient) getServer(ctx context.Context, serverNumber int) (*HetznerRobotServer, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/server/%d", c.url, serverNumber), nil, []int{http.StatusOK, http.StatusAccepted})
	if err != nil {
		return nil, err
	}

	serverResponse := HetznerRobotServerResponse{}
	if err = json.Unmarshal(res, &serverResponse); err != nil {
		return nil, err
	}
	return &serverResponse.Server, nil
}

func (c *HetznerRobotClient) getServers(ctx context.Context) ([]HetznerRobotServer, error) {
	res, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/server", c.url), nil, []int{http.StatusOK, http.StatusAccepted})
	if err != nil {
		return nil, err
	}

	// Debug: Print raw response
	fmt.Printf("Raw API response: %s\n", string(res))

	// Unmarshal into an array of server objects
	var serverObjects []map[string]interface{}
	if err = json.Unmarshal(res, &serverObjects); err != nil {
		return nil, fmt.Errorf("failed to unmarshal server objects: %w", err)
	}

	// Convert each server object to our server struct
	servers := make([]HetznerRobotServer, len(serverObjects))
	for i, obj := range serverObjects {
		// Get the server data from the "server" field
		serverData, ok := obj["server"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid server object structure at index %d", i)
		}

		// Convert the server data to JSON
		serverBytes, err := json.Marshal(serverData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal server data: %w", err)
		}

		// Unmarshal into our server struct
		if err = json.Unmarshal(serverBytes, &servers[i]); err != nil {
			return nil, fmt.Errorf("failed to unmarshal server data: %w", err)
		}
	}

	// Debug: Print number of servers found and their details
	fmt.Printf("Number of servers found: %d\n", len(servers))
	for i, server := range servers {
		fmt.Printf("Server %d: %+v\n", i, server)
	}

	return servers, nil
}
