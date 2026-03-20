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
)

// Public variables
var Action      string
var Armor       int
var Attack      string
var Damage      int
var Desc1       string
var Desc2       string
var Desc3       string
var ExpPoints   int
var Faction     string
var HitPoints   int
var Hurt        bool
var Level       int
var Loot        string
var MobileFile *os.File
var MobileId    string
var MobNbr      string
var Names       string
var Talk        string

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

var MobScanner *bufio.Scanner

// Mobile constructor
func MobileConstructor(MobileIdParm string) *Mobile {
  var pMobile *Mobile

  Action    = ""
  Armor     = 0
  Attack    = ""
  Damage    = 0
  Desc1     = ""
  Desc2     = ""
  Desc3     = ""
  ExpPoints = 0
  Faction   = ""
  HitPoints = 0
  Hurt      = false
  Level     = 0
  Loot      = ""
  MobileId  = ""
  MobNbr    = ""
  Names     = ""
  Talk = ""
  OpenMobFile(MobileIdParm)
  ParseMobStuff()
  CloseMobFile()
  pMobile = &Mobile{}
  pMobile.Action     = Action
  pMobile.Armor      = Armor
  pMobile.Attack     = Attack
  pMobile.Damage     = Damage
  pMobile.Desc1      = Desc1
  pMobile.Desc2      = Desc2
  pMobile.Desc3      = Desc3
  pMobile.ExpPoints  = ExpPoints
  pMobile.Faction    = Faction
  pMobile.HitPoints  = HitPoints
  pMobile.Level      = Level
  pMobile.Loot       = Loot
  pMobile.MobileId   = MobileId
  pMobile.Names      = Names
  pMobile.Talk       = Talk
  pMobile.MobileFile = MobileFile
  pMobile.Hurt       = false
  pMobile.MobNbr     = ""
  return pMobile
}

