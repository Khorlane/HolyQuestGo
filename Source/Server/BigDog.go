//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      BigDog.go                                             *
// Usage:     Starting point for all HolyQuest stuff                *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// BigDog is the main entry point for the OMugs server
func BigDog() {
  var EventTick      int
  var GoGoGoFileName string
  var MobHealTick    int
  var StopItFileName string

  PrintIt("OMugs Starting")
  HomeDir = GetHomeDir()
  if err := ChgDir(HomeDir); err != nil {
    // Change directory failed
    PrintIt("BigDog() - Change directory to HomeDir failed")
    os.Exit(1)
  }
  // Set Go Stop, force Go status
  StopItFileName = CONTROL_DIR + "StopIt"
  GoGoGoFileName = CONTROL_DIR + "GoGoGo"
  if FileExist(StopItFileName) {
    // If StopIt file exists, Rename it to GoGoGo
    if err := Rename(StopItFileName, GoGoGoFileName); err != nil {
      PrintIt("BigDog() - Rename of 'StopIt' to 'GoGoGo' failed!")
      os.Exit(1)
    }
  }
  // Log game startup
  OpenLogFile()
  LogBuf = "OMugs version " + VERSION + " has started"
  LogIt(LogBuf)
  LogBuf = "Home directory is " + HomeDir
  LogIt(LogBuf)
  // Initialize
  EventTick        = EVENT_TICK
  MobHealTick      = 0
  StateConnections = true
  StateRunning     = true
  StateStopping    = false
  PACMN            = 1.0 / MAC * MDRP / 100.0
  if ValErr := ValidateIt("All"); ValErr {
    // Validation failed
    LogBuf = "OMugs has stopped"
    LogIt(LogBuf)
    CloseLogFile()
    return
  }
  // Validation was ok, so open port, init, play on
  SockOpenPort(PORT_NBR)
  InitDescriptor()
  CalendarConstructor()
  for StateRunning {
    // Game runs until it is stopped
    Sleep(MILLI_SECONDS_TO_SLEEP)
    AdvanceTime()
    if !StateStopping {
      // Game is not stopping, but should it be?
      if FileExist(StopItFileName) {
        // StopIt file was found, Stop the game
        StateStopping = true
        LogBuf = "Game is stopping"
        LogIt(LogBuf)
      }
    }
    if !StateStopping {
      // No new connections after stop command
      SockCheckForNewConnections()
      if StateConnections && DnodeCount == 1 {
        // No players connected
        LogBuf = "No Connections - going to sleep"
        LogIt(LogBuf)
        StateConnections = false
      }
    }
    if StateConnections {
      // One or more players are connected
      SockRecv()
      EventTick++
      if EventTick >= EVENT_TICK {
        // Time to process events
        EventTick = 0
        Events()
      }
      MobHealTick++
      if MobHealTick >= MOB_HEAL_TICK {
        // Time to heal mobs
        MobHealTick = 0
        HealMobiles()
      }
    } else {
      // No connections
      if StateStopping {
        // Game is stopping
        StateRunning = false
      }
    }
  }
  // Game has stopped so clean up
  ClearDescriptor()
  SockClosePort(PORT_NBR)
  CalendarDestructor()
  LogBuf = "OMugs has stopped"
  LogIt(LogBuf)
  CloseLogFile()
}

// Put one shot test code here, great for debugging. Always runs first
func AppTestCode() {
}

// Change the current working directory
func ChgDir(Dir string) error {
  return os.Chdir(Dir)
}

// Check if a file exists
func FileExist(Name string) bool {
  _, err := os.Stat(Name)
  return !os.IsNotExist(err)
}

// Rename a file
func Rename(File1, File2 string) error {
  return os.Rename(File1, File2)
}

// Remove deletes a file and returns an error value
func Remove(File1 string) error {
  err := os.Remove(File1)
  return err
}

// Return the number of seconds since the epoch
func GetTimeSeconds() int64 {
  Now := time.Now()
  Epoch := Now.Unix()
  return Epoch
}

// Print a message to stdout
func PrintIt(Message string) {
  Message = "\r\n" + Message + "\r\n"
  println(Message)
}

// Pause execution for the specified duration
func Sleep(Milliseconds int) {
  time.Sleep(time.Duration(Milliseconds) * time.Millisecond)
}

//-----------------------------------------------------------------------------
// String utility functions
//-----------------------------------------------------------------------------

// Count characters in a string
func StrCountChar(Str1 string, c rune) int {
  Count := 0
  for _, Char := range Str1 {
    if Char == c {
      Count++
    }
  }
  return Count
}

