package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
)

type spaceAPIValidatorFormat struct {
	Data string `json:"data"`
}

func TestSpaceApi(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/api/v2.0/SpaceAPI", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	spaceAPIRouteHandler(c)

	if http.StatusOK != rec.Code {
		t.Errorf("expected status OK, got %v", rec.Code)
	}

	bodyBytes, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Errorf("Couldn't read Space Api response. Error: %v", err)
	}

	var spaceDescriptor SpaceDescriptor
	err = json.Unmarshal(bodyBytes, &spaceDescriptor)
	if err != nil {
		t.Errorf("Responce from Space Api isn't couldn't be assigned to SpaceDescriptor. Error: %v, %v", err, &spaceDescriptor)
	}

	spaceDescriptorStringified, err := json.Marshal(&spaceDescriptor)
	if err != nil {
		t.Errorf("Got error while trying to convert spaceDescriptor to json . Error: %v", err)
	}

	requestData, err := json.Marshal(spaceAPIValidatorFormat{string(spaceDescriptorStringified)})
	if err != nil {
		t.Errorf("Got error while trying to create request for spaceDescriptor. Error: %v", err)
	}

	res, err := http.Post("https://validator.spacedirectory.org/v1/validate/", "application/json", bytes.NewBuffer(requestData))
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", res.StatusCode)
	}
	if err != nil {
		t.Errorf("Got error while trying to validate spaceDescriptor. Error: %v", err)
	}

}
