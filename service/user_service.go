package service

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/MMO-server/mgr"
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/global/variable"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/network/message_router"
	"slices"
	"strings"
	"unicode/utf8"
)

var (
	singleUserService = singleton.Singleton{}
)

// UserService 玩家服务
// 注册，登录，创建角色，进入游戏
type UserService struct {
}

func GetUserServiceInstance() *UserService {
	instance, _ := singleton.GetOrDo[*UserService](&singleUserService, func() (*UserService, error) {
		return &UserService{}, nil
	})
	return instance
}

func (us *UserService) Start() {
	network.GetMessageRouterInstance().Subscribe("proto.GameEnterRequest", message_router.MessageHandler{Op: us.gameEnterRequest})
	network.GetMessageRouterInstance().Subscribe("proto.UserLoginRequest", message_router.MessageHandler{Op: us.userLoginRequest})
	network.GetMessageRouterInstance().Subscribe("proto.UserRegisterRequest", message_router.MessageHandler{Op: us.userRegisterRequest})
	network.GetMessageRouterInstance().Subscribe("proto.CharacterDeleteRequest", message_router.MessageHandler{Op: us.characterDeleteRequest})
	network.GetMessageRouterInstance().Subscribe("proto.CharacterListRequest", message_router.MessageHandler{Op: us.characterListRequest})
	network.GetMessageRouterInstance().Subscribe("proto.CharacterCreateRequest", message_router.MessageHandler{Op: us.characterCreateRequest})
	network.GetMessageRouterInstance().Subscribe("proto.ReviveRequest", message_router.MessageHandler{Op: us.reviveRequest})
	network.GetMessageRouterInstance().Subscribe("proto.PickupItemRequest", message_router.MessageHandler{Op: us.pickupItemRequest})
	network.GetMessageRouterInstance().Subscribe("proto.InventoryRequest", message_router.MessageHandler{Op: us.inventoryRequest})
	//物品放置请求
	network.GetMessageRouterInstance().Subscribe("proto.ItemPlacementRequest", message_router.MessageHandler{Op: us.itemPlacementRequest})
	//使用物品
	network.GetMessageRouterInstance().Subscribe("proto.ItemUseRequest", message_router.MessageHandler{Op: us.itemUseRequest})
	//丢弃物品
	network.GetMessageRouterInstance().Subscribe("proto.ItemDiscardRequest", message_router.MessageHandler{Op: us.itemDiscardRequest})
}

func (us *UserService) itemDiscardRequest(msg message_router.Msg) {
	rsq := msg.Message.(*proto.ItemDiscardRequest)
	chr, ok := model.GetUnit(int(rsq.EntityId)).(*model.Character)
	if !ok {
		return
	}
	chr.Knapsack.Discard(int(rsq.SlotIndex), int(rsq.Count))
	chr.SendInventory(true, false, false)
}

func (us *UserService) itemUseRequest(msg message_router.Msg) {
	rsq := msg.Message.(*proto.ItemUseRequest)
	chr, ok := model.GetUnit(int(rsq.EntityId)).(*model.Character)
	if !ok {
		return
	}
	chr.UseItem(int(rsq.SlotIndex))
}

func (us *UserService) itemPlacementRequest(msg message_router.Msg) {
	rsq := msg.Message.(*proto.ItemPlacementRequest)
	chr, ok := model.GetUnit(int(rsq.EntityId)).(*model.Character)
	if !ok {
		return
	}
	chr.Knapsack.Exchange(int(rsq.OriginIndex), int(rsq.TargetIndex))
	//发送背包数据
	chr.SendInventory(true, false, false)
}

func (us *UserService) inventoryRequest(msg message_router.Msg) {
	rsq := msg.Message.(*proto.InventoryRequest)
	chr, ok := model.GetUnit(int(rsq.EntityId)).(*model.Character)
	if !ok {
		return
	}
	//发送背包数据
	chr.SendInventory(rsq.QueryKnapsack, rsq.QueryWarehouse, rsq.QueryEquipment)
}

