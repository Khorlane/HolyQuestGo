//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Mobile.go                                             *
// Usage:     Manages mobile entities in the game world             *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "strconv"
  "strings"
)

// Public variables
var (
  Action      string
  Armor       int
  Attack      string
  Damage      int
  Desc1       string
  Desc2       string
  Desc3       string
  ExpPoints   int
  Faction     string
  HitPoints   int
  Hurt        bool
  Level       int
  Loot        string
  MobileFile *os.File
  MobileId    string
  MobNbr      string
  Names       string
  Talk        string
)

// Mobile struct definition
type Mobile struct {
  Action      string
  Armor       int
  Attack      string
  Damage      int
  Desc1       string
  Desc2       string
  Desc3       string
  ExpPoints   int
  Faction     string
  HitPoints   int
  Hurt        bool
  Level       int
  Loot        string
  MobileFile *os.File
  MobileId    string
  MobNbr      string
  Names       string
  Talk        string
  Stuff       string
}

var mobileFileName string
var MobScanner = bufio.NewScanner(MobileFile)

// Initialize a new Mobile instance
func MobileConstructor(MobileId string) *Mobile {
  pMobile := &Mobile{}
  OpenMobFile(MobileId)
  ParseMobStuff()
  CloseMobFile()
  pMobile.Hurt = false
  pMobile.MobNbr = ""
  return pMobile
}

// Add a mobile to a room
func AddMobToRoom(RoomId, MobileId string) {
  var BytesInFile        int
  var NewRoomMobFile     bool
  var MobCount           int
  var MobileIdAdded      bool
  var MobileIdCheck      string
  var RoomMobFileName    string
  var RoomMobTmpFileName string

  UpdateMobInWorld(MobileId, "add")
  MobileId = strings.ToLower(MobileId)
  // Open RoomMob file
  RoomMobFileName = ROOM_MOB_DIR + RoomId + ".txt"
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    NewRoomMobFile = true
  }
  defer RoomMobFile.Close()
  // Open temp RoomMob file
  if RoomId == "" {
    LogIt("AddMobToRoom - RoomId is blank")
    os.Exit(1)
  }
  RoomMobTmpFileName = ROOM_MOB_DIR + RoomId + ".tmp.txt"
  RoomMobTmpFile, err := os.Create(RoomMobTmpFileName)
  if err != nil {
    LogIt("AddMobToRoom - Open RoomMob temp file failed")
    os.Exit(1)
  }
  defer RoomMobTmpFile.Close()
  if NewRoomMobFile {
    // New room mobile file, write the mobile and return
    TmpStr := "1 " + MobileId + "\n"
    RoomMobTmpFile.WriteString(TmpStr)
    os.Rename(RoomMobTmpFileName, RoomMobFileName)
    return
  }
  // Write temp RoomMob file
  Scanner := bufio.NewScanner(RoomMobFile)
  for Scanner.Scan() {
    Stuff := Scanner.Text()
    if MobileIdAdded {
      // New mobile has been written, just write the rest of the mobiles
      RoomMobTmpFile.WriteString(Stuff + "\n")
      continue
    }
    MobileIdCheck = StrGetWord(Stuff, 2)
    if MobileId < MobileIdCheck {
      // Add new mobile in alphabetical order
      TmpStr := "1 " + MobileId + "\n"
      RoomMobTmpFile.WriteString(TmpStr)
      MobileIdAdded = true
      RoomMobTmpFile.WriteString(Stuff + "\n")
      continue
    }
    if MobileId == MobileIdCheck {
      // Existing mobile same as new mobile, add 1 to count
      MobCount, _ = strconv.Atoi(StrGetWord(Stuff, 1))
      MobCount++
      TmpStr := strconv.Itoa(MobCount) + " " + MobileId + "\n"
      RoomMobTmpFile.WriteString(TmpStr)
      MobileIdAdded = true
      continue
    }
    // None of the above conditions satisfied, just write it
    RoomMobTmpFile.WriteString(Stuff + "\n")
  }
  if !MobileIdAdded {
    // New mobile goes at the end
    TmpStr := "1 " + MobileId + "\n"
    RoomMobTmpFile.WriteString(TmpStr)
  }
  BytesInFile = StrGetLength(RoomMobTmpFileName) // TODO - steve - What is this doing?
  RoomMobFile.Close()
  RoomMobTmpFile.Close()
  os.Remove(RoomMobFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    os.Rename(RoomMobTmpFileName, RoomMobFileName)
  } else {
    // If the file is empty, delete it and abort ... it should never be empty
    os.Remove(RoomMobTmpFileName)
    LogIt("AddMobToRoom - RoomMob file size is not > 0!!")
    os.Exit(1)
  }
}

