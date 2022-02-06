# tenants

The goal for the *tenants* incubator is the design a way to get infrastructure for multi tenant that also can use shared infrastructure through. *Multi-Tenant through the IaC deployemt pipeline*

## examples

| name | guid | sha256 | hash |
| ---- | ---- | ------ | ---- |
| tenant-a | 4783eb5d-956e-4b58-a214-790fe7ef61ba | 71746f6c5d1c812640283718d5af5042 | af5042 |
| tenant-b | a79a8c37-9d0b-4b1d-a051-ba339eea6f21 | dc634b2c7ade59c80fa0df5771f09faf | f09faf |
| tenant-c | a869f326-7d83-4d1f-8af1-aef8fb11e57f | df3ebcc2b38304a9f663cefa49430e00 | 430e00 |

# options

## Option A

The goal with test **A** is to have a multiple tenants with setup shared infrastructure.

* The shared infra structure is project and has multiple stacks the output references to use in other stacks
* The tenant is a project with multiple stacks with one main called `meta` that the child stack uses


## Option B

*TBD, but the idea is it the same as **A** but it uses more component(s) instead of many stacks, try to reuse the code.*

## Option C

*TBD*


# Resources

* https://pkg.go.dev/github.com/pulumi/pulumi/sdk/v3@v3.24.1/go/pulumi
* https://pkg.go.dev/github.com/pulumi/pulumi/sdk/v3@v3.24.1/go/pulumi/config
* https://github.com/pulumi/examples/blob/master/aws-py-stackreference
* https://github.com/pulumi/examples/blob/master/aws-ts-stackreference
* https://github.com/pulumi/examples/tree/master/aws-ts-stackreference-architecture
* https://www.pulumi.com/docs/intro/concepts/stack/
* https://www.pulumi.com/docs/intro/concepts/config/
* https://www.pulumi.com/docs/intro/concepts/inputs-outputs/#convert-input-to-output-through-interpolation
* https://pkg.go.dev/github.com/pulumi/pulumi/sdk/v3/go/pulumi
* https://www.pulumi.com/learn/building-with-pulumi/stack-references/
* https://golangexample.com/a-tool-to-generate-pulumi-package-schemas-from-go-type-definitions/
* https://dev.to/aws-builders/infrastructure-as-code-on-aws-using-go-and-pulumi-gn5
* https://github.com/pulumi/pulumi/issues/3942
* https://www.leebriggs.co.uk/blog/2021/05/09/pulumi-apply.html
