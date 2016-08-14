package shop

// ManagerProducts informer of products by id
type ManagerProducts interface {

	// GetInfo get info of product by id
	GetInfo(string) (*ProductInfo, error)
}
