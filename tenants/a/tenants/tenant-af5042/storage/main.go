package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		tenantStackID := cfg.Require("tenantStack")
		// tenantStack := pulumi.StackReference(tenantStackID)
		tenantStack, err := pulumi.NewStackReference(ctx, tenantStackID, nil)
		if err != nil {
			panic(err)
		}

		tenant := tenantStack.GetOutput(pulumi.String("tenant"))

		bucketName := pulumi.Sprintf("tenant-%s-storage-1", tenant)

		fmt.Println(bucketName)

		panic(bucketName)
		// bucketArn, err := arn.Parse(bucketName)
		// if err != nil {
		// 	fmt.Sprintf("Cannot validate %v", bucketArn.String())
		// 	panic(err)
		// }

		// Create an AWS resource (S3 Bucket)
		// bucket, err := s3.NewBucket(ctx, bucketName, nil)
		// if err != nil {
		// 	return err
		// }

		// // Export the name of the bucket
		// ctx.Export("bucketName", bucket.ID())
		return nil
	})
}
