package statemachine

import (
	"bytes"
	"errors"
	"fmt"
)

/*
type transition struct {
	event   string
	toState string
}
*/

type StateMachine struct {
	smc                      *StateMachineConfig
	states                   []string
	transitions              map[string]map[string]string
	currentState             string
	panicOnIllegalTransition bool
	onEnters                 map[string]func(string, interface{})
	onExits                  map[string]func(string, interface{})
}

func NewStateMachine() *StateMachine {
	sm := StateMachine{
		states:      make([]string, 0),
		transitions: make(map[string]map[string]string),
		onEnters:    make(map[string]func(string, interface{})),
		onExits:     make(map[string]func(string, interface{})),
	}
	smc := NewStateMachineConfig(&sm)
	sm.smc = &smc

	return &sm
}

func (sm *StateMachine) containsState(state string) bool {
	contains := false
	for _, b := range sm.states {
		if b == state {
			contains = true
			break
		}
	}
	return contains
}

func (sm *StateMachine) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("StateMachine representation\n"))
	for _, state := range sm.states {
		format := "\t%s\n"
		if sm.GetCurrentState() == state {
			format = "\t*%s*\n"
		}
		buf.WriteString(fmt.Sprintf(format, state))
		for e, s := range sm.transitions[state] {
			buf.WriteString(fmt.Sprintf("\t\t%s -> %s\n", e, s))
		}
	}
	return string(buf.Bytes())
}

func (sm *StateMachine) Configure() *StateMachineConfig {
	return sm.smc
}

func (sm *StateMachine) GetCurrentState() string {
	if sm.currentState == "" {
		panic(errors.New("state machine not initialized"))
	}
	return sm.currentState
}

func (sm *StateMachine) Trigger(event string, payload interface{}) {
	if event == "" {
		panic(errors.New("cannot trigger with empty event"))
	}

	allowedTransitions := sm.transitions[sm.currentState]
	nextState, ok := allowedTransitions[event]
	if !ok && sm.panicOnIllegalTransition {
		panic(fmt.Errorf("attempted illegal transition: %s + %s", sm.currentState, event))
	}

	//fire exit event
	onExit, ok := sm.onExits[sm.currentState]
	if ok {
		onExit(sm.currentState, payload)
	}

	sm.currentState = nextState

	//fire enter event
	onEnter, ok := sm.onEnters[sm.currentState]
	if ok {
		onEnter(sm.currentState, payload)
	}
}

type StateMachineConfig struct {
	sm *StateMachine
}

func NewStateMachineConfig(sm *StateMachine) StateMachineConfig {
	return StateMachineConfig{sm: sm}
}

func (smc *StateMachineConfig) AddState(state string) *StateMachineConfig {
	if state == "" {
		panic(errors.New("non-empty argument expected: state"))
	}

	if smc.sm.containsState(state) {
		panic(fmt.Errorf("state machine already contains state: %s", state))
	}

	smc.sm.states = append(smc.sm.states, state)

	//set current state if not set
	if smc.sm.currentState == "" {
		smc.sm.currentState = state
	}

	return smc
}

func (smc *StateMachineConfig) AddStates(states ...string) *StateMachineConfig {
	if len(states) == 0 {
		panic(errors.New("non-empty argument expected: states"))
	}

	for _, state := range states {
		smc.AddState(state)
	}

	//set current state if not set
	if smc.sm.currentState == "" {
		smc.sm.currentState = states[0]
	}

	return smc
}

func (smc *StateMachineConfig) SetInitialState(state string) *StateMachineConfig {
	if state == "" {
		panic(errors.New("non-empty argument expected: state"))
	}

	if !smc.sm.containsState(state) {
		panic(fmt.Errorf("state machine does not contain state: %s", state))
	}

	smc.sm.currentState = state
	return smc
}

func (smc *StateMachineConfig) AddTransition(fromState string, toState string, event string) *StateMachineConfig {

	if !smc.sm.containsState(fromState) {
		panic(fmt.Errorf("unknown state: %s", fromState))
	}

	if !smc.sm.containsState(toState) {
		panic(fmt.Errorf("unknown state: %s", toState))
	}

	//get transition slice for the "from" state
	tmap := smc.sm.transitions[fromState]
	//make if null
	if tmap == nil {
		tmap = make(map[string]string)

	}
	//append
	tmap[event] = toState

	//add
	smc.sm.transitions[fromState] = tmap
	return smc
}

func (smc *StateMachineConfig) OnEnter(state string, f func(string, interface{})) *StateMachineConfig {
	if !smc.sm.containsState(state) {
		panic(fmt.Errorf("unknown state: %s", state))
	}
	smc.sm.onEnters[state] = f
	return smc
}

func (smc *StateMachineConfig) OnExit(state string, f func(string, interface{})) *StateMachineConfig {
	if !smc.sm.containsState(state) {
		panic(fmt.Errorf("unknown state: %s", state))
	}
	smc.sm.onExits[state] = f
	return smc
}

func (smc *StateMachineConfig) PanicOnIllegalTransition(truefalse bool) {
	smc.sm.panicOnIllegalTransition = truefalse
}
