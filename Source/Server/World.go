//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      World.go                                              *
// Usage:     Manages world events and mobile movements             *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

// Create 'spawn mobile' events
func CreateSpawnMobileEvents() {
	var ControlMobSpawnFile      *os.File
	var ControlMobSpawnFileName   string
	var Count                     int
	var CurrentTime               int
	var Days                      int
	var EventFile                *os.File
	var EventFileName             string
	var EventTime                 string
	var Hours                     int
	var Limit                     int
	var Minutes                   int
	var MobileId                  string
	var Months                    int
	var RoomId                    string
	var Seconds                   int
	var Weeks                     int
	var WorldMobileFile          *os.File
	var WorldMobileFileName       string
	var Years                     int
	var DirEntries              []os.DirEntry
	var err                       error

	if ChgDir(WORLD_MOBILES_DIR) != nil {
		LogIt("World::CreateSpawnMobileEvents - Change directory to WORLD_MOBILES_DIR failed")
		os.Exit(1)
	}
	DirEntries, err = os.ReadDir("./")
	if err != nil {
		LogIt("World::CreateSpawnMobileEvents - Change directory to WORLD_MOBILES_DIR failed")
		os.Exit(1)
	}
	if ChgDir(HomeDir) != nil {
		LogBuf  = "World::CreateSpawnMobileEvents - Change directory to HomeDir failed: "
		LogBuf += HomeDir
		LogIt(LogBuf)
		os.Exit(1)
	}
	for _, DirEntry := range DirEntries {
		if DirEntry.IsDir() {
			// Skip directories
			continue
		}
		WorldMobileFileName = DirEntry.Name()
		MobileId = StrLeft(WorldMobileFileName, StrGetLength(WorldMobileFileName)-4)
		if MobileId == "ReadMe" {
			continue
		}
		//* Have we already created a spawn event for this MobileId?
		ControlMobSpawnFileName  = CONTROL_MOB_SPAWN_DIR
		ControlMobSpawnFileName += MobileId
		ControlMobSpawnFile, err = os.Open(ControlMobSpawnFileName)
		if err == nil {
			// The NoMoreSpawnEventsFlag is set for this mobile
			ControlMobSpawnFile.Close()
			continue
		}
		//* Check MaxInWorld against actual 'in world' count
		WorldMobileFileName  = WORLD_MOBILES_DIR + WorldMobileFileName
		WorldMobileFile, err = os.Open(WorldMobileFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogBuf  = "World::CreateSpawnMobileEvents - Open World Mobile file failed: "
			LogBuf += WorldMobileFileName
			LogBuf += " - "
			LogBuf += err.Error()
			LogIt(LogBuf)
			os.Exit(1)
		}
		Scanner := bufio.NewScanner(WorldMobileFile)
		Scanner.Scan()
		Stuff = Scanner.Text()
		if StrGetWord(Stuff, 1) != "MaxInWorld:" {
			// World mobile file format error MaxInWorld
			LogIt("World::CreateSpawnMobileEvents - World mobile file format error MaxInWorld")
			os.Exit(1)
		}
		Count = CountMob(MobileId)
		Limit = StrToInt(StrGetWord(Stuff, 2))
		if Count >= Limit {
			// No spawn event needed
			WorldMobileFile.Close()
			continue
		}
		//*******************************
		//* Create 'spawn mobile' event *
		//*******************************
		Scanner.Scan()
		Stuff = Scanner.Text()
		if StrGetWord(Stuff, 1) != "RoomId:" {
			// World mobile file format error RoomId
			LogIt("World::CreateSpawnMobileEvents - World mobile file format error RoomId")
			os.Exit(1)
		}
		RoomId = StrGetWord(Stuff, 2)
		Scanner.Scan()
		Stuff = Scanner.Text()
		if StrGetWord(Stuff, 1) != "Interval:" {
			// World mobile file format error Interval
			LogIt("World::CreateSpawnMobileEvents - World mobile file format error Interval")
			os.Exit(1)
		}
		Seconds = StrToInt(StrGetWord(Stuff, 2)) * 1
		Minutes = StrToInt(StrGetWord(Stuff, 3)) * 60
		Hours   = StrToInt(StrGetWord(Stuff, 4)) * 3600
		Days    = StrToInt(StrGetWord(Stuff, 5)) * 86400
		Weeks   = StrToInt(StrGetWord(Stuff, 6)) * 604800
		Months  = StrToInt(StrGetWord(Stuff, 7)) * 2592000
		Years   = StrToInt(StrGetWord(Stuff, 8)) * 31104000
		CurrentTime  = int(GetTimeSeconds())
		CurrentTime += Seconds
		CurrentTime += Minutes
		CurrentTime += Hours
		CurrentTime += Days
		CurrentTime += Weeks
		CurrentTime += Months
		CurrentTime += Years
		SprintfBuf  := fmt.Sprintf("%d", CurrentTime)
		EventTime    = SprintfBuf
		EventFileName  = CONTROL_EVENTS_DIR
		EventFileName += "M"
		EventFileName += EventTime
		EventFileName += ".txt"
		EventFile, err = os.OpenFile(EventFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// Open for append failed
			LogBuf  = "World::CreateSpawnMobileEvents - Open Events file failed - append: "
			LogBuf += EventFileName
			LogBuf += " - "
			LogBuf += err.Error()
			LogIt(LogBuf)
			os.Exit(1)
		}
		for Count < Limit {
			TmpStr  = MobileId
			TmpStr += " "
			TmpStr += RoomId
			TmpStr += "\r\n"
			EventFile.WriteString(TmpStr + "\n")
			Count++
		}
		EventFile.Close()
		WorldMobileFile.Close()
		// Set the NoMoreSpawnEventsFlag for this mobile
		ControlMobSpawnFile, err = os.Create(ControlMobSpawnFileName)
		if err != nil {
			// Create file failed
			LogBuf  = "World::CreateSpawnMobileEvents - Create Control Mobile Spawn file failed: "
			LogBuf += ControlMobSpawnFileName
			LogBuf += " - "
			LogBuf += err.Error()
			LogIt(LogBuf)
			os.Exit(1)
		}
		ControlMobSpawnFile.Close()
	}
	if ChgDir(HomeDir) != nil {
		LogIt("World::CreateSpawnMobileEvents - Change directory to HomeDir failed")
		os.Exit(1)
	}
}

