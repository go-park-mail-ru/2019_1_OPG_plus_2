package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"

	a "github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/adapters"
	"github.com/go-park-mail-ru/2019_1_OPG_plus_2/internal/pkg/models"
)

var baseUrl = "localhost:8002/api"

type TestParams struct {
	isAuth  bool
	muxVars map[string]string
	jwt     models.JwtData
	method  string
	url     string
}

type TestCase struct {
	handler      http.HandlerFunc
	params       TestParams
	inputMessage json.RawMessage
	expStatus    int
	expMessage   interface{}
}

func testInitial() {
	a.SetStorages(a.NewMockCoreStorage(), a.NewMockAuthStorage())
	a.SetHandlers(NewUserHandlers(), NewAuthHandlers())
}

func testCaseInitial(tCase TestCase) (*httptest.ResponseRecorder, *http.Request) {
	testParams := tCase.params
	url := baseUrl + testParams.url
	req := httptest.NewRequest(testParams.method, url, bytes.NewReader(tCase.inputMessage))
	w := httptest.NewRecorder()
	ctx := req.Context()

	data := testParams.jwt
	ctx = context.WithValue(ctx, "isAuth", testParams.isAuth)
	ctx = context.WithValue(ctx, "jwtData", data)

	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, tCase.params.muxVars)

	return w, req
}
