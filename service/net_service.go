package service

import (
	"github.com/NumberMan1/MMO-server/config"
	core2 "github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/network/message_router"
	"strconv"
	"sync"
	"time"
)

// NetService 网络服务
type NetService struct {
	tcpServer *network.TcpServer
	//记录conn最后一次心跳包的时间
	//heartBeatPairs map[network.Connection]time.Time
	heartBeatPairs *sync.Map
	heartTicker    *time.Ticker
	cancel         chan struct{}
}

func NewNetService() *NetService {
	//server, _ := network.NewTcpServer("127.0.0.1:32510")
	server, _ := network.NewTcpServer(config.ServerConfig.Server.Host + ":" + strconv.FormatInt(int64(config.ServerConfig.Server.Port), 10))
	n := &NetService{
		tcpServer: server,
		//heartBeatPairs: map[network.Connection]time.Time{},
		heartBeatPairs: &sync.Map{},
		heartTicker:    time.NewTicker(5 * time.Second),
		cancel:         make(chan struct{}, 1),
	}
	server.SetConnectedCallback(network.TcpServerConnectedCallback{Op: n.onClientConnected})
	server.SetDisconnectedCallback(network.TcpServerDisconnectedCallback{Op: n.onDisconnected})
	return n
}

func (n *NetService) Start() {
	//启动网络监听，指定消息包装类型
	n.tcpServer.Start()
	//启动消息分发器
	network.GetMessageRouterInstance().Start(config.ServerConfig.Server.WorkerCount)
	network.GetMessageRouterInstance().Subscribe("proto.HeartBeatRequest", message_router.MessageHandler{Op: n.heartBeatRequest})
	go n.timerCallback()
}

func (n *NetService) Stop() {
	network.GetMessageRouterInstance().Off("proto.HeartBeatRequest", message_router.MessageHandler{Op: n.heartBeatRequest})
	err := n.tcpServer.Stop()
	if err != nil {
		return
	}
	n.heartTicker.Stop()
	n.cancel <- struct{}{}
}

// 收到心跳包
func (n *NetService) heartBeatRequest(msg message_router.Msg) {
	//n.heartBeatPairs[msg.Sender] = time.Now()
	n.heartBeatPairs.Store(msg.Sender.(network.Connection), time.Now())
	p := &proto.HeartBeatResponse{}
	msg.Sender.(network.Connection).Send(p)
	//更新Session的心跳
	session := msg.Sender.(network.Connection).Get("Session")
	if session != nil {
		session.(*core2.Session).HeartTime = time.Now()
	}
}

func (n *NetService) timerCallback() {
	for {
		select {
		case <-n.heartTicker.C:
			now := time.Now()
			n.heartBeatPairs.Range(func(k, v any) bool {
				tp := v.(time.Time)
				conn := k.(network.Connection)
				cha := now.Sub(tp)
				//关闭超时的客户端连接
				if cha.Seconds() > (10 * time.Second).Seconds() {
					conn.Close()
					n.heartBeatPairs.Delete(conn)
				}
				return true
			})
		case <-n.cancel:
			return
		}
	}
}

// 当客户端接入
func (n *NetService) onClientConnected(conn network.Connection) {
	logger.SLCInfo("客户端接入")
	n.heartBeatPairs.Store(conn, time.Now())
}

func (n *NetService) onDisconnected(conn network.Connection) {
	n.heartBeatPairs.Delete(conn)
	logger.SLCInfo("连接断开:%v", conn.Socket().RemoteAddr().String())
	session := conn.Get("Session")
	if session != nil {
		session.(*core2.Session).Character = nil
	}
	//character := conn.Get("Session").(*core2.Session).Character
	//if character != nil {
	//	space := character.Space()
	//	if space != nil {
	//		space.EntityLeave(character)
	//		core2.GetCharacterManagerInstance().RemoveCharacter(character.Id())
	//	}
	//}
}
