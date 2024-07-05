package define

type ItemDefine struct {
	ID          int     `json:"ID" bson:"id"`                   //物品ID
	Name        string  `json:"Name" bson:"name"`               //名称
	ItemType    string  `json:"ItemType" bson:"item_type"`      //物品类别
	EquipsType  string  `json:"EquipsType" bson:"equips_type"`  //装备类型
	Quality     string  `json:"Quality" bson:"quality"`         //品质
	Description string  `json:"Description" bson:"description"` //介绍
	Capicity    int     `json:"Capicity" bson:"capicity"`       //堆叠上限
	BuyPrice    int     `json:"BuyPrice" bson:"buy_price"`      //买入价格
	SellPrice   int     `json:"SellPrice" bson:"sell_price"`    //卖出价格
	Icon        string  `json:"Icon" bson:"icon"`               //图标资源
	Model       string  `json:"Model" bson:"model"`             //场景模型
	Speed       float32 `json:"Speed" bson:"speed"`             //移动速度
	HPMax       float32 `json:"HPMax" bson:"hp_max"`            //生命上限
	MPMax       float32 `json:"MPMax" bson:"mp_max"`            //法力上限
	AD          float32 `json:"AD" bson:"ad"`                   //物攻
	AP          float32 `json:"AP" bson:"ap"`                   //魔攻
	DEF         float32 `json:"DEF" bson:"def"`                 //物防
	MDEF        float32 `json:"MDEF" bson:"mdef"`               //魔防
	CRI         float32 `json:"CRI" bson:"cri"`                 //暴击率
	CRD         float32 `json:"CRD" bson:"crd"`                 //暴击伤害
	HitRate     float32 `json:"HitRate" bson:"hit_rate"`        //命中率
	DodgeRate   float32 `json:"DodgeRate" bson:"dodge_rate"`    //闪避率
	HpRegen     float32 `json:"HpRegen" bson:"hp_regen"`        //生命恢复/秒
	HpSteal     float32 `json:"HpSteal" bson:"hp_steal"`        //伤害吸血
	STR         float32 `json:"STR" bson:"str"`                 //力量
	INT         float32 `json:"INT" bson:"int"`                 //智力
	AGI         float32 `json:"AGI" bson:"agi"`                 //敏捷
}

func (i *ItemDefine) GetId() int {
	return i.ID
}
