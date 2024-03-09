package service

import (
	core2 "github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/network/core"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"time"
)

// NetService 网络服务
type NetService struct {
	tcpServer *core.TcpServer
	//记录conn最后一次心跳包的时间
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
	//启动网络监听，指定消息包装类型
	n.tcpServer.Start()
	//启动消息分发器
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

// 收到心跳包
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
				//关闭超时的客户端连接
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

// 当客户端接入
func (n *NetService) onClientConnected(conn network.Connection) {
	logger.SLCInfo("客户端接入")
	n.heartBeatPairs[conn] = time.Now()
	conn.Set("Session", core2.NewSession())
}

func (n *NetService) onDisconnected(conn network.Connection) {
	delete(n.heartBeatPairs, conn)
	logger.SLCInfo("连接断开:%v", conn.Socket().RemoteAddr().String())
	character := conn.Get("Session").(*core2.Session).Character
	if character != nil {
		space := character.Space()
		if space != nil {
			space.EntityLeave(character)
			core2.GetCharacterManagerInstance().RemoveCharacter(character.Id())
		}
	}
}