func (us *UserService) pickupItemRequest(msg message_router.Msg) {
	s1 := msg.Sender.(network.Connection).Get("Session").(*model.Session)
	if s1 == nil {
		return
	}
	chr := s1.Character
	units := model.RangeUnit(chr.Position(), chr.Space().Id, 3000)
	items := make([]*model.ItemEntity, 0)
	for e := units.Front(); e != nil; e = e.Next() {
		if ie, ok := e.Value.(*model.ItemEntity); ok {
			items = append(items, ie)
		}
	}
	if len(items) == 0 {
		return
	}
	itemEntity := slices.MinFunc(items, func(a, b *model.ItemEntity) int {
		f := vector3.GetDistance(a.Position(), chr.Position()) > vector3.GetDistance(b.Position(), chr.Position())
		if f {
			return 1
		} else {
			return 0
		}
	})
	//如果添加失败则结束
	if !chr.Knapsack.AddItem(itemEntity.Item().Id(), itemEntity.Item().Amount()) {
		return
	}
	//物品模型移出场景
	chr.Space().EntityLeave(itemEntity)
	mgr.GetEntityManagerInstance().RemoveEntity(chr.Space().Id, itemEntity)
	logger.SLCInfo("玩家拾取物品Chr[%v],背包[%v]", chr.CharacterId(), chr.Knapsack.InventoryInfo())
	//发送背包数据
	chr.SendInventory(true, false, false)
}

func (us *UserService) reviveRequest(msg message_router.Msg) {
	message := msg.Message.(*proto.ReviveRequest)
	actor := model.GetUnit(int(message.GetEntityId()))
	if chr, ok := actor.(*model.Character); ok && chr.IsDeath() && chr.Conn == msg.Sender {
		sp := GetSpaceServiceInstance().GetSpace(1)
		model.TeleportSpace(sp, vector3.Zero3(), vector3.Zero3(), chr)
		chr.Revive()
	}
}

func (us *UserService) userRegisterRequest(msg message_router.Msg) {
	message := msg.Message.(*proto.UserRegisterRequest)
	conn := msg.Sender.(network.Connection)
	var num int64
	variable.GDb.Model(&database.DbPlayer{}).Where("username = ?", message.Username).Count(&num)
	logger.SLCInfo("新用户注册:%s", message.Username)
	resp := &proto.UserRegisterResponse{}
	if num > 0 {
		resp.Code = 1
		resp.Message = "用户名已存在"
	} else {
		dbPlayer := &database.DbPlayer{
			Username: message.Username,
			Password: message.Password,
		}
		variable.GDb.Save(dbPlayer)
		resp.Code = 6
		resp.Message = "注册成功"
	}
	conn.Send(resp)
}

// 删除角色的请求
func (us *UserService) characterDeleteRequest(msg message_router.Msg) {
	conn := msg.Sender.(network.Connection)
	player := conn.Get("Session").(*model.Session).DbPlayer
	variable.GDb.Where("id = ?", msg.Message.(*proto.CharacterDeleteRequest).CharacterId).
		Where("player_id = ?", player.ID).
		Delete(&database.DbCharacter{})
	//给客户端响应
	rsp := &proto.CharacterDeleteResponse{
		Success: true,
		Message: "执行完成",
	}
	conn.Send(rsp)
}

// 查询角色列表的请求
func (us *UserService) characterListRequest(msg message_router.Msg) {
	conn := msg.Sender.(network.Connection)
	player := conn.Get("Session").(*model.Session).DbPlayer
	characters := make([]database.DbCharacter, 0)
	//从数据库查询出当前玩家的全部角色
	variable.GDb.Where("player_id = ?", player.ID).Find(&characters)
	rsp := &proto.CharacterListResponse{CharacterList: make([]*proto.NetActor, 0)}
	for _, character := range characters {
		rsp.CharacterList = append(rsp.CharacterList, &proto.NetActor{
			Id: int32(character.ID),
			//EntityId: character.EntityId,
			Tid:     int32(character.JobId),
			Name:    character.Name,
			Level:   int32(character.Level),
			Exp:     int64(character.Exp),
			SpaceId: int32(character.SpaceId),
			Gold:    character.Gold,
			//Entity:   nil,
		})
	}
	msg.Sender.(network.Connection).Send(rsp)
}

