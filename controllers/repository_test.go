package controllers

import (
	"strings"
	"testing"
	"time"

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

func TestGetImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewRepositoryController(mockECR)

	validateDescribeImagesInput := func(input *ecr.DescribeImagesInput) {
		if v, want := aws.StringValue(input.RepositoryName), "user/test"; v != want {
			t.Errorf("Name was '%v', expected '%v'", v, want)
		}
		if v, want := aws.StringValue(input.ImageIds[0].ImageTag), "test"; v != want {
			t.Errorf("Tag was '%v', expected '%v'", v, want)
		}
	}

	tag := ""

	tags := []*string{}
	tags = append(tags, &tag)

	detail := &ecr.ImageDetail{}
	detail.SetImageDigest("")
	detail.SetImagePushedAt(time.Time{})
	detail.SetImageSizeInBytes(0)
	detail.SetImageTags(tags)
	detail.SetRepositoryName("")

	output := &ecr.DescribeImagesOutput{}
	output.SetImageDetails([]*ecr.ImageDetail{detail})

	mockECR.EXPECT().
		DescribeImages(gomock.Any()).
		Do(validateDescribeImagesInput).
		Return(output, nil)

	c := generateContext(t, nil, map[string]string{"tag": "test", "name": "test", "owner": "user"})
	if _, err := controller.GetImage(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewRepositoryController(mockECR)

	validateBatchDeleteImageInput := func(input *ecr.BatchDeleteImageInput) {
		if v, want := aws.StringValue(input.RepositoryName), "user/test"; v != want {
			t.Errorf("Name was '%v', expected '%v'", v, want)
		}

		if v, want := aws.StringValue(input.ImageIds[0].ImageTag), "test"; v != want {
			t.Errorf("Tag was '%v', expected '%v'", v, want)
		}
	}

	output := &ecr.BatchDeleteImageOutput{}

	mockECR.EXPECT().
		BatchDeleteImage(gomock.Any()).
		Do(validateBatchDeleteImageInput).
		Return(output, nil)

	c := generateContext(t, nil, map[string]string{"tag": "test", "name": "test", "owner": "user"})
	if _, err := controller.DeleteImage(c); err != nil {
		t.Fatal(err)
	}
}

func TestListImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewRepositoryController(mockECR)

	numRepos := 3

	populateRepos := func(input *ecr.DescribeRepositoriesInput, fn func(output *ecr.DescribeRepositoriesOutput, lastPage bool) bool) error {
		repo := &ecr.Repository{}
		repo.SetRepositoryName("test")

		repos := []*ecr.Repository{}
		repos = append(repos, repo)

		output := &ecr.DescribeRepositoriesOutput{}
		output.SetRepositories(repos)

		for i := 0; i < numRepos; i++ {
			fn(output, true)
		}

		return nil
	}

	mockECR.EXPECT().
		DescribeRepositoriesPages(gomock.Any(), gomock.Any()).
		Do(populateRepos).
		Return(nil)

	for i := 0; i < numRepos; i++ {
		mockECR.EXPECT().
			DescribeImagesPages(gomock.Any(), gomock.Any()).
			Return(nil)
	}

	c := generateContext(t, nil, nil)
	if _, err := controller.ListImages(c); err != nil {
		t.Fatal(err)
	}
}
