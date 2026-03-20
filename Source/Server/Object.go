//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Object.go                                             *
// Usage:     Manages game objects and their interactions            *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
	"bufio"
	"os"
	"strconv"
)

var ObjectFile     *os.File
var ObjectFileName  string
var ObjectId        string
var pObject        *Object

// Object represents a single instance of an object in the game.
type Object struct {
  ObjectId          string
  ArmorValue        int
  ArmorWear         string
  ContainerCapacity int
  Cost              int
  Count             string
  Desc1             string
  Desc2             string
  Desc3             string
  DrinkPct          int
  FoodPct           int
  LightHours        int
  Names             string
  Type              string
  WeaponType        string
  WeaponDamage      int
  WearPosition      string
  Weight            int
}

var ObjScanner *bufio.Scanner

// Create a new object
func ObjectConstructor(ObjectIdParm string) {
  DEBUGIT(5)
  // Init variables
  pObject = &Object{
    ObjectId:          ObjectIdParm,
    ArmorValue:        0,
    ArmorWear:         "",
    ContainerCapacity: 0,
    Cost:              0,
    Count:             "1",
    Desc1:             "",
    Desc2:             "",
    Desc3:             "",
    DrinkPct:          0,
    FoodPct:           0,
    LightHours:        0,
    Names:             "",
    Type:              "",
    WeaponType:        "",
    WeaponDamage:      0,
    WearPosition:      "",
    Weight:            0,
  }
  // Construct object
  OpenObjectFile(ObjectIdParm)
  ParseObjectStuff()
  CloseObjectFile()
}