// Add a mobile to a room
func AddMobToRoom(RoomId, MobileId string) {
  var BytesInFile         int64
  var NewRoomMobFile      bool
  var MobCount            int
  var MobileIdAdded       bool
  var MobileIdCheck       string
  var RoomMobFile        *os.File
  var RoomMobFileName     string
  var RoomMobTmpFile     *os.File
  var RoomMobTmpFileName  string

  UpdateMobInWorld(MobileId, "add")
  MobileId = StrMakeLower(MobileId)
  // Open RoomMob file
  RoomMobFileName  = ROOM_MOB_DIR
  RoomMobFileName += RoomId
  RoomMobFileName += ".txt"
  NewRoomMobFile = false
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    NewRoomMobFile = true
  }
  // Open temp RoomMob file
  RoomMobTmpFileName  = ROOM_MOB_DIR
  RoomMobTmpFileName += RoomId
  if RoomId == "" {
    LogBuf = "Mobile::AddMobToRoom - RoomId is blank"
    LogIt(LogBuf)
    os.Exit(1)
  }
  RoomMobTmpFileName += ".tmp.txt"
  RoomMobTmpFile, err = os.Create(RoomMobTmpFileName)
  if err != nil {
    LogBuf = "Mobile::AddMobToRoom - Open RoomMob temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  if NewRoomMobFile {
    // New room mobile file, write the mobile and return
    TmpStr = "1 "
    TmpStr += MobileId
    TmpStr += "\n"
    RoomMobTmpFile.WriteString(TmpStr)
    RoomMobTmpFile.Close()
    Rename(RoomMobTmpFileName, RoomMobFileName)
    return
  }
  // Write temp RoomMob file
  MobileIdAdded = false
  Scanner := bufio.NewScanner(RoomMobFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    if MobileIdAdded {
      // New mobile has been written, just write the rest of the mobiles
      Stuff += "\n"
      RoomMobTmpFile.WriteString(Stuff)
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    MobileIdCheck = StrGetWord(Stuff, 2)
    if MobileId < MobileIdCheck {
      // Add new mobile in alphabetical order
      TmpStr = "1 "
      TmpStr += MobileId
      TmpStr += "\n"
      RoomMobTmpFile.WriteString(TmpStr)
      MobileIdAdded = true
      Stuff += "\n"
      RoomMobTmpFile.WriteString(Stuff)
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    if MobileId == MobileIdCheck {
      // Existing mobile same as new mobile, add 1 to count
      MobCount = StrToInt(StrGetWord(Stuff, 1))
      MobCount++
      Buf = fmt.Sprintf("%d", MobCount)
      TmpStr  = Buf
      TmpStr += " "
      TmpStr += MobileId
      TmpStr += "\n"
      RoomMobTmpFile.WriteString(TmpStr)
      MobileIdAdded = true
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    // None of the above conditions satisfied, just write it
    Stuff += "\n"
    RoomMobTmpFile.WriteString(Stuff)
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  if !MobileIdAdded {
    // New mobile goes at the end
    TmpStr  = "1 "
    TmpStr += MobileId
    TmpStr += "\n"
    RoomMobTmpFile.WriteString(TmpStr)
    MobileIdAdded = true
  }
  BytesInFile = int64(StrGetLength(RoomMobTmpFileName)) // TODO - steve - What is this doing?
  RoomMobFile.Close()
  RoomMobTmpFile.Close()
  Remove(RoomMobFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    Rename(RoomMobTmpFileName, RoomMobFileName)
  } else {
    // If the file is empty, delete it for and abort ... it should never be empty
    Remove(RoomMobTmpFileName)
    LogBuf = "Mobile::AddMobToRoom - RoomMob file size is not > 0!!"
    LogIt(LogBuf)
    os.Exit(1)
  }
}

// Count the number of a specific mobile in the world
func CountMob(MobileId string) int {
  var MobInWorldCount     int
  var MobInWorldFile     *os.File
  var MobInWorldFileName  string

  // Open Mobile InWorld file
  MobInWorldFileName  = CONTROL_MOB_INWORLD_DIR
  MobInWorldFileName += MobileId
  MobInWorldFileName += ".txt"
  MobInWorldFile, err := os.Open(MobInWorldFileName)
  if err == nil {
    // Get current count
    Scanner := bufio.NewScanner(MobInWorldFile)
    Scanner.Scan()
    Stuff = Scanner.Text()
    MobInWorldCount = StrToInt(Stuff)
    MobInWorldFile.Close()
  } else {
    // No file, so count is zero
    MobInWorldCount = 0
  }
  return MobInWorldCount
}

// Create a mobile player file
func CreateMobPlayer(PlayerName, MobileId string) {
  var NewFile            bool
  var MobPlayerFile     *os.File
  var MobPlayerFileName  string

  NewFile = true
  MobPlayerFileName  = MOB_PLAYER_DIR
  MobPlayerFileName += PlayerName
  MobPlayerFileName += ".txt"
  MobPlayerFile, err := os.Open(MobPlayerFileName)
  if err == nil {
    NewFile = false
    MobPlayerFile.Close()
  }
  if NewFile {
    // Create new file
    MobPlayerFile, err = os.Create(MobPlayerFileName)
    if err != nil {
      LogBuf = "Mobile::CreateMobPlayer - Open MobPlayerile file failed 1"
      LogIt(LogBuf)
      os.Exit(1)
    }
  } else {
    // Use existing file
    MobPlayerFile, err = os.OpenFile(MobPlayerFileName, os.O_RDWR, 0666)
    if err != nil {
      LogBuf = "Mobile::CreateMobPlayer - Open MobPlayerile file failed 2"
      LogIt(LogBuf)
      os.Exit(1)
    }
  }
  if !NewFile {
    MobPlayerFile.Seek(0, io.SeekEnd)
    MobPlayerFile.WriteString("\r\n")
  }
  MobPlayerFile.WriteString(MobileId + "\n")
  MobPlayerFile.Close()
}

// Write to a mobile statistics file
func CreateMobStatsFileWrite(Directory, MobileIdForFight, Stuff string) {
  var AfxMessage        string
  var MobStatsFile     *os.File
  var MobStatsFileName  string

  _ = AfxMessage
  MobStatsFileName  = Directory
  MobStatsFileName += MobileIdForFight
  MobStatsFileName += ".txt"
  MobStatsFile, err := os.Create(MobStatsFileName)
  if err != nil {
    // Open file failed
    LogBuf = "Mobile::CreateMobStatsFileWrite - Open for " + Directory + " " + MobileIdForFight + " failed."
    LogIt(LogBuf)
    os.Exit(1)
  }
  MobStatsFile.WriteString(Stuff + "\n")
  MobStatsFile.Close()
}

// Create a player-mob relationship file
func CreatePlayerMob(PlayerName, MobileId string) {
  var PlayerMobFile     *os.File
  var PlayerMobFileName  string

  PlayerMobFileName  = PLAYER_MOB_DIR
  PlayerMobFileName += PlayerName
  PlayerMobFileName += ".txt"
  PlayerMobFile, err := os.Create(PlayerMobFileName)
  if err != nil {
    LogBuf = "Mobile::CreatePlayerMob - Open PlayerMob file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  PlayerMobFile.WriteString(MobileId + "\n")
  PlayerMobFile.Close()
}

// Delete a player-mob relationship file
func DeleteMobPlayer(PlayerName, MobileId string) {
  var BytesInFile           int64
  var MobileIdDeleted       bool
  var MobileIdCheck         string
  var MobPlayerFileName     string
  var MobPlayerFileNameTmp  string
  var MobPlayerFile        *os.File
  var MobPlayerFileTmp     *os.File

  MobileId = StrMakeLower(MobileId)
  // Open MobPlayer file
  MobPlayerFileName = MOB_PLAYER_DIR
  MobPlayerFileName += PlayerName
  MobPlayerFileName += ".txt"
  MobPlayerFile, err := os.Open(MobPlayerFileName)
  if err != nil {
    // MobPlayer player file does not exist
    return
  }
  if MobileId == "file" {
    // Delete the file
    MobPlayerFile.Close()
    Remove(MobPlayerFileName)
    return
  }
  // Open temp MobPlayer file
  MobPlayerFileNameTmp = MOB_PLAYER_DIR
  MobPlayerFileNameTmp += PlayerName
  MobPlayerFileNameTmp += ".tmp.txt"
  MobPlayerFileTmp, err = os.Create(MobPlayerFileNameTmp)
  if err != nil {
    LogBuf = "Mobile::DeleteMobPlayer - Open MobPlayer temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Write temp MobPlayer file
  MobileIdDeleted = false
  Scanner := bufio.NewScanner(MobPlayerFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    if MobileIdDeleted {
      // Mobile has been deleted, just write the rest of the mobiles
      Stuff += "\n"
      MobPlayerFileTmp.WriteString(Stuff)
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    MobileIdCheck = StrGetWord(Stuff, 1)
    MobileIdCheck = StrMakeLower(MobileIdCheck)
    if MobileId == MobileIdCheck {
      // Found it, delete it
      MobileIdDeleted = true
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    // None of the above conditions satisfied, just write it
    Stuff += "\n"
    MobPlayerFileTmp.WriteString(Stuff)
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  FileInfo, err := MobPlayerFileTmp.Stat()
  if err == nil {
    BytesInFile = FileInfo.Size()
  }
  MobPlayerFile.Close()
  MobPlayerFileTmp.Close()
  Remove(MobPlayerFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    Rename(MobPlayerFileNameTmp, MobPlayerFileName)
  } else {
    // If the file is empty, delete it
    Remove(MobPlayerFileNameTmp)
  }
}

// Delete mobile statistics files
func DeleteMobStats(MobileId string) {
  var MobStatsFileName  string
  var PlayerMobFileName string

  _ = PlayerMobFileName
  // Delete 'MobStats' Armor file
  MobStatsFileName  = MOB_STATS_ARM_DIR
  MobStatsFileName += MobileId
  MobStatsFileName += ".txt"
  Remove(MobStatsFileName)
  // Delete 'MobStats' Attack file
  MobStatsFileName  = MOB_STATS_ATK_DIR
  MobStatsFileName += MobileId
  MobStatsFileName += ".txt"
  Remove(MobStatsFileName)
  // Delete 'MobStats' Damage file
  MobStatsFileName  = MOB_STATS_DMG_DIR
  MobStatsFileName += MobileId
  MobStatsFileName += ".txt"
  Remove(MobStatsFileName)
  // Delete 'MobStats' Desc1 file
  MobStatsFileName  = MOB_STATS_DSC_DIR
  MobStatsFileName += MobileId
  MobStatsFileName += ".txt"
  Remove(MobStatsFileName)
  // Delete 'MobStats' ExpPoints file
  MobStatsFileName  = MOB_STATS_EXP_DIR
  MobStatsFileName += MobileId
  MobStatsFileName += ".txt"
  Remove(MobStatsFileName)
  // Delete 'MobStats' HitPoints file
  MobStatsFileName = MOB_STATS_HPT_DIR
  MobStatsFileName += MobileId
  MobStatsFileName += ".txt"
  Remove(MobStatsFileName)
  // Delete 'MobStats' Loot file
  MobStatsFileName  = MOB_STATS_LOOT_DIR
  MobStatsFileName += MobileId
  MobStatsFileName += ".txt"
  Remove(MobStatsFileName)
  // Delete 'MobStats' Room file
  MobStatsFileName  = MOB_STATS_ROOM_DIR
  MobStatsFileName += MobileId
  MobStatsFileName += ".txt"
  Remove(MobStatsFileName)
}

// Delete player-mob relationship file
func DeletePlayerMob(PlayerName string) {
  var MobStatsFileName  string
  var PlayerMobFileName string

  _ = MobStatsFileName
  // Delete 'PlayerMob' file
  PlayerMobFileName  = PLAYER_MOB_DIR
  PlayerMobFileName += PlayerName
  PlayerMobFileName += ".txt"
  Remove(PlayerMobFileName)
}

// Check if a mobile is in the room by its name
func IsMobInRoom(MobileName string) *Mobile {
  var pMobile         *Mobile
  var NamesCheck       string
  var MobileHurt       bool
  var MobileId         string
  var MobileIdCheck    string
  var MobileIdHurt     string
  var MobileNameCheck  string
  var MobNbr           string
  var PositionOfDot    int
  var RoomMobFile     *os.File
  var RoomMobFileName  string

  _ = MobileIdCheck
  _ = MobileIdHurt
  _ = MobileNameCheck

  // Open RoomMob file
  RoomMobFileName  = ROOM_MOB_DIR
  RoomMobFileName += pDnodeActor.pPlayer.RoomId
  RoomMobFileName += ".txt"
  //*******************************
  //* Try matching using MobileId *
  //*******************************
  // Try matching using MobileId
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    // Room has no mobiles
    return nil
  }
  Scanner := bufio.NewScanner(RoomMobFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // Process each mobile in the room
    MobileId = StrGetWord(Stuff, 2)
    if MobileId == MobileName {
      // This mobile is a match
      RoomMobFile.Close()
      PositionOfDot = StrFindFirstChar(MobileId, '.')
      MobileHurt = false
      if PositionOfDot > 1 {
        // Mobile is hurt but not fighting
        MobileHurt = true
        MobileIdHurt = MobileId
        MobNbr = StrRight(MobileId, StrGetLength(MobileId)-PositionOfDot-1)
        MobileId = StrLeft(MobileId, PositionOfDot)
      }
      pMobile        = MobileConstructor(MobileId)
      pMobile.Hurt   = MobileHurt
      pMobile.MobNbr = MobNbr
      return pMobile
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  RoomMobFile.Close()
  //***************************************************
  //* No match found, try getting match using 'names' *
  //***************************************************
  // No match found, try getting match using 'names'
  RoomMobFile, err = os.Open(RoomMobFileName)
  if err != nil {
    // Room has no mobiles
    return nil
  }
  Scanner = bufio.NewScanner(RoomMobFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // Process each mobile in the room
    MobileId      = StrGetWord(Stuff, 2)
    PositionOfDot = StrFindFirstChar(MobileId, '.')
    MobileHurt     = false
    if PositionOfDot > 1 {
      // Mobile is hurt but not fighting
      MobileHurt   = true
      MobileIdHurt = MobileId
      MobNbr       = StrRight(MobileId, StrGetLength(MobileId)-PositionOfDot-1)
      MobileId     = StrLeft(MobileId, PositionOfDot)
    }
    pMobile        = MobileConstructor(MobileId)
    pMobile.Hurt   = MobileHurt
    pMobile.MobNbr = MobNbr
    if pMobile.Hurt {
      // Mobile is hurt
      if MobNbr == MobileName {
        // Kill nnn was entered, where nnn is the MobNbr
        RoomMobFile.Close()
        return pMobile
      }
    }
    NamesCheck = pMobile.Names
    NamesCheck = StrMakeLower(NamesCheck)
    if StrIsWord(MobileName, NamesCheck) {
      // This mobile is a match
      RoomMobFile.Close()
      return pMobile
    } else {
      // This mobile doesn't match
      pMobile = nil
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  RoomMobFile.Close()
  return nil
}

// Get the first description of a mobile
func GetMobDesc1(MobileId string) string {
  var Desc1           string
  var MobileFile     *os.File
  var MobileFileName  string
  var PositionOfDot   int

  PositionOfDot = StrFindFirstChar(MobileId, '.')
  if PositionOfDot > 1 {
    // Mobile is hurt but not fighting
    MobileId = StrLeft(MobileId, PositionOfDot)
  }
  MobileFileName  = MOBILES_DIR
  MobileFileName += MobileId
  MobileFileName += ".txt"
  MobileFile, err := os.Open(MobileFileName)
  if err != nil {
    LogBuf = "Mobile::GetMobDesc1 - Mobile does not exist!"
    LogIt(LogBuf)
    os.Exit(1)
  }
  Stuff = ""
  Scanner := bufio.NewScanner(MobileFile)
  for StrLeft(Stuff, 6) != "Desc1:" {
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  Desc1 = StrRight(Stuff, StrGetLength(Stuff)-6)
  Desc1 = StrTrimLeft(Desc1)
  MobileFile.Close()
  return Desc1
}

// Check if a mobile ID is in a room
func IsMobileIdInRoom(RoomId, MobileId string) bool {
  var MobileIdCheck    string
  var RoomMobFile     *os.File
  var RoomMobFileName  string

  // Open RoomMob file
  RoomMobFileName = ROOM_MOB_DIR
  RoomMobFileName += RoomId
  RoomMobFileName += ".txt"
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    // Room has no mobiles
    return false
  }
  Scanner := bufio.NewScanner(RoomMobFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // Process each mobile in the room
    MobileIdCheck = StrGetWord(Stuff, 2)
    if MobileId == MobileIdCheck {
      // Found matching mobile
      RoomMobFile.Close()
      return true
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  // No matching mobile found
  RoomMobFile.Close()
  return false
}

// Check if a mobile is valid by its ID
func IsMobValid(MobileId string) *Mobile {
  var pMobile        *Mobile
  var MobileFileName  string
  var MobileFile     *os.File

  _ = MobileFile

  MobileFileName = MOBILES_DIR
  MobileFileName += MobileId
  MobileFileName += ".txt"
  if FileExist(MobileFileName) {
    pMobile = MobileConstructor(MobileId)
    return pMobile
  } else {
    return nil
  }
}

// Put a mobile back in the room
func PutMobBackInRoom(PlayerName, RoomIdBeforeFleeing string) {
  var MobHitPointsLeft           string
  var MobHitPointsTotal          string
  var MobileId                   string
  var MobPlayerFile             *os.File
  var MobPlayerFileName          string
  var MobStatsHitPointsFile     *os.File
  var MobStatsHitPointsFileName  string
  var PositionOfDot              int

  // Open MobPlayer file
  MobPlayerFileName  = MOB_PLAYER_DIR
  MobPlayerFileName += PlayerName
  MobPlayerFileName += ".txt"
  MobPlayerFile, err := os.Open(MobPlayerFileName)
  if err != nil {
    // No mobiles to put back, someone else may be fighting the mob
    return
  }
  // For each mobile still in MobPlayer file (non-fighting mobiles), put it back in room
  Scanner := bufio.NewScanner(MobPlayerFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    MobileId = StrMakeLower(StrGetWord(Stuff, 1))
    // Read mobile stats hit points file
    MobStatsHitPointsFileName = MOB_STATS_HPT_DIR
    MobStatsHitPointsFileName += MobileId
    MobStatsHitPointsFileName += ".txt"
    MobStatsHitPointsFile, err = os.Open(MobStatsHitPointsFileName)
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
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  MobPlayerFile.Close()
  Remove(MobPlayerFileName)
}

// Remove a mobile from room
func RemoveMobFromRoom(RoomId, MobileId string) {
  var BytesInFile         int64
  var MobileIdRemoved     bool
  var MobileIdCheck       string
  var MobCount            int
  var RoomMobFile        *os.File
  var RoomMobFileName     string
  var RoomMobTmpFile     *os.File
  var RoomMobTmpFileName  string

  UpdateMobInWorld(MobileId, "remove")
  MobileId = StrMakeLower(MobileId)
  // Open RoomMob file
  RoomMobFileName  = ROOM_MOB_DIR
  RoomMobFileName += RoomId
  RoomMobFileName += ".txt"
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    LogBuf = "Mobile::RemoveMobFromRoom - Open RoomMob file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Open temp RoomMob file
  RoomMobTmpFileName  = ROOM_MOB_DIR
  RoomMobTmpFileName += RoomId
  if RoomId == "" {
    LogBuf = "RoomId is blank 2"
    LogIt(LogBuf)
    os.Exit(1)
  }
  RoomMobTmpFileName += ".tmp.txt"
  RoomMobTmpFile, err = os.Create(RoomMobTmpFileName)
  if err != nil {
    LogBuf = "Mobile::RemoveMobFromRoom - Open RoomMob temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Write temp RoomMob file
  MobileIdRemoved = false
  Scanner := bufio.NewScanner(RoomMobFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    if MobileIdRemoved {
      // Mobile has been removed, just write the rest of the mobiles
      Stuff += "\n"
      RoomMobTmpFile.WriteString(Stuff)
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    MobileIdCheck = StrGetWord(Stuff, 2)
    if MobileId == MobileIdCheck {
      // Found it, subtract 1 from count
      MobCount = StrToInt(StrGetWord(Stuff, 1))
      MobCount--
      MobileIdRemoved = true
      if MobCount > 0 {
        Buf = fmt.Sprintf("%d", MobCount)
        TmpStr = Buf
        MobileId = TmpStr + " " + MobileId
        MobileId += "\n"
        RoomMobTmpFile.WriteString(MobileId)
      }
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    // None of the above conditions satisfied, just write it
    Stuff += "\n"
    RoomMobTmpFile.WriteString(Stuff)
    Scanner.Scan()
    Stuff = Scanner.Text()
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
  Remove(RoomMobFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    Rename(RoomMobTmpFileName, RoomMobFileName)
  } else {
    // If the file is empty, delete it
    Remove(RoomMobTmpFileName)
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
  var RoomMobFile          *os.File
  var RoomMobFileName       string

  // Open RoomMob file
  RoomMobFileName  = ROOM_MOB_DIR
  RoomMobFileName += pDnode.pPlayer.RoomId
  RoomMobFileName += ".txt"
  RoomMobFile, err := os.Open(RoomMobFileName)
  if err != nil {
    // No mobiles in room to display
    return
  }
  Scanner := bufio.NewScanner(RoomMobFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    MobileCount   = StrGetWord(Stuff, 1)
    MobileId      = StrGetWord(Stuff, 2)
    PositionOfDot = StrFindFirstChar(MobileId, '.')
    MobileHurt    = false
    MobNbr = ""
    if PositionOfDot > 1 {
      // Mobile is hurt but not fighting
      MobileHurt   = true
      MobileIdHurt = MobileId
      MobNbr       = StrRight(MobileId, StrGetLength(MobileId)-PositionOfDot-1)
      MobileId     = StrLeft(MobileId, PositionOfDot)
    }
    pMobile = MobileConstructor(MobileId)
    pMobile.Hurt   = MobileHurt
    pMobile.MobNbr = MobNbr
    if MobileHurt {
      // Mobile is hurt
      pDnode.PlayerOut += "\r\n"
      pDnode.PlayerOut += "&W"
      pDnode.PlayerOut += "You see "
      pDnode.PlayerOut += pMobile.Desc1
      pDnode.PlayerOut += ", "
      pDnode.PlayerOut += "&M"
      pDnode.PlayerOut += "wounded"
      pDnode.PlayerOut += "&W"
      pDnode.PlayerOut += ", trying to hide."
      pDnode.PlayerOut += "&B"
      pDnode.PlayerOut += " ("
      pDnode.PlayerOut += MobileIdHurt
      pDnode.PlayerOut += ")"
      pDnode.PlayerOut += "&N"
    } else {
      // Mobile is not hurt
      pDnode.PlayerOut += "\r\n"
      pDnode.PlayerOut += "&W"
      pDnode.PlayerOut += "(" + MobileCount + ") "
      pDnode.PlayerOut += pMobile.Desc2
      pDnode.PlayerOut += "&N"
    }
    // Check for AGGRO mobs
    if StrIsWord("Aggro", pMobile.Action) {
      // Attack player
      j = StrToInt(MobileCount)
      for i = 1; i <= j; i++ {
        MobileIdsToBeRemoved += MobAttacks(pMobile)
        MobileIdsToBeRemoved += " "
      }
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  RoomMobFile.Close()
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
  RoomId     = pDnodeActor.pPlayer.RoomId
  //*****************
  //* Send messages *
  //*****************
  // Determine appropriate phrase
  if !pDnodeActor.PlayerStateFighting {
    // Phrases for starting a fight
    PhrasePlayer = " starts a fight with you!"
    PhraseAll    = " starts a fight with "
  } else {
    // Phrases for mob attacking a player already fighting
    PhrasePlayer = " attacks you!"
    PhraseAll    = " attacks "
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
  //*****************
  //* Start a fight *
  //*****************
  // Start a fight
  if !pMobile.Hurt {
    //  Mobile not hurt
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
func WhereMob(MobileIdSearch string) {
  var FileName         string
  var MobileHurt       bool
  var MobileId         string
  var PositionOfDot    int
  var RoomMobFile     *os.File
  var RoomMobFileName  string
  var RoomName         string

  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Mobiles"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "-------"
  pDnodeActor.PlayerOut += "\r\n"
  if ChgDir(ROOM_MOB_DIR) != nil {
    LogBuf = "Mobile::WhereMob - Change directory to ROOM_MOB_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogBuf = "Mobile::WhereMob - Failed to read ROOM_MOB_DIR"
    LogIt(LogBuf)
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    FileName = DirEntry.Name()
    // Open RoomMob file
    RoomMobFileName = FileName
    RoomMobFile, err = os.Open(RoomMobFileName)
    if err != nil {
      LogBuf = "Mobile::WhereMob - Open RoomMob file failed"
      LogIt(LogBuf)
      os.Exit(1)
    }
    RoomName = StrLeft(FileName, StrGetLength(FileName)-4)
    Scanner := bufio.NewScanner(RoomMobFile)
    Scanner.Scan()
    Stuff = Scanner.Text()
    for Stuff != "" {
      MobileId = StrGetWord(Stuff, 2)
      PositionOfDot = StrFindFirstChar(MobileId, '.')
      MobileHurt = false
      if PositionOfDot > 1 {
        // Mobile is hurt but not fighting
        MobileHurt = true
        MobileId = StrLeft(MobileId, PositionOfDot)
      }
      if MobileId == MobileIdSearch {
        pDnodeActor.PlayerOut += RoomName
        pDnodeActor.PlayerOut += " "
        if MobileHurt {
          pDnodeActor.PlayerOut += "&R"
        }
        pDnodeActor.PlayerOut += Stuff
        pDnodeActor.PlayerOut += "&N"
        pDnodeActor.PlayerOut += "\r\n"
      }
      Scanner.Scan()
      Stuff = Scanner.Text()
    }
    RoomMobFile.Close()
  }
  if ChgDir(HomeDir) != nil {
    LogBuf = "Mobile::WhereMob - Change directory to HomeDir failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
}

// Update the count of a mobile in the world
func UpdateMobInWorld(MobileId string, AddRemove string) {
  var MobInWorldCount     int
  var MobInWorldFile     *os.File
  var MobInWorldFileName  string
  var PositionOfDot       int

  MobInWorldCount = 0
  PositionOfDot = StrFindFirstChar(MobileId, '.')
  if PositionOfDot > 1 {
    // Get MobileId
    MobileId = StrLeft(MobileId, PositionOfDot)
  }
  // Open Mobile InWorld file
  MobInWorldFileName = CONTROL_MOB_INWORLD_DIR
  MobInWorldFileName += MobileId
  MobInWorldFileName += ".txt"
  MobInWorldFile, err := os.Open(MobInWorldFileName)
  if err == nil {
    // Get current count
    Scanner := bufio.NewScanner(MobInWorldFile)
    Scanner.Scan()
    Stuff = Scanner.Text()
    MobInWorldCount = StrToInt(Stuff)
    MobInWorldFile.Close()
  }
  // Create Mobiles InWorld file, doesn't matter if it already exists
  MobInWorldFile, err = os.Create(MobInWorldFileName)
  if err != nil {
    LogBuf = "Mobile::UpdateMobInWorld - Open Mobiles InWorld file failed for: " + MobInWorldFileName
    LogIt(LogBuf)
    return
  }
  if AddRemove == "add" {
    // Mobile is being added to the world
    MobInWorldCount++
  } else {
    // Mobile is being removed from the world
    MobInWorldCount--
  }
  Buf = fmt.Sprintf("%d", MobInWorldCount)
  TmpStr = Buf
  MobInWorldFile.WriteString(TmpStr + "\n")
  MobInWorldFile.Close()
}

// Create mobile statistics file
func CreateMobStatsFile(RoomId string) {
  var MobileIdForFight string

  MobileIdForFight = MobileId + "." + MobNbr
  // HitPoints
  Buf = fmt.Sprintf("%d", HitPoints)
  TmpStr = Buf
  Stuff = TmpStr
  Stuff += " "
  Stuff += TmpStr
  CreateMobStatsFileWrite(MOB_STATS_HPT_DIR, MobileIdForFight, Stuff)
  // Armor
  Buf = fmt.Sprintf("%d", Armor)
  Stuff = Buf
  CreateMobStatsFileWrite(MOB_STATS_ARM_DIR, MobileIdForFight, Stuff)
  // Attack
  Stuff= Attack
  CreateMobStatsFileWrite(MOB_STATS_ATK_DIR, MobileIdForFight, Stuff)
  // Damage
  Buf = fmt.Sprintf("%d", Damage)
  Stuff = Buf
  CreateMobStatsFileWrite(MOB_STATS_DMG_DIR, MobileIdForFight, Stuff)
  // Desc1
  Stuff = Desc1
  CreateMobStatsFileWrite(MOB_STATS_DSC_DIR, MobileIdForFight, Stuff)
  // ExpPoints
  Buf = fmt.Sprintf("%d", ExpPoints)
  Stuff = Buf
  Buf = fmt.Sprintf("%d", Level)
  TmpStr = Buf
  Stuff += " "
  Stuff += TmpStr
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
  OpenMobFile(MobileId)
  for Stuff != "Desc3:" {
    ReadMobLine() // Do not use ReadLine() here
  }
  // Mobile Description 3
  ReadMobLine() // Do not use ReadLine() here
  for Stuff != "End Desc3" {
    pDnodeActor.PlayerOut += Stuff
    pDnodeActor.PlayerOut += "\r\n"
    ReadMobLine() // Do not use ReadLine() here
  }
  pDnodeActor.PlayerOut += "&N"
  CloseMobFile()
}

// Get the next mobile number
func GetNextMobNbr() {
  var NextMobNbr          string
  var NextMobNbrFile     *os.File
  var NextMobNbrFileName  string
  var NextMobNbrInteger   int
  var NextMobNbrString    string

  // Read next mobile number file
  NextMobNbrFileName  = CONTROL_DIR
  NextMobNbrFileName += "NextMobileNumber"
  NextMobNbrFileName += ".txt"
  NextMobNbrFile, err := os.Open(NextMobNbrFileName)
  if err != nil {
    LogBuf = "Mobile::GetNextMobNbr - Open NextMobileNumber file failed (read)"
    LogIt(LogBuf)
    os.Exit(1)
  }
  Scanner := bufio.NewScanner(NextMobNbrFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  NextMobNbrFile.Close()
  // Increment next mobile number
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  NextMobNbr = Stuff
  NextMobNbrInteger = StrToInt(Stuff)
  NextMobNbrInteger++
  Buf = fmt.Sprintf("%d", NextMobNbrInteger)
  NextMobNbrString = Buf
  // Write next mobile number file
  NextMobNbrFile, err = os.Create(NextMobNbrFileName)
  if err != nil {
    LogBuf = "Mobile::GetNextMobNbr - Open NextMobileNumber file failed (write)"
    LogIt(LogBuf)
    os.Exit(1)
  }
  NextMobNbrFile.WriteString(NextMobNbrString + "\n")
  NextMobNbrFile.Close()
  // Set mobile number
  MobNbr = NextMobNbr
}

// Generate a message for a mobile to say
func MobTalk(pMobile *Mobile) string {
  var MobTalkFile     *os.File
  var MobTalkFileName  string
  var MobileMsg        string
  var MsgCount         int
  var RndMsgNbr        int

  _ = pMobile

  //******************************
  //* Open and read message file *
  //******************************
  // Open and read message file
  MobTalkFileName  = TALK_DIR
  MobTalkFileName += Talk
  MobTalkFileName += ".txt"
  MobTalkFile, err := os.Open(MobTalkFileName)
  if err != nil {
    // Open failed
    if Talk != "None" {
      // Talk is not 'None', so file should exist
      LogBuf = "Mobile::MobTalk - Failed to open "
      LogBuf += MobTalkFileName
      LogIt(LogBuf)
    }
    MobileMsg = "You are ignored.\r\n"
    return MobileMsg
  }
  // Mobile is going to talk
  MobileMsg = "&W"
  MobileMsg += StrMakeFirstUpper(Desc1)
  MobileMsg += " says:"
  MobileMsg += "&N"
  MobileMsg += "\r\n"
  // Select random message number
  Scanner := bufio.NewScanner(MobTalkFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  MsgCount = StrToInt(StrGetWord(Stuff, 4))
  RndMsgNbr = GetRandomNumber(MsgCount)
  // Search for selected message number
  Scanner.Scan()
  Stuff = Scanner.Text()
  for StrToInt(StrGetWord(Stuff, 2)) != RndMsgNbr {
    // Find the selected message
    if !Scanner.Scan() {
      // End of file and message was not found
      Buf = fmt.Sprintf("%d", RndMsgNbr)
      TmpStr = Buf
      LogBuf = "Mobile::MobTalk - Failed to find message "
      LogBuf += TmpStr
      LogBuf += " "
      LogBuf += MobTalkFileName
      LogIt(LogBuf)
      MobTalkFile.Close()
      MobileMsg = "You are ignored.\r\n"
      return MobileMsg
    }
    Stuff = Scanner.Text()
  }
  // Message found
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "End of Message" {
    // Read the message
    if !Scanner.Scan() {
      // End of file and normal end of message not found
      Buf = fmt.Sprintf("%d", RndMsgNbr)
      TmpStr = Buf
      LogBuf = "Mobile::MobTalk - Unexpect end of file reading message "
      LogBuf += TmpStr
      LogBuf += " "
      LogBuf += MobTalkFileName
      LogIt(LogBuf)
      MobTalkFile.Close()
      MobileMsg = "You are ignored.\r\n"
      return MobileMsg
    }
    MobileMsg += Stuff
    MobileMsg += "\r\n"
    Stuff = Scanner.Text()
  }
  MobTalkFile.Close()
  return MobileMsg
}

// Close the mobile file
func CloseMobFile() {
  MobileFile.Close()
  MobileFile = nil
}

// Open the file for a given mobile ID
func OpenMobFile(MobileId string) {
  var MobileFileName string

  MobileFileName = MOBILES_DIR
  MobileFileName += MobileId
  MobileFileName += ".txt"
  MobileFile, err := os.Open(MobileFileName)
  if err != nil {
    LogBuf = "Mobile::OpenFile - Mobile does not exist!"
    LogIt(LogBuf)
    os.Exit(1)
  }
  MobScanner = bufio.NewScanner(MobileFile)
}

// Parse mobile stuff
func ParseMobStuff() {
  ReadMobLine()
  for Stuff != "" {
    if StrLeft(Stuff, 9) == "MobileId:" {
      MobileId = StrRight(Stuff, StrGetLength(Stuff)-9)
      MobileId = StrTrimLeft(MobileId)
    } else
    if StrLeft(Stuff, 6) == "Names:" {
      Names = StrRight(Stuff, StrGetLength(Stuff)-6)
      Names = StrTrimLeft(Names)
    } else
    if StrLeft(Stuff, 6) == "Desc1:" {
      Desc1 = StrRight(Stuff, StrGetLength(Stuff)-6)
      Desc1 = StrTrimLeft(Desc1)
    } else
    if StrLeft(Stuff, 6) == "Desc2:" {
      Desc2 = StrRight(Stuff, StrGetLength(Stuff)-6)
      Desc2 = StrTrimLeft(Desc2)
    } else
    if StrLeft(Stuff, 6) == "Desc3:" {
      // Desc3 can be multi-line and is dealt with in 'ExamineMob'
    } else
    if StrLeft(Stuff, 7) == "Action:" {
      Action = StrRight(Stuff, StrGetLength(Stuff)-7)
      Action = StrTrimLeft(Action)
    } else
    if StrLeft(Stuff, 8) == "Faction:" {
      Faction = StrRight(Stuff, StrGetLength(Stuff)-8)
      Faction = StrTrimLeft(Faction)
    } else
    if StrLeft(Stuff, 6) == "Level:" {
      Level = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-6))
    } else
    if StrLeft(Stuff, 10) == "HitPoints:" {
      HitPoints = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-10))
      HitPoints += Level * MOB_HPT_PER_LEVEL
    } else
    if StrLeft(Stuff, 6) == "Armor:" {
      Armor = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-6))
      Armor += Level * MOB_ARM_PER_LEVEL
    } else
    if StrLeft(Stuff, 7) == "Attack:" {
      Attack = StrRight(Stuff, StrGetLength(Stuff)-7)
      Attack = StrTrimLeft(Attack)
      Attack = StrMakeLower(Attack)
    } else
    if StrLeft(Stuff, 7) == "Damage:" {
      Damage = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-7))
      Damage += Level * MOB_DMG_PER_LEVEL
    } else
    if StrLeft(Stuff, 10) == "ExpPoints:" {
      ExpPoints = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-10))
      ExpPoints += Level * MOB_EXP_PER_LEVEL
    } else
    if StrLeft(Stuff, 5) == "Loot:" {
      Loot = StrRight(Stuff, StrGetLength(Stuff)-5)
      Loot = StrTrimLeft(Loot)
    } else
    if StrLeft(Stuff, 5) == "Talk:" {
      Talk = StrRight(Stuff, StrGetLength(Stuff)-5)
      Talk = StrTrimLeft(Talk)
    }
    ReadMobLine()
  }
}

// Read a line from Mobile file
func ReadMobLine() {
  Stuff = ""
  if MobScanner != nil && MobScanner.Scan() {
    Stuff = MobScanner.Text()
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
  }
}
