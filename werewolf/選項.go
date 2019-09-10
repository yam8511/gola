package werewolf

// RuleOptions 角色選單
func RuleOptions() map[RULE]GROUP {
	return map[RULE]GROUP{
		平民:  人質,
		預言家: 神職,
		女巫:  神職,
		獵人:  神職,
		騎士:  神職,
		狼人:  狼職,
		狼王:  狼職,
		雪狼:  狼職,
	}
}