// Check 'spawn mobile' events
func CheckSpawnMobileEvents() {
	var CheckTime                 string
	var ControlMobSpawnFileName   string
	var DirEntries              []os.DirEntry
	var EventFile                *os.File
	var EventFileName             string
	var EventTime                 string
	var MobileId                  string
	var RoomId                    string
	var err                       error

	Buf = fmt.Sprintf("%d", GetTimeSeconds())
	CheckTime = Buf
	if ChgDir(CONTROL_EVENTS_DIR) != nil {
		// Change directory failed
		LogIt("World::CheckSpawnMobileEvents - Change directory to CONTROL_EVENTS_DIR failed")
		os.Exit(1)
	}
	DirEntries, err = os.ReadDir("./")
	if err != nil {
		LogIt("World::CheckSpawnMobileEvents - Change directory to CONTROL_EVENTS_DIR failed")
		os.Exit(1)
	}
	if ChgDir(HomeDir) != nil {
		LogBuf  = "World::CheckSpawnMobileEvents - Change directory to HomeDir failed: "
		LogBuf += HomeDir
		LogIt(LogBuf)
		os.Exit(1)
	}
	for _, DirEntry := range DirEntries {
		if DirEntry.IsDir() {
			// Skip directories
			continue
		}
		EventFileName = DirEntry.Name()
		if !strings.HasPrefix(EventFileName, "M") {
			// Event files starting with 'M' are 'spawn mobile' events
			continue
		}
		// Is it time for this event
		EventTime = StrLeft(EventFileName, StrGetLength(EventFileName)-4)
		EventTime = StrRight(EventTime, StrGetLength(EventTime)-1)
		if EventTime > CheckTime {
			// Event is in the future, so skip it
			continue
		}
		// Event's time has arrived
		EventFileName  = CONTROL_EVENTS_DIR + EventFileName
		EventFile, err = os.Open(EventFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogBuf  = "World::CheckSpawnMobileEvents - Open Events file failed: "
			LogBuf += EventFileName
			LogBuf += " - "
			LogBuf += err.Error()
			LogIt(LogBuf)
			os.Exit(1)
		}
		Scanner := bufio.NewScanner(EventFile)
		Scanner.Scan()
		Stuff = Scanner.Text()
		for Stuff != "" {
			// Get RoomId, MobileId, then spawn the mob
			MobileId = StrGetWord(Stuff, 1)
			RoomId   = StrGetWord(Stuff, 2)
			SpawnMobile(MobileId, RoomId)
			// Remove the NoMoreSpawnEventsFlag for this mobile
			// This is overkill, attempts to remove same flag over and over
			ControlMobSpawnFileName = CONTROL_MOB_SPAWN_DIR
			ControlMobSpawnFileName += MobileId
			Remove(ControlMobSpawnFileName)
			Scanner.Scan()
			Stuff = Scanner.Text()
		}
		// Event completed, remove it
		EventFile.Close()
		Remove(EventFileName)
	}
	if ChgDir(HomeDir) != nil {
		// Change directory failed
		LogBuf  = "World::CheckSpawnMobileEvents - Change directory to HomeDir failed: "
		LogBuf += HomeDir
		LogIt(LogBuf)
		os.Exit(1)
	}
}

