package packet

type Product struct {
	ProductId   uint32
	Title       string
	Point       uint32
	ProductType uint8
	Value       string
	Stock       uint32
	ExpiryDate  int64
}

type Products struct {
	Data []*Product
}

type ProductWithImage struct {
	Product
	ImageId uint32
}

type ProducesWithImage struct {
	Data []*ProductWithImage
}
