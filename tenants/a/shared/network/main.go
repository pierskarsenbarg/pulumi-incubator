package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		environment := ctx.Stack()

		cfg := config.New(ctx, "")
		name := cfg.Require("name")
		//region := cfg.Require("region")
		cidrBlockVpcRange := cfg.Require("cidrBlockVpc")
		cidrBlockPublicRange := cfg.Require("cidrBlockPublic")
		cidrBlockPrivateRange := cfg.Require("cidrBlockPrivate")

		vpcName := fmt.Sprintf("%s-%s", name, environment)
		vpcPublicName := fmt.Sprintf("%s-public", vpcName)
		vpcPrivateName := fmt.Sprintf("%s-private", vpcName)

		tags := pulumi.StringMap{
			"Name": pulumi.String(vpcName),
		}
		tagsPublic := pulumi.StringMap{
			"Name": pulumi.String(vpcPublicName),
		}
		tagsPrivate := pulumi.StringMap{
			"Name": pulumi.String(vpcPrivateName),
		}

		cidrBlockVPC := pulumi.String(cidrBlockVpcRange)
		cidrBlockPublic := pulumi.String(cidrBlockPublicRange)
		cidrBlockPrivate := pulumi.String(cidrBlockPrivateRange)

		// VPC
		vpc, err := ec2.NewVpc(ctx, vpcName, &ec2.VpcArgs{
			CidrBlock: cidrBlockVPC,
			Tags:      tags,
		})
		if err != nil {
			return err
		}

		// Subnets
		publicSubnet, err := ec2.NewSubnet(ctx, vpcPublicName, &ec2.SubnetArgs{
			VpcId:     vpc.ID(),
			CidrBlock: cidrBlockPublic,
			Tags:      tagsPublic,
		})
		if err != nil {
			return err
		}

		privateSubnet, err := ec2.NewSubnet(ctx, vpcPrivateName, &ec2.SubnetArgs{
			VpcId:     vpc.ID(),
			CidrBlock: cidrBlockPrivate,
			Tags:      tagsPrivate,
		})
		if err != nil {
			return err
		}

		// Internet Gateways
		publicInternetGateway, err := ec2.NewInternetGateway(ctx, fmt.Sprintf("%s-ig", vpcPublicName), &ec2.InternetGatewayArgs{
			VpcId: vpc.ID(),
			Tags:  tagsPublic,
		})
		if err != nil {
			return err
		}

		privateNatGateway, err := ec2.NewNatGateway(ctx, fmt.Sprintf("%s-nat-gw", vpcPrivateName), &ec2.NatGatewayArgs{
			ConnectivityType: pulumi.String("private"),
			SubnetId:         privateSubnet.ID(),
			Tags:             tagsPrivate,
		})
		if err != nil {
			return err
		}

		// Route Tables
		publicRouteTable, err := ec2.NewRouteTable(ctx, vpcPublicName, &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Tags:  tagsPublic,
		})
		if err != nil {
			return err
		}

		privateRouteTable, err := ec2.NewRouteTable(ctx, vpcPrivateName, &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Tags:  tagsPrivate,
		})
		if err != nil {
			return err
		}

		// Routes
		ec2.NewRoute(ctx, fmt.Sprintf("%s-route-public", vpcPublicName), &ec2.RouteArgs{
			RouteTableId:         publicRouteTable.ID(),
			DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
			GatewayId:            publicInternetGateway.ID(),
		})
		ec2.NewRoute(ctx, fmt.Sprintf("%s-route-private", vpcPrivateName), &ec2.RouteArgs{
			RouteTableId:         privateRouteTable.ID(),
			DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
			NatGatewayId:         privateNatGateway.ID(),
		})

		// Route Table Association
		ec2.NewRouteTableAssociation(ctx, vpcPublicName, &ec2.RouteTableAssociationArgs{
			RouteTableId: publicRouteTable.ID(),
			SubnetId:     publicSubnet.ID(),
		})
		ec2.NewRouteTableAssociation(ctx, vpcPrivateName, &ec2.RouteTableAssociationArgs{
			RouteTableId: privateRouteTable.ID(),
			SubnetId:     privateSubnet.ID(),
		})

		securityGroup, err := ec2.NewSecurityGroup(ctx, fmt.Sprintf("%s-allow-tls", vpcName), &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:    pulumi.String("tcp"),
					FromPort:    pulumi.Int(443),
					ToPort:      pulumi.Int(443),
					Description: pulumi.String("TLS from VPC"),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
					Ipv6CidrBlocks: pulumi.StringArray{
						pulumi.String("::/0"),
					},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					FromPort: pulumi.Int(0),
					ToPort:   pulumi.Int(0),
					Protocol: pulumi.String("-1"),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
					Ipv6CidrBlocks: pulumi.StringArray{
						pulumi.String("::/0"),
					},
				},
			},
			VpcId: vpc.ID(),
		})
		if err != nil {
			return err
		}

		ctx.Export("vpcId", vpc.ID())
		ctx.Export("securityGroupArn", securityGroup.Arn)

		return nil
	})
}
