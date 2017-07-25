package controllers

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/d.ims.io/mock"
	"github.com/quintilesims/d.ims.io/models"
)

func TestCreateRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewRepositoryController(mockECR)

	validateCreateRepositoryInput := func(input *ecr.CreateRepositoryInput) {
		if v, want := aws.StringValue(input.RepositoryName), "user/test"; v != want {
			t.Errorf("Name was '%v', expected '%v'", v, want)
		}
	}

	mockECR.EXPECT().
		CreateRepository(gomock.Any()).
		Do(validateCreateRepositoryInput).
		Return(&ecr.CreateRepositoryOutput{}, nil)

	c := generateContext(t, models.CreateRepositoryRequest{Name: "test"}, map[string]string{"owner": "user"})
	if _, err := controller.CreateRepository(c); err != nil {
		t.Fatal(err)
	}
}

func TestCreateRepositoryFailsWithSlashes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewRepositoryController(mockECR)

	c := generateContext(t, models.CreateRepositoryRequest{Name: "slash/test"}, map[string]string{"owner": "user"})
	_, err := controller.CreateRepository(c)
	if !strings.ContainsAny(err.Error(), "cannot contain '/' characters") {
		t.Fatal(err)
	}
}

func TestDeleteRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewRepositoryController(mockECR)

	validateDeleteRepositoryInput := func(input *ecr.DeleteRepositoryInput) {
		if v, want := aws.StringValue(input.RepositoryName), "user/test"; v != want {
			t.Errorf("Name was '%v', expected '%v'", v, want)
		}

		if v, want := aws.BoolValue(input.Force), true; v != want {
			t.Errorf("Force was '%v', expected '%v'", v, want)
		}
	}

	mockECR.EXPECT().
		DeleteRepository(gomock.Any()).
		Do(validateDeleteRepositoryInput).
		Return(&ecr.DeleteRepositoryOutput{}, nil)

	c := generateContext(t, nil, map[string]string{"name": "test", "owner": "user"})
	if _, err := controller.DeleteRepository(c); err != nil {
		t.Fatal(err)
	}
}

func TestGetRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewRepositoryController(mockECR)

	validateDescribeImagesInput := func(input *ecr.DescribeImagesInput, fn func(*ecr.DescribeImagesOutput, bool) bool) {
		if v, want := aws.StringValue(input.RepositoryName), "user/test"; v != want {
			t.Errorf("Name was '%v', expected '%v'", v, want)
		}
	}

	mockECR.EXPECT().
		DescribeImagesPages(gomock.Any(), gomock.Any()).
		Do(validateDescribeImagesInput).
		Return(nil)

	c := generateContext(t, nil, map[string]string{"name": "test", "owner": "user"})
	if _, err := controller.GetRepository(c); err != nil {
		t.Fatal(err)
	}
}

func TestListRepositories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewRepositoryController(mockECR)

	validateDescribeRepositoriesInput := func(input *ecr.DescribeRepositoriesInput, fn func(*ecr.DescribeRepositoriesOutput, bool) bool) {
		// no fields to validate
	}

	mockECR.EXPECT().
		DescribeRepositoriesPages(gomock.Any(), gomock.Any()).
		Do(validateDescribeRepositoriesInput).
		Return(nil)

	c := generateContext(t, nil, nil)
	if _, err := controller.ListRepositories(c); err != nil {
		t.Fatal(err)
	}
}
