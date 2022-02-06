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

		tenantStackID := cfg.Require("tenantStack")
		fmt.Println(tenantStackID)

		// tenantStack := pulumi.StackReference(tenantStackID)
		tenantStack, err := pulumi.NewStackReference(ctx, tenantStackID, nil)
		if err != nil {
			fmt.Println("error to get stack")
			panic(err)
		}

		tenant := tenantStack.GetStringOutput(pulumi.String("tenant"))
		environment := tenantStack.GetStringOutput(pulumi.String("environment"))

		instanceType := cfg.Require("instanceType")
		amiId := cfg.Require("ami")
		keyName := cfg.Require("keyName")

		sharedNetworkStackID := cfg.Require("sharedNetworkStack")
		fmt.Println(sharedNetworkStackID)

		sharedNetworkStack, err := pulumi.NewStackReference(ctx, sharedNetworkStackID, nil)
		if err != nil {
			fmt.Println("error to get network stack")
			panic(err)
		}

		networkPrivateSubnetZoneC := sharedNetworkStack.GetStringOutput(pulumi.String("privateSubnetZoneC"))
		networkSecurityGroupArn := sharedNetworkStack.GetStringOutput(pulumi.String("securityGroupArn"))

		computeTestId := pulumi.All(tenant, environment, instanceType, networkPrivateSubnetZoneC, networkSecurityGroupArn, amiId, keyName).ApplyT(func(args []interface{}) (pulumi.IDOutput, error) {
			hostName := fmt.Sprintf("tenant-%s-%v-test-host-1", args[0], args[1])

			server, err := ec2.NewInstance(ctx, hostName, &ec2.InstanceArgs{
				InstanceType: pulumi.String(args[2].(string)),
				SubnetId:     pulumi.String(args[3].(string)),
				//VpcSecurityGroupIds: pulumi.StringArray{pulumi.String(args[4].(string))},
				Ami: pulumi.String(args[5].(string)),
				//KeyName: pulumi.String(args[6].(string)),
				//DisableApiTermination: pulumi.Bool(true),
				Tags: pulumi.StringMap{
					"Name": pulumi.String(hostName),
				},
			})
			if err != nil {
				return pulumi.IDOutput{}, err
			}

			return server.ID(), nil
		})
		if err != nil {
			return err
		}

		ctx.Export("computeTestId", computeTestId)

		return nil
	})
}
