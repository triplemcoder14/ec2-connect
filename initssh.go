package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2" // Import the EC2 service
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/triplemcoder14/ec2-connect/helpers"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

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

	var instanceIDs []string
	for _, reservation := range result.Reservations {
		for _, inst := range reservation.Instances {
			if inst.InstanceId == nil {
				continue
			}
			tagName := helpers.GetTagName(&inst)
			instanceID := *inst.InstanceId
			instanceIDs = append(instanceIDs, instanceID)
			log.Printf("Instance ID: %s, Name: %s", instanceID, tagName)
		}
	}

	if len(instanceIDs) == 0 {
		log.Println("No running instances found.")
		return
	}

	var selectedInstanceID string
	fmt.Print("Enter the Instance ID to connect: ")
	fmt.Scanln(&selectedInstanceID)

	if !helpers.Contains(instanceIDs, selectedInstanceID) {
		log.Fatalf("Invalid Instance ID: %s. Please try again.", selectedInstanceID)
	}

	// SSH connection logic
	user := "ubuntu"                   // Change this if you're using a different user
	keyPath := "/path/to/your/key.pem" // Specify the path to your SSH key

	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	address := fmt.Sprintf("%s:22", selectedInstanceID) // Replace with the public DNS if needed

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // WARNING: This is insecure; use a proper host key verification in production.
	}

	fmt.Printf("Connecting to instance %s...\n", selectedInstanceID)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer client.Close()

	fmt.Println("Connected successfully!")

	// Now you can start a new session
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("failed to create session: %v", err)
	}
	defer session.Close()

	// Example: Run a command on the remote instance
	output, err := session.Output("whoami") // Change the command as needed
	if err != nil {
		log.Fatalf("failed to run command: %v", err)
	}

	fmt.Printf("Output: %s\n", output)
}



// package main

// import (
// 	"context"
// 	"log"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/ec2" // Import the ec2 service
// 	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
// 	"github.com/triplemcoder14/ec2-connect/helpers"
// )

// func main() {
// 	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
// 	if err != nil {
// 		log.Fatalf("unable to load SDK config, %v", err)
// 	}

// 	svc := ec2.NewFromConfig(cfg)

// 	filters := []ec2Types.Filter{
// 		{
// 			Name:   aws.String("instance-state-name"),
// 			Values: []string{"running"},
// 		},
// 	}

// 	input := &ec2.DescribeInstancesInput{
// 		Filters: filters,
// 	}

// 	// cmd.Stdin = os.Stdin
// 	// cmd.Stderr = os.Stderr
// 	// cmd.Stdout = os.Stdout

// 	result, err := svc.DescribeInstances(context.TODO(), input)
// 	if err != nil {
// 		log.Fatalf("failed to describe instances, %v", err)
// 	}

// 	for _, reservation := range result.Reservations {
// 		for _, inst := range reservation.Instances {
// 			if inst.InstanceId == nil {
// 				continue
// 			}
// 			tagName := helpers.GetTagName(&inst)
// 			log.Printf("Instance ID: %s, Name: %s", *inst.InstanceId, tagName)
// 		}
// 	}
// }
