package syslog

import (
	"errors"
	"flag"
	"log"
	"log/syslog"
	"strings"

	"github.com/tgragnato/goflow/transport"
)

type customWriter struct {
	writer *syslog.Writer
}

func (cw *customWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	parts := strings.SplitN(message, " ", 3)
	if len(parts) == 3 {
		// Remove the timestamp from the log message (YYYY/MM/DD HH:MM:SS message)
		message = parts[2]
		// Remove leading and trailing whitespaces
		message = strings.TrimSpace(message)
	}
	return cw.writer.Write([]byte(message))
}

type SyslogDriver struct {
	protocol string
	address  string
}

func (s *SyslogDriver) Prepare() error {
	flag.StringVar(&s.protocol, "transport.syslog.protocol", "udp", "Replace udp with tcp if your syslog server uses TCP")
	flag.StringVar(&s.address, "transport.syslog.address", "localhost:514", "Remote syslog server address")
	return nil
}

func (s *SyslogDriver) Init() error {
	remoteSyslog, err := syslog.Dial(s.protocol, s.address, syslog.LOG_INFO|syslog.LOG_LOCAL0, "goflow")
	if err != nil {
		return errors.New("Failed to connect to remote syslog server: " + err.Error())
	}

	log.SetOutput(&customWriter{writer: remoteSyslog})
	return nil
}

func (s *SyslogDriver) Send(key, data []byte) error {
	log.Println(string(data))
	return nil
}

func (s *SyslogDriver) Close() error {
	return nil
}

func init() {
	s := &SyslogDriver{}
	transport.RegisterTransportDriver("syslog", s)
}
