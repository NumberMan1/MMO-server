package service

import (
	"github.com/NumberMan1/MMO-server/mgr"
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/network/core"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"time"
)

type NetService struct {
	tcpServer      *core.TcpServer
	heartBeatPairs map[network.Connection]time.Time
	heartTicker    *time.Ticker
	cancel         chan struct{}
}

func NewNetService() *NetService {
	server, _ := core.NewTcpServer("127.0.0.1:32510")
	n := &NetService{
		tcpServer:      server,
		heartBeatPairs: map[network.Connection]time.Time{},
		heartTicker:    time.NewTicker(5 * time.Second),
		cancel:         make(chan struct{}, 1),
	}
	server.SetConnectedCallback(core.TcpServerConnectedCallback{Op: n.onClientConnected})
	server.SetDisconnectedCallback(core.TcpServerDisconnectedCallback{Op: n.onDisconnected})
	return n
}

func (n *NetService) Start() {
	n.tcpServer.Start()
	network.GetMessageRouterInstance().Start(4)
	network.GetMessageRouterInstance().Subscribe("proto.HeartBeatRequest", network.MessageHandler{Op: n.heartBeatRequest})
	go n.timerCallback()
}

func (n *NetService) Stop() {
	network.GetMessageRouterInstance().Off("proto.HeartBeatRequest", network.MessageHandler{Op: n.heartBeatRequest})
	err := n.tcpServer.Stop()
	if err != nil {
		return
	}
	n.heartTicker.Stop()
	n.cancel <- struct{}{}
}

func (n *NetService) heartBeatRequest(msg network.Msg) {
	n.heartBeatPairs[msg.Sender] = time.Now()
	p := &proto.HeartBeatResponse{}
	msg.Sender.Send(p)
}

func (n *NetService) timerCallback() {
	for {
		select {
		case <-n.heartTicker.C:
			now := time.Now()
			for conn, tp := range n.heartBeatPairs {
				cha := now.Sub(tp)
				if cha.Seconds() > (10 * time.Second).Seconds() {
					conn.Close()
					delete(n.heartBeatPairs, conn)
				}
			}
		case <-n.cancel:
			return
		}
	}
}

func (n *NetService) onClientConnected(conn network.Connection) {
	logger.SLCInfo("客户端接入")
	n.heartBeatPairs[conn] = time.Now()
}

func (n *NetService) onDisconnected(conn network.Connection) {
	delete(n.heartBeatPairs, conn)
	logger.SLCInfo("连接断开:%v", conn.Socket().RemoteAddr().String())
	character := conn.Get("Character")
	if character != nil {
		space := character.(*model.Character).Space
		if space != nil {
			co := conn.Get("Character").(*model.Character)
			space.CharacterLeave(conn, co)
			mgr.GetCharacterManagerInstance().RemoveCharacter(character.(*model.Character).Id)
		}
	}
}