// Handle world events
func Events() {
	CreateSpawnMobileEvents()
	CheckSpawnMobileEvents()
	MakeMobilesMove()
}

// Heal mobiles
func HealMobiles() {
	var DirEntries                []os.DirEntry
	var MobFighting                 bool
	var MobileId                    string
	var MobStatsHitPointsFileName   string
	var PositionOfDot               int
	var RoomId                      string
	var err                         error

	if ChgDir(MOB_STATS_HPT_DIR) != nil {
		// Change directory failed
		LogIt("World::HealMobiles - Change directory to MOB_STATS_HPT_DIR failed")
		os.Exit(1)
	}
	//*************************
	//* Heal no-fighting mobs *
	//*************************
	// Get a list of all MobStats\HitPoints files
	DirEntries, err = os.ReadDir("./")
	if err != nil {
		LogIt("World::HealMobiles - Change directory to MOB_STATS_HPT_DIR failed")
		os.Exit(1)
	}
	for _, DirEntry := range DirEntries {
		if DirEntry.IsDir() {
			// Skip directories
			continue
		}
		MobStatsHitPointsFileName = DirEntry.Name()
		MobileId    = StrLeft(MobStatsHitPointsFileName, StrGetLength(MobStatsHitPointsFileName)-4)
		MobFighting = HealMobilesFightCheck("MobPlayer", MobileId)
		if MobFighting {
			// Mobile is fighting, no heal
			continue
		}
		MobFighting = HealMobilesFightCheck("PlayerMob", MobileId)
		if MobFighting {
			// Mobile is fighting, no heal
			continue
		}
		//*******************
		//* Heal the mobile *
		//*******************
		RoomId = GetMobileRoom(MobileId)
		RemoveMobFromRoom(RoomId, MobileId)
		DeleteMobStats(MobileId)
		PositionOfDot = StrFindFirstChar(MobileId, '.')
		MobileId      = StrLeft(MobileId, PositionOfDot)
		AddMobToRoom(RoomId, MobileId)
	}
	if ChgDir(HomeDir) != nil {
		// Change directory failed
		LogIt("World::HealMobiles - Change directory to HomeDir failed")
		os.Exit(1)
	}
}

