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
