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

// Get RoomName
func GetRoomName(RoomId string) string {
  var RoomFileName string
  var RoomName     string

  RoomFileName = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  f, err := os.Open(RoomFileName)
  if err != nil {
    LogIt("Room::GetRoomName - Room does not exist")
    os.Exit(1) // _endthread()
  }

  scanner := bufio.NewScanner(f)
  if !scanner.Scan() {
    LogIt("Room::GetRoomName - RoomName: not found")
    os.Exit(1)
  }
  Stuff = scanner.Text()
  if !scanner.Scan() {
    LogIt("Room::GetRoomName - RoomName: not found")
    os.Exit(1)
  }
  Stuff = scanner.Text()
  if !scanner.Scan() {
    LogIt("Room::GetRoomName - RoomName: not found")
    os.Exit(1)
  }
  Stuff = scanner.Text()
  if !scanner.Scan() {
    LogIt("Room::GetRoomName - RoomName: not found")
    os.Exit(1)
  }
  Stuff = scanner.Text()
  if StrLeft(Stuff, 9) != "RoomName:" {
    LogIt("Room::GetRoomName - RoomName: not found")
    os.Exit(1)
  }
  RoomName = StrGetWords(Stuff, 2)
  RoomName = StrTrimLeft(RoomName)
  RoomName = StrTrimRight(RoomName)
  f.Close()
  return RoomName
}

// Get the list of exits that mobiles are allowed to use
func GetValidMobRoomExits(RoomId string) string {
  var ExitToRoomId  string
  var RoomFileName  string
  var ValidMobExits string

  RoomFileName = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  f, err := os.Open(RoomFileName)
  if err != nil {
    LogIt("Room::GetValidMobRoomExits - Room does not exist")
    os.Exit(1) // _endthread()
  }

  ValidMobExits = ""
  Stuff = "Not Done"
  Scanner := bufio.NewScanner(f)
  for Stuff != "End of Exits" {
    if !Scanner.Scan() {
      break
    }
    Stuff = Scanner.Text()
    if StrLeft(Stuff, 13) == "ExitToRoomId:" {
      ExitToRoomId = StrGetWord(Stuff, 2)
      if ExitToRoomId == "VineyardPath382" {
        //Success = 100;
        _ = 0
      }
      if !IsRoomType(ExitToRoomId, "NoNPC") {
        ValidMobExits += ExitToRoomId
        ValidMobExits += " "
      }
    }
  }
  ValidMobExits = StrTrimRight(ValidMobExits)
  f.Close()
  return ValidMobExits
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