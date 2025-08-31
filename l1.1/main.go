package main

import "fmt"

type Human struct {
	Name       string
	Age        int
	Weight     float32
	HairColour string
}

func (h Human) Say(phrase string) {
	fmt.Println(h.Name, "says:", phrase)
}

func (h *Human) Eat() {
	h.Weight += 0.5
	fmt.Println(h.Name, "weights", h.Weight, "kg")
}

func (h *Human) ChangeHairColour(colour string) {
	h.HairColour = colour
	fmt.Println(h.Name, "has", h.HairColour, "hair")

}

func (h *Human) Birthday() {
	h.Age++
	fmt.Println(h.Name, "is now", h.Age, "years old")
}

type Action struct {
	Human
	IsDoing string
}

func (a Action) Act() {
	fmt.Println(a.Name, a.IsDoing)
}

func main() {
	h := Human{Name: "Egor", Age: 20, Weight: 60, HairColour: "brown"}
	a := Action{
		Human:   h,
		IsDoing: "coding",
	}
	a.Say("Hi")
	a.Eat()
	a.ChangeHairColour("blonde")
	a.Birthday()
	a.Act()

}
