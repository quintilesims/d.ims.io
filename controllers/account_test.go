package controllers

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/d.ims.io/mock"
	"github.com/quintilesims/d.ims.io/models"
)

func TestGrantAccessInputValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	mockAccountManager := mock.NewMockAccountManager(ctrl)
	controller := NewAccountController(mockECR, mockAccountManager)

	c := generateContext(t, models.GrantAccessRequest{Account: ""}, nil)
	if _, err := controller.GrantAccess(c); err == nil {
		t.Fatal("expected error when GrantAccessRequest.Account is empty")
	}
}

func TestGrantAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	mockAccountManager := mock.NewMockAccountManager(ctrl)
	controller := NewAccountController(mockECR, mockAccountManager)

	fnListRepos := func(input *ecr.DescribeRepositoriesInput, fn func(output *ecr.DescribeRepositoriesOutput, lastPage bool) bool) error {
		output := &ecr.DescribeRepositoriesOutput{
			Repositories: []*ecr.Repository{
				{RepositoryName: aws.String("user/name-*")},
				{RepositoryName: aws.String("user/name-*")},
			},
		}

		fn(output, false)
		return nil
	}

	mockECR.EXPECT().
		DescribeRepositoriesPages(gomock.Any(), gomock.Any()).
		Do(fnListRepos).
		Return(nil)

	mockAccountManager.EXPECT().
		Accounts().
		Return([]string{"account-id"}, nil)

	getPolicyInput := &ecr.GetRepositoryPolicyInput{}
	getPolicyInput.SetRepositoryName("user/name-*")
	mockECR.EXPECT().
		GetRepositoryPolicy(getPolicyInput).
		Return(&ecr.GetRepositoryPolicyOutput{}, nil).
		Times(2)

	policyDoc := models.PolicyDocument{}
	policyDoc.AddAWSAccountPrincipal("account-id")
	setPolicyInput := &ecr.SetRepositoryPolicyInput{}
	setPolicyInput.SetRepositoryName("user/name-*")
	setPolicyInput.SetPolicyText(policyDoc.RenderPolicyText())
	mockECR.EXPECT().
		SetRepositoryPolicy(setPolicyInput).
		Return(&ecr.SetRepositoryPolicyOutput{}, nil).
		Times(2)

	mockAccountManager.EXPECT().
		GrantAccess(gomock.Any()).
		Return(nil)

	c := generateContext(t, models.GrantAccessRequest{Account: "account-id"}, nil)
	if _, err := controller.GrantAccess(c); err != nil {
		t.Fatal(err)
	}
}

func TestRevokeAccessInputValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	mockAccountManager := mock.NewMockAccountManager(ctrl)
	controller := NewAccountController(mockECR, mockAccountManager)

	c := generateContext(t, nil, nil)
	if _, err := controller.RevokeAccess(c); err == nil {
		t.Fatal("expected error when id account id is not specified")
	}
}

func TestRevokeAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	mockAccountManager := mock.NewMockAccountManager(ctrl)
	controller := NewAccountController(mockECR, mockAccountManager)

	fnListRepos := func(input *ecr.DescribeRepositoriesInput, fn func(output *ecr.DescribeRepositoriesOutput, lastPage bool) bool) error {
		output := &ecr.DescribeRepositoriesOutput{
			Repositories: []*ecr.Repository{
				{RepositoryName: aws.String("user/name-1")},
			},
		}

		fn(output, false)
		return nil
	}

	mockECR.EXPECT().
		DescribeRepositoriesPages(gomock.Any(), gomock.Any()).
		Do(fnListRepos).
		Return(nil)

	policyDoc := &models.PolicyDocument{}
	policyDoc.AddAWSAccountPrincipal("account-id")

	getPolicyOutput := &ecr.GetRepositoryPolicyOutput{}
	getPolicyOutput.SetPolicyText(policyDoc.RenderPolicyText())

	getPolicyInput := &ecr.GetRepositoryPolicyInput{}
	getPolicyInput.SetRepositoryName("user/name-1")
	mockECR.EXPECT().
		GetRepositoryPolicy(getPolicyInput).
		Return(getPolicyOutput, nil)

	setPolicyInput := &ecr.SetRepositoryPolicyInput{}
	setPolicyInput.SetRepositoryName("user/name-1")
	setPolicyInput.SetPolicyText("")
	mockECR.EXPECT().
		SetRepositoryPolicy(setPolicyInput).
		Return(&ecr.SetRepositoryPolicyOutput{}, nil)

	mockAccountManager.EXPECT().
		RevokeAccess(gomock.Any()).
		Return(nil)

	c := generateContext(t, nil, map[string]string{"id": "account-id"})
	if _, err := controller.RevokeAccess(c); err != nil {
		t.Fatal(err)
	}
}

func TestAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	mockAccountManager := mock.NewMockAccountManager(ctrl)
	controller := NewAccountController(mockECR, mockAccountManager)

	mockAccountManager.EXPECT().
		Accounts().
		Return([]string{}, nil)

	c := generateContext(t, nil, nil)
	if _, err := controller.ListAccounts(c); err != nil {
		t.Fatal(err)
	}
}
