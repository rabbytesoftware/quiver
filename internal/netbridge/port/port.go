package port

import (
	"fmt"
	"time"
)

type Port struct {
	PortName    string    `json:"name"`
	PortNumber  uint16    `json:"port"`
	HostAddr    string    `json:"host"`
	PortProto   string    `json:"protocol"`
	CreateTime  time.Time `json:"created_at"`
	Desc        string    `json:"description"`
}

func NewPort(
	name 		string,
	port 		uint16, 
	host 		string, 
	protocol 	string,
) Port {
	return Port{
		PortName:    name,
		PortNumber:  port,
		HostAddr:    host,
		PortProto:   protocol,
		CreateTime:  time.Now(),
		Desc:        "",
	}
}

func NewPortWithDescription(
	name 		string,
	port 		uint16, 
	host 		string, 
	protocol 	string,
	description string,
) Port {
	return Port{
		PortName:    name,
		PortNumber:  port,
		HostAddr:    host,
		PortProto:   protocol,
		CreateTime:  time.Now(),
		Desc:        description,
	}
}

func (p Port) Name() string {
	return p.PortName
}

func (p Port) Port() uint16 {
	return p.PortNumber
}

func (p Port) Host() string {
	return p.HostAddr
}

func (p Port) Protocol() string {
	return p.PortProto
}

func (p Port) CreatedAt() time.Time {
	return p.CreateTime
}

func (p Port) Description() string {
	return p.Desc
}

func (p Port) String() string {
	return fmt.Sprintf("%s:%d/%s", p.HostAddr, p.PortNumber, p.PortProto)
}
