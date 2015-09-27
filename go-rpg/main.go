package main

import (
	"fmt"

	"./characters"
)

func main() {
	fmt.Println("Hello World!")
	myHero := characters.Hero{Health: 100, Name: "Casey", Attack: 12, Defense: 4}
	aMonster := characters.Monster{Health: 6, Name: "Skeleton", Defense: 12, Attack: 6}

	fmt.Println("Hero has encountered a Skeleman!")

	fmt.Printf("%s's HP Remaining: %d\n", myHero.Name, myHero.Health)
	fmt.Printf("%s's HP Remaining: %d\n", aMonster.Name, aMonster.Health)
	myHero.Fight(&aMonster)
	fmt.Printf("%s's HP Remaining: %d\n", myHero.Name, myHero.Health)
	fmt.Printf("%s's HP Remaining: %d\n", aMonster.Name, aMonster.Health)
}
