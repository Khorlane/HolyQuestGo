//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      World.go                                          	  *
// Usage:     Manages world events and mobile movements              *
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
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Create 'spawn mobile' events
func CreateSpawnMobileEvents() {
	if err := ChgDir(WORLD_MOBILES_DIR); err != nil {
		LogIt("CreateSpawnMobileEvents - Change directory to WORLD_MOBILES_DIR failed")
		os.Exit(1)
	}
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		WorldMobileFileName := info.Name()
		MobileID := strings.TrimSuffix(WorldMobileFileName, filepath.Ext(WorldMobileFileName))
		if MobileID == "ReadMe" {
			return nil
		}
		ControlMobSpawnFileName := CONTROL_MOB_SPAWN_DIR + MobileID
		if _, err := os.Stat(ControlMobSpawnFileName); err == nil {
			return nil
		}
		WorldMobileFilePath := WORLD_MOBILES_DIR + WorldMobileFileName
		file, err := os.Open(WorldMobileFilePath)
		if err != nil {
			LogIt("CreateSpawnMobileEvents - Open World Mobile file failed")
			os.Exit(1)
		}
		defer file.Close()
		var Stuff string
		fmt.Fscanln(file, &Stuff)
		if !strings.HasPrefix(Stuff, "MaxInWorld:") {
			LogIt("CreateSpawnMobileEvents - World mobile file format error MaxInWorld")
			os.Exit(1)
		}
		// Placeholder for CountMob and additional logic
		return nil
	})
	if err != nil {
		LogIt("CreateSpawnMobileEvents - Error walking through files")
		os.Exit(1)
	}
	if err := ChgDir(HomeDir); err != nil {
		LogIt("CreateSpawnMobileEvents - Change directory to HomeDir failed")
		os.Exit(1)
	}
}


// Check 'spawn mobile' events
func CheckSpawnMobileEvents() {
	CheckTime := fmt.Sprintf("%d", GetTimeSeconds())
	if err := ChgDir(CONTROL_EVENTS_DIR); err != nil {
		LogIt("CheckSpawnMobileEvents - Change directory to CONTROL_EVENTS_DIR failed")
		os.Exit(1)
	}
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		EventFileName := info.Name()
		if !strings.HasPrefix(EventFileName, "M") {
			return nil
		}
		EventTime := strings.TrimSuffix(EventFileName, filepath.Ext(EventFileName))
		EventTime = EventTime[1:] // Remove the 'M' prefix
		if EventTime > CheckTime {
			return nil
		}
		EventFilePath := CONTROL_EVENTS_DIR + EventFileName
		file, err := os.Open(EventFilePath)
		if err != nil {
			LogIt("CheckSpawnMobileEvents - Open Events file failed")
			os.Exit(1)
		}
		defer file.Close()
		var Stuff string
		for fmt.Fscanln(file, &Stuff); Stuff != ""; fmt.Fscanln(file, &Stuff) {
			MobileID := strings.Fields(Stuff)[0]
			RoomID := strings.Fields(Stuff)[1]
			SpawnMobile(MobileID, RoomID)
		}
		file.Close()
		if err := os.Remove(EventFilePath); err != nil {
			LogIt("CheckSpawnMobileEvents - Remove Events file failed")
			os.Exit(1)
		}
		return nil
	})
	if err != nil {
		LogIt("CheckSpawnMobileEvents - Error walking through files")
		os.Exit(1)
	}
	if err := ChgDir(HomeDir); err != nil {
		LogIt("CheckSpawnMobileEvents - Change directory to HomeDir failed")
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
	if err := ChgDir(MOB_STATS_HPT_DIR); err != nil {
		LogIt("HealMobiles - Change directory to MOB_STATS_HPT_DIR failed")
		os.Exit(1)
	}
	// Get a list of all MobStats\HitPoints files
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		MobStatsHitPointsFileName := info.Name()
		MobileID := strings.TrimSuffix(MobStatsHitPointsFileName, filepath.Ext(MobStatsHitPointsFileName))
		if HealMobilesFightCheck("MobPlayer", MobileID) || HealMobilesFightCheck("PlayerMob", MobileID) {
			return nil
		}
		RoomID := GetMobileRoom(MobileID)
		RemoveMobFromRoom(RoomID, MobileID)
		DeleteMobStats(MobileID)
		PositionOfDot := strings.Index(MobileID, ".")
		if PositionOfDot > 0 {
			MobileID = MobileID[:PositionOfDot]
		}
		AddMobToRoom(RoomID, MobileID)
		return nil
	})
	if err != nil {
		LogIt("HealMobiles - Error walking through files")
		os.Exit(1)
	}
	if err := ChgDir(HomeDir); err != nil {
		LogIt("HealMobiles - Change directory to HomeDir failed")
		os.Exit(1)
	}
}