// Count the number of a specific mobile in the world
func CountMob(MobileId string) int {
  var MobInWorldCount int

  // Open Mobile InWorld file
  MobInWorldFileName := CONTROL_MOB_INWORLD_DIR + MobileId + ".txt"
  MobInWorldFile, err := os.Open(MobInWorldFileName)
  if err == nil {
    // Get current count
    Scanner := bufio.NewScanner(MobInWorldFile)
    if Scanner.Scan() {
      MobInWorldCount, _ = strconv.Atoi(Scanner.Text())
    }
    MobInWorldFile.Close()
  } else {
    // No file, so count is zero
    MobInWorldCount = 0
  }
  return MobInWorldCount
}

// Create a mobile player file
func CreateMobPlayer(PlayerName, MobileId string) {
  NewFile := true
  MobPlayerFileName := MOB_PLAYER_DIR + PlayerName + ".txt"
  // Check if file exists
  MobPlayerFile, err := os.OpenFile(MobPlayerFileName, os.O_RDWR|os.O_CREATE, 0666)
  if err != nil {
    LogIt("CreateMobPlayer - Open MobPlayer file failed 1")
    os.Exit(1)
  }
  defer MobPlayerFile.Close()
  FileInfo, _ := MobPlayerFile.Stat()
  if FileInfo.Size() > 0 {
    NewFile = false
  }
  if !NewFile {
    MobPlayerFile.Seek(0, io.SeekEnd)
    MobPlayerFile.WriteString("\r\n")
  }
  MobPlayerFile.WriteString(MobileId + "\n")
}

// Write to a mobile statistics file
func CreateMobStatsFileWrite(Directory, MobileIdForFight, Stuff string) {
  MobStatsFileName := Directory + MobileIdForFight + ".txt"
  MobStatsFile, err := os.Create(MobStatsFileName)
  if err != nil {
    LogIt("CreateMobStatsFileWrite - Open for " + Directory + " " + MobileIdForFight + " failed.")
    os.Exit(1)
  }
  defer MobStatsFile.Close()
  MobStatsFile.WriteString(Stuff + "\n")
}

// Create a player-mob relationship file
func CreatePlayerMob(PlayerName, MobileId string) {
  PlayerMobFileName := PLAYER_MOB_DIR + PlayerName + ".txt"
  PlayerMobFile, err := os.Create(PlayerMobFileName)
  if err != nil {
    LogIt("CreatePlayerMob - Open PlayerMob file failed")
    os.Exit(1)
  }
  defer PlayerMobFile.Close()
  PlayerMobFile.WriteString(MobileId + "\n")
}

// Delete a player-mob relationship file
func DeleteMobPlayer(PlayerName, MobileId string) {
  MobileId = strings.ToLower(MobileId)
  // Open MobPlayer file
  MobPlayerFileName := MOB_PLAYER_DIR + PlayerName + ".txt"
  MobPlayerFile, err := os.Open(MobPlayerFileName)
  if err != nil {
    // MobPlayer player file does not exist
    return
  }
  defer MobPlayerFile.Close()
  if MobileId == "file" {
    // Delete the file
    os.Remove(MobPlayerFileName)
    return
  }
  // Open temp MobPlayer file
  MobPlayerFileNameTmp := MOB_PLAYER_DIR + PlayerName + ".tmp.txt"
  MobPlayerFileTmp, err := os.Create(MobPlayerFileNameTmp)
  if err != nil {
    LogIt("DeleteMobPlayer - Open MobPlayer temp file failed")
    os.Exit(1)
  }
  defer MobPlayerFileTmp.Close()
  // Write temp MobPlayer file
  MobileIdDeleted := false
  Scanner := bufio.NewScanner(MobPlayerFile)
  for Scanner.Scan() {
    Stuff := Scanner.Text()
    if MobileIdDeleted {
      // Mobile has been deleted, just write the rest of the mobiles
      MobPlayerFileTmp.WriteString(Stuff + "\n")
      continue
    }
    MobileIdCheck := strings.ToLower(StrGetWord(Stuff, 1))
    if MobileId == MobileIdCheck {
      // Found it, delete it
      MobileIdDeleted = true
      continue
    }
    // None of the above conditions satisfied, just write it
    MobPlayerFileTmp.WriteString(Stuff + "\n")
  }
  FileInfo, err := MobPlayerFileTmp.Stat()
  var BytesInFile int64
  if err == nil {
    BytesInFile = FileInfo.Size()
  } else {
    BytesInFile = 0
  }
  os.Remove(MobPlayerFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    os.Rename(MobPlayerFileNameTmp, MobPlayerFileName)
  } else {
    // If the file is empty, delete it
    os.Remove(MobPlayerFileNameTmp)
  }
}