// See if mobile is fighting
func HealMobilesFightCheck(Dir, MobileId string) bool {
	var DirEntries        []os.DirEntry
	var MobFighting         bool
	var MobPlayerFile      *os.File
	var MobPlayerFileName   string
	var err error

	MobFighting = false
	if Dir == "MobPlayer" {
		// Checking MobPlayer
		if ChgDir(MOB_PLAYER_DIR) != nil {
			// Change directory failed
			LogIt("World::HealMobilesFightCheck - Change directory to MOB_PLAYER_DIR failed")
			os.Exit(1)
		}
	}
	if Dir == "PlayerMob" {
		// Checking PlayerMob
		if ChgDir(PLAYER_MOB_DIR) != nil {
			// Change directory failed
			LogIt("World::HealMobilesFightCheck - Change directory to PLAYER_MOB_DIR failed")
			os.Exit(1)
		}
	}
	// Get a list of all MobPlayer files
	DirEntries, err = os.ReadDir("./")
	if err != nil {
		TmpStr = "World::HealMobilesFightCheck - Open "
		TmpStr += Dir
		TmpStr += " file failed"
		LogIt(TmpStr)
		os.Exit(1)
	}
	for _, DirEntry := range DirEntries {
		if DirEntry.IsDir() {
			// Skip directories
			continue
		}
		MobPlayerFileName = DirEntry.Name()
		// Set file name based on Dir
		if Dir == "MobPlayer" {
			// Checking MobPlayer
			MobPlayerFileName = MOB_PLAYER_DIR + MobPlayerFileName
		}
		if Dir == "PlayerMob" {
			// Checking PlayerMob
			MobPlayerFileName = PLAYER_MOB_DIR + MobPlayerFileName
		}
		MobPlayerFile, err = os.Open(MobPlayerFileName)
		if err != nil {
			// Failed to open MobPlayer or MobPlayer file
			TmpStr  = "World::HealMobilesFightCheck - Open "
			TmpStr += Dir
			TmpStr += " file failed"
			LogIt(TmpStr)
			os.Exit(1)
		}
		Scanner := bufio.NewScanner(MobPlayerFile)
		Scanner.Scan()
		Stuff = Scanner.Text()
		for Stuff != "" {
			// Read all lines
			if Stuff == MobileId {
				// A match means the mobile is fighting
				MobFighting = true
			}
			Scanner.Scan()
			Stuff = Scanner.Text()
		}
		MobPlayerFile.Close()
	}
	if ChgDir(HomeDir) != nil {
		// Change directory failed
		LogIt("World::HealMobilesFightCheck - Change directory to HomeDir failed")
		os.Exit(1)
	}
	return MobFighting
}

// Yep, believe it or not, this makes the mobs move
func MakeMobilesMove() {
	var RoomMobListFile     *os.File
	var RoomMobListFileName  string
	var RoomMobMoveFile     *os.File
	var RoomMobMoveFileName  string
	var Success1             bool
	var Success2             bool

	//********************************
	//* Check for existance of files *
	//********************************
	Success1 = false
	Success2 = false
	RoomMobListFileName  = CONTROL_DIR
	RoomMobListFileName += "RoomMobList.txt"
	RoomMobListFile, err := os.Open(RoomMobListFileName)
	if err == nil {
		Success1 = true
	}
	RoomMobMoveFileName  = CONTROL_DIR
	RoomMobMoveFileName += "RoomMobMove.txt"
	RoomMobMoveFile, err = os.Open(RoomMobMoveFileName)
	if err == nil {
		Success2 = true
	}
	if Success1 {
		// RoomMobList file exists, but is it empty?
		FileInfo, _ := RoomMobListFile.Stat()
		if FileInfo.Size() == 0 {
			// Nothing in the MobList file
			Success1 = false
			RoomMobListFile.Close()
			Remove(RoomMobListFileName)
		}
	}
	if Success2 {
		// RoomMobMove file exists, but is it empty?
		FileInfo, _ := RoomMobMoveFile.Stat()
		if FileInfo.Size() == 0 {
			// Nothing in the MobMove file
			Success2 = false
			RoomMobMoveFile.Close()
			Remove(RoomMobMoveFileName)
		}
	}
	//***********************************
	//* Determine which file to process *
	//***********************************
	if Success2 {
		// Process RoomMobMove until empty
		MakeMobilesMove3()
		return
	}
	if Success1 {
		// Process RoomMobList until empty, creates RoomMoveMove
		MakeMobilesMove2()
		return
	}
	// Create RoomMobList,neither RoomMobList or RoomMobMove exist
	MakeMobilesMove1()
}

