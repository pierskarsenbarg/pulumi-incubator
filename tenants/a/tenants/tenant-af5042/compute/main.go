package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		instanceType := cfg.Require("instanceType")
		availabilityZone := cfg.Require("availabilityZone")

		sharedNetworkStackID := cfg.Require("sharedNetworkStack")
		fmt.Println(sharedNetworkStackID)

		sharedNetworkStack, err := pulumi.NewStackReference(ctx, sharedNetworkStackID, nil)
		if err != nil {
			fmt.Println("error to get network stack")
			panic(err)
		}

		networkVpcId := sharedNetworkStack.GetStringOutput(pulumi.String("vpcId"))
		networkSecurityGroupArn := sharedNetworkStack.GetStringOutput(pulumi.String("securityGroupArn"))

		computeTestArn := pulumi.All(networkVpcId, networkSecurityGroupArn, availabilityZone, instanceType).ApplyT(func(args []interface{}) (pulumi.IDOutput, pulumi.StringOutput, error) {
			hostName := fmt.Sprintf("tenant-%s-%v-test-host-1", args[0], args[1])
			host, err := ec2.NewDedicatedHost(ctx, hostName, &ec2.DedicatedHostArgs{
				AvailabilityZone: pulumi.String(availabilityZone),
				InstanceType:     pulumi.String(instanceType),
			})
			if err != nil {
				return pulumi.IDOutput{}, pulumi.StringOutput{}, err
			}
			return host.ID(), host.Arn, nil

		})
		if err != nil {
			return err
		}

		ctx.Export("computeTestArn", computeTestArn)

		return nil
	})
}