// See if mobile is fighting
func HealMobilesFightCheck(dir, mobileID string) bool {
	MobFighting := false
	switch dir {
		case "MobPlayer":
			if err := ChgDir(MOB_PLAYER_DIR); err != nil {
				LogIt("HealMobilesFightCheck - Change directory to MOB_PLAYER_DIR failed")
				os.Exit(1)
			}
		case "PlayerMob":
			if err := ChgDir(PLAYER_MOB_DIR); err != nil {
				LogIt("HealMobilesFightCheck - Change directory to PLAYER_MOB_DIR failed")
				os.Exit(1)
			}
	}
	err := filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		MobPlayerFileName := info.Name()
		file, err := os.Open(MobPlayerFileName)
		if err != nil {
			LogIt("HealMobilesFightCheck - Open " + dir + " file failed")
			os.Exit(1)
		}
		defer file.Close()
		var Stuff string
		for fmt.Fscanln(file, &Stuff); Stuff != ""; fmt.Fscanln(file, &Stuff) {
			if Stuff == mobileID {
				MobFighting = true
			}
		}
		return nil
	})
	if err != nil {
		LogIt("HealMobilesFightCheck - Error walking through files")
		os.Exit(1)
	}
	if err := ChgDir(HomeDir); err != nil {
		LogIt("HealMobilesFightCheck - Change directory to HomeDir failed")
		os.Exit(1)
	}
	return MobFighting
}

// Yep, believe it or not, this makes the mobs move
func MakeMobilesMove() {
	RoomMobListFileName := CONTROL_DIR + "RoomMobList.txt"
	RoomMobMoveFileName := CONTROL_DIR + "RoomMobMove.txt"
	RoomMobListFile, err1 := os.Open(RoomMobListFileName)
	if err1 == nil {
		defer RoomMobListFile.Close()
	}
	RoomMobMoveFile, err2 := os.Open(RoomMobMoveFileName)
	if err2 == nil {
		defer RoomMobMoveFile.Close()
	}
	if err1 == nil {
		if stat, _ := RoomMobListFile.Stat(); stat.Size() == 0 {
			RoomMobListFile.Close()
			os.Remove(RoomMobListFileName)
			err1 = os.ErrNotExist
		}
	}
	if err2 == nil {
		if stat, _ := RoomMobMoveFile.Stat(); stat.Size() == 0 {
			RoomMobMoveFile.Close()
			os.Remove(RoomMobMoveFileName)
			err2 = os.ErrNotExist
		}
	}
	if err2 == nil {
		MakeMobilesMove3()
		return
	}
	if err1 == nil {
		MakeMobilesMove2()
		return
	}
	MakeMobilesMove1()
}

