package td

type TowerDefense struct {
	ready  chan struct{}
}

func NewTowerDefense() *TowerDefense {
	return &TowerDefense{
		ready: make(chan struct{}, 0),
	}
}



