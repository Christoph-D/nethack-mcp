package tmux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

type MapData struct {
	Turn         int64         `json:"turn"`
	DungeonLevel int           `json:"dungeon_level"`
	Cursor       *string       `json:"cursor,omitempty"`
	Tiles        Tiles         `json:"tiles"`
	Monsters     []interface{} `json:"monsters,omitempty"`
	Items        []interface{} `json:"items,omitempty"`
	Hero         string        `json:"hero"`
}

type Tiles struct {
	Walls           [][]int       `json:"walls,omitempty"`
	StoneWalls      [][]int       `json:"stone_walls,omitempty"`
	OpenSpaces      [][]int       `json:"open_spaces,omitempty"`
	Air             [][]int       `json:"air,omitempty"`
	Cloud           [][]int       `json:"cloud,omitempty"`
	Special         []interface{} `json:"special,omitempty"`
	UnexploredTiles [][]int       `json:"unexplored_tiles,omitempty"`
}

func tilePosToSet(arr [][]int) map[string]bool {
	set := make(map[string]bool)
	for _, pos := range arr {
		if len(pos) == 2 {
			set[fmt.Sprintf("%d,%d", pos[0], pos[1])] = true
		}
	}
	return set
}

func compareTilePosArrays(a, b [][]int) bool {
	aSet := tilePosToSet(a)
	bSet := tilePosToSet(b)

	if len(aSet) != len(bSet) {
		return false
	}

	for k := range aSet {
		if !bSet[k] {
			return false
		}
	}
	return true
}

func compareInterfaceArrays(a, b []interface{}) bool {
	return reflect.DeepEqual(a, b)
}

func compareTiles(a, b Tiles) map[string]bool {
	changed := make(map[string]bool)

	changed["walls"] = !compareTilePosArrays(a.Walls, b.Walls)
	changed["stone_walls"] = !compareTilePosArrays(a.StoneWalls, b.StoneWalls)
	changed["open_spaces"] = !compareTilePosArrays(a.OpenSpaces, b.OpenSpaces)
	changed["air"] = !compareTilePosArrays(a.Air, b.Air)
	changed["cloud"] = !compareTilePosArrays(a.Cloud, b.Cloud)
	changed["special"] = !compareInterfaceArrays(a.Special, b.Special)
	changed["unexplored_tiles"] = !compareTilePosArrays(a.UnexploredTiles, b.UnexploredTiles)

	return changed
}

func generateDiff(current, previous *MapData) string {
	if previous == nil {
		b, _ := json.Marshal(current)
		return string(b)
	}

	var output strings.Builder
	output.WriteString("{\n")

	output.WriteString(fmt.Sprintf("  \"turn\": %d,\n", current.Turn))
	output.WriteString(fmt.Sprintf("  \"dungeon_level\": %d,\n", current.DungeonLevel))

	if current.Cursor != nil {
		output.WriteString(fmt.Sprintf("  \"cursor\": \"%s\",\n", *current.Cursor))
	}

	output.WriteString("  \"tiles\": {\n")

	tileChanges := compareTiles(current.Tiles, previous.Tiles)
	tilesFields := []struct {
		name  string
		value interface{}
		omit  bool
	}{
		{"walls", current.Tiles.Walls, len(current.Tiles.Walls) == 0},
		{"stone_walls", current.Tiles.StoneWalls, len(current.Tiles.StoneWalls) == 0},
		{"open_spaces", current.Tiles.OpenSpaces, len(current.Tiles.OpenSpaces) == 0},
		{"air", current.Tiles.Air, len(current.Tiles.Air) == 0},
		{"cloud", current.Tiles.Cloud, len(current.Tiles.Cloud) == 0},
		{"special", current.Tiles.Special, len(current.Tiles.Special) == 0},
		{"unexplored_tiles", current.Tiles.UnexploredTiles, len(current.Tiles.UnexploredTiles) == 0},
	}

	var tileWritten bool
	for _, field := range tilesFields {
		if field.omit {
			continue
		}

		if tileWritten {
			output.WriteString(",\n")
		}
		tileWritten = true

		if !tileChanges[field.name] {
			output.WriteString(fmt.Sprintf("    // %s: unchanged", field.name))
			continue
		}

		jsonBytes, _ := json.Marshal(field.value)
		output.WriteString(fmt.Sprintf("    \"%s\": %s", field.name, string(jsonBytes)))
	}

	if tileWritten {
		output.WriteString("\n  },\n")
	} else {
		output.WriteString("\n  },\n")
	}

	if len(current.Monsters) > 0 {
		if !reflect.DeepEqual(current.Monsters, previous.Monsters) {
			jsonBytes, _ := json.Marshal(current.Monsters)
			output.WriteString(fmt.Sprintf("  \"monsters\": %s,\n", string(jsonBytes)))
		} else {
			output.WriteString("  // monsters: unchanged,\n")
		}
	} else {
		output.WriteString("  // monsters: unchanged,\n")
	}

	if len(current.Items) > 0 {
		if !reflect.DeepEqual(current.Items, previous.Items) {
			jsonBytes, _ := json.Marshal(current.Items)
			output.WriteString(fmt.Sprintf("  \"items\": %s,\n", string(jsonBytes)))
		} else {
			output.WriteString("  // items: unchanged,\n")
		}
	} else {
		output.WriteString("  // items: unchanged,\n")
	}

	output.WriteString(fmt.Sprintf("  \"hero\": \"%s\"\n", current.Hero))
	output.WriteString("}")

	return output.String()
}

