package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2" // Import the ec2 service
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

	// cmd.Stdin = os.Stdin
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout

	result, err := svc.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Fatalf("failed to describe instances, %v", err)
	}

	for _, reservation := range result.Reservations {
		for _, inst := range reservation.Instances {
			if inst.InstanceId == nil {
				continue
			}
			tagName := helpers.GetTagName(&inst)
			log.Printf("Instance ID: %s, Name: %s", *inst.InstanceId, tagName)
		}
	}
}
