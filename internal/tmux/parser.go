package tmux

import (
	"regexp"
	"strings"
)

type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type Entity struct {
	Symbol string   `json:"symbol"`
	Pos    Position `json:"pos"`
	RelPos Position `json:"rel_pos"` // Relative to player
	Name   string   `json:"name,omitempty"`
}

type PlayerStatus struct {
	Name     string `json:"name"`
	St       string `json:"st"`
	Dx       string `json:"dx"`
	Co       string `json:"co"`
	In       string `json:"in"`
	Wi       string `json:"wi"`
	Ch       string `json:"ch"`
	Align    string `json:"align"`
	Dlvl     string `json:"dlvl"`
	Gold     string `json:"gold"`
	HP       string `json:"hp"`
	HPMax    string `json:"hp_max"`
	Pw       string `json:"pw"`
	PwMax    string `json:"pw_max"`
	AC       string `json:"ac"`
	Xp       string `json:"xp"`
	Hunger   string `json:"hunger,omitempty"`
	Encumber string `json:"encumber,omitempty"`
}

type Output struct {
	RawScreen string       `json:"-"`
	Screen    []string     `json:"screen"`
	PlayerPos Position     `json:"player_pos"`
	Status    PlayerStatus `json:"status"`
	Entities  []Entity     `json:"entities"`
	LocalView []string     `json:"local_view"` // 9x9 grid
	Messages  []string     `json:"messages"`
}

func Parse(raw string) Output {
	lines := strings.Split(raw, "\n")
	// Clean up carriage returns just in case
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, "\r")
	}
	// Remove empty last line if it exists (common with split)
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Smarter separation of Grid vs Status
	// Look for "Dlvl:"
	var statusLineIdx int = -1
	for i, line := range lines {
		if strings.Contains(line, "Dlvl:") {
			statusLineIdx = i
			break
		}
	}

	var gridLines []string
	var statusLines []string

	if statusLineIdx != -1 {
		// Status lines are usually this one and the one before it?
		// Or Dlvl is the last line? In the example:
		// Robot the Evoker ...
		// Dlvl:1 ...

		// So if we found Dlvl at `i`, then `i-1` is also status.
		startStatus := statusLineIdx - 1
		if startStatus < 0 {
			startStatus = 0
		}

		statusLines = lines[startStatus:]

		// Grid is everything before startStatus.
		// We might trim empty lines at the end of grid.
		gridLines = lines[:startStatus]
	} else {
		// Fallback
		if len(lines) >= 2 {
			gridLines = lines[:len(lines)-2]
			statusLines = lines[len(lines)-2:]
		} else {
			gridLines = lines
		}
	}

	// Find Player
	playerPos := Position{-1, -1}
	for r, line := range gridLines {
		if idx := strings.Index(line, "@"); idx != -1 {
			playerPos = Position{r, idx}
			break
		}
	}

	// Parse Entities
	var entities []Entity

	ignored := "-|# ." // Walls, corridors, floor.

	for r, line := range gridLines {
		for c, char := range line {
			s := string(char)
			if s == " " {
				continue
			}
			if strings.Contains(ignored, s) {
				continue
			}

			// If it's the player, we already tracked it, but maybe add to entities too?
			if s == "@" {
				continue
			}

			e := Entity{
				Symbol: s,
				Pos:    Position{r, c},
			}
			if playerPos.Row != -1 {
				e.RelPos = Position{r - playerPos.Row, c - playerPos.Col}
			}

			// Simple name guessing
			switch s {
			case "<":
				e.Name = "staircase up"
			case ">":
				e.Name = "staircase down"
			case "+":
				e.Name = "closed door"
			case "$":
				e.Name = "gold"
			case "{":
				e.Name = "fountain"
			case "_":
				e.Name = "altar"
			}

			entities = append(entities, e)
		}
	}

	// Local View (9x9)
	var localView []string
	viewRadius := 4 // 4+1+4 = 9
	if playerPos.Row != -1 {
		for r := playerPos.Row - viewRadius; r <= playerPos.Row+viewRadius; r++ {
			var rowStr string
			for c := playerPos.Col - viewRadius; c <= playerPos.Col+viewRadius; c++ {
				if r >= 0 && r < len(gridLines) && c >= 0 && c < len(gridLines[r]) {
					rowStr += string(gridLines[r][c])
				} else {
					rowStr += " " // Void
				}
			}
			localView = append(localView, rowStr)
		}
	}

	// Parse Status
	status := parseStatus(statusLines)

	return Output{
		RawScreen: raw,
		PlayerPos: playerPos,
		Status:    status,
		Entities:  entities,
		LocalView: localView,
	}
}

func parseStatus(lines []string) PlayerStatus {
	s := PlayerStatus{}
	if len(lines) == 0 {
		return s
	}

	// Line 1: Robot the Evoker St:9 Dx:14 Co:15 ...
	// Regex or simple split?
	// Let's use simple splitting and keyword matching for robustness
	fullText := strings.Join(lines, " ")

	// Split by double space to get name? Or just take everything before "St:"
	if idx := strings.Index(fullText, "St:"); idx != -1 {
		s.Name = strings.TrimSpace(fullText[:idx])
	}

	reVal := regexp.MustCompile(`([A-Za-z0-9$]+):(\d+(?:\(\d+\))?)`)

	matches := reVal.FindAllStringSubmatch(fullText, -1)

	for _, m := range matches {
		key := m[1]
		val := m[2]

		switch key {
		case "St":
			s.St = val
		case "Dx":
			s.Dx = val
		case "Co":
			s.Co = val
		case "In":
			s.In = val
		case "Wi":
			s.Wi = val
		case "Ch":
			s.Ch = val
		case "Dlvl":
			s.Dlvl = val
		case "$":
			s.Gold = val
		case "HP":

			parts := strings.Split(val, "(")
			s.HP = parts[0]
			if len(parts) > 1 {
				s.HPMax = strings.TrimRight(parts[1], ")")
			}
		case "Pw":
			parts := strings.Split(val, "(")
			s.Pw = parts[0]
			if len(parts) > 1 {
				s.PwMax = strings.TrimRight(parts[1], ")")
			}
		case "AC":
			s.AC = val
		case "Xp":
			s.Xp = val
		}
	}

	// Hunger/Encumbrance often plain words
	keywords := []string{"Hungry", "Weak", "Fainting", "Satiated", "Burdened", "Stressed", "Strained", "Overtaxed", "Overloaded"}
	for _, k := range keywords {
		if strings.Contains(fullText, k) {
			if k == "Hungry" || k == "Weak" || k == "Fainting" || k == "Satiated" {
				s.Hunger = k
			} else {
				s.Encumber = k
			}
		}
	}

	return s
}
