package controllers

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/quintilesims/d.ims.io/mock"
)

func TestGrantAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	mockAccountManager := mock.NewMockAccountManager(ctrl)
	controller := NewAccountController(mockECR, mockAccountManager)

	mockECR.EXPECT().
		DescribeRepositoriesPages(gomock.Any(), gomock.Any()).
		Return(nil)

	mockAccountManager.EXPECT().
		Accounts().
		Return([]string{}, nil)

	mockAccountManager.EXPECT().
		GrantAccess(gomock.Any()).
		Return(nil)

	c := generateContext(t, nil, nil)
	if _, err := controller.GrantAccess(c); err != nil {
		t.Fatal(err)
	}
}

func TestRevokeAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockECR := mock.NewMockECRAPI(ctrl)
	mockAccountManager := mock.NewMockAccountManager(ctrl)
	controller := NewAccountController(mockECR, mockAccountManager)

	mockECR.EXPECT().
		DescribeRepositoriesPages(gomock.Any(), gomock.Any()).
		Return(nil)

	mockAccountManager.EXPECT().
		RevokeAccess(gomock.Any()).
		Return(nil)

	c := generateContext(t, nil, nil)
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
