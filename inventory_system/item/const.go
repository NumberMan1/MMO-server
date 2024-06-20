package item

type ItemType int

const (
	ItemType_Material   ItemType = iota //材料&道具
	ItemType_Consumable                 //消耗品
	ItemType_Equipment                  //武器&装备
)

type Quality int

const (
	Quality_Common    = iota //普通
	Quality_Uncommon         //非凡
	Quality_Rare             //稀有
	Quality_Epic             //史诗
	Quality_Legendary        //传说
	Quality_Artifact         //神器
)
