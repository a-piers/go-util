package mail

import (
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"
)

// OptionsFn is the option function type
type OptionsFn func(*Server)

// WithoutSSL disables SSL check for email server
func WithoutSSL() OptionsFn {
	return func(s *Server) {
		s.Dailer.SSL = false
	}
}

// SkipSSLVerify will not verify emial server certificate
func SkipSSLVerify() OptionsFn {
	return func(s *Server) {
		s.Dailer.TLSConfig.InsecureSkipVerify = true
	}
}

// Server represents the email server
type Server struct {
	Dailer *gomail.Dialer
}

// Message represents the email details
type Message struct {
	From     string
	To       []string
	Subject  string
	BodyText string
	BodyHTML string
}

// NewServer creates new Mail server connection
func NewServer(host string, port int, username, password string, options ...OptionsFn) (s *Server) {
	d := gomail.NewDialer(host, port, username, password)
	s = &Server{
		// Dialer: &gomail.Dialer{
		Dailer: d,
	}

	for _, optionFn := range options {
		optionFn(s)
	}

	return s
}

// Send function will send a message via the server using the Message details
func (s *Server) Send(m *Message) error {

	message := gomail.NewMessage()

	// From address
	message.SetHeader("From", m.From)

	// To address/recipients
	// m.SetAddressHeader("To", r.Address, r.Name)
	message.SetHeader("To", strings.Join(m.To, ";"))

	// Subject of message
	message.SetHeader("Subject", m.Subject)

	// Body of message
	if m.BodyText != "" {
		message.SetBody("text/plain", m.BodyText)
	} else {
		message.SetBody("text/html", m.BodyHTML)
	}

	// Dailer to send message
	// d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	// Send the message using the dailer.
	if err := s.Dailer.DialAndSend(message); err != nil {
		return fmt.Errorf("Unable to send message: %s", err)
	}
	return nil
}
