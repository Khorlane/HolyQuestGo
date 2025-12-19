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

// Get a random number between 1 and Limit
func GetRandomNumber(Limit int) int {
  return rand.Intn(Limit) + 1
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

func PronounSubstitute(MsgText string) string {
  var PronounHeShe string
  var PronounHimHer string
  var PronounHisHers string
  var PronounHimselfHerself string

  if pDnodeSrc.pPlayer.Sex == "M" {
    PronounHeShe          = "he"
    PronounHimHer         = "him"
    PronounHisHers        = "his"
    PronounHimselfHerself = "himself"
  } else {
    PronounHeShe          = "she"
    PronounHimHer         = "her"
    PronounHisHers        = "her"
    PronounHimselfHerself = "herself"
  }
  StrReplace(&MsgText, "$P",               pDnodeSrc.PlayerName)
  StrReplace(&MsgText, "$pHeShe",          PronounHeShe)
  StrReplace(&MsgText, "$pHimHer",         PronounHimHer)
  StrReplace(&MsgText, "$pHisHers",        PronounHisHers)
  StrReplace(&MsgText, "$pHimselfHerself", PronounHimselfHerself)
  if pDnodeTgt != nil {
    if pDnodeTgt.pPlayer.Sex == "M" {
      PronounHeShe          = "he"
      PronounHimHer         = "him"
      PronounHisHers        = "his"
      PronounHimselfHerself = "himself"
    } else {
      PronounHeShe          = "she"
      PronounHimHer         = "her"
      PronounHisHers        = "her"
      PronounHimselfHerself = "herself"
    }
    StrReplace(&MsgText, "$T",               pDnodeTgt.PlayerName)
    StrReplace(&MsgText, "$tHeShe",          PronounHeShe)
    StrReplace(&MsgText, "$tHimHer",         PronounHimHer)
    StrReplace(&MsgText, "$tHisHers",        PronounHisHers)
    StrReplace(&MsgText, "$tHimselfHerself", PronounHimselfHerself)
  }
  return MsgText
}

func TranslateWord(Word string) string {
  if Word == "n"            { return "go north" }
  if Word == "s"            { return "go south" }
  if Word == "e"            { return "go east" }
  if Word == "w"            { return "go west" }

  if Word == "ne"           { return "go northeast" }
  if Word == "se"           { return "go southeast" }
  if Word == "sw"           { return "go southwest" }
  if Word == "nw"           { return "go northwest" }

  if Word == "u"            { return "go up" }
  if Word == "d"            { return "go down" }

  if Word == "con"          { return "consider" }
  if Word == "des"          { return "destroy" }
  if Word == "em"           { return "emote" }
  if Word == "eq"           { return "equipment" }
  if Word == "i"            { return "inventory" }
  if Word == "k"            { return "kill" }
  if Word == "l"            { return "look" }

  if Word == "obj"          { return "object" }
  if Word == "mob"          { return "mobile" }

  if Word == "north"        { return "go north" }
  if Word == "south"        { return "go south" }
  if Word == "east"         { return "go east" }
  if Word == "west"         { return "go west" }

  if Word == "northeast"    { return "go northeast" }
  if Word == "southeast"    { return "go southeast" }
  if Word == "southwest"    { return "go southwest" }
  if Word == "northwest"    { return "go northwest" }

  if Word == "up"           { return "go up" }
  if Word == "down"         { return "go down" }

  if Word == "head"         { return "01" }
  if Word == "earleft"      { return "02" }
  if Word == "earright"     { return "03" }
  if Word == "neck"         { return "04" }
  if Word == "shoulders"    { return "05" }
  if Word == "chest"        { return "06" }
  if Word == "back"         { return "07" }
  if Word == "arms"         { return "08" }
  if Word == "wristleft"    { return "09" }
  if Word == "wristright"   { return "10" }
  if Word == "hands"        { return "11" }
  if Word == "fingerleft"   { return "12" }
  if Word == "fingerright"  { return "13" }
  if Word == "shield"       { return "14" }
  if Word == "waist"        { return "15" }
  if Word == "legs"         { return "16" }
  if Word == "ankleleft"    { return "17" }
  if Word == "ankleright"   { return "18" }
  if Word == "feet"         { return "19" }
  if Word == "wielded"      { return "20" }
  if Word == "light"        { return "21" }

  if Word == "01"           { return "Worn on head:            " }
  if Word == "02"           { return "Worn in left ear:        " }
  if Word == "03"           { return "Worn in right ear:       " }
  if Word == "04"           { return "Worn around neck:        " }
  if Word == "05"           { return "Worn on shoulders:       " }
  if Word == "06"           { return "Worn on chest:           " }
  if Word == "07"           { return "Worn on back:            " }
  if Word == "08"           { return "Worn on arms:            " }
  if Word == "09"           { return "Worn on left wrist       " }
  if Word == "10"           { return "Worn on right wrist      " }
  if Word == "11"           { return "Worn on hands:           " }
  if Word == "12"           { return "Worn on left finger:     " }
  if Word == "13"           { return "Worn on right finger:    " }
  if Word == "14"           { return "Worn as a shield:        " }
  if Word == "15"           { return "Worn around waist:       " }
  if Word == "16"           { return "Worn on legs:            " }
  if Word == "17"           { return "Worn around left ankle:  " }
  if Word == "18"           { return "Worn around right ankle: " }
  if Word == "19"           { return "Worn on feet:            " }
  if Word == "20"           { return "Weapon wielded:          " }
  if Word == "21"           { return "Using as a light:        " }

  if Word == "hand"         { return "slap" }
  if Word == "axe"          { return "chop" }
  if Word == "club"         { return "bash" }
  if Word == "dagger"       { return "stab" }
  if Word == "hammer"       { return "pound" }
  if Word == "spear"        { return "pierce" }
  if Word == "staff"        { return "whack" }
  if Word == "sword"        { return "slash" }

  return Word
}