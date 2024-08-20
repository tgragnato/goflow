package syslog

import (
	"errors"
	"flag"
	"log"
	"log/syslog"

	"github.com/tgragnato/goflow/transport"
)

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

	log.SetOutput(remoteSyslog)
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
