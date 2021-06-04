package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aws/atest"
)

func TestAccAWSSQSQueuePolicy_basic(t *testing.T) {
	var queueAttributes map[string]*string
	resourceName := "aws_sqs_queue_policy.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { atest.PreCheck(t) },
		ErrorCheck:   atest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    atest.Providers,
		CheckDestroy: testAccCheckAWSSQSQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSQSPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSQSQueueExists("aws_sqs_queue.test", &queueAttributes),
					testAccCheckAWSSQSQueueDefaultAttributes(&queueAttributes),
					resource.TestMatchResourceAttr("aws_sqs_queue_policy.test", "policy",
						regexp.MustCompile("^{\"Version\":\"2012-10-17\".+")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccAWSSQSPolicyConfigBasic(rName),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "policy", "aws_sqs_queue.test", "policy"),
				),
			},
		},
	})
}

func TestAccAWSSQSQueuePolicy_disappears_queue(t *testing.T) {
	var queueAttributes map[string]*string
	resourceName := "aws_sqs_queue_policy.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { atest.PreCheck(t) },
		ErrorCheck:   atest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    atest.Providers,
		CheckDestroy: testAccCheckAWSSQSQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSQSPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSQSQueueExists("aws_sqs_queue.test", &queueAttributes),
					testAccCheckAWSSQSQueueDefaultAttributes(&queueAttributes),
					atest.CheckDisappears(atest.Provider, resourceAwsSqsQueue(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSSQSQueuePolicy_disappears(t *testing.T) {
	var queueAttributes map[string]*string
	resourceName := "aws_sqs_queue_policy.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { atest.PreCheck(t) },
		ErrorCheck:   atest.ErrorCheck(t, sqs.EndpointsID),
		Providers:    atest.Providers,
		CheckDestroy: testAccCheckAWSSQSQueueDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSQSPolicyConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSQSQueueExists("aws_sqs_queue.test", &queueAttributes),
					testAccCheckAWSSQSQueueDefaultAttributes(&queueAttributes),
					atest.CheckDisappears(atest.Provider, resourceAwsSqsQueuePolicy(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAWSSQSPolicyConfigBasic(rName string) string {
	return fmt.Sprintf(`
resource "aws_sqs_queue" "test" {
  name = %[1]q
}

resource "aws_sqs_queue_policy" "test" {
  queue_url = aws_sqs_queue.test.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "sqspolicy",
  "Statement": [
    {
      "Sid": "First",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "${aws_sqs_queue.test.arn}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_sqs_queue.test.arn}"
        }
      }
    }
  ]
}
POLICY
}
`, rName)
}