// Delete mobile statistics files
func DeleteMobStats(MobileId string) {
  // Delete 'MobStats' Armor file
  MobStatsFileName := MOB_STATS_ARM_DIR + MobileId + ".txt"
  os.Remove(MobStatsFileName)
  // Delete 'MobStats' Attack file
  MobStatsFileName = MOB_STATS_ATK_DIR + MobileId + ".txt"
  os.Remove(MobStatsFileName)
  // Delete 'MobStats' Damage file
  MobStatsFileName = MOB_STATS_DMG_DIR + MobileId + ".txt"
  os.Remove(MobStatsFileName)
  // Delete 'MobStats' Desc1 file
  MobStatsFileName = MOB_STATS_DSC_DIR + MobileId + ".txt"
  os.Remove(MobStatsFileName)
  // Delete 'MobStats' ExpPoints file
  MobStatsFileName = MOB_STATS_EXP_DIR + MobileId + ".txt"
  os.Remove(MobStatsFileName)
  // Delete 'MobStats' HitPoints file
  MobStatsFileName = MOB_STATS_HPT_DIR + MobileId + ".txt"
  os.Remove(MobStatsFileName)
  // Delete 'MobStats' Loot file
  MobStatsFileName = MOB_STATS_LOOT_DIR + MobileId + ".txt"
  os.Remove(MobStatsFileName)
  // Delete 'MobStats' Room file
  MobStatsFileName = MOB_STATS_ROOM_DIR + MobileId + ".txt"
  os.Remove(MobStatsFileName)
}

// Delete player-mob relationship file
func DeletePlayerMob(PlayerName string) {
  // Delete 'PlayerMob' file
  PlayerMobFileName := PLAYER_MOB_DIR + PlayerName + ".txt"
  os.Remove(PlayerMobFileName)
}

// Check if a mobile is in the room by its name
func IsMobInRoom(mobileName string) *Mobile {
  var pMobile         *Mobile
  var NamesCheck       string
  var MobileHurt       bool
  var MobileId         string
  var MobNbr           string
  var PositionOfDot    int
  var RoomMobFile     *os.File
  var RoomMobFileName  string
  var Stuff            string

  // Open RoomMob file
  RoomMobFileName = ROOM_MOB_DIR + pDnodeActor.pPlayer.RoomId + ".txt"
  // Try matching using MobileId
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    // Room has no mobiles
    return nil
  }
  defer RoomMobFile.Close()
  scanner := bufio.NewScanner(RoomMobFile)
  for scanner.Scan() {
    Stuff = scanner.Text()
    if Stuff == "" {
      continue
    }
    MobileId = StrGetWord(Stuff, 2)
    if MobileId == mobileName {
      // This mobile is a match
      PositionOfDot = StrFindFirstChar(MobileId, '.')
      MobileHurt = false
      if PositionOfDot > 1 {
        // Mobile is hurt but not fighting
        MobileHurt = true
        MobNbr = StrRight(MobileId, StrGetLength(MobileId)-PositionOfDot-1)
        MobileId = StrLeft(MobileId, PositionOfDot)
      }
      pMobile = &Mobile{MobileId: MobileId, Hurt: MobileHurt, MobNbr: MobNbr}
      return pMobile
    }
  }
  // No match found, try getting match using 'names'
  RoomMobFile.Seek(0, 0) // Reset file pointer
  scanner = bufio.NewScanner(RoomMobFile)
  for scanner.Scan() {
    Stuff = scanner.Text()
    if Stuff == "" {
      continue
    }
   MobileId = StrGetWord(Stuff, 2)
    PositionOfDot = StrFindFirstChar(MobileId, '.')
    MobileHurt = false
    if PositionOfDot > 1 {
      // Mobile is hurt but not fighting
      MobileHurt = true
      MobNbr = StrRight(MobileId, StrGetLength(MobileId)-PositionOfDot-1)
      MobileId = StrLeft(MobileId, PositionOfDot)
    }
    pMobile = &Mobile{MobileId: MobileId, Hurt: MobileHurt, MobNbr: MobNbr}
    if pMobile.Hurt {
      // Mobile is hurt
      if MobNbr == mobileName {
        // Kill nnn was entered, where nnn is the MobNbr
        return pMobile
      }
    }
    NamesCheck = pMobile.Names
    NamesCheck = StrMakeLower(NamesCheck)
    if StrIsWord(mobileName, NamesCheck) {
      // This mobile is a match
      return pMobile
    }
  }
  return nil
}

// Get the first description of a mobile
func GetMobDesc1(MobileId string) string {
  PositionOfDot := strings.Index(MobileId, ".")
  if PositionOfDot > 1 {
    // Mobile is hurt but not fighting
    MobileId = MobileId[:PositionOfDot]
  }
  MobileFileName := MOBILES_DIR + MobileId + ".txt"
  MobileFile, err := os.Open(MobileFileName)
  if err != nil {
    LogIt("GetMobDesc1 - Mobile does not exist!")
    return ""
  }
  defer MobileFile.Close()
  Scanner := bufio.NewScanner(MobileFile)
  for Scanner.Scan() {
    Stuff := Scanner.Text()
    if strings.HasPrefix(Stuff, "Desc1:") {
      Desc1 := strings.TrimSpace(Stuff[6:])
      return Desc1
    }
  }
  return ""
}

