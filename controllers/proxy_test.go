package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/d.ims.io/mock"
	"testing"
	"github.com/quintilesims/d.ims.io/controllers/proxy"
	"net/http"
)

func TestProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testProxy := proxy.ProxyFunc(func(token string, w http.ResponseWriter, r *http.Request){
	 	if v, want := token, "token"; v != want {
                        t.Errorf("Token was '%v', expected '%v'", v, want)
                }
	})

	mockECR := mock.NewMockECRAPI(ctrl)
	controller := NewProxyController(mockECR, testProxy)

	authData := []*ecr.AuthorizationData{
		{AuthorizationToken: aws.String("token")},
	}

	mockECR.EXPECT().
		GetAuthorizationToken(gomock.Any()).
		Return(&ecr.GetAuthorizationTokenOutput{AuthorizationData: authData}, nil)

	c := generateContext(t, nil, nil)
	resp, err := controller.DoProxy(c)
	if err != nil {
		t.Fatal(err)
	}
	
	// run the test proxy
	resp.Write(nil, nil)
}
