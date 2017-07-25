package main

import (
	"fmt"
	"github.com/teo-mateo/gosm/statemachine"
)

const (
	State1 = "STATE1"
	State2 = "STATE2"
	State3 = "STATE3"
)

const (
	e1 = "e1"
	e2 = "e2"
	e3 = "e3"
	e4 = "e4"
)

func main() {
	test1()
}

func test1() {
	sm := statemachine.NewStateMachine()
	sm.Configure().
		AddStates(State1, State2, State3).
		AddTransition(State1, State2, e1).
		AddTransition(State2, State3, e2).
		AddTransition(State3, State2, e1).
		AddTransition(State3, State1, e3).
		AddTransition(State1, State3, e4).
		OnEnter(State2, func(state string, payload interface{}) {
			fmt.Println("Entering", state)
		}).
		OnEnter(State3, func(state string, payload interface{}) {
			fmt.Println("Entering", state)
		})

	fmt.Println(sm)
	sm.Trigger(e1, nil)
	sm.Trigger(e2, nil)
	fmt.Println(sm)
}
