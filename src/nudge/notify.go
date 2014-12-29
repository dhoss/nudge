package nudge

import (
  "log"
  "net/smtp"
  "bytes"
  "fmt"
)

var recipients []string
var fullMessage string

func RecipientList(r []string) (){
  recipients = r
}

func SendEmailAlert(message string) error {
  // make me a config option
  c, err := smtp.Dial("mail.fart.com:25")
  if err != nil {
    log.Fatal(err)
    return err
  }

  // make me a config option
  if err := c.Mail("alert-thingy@fart.com"); err != nil {
    log.Fatal(err)
    return err
  }

  for _, address := range recipients {
    log.Printf("Adding %s to recipient list", address)
    if err := c.Rcpt(address); err != nil {
      log.Fatal(err)
      return err
    }
    fullMessage += fmt.Sprintf("To: %s\r\n", address)
  }

  wc, err := c.Data()
  if err != nil {
    log.Fatal(err)
    return err
  }

  // clean me up to do something like this: https://gist.github.com/andelf/5004821
  fullMessage += fmt.Sprintf("Subject: Fart Error\r\n\r\n%s", message)

  buf := bytes.NewBufferString(fullMessage)
  if _,err = buf.WriteTo(wc); err != nil {
    log.Fatal(err)
    return err
  }

  err = wc.Close()
  if err != nil {
    log.Fatal(err)
    return err
  }

  err = c.Quit()
  if err != nil {
    log.Fatal(err)
    return err
  }
  return nil
}