// Check if a mobile ID is in a room
func IsMobileIdInRoom(RoomId, MobileId string) bool {
  RoomMobFileName := ROOM_MOB_DIR + RoomId + ".txt"
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    // Room has no mobiles
    return false
  }
  defer RoomMobFile.Close()
  Scanner := bufio.NewScanner(RoomMobFile)
  for Scanner.Scan() {
    Stuff := Scanner.Text()
    MobileIdCheck := StrGetWord(Stuff, 2)
    if MobileId == MobileIdCheck {
      // Found matching mobile
      return true
    }
  }
  // No matching mobile found
  return false
}

// Check if a mobile is valid by its ID
func IsMobValid(mobileId string) *Mobile {
  var pMobile       *Mobile
  var mobileFileName string

  mobileFileName = MOBILES_DIR + mobileId + ".txt"
  if FileExist(mobileFileName) {
    pMobile = &Mobile{MobileId: mobileId}
    return pMobile
  } else {
    return nil
  }
}

// Put a mobile back in the room
func PutMobBackInRoom(PlayerName, RoomIdBeforeFleeing string) {
   var MobHitPointsLeft          string
   var MobHitPointsTotal         string
   var MobileId                  string
   var MobPlayerFileName         string
   var MobStatsHitPointsFileName string
   var PositionOfDot             int

  // Open MobPlayer file
  MobPlayerFileName = MOB_PLAYER_DIR + PlayerName + ".txt"
  MobPlayerFile, err := os.Open(MobPlayerFileName)
  if err != nil {
    // No mobiles to put back, someone else may be fighting the mob
    return
  }
  defer MobPlayerFile.Close()
  // For each mobile still in MobPlayer file (non-fighting mobiles), put it back in room
  Scanner := bufio.NewScanner(MobPlayerFile)
  for Scanner.Scan() {
    Stuff := Scanner.Text()
    MobileId = StrMakeLower(StrGetWord(Stuff, 1))
    // Read mobile stats hit points file
    MobStatsHitPointsFileName = MOB_STATS_HPT_DIR + MobileId + ".txt"
    MobStatsHitPointsFile, err := os.Open(MobStatsHitPointsFileName)
    if err != nil {
      LogBuf = "Mobile::PutMobBackInRoom - Open MobStatsHitPointsFile file failed (read)"
      LogIt(LogBuf)
      os.Exit(1)
    }
    ScannerStats := bufio.NewScanner(MobStatsHitPointsFile)
    ScannerStats.Scan()
    Stuff = ScannerStats.Text()
    MobStatsHitPointsFile.Close()
    MobHitPointsTotal = StrGetWord(Stuff, 1)
    MobHitPointsLeft = StrGetWord(Stuff, 2)
    if MobHitPointsTotal == MobHitPointsLeft {
      // Mobile is not hurt
      DeleteMobStats(MobileId)
      PositionOfDot = StrFindFirstChar(MobileId, '.')
      if PositionOfDot > 1 {
        // Get MobileId
        MobileId = StrLeft(MobileId, PositionOfDot)
      }
    }
    AddMobToRoom(RoomIdBeforeFleeing, MobileId)
    UpdateMobInWorld(MobileId, "remove")
  }
  // Remove MobPlayer file
  os.Remove(MobPlayerFileName)
}

