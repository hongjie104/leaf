package gate

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/hongjie104/leaf/chanrpc"
	"github.com/hongjie104/leaf/log"
	"github.com/hongjie104/leaf/network"
)

// Gate Gate
type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor
	AgentChanRPC    *chanrpc.Server

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool
}

// Run Run
func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

// OnDestroy OnDestroy
func (gate *Gate) OnDestroy() {}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debugf("read message: %v", err)
			break
		}

		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Debugf("unmarshal message error: %v", err)
				break
			}
			// if conf.RunMode == "debug" {
			// 	log.Debugf("receive msg = %s\n", string(data))
			// }
			err = a.gate.Processor.Route(msg, a)
			if err != nil {
				log.Debugf("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		err := a.gate.AgentChanRPC.Call0("CloseAgent", a)
		if err != nil {
			log.Errorf("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Errorf("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		// if conf.RunMode == "debug" {
		// 	log.Debugf("send msg, id = %s, data = %s\n", string(data[0]), string(data[1]))
		// }
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Errorf("write message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		a.log(msg)
	}
}

func (a *agent) log(m interface{}) {
	// roleID := ""
	// userData := a.UserData()
	// if userData != nil {
	// 	r := userData.(*data.Role)
	// 	roleID = r.ID.Hex()
	// }
	t := reflect.TypeOf(m).Elem()
	if t.Name() != "S2C_SystemTime" {
		v := reflect.ValueOf(m).Elem()
		key := ""
		var tmp []string
		values := ""
		for i := 0; i < t.NumField(); i++ {
			tmp = strings.Split(t.Field(i).Tag.Get("sproto"), ",")
			key = strings.Split(tmp[len(tmp)-1], "=")[1]
			values += fmt.Sprintf("key = %s val = %v,", key, v.FieldByName(t.Field(i).Name))
		}
		log.Debugf(fmt.Sprintf("[send msg]msg=%s,%s", t.Name(), values))
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
