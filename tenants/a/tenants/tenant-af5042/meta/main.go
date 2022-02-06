package main

import (
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		projectName := ctx.Project()

		cfg := config.New(ctx, "")
		tenant := cfg.Require("tenant")
		tenantID := cfg.Require("tenantId")
		region := cfg.Require("region")
		size := cfg.Require("size")
		market := cfg.Require("market")

		// ctx.Export("info", Tenant{
		// 	ProjectName: pulumi.String(projectName),
		// 	TenantID:    pulumi.String(tenantID),
		// })

		ctx.Export("projectName", pulumi.String(projectName))
		ctx.Export("tenant", pulumi.String(tenant))
		ctx.Export("tenantId", pulumi.String(tenantID))
		ctx.Export("size", pulumi.String(size))
		ctx.Export("market", pulumi.String(market))
		ctx.Export("region", pulumi.String(region))

		return nil
	})
}

type Tenant struct {
	ProjectName pulumi.StringInput `pulumi:"projectName"`
	Tenant      pulumi.StringInput `pulumi:"tenant"`
	TenantID    pulumi.StringInput `pulumi:"tenantId"`
	Name        pulumi.StringInput `pulumi:"name"`
	Size        pulumi.StringInput `pulumi:"size"`
	Market      pulumi.StringInput `pulumi:"market"`
	Region      pulumi.StringInput `pulumi:"region"`
}

func (Tenant) ElementType() reflect.Type {
	return nestedType
}

type TenantNested struct {
	ProjectName string `pulumi:"projectName"`
	Tenant      string `pulumi:"tenant"`
	TenantID    string `pulumi:"tenantId"`
	Name        string `pulumi:"name"`
	Size        string `pulumi:"size"`
	Market      string `pulumi:"market"`
	Region      string `pulumi:"region"`
}

var nestedType = reflect.TypeOf((*TenantNested)(nil)).Elem()