// Build file containing RoomMob file list
func MakeMobilesMove1() {
	var DirEntries          []os.DirEntry
	var RoomMobList         []string
	var RoomMobFileName       string
	var RoomMobListFile      *os.File
	var RoomMobListFileName   string
	var err                   error

	// Open MakeMobList file
	RoomMobListFileName  = CONTROL_DIR
	RoomMobListFileName += "RoomMobList.txt"
	RoomMobListFile, err = os.Create(RoomMobListFileName)
	if err != nil {
		// Failed to open RoomMobMove file
		LogIt("World::MakeMobilesMove1 - Create RoomMobList file failed")
		os.Exit(1)
	}
	if ChgDir(ROOM_MOB_DIR) != nil {
		// Change directory failed
		LogIt("World::MakeMobilesMove1 - Change directory to ROOM_MOB_DIR failed")
		os.Exit(1)
	}
	// Get a list of all RoomMob files
	DirEntries, err = os.ReadDir("./")
	if err != nil {
		LogIt("World::MakeMobilesMove1 - Change directory to ROOM_MOB_DIR failed")
		os.Exit(1)
	}
	for _, DirEntry := range DirEntries {
		if DirEntry.IsDir() {
			// Skip directories
			continue
		}
		RoomMobFileName = DirEntry.Name()
		if StrFind(RoomMobFileName, "Spawn") == -1 {
			// Not a spawn room, Random position in list
			TmpStr = fmt.Sprintf("%05d", rand.Intn(100000))
		} else {
			// Force 'spawn' rooms to be first in list
			TmpStr = "00000"
		}
		TmpStr     += " "
		TmpStr     += RoomMobFileName
		RoomMobList = append(RoomMobList, TmpStr)
	}
	// sort em
	sort.Strings(RoomMobList)
	// Write em
	for _, Item := range RoomMobList {
		TmpStr  = Item
		TmpStr  = StrGetWord(TmpStr, 2)
		TmpStr += "\n"
		RoomMobListFile.WriteString(TmpStr)
	}
	RoomMobListFile.Close()
	if len(RoomMobList) == 0 {
		// No mobiles are moving, MobMove file is empty
		Remove(RoomMobListFileName)
	}
	if ChgDir(HomeDir) != nil {
		// Change to home directory failed
		LogIt("World::MakeMobilesMove1 - Change directory to HomeDir failed")
		os.Exit(1)
	}
}

