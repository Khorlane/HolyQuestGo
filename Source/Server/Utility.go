//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Utility.go                                            *
// Usage:     Provides utility functions for the game server        *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
)

// Calculate a percentage
func CalcPct(dividend, divisor int) int {
	x := float64(dividend)
	y := float64(divisor)
	z := x / y * 100.0
	return int(z)
}

// Insert commas into a numeric string
func FormatCommas(input string) string {
	x := len(input)
	if x < 4 {
		return input
	}
	y := x / 3
	z := x % 3
	if z == 0 {
		y--
	}
	j := x - 3
	for i := y; i > 0; i-- {
		input = input[:j] + "," + input[j:]
		j -= 3
	}
	return input
}

// Get home directory
func GetHomeDir() string {
	homeDirFileName := "HomeDir.txt"
	file, err := os.Open(homeDirFileName)
	if err != nil {
		log.Fatalf("GetHomeDir - Open HomeDir file failed: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

// Get a random number
func GetRandomNumber(limit int) int {
	randomNumber := min(rand.Intn(limit) + 1, limit)
	return randomNumber
}

// Get a SQL statement
func GetSqlStmt(sqlStmtId string) string {
	sqlStmtFileName := SQL_DIR + sqlStmtId + ".txt"
	file, err := os.Open(sqlStmtFileName)
	if err != nil {
		log.Fatalf("GetSqlStmt - Open SqlStmt file failed: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sqlStmt strings.Builder
	for scanner.Scan() {
		sqlStmt.WriteString(scanner.Text() + " ")
	}
	return strings.TrimSpace(sqlStmt.String())
}

// Translate a word
func TranslateWord(word string) string {
	translations := map[string]string{
		"n":  "go north",
		"s":  "go south",
		"e":  "go east",
		"w":  "go west",
		"ne": "go northeast",
		"se": "go southeast",
		"sw": "go southwest",
		"nw": "go northwest",
		"u":  "go up",
		"d":  "go down",
		"con": "consider",
		"des": "destroy",
		"em":  "emote",
		"eq":  "equipment",
		"i":   "inventory",
		"k":   "kill",
		"l":   "look",
		"obj": "object",
		"mob": "mobile",
		"north": "go north",
		"south": "go south",
		"east": "go east",
		"west": "go west",
		"northeast": "go northeast",
		"southeast": "go southeast",
		"southwest": "go southwest",
		"northwest": "go northwest",
		"up": "go up",
		"down": "go down",
	}
	if translation, exists := translations[word]; exists {
		return translation
	}
	return word
}
