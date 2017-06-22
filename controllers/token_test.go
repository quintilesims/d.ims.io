package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/d.ims.io/mock"
)

func TestCreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenManager := mock.NewMockTokenManager(ctrl)
	controller := NewTokenController(mockTokenManager)

	mockTokenManager.EXPECT().
		CreateToken(gomock.Any()).
		Return("", nil)

	c := generateContext(t, nil, nil)
	if _, err := controller.CreateToken(c); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenManager := mock.NewMockTokenManager(ctrl)
	controller := NewTokenController(mockTokenManager)

	mockTokenManager.EXPECT().
		DeleteToken("test").
		Return(nil)

	c := generateContext(t, nil, map[string]string{"token": "test"})
	if _, err := controller.DeleteToken(c); err != nil {
		t.Fatal(err)
	}
}

/*
func TestDeleteToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewTokenController(mockECR)

	validateDeleteTokenInput := func(input *ecr.DeleteTokenInput) {
		if v, want := aws.StringValue(input.TokenName), "test"; v != want {
			t.Errorf("Name was '%v', expected '%v'", v, want)
		}

		if v, want := aws.BoolValue(input.Force), true; v != want {
			t.Errorf("Force was '%v', expected '%v'", v, want)
		}
	}

	mockECR.EXPECT().
		DeleteToken(gomock.Any()).
		Do(validateDeleteTokenInput).
		Return(&ecr.DeleteTokenOutput{}, nil)

	c := generateContext(t, nil, map[string]string{"name": "test"})
	if _, err := controller.DeleteToken(c); err != nil {
		t.Fatal(err)
	}
}
*/
