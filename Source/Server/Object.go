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

// Global instance of Object
var pObject = &Object{
  ObjectId:          "",
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

var ObjScanner = bufio.NewScanner(ObjectFile)

// Create a new object
func ObjectConstructor() {
  // Init object structure
  pObject.ArmorValue        = 0;
  pObject.ArmorWear         = "";
  pObject.ContainerCapacity = 0;
  pObject.Cost              = 0;
  pObject.Count             = "1";
  pObject.Desc1             = "";
  pObject.Desc2             = "";
  pObject.Desc3             = "";
  pObject.DrinkPct          = 0;
  pObject.FoodPct           = 0;
  pObject.LightHours        = 0;
  pObject.Names             = "";
  pObject.Type              = "";
  pObject.WeaponType        = "";
  pObject.WeaponDamage      = 0;
  pObject.WearPosition      = "";
  pObject.Weight            = 0;
  // Populate object structure
  OpenObjectFile(ObjectId);
  ParseObjectStuff();
  CloseObjectFile();
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
  var err                   error

  WearWieldFailed = false
  ObjectId = StrMakeLower(ObjectId)
  // Open PlayerEqu file
  PlayerEquFileName = PLAYER_EQU_DIR + pDnodeActor.PlayerName + ".txt"
  NewPlayerEquFile = false
  PlayerEquFile, err = os.Open(PlayerEquFileName)
  if err != nil {
    NewPlayerEquFile = true
  }
  // Open temp PlayerEqu file
  PlayerEquFileNameTmp = PLAYER_EQU_DIR + pDnodeActor.PlayerName + ".tmp.txt"
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
    os.Rename(PlayerEquFileNameTmp, PlayerEquFileName)
    return WearWieldFailed
  }
  defer PlayerEquFile.Close()
  // Write temp PlayerEqu file
  ObjectIdAdded = false
  Scanner := bufio.NewScanner(PlayerEquFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    if ObjectIdAdded {
      // New object has been written, just write the rest of the objects
      PlayerEquFileTmp.WriteString(Stuff + "\n")
      continue
    }
    WearPositionCheck = StrGetWord(Stuff, 1)
    if WearPosition < WearPositionCheck {
      // Add new object in alphabetical order by translated WearPosition
      ObjectId = WearPosition + " " + ObjectId
      PlayerEquFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      PlayerEquFileTmp.WriteString(Stuff + "\n")
      continue
    }
    if WearPosition == WearPositionCheck {
      // Already wearing an object in that position
      WearWieldFailed = true
      ObjectIdAdded = true // Not really added
      PlayerEquFileTmp.WriteString(Stuff + "\n")
      continue
    }
    // None of the above conditions satisfied, just write it
    PlayerEquFileTmp.WriteString(Stuff + "\n")
  }
  if !ObjectIdAdded {
    // New object is alphabetically last
    ObjectId = WearPosition + " " + ObjectId
    PlayerEquFileTmp.WriteString(ObjectId + "\n")
    ObjectIdAdded = true
  }
  PlayerEquFileTmp.Close()
  os.Remove(PlayerEquFileName)
  os.Rename(PlayerEquFileNameTmp, PlayerEquFileName)
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
  var err                   error

  pDnodeTgt = pDnodeTgt1
  ObjectId = StrMakeLower(ObjectId)
  // Open PlayerObj file
  PlayerObjFileName = PLAYER_OBJ_DIR + pDnodeTgt.PlayerName + ".txt"
  NewPlayerObjFile = false
  PlayerObjFile, err = os.Open(PlayerObjFileName)
  if err != nil {
    NewPlayerObjFile = true
  }
  // Open temp PlayerObj file
  PlayerObjFileNameTmp = PLAYER_OBJ_DIR + pDnodeTgt.PlayerName + ".tmp.txt"
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
    os.Rename(PlayerObjFileNameTmp, PlayerObjFileName)
    return
  }
  defer PlayerObjFile.Close()
  // Write temp PlayerObj file
  ObjectIdAdded = false
  Scanner := bufio.NewScanner(PlayerObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    if ObjectIdAdded {
      // New object has been written, just write the rest of the objects
      PlayerObjFileTmp.WriteString(Stuff + "\n")
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId < ObjectIdCheck {
      // Add new object in alphabetical order
      ObjectId = "1 " + ObjectId
      PlayerObjFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      PlayerObjFileTmp.WriteString(Stuff + "\n")
      continue
    }
    if ObjectId == ObjectIdCheck {
      // Existing object same as new object, add 1 to count
      ObjCount = StrToInt(StrGetWord(Stuff, 1))
      ObjCount++
      ObjectId = strconv.Itoa(ObjCount) + " " + ObjectId
      PlayerObjFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      continue
    }
    // None of the above conditions satisfied, just write it
    PlayerObjFileTmp.WriteString(Stuff + "\n")
  }
  if !ObjectIdAdded {
    // New object is alphabetically last
    ObjectId = "1 " + ObjectId
    PlayerObjFileTmp.WriteString(ObjectId + "\n")
    ObjectIdAdded = true
  }
  PlayerObjFileTmp.Close()
  os.Remove(PlayerObjFileName)
  os.Rename(PlayerObjFileNameTmp, PlayerObjFileName)
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
  var err                 error

  ObjectId = StrMakeLower(ObjectId)
  // Open RoomObj file
  RoomObjFileName = ROOM_OBJ_DIR + RoomId + ".txt"
  NewRoomObjFile = false
  RoomObjFile, err = os.Open(RoomObjFileName)
  if err != nil {
    NewRoomObjFile = true
  }
  // Open temp RoomObj file
  RoomObjFileNameTmp = ROOM_OBJ_DIR + RoomId + ".tmp.txt"
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
    err = os.Rename(RoomObjFileNameTmp, RoomObjFileName)
    if err != nil {
      LogBuf = "Object::AddObjToRoom - Rename RoomObj temp file failed"
      LogIt(LogBuf)
      os.Exit(1)
    }
    return
  }
  defer RoomObjFile.Close()
  // Write temp RoomObj file
  ObjectIdAdded = false
  Scanner := bufio.NewScanner(RoomObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    if ObjectIdAdded {
      // New object has been written, just write the rest of the objects
      RoomObjFileTmp.WriteString(Stuff + "\n")
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId < ObjectIdCheck {
      // Add new object in alphabetical order
      ObjectId = "1 " + ObjectId
      RoomObjFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      RoomObjFileTmp.WriteString(Stuff + "\n")
      continue
    }
    if ObjectId == ObjectIdCheck {
      // Existing object same as new object, add 1 to count
      ObjCount = StrToInt(StrGetWord(Stuff, 1))
      ObjCount++
      ObjectId = strconv.Itoa(ObjCount) + " " + ObjectId
      RoomObjFileTmp.WriteString(ObjectId + "\n")
      ObjectIdAdded = true
      continue
    }
    // None of the above conditions satisfied, just write it
    RoomObjFileTmp.WriteString(Stuff + "\n")
  }
  if !ObjectIdAdded {
    // New object is alphabetically last
    ObjectId = "1 " + ObjectId
    RoomObjFileTmp.WriteString(ObjectId + "\n")
    ObjectIdAdded = true
  }
  RoomObjFileTmp.Close()
  os.Remove(RoomObjFileName)
  os.Rename(RoomObjFileNameTmp, RoomObjFileName)
}

