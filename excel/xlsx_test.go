package excel

import "testing"

type Address struct {
	City string
	Addr string
}

type Student struct {
	ID   int
	Name string
	Addr Address
	Address
	Score  float32
	Status bool
}

func TestXlsx(t *testing.T) {
	v := []Student{
		{
			ID:   1,
			Name: "张三",
			Address: Address{
				City: "河南",
				Addr: "郑州",
			},
			Addr: Address{
				City: "河南",
				Addr: "郑州",
			},
			Score:  60.0,
			Status: true,
		},
		{
			ID:   2,
			Name: "李四",
			Address: Address{
				City: "河北",
				Addr: "石家庄",
			},
			Addr: Address{
				City: "河北",
				Addr: "石家庄",
			},
			Score:  95.3,
			Status: false,
		},
	}
	Sheet(&Student{}, "student").Write(&v).ToFile("./")
}