// Add an object to player's equipment
func AddObjToPlayerEqu(WearPosition string, ObjectId string) bool {
  var NewPlayerEquFile      bool
  var ObjectIdAdded         bool
  var PlayerEquFileName     string
  var PlayerEquFileNameTmp  string
  var PlayerEquFile        *os.File
  var PlayerEquFileTmp     *os.File
  var WearPositionCheck     string
  var WearWieldFailed       bool

  DEBUGIT(5)
  WearWieldFailed = false
  ObjectId = StrMakeLower(ObjectId)
  // Open PlayerEqu file
  PlayerEquFileName = PLAYER_EQU_DIR
  PlayerEquFileName += pDnodeActor.PlayerName
  PlayerEquFileName += ".txt"
  NewPlayerEquFile = false
  PlayerEquFile, err := os.Open(PlayerEquFileName)
  if err != nil {
    NewPlayerEquFile = true
  }
  // Open temp PlayerEqu file
  PlayerEquFileNameTmp = PLAYER_EQU_DIR
  PlayerEquFileNameTmp += pDnodeActor.PlayerName
  PlayerEquFileNameTmp += ".tmp.txt"
  PlayerEquFileTmp, err = os.Create(PlayerEquFileNameTmp)
  if err != nil {
    LogBuf = "Object::AddObjToPlayerEqu - Open PlayerEqu temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  WearPosition = TranslateWord(WearPosition)
  if NewPlayerEquFile {
    // New player equipment file, write the object and return
    ObjectId = WearPosition + " " + ObjectId
    PlayerEquFileTmp.WriteString(ObjectId + "\n")
    PlayerEquFileTmp.Close()
    Rename(PlayerEquFileNameTmp, PlayerEquFileName)
    return WearWieldFailed
  }
  // Write temp PlayerEqu file
  ObjectIdAdded = false
  Scanner := bufio.NewScanner(PlayerEquFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    if ObjectIdAdded {
      // New object has been written, just write the rest of the objects
      PlayerEquFileTmp.WriteString(Stuff)
      PlayerEquFileTmp.WriteString("\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    WearPositionCheck = StrGetWord(Stuff, 1)
    if WearPosition < WearPositionCheck {
      // Add new object in alphabetical order by translated WearPosition
      ObjectId = WearPosition + " " + ObjectId
      PlayerEquFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      PlayerEquFileTmp.WriteString(Stuff + "\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    if WearPosition == WearPositionCheck {
      // Already wearing an object in that position
      WearWieldFailed = true
      ObjectIdAdded = true // Not really added
      PlayerEquFileTmp.WriteString(Stuff + "\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    // None of the above conditions satisfied, just write it
    PlayerEquFileTmp.WriteString(Stuff + "\n")
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  if !ObjectIdAdded {
    // New object is alphabetically last
    ObjectId = WearPosition + " " + ObjectId
    PlayerEquFileTmp.WriteString(ObjectId + "\n")
    ObjectIdAdded = true
  }
  PlayerEquFile.Close()
  PlayerEquFileTmp.Close()
  Remove(PlayerEquFileName)
  Rename(PlayerEquFileNameTmp, PlayerEquFileName)
  return WearWieldFailed
}

// Add an object to player's inventory
func AddObjToPlayerInv(pDnodeTgt1 *Dnode, ObjectId string) {
  var NewPlayerObjFile      bool
  var ObjectIdAdded         bool
  var ObjectIdCheck         string
  var ObjCount              int
  var PlayerObjFileName     string
  var PlayerObjFileNameTmp  string
  var PlayerObjFile        *os.File
  var PlayerObjFileTmp     *os.File

  DEBUGIT(5)
  pDnodeTgt = pDnodeTgt1
  ObjectId = StrMakeLower(ObjectId)
  // Open PlayerObj file
  PlayerObjFileName = PLAYER_OBJ_DIR
  PlayerObjFileName += pDnodeTgt.PlayerName
  PlayerObjFileName += ".txt"
  NewPlayerObjFile = false
  PlayerObjFile, err := os.Open(PlayerObjFileName)
  if err != nil {
    NewPlayerObjFile = true
  }
  // Open temp PlayerObj file
  PlayerObjFileNameTmp = PLAYER_OBJ_DIR
  PlayerObjFileNameTmp += pDnodeTgt.PlayerName
  PlayerObjFileNameTmp += ".tmp.txt"
  PlayerObjFileTmp, err = os.Create(PlayerObjFileNameTmp)
  if err != nil {
    LogBuf = "Object::AddObjToPlayerInv - Open PlayerObj temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  if NewPlayerObjFile {
    // New player inventory file, write the object and return
    ObjectId = "1 " + ObjectId
    PlayerObjFileTmp.WriteString(ObjectId + "\n")
    PlayerObjFileTmp.Close()
    Rename(PlayerObjFileNameTmp, PlayerObjFileName)
    return
  }
  // Write temp PlayerObj file
  ObjectIdAdded = false
  Scanner := bufio.NewScanner(PlayerObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    if ObjectIdAdded {
      // New object has been written, just write the rest of the objects
      PlayerObjFileTmp.WriteString(Stuff + "\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId < ObjectIdCheck {
      // Add new object in alphabetical order
      ObjectId = "1 " + ObjectId
      PlayerObjFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      PlayerObjFileTmp.WriteString(Stuff + "\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    if ObjectId == ObjectIdCheck {
      // Existing object same as new object, add 1 to count
      ObjCount = StrToInt(StrGetWord(Stuff, 1))
      ObjCount++
      Buf = strconv.Itoa(ObjCount)
      TmpStr = Buf
      ObjectId = TmpStr + " " + ObjectId
      PlayerObjFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    // None of the above conditions satisfied, just write it
    PlayerObjFileTmp.WriteString(Stuff + "\n")
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  if !ObjectIdAdded {
    // New object is alphabetically last
    ObjectId = "1 " + ObjectId
    PlayerObjFileTmp.WriteString(ObjectId + "\n")
    ObjectIdAdded = true
  }
  PlayerObjFile.Close()
  PlayerObjFileTmp.Close()
  Remove(PlayerObjFileName)
  Rename(PlayerObjFileNameTmp, PlayerObjFileName)
}

// Add an object to room
func AddObjToRoom(RoomId string, ObjectId string) {
  var NewRoomObjFile      bool
  var ObjectIdAdded       bool
  var ObjectIdCheck       string
  var ObjCount            int
  var RoomObjFileName     string
  var RoomObjFileNameTmp  string
  var RoomObjFile        *os.File
  var RoomObjFileTmp     *os.File

  DEBUGIT(5)
  ObjectId = StrMakeLower(ObjectId)
  // Open RoomObj file
  RoomObjFileName = ROOM_OBJ_DIR
  RoomObjFileName += RoomId
  RoomObjFileName += ".txt"
  NewRoomObjFile = false
  RoomObjFile, err := os.Open(RoomObjFileName)
  if err != nil {
    NewRoomObjFile = true
  }
  // Open temp RoomObj file
  RoomObjFileNameTmp = ROOM_OBJ_DIR
  RoomObjFileNameTmp += RoomId
  RoomObjFileNameTmp += ".tmp.txt"
  RoomObjFileTmp, err = os.Create(RoomObjFileNameTmp)
  if err != nil {
    LogBuf = "Object::AddObjToRoom - Open RoomObj temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  if NewRoomObjFile {
    // New room object file, write the object and return
    ObjectId = "1 " + ObjectId
    RoomObjFileTmp.WriteString(ObjectId + "\n")
    RoomObjFileTmp.Close()
    err = Rename(RoomObjFileNameTmp, RoomObjFileName)
    if err != nil {
      LogBuf = "Object::AddObjToRoom - Rename RoomObj temp file failed"
      LogIt(LogBuf)
      os.Exit(1)
    }
    return
  }
  // Write temp RoomObj file
  ObjectIdAdded = false
  Scanner := bufio.NewScanner(RoomObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    if ObjectIdAdded {
      // New object has been written, just write the rest of the objects
      RoomObjFileTmp.WriteString(Stuff + "\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId < ObjectIdCheck {
      // Add new object in alphabetical order
      ObjectId = "1 " + ObjectId
      RoomObjFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      RoomObjFileTmp.WriteString(Stuff + "\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    if ObjectId == ObjectIdCheck {
      // Existing object same as new object, add 1 to count
      ObjCount = StrToInt(StrGetWord(Stuff, 1))
      ObjCount++
      Buf = strconv.Itoa(ObjCount)
      TmpStr = Buf
      ObjectId = TmpStr + " " + ObjectId
      RoomObjFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    // None of the above conditions satisfied, just write it
    RoomObjFileTmp.WriteString(Stuff + "\n")
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  if !ObjectIdAdded {
    // New object is alphabetically last
    ObjectId = "1 " + ObjectId
    RoomObjFileTmp.WriteString(ObjectId + "\n")
    ObjectIdAdded = true
  }
  RoomObjFile.Close()
  RoomObjFileTmp.Close()
  Remove(RoomObjFileName)
  Rename(RoomObjFileNameTmp, RoomObjFileName)
}

// Calculate player armor class
func CalcPlayerArmorClass(pPlayer *Player) int {
  var ArmorClass         int
  var ObjectId           string
  var PlayerEquFile     *os.File
  var PlayerEquFileName  string
  var WearPosition       string
  var err                error

  _ = WearPosition

  DEBUGIT(5)
  ArmorClass = 0
  // Open PlayerEqu file
  PlayerEquFileName = PLAYER_EQU_DIR
  PlayerEquFileName += pPlayer.Name
  PlayerEquFileName += ".txt"
  PlayerEquFile, err = os.Open(PlayerEquFileName)
  if err != nil {
    // No player equipment
    return ArmorClass
  }
  Scanner := bufio.NewScanner(PlayerEquFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    ObjectId = StrGetWord(Stuff, 2)
    pObject = nil
    ObjectConstructor(ObjectId)
    ArmorClass += pObject.ArmorValue
    pObject = nil
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  PlayerEquFile.Close()
  return ArmorClass
}

// Is object in player's equipment?
func IsObjInPlayerEqu(ObjectName string) {
  var NamesCheck         string
  var ObjectId           string
  var ObjectIdCheck      string
  var ObjectNameCheck    string
  var PlayerEquFileName  string
  var PlayerEquFile     *os.File

  _ = ObjectIdCheck
  _ = ObjectNameCheck

  DEBUGIT(5)
  // Open PlayerEqu file
  PlayerEquFileName = PLAYER_EQU_DIR
  PlayerEquFileName += pDnodeActor.PlayerName
  PlayerEquFileName += ".txt"
  //*******************************
  //* Try matching using ObjectId *
  //*******************************
  PlayerEquFile, err := os.Open(PlayerEquFileName)
  if err != nil {
    // Player has no objects
    pObject = nil
    return
  }
  Scanner := bufio.NewScanner(PlayerEquFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // For each player equipment object
    ObjectId = StrGetWord(Stuff, 2)
    ObjectName = StrMakeLower(ObjectName)
    if ObjectName == ObjectId {
      // Found a match
      pObject = nil
      ObjectConstructor(ObjectId)
      if pObject != nil {
        // Object exists
        return
      } else {
        // Object does not exist, Log it
        LogBuf = ObjectId
        LogBuf += " is an invalid item found in player equipment - "
        LogBuf += "Object::IsObjInPlayerEqu"
        LogIt(LogBuf)
        pObject = nil
        return
      }
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  PlayerEquFile.Close()
  //***************************************************
  //* No match found, try getting match using 'names' *
  //***************************************************
  PlayerEquFile, err = os.Open(PlayerEquFileName)
  if err != nil {
    // Player has no objects
    pObject = nil
    return
  }
  Scanner = bufio.NewScanner(PlayerEquFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // For each player equipment object
    ObjectId = StrGetWord(Stuff, 2)
    pObject = nil
    ObjectConstructor(ObjectId)
    if pObject == nil {
      // Object does not exist, Log it
      LogBuf = ObjectId
      LogBuf += " is an invalid item found in player equipment - "
      LogBuf += "Object::IsObjInPlayerEqu"
      LogIt(LogBuf)
      pObject = nil
      return
    }
    NamesCheck = pObject.Names
    NamesCheck = StrMakeLower(NamesCheck)
    if StrIsWord(ObjectName, NamesCheck) {
      // Match
      return
    } else {
      // No match
      pObject = nil
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  PlayerEquFile.Close()
  // Object not found in player's inventory
  pObject = nil
}

// Is object in player's inventory?
func IsObjInPlayerInv(ObjectName string) {
  var NamesCheck         string
  var ObjectId           string
  var ObjectIdCheck      string
  var ObjectNameCheck    string
  var PlayerObjFileName  string
  var PlayerObjFile     *os.File

  _ = ObjectIdCheck
  _ = ObjectNameCheck

  DEBUGIT(5)
  // Open PlayerObj file
  PlayerObjFileName = PLAYER_OBJ_DIR
  PlayerObjFileName += pDnodeActor.PlayerName
  PlayerObjFileName += ".txt"
  //*******************************
  //* Try matching using ObjectId *
  //*******************************
  PlayerObjFile, err := os.Open(PlayerObjFileName)
  if err != nil {
    // Player has no objects
    pObject = nil
    return
  }
  Scanner := bufio.NewScanner(PlayerObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // For all items in player inventory
    ObjectId = StrGetWord(Stuff, 2)
    ObjectName = StrMakeLower(ObjectName)
    if ObjectName == ObjectId {
      // Found a match
      pObject = nil
      ObjectConstructor(ObjectId)
      if pObject != nil {
        // Object exists
        pObject.Count = StrGetWord(Stuff, 1)
        return
      } else {
        // Object does not exist, Log it
        LogBuf = ObjectId
        LogBuf += " is an invalid item found in player inventory - "
        LogBuf += "Object::IsObjInPlayerInv"
        LogIt(LogBuf)
        pObject = nil
        return
      }
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  PlayerObjFile.Close()
  //***************************************************
  //* No match found, try getting match using 'names' *
  //***************************************************
  PlayerObjFile, err = os.Open(PlayerObjFileName)
  if err != nil {
    // Player has no objects
    pObject = nil
    return
  }
  Scanner = bufio.NewScanner(PlayerObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // For all items in player inventory
    ObjectId = StrGetWord(Stuff, 2)
    pObject = nil
    ObjectConstructor(ObjectId)
    if pObject == nil {
      // Object does not exist, Log it
      LogBuf = ObjectId
      LogBuf += " is an invalid item found in player inventory - "
      LogBuf += "Object::IsObjInPlayerInv"
      LogIt(LogBuf)
      pObject = nil
      return
    }
    pObject.Count = StrGetWord(Stuff, 1)
    NamesCheck = pObject.Names
    NamesCheck = StrMakeLower(NamesCheck)
    if StrIsWord(ObjectName, NamesCheck) {
      // Match
      return
    } else {
      // No match
      pObject = nil
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  PlayerObjFile.Close()
  // Object not found in player's inventory
}

// Is object in room
func IsObjInRoom(ObjectName string) {
  var NamesCheck       string
  var ObjectId         string
  var ObjectIdCheck    string
  var ObjectNameCheck  string
  var RoomObjFileName  string
  var RoomObjFile     *os.File

  _ = ObjectIdCheck
  _ = ObjectNameCheck

  DEBUGIT(5)
  // Open RoomObj file
  RoomObjFileName = ROOM_OBJ_DIR
  RoomObjFileName += pDnodeActor.pPlayer.RoomId
  RoomObjFileName += ".txt"
  //*******************************
  //* Try matching using ObjectId *
  //*******************************
  RoomObjFile, err := os.Open(RoomObjFileName)
  if err != nil {
    // Room has no objects
    pObject = nil
    return
  }
  Scanner := bufio.NewScanner(RoomObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // For each item in room
    ObjectId = StrGetWord(Stuff, 2)
    ObjectName = StrMakeLower(ObjectName)
    if ObjectName == ObjectId {
      // Found a match
      pObject = nil
      ObjectConstructor(ObjectId)
      if pObject != nil {
        // Object exists
        return
      } else {
        // Object does not exist, Log it
        LogBuf = ObjectId
        LogBuf += " is an invalid item found in room - "
        LogBuf += "Object::IsObjInRoom"
        LogIt(LogBuf)
        pObject = nil
        return
      }
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  RoomObjFile.Close()
  //***************************************************
  //* No match found, try getting match using 'names' *
  //***************************************************
  RoomObjFile, err = os.Open(RoomObjFileName)
  if err != nil {
    // Room has no objects
    pObject = nil
    return
  }
  Scanner = bufio.NewScanner(RoomObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // For each item in room
    ObjectId = StrGetWord(Stuff, 2)
    pObject = nil
    ObjectConstructor(ObjectId)
    if pObject == nil {
      // Object does not exist, Log it
      LogBuf = ObjectId
      LogBuf += " is an invalid item found in room - "
      LogBuf += "Object::IsObjInRoom"
      LogIt(LogBuf)
      pObject = nil
      return
    }
    NamesCheck = pObject.Names
    NamesCheck = StrMakeLower(NamesCheck)
    if StrIsWord(ObjectName, NamesCheck) {
      // Match
      return
    } else {
      // No match
      pObject = nil
    }
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  RoomObjFile.Close()
  // Object not found in room
  pObject = nil
}

// Is this a valid object?
func IsObject(ObjectId string) {
  var ObjectFileName  string
  var ObjectFile     *os.File

  _ = ObjectFile

  DEBUGIT(5)
  ObjectFileName = OBJECTS_DIR
  ObjectFileName += ObjectId
  ObjectFileName += ".txt"
  if FileExist(ObjectFileName) {
    pObject = nil
    ObjectConstructor(ObjectId)
    return
  } else {
    pObject = nil
    return
  }
}

// Remove an object from player's equipment
func RemoveObjFromPlayerEqu(ObjectId string) {
  var BytesInFile           int64
  var ObjectIdRemoved       bool
  var ObjectIdCheck         string
  var PlayerEquFileName     string
  var PlayerEquFileNameTmp  string
  var PlayerEquFile        *os.File
  var PlayerEquFileTmp     *os.File

  DEBUGIT(5)
  ObjectId = StrMakeLower(ObjectId)
  // Open PlayerEqu file
  PlayerEquFileName = PLAYER_EQU_DIR
  PlayerEquFileName += pDnodeActor.PlayerName
  PlayerEquFileName += ".txt"
  PlayerEquFile, err := os.Open(PlayerEquFileName)
  if err != nil {
    LogBuf = "Object::RemoveObjFromPlayerEqu - Open PlayerEqu file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Open temp PlayerEqu file
  PlayerEquFileNameTmp = PLAYER_EQU_DIR
  PlayerEquFileNameTmp += pDnodeActor.PlayerName
  PlayerEquFileNameTmp += ".tmp.txt"
  PlayerEquFileTmp, err = os.Create(PlayerEquFileNameTmp)
  if err != nil {
    LogBuf = "Object::RemoveObjFromPlayerEqu - Open PlayerEqu temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Write temp PlayerEqu file
  ObjectIdRemoved = false
  Scanner := bufio.NewScanner(PlayerEquFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    if ObjectIdRemoved {
      // Object has been removed, just write the rest of the objects
      PlayerEquFileTmp.WriteString(Stuff + "\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId == ObjectIdCheck {
      // Found it, skipping it will remove it from the file
      ObjectIdRemoved = true
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    // None of the above conditions satisfied, just write it
    PlayerEquFileTmp.WriteString(Stuff + "\n")
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  if !ObjectIdRemoved {
    // Object not removed, this is definitely BAD!
    LogBuf = "Object::RemoveObjFromPlayerEqu - Object not removed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  FileInfo, err := PlayerEquFileTmp.Stat()
  if err == nil {
    BytesInFile = FileInfo.Size()
  } else {
    BytesInFile = 0
  }
  PlayerEquFile.Close()
  PlayerEquFileTmp.Close()
  Remove(PlayerEquFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    Rename(PlayerEquFileNameTmp, PlayerEquFileName)
  } else {
    // If the file is empty, delete it
    Remove(PlayerEquFileNameTmp)
  }
}

// Remove an object from player's inventory
func RemoveObjFromPlayerInv(ObjectId string, Count int) {
  var BytesInFile           int64
  var ObjectIdRemoved       bool
  var ObjectIdCheck         string
  var ObjCount              int
  var PlayerObjFileName     string
  var PlayerObjFileNameTmp  string
  var PlayerObjFile        *os.File
  var PlayerObjFileTmp     *os.File

  DEBUGIT(5)
  ObjectId = StrMakeLower(ObjectId)
  // Open PlayerObj file
  PlayerObjFileName = PLAYER_OBJ_DIR
  PlayerObjFileName += pDnodeActor.PlayerName
  PlayerObjFileName += ".txt"
  PlayerObjFile, err := os.Open(PlayerObjFileName)
  if err != nil {
    LogBuf = "Object::RemoveObjFromPlayerInv - Open PlayerObj file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Open temp PlayerObj file
  PlayerObjFileNameTmp = PLAYER_OBJ_DIR
  PlayerObjFileNameTmp += pDnodeActor.PlayerName
  PlayerObjFileNameTmp += ".tmp.txt"
  PlayerObjFileTmp, err = os.Create(PlayerObjFileNameTmp)
  if err != nil {
    LogBuf = "Object::RemoveObjFromPlayerInv - Open PlayerObj temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Write temp PlayerObj file
  ObjectIdRemoved = false
  Scanner := bufio.NewScanner(PlayerObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    if ObjectIdRemoved {
      // Object has been removed, just write the rest of the objects
      PlayerObjFileTmp.WriteString(Stuff + "\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId == ObjectIdCheck {
      // Found it, subtract 'count' from ObjCount
      ObjCount = StrToInt(StrGetWord(Stuff, 1))
      ObjCount -= Count
      ObjectIdRemoved = true
      if ObjCount > 0 {
        Buf = strconv.Itoa(ObjCount)
        TmpStr = Buf
        ObjectId = TmpStr + " " + ObjectId
        PlayerObjFileTmp.WriteString(ObjectId + "\n")
      }
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    // None of the above conditions satisfied, just write it
    PlayerObjFileTmp.WriteString(Stuff + "\n")
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  if !ObjectIdRemoved {
    // Object not removed, this is definitely BAD!
    LogBuf = "Object::RemoveObjFromPlayerInv - Object not removed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  FileInfo, err := PlayerObjFileTmp.Stat()
  if err == nil {
    BytesInFile = FileInfo.Size()
  } else {
    BytesInFile = 0
  }
  PlayerObjFile.Close()
  PlayerObjFileTmp.Close()
  Remove(PlayerObjFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    Rename(PlayerObjFileNameTmp, PlayerObjFileName)
  } else {
    // If the file is empty, delete it
    Remove(PlayerObjFileNameTmp)
  }
}

// Remove an object from room
func RemoveObjFromRoom(ObjectId string) {
  var BytesInFile         int64
  var ObjectIdRemoved     bool
  var ObjectIdCheck       string
  var ObjCount            int
  var RoomObjFileName     string
  var RoomObjFileNameTmp  string
  var RoomObjFile        *os.File
  var RoomObjFileTmp     *os.File

  DEBUGIT(5)
  ObjectId = StrMakeLower(ObjectId)
  // Open RoomObj file
  RoomObjFileName = ROOM_OBJ_DIR
  RoomObjFileName += pDnodeActor.pPlayer.RoomId
  RoomObjFileName += ".txt"
  RoomObjFile, err := os.Open(RoomObjFileName)
  if err != nil {
    LogBuf = "Object::RemoveObjFromRoom - Open RoomObj file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Open temp RoomObj file
  RoomObjFileNameTmp = ROOM_OBJ_DIR
  RoomObjFileNameTmp += pDnodeActor.pPlayer.RoomId
  RoomObjFileNameTmp += ".tmp.txt"
  RoomObjFileTmp, err = os.Create(RoomObjFileNameTmp)
  if err != nil {
    LogBuf = "Object::RemoveObjFromRoom - Open RoomObj temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Write temp RoomObj file
  ObjectIdRemoved = false
  Scanner := bufio.NewScanner(RoomObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    if ObjectIdRemoved {
      // Object has been removed, just write the rest of the objects
      RoomObjFileTmp.WriteString(Stuff + "\n")
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId == ObjectIdCheck {
      // Found it, subtract 1 from count
      ObjCount = StrToInt(StrGetWord(Stuff, 1))
      ObjCount--
      ObjectIdRemoved = true
      if ObjCount > 0 {
        Buf = strconv.Itoa(ObjCount)
        TmpStr = Buf
        ObjectId = TmpStr + " " + ObjectId
        RoomObjFileTmp.WriteString(ObjectId + "\n")
      }
      Scanner.Scan()
      Stuff = Scanner.Text()
      continue
    }
    // None of the above conditions satisfied, just write it
    RoomObjFileTmp.WriteString(Stuff + "\n")
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  if !ObjectIdRemoved {
    // Object not removed, this is definitely BAD!
    LogBuf = "Object::RemoveObjFromRoom - Object not removed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  FileInfo, err := RoomObjFileTmp.Stat()
  if err == nil {
    BytesInFile = FileInfo.Size()
  } else {
    BytesInFile = 0
  }
  RoomObjFile.Close()
  RoomObjFileTmp.Close()
  Remove(RoomObjFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    Rename(RoomObjFileNameTmp, RoomObjFileName)
  } else {
    // If the file is empty, delete it
    Remove(RoomObjFileNameTmp)
  }
}

// Show player equipment
func ShowPlayerEqu(pDnodeTgt1 *Dnode) {
  var ObjectId           string
  var PlayerEquFile     *os.File
  var PlayerEquFileName  string
  var WearPosition       string

  DEBUGIT(5)
  pDnodeTgt = pDnodeTgt1
  // Open PlayerEqu file
  PlayerEquFileName = PLAYER_EQU_DIR
  PlayerEquFileName += pDnodeTgt.PlayerName
  PlayerEquFileName += ".txt"
  PlayerEquFile, err := os.Open(PlayerEquFileName)
  if err != nil {
    // No player equipment
    if pDnodeActor == pDnodeTgt {
      // Player is checking their own equipment
      pDnodeActor.PlayerOut += "\r\n"
      pDnodeActor.PlayerOut += "You have no equipment!\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeTgt.pPlayer)
      return
    } else {
      // A player is looking at another player
      pDnodeActor.PlayerOut += "\r\n"
      pDnodeActor.PlayerOut += pDnodeTgt.PlayerName
      pDnodeActor.PlayerOut += " has no equipment!"
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Equipment\r\n"
  pDnodeActor.PlayerOut += "---------\r\n"
  Scanner := bufio.NewScanner(PlayerEquFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    WearPosition = StrGetWord(Stuff, 1)
    WearPosition = TranslateWord(WearPosition)
    ObjectId = StrGetWord(Stuff, 2)
    pObject = nil
    ObjectConstructor(ObjectId)
    pDnodeActor.PlayerOut += WearPosition
    pDnodeActor.PlayerOut += pObject.Desc1
    pDnodeActor.PlayerOut += "\r\n"
    pObject = nil
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  PlayerEquFile.Close()
}

// Show player inventory
func ShowPlayerInv() {
  var ObjectCount        string
  var ObjectId           string
  var PlayerObjFile     *os.File
  var PlayerObjFileName  string

  DEBUGIT(5)
  // Open PlayerObj file
  PlayerObjFileName = PLAYER_OBJ_DIR
  PlayerObjFileName += pDnodeActor.PlayerName
  PlayerObjFileName += ".txt"
  PlayerObjFile, err := os.Open(PlayerObjFileName)
  if err != nil {
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "It is sad, but you have nothing in your inventory.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Inventory\r\n"
  pDnodeActor.PlayerOut += "---------\r\n"
  Scanner := bufio.NewScanner(PlayerObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    ObjectCount = StrGetWord(Stuff, 1)
    ObjectId = StrGetWord(Stuff, 2)
    pObject = nil
    ObjectConstructor(ObjectId)
    pDnodeActor.PlayerOut += "(" + ObjectCount + ") "
    pDnodeActor.PlayerOut += pObject.Desc1
    pDnodeActor.PlayerOut += "\r\n"
    pObject = nil
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  PlayerObjFile.Close()
}

// Show objects in room
func ShowObjsInRoom(pDnode *Dnode) {
  var ObjectCount      string
  var ObjectId         string
  var RoomObjFile     *os.File
  var RoomObjFileName  string

  DEBUGIT(5)
  // Open RoomObj file
  RoomObjFileName = ROOM_OBJ_DIR
  RoomObjFileName += pDnode.pPlayer.RoomId
  RoomObjFileName += ".txt"
  RoomObjFile, err := os.Open(RoomObjFileName)
  if err != nil {
    // No objects in room to display
    return
  }
  Scanner := bufio.NewScanner(RoomObjFile)
  Scanner.Scan()
  Stuff = Scanner.Text()
  for Stuff != "" {
    // For each object in the room
    ObjectCount = StrGetWord(Stuff, 1)
    ObjectId = StrGetWord(Stuff, 2)
    pObject = nil
    ObjectConstructor(ObjectId)
    pObject.Type = StrMakeLower(pObject.Type)
    pDnode.PlayerOut += "\r\n"
    if pObject.Type != "notake" {
      // Should be only 1 NoTake type object in a room, like signs or statues
      pDnode.PlayerOut += "(" + ObjectCount + ") "
    }
    pDnode.PlayerOut += pObject.Desc2
    pObject = nil
    Scanner.Scan()
    Stuff = Scanner.Text()
  }
  RoomObjFile.Close()
}

// Find an object where ever it is
func WhereObj(ObjectIdSearch string) {
  DEBUGIT(5)
  WhereObjPlayerEqu(ObjectIdSearch)
  WhereObjPlayerObj(ObjectIdSearch)
  WhereObjRoomObj(ObjectIdSearch)
}

// Where is object in PlayerEqu
func WhereObjPlayerEqu(ObjectIdSearch string) {
  var FileName           string
  var ObjectId           string
  var PlayerEquFileName  string
  var PlayerEquFile     *os.File
  var PlayerName         string

  DEBUGIT(5)
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Objects in player equipment"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "---------------------------"
  pDnodeActor.PlayerOut += "\r\n"
  if ChgDir(PLAYER_EQU_DIR) != nil {
    LogBuf = "Object::WhereObjPlayerEqu - Change directory to PLAYER_EQU_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogBuf = "Object::WhereObjPlayerEqu - Change directory to PLAYER_EQU_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    FileName = DirEntry.Name()
    // Open PlayerEqu file
    PlayerEquFileName = FileName
    PlayerEquFile, err = os.Open(PlayerEquFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogBuf = "Object::WhereObjPlayerEqu - Open PlayerEqu file failed"
      LogIt(LogBuf)
      os.Exit(1)
    }
    PlayerName = StrLeft(FileName, StrGetLength(FileName)-4)
    Scanner := bufio.NewScanner(PlayerEquFile)
    Scanner.Scan()
    Stuff = Scanner.Text()
    for Stuff != "" {
      ObjectId = StrGetWord(Stuff, 2)
      if ObjectId == ObjectIdSearch {
        pDnodeActor.PlayerOut += PlayerName
        pDnodeActor.PlayerOut += " "
        pDnodeActor.PlayerOut += Stuff
        pDnodeActor.PlayerOut += "&N"
        pDnodeActor.PlayerOut += "\r\n"
      }
      Scanner.Scan()
      Stuff = Scanner.Text()
    }
    PlayerEquFile.Close()
  }
  if ChgDir(HomeDir) != nil {
    LogBuf = "Object::WhereObjPlayerEqu - Change directory to HomeDir failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
}

// Where is object in PlayerObj
func WhereObjPlayerObj(ObjectIdSearch string) {
  var FileName           string
  var ObjectId           string
  var PlayerObjFileName  string
  var PlayerObjFile     *os.File
  var PlayerName         string

  DEBUGIT(5)
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Objects in player inventory"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "---------------------------"
  pDnodeActor.PlayerOut += "\r\n"
  if ChgDir(PLAYER_OBJ_DIR) != nil {
    LogBuf = "Object::WhereObjPlayerObj - Change directory to PLAYER_OBJ_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogBuf = "Object::WhereObjPlayerObj - Change directory to PLAYER_OBJ_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    FileName = DirEntry.Name()
    // Open PlayerObj file
    PlayerObjFileName = FileName
    PlayerObjFile, err = os.Open(PlayerObjFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogBuf = "Object::WhereObjPlayerObj - Open PlayerObj file failed"
      LogIt(LogBuf)
      os.Exit(1)
    }
    PlayerName = StrLeft(FileName, StrGetLength(FileName)-4)
    Scanner := bufio.NewScanner(PlayerObjFile)
    Scanner.Scan()
    Stuff = Scanner.Text()
    for Stuff != "" {
      ObjectId = StrGetWord(Stuff, 2)
      if ObjectId == ObjectIdSearch {
        pDnodeActor.PlayerOut += PlayerName
        pDnodeActor.PlayerOut += " "
        pDnodeActor.PlayerOut += Stuff
        pDnodeActor.PlayerOut += "&N"
        pDnodeActor.PlayerOut += "\r\n"
      }
      Scanner.Scan()
      Stuff = Scanner.Text()
    }
    PlayerObjFile.Close()
  }
  if ChgDir(HomeDir) != nil {
    LogBuf = "Object::WhereObjPlayerObj - Change directory to HomeDir failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
}

// Where is object in RoomObj
func WhereObjRoomObj(ObjectIdSearch string) {
  var FileName         string
  var ObjectId         string
  var RoomName         string
  var RoomObjFileName  string
  var RoomObjFile     *os.File

  DEBUGIT(5)
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Objects in rooms"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "----------------"
  pDnodeActor.PlayerOut += "\r\n"
  if ChgDir(ROOM_OBJ_DIR) != nil {
    LogBuf = "Object::WhereObjRoomObj - Change directory to ROOM_OBJ_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  DirEntries, err := os.ReadDir("./")
  if err != nil {
    LogBuf = "Object::WhereObjRoomObj - Change directory to ROOM_OBJ_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  for _, DirEntry := range DirEntries {
    if DirEntry.IsDir() {
      continue
    }
    FileName = DirEntry.Name()
    // Open RoomObj file
    RoomObjFileName = FileName
    RoomObjFile, err = os.Open(RoomObjFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogBuf = "Object::WhereObj - Open RoomObj file failed"
      LogIt(LogBuf)
      os.Exit(1)
    }
    RoomName = StrLeft(FileName, StrGetLength(FileName)-4)
    Scanner := bufio.NewScanner(RoomObjFile)
    Scanner.Scan()
    Stuff = Scanner.Text()
    for Stuff != "" {
      // For each room object
      ObjectId = StrGetWord(Stuff, 2)
      if ObjectId == ObjectIdSearch {
        // Match
        pDnodeActor.PlayerOut += RoomName
        pDnodeActor.PlayerOut += " "
        pDnodeActor.PlayerOut += Stuff
        pDnodeActor.PlayerOut += "&N"
        pDnodeActor.PlayerOut += "\r\n"
      }
      Scanner.Scan()
      Stuff = Scanner.Text()
    }
    RoomObjFile.Close()
  }
  if ChgDir(HomeDir) != nil {
    LogBuf = "Object::WhereObj - Change directory to HomeDir failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
}

// Examine object
func ExamineObj(ObjectId string) {
  DEBUGIT(5)
  OpenObjectFile(ObjectId)
  for Stuff != "Desc3:" {
    ReadObjLine() // Do not use ReadLine() here
  }
  // Object Description 3
  ReadObjLine() // Do not use ReadLine() here
  for Stuff != "End Desc3" {
    pDnodeActor.PlayerOut += Stuff
    pDnodeActor.PlayerOut += "\r\n"
    ReadObjLine() // Do not use ReadLine() here
  }
  pDnodeActor.PlayerOut += "&N"
  CloseObjectFile()
}

// Close object file
func CloseObjectFile() {
  DEBUGIT(5)
  ObjectFile.Close()
  ObjectFile = nil
}

// Open Object file
func OpenObjectFile(ObjectId string) {
  var ObjectFileName string

  DEBUGIT(5)
  ObjectFileName = OBJECTS_DIR
  ObjectFileName += ObjectId
  ObjectFileName += ".txt"
  ObjectFile, err := os.Open(ObjectFileName)
  if err != nil {
    LogBuf = "Object::OpenFile - Object does not exist!"
    LogIt(LogBuf)
    os.Exit(1)
  }
  ObjScanner = bufio.NewScanner(ObjectFile)
}

// Parse object stuff
func ParseObjectStuff() {
  DEBUGIT(5)
  ReadObjLine()
  for Stuff != "" {
    if StrLeft(Stuff, 9) == "ObjectId:" {
      pObject.ObjectId = StrRight(Stuff, StrGetLength(Stuff)-9)
      pObject.ObjectId = StrTrimLeft(pObject.ObjectId)
    } else
    if StrLeft(Stuff, 6) == "Names:" {
      pObject.Names = StrRight(Stuff, StrGetLength(Stuff)-6)
      pObject.Names = StrTrimLeft(pObject.Names)
    } else
    if StrLeft(Stuff, 6) == "Desc1:" {
      pObject.Desc1 = StrRight(Stuff, StrGetLength(Stuff)-6)
      pObject.Desc1 = StrTrimLeft(pObject.Desc1)
    } else
    if StrLeft(Stuff, 6) == "Desc2:" {
      pObject.Desc2 = StrRight(Stuff, StrGetLength(Stuff)-6)
      pObject.Desc2 = StrTrimLeft(pObject.Desc2)
    } else
    if StrLeft(Stuff, 6) == "Desc3:" {
      // Desc3 can be multi-line and is dealt with in 'ExamineObj'
    } else
    if StrLeft(Stuff, 7) == "Weight:" {
      pObject.Weight = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-7))
    } else
    if StrLeft(Stuff, 5) == "Cost:" {
      pObject.Cost = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-5))
    } else
    if StrLeft(Stuff, 5) == "Type:" {
      pObject.Type = StrRight(Stuff, StrGetLength(Stuff)-5)
      pObject.Type = StrTrimLeft(pObject.Type)
    } else
    if StrLeft(Stuff, 11) == "ArmorValue:" {
      pObject.ArmorValue = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-11))
    } else
    if StrLeft(Stuff, 10) == "ArmorWear:" {
      pObject.ArmorWear = StrRight(Stuff, StrGetLength(Stuff)-10)
      pObject.ArmorWear = StrTrimLeft(pObject.ArmorWear)
      pObject.WearPosition = pObject.ArmorWear
      pObject.WearPosition = StrMakeLower(pObject.WearPosition)
    } else
    if StrLeft(Stuff, 11) == "WeaponType:" {
      pObject.WeaponType = StrRight(Stuff, StrGetLength(Stuff)-11)
      pObject.WeaponType = StrTrimLeft(pObject.WeaponType)
      pObject.ArmorWear = "wielded"
      pObject.WearPosition = "wielded"
    } else
    if StrLeft(Stuff, 13) == "WeaponDamage:" {
      pObject.WeaponDamage = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-13))
    } else
    if StrLeft(Stuff, 8) == "FoodPct:" {
      pObject.FoodPct = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-8))
    } else
    if StrLeft(Stuff, 9) == "DrinkPct:" {
      pObject.DrinkPct = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-9))
    }
    ReadObjLine()
  }
}

// Read a line from the object file
func ReadObjLine() {
  DEBUGIT(5)
  Stuff = ""
  if ObjScanner.Scan() {
    Stuff = ObjScanner.Text()
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
  }
}