// Calculate player armor class
func CalcPlayerArmorClass() int {
  var ArmorClass        int
  var PlayerEquFile    *os.File
  var PlayerEquFileName string
  var err               error

  ArmorClass = 0
  // Open PlayerObj file
  PlayerEquFileName = PLAYER_EQU_DIR + pDnodeActor.PlayerName + ".txt"
  PlayerEquFile, err = os.Open(PlayerEquFileName)
  if err != nil {
    // No player equipment
    return ArmorClass
  }
  defer PlayerEquFile.Close()
  Scanner := bufio.NewScanner(PlayerEquFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    ObjectId = StrGetWord(Stuff, 2)
    ObjectConstructor()
    ArmorClass += pObject.ArmorValue
  }
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
  var Scanner           *bufio.Scanner
  var err                error

  _ = ObjectIdCheck
  _ = ObjectNameCheck

  // Open PlayerEqu file
  PlayerEquFileName = PLAYER_EQU_DIR + pDnodeActor.PlayerName + ".txt"
  // Try matching using ObjectId
  PlayerEquFile, err = os.Open(PlayerEquFileName)
  if err != nil {
    // Player has no objects
    pObject = nil
    return
  }
  ObjectName = StrMakeLower(ObjectName)
  Scanner = bufio.NewScanner(PlayerEquFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    ObjectId = StrGetWord(Stuff, 2)
    if ObjectName == ObjectId {
      ObjectConstructor()
      return
    }
  }
  PlayerEquFile.Close()
  // No match found, try getting match using 'names'
  PlayerEquFile, err = os.Open(PlayerEquFileName)
  if err != nil {
    // Player has no objects
    pObject = nil
    return
  }
  Scanner = bufio.NewScanner(PlayerEquFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    ObjectId = StrGetWord(Stuff, 2)
    ObjectConstructor()
    NamesCheck = pObject.Names
    NamesCheck = StrMakeLower(NamesCheck)
    if StrIsWord(ObjectName, NamesCheck) {
      return
    } else {
      pObject = nil
    }
  }
  PlayerEquFile.Close()
  // Object not found in player's inventory
  pObject = nil
  return
}