// Remove a mobile from room
func RemoveMobFromRoom(RoomId, MobileId string) {
  var BytesInFile int64
  var MobileIdRemoved bool
  var MobileIdCheck string
  var MobCount int
  var RoomMobFileName string
  var RoomMobTmpFileName string

  UpdateMobInWorld(MobileId, "remove")
  MobileId = StrMakeLower(MobileId)
  // Open RoomMob file
  RoomMobFileName = ROOM_MOB_DIR + RoomId + ".txt"
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    LogBuf = "Mobile::RemoveMobFromRoom - Open RoomMob file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  defer RoomMobFile.Close()
  // Open temp RoomMob file
  RoomMobTmpFileName = ROOM_MOB_DIR + RoomId + ".tmp.txt"
  if RoomId == "" {
    LogBuf = "RoomId is blank 2"
    LogIt(LogBuf)
    os.Exit(1)
  }
  RoomMobTmpFile, err := os.Create(RoomMobTmpFileName)
  if err != nil {
    LogBuf = "Mobile::RemoveMobFromRoom - Open RoomMob temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  defer RoomMobTmpFile.Close()
  // Write temp RoomMob file
  MobileIdRemoved = false
  Scanner := bufio.NewScanner(RoomMobFile)
  for Scanner.Scan() {
    Stuff := Scanner.Text()
    if MobileIdRemoved {
      // Mobile has been removed, just write the rest of the mobiles
      RoomMobTmpFile.WriteString(Stuff + "\n")
      continue
    }
    MobileIdCheck = StrGetWord(Stuff, 2)
    if MobileId == MobileIdCheck {
      // Found it, subtract 1 from count
      MobCount, _ = strconv.Atoi(StrGetWord(Stuff, 1))
      MobCount--
      MobileIdRemoved = true
      if MobCount > 0 {
        TmpStr := fmt.Sprintf("%d %s", MobCount, MobileId)
        RoomMobTmpFile.WriteString(TmpStr + "\n")
      }
      continue
    }
    // None of the above conditions satisfied, just write it
    RoomMobTmpFile.WriteString(Stuff + "\n")
  }
  if !MobileIdRemoved {
    // Mobile not removed, this is definitely BAD!
    LogBuf = "Mobile::RemoveMobFromRoom - Mobile not removed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  FileInfo, err := RoomMobTmpFile.Stat()
  if err == nil {
    BytesInFile = FileInfo.Size()
  } else {
    BytesInFile = 0
  }
  RoomMobFile.Close()
  RoomMobTmpFile.Close()
  os.Remove(RoomMobFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    os.Rename(RoomMobTmpFileName, RoomMobFileName)
  } else {
    // If the file is empty, delete it
    os.Remove(RoomMobTmpFileName)
  }
}

// Show room mobiles
func ShowMobsInRoom(pDnode *Dnode) {
  var pMobile              *Mobile
  var i, j                  int
  var MobileCount           string
  var MobileHurt            bool
  var MobileId              string
  var MobileIdsToBeRemoved  string
  var MobileIdHurt          string
  var MobNbr                string
  var PositionOfDot         int
  var RemoveMobCount        int
  var RoomMobFileName       string

  // Open RoomMob file
  RoomMobFileName = ROOM_MOB_DIR + pDnode.pPlayer.RoomId + ".txt"
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    // No mobiles in room to display
    return
  }
  defer RoomMobFile.Close()
  Scanner := bufio.NewScanner(RoomMobFile)
  for Scanner.Scan() {
    Stuff := Scanner.Text()
    MobileCount = StrGetWord(Stuff, 1)
    MobileId = StrGetWord(Stuff, 2)
    PositionOfDot = StrFindFirstChar(MobileId, '.')
    MobileHurt = false
    if PositionOfDot > 1 {
      // Mobile is hurt but not fighting
      MobileHurt = true
      MobileIdHurt = MobileId
      MobNbr = StrRight(MobileId, StrGetLength(MobileId)-PositionOfDot-1)
      MobileId = StrLeft(MobileId, PositionOfDot)
    }
    pMobile = MobileConstructor(MobileId)
    pMobile.Hurt = MobileHurt
    pMobile.MobNbr = MobNbr
    if MobileHurt {
      // Mobile is hurt
      pDnode.PlayerOut += "\r\n"
      pDnode.PlayerOut += "&WYou see " + pMobile.Desc1 + ", &Mwounded&W, trying to hide."
      pDnode.PlayerOut += "&B (" + MobileIdHurt + ")&N"
    } else {
      // Mobile is not hurt
      pDnode.PlayerOut += "\r\n"
      pDnode.PlayerOut += "&W(" + MobileCount + ") " + pMobile.Desc2 + "&N"
    }
    // Check for AGGRO mobs
    if StrIsWord("Aggro", pMobile.Action) {
      // Attack player
      j, _ = strconv.Atoi(MobileCount)
      for i = 1; i <= j; i++ {
        MobileIdsToBeRemoved += MobAttacks(pMobile) + " "
      }
    }
  }
  // Remove mobs, that attacked a player, from room
  RemoveMobCount = StrCountWords(MobileIdsToBeRemoved)
  for i = 1; i <= RemoveMobCount; i++ {
    MobileId = StrGetWord(MobileIdsToBeRemoved, i)
    RemoveMobFromRoom(pDnode.pPlayer.RoomId, MobileId)
  }
}

// Handle the logic for a mobile attacking a player
func MobAttacks(pMobile *Mobile) string {
  var KillMsg             string
  var MobileId            string
  var MobileIdToBeRemoved string
  var PhraseAll           string
  var PhrasePlayer        string
  var PlayerName          string
  var RoomId              string

  PlayerName = pDnodeActor.PlayerName
  RoomId = pDnodeActor.pPlayer.RoomId
  // Determine appropriate phrase
  if !pDnodeActor.PlayerStateFighting {
    // Phrases for starting a fight
    PhrasePlayer = " starts a fight with you!"
    PhraseAll = " starts a fight with "
  } else {
    // Phrases for mob attacking a player already fighting
    PhrasePlayer = " attacks you!"
    PhraseAll = " attacks "
  }
  // Send message to player
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "&R"
  pDnodeActor.PlayerOut += pMobile.Desc1
  pDnodeActor.PlayerOut += PhrasePlayer
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += pDnodeActor.pPlayer.Output
  // Send message to room
  KillMsg = "&R"
  KillMsg += pMobile.Desc1
  KillMsg += PhraseAll
  KillMsg += PlayerName
  KillMsg += "!"
  KillMsg += "&N"
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(RoomId, KillMsg)
  // Start a fight
  if !pMobile.Hurt {
    // Mobile not hurt
    GetNextMobNbr()
    CreateMobStatsFile(RoomId)
    MobileId = pMobile.MobileId
    MobileIdToBeRemoved = MobileId
    MobileId = pMobile.MobileId + "." + pMobile.MobNbr
  } else {
    // Mobile is hurt
    MobileId = pMobile.MobileId + "." + pMobile.MobNbr
    MobileIdToBeRemoved = MobileId
  }
  if !pDnodeActor.PlayerStateFighting {
    // Set player and mobile to fight
    CreatePlayerMob(PlayerName, MobileId)
    CreateMobPlayer(PlayerName, MobileId)
    pDnodeActor.PlayerStateFighting = true
  } else {
    // Player is fighting, this mob is an 'add'
    CreateMobPlayer(PlayerName, MobileId)
  }
  return MobileIdToBeRemoved
}

// Search all rooms for a specific mobile        
func WhereMob(mobileIdSearch string) {
  var FileName        string
  var MobileHurt      bool
  var MobileId        string
  var PositionOfDot   int
  var RoomMobFileName string
  var RoomName        string

  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Mobiles"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "-------"
  pDnodeActor.PlayerOut += "\r\n"
  if ChgDir(ROOM_MOB_DIR) != nil {
    logBuf := "Mobile::WhereMob - Change directory to ROOM_MOB_DIR failed"
    LogIt(logBuf)
    os.Exit(1)
  }
  entries, err := os.ReadDir("./")
  if err != nil {
    logBuf := "Mobile::WhereMob - Failed to read ROOM_MOB_DIR"
    LogIt(logBuf)
    os.Exit(1)
  }
  for _, entry := range entries {
    if entry.IsDir() {
      continue
    }
    FileName = entry.Name()
    RoomMobFileName = FileName

    fileContent, err := os.ReadFile(RoomMobFileName)
    if err != nil {
      logBuf := "Mobile::WhereMob - Open RoomMob file failed"
      LogIt(logBuf)
      os.Exit(1)
    }
    RoomName = StrLeft(FileName, StrGetLength(FileName)-4)
    lines := strings.Split(string(fileContent), "\n")
    for _, line := range lines {
      if line == "" {
        continue
      }
      MobileId = StrGetWord(line, 2)
      PositionOfDot = strings.Index(MobileId, ".")
      MobileHurt = false

      if PositionOfDot > 1 {
        MobileHurt = true
        MobileId = StrLeft(MobileId, PositionOfDot)
      }
      if MobileId == mobileIdSearch {
        pDnodeActor.PlayerOut += RoomName + " "
        if MobileHurt {
          pDnodeActor.PlayerOut += "&R"
        }
        pDnodeActor.PlayerOut += line + "&N\r\n"
      }
    }
  }
  if ChgDir(HomeDir) != nil {
    logBuf := "Mobile::WhereMob - Change directory to HomeDir failed"
    LogIt(logBuf)
    os.Exit(1)
  }
}

// Update the count of a mobile in the world
func UpdateMobInWorld(mobileId string, addRemove string) {
  mobInWorldCount := 0
  // Find the position of the first dot in MobileId
  positionOfDot := strings.Index(mobileId, ".")
  if positionOfDot > 1 {
    // Get MobileId
    mobileId = StrLeft(mobileId, positionOfDot)
  }
  // Open Mobile InWorld file
  mobInWorldFileName := CONTROL_MOB_INWORLD_DIR + mobileId + ".txt"
  fileContent, err := os.ReadFile(mobInWorldFileName)
  if err == nil {
    // Get current count
    mobInWorldCount, _ = strconv.Atoi(strings.TrimSpace(string(fileContent)))
  }
  // Create or overwrite Mobiles InWorld file
  file, err := os.Create(mobInWorldFileName)
  if err != nil {
    logBuf := "Mobile::UpdateMobInWorld - Open Mobiles InWorld file failed for: " + mobInWorldFileName
    LogIt(logBuf)
    return
  }
  defer file.Close()
  if addRemove == "add" {
    // Mobile is being added to the world
    mobInWorldCount++
  } else {
    // Mobile is being removed from the world
    mobInWorldCount--
  }
  // Write the updated count to the file
  file.WriteString(strconv.Itoa(mobInWorldCount) + "\n")
}

// Create mobile statistics file
func CreateMobStatsFile(RoomId string) {
  MobileIdForFight := MobileId + "." + MobNbr
  // HitPoints
  Stuff := fmt.Sprintf("%d %d", HitPoints, HitPoints)
  CreateMobStatsFileWrite(MOB_STATS_HPT_DIR, MobileIdForFight, Stuff)
  // Armor
  Stuff = fmt.Sprintf("%d", Armor)
  CreateMobStatsFileWrite(MOB_STATS_ARM_DIR, MobileIdForFight, Stuff)
  // Attack
  Stuff = Attack
  CreateMobStatsFileWrite(MOB_STATS_ATK_DIR, MobileIdForFight, Stuff)
  // Damage
  Stuff = fmt.Sprintf("%d", Damage)
  CreateMobStatsFileWrite(MOB_STATS_DMG_DIR, MobileIdForFight, Stuff)
  // Desc1
  Stuff = Desc1
  CreateMobStatsFileWrite(MOB_STATS_DSC_DIR, MobileIdForFight, Stuff)
  // ExpPoints
  Stuff = fmt.Sprintf("%d %d", ExpPoints, Level)
  CreateMobStatsFileWrite(MOB_STATS_EXP_DIR, MobileIdForFight, Stuff)
  // Loot
  Stuff = Loot
  CreateMobStatsFileWrite(MOB_STATS_LOOT_DIR, MobileIdForFight, Stuff)
  // Room
  Stuff = RoomId
  CreateMobStatsFileWrite(MOB_STATS_ROOM_DIR, MobileIdForFight, Stuff)
}

// Examine a mobile
func ExamineMob(MobileId string) {
  mobile := &Mobile{}
  OpenMobFile(MobileId)
  defer CloseMobFile()
  MobileFile := mobile.MobileFile
  scanner := bufio.NewScanner(MobileFile)
  for scanner.Scan() {
    Stuff := scanner.Text()
    if Stuff == "Desc3:" {
      break
    }
  }
  // Mobile Description 3
  for scanner.Scan() {
    Stuff := scanner.Text()
    if Stuff == "End Desc3" {
      break
    }
    pDnodeActor.PlayerOut += Stuff + "\r\n"
  }
  pDnodeActor.PlayerOut += "&N"
}

// Get the next mobile number
func GetNextMobNbr() {
  NextMobNbrFileName := CONTROL_DIR + "NextMobileNumber.txt"
  // Read next mobile number file
  NextMobNbrFile, err := os.Open(NextMobNbrFileName)
  if err != nil {
    LogIt("GetNextMobNbr - Open NextMobileNumber file failed (read)")
    os.Exit(1)
  }
  Scanner := bufio.NewScanner(NextMobNbrFile)
  Scanner.Scan()
  Stuff := strings.TrimSpace(Scanner.Text())
  NextMobNbrFile.Close()
  // Increment next mobile number
  NextMobNbrInteger, _ := strconv.Atoi(Stuff)
  NextMobNbrInteger++
  NextMobNbrString := strconv.Itoa(NextMobNbrInteger)
  // Write next mobile number file
  NextMobNbrFile, err = os.Create(NextMobNbrFileName)
  if err != nil {
    LogIt("GetNextMobNbr - Open NextMobileNumber file failed (write)")
    os.Exit(1)
  }
  defer NextMobNbrFile.Close()
  NextMobNbrFile.WriteString(NextMobNbrString + "\n")
  // Set mobile number
  MobNbr = Stuff
}

// Generate a message for a mobile to say
func MobTalk(pMobile *Mobile) string {
  var MobTalkFileName string
  var MobileMsg       string
  var MsgCount        int
  var RndMsgNbr       int

  // Open and read message file
  MobTalkFileName = TALK_DIR + Talk + ".txt"
  mobTalkFile, err := os.Open(MobTalkFileName)
  if err != nil {
    if Talk != "None" {
      // Log failure to open file
      LogBuf = "Mobile::MobTalk - Failed to open " + MobTalkFileName
      LogIt(LogBuf)
      os.Exit(1)
    }
    return "You are ignored.\r\n"
  }
  defer mobTalkFile.Close()
  // Mobile is going to talk
  MobileMsg = "&W" + StrMakeFirstUpper(Desc1) + " says:" + "&N" + "\r\n"
  // Select random message number
  scanner := bufio.NewScanner(mobTalkFile)
  scanner.Scan()
  Stuff := scanner.Text()
  MsgCount, _ = strconv.Atoi(StrGetWord(Stuff, 4))
  RndMsgNbr = GetRandomNumber(MsgCount)
  // Search for selected message number
  scanner.Scan()
  Stuff = scanner.Text()
  for {
    msgNumber, err := strconv.Atoi(StrGetWord(Stuff, 2))
    if err != nil || msgNumber == RndMsgNbr {
      break
    }
    if !scanner.Scan() {
      // End of file and message was not found
      LogBuf = fmt.Sprintf("Mobile::MobTalk - Failed to find message %d %s", RndMsgNbr, MobTalkFileName)
      LogIt(LogBuf)
      return "You are ignored.\r\n"
    }
    Stuff = scanner.Text()
  }
  // Message found
  scanner.Scan()
  Stuff = scanner.Text()
  for Stuff != "End of Message" {
    if !scanner.Scan() {
      // End of file and normal end of message not found
      LogBuf = fmt.Sprintf("Mobile::MobTalk - Unexpected end of file reading message %d %s", RndMsgNbr, MobTalkFileName)
      LogIt(LogBuf)
      return "You are ignored.\r\n"
    }
    MobileMsg += Stuff + "\r\n"
    Stuff = scanner.Text()
  }
  return MobileMsg
}

// Close the mobile file
func CloseMobFile() {
  MobileFile.Close()
}

// Open the file for a given mobile ID
func OpenMobFile(mobileId string) {
  mobileFileName = MOBILES_DIR + mobileId + ".txt"
  file, err := os.Open(mobileFileName)
  if err != nil {
    LogBuf = "Mobile::OpenFile - Mobile does not exist!"
    LogIt(LogBuf)
    os.Exit(1)
  }
  MobileFile = file
}

// Parse mobile stuff
func ParseMobStuff() {
  ReadMobLine()
  for Stuff != "" {
    if StrLeft(Stuff, 9) == "MobileId:" {
      MobileId = StrTrimLeft(StrRight(Stuff, StrGetLength(Stuff)-9))
    } else if StrLeft(Stuff, 6) == "Names:" {
      Names = StrTrimLeft(StrRight(Stuff, StrGetLength(Stuff)-6))
    } else if StrLeft(Stuff, 6) == "Desc1:" {
      Desc1 = StrTrimLeft(StrRight(Stuff, StrGetLength(Stuff)-6))
    } else if StrLeft(Stuff, 6) == "Desc2:" {
      Desc2 = StrTrimLeft(StrRight(Stuff, StrGetLength(Stuff)-6))
    } else if StrLeft(Stuff, 6) == "Desc3:" {
      // Desc3 can be multi-line and is dealt with in 'ExamineMob'
    } else if StrLeft(Stuff, 7) == "Action:" {
      Action = StrTrimLeft(StrRight(Stuff, StrGetLength(Stuff)-7))
    } else if StrLeft(Stuff, 8) == "Faction:" {
      Faction = StrTrimLeft(StrRight(Stuff, StrGetLength(Stuff)-8))
    } else if StrLeft(Stuff, 6) == "Level:" {
      Level, _ = strconv.Atoi(StrRight(Stuff, StrGetLength(Stuff)-6))
    } else if StrLeft(Stuff, 10) == "HitPoints:" {
      HitPoints, _ = strconv.Atoi(StrRight(Stuff, StrGetLength(Stuff)-10))
      HitPoints += Level * MOB_HPT_PER_LEVEL
    } else if StrLeft(Stuff, 6) == "Armor:" {
      Armor, _ = strconv.Atoi(StrRight(Stuff, StrGetLength(Stuff)-6))
      Armor += Level * MOB_ARM_PER_LEVEL
    } else if StrLeft(Stuff, 7) == "Attack:" {
      Attack = StrMakeLower(StrTrimLeft(StrRight(Stuff, StrGetLength(Stuff)-7)))
    } else if StrLeft(Stuff, 7) == "Damage:" {
      Damage, _ = strconv.Atoi(StrRight(Stuff, StrGetLength(Stuff)-7))
      Damage += Level * MOB_DMG_PER_LEVEL
    } else if StrLeft(Stuff, 10) == "ExpPoints:" {
      ExpPoints, _ = strconv.Atoi(StrRight(Stuff, StrGetLength(Stuff)-10))
      ExpPoints += Level * MOB_EXP_PER_LEVEL
    } else if StrLeft(Stuff, 5) == "Loot:" {
      Loot = StrTrimLeft(StrRight(Stuff, StrGetLength(Stuff)-5))
    } else if StrLeft(Stuff, 5) == "Talk:" {
      Talk = StrTrimLeft(StrRight(Stuff, StrGetLength(Stuff)-5))
    }
    ReadMobLine()
  }
}

// Read a line from Mobile file
func ReadMobLine() {
  if MobScanner.Scan() {
    Stuff = StrTrimLeft(MobScanner.Text())
    Stuff = StrTrimRight(Stuff)
  }
}