package nudge

import (
  "log"
  "testing"
  "net/http"
  "net/http/httptest"
  "github.com/stretchr/testify/assert"
  "fmt"
)

var successString = `<SchedulerService status="success" error="" info=""></SchedulerService>`
var errorString   = `<SchedulerService status="error" error="" info=""></SchedulerService>`

func TestRequestHandler(t *testing.T) {
  assert := assert.New(t)
  
  r := RequestHandler{Url: "http://example.com"}
  if assert.NotNil(r) {
    assert.Equal("http://example.com", r.Url)
  }
}

func TestHandleResponse(t *testing.T) {
  assert := assert.New(t)
  ts := SpawnTestServer(successString)
  defer ts.Close()
  r := RequestHandler{Url: ts.URL}
  res, err := http.Get(ts.URL)
  if err != nil {
    log.Fatal(err)
  }
  result, err := r.HandleResponse(res)
  if err != nil {
    log.Fatal(err)
  }
  assert.Equal("success", result.SchedulerService.Status)
}

func TestEnqueue(t *testing.T) {
  assert := assert.New(t)
  ts := SpawnTestServer(successString)
  defer ts.Close()
  r := NewRequest(ts.URL)
  result, err := r.Enqueue()
  if err != nil {
    log.Fatal(err)
  }
  assert.Equal("success", result.SchedulerService.Status)
}

func TestErrorAlert(t *testing.T) {
  assert := assert.New(t)
  ts     := SpawnTestServer(errorString)
  defer ts.Close()
  RecipientList([]string{"pffft@fart.com"})
  err := SendEmailAlert("testing email alert")
  if err != nil {
    log.Fatal(err)
  }
  assert.Equal(err, nil)
}

func SpawnTestServer(response string) *httptest.Server {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, fmt.Sprintf(`
      <?xml version="1.0" encoding="UTF-8"?>
      <string xmlns="http://fart.com/">
        %s
      </string>
      `, response))
    }))
    return ts
}
