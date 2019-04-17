package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"

	a "2019_1_OPG_plus_2/internal/pkg/adapters"
	"2019_1_OPG_plus_2/internal/pkg/models"
)

var baseUrl = "localhost:8002/api"

type testParams struct {
	isAuth  bool
	muxVars map[string]string
	jwt     models.JwtData
	method  string
	url     string
}

type testCase struct {
	handler       http.HandlerFunc
	params        testParams
	inputMessage  json.RawMessage
	outputMessage interface{}
	expStatus     int
	expMessage    interface{}
}

func testsInitial() {
	a.SetStorages(a.NewMockCoreStorage(), a.NewMockAuthStorage())
	a.SetHandlers(NewUserHandlers(), NewAuthHandlers(), NewVkAuthHandlers())
}

func testCaseInitial(tCase *testCase) (*httptest.ResponseRecorder, *http.Request) {
	params := tCase.params
	url := baseUrl + params.url
	r := httptest.NewRequest(params.method, url, bytes.NewReader(tCase.inputMessage))
	w := httptest.NewRecorder()
	ctx := r.Context()

	data := params.jwt
	ctx = context.WithValue(ctx, "isAuth", params.isAuth)
	ctx = context.WithValue(ctx, "jwtData", data)

	r = r.WithContext(ctx)

	r = mux.SetURLVars(r, tCase.params.muxVars)

	return w, r
}

func test(t *testing.T, tCase *testCase) (*httptest.ResponseRecorder, *http.Request) {
	w, r := testCaseInitial(tCase)
	tCase.handler(w, r)

	if w.Code != tCase.expStatus {
		t.Errorf("Wrong Status:\n\tGot: %d\n\tExpected: %d\n", w.Code, tCase.expStatus)
	}

	_ = json.NewDecoder(w.Body).Decode(tCase.outputMessage)

	if !reflect.DeepEqual(tCase.outputMessage, tCase.expMessage) {
		t.Errorf("Wrong Body:\n\tGot: %v\n\tExpected: %v\n", tCase.outputMessage, tCase.expMessage)
	}

	return w, r
}