// Is object in player's inventory?
func IsObjInPlayerInv(ObjectName string) {
  var NamesCheck         string
  var ObjectId           string
  var PlayerObjFileName  string
  var PlayerObjFile     *os.File
  var Scanner           *bufio.Scanner
  var err                error

  // Open PlayerObj file
  PlayerObjFileName = PLAYER_OBJ_DIR + pDnodeActor.PlayerName + ".txt"
  // Try matching using ObjectId
  PlayerObjFile, err = os.Open(PlayerObjFileName)
  if err != nil {
    // Player has no objects
    pObject = nil
    return
  }
  ObjectName = StrMakeLower(ObjectName)
  Scanner = bufio.NewScanner(PlayerObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    ObjectId = StrGetWord(Stuff, 2)
    if ObjectName == ObjectId {
      ObjectConstructor()
      pObject.Count = StrGetWord(Stuff, 1)
      return
    }
  }
  PlayerObjFile.Close()
  // No match found, try getting match using 'names'
  PlayerObjFile, err = os.Open(PlayerObjFileName)
  if err != nil {
    // Player has no objects
    pObject = nil
    return
  }
  Scanner = bufio.NewScanner(PlayerObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    ObjectId = StrGetWord(Stuff, 2)
    ObjectConstructor()
    pObject.Count = StrGetWord(Stuff, 1)
    NamesCheck = pObject.Names
    NamesCheck = StrMakeLower(NamesCheck)
    if StrIsWord(ObjectName, NamesCheck) {
      return
    } else {
      pObject = nil
    }
  }
  PlayerObjFile.Close()
  // Object not found in player's inventory
  return
}

// Is object in room
func IsObjInRoom(ObjectName string) {
  var NamesCheck       string
  var ObjectId         string
  var RoomObjFileName  string
  var RoomObjFile     *os.File
  var Scanner         *bufio.Scanner
  var err              error

  // Open RoomObj file
  RoomObjFileName = ROOM_OBJ_DIR + pDnodeActor.pPlayer.RoomId + ".txt"
  // Try matching using ObjectId
  RoomObjFile, err = os.Open(RoomObjFileName)
  if err != nil {
    // Room has no objects
    pObject = nil
    return
  }
  ObjectName = StrMakeLower(ObjectName)
  Scanner = bufio.NewScanner(RoomObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    ObjectId = StrGetWord(Stuff, 2)
    if ObjectName == ObjectId {
      ObjectConstructor()
      return
    }
  }
  RoomObjFile.Close()
  // No match found, try getting match using 'names'
  RoomObjFile, err = os.Open(RoomObjFileName)
  if err != nil {
    // Room has no objects
    pObject = nil
    return
  }
  Scanner = bufio.NewScanner(RoomObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    ObjectId = StrGetWord(Stuff, 2)
    ObjectConstructor()
    NamesCheck = pObject.Names
    NamesCheck = StrMakeLower(NamesCheck)
    if StrIsWord(ObjectName, NamesCheck) {
      return
    } else {
      pObject = nil
    }
  }
  RoomObjFile.Close()
  // Object not found in room
  pObject = nil
  return
}