// Build file containing RoomMob file list
func MakeMobilesMove1() {
	RoomMobListFileName := CONTROL_DIR + "RoomMobList.txt"
	RoomMobListFile, err := os.Create(RoomMobListFileName)
	if err != nil {
		LogIt("MakeMobilesMove1 - Create RoomMobList file failed")
		os.Exit(1)
	}
	defer RoomMobListFile.Close()
	if err := ChgDir(ROOM_MOB_DIR); err != nil {
		LogIt("MakeMobilesMove1 - Change directory to ROOM_MOB_DIR failed")
		os.Exit(1)
	}
	RoomMobList := []string{}
	err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		RoomMobFileName := info.Name()
		if !strings.Contains(RoomMobFileName, "Spawn") {
			TmpStr := fmt.Sprintf("%05d", rand.Intn(100000))
			TmpStr += " " + RoomMobFileName
			RoomMobList = append(RoomMobList, TmpStr)
		} else {
			TmpStr := "00000 " + RoomMobFileName
			RoomMobList = append(RoomMobList, TmpStr)
		}
		return nil
	})
	if err != nil {
		LogIt("MakeMobilesMove1 - Error walking through files")
		os.Exit(1)
	}
	sort.Strings(RoomMobList)
	for _, item := range RoomMobList {
		TmpStr := strings.Fields(item)[1] + "\n"
		RoomMobListFile.WriteString(TmpStr)
	}
	if len(RoomMobList) == 0 {
		os.Remove(RoomMobListFileName)
	}
	if err := ChgDir(HomeDir); err != nil {
		LogIt("MakeMobilesMove1 - Change directory to HomeDir failed")
		os.Exit(1)
	}
}

// Build file containing mobiles to be moved
func MakeMobilesMove2() {
	RoomMobListFileName := CONTROL_DIR + "RoomMobList.txt"
	RoomMobListTempFileName := CONTROL_DIR + "RoomMobListTemp.txt"
	RoomMobMoveFileName := CONTROL_DIR + "RoomMobMove.txt"
	RoomMobListFile, err := os.Open(RoomMobListFileName)
	if err != nil {
		LogIt("MakeMobilesMove2 - Open RoomMobList file failed")
		os.Exit(1)
	}
	defer RoomMobListFile.Close()
	RoomMobListTempFile, err := os.Create(RoomMobListTempFileName)
	if err != nil {
		LogIt("MakeMobilesMove2 - Create RoomMobListTemp file failed")
		os.Exit(1)
	}
	defer RoomMobListTempFile.Close()
	RoomMobMoveFile, err := os.Create(RoomMobMoveFileName)
	if err != nil {
		LogIt("MakeMobilesMove2 - Create RoomMobMove file failed")
		os.Exit(1)
	}
	defer RoomMobMoveFile.Close()
	TimerStart := time.Now()
	TimerStop := TimerStart.Add(100 * time.Millisecond)
	Scanner := bufio.NewScanner(RoomMobListFile)
	for Scanner.Scan() {
		RoomMobFileName := Scanner.Text()
		if time.Now().After(TimerStop) {
			RoomMobListTempFile.WriteString(RoomMobFileName + "\n")
			continue
		}
		RoomID = strings.TrimSuffix(RoomMobFileName, filepath.Ext(RoomMobFileName))
		RoomMobFilePath := ROOM_MOB_DIR + RoomMobFileName
		RoomMobFile, err := os.Open(RoomMobFilePath)
		if err != nil {
			continue
		}
		defer RoomMobFile.Close()
		RoomScanner := bufio.NewScanner(RoomMobFile)
		for RoomScanner.Scan() {
			Stuff = RoomScanner.Text()
			// Placeholder for processing each mobile in the room
		}
	}
	if err := Scanner.Err(); err != nil {
		LogIt("MakeMobilesMove2 - Error scanning RoomMobList file")
		os.Exit(1)
	}
	if err := os.Remove(RoomMobListFileName); err != nil {
		LogIt("MakeMobilesMove2 - Remove RoomMobList file failed")
		os.Exit(1)
	}
	if err := os.Rename(RoomMobListTempFileName, RoomMobListFileName); err != nil {
		LogIt("MakeMobilesMove2 - Rename RoomMobListTemp file failed")
		os.Exit(1)
	}
}