// Count the number of words in a string
func StrCountWords(Str1 string) int {
  // Squeeze spaces in the string
  S1 := StrSqueeze(Str1)
  // Split the string into words
  Words := strings.Fields(S1)
  // Return the number of words
  return len(Words)
}

// Delete characters from a string
func StrDelete(s string, Pos int, Length int) string {
  if Pos < 0 || Pos >= len(s) || Length <= 0 {
    return s
  }
  End := Pos + Length
  if End > len(s) {
    End = len(s)
  }
  return s[:Pos] + s[End:]
}

// Delete the word specified by WordNbr, squeezing string first
func StrDeleteWord(Str1 string, WordNbr int) string {
  S1 := StrSqueeze(Str1)
  if S1 == "" {
    return ""
  }
  Words := strings.Split(S1, " ")
  Index := WordNbr - 1
  if Index < 0 || Index >= len(Words) {
    return Str1
  }
  Words = append(Words[:Index], Words[Index+1:]...)
  Str1 = ""
  for _, s := range Words {
    Str1 += s
    Str1 += " "
  }
  return Str1
}

// Find first occurrence of Needle in HayStack
func StrFind(HayStack, Needle string) int {
  return strings.Index(HayStack, Needle)
}

// Find first occurrence of Needle in HayStack after Pos
func StrFindAfterPos(HayStack, Needle string, Pos int) int {
  if Pos < 0 || Pos >= len(HayStack) {
    return -1
  }
  return strings.Index(HayStack[Pos:], Needle) + Pos
}

// Find the first occurrence of a character in a string
func StrFindFirstChar(Str1 string, c rune) int {
  Idx := strings.IndexRune(Str1, c)
  if Idx == -1 {
    return 0
  }
  return Idx
}

// Find one of the characters specified in Needle in the HayStack
func StrFindOneOf(HayStack, Needle string) int {
  if HayStack == "" || Needle == "" {
    return -1
  }
  for i, r := range HayStack {
    if strings.ContainsRune(Needle, r) {
      return i
    }
  }
  return -1
}

// Return the character in Str1 at Position
func StrGetAt(Str1 string, Position int) byte {
  if Position < 0 || Position >= len(Str1) {
    return 0
  }
  return Str1[Position]
}

// Get the length of a string
func StrGetLength(Str1 string) int {
  return len(Str1)
}

// Get a specific word from a string
func StrGetWord(Str1 string, WordNbr int) string {
  var Word string
  var i    int

  iss := strings.Fields(Str1)
  for _, Word = range iss {
    i++
    if i == WordNbr {
      return Word
    }
  }
  return ""
}

// Get the position and length of a word from String
func StrGetWordPosLen(Str string, WordNbr int) string {
  if StrWordCount(Str) < WordNbr {
    return "0 0"
  }
  for {
    Old := Str
    StrReplace(&Str, "  ", " ")
    if Old == Str {
      break
    }
  }
  Str = StrTrimLeft(Str)
  Str = StrTrimRight(Str)
  Str = " " + Str + " "
  z := StrGetLength(Str)
  i := 0
  x := 0
  y := 0
  for j := 1; j <= WordNbr; j++ {
    x = StrFindAfterPos(Str, " ", i)
    i = x + 1
    y = StrFindAfterPos(Str, " ", i)
    i = y
    if i >= z {
      return "0 0"
    }
  }
  return fmt.Sprintf("%d %d", x+1, y-x-1)
}

// Get the rest of the Words in a string starting with the Word indicated by WordNbr
func StrGetWords(Str1 string, WordNbr int) string {
  iss := strings.Fields(Str1)
  if len(iss) == 0 {
    return ""
  }
  i := 1
  for idx, Word := range iss {
    _ = Word
    i++
    if i == WordNbr {
      if idx+1 >= len(iss) {
        return ""
      }
      Rest := strings.Join(iss[idx+1:], " ")
      Rest = StrTrimLeft(Rest)
      return Rest
    }
  }
  return ""
}

// Insert string Str2 into string Str1 at Position
func StrInsert(Str1 string, Position int, Str2 string) string {
  if Position < 0 {
    Position = 0
  }
  if Position > len(Str1) {
    Position = len(Str1)
  }
  return Str1[:Position] + Str2 + Str1[Position:]
}

// Insert character c into string Str1 at Position
func StrInsertChar(Str1 string, Position int, c byte) string {
  if Position < 0 {
    Position = 0
  }
  if Position > len(Str1) {
    Position = len(Str1)
  }
  return Str1[:Position] + string(c) + Str1[Position:]
}

// Is word 'not in' word list?
func StrIsNotWord(Word, WordList string) bool {
  if StrGetLength(Word) == 0 {
    return true
  }
  n := StrWordCount(WordList)
  Found := false
  for i := 1; i <= n; i++ {
    s := StrGetWord(WordList, i)
    if Word == s {
      Found = true
    }
  }
  return !Found
}

