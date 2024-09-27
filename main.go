package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

type Counter struct {
	MaxAttacks int
	MaxMoves   int
	Attacks    int
	Moves      int
}

type ActionRequest struct {
	Action string `json:"action"`
}

type ActionResponse struct {
	Message string `json:"message"`
}

type Actions []Action

type Action struct {
	Unit   string `json:"unit"`
	Action string `json:"action"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

type GameState struct {
	AttackActionsAvailable int       `json:"attackActionsAvailable"`
	BoardSize              BoardSize `json:"boardSize"`
	EnemyUnits             []Unit    `json:"enemyUnits"`
	FriendlyUnits          []Unit    `json:"friendlyUnits"`
	MoveActionsAvailable   int       `json:"moveActionsAvailable"`
	Player1                string    `json:"player1"`
	Player2                string    `json:"player2"`
	TurnNumber             int       `json:"turnNumber"`
	UUID                   string    `json:"uuid"`
	YourID                 string    `json:"yourId"`
}

func (gs *GameState) OccupiedCords() []Cord {
	var occupiedCells []Cord

	for _, e := range gs.EnemyUnits {
		occupiedCells = append(occupiedCells, Cord{X: e.X, Y: e.Y})
	}
	for _, f := range gs.FriendlyUnits {
		occupiedCells = append(occupiedCells, Cord{X: f.X, Y: f.Y})
	}
	return occupiedCells
}

func (gs *GameState) EnemyCords() []Cord {
	var occupiedCells []Cord

	for _, e := range gs.EnemyUnits {
		occupiedCells = append(occupiedCells, Cord{X: e.X, Y: e.Y})
	}
	return occupiedCells
}

type BoardSize struct {
	H int `json:"h"`
	W int `json:"w"`
}

type Unit struct {
	Armor          int    `json:"armor"`
	AttackStrength int    `json:"attackStrength"`
	Attacks        int    `json:"attacks"`
	Health         int    `json:"health"`
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	MaxHealth      int    `json:"maxHealth"`
	Moves          int    `json:"moves"`
	Side           string `json:"side"`
	X              int    `json:"x"`
	Y              int    `json:"y"`
	Range          int    `json:"range"`
}

func (u *Unit) isNextToEnemy(gs GameState) (bool, Cord) {
	for _, e := range gs.EnemyUnits {
		if (u.X == e.X && u.Y == e.Y+1) || (u.X == e.X && u.Y == e.Y-1) || (u.X == e.X+1 && u.Y == e.Y) || (u.X == e.X-1 && u.Y == e.Y) {
			return true, Cord{e.X, e.Y}
		}
	}
	return false, Cord{0, 0}
}

type Cords []Cord

type Cord struct {
	X int
	Y int
}

func determineDirection(A Cord, B Cord) string {
	x1, y1 := A.X, A.Y
	x2, y2 := B.X, B.Y

	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 && dy == 0 {
		return "Same"
	}

	if abs(dx) >= abs(dy) {
		if dx > 0 {
			return "Right"
		} else {
			return "Left"
		}
	} else {
		if dy > 0 {
			return "Up"
		} else {
			return "Down"
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (u *Unit) canMoveTo(gs GameState, target Cord) bool {
	for _, e := range gs.OccupiedCords() {
		if target == e {
			return false
		}
	}
	return true
}

func (u *Unit) wizardAction(gs GameState, counter *Counter) Action {
	for _, e := range gs.EnemyUnits {
		if counter.Attacks < counter.MaxAttacks {
			counter.Attacks += 1
			return Action{
				Unit:   u.ID,
				X:      e.X,
				Y:      e.Y,
				Action: "attack",
			}
		}
	}

	return Action{}
}

func (u *Unit) archerAction(gs GameState, counter *Counter) Action {
	var attacks int = 0
	for _, e := range gs.EnemyUnits {
		if manhattanDistance(Cord{X: u.X, Y: u.Y}, Cord{X: e.X, Y: e.Y}) <= 4 {
			// if u.X == e.X || u.Y == e.Y {
			if counter.Attacks < counter.MaxAttacks && attacks <= 1 {
				counter.Attacks += 1
				attacks += 1
				return Action{
					Unit:   u.ID,
					X:      e.X,
					Y:      e.Y,
					Action: "attack",
				}
			}
			//}
		}
	}

	nextCord := closestEnemy(Cord{X: u.X, Y: u.Y}, gs.EnemyCords())
	moveDirection := determineDirection(Cord{X: u.X, Y: u.Y}, nextCord)

	var action Action

	switch moveDirection {
	case "Right":
		action.Unit = u.ID
		action.X = u.X + 1
		action.Y = u.Y
		action.Action = "move"
	case "Left":
		action.Unit = u.ID
		action.X = u.X - 1
		action.Y = u.Y
		action.Action = "move"
	case "Up":
		action.Unit = u.ID
		action.X = u.X
		action.Y = u.Y + 1
		action.Action = "move"
	case "Down":
		action.Unit = u.ID
		action.X = u.X
		action.Y = u.Y - 1
		action.Action = "move"
	}

	if u.canMoveTo(gs, Cord{X: action.X, Y: action.Y}) {
		counter.Moves += 1
		return action
	}

	return Action{}
}

func (u *Unit) move(gs GameState) Action {
	nextCord := closestEnemy(Cord{X: u.X, Y: u.Y}, gs.EnemyCords())
	moveDirection := determineDirection(Cord{X: u.X, Y: u.Y}, nextCord)

	var action Action

	switch moveDirection {
	case "Right":
		action.Unit = u.ID
		action.X = u.X + 1
		action.Y = u.Y
		action.Action = "move"
	case "Left":
		action.Unit = u.ID
		action.X = u.X - 1
		action.Y = u.Y
		action.Action = "move"
	case "Up":
		action.Unit = u.ID
		action.X = u.X
		action.Y = u.Y + 1
		action.Action = "move"
	case "Down":
		action.Unit = u.ID
		action.X = u.X
		action.Y = u.Y - 1
		action.Action = "move"
	}
	return action
}

func (u *Unit) attack(gs GameState) Action {
	nextCord := closestEnemy(Cord{X: u.X, Y: u.Y}, gs.EnemyCords())
	if manhattanDistance(Cord{X: u.X, Y: u.Y}, nextCord) <= u.Range {
		return Action{
			Unit:   u.ID,
			X:      nextCord.X,
			Y:      nextCord.Y,
			Action: "attack",
		}
	}
	return Action{}
}

func (u *Unit) moveUnit(gs GameState, counter *Counter) Action {
	nextCord := closestEnemy(Cord{X: u.X, Y: u.Y}, gs.EnemyCords())
	moveDirection := determineDirection(Cord{X: u.X, Y: u.Y}, nextCord)

	var action Action

	switch moveDirection {
	case "Right":
		action.Unit = u.ID
		action.X = u.X + 1
		action.Y = u.Y
		action.Action = "move"
	case "Left":
		action.Unit = u.ID
		action.X = u.X - 1
		action.Y = u.Y
		action.Action = "move"
	case "Up":
		action.Unit = u.ID
		action.X = u.X
		action.Y = u.Y + 1
		action.Action = "move"
	case "Down":
		action.Unit = u.ID
		action.X = u.X
		action.Y = u.Y - 1
		action.Action = "move"
	}

	if u.canMoveTo(gs, Cord{X: action.X, Y: action.Y}) {
		u.X = action.X
		u.Y = action.Y
		return action
	}

	nextTo, ec := u.isNextToEnemy(gs)
	if nextTo {
		log.Println("Next to enemy attack", ec)
		if counter.Attacks < counter.MaxAttacks {
			counter.Attacks += 1
			return Action{
				Unit:   u.ID,
				X:      ec.X,
				Y:      ec.Y,
				Action: "attack",
			}
		}
	}

	return Action{}
}

func manhattanDistance(a, b Cord) int {
	return int(math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y)))
}

func closestEnemy(unit Cord, enemies Cords) Cord {
	minDist := math.MaxInt64
	closest := enemies[0]

	for _, enemy := range enemies {
		dist := manhattanDistance(unit, enemy)
		if dist < minDist {
			minDist = dist
			closest = enemy
		}
	}
	return closest
}

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/", fileServer)

	mux.HandleFunc("/action/{kind}", handlerAction)
	mux.HandleFunc("/check", handlerServerCheck)

	http.ListenAndServe("localhost:8080", mux)
}

func handlerWeFuckedUp(w http.ResponseWriter, r *http.Request) {
	gs := GameState{}
	err := json.NewDecoder(r.Body).Decode(&gs)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	counter := &Counter{MaxAttacks: gs.AttackActionsAvailable, MaxMoves: gs.MoveActionsAvailable, Attacks: 0, Moves: 0}

	var actions Actions

	for i := 0; i <= gs.AttackActionsAvailable; i++ {
	}

	fmt.Println("Counter Stats:", counter.MaxAttacks, counter.MaxMoves, counter.Attacks, counter.Moves)
	json.NewEncoder(w).Encode(actions)
}

func handlerServerCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Got server check...")

	gs := GameState{}

	err := json.NewDecoder(r.Body).Decode(&gs)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	counter := &Counter{MaxAttacks: gs.AttackActionsAvailable, MaxMoves: gs.MoveActionsAvailable, Attacks: 0, Moves: 0}

	log.Println("Initial Counter Values:", counter)
	var actions Actions

	for i := 0; i < gs.AttackActionsAvailable; i++ {
		for _, u := range gs.FriendlyUnits {
			if u.Kind == "archer" {
				for i := 0; i < u.Attacks; i++ {
					action := u.archerAction(gs, counter)
					actions = append(actions, action)
				}
			}
		}

		if counter.Attacks >= counter.MaxAttacks {
			break
		}

	}

	for i := 0; i < gs.MoveActionsAvailable; i++ {
		for _, u := range gs.FriendlyUnits {
			if u.Kind == "barbarian" {
				action := u.moveUnit(gs, counter)
				actions = append(actions, action)
			}
		}

		for _, u := range gs.FriendlyUnits {
			if u.Kind == "knight" {
				action := u.moveUnit(gs, counter)
				actions = append(actions, action)
			}
		}
		for _, u := range gs.FriendlyUnits {
			if u.Kind == "warrior" {
				action := u.moveUnit(gs, counter)
				actions = append(actions, action)
			}
		}

	}

	var factions Actions
	for _, action := range actions {
		if action.Action != "" {
			factions = append(factions, action)
		}
	}

	fmt.Println("Counter Stats:", counter.MaxAttacks, counter.MaxMoves, counter.Attacks, counter.Moves)
	fmt.Println("Number of actions:", len(factions))
	json.NewEncoder(w).Encode(factions)
}

func handlerAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the kind parameter from the URL
	kind := r.PathValue("kind")

	// Set the content type to JSON for the response
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON body
	var actionReq ActionRequest
	err := json.NewDecoder(r.Body).Decode(&actionReq)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Use the parsed data (for example, print it or process it)
	fmt.Printf("Received action, message: %s, kind: %v\n", actionReq.Action, kind)

	// Create a response message
	response := ActionResponse{
		Message: "Action received: " + actionReq.Action,
	}

	// Encode the response as JSON and send it
	json.NewEncoder(w).Encode(response)
}
