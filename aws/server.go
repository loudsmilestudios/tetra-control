package aws

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/loudsmilestudios/TetraControl/core"
)

type serverData struct {
	IP     string
	Port   uint16
	Status string
}

// Server represents a game server running on AWS
type Server struct {
	Identifier string            `json:"identifier"`
	TaskArn    string            `json:"task"`
	Metadata   map[string]string `json:"metadata"`
	data       *serverData
}

// GetIP returns the IP address of a server
func (server *Server) GetIP() (string, error) {
	data, err := server.getData()
	return data.IP, err
}

// GetPort returns the port of a server
func (server *Server) GetPort() (uint16, error) {
	data, err := server.getData()
	return data.Port, err
}

// GetIdentifier returns the identifier of a server
func (server *Server) GetIdentifier() (string, error) {
	return server.Identifier, nil
}

// GetState returns the current state of a server
func (server *Server) GetState() (core.ServerState, error) {
	data, err := server.getData()
	if err != nil {
		return core.Unknown, err
	}

	switch data.Status {
	case "PROVISIONING", "PENDING", "ACTIVATING":
		return core.Starting, nil
	case "RUNNING":
		return core.Active, nil
	case "DEACTIVATING", "STOPPING", "DEPROVISIONING", "STOPPED":
		return core.Exiting, nil
	}

	log.Printf("server %v is in unhandled %v status", server.Identifier, data.Status)
	return core.Error, nil
}

// GetMetadata returns the metadata of a server
func (server *Server) GetMetadata() (map[string]string, error) {
	return server.Metadata, nil
}

func (server *Server) getData() (*serverData, error) {
	// If data has already been loaded
	// return that
	if server.data != nil {
		return server.data, nil
	}

	// Look up task information
	result, err := ecsClient.DescribeTasks(&ecs.DescribeTasksInput{
		Cluster: &config.ecsCluster,
		Tasks:   []*string{&server.TaskArn},
	})
	if err != nil {
		return nil, err
	}
	if len(result.Tasks) <= 0 {
		return nil, errors.New("task could not be found")
	}

	// Find game server container
	for _, container := range result.Tasks[0].Containers {
		if container.Name == &config.GameServerContainer {

			// Look up network interface associated with container's private address
			inferfaces, err := ec2Client.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
				Filters: []*ec2.Filter{
					{
						Name: core.Strpointer("addresses.private-ip-address"),
						Values: []*string{
							container.NetworkInterfaces[0].PrivateIpv4Address,
						},
					},
				},
			})
			if err != nil {
				return nil, err
			}

			// If network interface exists, return information
			if len(inferfaces.NetworkInterfaces) > 0 {
				server.data = &serverData{
					IP:     *inferfaces.NetworkInterfaces[0].Association.PublicIp,
					Status: *result.Tasks[0].LastStatus,
					Port:   config.GameServerPort,
				}
				return server.data, nil
			}
		}
	}

	// Return just status if no connection info found
	return &serverData{Status: *result.Tasks[0].LastStatus}, nil
}