// Check if a word exists in a word list
func StrIsWord(Word, WordList string) bool {
  if len(Word) == 0 {
    // Word is null, so it can't be in word list
    return false
  }
  n := StrWordCount(WordList)
  for i := 1; i <= n; i++ {
    Str1 := StrGetWord(WordList, i)
    if Word == Str1 {
      // Word was found in word list
      return true
    }
  }
  // Word was not found in word list
  return false
}

// Get the left portion of a string
func StrLeft(Str1 string, Length int) string {
  if Length > len(Str1) {
    return Str1
  }
  return Str1[:Length]
}

// Make the first letter of a string lowercase
func StrMakeFirstLower(Str1 string) string {
  if len(Str1) == 0 {
    return Str1
  }
  return strings.ToLower(Str1[:1]) + Str1[1:]
}

// Make the first letter of a string uppercase
func StrMakeFirstUpper(Str1 string) string {
  if len(Str1) == 0 {
    return Str1
  }
  return strings.ToUpper(string(Str1[0])) + Str1[1:]
}

// Lower case the whole string
func StrMakeLower(Str1 string) string {
  return strings.ToLower(Str1)
}

// Uppercase the whole string
func StrMakeUpper(Str1 string) string {
  return strings.ToUpper(Str1)
}

// Remove all occurrences of a character from a string
func StrRemove(Str1 string, c byte) string {
  Out := make([]byte, 0, len(Str1))
  for i := 0; i < len(Str1); i++ {
    if Str1[i] != c {
      Out = append(Out, Str1[i])
    }
  }
  return string(Out)
}

// Replace all occurrences of a substring with another substring in a string
func StrReplace(Str *string, From string, To string) {
  if From == "" {
    return
  }
  StartPos := 0
  for {
    Index := strings.Index((*Str)[StartPos:], From)
    if Index == -1 {
      break
    }
    StartPos += Index
    *Str = (*Str)[:StartPos] + To + (*Str)[StartPos+len(From):]
    StartPos += len(To)
  }
}

// Get the right portion of a string
func StrRight(Str1 string, Length int) string {
  if Str1 == "" {
    return ""
  }
  if Length > len(Str1) {
    return Str1
  }
  return Str1[len(Str1)-Length:]
}

// Replace the character in Str1 at Position
func StrSetAt(Str1 string, Position int, c byte) string {
  if Position < 0 || Position >= len(Str1) {
    return Str1
  }
  b := []byte(Str1)
  b[Position] = c
  return string(b)
}

// Remove leading, trailing, and extra spaces
func StrSqueeze(Str1 string) string {
  // Trim leading and trailing spaces
  Str1 = StrTrimLeft(Str1)
  Str1 = StrTrimRight(Str1)
  // Replace consecutive spaces with a single space
  for strings.Contains(Str1, "  ") {
    Str1 = strings.ReplaceAll(Str1, "  ", " ")
  }
  return Str1
}

// Convert a string to an integer (Provides stoi functionality)
func StrToInt(Str string) int {
  Nbr, err := strconv.Atoi(Str)
  if err != nil {
    return 0
  }
  return Nbr
}

// Remove leading whitespace
func StrTrimLeft(Str1 string) string {
  if Str1 == "" {
    return ""
  }
  First := strings.IndexFunc(Str1, func(r rune) bool {
    return r != ' ' && r != '\r' && r != '\n'
  })
  if First == -1 {
    return ""
  }
  return Str1[First:]
}

// Remove trailing whitespace
func StrTrimRight(Str1 string) string {
  Last := strings.LastIndexFunc(Str1, func(r rune) bool {
    return r != ' ' && r != '\r' && r != '\n'
  })
  if Last == -1 {
    return ""
  }
  return Str1[:Last+1]
}

// Get the string in StrVector1 at Position
func StrVectorGetAt(StrVector1 []string, Position int) string {
  if Position < 0 || Position >= len(StrVector1) {
    return ""
  }
  return StrVector1[Position]
}

// Replace the string in StrVector1 at Position
func StrVectorSetAt(StrVector1 *[]string, Position int, Str1 string) {
  x := len(*StrVector1) - 1
  if Position > x {
    *StrVector1 = append(*StrVector1, Str1)
    return
  }
  (*StrVector1)[Position] = Str1
}

// Count the number of words in a string
func StrWordCount(Str1 string) int {
  T := Str1
  NWords := 0
  if len(T) > 0 && T[len(T)-1] != ' ' {
    NWords = 1
  }
  for s := len(T) - 1; s > 0; s-- {
    if T[s] == ' ' && T[s-1] != ' ' {
      NWords++
    }
  }
  return NWords
}
