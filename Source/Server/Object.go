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
  // TODO: implement function logic
}

// Is object in player's inventory?
func IsObjInPlayerInv(ObjectName string) {
  // TODO: implement function logic
}

// Is object in room?
func IsObjInRoom(ObjectName string) {
  // TODO: implement function logic
}

// Is this a valid object?
func IsObject(ObjectId string) {
  // TODO: implement function logic
}

// Remove an object from player's equipment
func RemoveObjFromPlayerEqu(ObjectId string) {
  // TODO: implement function logic
}

// Remove an object from player's inventory
func RemoveObjFromPlayerInv(ObjectId string, Count int) {
  // TODO: implement function logic
}

// Remove an object from room
func RemoveObjFromRoom(ObjectId string) {
  // TODO: implement function logic
}

// Show player equipment
func ShowPlayerEqu(pDnodeTgt1 *Dnode) {
  // TODO: implement function logic
}

// Show player inventory
func ShowPlayerInv() {
  // TODO: implement function logic
}

// Show objects in room
func ShowObjsInRoom(pDnode *Dnode) {
  // TODO: implement function logic
}

// Find an object wherever it is
func WhereObj(ObjectIdSearch string) {
  WhereObjPlayerEqu(ObjectIdSearch)
  WhereObjPlayerObj(ObjectIdSearch)
  WhereObjRoomObj(ObjectIdSearch)
}

// Where is object in PlayerEqu
func WhereObjPlayerEqu(ObjectIdSearch string) {
  // TODO: implement function logic
}

// Where is object in PlayerObj
func WhereObjPlayerObj(ObjectIdSearch string) {
  // TODO: implement function logic
}

// Where is object in RoomObj
func WhereObjRoomObj(ObjectIdSearch string) {
  // TODO: implement function logic
}

// Examine object
func ExamineObj(ObjectId string) {
  // TODO: implement function logic
}

// Close object file
func CloseObjectFile() {
  ObjectFile.Close()
}

// Open Object file
func OpenObjectFile(ObjectId string) {
  // TODO: implement function logic
}

// Parse object stuff
func ParseObjectStuff() {
  // TODO: implement function logic
}

// Read a line from the object file
func ReadObjLine() {
  if ObjScanner.Scan() {
    Stuff = ObjScanner.Text()
    Stuff = StrTrimLeft(Stuff)
    Stuff = StrTrimRight(Stuff)
  }
}