// Build file containing mobiles to be moved
func MakeMobilesMove2() {
	var ControlMobNoMoveFile     *os.File
	var ControlMobNoMoveFileName  string
	var ExitCount                 int
	var ExitNumber                int
	var ExitToRoomId              string
	var i                         int
	var MobCount                  int
	var MobileId                  string
	var MobileIdCheck             string
	var MobListNotCompleted       bool
	var PositionOfDot             int
	var RandomPct                 int
	var RoomId                    string
	var RoomMobFile              *os.File
	var RoomMobFileName           string
	var RoomMobListFile          *os.File
	var RoomMobListFileName       string
	var RoomMobListTempFile      *os.File
	var RoomMobListTempFileName   string
	var RoomMobMoveFile          *os.File
	var RoomMobMoveFileName       string
	var TimerStart                time.Time
	var TimerStop                 time.Time
	var ValidMobRoomExits         string
	var err                       error

	// Open MakeMobList file
	RoomMobListFileName  = CONTROL_DIR
	RoomMobListFileName += "RoomMobList.txt"
	RoomMobListFile, err = os.Open(RoomMobListFileName)
	if err != nil {
		// Failed to open RoomMobList file
		LogIt("World::MakeMobilesMove1 - Create RoomMobList file failed")
		os.Exit(1)
	}
	// Open MakeMobListTemp file
	RoomMobListTempFileName  = CONTROL_DIR
	RoomMobListTempFileName += "RoomMobListTemp.txt"
	RoomMobListTempFile, err = os.Create(RoomMobListTempFileName)
	if err != nil {
		// Failed to open RoomMobListTemp file
		LogIt("World::MakeMobilesMove2 - Create RoomMobListTemp file failed")
		os.Exit(1)
	}
	// Open RoomMobMove file
	RoomMobMoveFileName  = CONTROL_DIR
	RoomMobMoveFileName += "RoomMobMove.txt"
	RoomMobMoveFile, err = os.Create(RoomMobMoveFileName)
	if err != nil {
		// Failed to open RoomMobMove file
		LogIt("World::MakeMobilesMove2 - Create RoomMobMove file failed")
		os.Exit(1)
	}
	//***************************
	//* Create RoomMobMove file *
	//***************************
	TimerStart = time.Now()
	TimerStop = TimerStart.Add(100 * time.Millisecond)
	Scanner := bufio.NewScanner(RoomMobListFile)
	for Scanner.Scan() {
		RoomMobFileName = Scanner.Text()
		if time.Now().After(TimerStop) {
			// Time to stop so cpu is not maxed
			MobListNotCompleted = true
			RoomMobListTempFile.WriteString(RoomMobFileName + "\n")
			continue
		}
		RoomId           = StrLeft(RoomMobFileName, StrGetLength(RoomMobFileName)-4)
		// Open RoomMob file
		RoomMobFileName  = ROOM_MOB_DIR + RoomMobFileName
		RoomMobFile, err = os.Open(RoomMobFileName)
		if err != nil {
			// No RoomMob file? Really, I guess all the mobs got themselves killed
			continue
		}
		RoomScanner := bufio.NewScanner(RoomMobFile)
		for RoomScanner.Scan() {
			Stuff = RoomScanner.Text()
			if Stuff == "" {
				continue
			}
			// For each mobile in room
			MobCount      = StrToInt(StrGetWord(Stuff, 1))
			MobileId      = StrGetWord(Stuff, 2)
			MobileIdCheck = MobileId
			PositionOfDot = StrFindFirstChar(MobileIdCheck, '.')
			if PositionOfDot > 1 {
				// Mobile is hurt but not fighting
				MobileIdCheck = StrLeft(MobileIdCheck, PositionOfDot)
			}
			//* Is the MobNoMoveFlag set?
			ControlMobNoMoveFileName  = CONTROL_MOB_NOMOVE_DIR
			ControlMobNoMoveFileName += MobileIdCheck
			ControlMobNoMoveFile, err = os.Open(ControlMobNoMoveFileName)
			if err == nil {
				// The MobNoMoveFlag is set for this mobile
				ControlMobNoMoveFile.Close()
			} else {
				// Mobile may move
				for i = 1; i <= MobCount; i++ {
					// For each mobile occurrence
					if StrFind(RoomId, "Spawn") == -1 {
						// Not a spawn room, Get random chance of mob moving
						RandomPct = GetRandomNumber(100)
					} else {
						// Force mobs in 'spawn' rooms to move
						RandomPct = -1
					}
					if RandomPct <= MOB_MOVE_PCT {
						// Mobile is to be moved
						ValidMobRoomExits = GetValidMobRoomExits(RoomId)
						ExitCount = StrCountWords(ValidMobRoomExits)
						if ExitCount > 0 {
							// Mob has at least one exit available
							ExitNumber = GetRandomNumber(ExitCount)
							ExitToRoomId = StrGetWord(ValidMobRoomExits, ExitNumber)
							if ExitToRoomId == "" {
								// Blow up for now, but we should LogThis, not blow up??
								LogIt("ExitToRoomId is blank zz")
								os.Exit(1)
							}
							TmpStr  = MobileId
							TmpStr += " "
							TmpStr += RoomId
							TmpStr += " "
							TmpStr += ExitToRoomId
							TmpStr += "\n"
							RoomMobMoveFile.WriteString(TmpStr)
						}
					}
				}
			}
		}
		RoomMobFile.Close()
	}
	if err := Scanner.Err(); err != nil {
		LogIt("World::MakeMobilesMove2 - Open RoomMobList file failed")
		os.Exit(1)
	}
	// Close files
	RoomMobMoveFile.Close()
	RoomMobListFile.Close()
	RoomMobListTempFile.Close()
	// Done with RoomMobList file, get rid of it
	Remove(RoomMobListFileName)
	if MobListNotCompleted {
		// Time ran out before MobList was completely processed
		err = Rename(RoomMobListTempFileName, RoomMobListFileName)
		if err != nil {
			// If rename fails, log the error and stop execution
			LogIt("World::MakeMobilesMove2 - Rename RoomMobListTemp file failed")
			os.Exit(1)
		}
	} else {
		// MobList was completely processed
		err = Remove(RoomMobListTempFileName)
		if err != nil {
			// If delete fails, log the error and stop execution
			LogIt("World::MakeMobilesMove2 - Remove RoomMobListTemp file failed")
			os.Exit(1)
		}
	}
}

