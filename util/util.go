package util

import (
	"encoding/json"
	"fmt"
	"racelogctl/internal"
)

func DuplicateArray(src []interface{}) []interface{} {
	ret := make([]interface{}, len(src))
	copy(ret, src)
	return ret
}

func GetIntValue(src interface{}) int {
	switch v := src.(type) {
	case int:
		return int(v)
	case float64:
		return int(v)
	default:
		panic(src)
	}
}

func PatchSession(src []interface{}, patchData []interface{}) []interface{} {

	workCopy := DuplicateArray(src)
	for _, delta := range patchData {
		workDelta := delta.([]interface{})
		col := GetIntValue(workDelta[0])
		value := workDelta[1]
		workCopy[col] = value
	}
	return workCopy
}

// src contains the previous (complete) data [[RUN 6 ...], [RUN 7 ...]]
// pathData is an 2d array: [[patchRow,patchCol,value],...]
// we create a copy each row and
func PatchCars(src [][]interface{}, patchData [][]interface{}) [][]interface{} {
	// fmt.Printf("patch2dData: %v", patchData)
	// we do a "deep copy" manually here. the 2nd dimesion needs own instances, too
	workCopy := make([][]interface{}, len(src))
	copy(workCopy, src)

	for i := 0; i < len(src); i++ {
		// workCopy[i] = make([]interface{}, len(src[i]))
		// copy(workCopy[i], src[i])
		workCopy[i] = DuplicateArray(src[i])
	}
	for _, delta := range patchData {

		row := GetIntValue(delta[0])
		col := GetIntValue(delta[1])
		value := delta[2]
		workCopy[row][col] = value
	}
	return workCopy
}

// state contains the previous (complete) data
// incoming is the data coming via WAMP message. Depending on State.Type different actions apply
// returns the next state
func ProcessDeltaStates(state, incoming internal.State) internal.State {
	// fmt.Printf("patch2dData: %v", state)
	s := internal.State{}
	switch incoming.Type {
	case 1:
		s = incoming
		return s
	case 8:
		s.Type = 1
		s.Timestamp = incoming.Timestamp
		s.Payload.Cars = PatchCars(state.Payload.Cars, incoming.Payload.Cars)
		s.Payload.Session = PatchSession(state.Payload.Session, incoming.Payload.Session)
		if len(incoming.Payload.Messages) > 0 {
			fmt.Printf("have %d messages \n", len(incoming.Payload.Messages))
			s.Payload.Messages = incoming.Payload.Messages // messages don't have delta processing by design
		} else {
			s.Payload.Messages = [][]interface{}{}
		}
	}
	return s
}

func ConvertJsonToGo(jsonData []byte) internal.State {
	s := internal.State{}
	// logger.Printf("jsonData: %v", string(jsonData))
	err := json.Unmarshal(jsonData, &s)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("s: %v\n", s)
	return s
}
