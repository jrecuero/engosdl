package engosdl_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	cp2 := components.NewKeyboard("player-keyboard", map[int]bool{})
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

func TestComponent_RegisterComponent(t *testing.T) {
	newComponent := struct {
		componentName string
		name          string
		life          int
	}{
		componentName: components.ComponentNameEntityStats,
		name:          "test-stats",
		life:          50,
	}
	// componentName := "entity-stats"
	// name := "test-stats"
	// life := 50
	entity := engosdl.NewEntity("entity")
	// for _, component := range engosdl.GetComponentManager().GetComponents() {
	// 	if component.GetName() == newComponent.componentName {
	// 		// instance := reflect.Zero(reflect.TypeOf(component))
	// 		// obj := instance.Interface()
	// 		if constructor := engosdl.GetComponentManager().Constructors[newComponent.componentName]; constructor != nil {
	// 			stats := constructor(newComponent.name, newComponent.life)
	// 			entity.AddComponent(stats)
	// 			entity.DoDump()
	// 			fmt.Println(entity)
	// 		}
	// 	}
	// }
	for k, constructor := range engosdl.GetComponentManager().Constructors {
		if k == newComponent.componentName {
			// if constructor := engosdl.GetComponentManager().Constructors[newComponent.componentName]; constructor != nil {
			stats := constructor(newComponent.name, newComponent.life)
			entity.AddComponent(stats)
			entity.DoDump()
			fmt.Println(entity)
			// }
		}
	}
}

func TestComponent_ReadJSON(t *testing.T) {
	file, err := ioutil.ReadFile("assets/test/entities.json")
	if err != nil {
		panic(err)
	}
	// entities := []*entityToMarshal{}
	entities := []*engosdl.EntityToUnmarshal{}
	err = json.Unmarshal([]byte(file), &entities)
	if err != nil {
		panic(err)
	}
	for _, instance := range entities {
		entity := engosdl.NewEntity("")
		entity.Unmarshal(instance)
		fmt.Println(entity)
	}
}
