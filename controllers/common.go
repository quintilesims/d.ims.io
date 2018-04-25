package controllers

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/quintilesims/d.ims.io/models"
	"github.com/zpatrick/fireball"
)

func listRepositories(e ecriface.ECRAPI) ([]string, error) {
	repositories := []string{}

	fn := func(output *ecr.DescribeRepositoriesOutput, lastPage bool) bool {
		for _, repository := range output.Repositories {
			repositories = append(repositories, aws.StringValue(repository.RepositoryName))
		}

		return !lastPage
	}

	if err := e.DescribeRepositoriesPages(&ecr.DescribeRepositoriesInput{}, fn); err != nil {
		return nil, err
	}

	return repositories, nil
}

func getRepositoryPolicy(e ecriface.ECRAPI, repositoryName string) (*models.PolicyDocument, error) {
	policyDoc := &models.PolicyDocument{}

	getPolicyInput := &ecr.GetRepositoryPolicyInput{}
	getPolicyInput.SetRepositoryName(repositoryName)
	if err := getPolicyInput.Validate(); err != nil {
		return nil, err
	}

	getPolicyOutput, err := e.GetRepositoryPolicy(getPolicyInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() != ecr.ErrCodeRepositoryPolicyNotFoundException {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if getPolicyOutput.PolicyText != nil && aws.StringValue(getPolicyOutput.PolicyText) != "" {
		if err := json.Unmarshal([]byte(aws.StringValue(getPolicyOutput.PolicyText)), policyDoc); err != nil {
			return nil, err
		}
	}

	return policyDoc, nil
}

func setRepositoryPolicy(e ecriface.ECRAPI, repositoryName, policyText string) error {
	policyInput := &ecr.SetRepositoryPolicyInput{}
	policyInput.SetPolicyText(policyText)
	policyInput.SetRepositoryName(repositoryName)
	if err := policyInput.Validate(); err != nil {
		return fireball.NewError(400, err, nil)
	}

	if _, err := e.SetRepositoryPolicy(policyInput); err != nil {
		return err
	}

	return nil
}

func removeFromRepositoryPolicy(e ecriface.ECRAPI, repositoryName string, accountID string) error {
	log.Printf("[DEBUG] removing account: '%s' from policy for repository: '%s'", accountID, repositoryName)
	policyDoc, err := getRepositoryPolicy(e, repositoryName)
	if err != nil {
		return err
	}

	statementCount := len(policyDoc.Statement)
	policyDoc.RemoveAWSAccountPrincipal(accountID)

	// check if there was a change in the policy document
	if len(policyDoc.Statement) == statementCount {
		return nil
	}

	return setRepositoryPolicy(e, repositoryName, policyDoc.RenderPolicyText())
}

func addToRepositoryPolicy(e ecriface.ECRAPI, repositoryName string, accounts []string) error {
	log.Printf("[DEBUG] adding accounts: '%v' from policy for repository: '%s'", accounts, repositoryName)
	if len(accounts) == 0 {
		return nil
	}

	policyDoc, err := getRepositoryPolicy(e, repositoryName)
	if err != nil {
		return err
	}

	statementCount := len(policyDoc.Statement)
	for _, accountID := range accounts {
		policyDoc.AddAWSAccountPrincipal(accountID)
	}

	// check if there was change in the policy document
	if len(policyDoc.Statement) == statementCount {
		return nil
	}

	return setRepositoryPolicy(e, repositoryName, policyDoc.RenderPolicyText())
}