// Is this a valid object?
func IsObject(ObjectId string) {
  var ObjectFileName string

  ObjectFileName = OBJECTS_DIR + ObjectId + ".txt"
  if FileExist(ObjectFileName) {
    pObject = &Object{ObjectId: ObjectId}
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
  var FileInfo              os.FileInfo
  var err                   error

  ObjectId = StrMakeLower(ObjectId)
  // Open PlayerEqu file
  PlayerEquFileName = PLAYER_EQU_DIR + pDnodeActor.PlayerName + ".txt"
  PlayerEquFile, err = os.Open(PlayerEquFileName)
  if err != nil {
    LogBuf = "Object::RemoveObjFromPlayerEqu - Open PlayerEqu file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  defer PlayerEquFile.Close()
  // Open temp PlayerEqu file
  PlayerEquFileNameTmp = PLAYER_EQU_DIR + pDnodeActor.PlayerName + ".tmp.txt"
  PlayerEquFileTmp, err = os.Create(PlayerEquFileNameTmp)
  if err != nil {
    LogBuf = "Object::RemoveObjFromPlayerEqu - Open PlayerEqu temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Write temp PlayerEqu file
  ObjectIdRemoved = false
  Scanner := bufio.NewScanner(PlayerEquFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    if ObjectIdRemoved {
      // Object has been removed, just write the rest of the objects
      PlayerEquFileTmp.WriteString(Stuff + "\n")
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId == ObjectIdCheck {
      // Found it, skipping it will remove it from the file
      ObjectIdRemoved = true
      continue
    }
    // None of the above conditions satisfied, just write it
    PlayerEquFileTmp.WriteString(Stuff + "\n")
  }
  if !ObjectIdRemoved {
    // Object not removed, this is definitely BAD!
    LogBuf = "Object::RemoveObjFromPlayerEqu - Object not removed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  FileInfo, err = PlayerEquFileTmp.Stat()
  if err == nil {
    BytesInFile = FileInfo.Size()
  } else {
    BytesInFile = 0
  }
  PlayerEquFileTmp.Close()
  os.Remove(PlayerEquFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    os.Rename(PlayerEquFileNameTmp, PlayerEquFileName)
  } else {
    // If the file is empty, delete it
    os.Remove(PlayerEquFileNameTmp)
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
  var FileInfo              os.FileInfo
  var err                   error

  ObjectId = StrMakeLower(ObjectId)
  // Open PlayerObj file
  PlayerObjFileName = PLAYER_OBJ_DIR + pDnodeActor.PlayerName + ".txt"
  PlayerObjFile, err = os.Open(PlayerObjFileName)
  if err != nil {
    LogBuf = "Object::RemoveObjFromPlayerInv - Open PlayerObj file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  defer PlayerObjFile.Close()
  // Open temp PlayerObj file
  PlayerObjFileNameTmp = PLAYER_OBJ_DIR + pDnodeActor.PlayerName + ".tmp.txt"
  PlayerObjFileTmp, err = os.Create(PlayerObjFileNameTmp)
  if err != nil {
    LogBuf = "Object::RemoveObjFromPlayerInv - Open PlayerObj temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Write temp PlayerObj file
  ObjectIdRemoved = false
  Scanner := bufio.NewScanner(PlayerObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    if ObjectIdRemoved {
      // Object has been removed, just write the rest of the objects
      PlayerObjFileTmp.WriteString(Stuff + "\n")
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId == ObjectIdCheck {
      // Found it, subtract 'count' from ObjCount
      ObjCount = StrToInt(StrGetWord(Stuff, 1))
      ObjCount -= Count
      ObjectIdRemoved = true
      if ObjCount > 0 {
        ObjectId = strconv.Itoa(ObjCount) + " " + ObjectId
        PlayerObjFileTmp.WriteString(ObjectId + "\n")
      }
      continue
    }
    // None of the above conditions satisfied, just write it
    PlayerObjFileTmp.WriteString(Stuff + "\n")
  }
  if !ObjectIdRemoved {
    // Object not removed, this is definitely BAD!
    LogBuf = "Object::RemoveObjFromPlayerInv - Object not removed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  FileInfo, err = PlayerObjFileTmp.Stat()
  if err == nil {
    BytesInFile = FileInfo.Size()
  } else {
    BytesInFile = 0
  }
  PlayerObjFileTmp.Close()
  os.Remove(PlayerObjFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    os.Rename(PlayerObjFileNameTmp, PlayerObjFileName)
  } else {
    // If the file is empty, delete it
    os.Remove(PlayerObjFileNameTmp)
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
  var FileInfo            os.FileInfo
  var err                 error

  ObjectId = StrMakeLower(ObjectId)
  // Open RoomObj file
  RoomObjFileName = ROOM_OBJ_DIR + pDnodeActor.pPlayer.RoomId + ".txt"
  RoomObjFile, err = os.Open(RoomObjFileName)
  if err != nil {
    LogBuf = "Object::RemoveObjFromRoom - Open RoomObj file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  defer RoomObjFile.Close()
  // Open temp RoomObj file
  RoomObjFileNameTmp = ROOM_OBJ_DIR + pDnodeActor.pPlayer.RoomId + ".tmp.txt"
  RoomObjFileTmp, err = os.Create(RoomObjFileNameTmp)
  if err != nil {
    LogBuf = "Object::RemoveObjFromRoom - Open RoomObj temp file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Write temp RoomObj file
  ObjectIdRemoved = false
  Scanner := bufio.NewScanner(RoomObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    if ObjectIdRemoved {
      // Object has been removed, just write the rest of the objects
      RoomObjFileTmp.WriteString(Stuff + "\n")
      continue
    }
    ObjectIdCheck = StrGetWord(Stuff, 2)
    if ObjectId == ObjectIdCheck {
      // Found it, subtract 1 from count
      ObjCount = StrToInt(StrGetWord(Stuff, 1))
      ObjCount--
      ObjectIdRemoved = true
      if ObjCount > 0 {
        ObjectId = strconv.Itoa(ObjCount) + " " + ObjectId
        RoomObjFileTmp.WriteString(ObjectId + "\n")
      }
      continue
    }
    // None of the above conditions satisfied, just write it
    RoomObjFileTmp.WriteString(Stuff + "\n")
  }
  if !ObjectIdRemoved {
    // Object not removed, this is definitely BAD!
    LogBuf = "Object::RemoveObjFromRoom - Object not removed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  FileInfo, err = RoomObjFileTmp.Stat()
  if err == nil {
    BytesInFile = FileInfo.Size()
  } else {
    BytesInFile = 0
  }
  RoomObjFileTmp.Close()
  os.Remove(RoomObjFileName)
  if BytesInFile > 0 {
    // If the file is not empty, rename it
    os.Rename(RoomObjFileNameTmp, RoomObjFileName)
  } else {
    // If the file is empty, delete it
    os.Remove(RoomObjFileNameTmp)
  }
}

// Show player equipment
func ShowPlayerEqu(pDnodeTgt1 *Dnode) {
  var PlayerEquFile     *os.File
  var PlayerEquFileName  string
  var WearPosition       string
  var Scanner           *bufio.Scanner
  var err                error

  pDnodeTgt = pDnodeTgt1
  // Open PlayerEqu file
  PlayerEquFileName = PLAYER_EQU_DIR + pDnodeTgt.PlayerName + ".txt"
  PlayerEquFile, err = os.Open(PlayerEquFileName)
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
  defer PlayerEquFile.Close()
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Equipment\r\n"
  pDnodeActor.PlayerOut += "---------\r\n"
  Scanner = bufio.NewScanner(PlayerEquFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    WearPosition = StrGetWord(Stuff, 1)
    WearPosition = TranslateWord(WearPosition)
    ObjectId = StrGetWord(Stuff, 2)
    ObjectConstructor()
    pDnodeActor.PlayerOut += WearPosition
    pDnodeActor.PlayerOut += pObject.Desc1
    pDnodeActor.PlayerOut += "\r\n"
    pObject = nil
  }
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Show player inventory
func ShowPlayerInv() {
  var ObjectCount        string
  var PlayerObjFile     *os.File
  var PlayerObjFileName  string
  var Scanner           *bufio.Scanner
  var err                error

  // Open PlayerObj file
  PlayerObjFileName = PLAYER_OBJ_DIR + pDnodeActor.PlayerName + ".txt"
  PlayerObjFile, err = os.Open(PlayerObjFileName)
  if err != nil {
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "It is sad, but you have nothing in your inventory.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  defer PlayerObjFile.Close()
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Inventory\r\n"
  pDnodeActor.PlayerOut += "---------\r\n"
  Scanner = bufio.NewScanner(PlayerObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    ObjectCount = StrGetWord(Stuff, 1)
    ObjectId = StrGetWord(Stuff, 2)
    ObjectConstructor()
    pDnodeActor.PlayerOut += "(" + ObjectCount + ") "
    pDnodeActor.PlayerOut += pObject.Desc1
    pDnodeActor.PlayerOut += "\r\n"
    pObject = nil
  }
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Show objects in room
func ShowObjsInRoom(pDnode *Dnode) {
  var ObjectCount      string
  var RoomObjFile     *os.File
  var RoomObjFileName  string
  var Scanner         *bufio.Scanner
  var err              error

  // Open RoomObj file
  RoomObjFileName = ROOM_OBJ_DIR + pDnode.pPlayer.RoomId + ".txt"
  RoomObjFile, err = os.Open(RoomObjFileName)
  if err != nil {
    // No objects in room to display
    return
  }
  defer RoomObjFile.Close()
  Scanner = bufio.NewScanner(RoomObjFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "" {
      continue
    }
    // For each object in the room
    ObjectCount = StrGetWord(Stuff, 1)
    ObjectId = StrGetWord(Stuff, 2)
    ObjectConstructor()
    pObject.Type = StrMakeLower(pObject.Type)
    pDnode.PlayerOut += "\r\n"
    if pObject.Type != "notake" {
      // Should be only 1 NoTake type object in a room, like signs or statues
      pDnode.PlayerOut += "(" + ObjectCount + ") "
    }
    pDnode.PlayerOut += pObject.Desc2
    pObject = nil
  }
}

// Find an object where ever it is
func WhereObj(ObjectIdSearch string) {
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
  var DirEntries         []os.DirEntry
  var Scanner           *bufio.Scanner
  var err                error

  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Objects in player equipment"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "---------------------------"
  pDnodeActor.PlayerOut += "\r\n"
  DirEntries, err = os.ReadDir(PLAYER_EQU_DIR)
  if err != nil {
    LogBuf = "Object::WhereObjPlayerEqu - Change directory to PLAYER_EQU_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  for _, entry := range DirEntries {
    if entry.IsDir() {
      continue
    }
    FileName = entry.Name()
    // Open PlayerEqu file
    PlayerEquFileName = PLAYER_EQU_DIR + FileName
    PlayerEquFile, err = os.Open(PlayerEquFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogBuf = "Object::WhereObjPlayerEqu - Open PlayerEqu file failed"
      LogIt(LogBuf)
      os.Exit(1)
    }
    PlayerName = StrLeft(FileName, StrGetLength(FileName)-4)
    Scanner = bufio.NewScanner(PlayerEquFile)
    for Scanner.Scan() {
      Stuff = Scanner.Text()
      if Stuff == "" {
        continue
      }
      ObjectId = StrGetWord(Stuff, 2)
      if ObjectId == ObjectIdSearch {
        pDnodeActor.PlayerOut += PlayerName
        pDnodeActor.PlayerOut += " "
        pDnodeActor.PlayerOut += Stuff
        pDnodeActor.PlayerOut += "&N"
        pDnodeActor.PlayerOut += "\r\n"
      }
    }
    PlayerEquFile.Close()
  }
}

// Where is object in PlayerObj
func WhereObjPlayerObj(ObjectIdSearch string) {
  var FileName           string
  var ObjectId           string
  var PlayerObjFileName  string
  var PlayerObjFile     *os.File
  var PlayerName         string
  var DirEntries         []os.DirEntry
  var Scanner           *bufio.Scanner
  var err                error

  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Objects in player inventory"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "---------------------------"
  pDnodeActor.PlayerOut += "\r\n"
  DirEntries, err = os.ReadDir(PLAYER_OBJ_DIR)
  if err != nil {
    LogBuf = "Object::WhereObjPlayerObj - Change directory to PLAYER_OBJ_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  for _, entry := range DirEntries {
    if entry.IsDir() {
      continue
    }
    FileName = entry.Name()
    // Open PlayerObj file
    PlayerObjFileName = PLAYER_OBJ_DIR + FileName
    PlayerObjFile, err = os.Open(PlayerObjFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogBuf = "Object::WhereObjPlayerObj - Open PlayerObj file failed"
      LogIt(LogBuf)
      os.Exit(1)
    }
    PlayerName = StrLeft(FileName, StrGetLength(FileName)-4)
    Scanner = bufio.NewScanner(PlayerObjFile)
    for Scanner.Scan() {
      Stuff = Scanner.Text()
      if Stuff == "" {
        continue
      }
      ObjectId = StrGetWord(Stuff, 2)
      if ObjectId == ObjectIdSearch {
        pDnodeActor.PlayerOut += PlayerName
        pDnodeActor.PlayerOut += " "
        pDnodeActor.PlayerOut += Stuff
        pDnodeActor.PlayerOut += "&N"
        pDnodeActor.PlayerOut += "\r\n"
      }
    }
    PlayerObjFile.Close()
  }
}

// Where is object in RoomObj
func WhereObjRoomObj(ObjectIdSearch string) {
  var FileName         string
  var ObjectId         string
  var RoomName         string
  var RoomObjFileName  string
  var RoomObjFile     *os.File
  var DirEntries       []os.DirEntry
  var Scanner         *bufio.Scanner
  var err              error

  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Objects in rooms"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "----------------"
  pDnodeActor.PlayerOut += "\r\n"
  DirEntries, err = os.ReadDir(ROOM_OBJ_DIR)
  if err != nil {
    LogBuf = "Object::WhereObjRoomObj - Change directory to ROOM_OBJ_DIR failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  for _, entry := range DirEntries {
    if entry.IsDir() {
      continue
    }
    FileName = entry.Name()
    // Open RoomObj file
    RoomObjFileName = ROOM_OBJ_DIR + FileName
    RoomObjFile, err = os.Open(RoomObjFileName)
    if err != nil {
      // File does not exist - Very bad!
      LogBuf = "Object::WhereObj - Open RoomObj file failed"
      LogIt(LogBuf)
      os.Exit(1)
    }
    RoomName = StrLeft(FileName, StrGetLength(FileName)-4)
    Scanner = bufio.NewScanner(RoomObjFile)
    for Scanner.Scan() {
      Stuff = Scanner.Text()
      if Stuff == "" {
        continue
      }
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
    }
    RoomObjFile.Close()
  }
}

// Examine object
func ExamineObj(ObjectId string) {
  var ObjectFileName  string
  var err             error

  ObjectFileName = OBJECTS_DIR + ObjectId + ".txt"
  ObjectFile, err = os.Open(ObjectFileName)
  if err != nil {
    LogBuf = "Object::ExamineObj - Open object file failed"
    LogIt(LogBuf)
    os.Exit(1)
  }
  ObjScanner = bufio.NewScanner(ObjectFile)
  Stuff = ""
  for Stuff != "Desc3:" {
    if !ObjScanner.Scan() {
      break
    }
    Stuff = ObjScanner.Text()
  }
  // Object Description 3
  if ObjScanner.Scan() {
    Stuff = ObjScanner.Text()
  }
  for Stuff != "End Desc3" {
    pDnodeActor.PlayerOut += Stuff
    pDnodeActor.PlayerOut += "\r\n"
    if !ObjScanner.Scan() {
      break
    }
    Stuff = ObjScanner.Text()
  }
  pDnodeActor.PlayerOut += "&N"
  CloseObjectFile()
}

// Close object file
func CloseObjectFile() {
  ObjectFile.Close()
}

// Open Object file
func OpenObjectFile(ObjectId string) {
  var ObjectFileName  string
  var err             error

  ObjectFileName = OBJECTS_DIR + ObjectId + ".txt"
  ObjectFile, err = os.Open(ObjectFileName)
  if err != nil {
    LogBuf = "Object::OpenFile - Object does not exist!"
    LogIt(LogBuf)
    os.Exit(1)
  }
  ObjScanner = bufio.NewScanner(ObjectFile)
}

// Parse object stuff
func ParseObjectStuff() {
  ReadObjLine()
  for Stuff != "" {
    if StrLeft(Stuff, 9) == "ObjectId:" {
      pObject.ObjectId     = StrRight(Stuff, StrGetLength(Stuff)-9)
      pObject.ObjectId     = StrTrimLeft(pObject.ObjectId)
    } else if StrLeft(Stuff, 6) == "Names:" {
      pObject.Names        = StrRight(Stuff, StrGetLength(Stuff)-6)
      pObject.Names        = StrTrimLeft(pObject.Names)
    } else if StrLeft(Stuff, 6) == "Desc1:" {
      pObject.Desc1        = StrRight(Stuff, StrGetLength(Stuff)-6)
      pObject.Desc1        = StrTrimLeft(pObject.Desc1)
    } else if StrLeft(Stuff, 6) == "Desc2:" {
      pObject.Desc2        = StrRight(Stuff, StrGetLength(Stuff)-6)
      pObject.Desc2        = StrTrimLeft(pObject.Desc2)
    } else if StrLeft(Stuff, 6) == "Desc3:" {
      // Desc3 can be multi-line and is dealt with in 'ExamineObj'
    } else if StrLeft(Stuff, 7) == "Weight:" {
      pObject.Weight       = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-7))
    } else if StrLeft(Stuff, 5) == "Cost:" {
      pObject.Cost         = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-5))
    } else if StrLeft(Stuff, 5) == "Type:" {
      pObject.Type         = StrRight(Stuff, StrGetLength(Stuff)-5)
      pObject.Type         = StrTrimLeft(pObject.Type)
    } else if StrLeft(Stuff, 11) == "ArmorValue:" {
      pObject.ArmorValue   = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-11))
    } else if StrLeft(Stuff, 10) == "ArmorWear:" {
      pObject.ArmorWear    = StrRight(Stuff, StrGetLength(Stuff)-10)
      pObject.ArmorWear    = StrTrimLeft(pObject.ArmorWear)
      pObject.WearPosition = pObject.ArmorWear
      pObject.WearPosition = StrMakeLower(pObject.WearPosition)
    } else if StrLeft(Stuff, 11) == "WeaponType:" {
      pObject.WeaponType = StrRight(Stuff, StrGetLength(Stuff)-11)
      pObject.WeaponType   = StrTrimLeft(pObject.WeaponType)
      pObject.ArmorWear    = "wielded"
      pObject.WearPosition = "wielded"
    } else if StrLeft(Stuff, 13) == "WeaponDamage:" {
      pObject.WeaponDamage = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-13))
    } else if StrLeft(Stuff, 8) == "FoodPct:" {
      pObject.FoodPct = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-8))
    } else if StrLeft(Stuff, 9) == "DrinkPct:" {
      pObject.DrinkPct     = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-9))
    }
    ReadObjLine()
  }
}

// Read a line from the object file
func ReadObjLine() {
  if ObjScanner.Scan() {
    Stuff = ObjScanner.Text()
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
  }
}