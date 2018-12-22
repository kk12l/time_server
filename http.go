package main

import (
  "net/http"
  "strings"
  "log"
  "time"
  "encoding/json"
  "strconv"
  "fmt"
)

var time_diff time.Duration

func time_from_json(f float64) (time.Time, error) {
  s := fmt.Sprintf("%013.6f", f)
  time_out, err := time.Parse("060102.150405", s)
  nanosecond := 1e6 * time.Duration(int64(f * 1e9) % 1e3)
  time_out.Add(nanosecond)
  return time_out, err
}

func time_to_json(t time.Time) float64 {
  time_out, _ := strconv.ParseFloat(t.Format("060102.150405"), 64)
  return time_out
}

func NowHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    w.WriteHeader(404)
    return
  }
  resp, _ := json.Marshal(map[string]float64{"time": time_to_json(time.Now().Add(time_diff).In(time.UTC))})
  w.Header().Set("Content-Type", "application/json")
  w.Write(resp)
}

func StringHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    w.WriteHeader(404)
    return
  }
  r.ParseForm()
  for k, v := range r.Form {
    if k == "time" {
      var resp []byte
      json_time, _ := strconv.ParseFloat(strings.Join(v, ""), 9)
      if parsed_time, parsed_err := time_from_json(json_time); parsed_err != nil {
        resp, _ = json.Marshal(map[string]string{"error": parsed_err.Error()})
        w.WriteHeader(406)
      } else {
        resp, _ = json.Marshal(map[string]string{"str": parsed_time.Format("20060102150405")})
      }
      w.Header().Set("Content-Type", "application/json")
      w.Write(resp)
      return
    }
  }
  w.WriteHeader(404)
  return
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    w.WriteHeader(404)
    return
  }
  r.ParseForm()
  var time_in time.Time
  var delta_in time.Duration
  for k, v := range r.Form {
    if k == "time" {
      var resp []byte
      var parsed_err error
      json_time, _ := strconv.ParseFloat(strings.Join(v, ""), 9)
      if time_in, parsed_err = time_from_json(json_time); parsed_err != nil {
        resp, _ = json.Marshal(map[string]string{"error": parsed_err.Error()})
        w.WriteHeader(406)
        w.Write(resp)
        return
      }
    }
    if k == "delta" {
      json_time, _ := strconv.ParseFloat(strings.Join(v, ""), 9)
      delta_year   := 1e9 * time.Duration(int64(json_time / 1e4)) * 365 * 24 * 60 * 60
      delta_month  := 1e9 * time.Duration(int64(json_time / 1e2) % 1e2) * 30 * 24 * 60 * 60
      delta_day    := 1e9 * time.Duration(int64(json_time) % 1e2) * 24 * 60 * 60
      delta_hour   := 1e9 * time.Duration(int64(json_time * 1e2) % 1e2) * 60 * 60
      delta_minute := 1e9 * time.Duration(int64(json_time * 1e4) % 1e2) * 60
      delta_second := 1e9 * time.Duration(int64(json_time * 1e6) % 1e2)
      delta_millisecond := 1e6 * time.Duration(int64(json_time * 1e9) % 1e3)
      delta_in = delta_year + delta_month + delta_day + delta_hour + delta_minute + delta_second + delta_millisecond;
    }
  }
  if time_in.IsZero() || delta_in == 0 {
    w.WriteHeader(406)
    return
  } else {
    var resp []byte
    resp, _ = json.Marshal(map[string]float64{"time": time_to_json(time_in.Add(delta_in))})
    w.WriteHeader(200)
    w.Write(resp)
    return
  }
  w.WriteHeader(404)
  return
}

func SetHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    w.WriteHeader(404)
    return
  }
  r.ParseForm()
  for k, v := range r.Form {
    if k == "time" {
      var resp []byte
      json_time, _ := strconv.ParseFloat(strings.Join(v, ""), 9)
      if parsed_time, parsed_err := time_from_json(json_time); parsed_err != nil {
        resp, _ = json.Marshal(map[string]string{"error": parsed_err.Error()})
        w.WriteHeader(406)
        w.Write(resp)
        return
      } else {
        time_diff = parsed_time.Sub(time.Now().In(time.UTC))
        w.WriteHeader(200)
        return
      }
    }
  }
  w.WriteHeader(404)
  return
}

func ResetHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    w.WriteHeader(404)
    return
  }
  time_diff = 0
  w.WriteHeader(200)
}

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/time/now",    NowHandler)
  mux.HandleFunc("/time/string", StringHandler)
  mux.HandleFunc("/time/add",    AddHandler)
  mux.HandleFunc("/time/set",    SetHandler)
  mux.HandleFunc("/time/reset",  ResetHandler)
  err := http.ListenAndServe(":9000", mux)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
