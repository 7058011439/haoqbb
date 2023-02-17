package SVar

import (
	"fmt"
	"testing"
)

func setSchool(s SVar) {
	s["name"] = "北京四中"
	s["class"] = "三连四班"
}

func TestNewSVar(t *testing.T) {
	s := NewSVar()
	s["age"] = 81
	s["score"] = 100
	setSchool(s.SVar("school"))
	s.SVar("addr").SVar("city")["area"] = "朝阳区"

	js := s.String()
	fmt.Println(js)
	fmt.Println(s["addr"].(SVar)["city"])
	fmt.Println(s.SVar("addr").SVar("city"))

	var pc []SVar
	for i := 0; i < 10; i++ {
		c := NewSVar()
		c["pc"] = "red"
		c["pn"] = i
		pc = append(pc, c)
	}
	s["card"] = pc

	for index, p := range s["card"].([]SVar) {
		fmt.Printf("key = %v, pc = %v, pn = %v\n", index, p["pc"], p["pn"])
	}

	var ps []SVar
	for i := 0; i < 10; i++ {
		c := NewSVar()
		c["x"] = "red"
		c["y"] = i
		ps = append(ps, c)
	}
	s["point"] = ps

	pl := s["point"].([]SVar)
	for index, p := range pl {
		fmt.Printf("key = %v, x = %v, y = %v\n", index, p["x"], p["y"])
	}

	s1 := NewSVarByData(js)
	fmt.Println(s1)
	fmt.Println("abcdefg")
}
