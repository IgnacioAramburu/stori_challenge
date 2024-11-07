package cdk

import (
	"log"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/jsii-runtime-go"
)

type LambdaStackProps struct {
	awscdk.StackProps
}

func NewLambdaStack(scope awscdk.Construct, id string, props *LambdaStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Lambda Layer
	layer := awslambda.NewLayerVersion(stack, jsii.String("layer"), &awslambda.LayerVersionProps{
		Code: awslambda.Code_FromAsset(jsii.String("layer.zip"), nil),
		CompatibleRuntimes: []awslambda.Runtime{
			awslambda.Runtime_GO_1_X(),
		},
	})

	// Lambda 1
	lambda1 := awslambda.NewFunction(stack, jsii.String("lbd_generate_data"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Handler: jsii.String("cmd/lbd_generate_data.HandleRequest"), // Assuming handler is in handler.go
		Code:    awslambda.Code_FromAsset(jsii.String("cmd/lbd_generate_data"), nil),
		Environment: map[string]*string{
			"LAYER_ARN": layer.LayerVersionArn(),
		},
		Layers: []awslambda.ILayerVersion{layer},
	})

	// Lambda 2
	lambda2 := awslambda.NewFunction(stack, jsii.String("lbd_send_summary_mail"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Handler: jsii.String("cmd/lbd_send_summary_mail.HandleRequest"), // Assuming handler is in handler.go
		Code:    awslambda.Code_FromAsset(jsii.String("cmd/lbd_send_summary_mail"), nil),
		Environment: map[string]*string{
			"LAYER_ARN": layer.LayerVersionArn(),
		},
		Layers: []awslambda.ILayerVersion{layer},
	})

	// Add IAM policies if necessary
	lambda1.Role().AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaBasicExecutionRole")))
	lambda2.Role().AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaBasicExecutionRole")))

	// Optional error handling
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error creating Lambda stack:", err)
		}
	}()

	return stack
}
