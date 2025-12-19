//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Room.go                                               *
// Usage:     Manages game rooms and their interactions             *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
  "bufio"
  "os"
)

// Get RoomId
func GetRoomId(RoomId string) string {
  var RoomFileName string

  RoomFileName = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  f, err := os.Open(RoomFileName)
  if err != nil {
    LogIt("Room::GetRoomId - Room does not exist")
    os.Exit(1) // _endthread()
  }
  defer f.Close()

  scanner := bufio.NewScanner(f)
  if !scanner.Scan() {
    LogIt("Room::GetRoomId - RoomId: not found")
    os.Exit(1)
  }
  Stuff = scanner.Text()
  if StrLeft(Stuff, 7) != "RoomId:" {
    LogIt("Room::GetRoomId - RoomId: not found")
    os.Exit(1)
  }
  RoomId = StrGetWord(Stuff, 2)
  return RoomId
}

func GetRoomName(RoomId string) string {
  return ""
}

func GetValidMobRoomExits(RoomId string) string {
  return ""
}

func IsExit(MudCmdIsExit string) bool {
  return false
}

func IsRoom(RoomId string) bool {
  return false
}

func IsRoomType(RoomId string, RoomType string) bool {
  return false
}

func ShowRoom(pDnode *Dnode) {
  // TODO: Implement ShowRoom logic
}

func CloseFile() {
  // TODO: Implement CloseFile logic
}

func MoveFollowers(pDnode *Dnode, ExitToRoomId string) {
  // TODO: Implement MoveFollowers logic
}

func MovePlayer(pDnode *Dnode, ExitToRoomId string) {
  // TODO: Implement MovePlayer logic
}

func OpenFile(pDnode *Dnode) bool {
  return false
}

func ShowRoomDesc(pDnode *Dnode) {
  // TODO: Implement ShowRoomDesc logic
}

func ShowRoomExitDesc() {
  // TODO: Implement ShowRoomExitDesc logic
}

func ShowRoomExits(pDnode *Dnode) {
  // TODO: Implement ShowRoomExits logic
}

func ShowRoomName(pDnode *Dnode) {
  // TODO: Implement ShowRoomName logic
}