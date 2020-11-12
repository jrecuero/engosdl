package engosdl_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jrecuero/engosdl"
	"github.com/jrecuero/engosdl/assets/components"
)

func TestComponent_ConvertToJSON(t *testing.T) {
	// component := engosdl.NewComponent("test")
	// fmt.Printf("%#+v\n", component)
	// fmt.Printf("%#+v\n", component.Object)
	// if result, err := json.Marshal(component); err == nil {
	// 	fmt.Printf("%s\n", result)
	// }
	// data := `{"id":"2","name":"test","loaded":false,"started":false,"entity":null,"active":true,"delegate":null,"registers":[]}`
	// byteString := []byte(data)
	// var comp engosdl.Component
	// if err := json.Unmarshal(byteString, &comp); err == nil {
	// 	fmt.Printf("%#+v\n", comp)
	// 	fmt.Printf("%#+v\n", comp.Object)
	// }

	player := engosdl.NewEntity("player")
	cp1 := components.NewEntityStats("player-stats", 100)
	cp2 := components.NewKeyboard("player-keyboard")
	player.AddComponent(cp1)
	player.AddComponent(cp2)
	fmt.Printf("%#+v\n", player)
	fmt.Printf("%#+v\n", player.Object)
	if result, err := json.Marshal(player); err != nil {
		fmt.Printf("result is: %s\n", result)
	} else {
		panic(err)
	}
}
