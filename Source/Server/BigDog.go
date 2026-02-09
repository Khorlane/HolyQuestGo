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
  var goGoGoFileName string
  var MobHealTick    int
  var stopItFileName string

  PrintIt("OMugs Starting")
  HomeDir = GetHomeDir()
  if err := ChgDir(HomeDir); err != nil {
    PrintIt("BigDog() - Change directory to HomeDir failed")
    os.Exit(1)
  }
  stopItFileName = CONTROL_DIR + "StopIt"
  goGoGoFileName = CONTROL_DIR + "GoGoGo"
  if FileExist(stopItFileName) {
    if err := Rename(stopItFileName, goGoGoFileName); err != nil {
      PrintIt("BigDog() - Rename of 'StopIt' to 'GoGoGo' failed!")
      os.Exit(1)
    }
  }
  OpenLogFile()
  LogBuf = "OMugs version " + VERSION + " has started"
  LogIt(LogBuf)
  LogBuf = "Home directory is " + HomeDir
  LogIt(LogBuf)
  EventTick = EVENT_TICK
  MobHealTick      = 0
  StateConnections = true
  StateRunning     = true
  StateStopping    = false
  if ValErr := ValidateIt("All"); ValErr {
    LogBuf = "OMugs has stopped"
    LogIt(LogBuf)
    CloseLogFile()
    return
  }
  SockOpenPort(PORT_NBR)
  InitDescriptor()
  CalendarConstructor()
  for StateRunning {
    time.Sleep(time.Duration(MILLI_SECONDS_TO_SLEEP) * time.Millisecond)
    AdvanceTime()
    if !StateStopping && FileExist(stopItFileName) {
      StateStopping = true
      LogBuf = "Game is stopping"
      LogIt(LogBuf)
    }
    if !StateStopping {
      SockCheckForNewConnections()
      if StateConnections && DnodeCount == 1 {
        LogBuf = "No Connections - going to sleep"
        LogIt(LogBuf)
        StateConnections = false
      }
    }
    if StateConnections {
      SockRecv()
      EventTick++
      if EventTick >= EVENT_TICK {
        EventTick = 0
        Events()
      }
      MobHealTick++
      if MobHealTick >= MOB_HEAL_TICK {
        MobHealTick = 0
        HealMobiles()
      }
    } else if StateStopping {
      StateRunning = false
    }
  }
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
func Rename(file1, file2 string) error {
  return os.Rename(file1, file2)
}

// Remove deletes a file and returns an error value
func Remove(file1 string) error {
  err := os.Remove(file1)
  return err
}

// Return the number of seconds since the epoch
func GetTimeSeconds() int64 {
  now := time.Now()
  epoch := now.Unix()
  return epoch
}

// Print a message to stdout
func PrintIt(message string) {
  message = "\r\n" + message + "\r\n"
  println(message)
}

// Pause execution for the specified duration
func Sleep(milliseconds int) {
  time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}

//-----------------------------------------------------------------------------
// String utility functions
//-----------------------------------------------------------------------------

// Count characters in a string
func StrCountChar(Str1 string, c rune) int {
  count := 0
  for _, char := range Str1 {
    if char == c {
      count++
    }
  }
  return count
}

// Count the number of words in a string
func StrCountWords(str1 string) int {
  // Squeeze spaces in the string
  s1 := StrSqueeze(str1)
  // Split the string into words
  words := strings.Fields(s1)
  // Return the number of words
  return len(words)
}

// Delete characters from a string
func StrDelete(s string, pos int, length int) string {
  if pos < 0 || pos >= len(s) || length <= 0 {
    return s
  }
  end := pos + length
  if end > len(s) {
    end = len(s)
  }
  return s[:pos] + s[end:]
}

// Delete the word specified by WordNbr, squeezing string first
func StrDeleteWord(Str1 string, WordNbr int) string {
  S1 := StrSqueeze(Str1)
  if S1 == "" {
    return ""
  }
  Words := strings.Split(S1, " ")
  index := WordNbr - 1
  if index < 0 || index >= len(Words) {
    return Str1
  }
  Words = append(Words[:index], Words[index+1:]...)
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
func StrFindFirstChar(str1 string, c rune) int {
  idx := strings.IndexRune(str1, c)
  if idx == -1 {
    return 0
  }
  return idx
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
func StrGetWordPosLen(str string, wordNbr int) string {
  if StrWordCount(str) < wordNbr {
    return "0 0"
  }
  for {
    old := str
    StrReplace(&str, "  ", " ")
    if old == str {
      break
    }
  }
  str = StrTrimLeft(str)
  str = StrTrimRight(str)
  str = " " + str + " "
  z := StrGetLength(str)
  i := 0
  x := 0
  y := 0
  for j := 1; j <= wordNbr; j++ {
    x = StrFindAfterPos(str, " ", i)
    i = x + 1
    y = StrFindAfterPos(str, " ", i)
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
      rest := strings.Join(iss[idx+1:], " ")
      rest = StrTrimLeft(rest)
      return rest
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
  found := false
  for i := 1; i <= n; i++ {
    s := StrGetWord(WordList, i)
    if Word == s {
      found = true
    }
  }
  return !found
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
  out := make([]byte, 0, len(Str1))
  for i := 0; i < len(Str1); i++ {
    if Str1[i] != c {
      out = append(out, Str1[i])
    }
  }
  return string(out)
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