func GetTarget() string {
	target := os.Getenv("NETHACK_TMUX_SESSION")
	if target == "" {
		return "nethack"
	}
	return target
}

func CapturePane(target string) (string, error) {
	args := []string{"capture-pane", "-p", "-t", target}
	cmd := exec.Command("tmux", args...)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Println(string(exitErr.Stderr))
			return "", fmt.Errorf("tmux capture-pane failed: %w", err)
		}
		return "", fmt.Errorf("tmux capture-pane failed: %w", err)
	}

	output := stdout.String()

	dumpFile := os.Getenv("NETHACK_DUMP_FILENAME")
	if dumpFile == "" {
		dumpFile = "/tmp/" + GetTarget() + "-map.json"
	}

	currentContent, err := os.ReadFile(dumpFile)
	if err != nil {
		return output, nil
	}

	var current MapData
	if err := json.Unmarshal(currentContent, &current); err != nil {
		return output, nil
	}

	previousFile := dumpFile + ".previous"
	var previous *MapData
	if prevContent, err := os.ReadFile(previousFile); err == nil {
		var prev MapData
		if err := json.Unmarshal(prevContent, &prev); err == nil {
			previous = &prev
		}
	}

	diff := generateDiff(&current, previous)
	output += "\n" + diff

	return output, nil
}

func SendKeys(target string, keys []string) (string, error) {
	if len(keys) > 5 {
		return "", fmt.Errorf("Please send at most 5 keys at a time")
	}

	for _, key := range keys {
		isValid := false

		if len(key) == 1 {
			isValid = true
		} else if strings.HasPrefix(key, "C-") && len(key) == 3 {
			isValid = true
		} else if strings.EqualFold(key, "Enter") || strings.EqualFold(key, "Space") || strings.EqualFold(key, "Escape") {
			isValid = true
		}

		if !isValid {
			return "", fmt.Errorf("invalid key '%s': must be a single character, C-<char> (e.g. C-x), or Enter/Space/Escape", key)
		}
	}

	// Normalize special keys to correct capitalization for tmux
	for i, key := range keys {
		if key == " " {
			keys[i] = "Space"
		} else if key == "\n" {
			keys[i] = "Enter"
		} else if strings.EqualFold(key, "Space") {
			keys[i] = "Space"
		} else if strings.EqualFold(key, "Enter") {
			keys[i] = "Enter"
		} else if strings.EqualFold(key, "Escape") {
			keys[i] = "Escape"
		}
	}

	dumpFile := os.Getenv("NETHACK_DUMP_FILENAME")
	if dumpFile == "" {
		dumpFile = "/tmp/" + GetTarget() + "-map.json"
	}
	previousFile := dumpFile + ".previous"
	os.Rename(dumpFile, previousFile)

	for _, key := range keys {
		cmd := exec.Command("tmux", "send-keys", "-t", target, key)

		err := cmd.Run()
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Println(string(exitErr.Stderr))
				return "", fmt.Errorf("tmux send-keys failed: %w", err)
			}
			return "", fmt.Errorf("tmux send-keys failed: %w", err)
		}

		time.Sleep(200 * time.Millisecond)
	}

	return CapturePane(target)
}
