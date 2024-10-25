package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/triplemcoder14/ec2-connect/helpers"
)

const defaultUser = "ubuntu"
const defaultKeyPath = "~/.ssh/id_rsa"

func main() {
	// using  flags for the cmd
	userPtr := flag.String("user", defaultUser, "SSH user for login")
	directoryPtr := flag.String("directory", defaultKeyPath, "Directory where SSH keys are stored")
	regionPtr := flag.String("region", "us-east-2", "AWS region to use")
	flag.Parse()

	// load the aws config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*regionPtr))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// set region if not provided via flags
	if *regionPtr == "" {
		*regionPtr = selectRegion()
	}
	cfg.Region = *regionPtr

	svc := ec2.NewFromConfig(cfg)

	filters := []ec2Types.Filter{
		{
			Name:   aws.String("instance-state-name"),
			Values: []string{"running"},
		},
	}

	input := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	result, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to describe instances, %v", err)
	}

	var instances []ec2Types.Instance
	for _, reservation := range result.Reservations {
		instances = append(instances, reservation.Instances...)
	}

	if len(instances) == 0 {
		log.Println("No running instances found.")
		return
	}

	// display available instances
	for _, inst := range instances {
		tagName := helpers.GetTagName(&inst)
		publicIP := ""
		if inst.PublicIpAddress != nil {
			publicIP = *inst.PublicIpAddress
		}
		log.Printf("Instance ID: %s, Name: %s, Public IP: %s", *inst.InstanceId, tagName, publicIP)
	}

	var selectedInstanceID string
	fmt.Print("Enter the Instance ID to connect: ")
	fmt.Scanln(&selectedInstanceID)

	var selectedPublicIP string
	for _, inst := range instances {
		if *inst.InstanceId == selectedInstanceID {
			if inst.PublicIpAddress != nil {
				selectedPublicIP = *inst.PublicIpAddress
			} else {
				log.Fatalf("Instance %s does not have a public IP address.", selectedInstanceID)
			}
			break
		}
	}

	if selectedPublicIP == "" {
		log.Fatalf("Invalid Instance ID: %s. Please try again.", selectedInstanceID)
	}

	// ssh connection logic
	keyPath := expandHomeDir(*directoryPtr) // expantion to home directory

	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	address := fmt.Sprintf("%s:22", selectedPublicIP) // using public ip for ssh connection

	sshConfig := &ssh.ClientConfig{
		User: *userPtr,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	fmt.Printf("Connecting to instance %s...\n", selectedInstanceID)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected successfully!")

	// now you can start a new session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}
	defer session.Close()

	//  run a command on the remote instance
	output, err := session.Output("whoami") // change the command as needed
	if err != nil {
		log.Fatalf("failed to run command: %v", err)
	}

	fmt.Printf("Output: %s\n", output)
}

// function to select a region using fzf
func selectRegion() string {
	cmd := exec.Command("fzf", "--header", "Select AWS Region:")
	cmd.Stdin = nil
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("failed to run fzf: %v", err)
	}
	return string(output)
}

// helper function to expand the home directory
func expandHomeDir(path string) string {
	if home := os.Getenv("HOME"); home != "" {
		path = filepath.Join(home, path[2:]) // strip off the ~ and join with home dir
	}
	return path
}
