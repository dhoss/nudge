package nudge

import (
    "io/ioutil"
    "net/http"
    "encoding/xml"
    "log"
    "bytes"
    "errors"
)

// takes a URL to send a blank POST request to
// e.g.  http://charge.myfamilysouth.com/schedulerservice.asmx/RunScheduledTasks
type RequestHandler struct {
  Url string
}

// parses the XML response
// e.g.:
//  <?xml version="1.0" encoding="UTF-8"?>
//  <string xmlns="http://fart.com/">
//      <SchedulerService status="success" error="" info=""></SchedulerService>
//  </string>
type SchedulerService struct {
  Status string `xml:"status,attr"`
  Error  string `xml:"error,attr"`
  Info   string `xml:"info,attr"`
}

type Result struct {
  XMLName xml.Name `xml:"string"`
  SchedulerService SchedulerService
}

// simply send an empty POST request to the URL provided
func (r *RequestHandler) Enqueue() (Result, error) {
  res, err := http.Post(r.Url, "text/plain", bytes.NewBufferString("just a string")); 
  if err != nil {
    log.Fatal("Rouse: ", err)
    return Result{}, err
  }
  handler, err := r.HandleResponse(res)
  if err != nil {
    log.Fatal("Rouse: ", err)
    return Result{}, err
  }
  return handler, nil
}

func (r *RequestHandler) HandleResponse(res *http.Response) (Result, error) {
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body); 
  if err != nil {
    log.Fatal("HandleResponse: ", err)
    return Result{}, err
  }
  v := Result{}
  if err := xml.Unmarshal([]byte(body),&v); err != nil {
    log.Fatal("HandleResponse: ", err)
    return Result{}, err
  }

  if v.SchedulerService.Error != "" {
    log.Fatal("HandleResponse: ", v.SchedulerService.Error)
    // make me config
    RecipientList([]string{"pfffft@fart.com"})
    err := SendEmailAlert(v.SchedulerService.Error)
    if err != nil {
      log.Fatal(err)
    }
    return Result{}, errors.New(v.SchedulerService.Error)
  }

  return v, nil
}

func NewRequest(url string) RequestHandler {
  return RequestHandler{Url: url}
}
