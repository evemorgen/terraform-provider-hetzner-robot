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

	// Unmarshal into the response struct
	var response HetznerRobotServersResponse
	if err = json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal servers response: %w", err)
	}

	// Debug: Print number of servers found and their details
	fmt.Printf("Number of servers found: %d\n", len(response.Server))
	for i, server := range response.Server {
		fmt.Printf("Server %d: %+v\n", i, server)
	}

	return response.Server, nil
}
