// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var SQSAllowEmptyValues = []string{"tags."}

type SQSGenerator struct {
	AWSService
}

// Generate TerraformResources from AWS API,
func (g *SQSGenerator) InitResources() error {
	var attrOutput string

	sess := g.generateSession()
	svc := sqs.New(sess)
	g.Resources = []terraform_utils.Resource{}

	queues, err := svc.ListQueues(&sqs.ListQueuesInput{})
	if err != nil {
		return err
	}

	// TODO: Pull this out to a function.
	for _, queueURL := range queues.QueueUrls {

		attrOutput, err := svc.GetQueueAttributes(&sqs.GetQueueAttributesInput{
			QueueUrl: queueURL,
		})
		if err != nil {
			return err
		}

		attributes := make(map[string]string)

		for k, v := range attrOutput.Attributes {
			attributes[k] = aws.String(v)
		}

		resource := terraform_utils.NewResource(
			queueURL,
			queueURL,
			"aws_sqs",
			"aws",
			attributes,
			SQSAllowEmptyValues,
			map[string]string{},
		)
	}

	// g.Resources = g.createResources(queues)
	g.PopulateIgnoreKeys()
	return nil

}
