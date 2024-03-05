package price

type Pair struct {
	FirstAsset  string `json:"first_asset"`
	SecondAsset string `json:"second_asset"`
}

type Cache struct {
	prices map[Pair]float64 // TODO: Change float64 to a struct that contains signed message and signature
}

func (c *Cache) GetPrice(pair Pair) float64 {
	return c.prices[pair]
}

func (c *Cache) SetPrice(pair Pair, price float64) {
	c.prices[pair] = price
}

func NewCache() *Cache {
	return &Cache{make(map[Pair]float64)}
}
