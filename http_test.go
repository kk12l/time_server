package main

import (
  "net/http"
  "net/http/httptest"
  "testing"
  "encoding/json"
  "time"
  "strconv"
)


func MakeRequest(t *testing.T, handler http.HandlerFunc, method string, req string) *httptest.ResponseRecorder {
  w := httptest.NewRecorder()
  r, err := http.NewRequest(method, req, nil)
  if err != nil {
    t.Fatal(err)
  }

  handler.ServeHTTP(w, r)
  return w
}

func TestNowHandler(t *testing.T) {
  var handler http.HandlerFunc = http.HandlerFunc(NowHandler)
  var request *httptest.ResponseRecorder

  request = MakeRequest(t, handler, http.MethodPost, "/time/now")
  if request.Code != http.StatusNotFound {
    t.Errorf("exp %d, got %d", http.StatusNotFound, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/now")
  if request.Code != http.StatusOK {
    t.Errorf("exp %d, got %d", http.StatusOK, request.Code)
  }

  var dat map[string]interface{}
  var err error
  var decoder *json.Decoder

  decoder = json.NewDecoder(request.Body)
  if err = decoder.Decode(&dat); err != nil {
    t.Fatal(err)
  }
  if dat["time"] == nil {
      t.Errorf("Expected time param")
  } else {
    resp_time := strconv.FormatFloat(dat["time"].(float64), 'f', 6, 64)
    if _, parsed_err := time.Parse("060102.150405", resp_time); parsed_err != nil {
      t.Errorf(parsed_err.Error())
    }
  }
}

func TestStringHandler(t *testing.T) {
  var handler http.HandlerFunc = http.HandlerFunc(StringHandler)
  var request *httptest.ResponseRecorder

  request = MakeRequest(t, handler, http.MethodPost, "/time/string?time=710822.010001")
  if request.Code != http.StatusNotFound {
    t.Errorf("exp %d, got %d", http.StatusNotFound, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/string?_time=710822.010001")
  if request.Code != http.StatusNotFound {
    t.Errorf("exp %d, got %d", http.StatusNotFound, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/string?time=822.010001499")
  if request.Code != http.StatusOK {
    t.Errorf("exp %d, got %d", http.StatusOK, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/string?time=710822.010001")
  if request.Code != http.StatusOK {
    t.Errorf("exp %d, got %d", http.StatusOK, request.Code)
  }

  var dat map[string]interface{}
  var err error
  var decoder *json.Decoder
  var resp_time_json string

  decoder = json.NewDecoder(request.Body)
  if err = decoder.Decode(&dat); err != nil {
    t.Fatal(err)
  }
  if dat["str"] == nil {
      t.Errorf("Expected str param")
  } else {
    resp_time_json = dat["str"].(string)
    if resp_time_json != "19710822010001" {
      t.Errorf("exp %s, got %s", "19710822010001", resp_time_json)
    }
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/string?time=711131.010001499")
  if request.Code != http.StatusNotAcceptable {
    t.Errorf("exp %d, got %d", http.StatusNotAcceptable, request.Code)
  }

  decoder = json.NewDecoder(request.Body)
  if err = decoder.Decode(&dat); err != nil {
    t.Fatal(err)
  }
  if dat["error"] == nil {
      t.Errorf("Expected error param")
  } else {
    resp_time_json = dat["error"].(string)
    if resp_time_json != `parsing time "711131.010001": day out of range` {
      t.Errorf("exp %s, got %s", `parsing time "711131.010001": day out of range`, resp_time_json)
    }
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/string?time=711130.240001")
  if request.Code != http.StatusNotAcceptable {
    t.Errorf("exp %d, got %d", http.StatusNotAcceptable, request.Code)
  }

  decoder = json.NewDecoder(request.Body)
  if err = decoder.Decode(&dat); err != nil {
    t.Fatal(err)
  }
  if dat["error"] == nil {
      t.Errorf("Expected error param")
  } else {
    resp_time_json = dat["error"].(string)
    if resp_time_json != `parsing time "711130.240001": hour out of range` {
      t.Errorf("exp %s, got %s", `parsing time "711130.240001": hour out of range`, resp_time_json)
    }
  }
}

func TestResetHandler(t *testing.T) {
  var handler http.HandlerFunc = http.HandlerFunc(ResetHandler)
  var handler_now http.HandlerFunc = http.HandlerFunc(NowHandler)
  var request *httptest.ResponseRecorder
  var request_now *httptest.ResponseRecorder

  request = MakeRequest(t, handler, http.MethodGet, "/time/reset")
  if request.Code != http.StatusNotFound {
    t.Errorf("exp %d, got %d", http.StatusNotFound, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodPost, "/time/reset")
  if request.Code != http.StatusOK {
    t.Errorf("exp %d, got %d", http.StatusOK, request.Code)
  }

  request_now = MakeRequest(t, handler_now, http.MethodGet, "/time/now")
  test_time := time_to_json(time.Now().In(time.UTC))
  if request_now.Code != http.StatusOK {
    t.Errorf("exp %d, got %d", http.StatusOK, request_now.Code)
  }

  var dat map[string]interface{}
  var err error
  var decoder *json.Decoder

  decoder = json.NewDecoder(request_now.Body)
  if err = decoder.Decode(&dat); err != nil {
    t.Fatal(err)
  }
  if dat["time"] == nil {
      t.Errorf("Expected time param")
  } else {
    resp_time := dat["time"].(float64)
    if resp_time != test_time {
      t.Errorf("exp %f, got %f", test_time, resp_time)
    }
  }
}

func TestSetHandler(t *testing.T) {
  var handler http.HandlerFunc = http.HandlerFunc(SetHandler)
  var handler_now http.HandlerFunc = http.HandlerFunc(NowHandler)
  var request *httptest.ResponseRecorder
  var request_now *httptest.ResponseRecorder

  request = MakeRequest(t, handler, http.MethodGet, "/time/set?time=710813.010000499")
  if request.Code != http.StatusNotFound {
    t.Errorf("exp %d, got %d", http.StatusNotFound, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodPost, "/time/set?_time=710813.010000499")
  if request.Code != http.StatusNotFound {
    t.Errorf("exp %d, got %d", http.StatusNotFound, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodPost, "/time/set?time=710813.010000499")
  if request.Code != http.StatusOK {
    t.Errorf("exp %d, got %d", http.StatusOK, request.Code)
  }

  request_now = MakeRequest(t, handler_now, http.MethodGet, "/time/now")
  test_time, _ := time_from_json(710813.010000499)
  if request_now.Code != http.StatusOK {
    t.Errorf("exp %d, got %d", http.StatusOK, request_now.Code)
  }

  var dat map[string]interface{}
  var err error
  var decoder *json.Decoder

  decoder = json.NewDecoder(request_now.Body)
  if err = decoder.Decode(&dat); err != nil {
    t.Fatal(err)
  }
  if dat["time"] == nil {
      t.Errorf("Expected time param")
  } else {
    resp_time, _ := time_from_json(dat["time"].(float64))
    if resp_time != test_time {
      t.Errorf("exp %s, got %s", test_time.Format("060102.150405"), resp_time.Format("060102.150405"))
    }
  }

  request = MakeRequest(t, handler, http.MethodPost, "/time/set?time=711131.010001")
  if request.Code != http.StatusNotAcceptable {
    t.Errorf("exp %d, got %d", http.StatusNotAcceptable, request.Code)
  }

  decoder = json.NewDecoder(request.Body)
  if err = decoder.Decode(&dat); err != nil {
    t.Fatal(err)
  }
  if dat["error"] == nil {
      t.Errorf("Expected error param")
  } else {
    resp_time_json := dat["error"].(string)
    if resp_time_json != `parsing time "711131.010001": day out of range` {
      t.Errorf("exp %s, got %s", `parsing time "711131.010001": day out of range`, resp_time_json)
    }
  }
}

func TestAddHandler(t *testing.T) {
  var handler http.HandlerFunc = http.HandlerFunc(AddHandler)
  var request *httptest.ResponseRecorder

  request = MakeRequest(t, handler, http.MethodPost, "/time/set?time=710813.010000499&delta=0.99")
  if request.Code != http.StatusNotFound {
    t.Errorf("exp %d, got %d", http.StatusNotFound, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/set?delta=0.99")
  if request.Code != http.StatusNotAcceptable {
    t.Errorf("exp %d, got %d", http.StatusNotAcceptable, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/set?time=710813.010000499")
  if request.Code != http.StatusNotAcceptable {
    t.Errorf("exp %d, got %d", http.StatusNotAcceptable, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/set?time=710813.010000499&delta=")
  if request.Code != http.StatusNotAcceptable {
    t.Errorf("exp %d, got %d", http.StatusNotAcceptable, request.Code)
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/set?time=711131.010001499&delta=0.99")
  if request.Code != http.StatusNotAcceptable {
    t.Errorf("exp %d, got %d", http.StatusNotAcceptable, request.Code)
  }

  var dat map[string]interface{}
  var err error
  var decoder *json.Decoder

  decoder = json.NewDecoder(request.Body)
  if err = decoder.Decode(&dat); err != nil {
    t.Fatal(err)
  }
  if dat["error"] == nil {
      t.Errorf("Expected error param")
  } else {
    resp_time_json := dat["error"].(string)
    if resp_time_json != `parsing time "711131.010001": day out of range` {
      t.Errorf("exp %s, got %s", `parsing time "711131.010001": day out of range`, resp_time_json)
    }
  }

  request = MakeRequest(t, handler, http.MethodGet, "/time/set?time=710813.010000499&delta=10101.999999")
  test_time, _ := time_from_json(720916.054039499)
  if request.Code != http.StatusOK {
    t.Errorf("exp %d, got %d", http.StatusOK, request.Code)
  }

  decoder = json.NewDecoder(request.Body)
  if err = decoder.Decode(&dat); err != nil {
    t.Fatal(err)
  }
  if dat["time"] == nil {
      t.Errorf("Expected time param")
  } else {
    resp_time, _ := time_from_json(dat["time"].(float64))
    if resp_time != test_time {
      t.Errorf("exp %s, got %s", test_time.Format("060102.150405"), resp_time.Format("060102.150405"))
    }
  }
}