// Yep, believe it or not, this makes the mobs move
func MakeMobilesMove3() {
	var (
		ArriveMsg, ExitToRoomId, LeaveMsg, MobileDesc1, MobileId, RoomId string
		MobMoveNotCompleted                                      bool
		RoomMobMoveFileName, RoomMobMoveTempFileName             string
		TimerStart, TimerStop                                   time.Time
	)
	MobMoveNotCompleted = false
	RoomMobMoveFileName = CONTROL_DIR + "RoomMobMove.txt"
	RoomMobMoveFile, err := os.Open(RoomMobMoveFileName)
	if err != nil {
		LogIt("MakeMobilesMove3 - Open RoomMobMove failed")
		return
	}
	defer RoomMobMoveFile.Close()
	RoomMobMoveTempFileName = CONTROL_DIR + "RoomMobMoveTemp.txt"
	RoomMobMoveTempFile, err := os.Create(RoomMobMoveTempFileName)
	if err != nil {
		LogIt("MakeMobilesMove3 - Open RoomMobMoveTemp failed")
		return
	}
	defer RoomMobMoveTempFile.Close()
	TimerStart = time.Now()
	TimerStop = TimerStart.Add(100 * time.Millisecond)
	Scanner := bufio.NewScanner(RoomMobMoveFile)
	for Scanner.Scan() {
		Stuff := Scanner.Text()
		if time.Now().After(TimerStop) {
			MobMoveNotCompleted = true
			RoomMobMoveTempFile.WriteString(Stuff + "\n")
			continue
		}
		MobileId = StrGetWord(Stuff, 1)
		RoomId = StrGetWord(Stuff, 2)
		ExitToRoomId = StrGetWord(Stuff, 3)
		if !IsMobileIdInRoom(RoomId, MobileId) {
			continue
		}
		MobileDesc1 = GetMobDesc1(MobileId)
		LeaveMsg = MobileDesc1 + " leaves."
		ArriveMsg = MobileDesc1 + " arrives."
		RemoveMobFromRoom(RoomId, MobileId)
		AddMobToRoom(ExitToRoomId, MobileId)
		SendToRoom(RoomId, LeaveMsg)
		SendToRoom(ExitToRoomId, ArriveMsg)
		PositionOfDot := strings.Index(MobileId, ".")
		if PositionOfDot > 1 {
			MobStatsFileName := MOB_STATS_ROOM_DIR + MobileId + ".txt"
			if err := os.Remove(MobStatsFileName); err != nil {
				LogIt("MakeMobilesMove3 - Remove MobStats Room file failed")
				return
			}
			CreateMobStatsFileWrite(MOB_STATS_ROOM_DIR, MobileId, ExitToRoomId)
		}
	}
	if err := Scanner.Err(); err != nil {
		LogIt("MakeMobilesMove3 - Error scanning RoomMobMove file")
		return
	}
	if err := os.Remove(RoomMobMoveFileName); err != nil {
		LogIt("MakeMobilesMove3 - Remove RoomMobMove file failed")
		return
	}
	if MobMoveNotCompleted {
		if err := os.Rename(RoomMobMoveTempFileName, RoomMobMoveFileName); err != nil {
			LogIt("MakeMobilesMove3 - Rename RoomMobMoveTemp file failed")
			return
		}
	} else {
		if err := os.Remove(RoomMobMoveTempFileName); err != nil {
			LogIt("MakeMobilesMove3 - Remove RoomMobMoveTemp file failed")
			return
		}
	}
}

// Spawn a mobile so players have something to whack!
func SpawnMobile(MobileId, RoomId string) {
	var (
		LogMessage, MobileAction, SpawnMsg string
	)

	//********************
	//* Spawn the mobile *
	//********************
	pMobile := IsMobValid(MobileId)
	if pMobile == nil {
		// Very bad, no such mobile
		LogMessage = "SpawnMobile - Mobile not found.\n"
		LogMessage += "MobileId: " + MobileId
		LogIt(LogMessage)
		return
	}
	AddMobToRoom(RoomId, MobileId)
	SpawnMsg = pMobile.Desc1 + " suddenly appears!"
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
	ControlMobNoMoveFileName := CONTROL_MOB_NOMOVE_DIR + MobileId
	ControlMobNoMoveFile, err := os.Create(ControlMobNoMoveFileName)
	if err != nil {
		LogIt("SpawnMobileNoMove - Create Control Mobile NoMove file failed")
		return
	}
	defer ControlMobNoMoveFile.Close()
}