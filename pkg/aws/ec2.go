package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (c *Client) GetInstanceID(nodeName string) (string, error) {
	svc := ec2.New(c.Session)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{aws.String(nodeName)},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		return "", err
	}

	var instanceID string
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			instanceID = *instance.InstanceId
			break
		}
		break
	}
	return instanceID, nil
}