// Yep, believe it or not, this makes the mobs move
func MakeMobilesMove3() {
	var ArriveMsg                string
	var ExitToRoomId             string
	var LeaveMsg                 string
	var MobileDesc1              string
	var MobileId                 string
	var RoomId                   string
	var MobMoveNotCompleted      bool
	var MobStatsFileName         string
	var RoomMobMoveFileName      string
	var RoomMobMoveTempFileName  string
	var TimerStart               time.Time
	var TimerStop                time.Time
	var PositionOfDot            int
	var RoomMobMoveFile         *os.File
	var RoomMobMoveTempFile     *os.File
	var err                      error

	//******************************
	//* Initization and open files *
	//******************************
	MobMoveNotCompleted  = false
	RoomMobMoveFileName  = CONTROL_DIR + "RoomMobMove.txt"
	RoomMobMoveFile, err = os.Open(RoomMobMoveFileName)
	if err != nil {
		// No RoomMobMove file, Ok, who delete the file when I wasn't looking?
		LogIt("World::MakeMobilesMove3 - Open RoomMobMove failed")
		os.Exit(1)
	}
	RoomMobMoveTempFileName = CONTROL_DIR + "RoomMobMoveTemp.txt"
	RoomMobMoveTempFile, err = os.Create(RoomMobMoveTempFileName)
	if err != nil {
		// RoomMobMoveTemp file failed to open
		LogIt("World::MakeMobilesMove3 - Open RoomMobMoveTemp failed")
		os.Exit(1)
	}
	//****************************
	//* Process RoomMobMove file *
	//****************************
	TimerStart = time.Now()
	TimerStop  = TimerStart.Add(100 * time.Millisecond)
	Scanner := bufio.NewScanner(RoomMobMoveFile)
	for Scanner.Scan() {
		Stuff = Scanner.Text()
		if time.Now().After(TimerStop) {
			// Time to stop so cpu is not maxed
			MobMoveNotCompleted = true
			RoomMobMoveTempFile.WriteString(Stuff + "\n")
			continue
		}
		MobileId     = StrGetWord(Stuff, 1)
		RoomId       = StrGetWord(Stuff, 2)
		ExitToRoomId = StrGetWord(Stuff, 3)
		if !IsMobileIdInRoom(RoomId, MobileId) {
			// Mob not in room anymore, prolly get itself killed, so can't be moved
			continue
		}
		MobileDesc1 = GetMobDesc1(MobileId)
		LeaveMsg    = MobileDesc1
		LeaveMsg   += " leaves."
		ArriveMsg   = MobileDesc1
		ArriveMsg  += " arrives."
		RemoveMobFromRoom(RoomId, MobileId)
		AddMobToRoom(ExitToRoomId, MobileId)
		pDnodeSrc   = nil
		pDnodeTgt   = nil
		SendToRoom(RoomId, LeaveMsg)
		SendToRoom(ExitToRoomId, ArriveMsg)
		PositionOfDot = StrFindFirstChar(MobileId, '.')
		if PositionOfDot > 1 {
			// Delete 'MobStats' Room file
			MobStatsFileName  = MOB_STATS_ROOM_DIR
			MobStatsFileName += MobileId
			MobStatsFileName += ".txt"
			err = Remove(MobStatsFileName)
			if err != nil {
				// If file remove fails, log the error and stop execution
				LogIt("World::MakeMobilesMove - Remove MobStats Room file failed")
				os.Exit(1)
			}
			// Write new RoomId into MobStats Room file
			CreateMobStatsFileWrite(MOB_STATS_ROOM_DIR, MobileId, ExitToRoomId)
		}
	}
	if err := Scanner.Err(); err != nil {
		LogIt("World::MakeMobilesMove3 - Open RoomMobMove failed")
		os.Exit(1)
	}
	// Close RoomMobMove files
	RoomMobMoveFile.Close()
	RoomMobMoveTempFile.Close()
	// Done with RoomMobMove file, get rid of it
	err = Remove(RoomMobMoveFileName)
	if err != nil {
		// If file remove fails, log the error and stop execution
		LogIt("World::MakeMobilesMove3 - Remove RoomMobMove file failed")
		os.Exit(1)
	}
	// Check whether or not mobs got moved
	if MobMoveNotCompleted {
		// Time ran out before all the mobs got moved, so moved the rest of them later
		err = Rename(RoomMobMoveTempFileName, RoomMobMoveFileName)
		if err != nil {
			// If rename fails, log the error and stop execution
			LogIt("World::MakeMobilesMove3 - Rename RoomMobMoveTemp file failed")
			os.Exit(1)
		}
	} else {
		// All mobs got moved, delete the temp file
		err = Remove(RoomMobMoveTempFileName)
		if err != nil {
			// If delete fails, log the error and stop execution
			LogIt("World::MakeMobilesMove3 - Remove RoomMobMoveTemp file failed")
			os.Exit(1)
		}
	}
}

