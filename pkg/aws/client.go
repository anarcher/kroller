package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Client struct {
	Session *session.Session
}

func NewClient(region string) (*Client, error) {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	c := Client{
		Session: s,
	}

	return &c, nil
}