// 创建角色
func (us *UserService) characterCreateRequest(msg message_router.Msg) {
	logger.SLCInfo("创建角色:%v", msg.Message)
	conn := msg.Sender.(network.Connection)
	rsp := &proto.ChracterCreateResponse{
		Success:   false,
		Character: nil,
	}
	player := conn.Get("Session").(*model.Session).DbPlayer
	if player == nil {
		// 未登录不能创建角色
		logger.SLCInfo("未登录不能创建角色")
		rsp.Message = "未登录不能创建角色"
		conn.Send(rsp)
		return
	}
	var num int64
	variable.GDb.Model(&database.DbCharacter{}).Where("player_id = ?", player.ID).Count(&num)
	if num >= 4 {
		// 角色数量最多4个
		logger.SLCInfo("角色数量最多4个")
		rsp.Message = "角色数量最多4个"
		conn.Send(rsp)
		return
	}
	msgTemp := msg.Message.(*proto.CharacterCreateRequest)
	nameLen := utf8.RuneCountInString(msgTemp.Name)
	// 判断角色名是否为空或包含非法字符如空格等
	if nameLen == 0 || strings.ContainsAny(msgTemp.Name, " \t\r\n\\") {
		logger.SLCInfo("创建角色失败，角色名不能为空或包含非法字符如空格等")
		rsp.Message = "判断角色名不能为空或包含非法字符如空格等"
		conn.Send(rsp)
		return
	}
	//角色名最长7个字
	if nameLen > 7 {
		logger.SLCInfo("创建角色失败，角色名不能超过7个字符")
		rsp.Message = "创建角色失败，角色名不能超过7个字符"
		conn.Send(rsp)
		return
	}
	//检验角色名是否存在
	variable.GDb.Model(&database.DbCharacter{}).Where("name = ?", msgTemp.Name).Count(&num)
	if num > 0 {
		logger.SLCInfo("创建角色失败，角色名已存在")
		rsp.Message = "创建角色失败，角色名已存在"
		conn.Send(rsp)
		return
	}
	//出生点坐标|
	birthPos := vector3.NewVector3(354947, 1660, 308498)
	dbCharacter := database.NewDbCharacter()
	dbCharacter.Name = msgTemp.Name
	dbCharacter.JobId = int(msgTemp.JobType)
	dbCharacter.SpaceId = 2
	dbCharacter.X = int(birthPos.X)
	dbCharacter.Y = int(birthPos.Y)
	dbCharacter.Z = int(birthPos.Z)
	dbCharacter.PlayerId = int(player.ID)
	tx := variable.GDb.Save(dbCharacter)
	if tx.RowsAffected > 0 {
		rsp.Success = true
		rsp.Message = "角色创建成功"
		conn.Send(rsp)
	}
}

func (us *UserService) userLoginRequest(msg message_router.Msg) {
	req := msg.Message.(*proto.UserLoginRequest)
	dbPlayer := &database.DbPlayer{}
	result := variable.GDb.Where("username = ? and password = ?", req.Username, req.Password).First(&dbPlayer)
	rsp := &proto.UserLoginResponse{}
	conn := msg.Sender.(network.Connection)
	if result.Error != nil {
		rsp.Success = false
		rsp.Message = result.Error.Error()
		return
	}
	if result.RowsAffected > 0 {
		rsp.Success = true
		rsp.Message = "登录成功"
		conn.Get("Session").(*model.Session).DbPlayer = dbPlayer //登录成功，在conn里记录用户信息
	} else {
		rsp.Success = false
		rsp.Message = "用户名或密码错误"
	}
	conn.Send(rsp)
}

func (us *UserService) gameEnterRequest(msg message_router.Msg) {
	rsq := msg.Message.(*proto.GameEnterRequest)
	conn := msg.Sender.(network.Connection)
	logger.SLCInfo("有玩家进入游戏,角色Id:%d", rsq.CharacterId)
	// 获取当前玩家
	player := conn.Get("Session").(*model.Session).DbPlayer
	// 查询数据库的角色
	dbRole := &database.DbCharacter{}
	variable.GDb.Where("player_id = ?", player.ID).
		Where("id = ?", rsq.CharacterId).First(dbRole)
	logger.SLCInfo("dbRole = %v", dbRole)
	// 把数据库角色变成游戏角色
	character := model.GetCharacterManagerInstance().CreateCharacter(dbRole)
	//角色与conn关联
	character.Conn = conn
	//角色存入session
	conn.Get("Session").(*model.Session).Character = character
	////通知玩家登录成功
	//response := &proto.GameEnterResponse{
	//	Success:   true,
	//	Entity:    character.EntityData(),
	//	Character: character.Info(),
	//}
	//msg.Sender.Send(response)
	//将新角色加入到地图
	space := GetSpaceServiceInstance().GetSpace(dbRole.SpaceId) //新手村
	space.EntityEnter(character)
}
