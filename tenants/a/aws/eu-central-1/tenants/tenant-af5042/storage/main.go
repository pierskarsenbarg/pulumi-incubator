package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
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

		// bucket := tenant.ApplyT(func(v string) (pulumi.IDOutput, error) {
		// 	bucketName := fmt.Sprintf("tenant-%s-bucket-1", v)
		// 	bucket, err := s3.NewBucket(ctx, bucketName, nil)
		// 	if err != nil {
		// 		return pulumi.IDOutput{}, err
		// 	}
		// 	return bucket.ID(), nil

		// })
		bucket := pulumi.All(tenant, environment).ApplyT(func(args []interface{}) (pulumi.IDOutput, error) {
			bucketName := fmt.Sprintf("tenant-%s-%v-bucket-1", args[0], args[1])
			bucket, err := s3.NewBucket(ctx, bucketName, nil)
			if err != nil {
				return pulumi.IDOutput{}, err
			}
			return bucket.ID(), nil

		})

		ctx.Export("bucketName", bucket)

		return nil
	})
}
