package raft

import "time"

type Config struct {
	Name      string        `required:"true"`
	Tick      time.Duration `default:"1s"`
	Timeout   time.Duration `default:"500ms"`
	Aggregate bool          `default:"true"`
	Leader    string        `required:"false"`
	Peers     []Peer        `required:"false"`
}

type Peer struct {
	PID      uint32
	Name     string
	BindAddr string
	Endpoint string
}
