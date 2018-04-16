package models

import (
	"encoding/json"
	"fmt"
)

type PolicyDocument struct {
	Version   string
	Statement []statement
}

type statement struct {
	Sid       string
	Effect    string
	Principal map[string]string
	Action    []string
}

func NewPolicyDocument() PolicyDocument {
	return PolicyDocument{
		Version:   "2008-10-17",
		Statement: []statement{},
	}
}

func newStatement(accountID string) statement {
	return statement{
		Sid:    accountID,
		Effect: "allow",
		Principal: map[string]string{
			"AWS": fmt.Sprintf("arn:aws:iam::%s:root", accountID),
		},
		Action: []string{
			"ecr:GetDownloadUrlForLayer",
			"ecr:BatchGetImage",
			"ecr:BatchCheckLayerAvailability",
		},
	}
}

func (p *PolicyDocument) AddAWSAccountPrincipal(accountID string) {
	permissionExists := false
	for _, s := range p.Statement {
		if s.Sid == accountID {
			permissionExists = true
			break
		}
	}

	if permissionExists {
		return
	}

	p.Statement = append(p.Statement, newStatement(accountID))
}

func (p *PolicyDocument) RenderPolicyText() string {
	d, err := json.Marshal(p)
	if err != nil {
		//not expected to ever reach this code as PolicyDocument struct can always be marshaled
		return ""
	}

	return string(d)
}
