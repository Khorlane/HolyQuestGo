package server

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// BigDog is the main entry point for the OMugs server
func BigDog() {
	var (
		EventTick        int
		MobHealTick      int
		whoIsOnlineTick  int
		stopItFileName   string
		goGoGoFileName   string
	)

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
	MobHealTick = 0
	whoIsOnlineTick = 0
	StateConnections = true
	StateRunning = true
	StateStopping = false
	// rand.Seed(time.Now().UnixNano()) // If Seed is not called, the generator is seeded randomly at program startup.

	if ValErr := ValidateIt("All"); ValErr {
		LogBuf = "OMugs has stopped"
		LogIt(LogBuf)
		CloseLogFile()
		return
	}

	SockOpenPort(PORT_NBR)
	InitDescriptor()

	for StateRunning {
		time.Sleep(MILLI_SECONDS_TO_SLEEP * time.Millisecond)
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

		whoIsOnlineTick++
		if whoIsOnlineTick >= WHO_IS_ONLINE_TICK {
			whoIsOnlineTick = 0
			pWhoIsOnline := NewWhoIsOnline(HomeDir)
			pWhoIsOnline.Destroy()
		}
	}

	ClearDescriptor()
	SockClosePort(PORT_NBR)
	pWhoIsOnline := NewWhoIsOnline(HomeDir)
	pWhoIsOnline.Destroy()
	LogBuf = "OMugs has stopped"
	LogIt(LogBuf)
	CloseLogFile()
}

// ChgDir changes the current working directory
func ChgDir(Dir string) error {
	return os.Chdir(Dir)
}

// Check if a file exists
func FileExist(Name string) bool {
	_, err := os.Stat(Name)
	return !os.IsNotExist(err)
}

// Rename renames a file
func Rename(file1, file2 string) error {
	return os.Rename(file1, file2)
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

// Create a new WhoIsOnline instance
func NewWhoIsOnline(homeDir string) *WhoIsOnline {
	// Placeholder for creating a new WhoIsOnline instance
	return &WhoIsOnline{}
}

// Clean up the WhoIsOnline instance
func (w *WhoIsOnline) Destroy() {
	// Placeholder for destroying the WhoIsOnline instance
}

// Represents the online players
type WhoIsOnline struct {
	// Placeholder for WhoIsOnline fields
}

// Return the number of seconds since the epoch
func GetTimeSeconds() int {
	now := time.Now()
	epoch := now.Unix()
	return int(epoch)
}

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
	return strings.IndexRune(str1, c)
}

// Get the length of a string
func StrGetLength(Str1 string) int {
	return len(Str1)
}

// Get a specific word from a string
func StrGetWord(Str1 string, WordNbr int) string {
	var (
		Word string
		i    int
	)

	iss := strings.Fields(Str1)
	for _, Word = range iss {
		i++
		if i == WordNbr {
			return Word
		}
	}
	return ""
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

// Convert a string to an integer
func StrToInt(Str string) int {
  Nbr, err := strconv.Atoi(Str)
  if err != nil {
    // Handle the error by returning 0 as a default value
    return 0
  }
  return Nbr
}
