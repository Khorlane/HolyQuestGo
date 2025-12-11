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
func NewObject() {
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
func AddObjToPlayerEqu(WearPosition string, ObjectId string) {
  // TODO: implement function logic
}

// Add an object to player's inventory
func AddObjToPlayerInv(pDnodeTgt1 *Dnode, ObjectId string) {
  // TODO: implement function logic
}

// Add an object to room
func AddObjToRoom(RoomId string, ObjectId string) {
  // TODO: implement function logic
}

// Calculate player armor class
func CalcPlayerArmorClass() int {
  var ArmorClass int
  return ArmorClass
}

// Close object file
func CloseObjectFile() {
  ObjectFile.Close()
}

// Examine object
func ExamineObj(ObjectId string) {
  // TODO: implement function logic
}

// Is this a valid object?
func IsObject(ObjectId string) {
  // TODO: implement function logic
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

// Show objects in room
func ShowObjsInRoom(pDnode *Dnode) {
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