// Spawn a mobile so players have something to whack!
func SpawnMobile(MobileId, RoomId string) {
	var LogMessage    string
	var MobileAction  string
	var SpawnMsg      string
	var pMobile      *Mobile

	//********************
	//* Spawn the mobile *
	//********************
	pMobile = IsMobValid(MobileId)
	if pMobile == nil {
		// Very bad, no such mobile
		LogMessage  = "World::SpawnMobile - Mobile not found."
		LogMessage += "\n"
		LogMessage += "MobileId: "
		LogMessage += MobileId
		LogIt(LogMessage)
		os.Exit(1)
	}
	AddMobToRoom(RoomId, MobileId)
	SpawnMsg  = pMobile.Desc1
	SpawnMsg += " suddenly appears!"
	pDnodeSrc = nil
	pDnodeTgt = nil
	SendToRoom(RoomId, SpawnMsg)
	MobileAction = pMobile.Action
	// Clean up
	pMobile = nil
	if StrIsWord("NoMove", MobileAction) {
		SpawnMobileNoMove(MobileId)
	}
}

// Make mobile stand still
func SpawnMobileNoMove(MobileId string) {
	var ControlMobNoMoveFile     *os.File
	var ControlMobNoMoveFileName  string
	var err error

	ControlMobNoMoveFileName  = CONTROL_MOB_NOMOVE_DIR
	ControlMobNoMoveFileName += MobileId
	ControlMobNoMoveFile, err = os.Create(ControlMobNoMoveFileName)
	if err != nil {
		LogIt("World::SpawnMobile - Create Control Mobile NoMove file failed")
		os.Exit(1)
	}
	ControlMobNoMoveFile.Close()
}
