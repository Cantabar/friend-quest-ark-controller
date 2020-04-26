package core

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/lhw/go-pkg-ark/arkrcon"
)

// AcsInstance instance of ACS instance
type AcsInstance struct {
	InstanceID                    *string
	Name, PublicIPAddress, Status string
	ActivePlayers                 int
}

// getRconPasswordFromEnv handles retrieving the Rcon Password from the Environment
func getRconPasswordFromEnv() string {
	return os.Getenv("ACSRCONPASS")
}

// GetAcsHostNames retrieves host names
func GetAcsHostNames(instances []AcsInstance) []string {
	var acsHosts []string
	for _, inst := range instances {
		acsHosts = append(acsHosts, inst.Name)
	}
	return acsHosts
}

// FindInstanceByAcsHost finds instance by host name
func FindInstanceByAcsHost(instances []AcsInstance, acsHost string) int {
	for i, inst := range instances {
		if strings.Compare(inst.Name, acsHost) == 0 {
			return i
		}
	}
	return -1
}

// SaveWorldByHost sends a saveworld command to the specified host
func SaveWorldByHost(host string) {
	var port string = "27020"
	var addressBuilder strings.Builder
	addressBuilder.WriteString(host)
	addressBuilder.WriteString(":")
	addressBuilder.WriteString(port)
	address := addressBuilder.String()
	rconPass := getRconPasswordFromEnv()
	conn, connErr := arkrcon.NewARKRconConnection(address, rconPass)
	cmdErr := conn.SaveWorld()
	if connErr != nil {
		fmt.Println(connErr)
	}
	if cmdErr != nil {
		fmt.Println(cmdErr)
	}
}

// GetActivePlayers gets active players on provided host
func GetActivePlayers(host string) int {
	if host == "" {
		return -1
	}
	var port string = "27020"
	var addressBuilder strings.Builder
	addressBuilder.WriteString(host)
	addressBuilder.WriteString(":")
	addressBuilder.WriteString(port)
	address := addressBuilder.String()
	rconPass := getRconPasswordFromEnv()
	conn, connErr := arkrcon.NewARKRconConnection(address, rconPass)
	players, playerCmdErr := conn.ListPlayers()
	if connErr != nil {
		fmt.Println(connErr)
		return -1
	}
	if playerCmdErr != nil {
		fmt.Println(playerCmdErr)
		return -1
	}
	return len(players)
}

// GetInstancesAndActivePlayers Gets all active players
func GetInstancesAndActivePlayers() []AcsInstance {
	instances := GetAcsInstances()
	var acsInstances []AcsInstance
	for _, inst := range instances {
		var acsInstance AcsInstance = inst
		acsInstance.ActivePlayers = GetActivePlayers(inst.PublicIPAddress)
		acsInstances = append(acsInstances, acsInstance)
	}
	return acsInstances
}

// GetAcsInstances gets acs instances
func GetAcsInstances() []AcsInstance {
	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")}))
	svc := ec2.New(newSession)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []*string{aws.String("ACS-Host")},
			},
		},
	}
	resp, err := svc.DescribeInstances(params)
	if err != nil {
		fmt.Println(err)
	}
	var acsInstances []AcsInstance
	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			acsInstance := AcsInstance{
				InstanceID:      inst.InstanceId,
				Name:            "",
				Status:          *inst.State.Name,
				PublicIPAddress: "",
			}
			if *inst.State.Code == 16 {
				acsInstance.PublicIPAddress = *inst.PublicIpAddress
			}
			for _, tag := range inst.Tags {
				if strings.Compare(*tag.Key, "ACS-Host") == 0 {
					acsInstance.Name = *tag.Value
				}
			}
			acsInstances = append(acsInstances, acsInstance)
		}
	}
	return acsInstances
}

// StartAcsInstance starts given acs instance
func StartAcsInstance(inst AcsInstance) {
	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")}))
	svc := ec2.New(newSession)
	params := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			inst.InstanceID,
		},
	}
	resp, err := svc.StartInstances(params)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}

// StopAcsInstance stops given acs instance
func StopAcsInstance(inst AcsInstance) {
	fmt.Println("Attempting to save world")
	SaveWorldByHost(inst.PublicIPAddress)
	fmt.Println("Stopping instance: ", inst.Name)
	newSession := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")}))
	svc := ec2.New(newSession)
	params := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			inst.InstanceID,
		},
	}
	resp, err := svc.StopInstances(params)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp)
	}
}

// OpenSSHToInstance opens an SSH terminal to instance
func OpenSSHToInstance(inst AcsInstance) {
	fmt.Println("Opening SSH Terminal to: ", inst.PublicIPAddress)
	host := "ubuntu@" + inst.PublicIPAddress
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", host)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}
