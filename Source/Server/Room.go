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

var RoomFile *os.File

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
    os.Exit(1)
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

// If valid room exit, then deal with it
func IsExit(MudCmdIsExit string) bool {
  var Found        bool
  var ExitLookup   string
  var ExitName     string
  var ExitToRoomId string

  Found = false
  if !OpenRoomFile(pDnodeActor) {
    LogIt("Room::IsExit - Room does not exist")
    os.Exit(1)
  }
  ExitLookup = StrGetWord(CmdStr, 2)
  ExitLookup = StrMakeLower(ExitLookup)
  ExitLookup = TranslateWord(ExitLookup)
  Stuff = "Not Done"
  Scanner := bufio.NewScanner(RoomFile)
  for Stuff != "End of Exits" {
    if !Scanner.Scan() {
      break
    }
    Stuff = Scanner.Text()
    if StrLeft(Stuff, 9) == "ExitName:" {
      ExitName = StrGetWord(Stuff, 2)
      ExitName = StrMakeLower(ExitName)
      ExitName = TranslateWord(ExitName)
      if ExitName == ExitLookup {
        Found = true
        Stuff = "End of Exits"
      }
    }
  }
  if Found {
    if IsSleeping() {
      CloseRoomFile()
      return true
    }
    if pDnodeActor.pPlayer.Position == "sit" {
      CloseRoomFile()
      pDnodeActor.PlayerOut += "You must be standing before you can move."
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return true
    }
    if MudCmdIsExit == "go" {
      if pDnodeActor.pPlayer.pPlayerFollowers[0] != nil {
        if MudCmd != "flee" {
          CloseRoomFile()
          pDnodeActor.PlayerOut += "Can't honor your command, you are following "
          pDnodeActor.PlayerOut += pDnodeActor.pPlayer.pPlayerFollowers[0].Name
          pDnodeActor.PlayerOut += ".\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          return true
        }
      }
      for StrLeft(Stuff, 13) != "ExitToRoomId:" {
        if !Scanner.Scan() {
          break
        }
        Stuff = Scanner.Text()
      }
      ExitToRoomId = StrGetWord(Stuff, 2)
      MovePlayer(pDnodeActor, ExitToRoomId)
      CloseRoomFile()
      ShowRoom(pDnodeActor)
      if PlayerRoomHasNotBeenHere(pDnodeActor.pPlayer) {
        pDnodeActor.PlayerOut += "\r\n"
        pDnodeActor.PlayerOut += "&YYou gain experience by exploring!&N"
        pDnodeActor.PlayerOut += "\r\n"
        CreatePrompt(pDnodeActor.pPlayer)
        pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
        GainExperience(pDnodeActor, 25)
      }
      if MudCmd != "flee" {
        MoveFollowers(pDnodeActor, ExitToRoomId)
      }
    } else {
      if MudCmdIsExit == "look" {
        ShowRoomExitDesc()
        CloseRoomFile()
      }
    }
    return true
  } else {
    CloseRoomFile()
    return false
  }
}

// Is this a valid room?
func IsRoom(RoomId string) bool {
  var RoomFileName string

  RoomFileName = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  f, err := os.Open(RoomFileName)
  if err != nil {
    return false
  }
  scanner := bufio.NewScanner(f)
  if !scanner.Scan() {
    f.Close()
    return false
  }
  Stuff = scanner.Text()
  f.Close()
  if StrLeft(Stuff, 7) != "RoomId:" {
    LogIt("Room::IsRoom - RoomId: not found")
    os.Exit(1)
  }
  Stuff = StrGetWord(Stuff, 2)
  if Stuff != RoomId {
    return false
  }
  return true
}

// Is the room type valid?
func IsRoomType(RoomId string, RoomType string) bool {
  var RoomFileName string

  RoomFileName = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  f, err := os.Open(RoomFileName)
  if err != nil {
    LogIt("Room::IsRoomType - Room does not exist")
    os.Exit(1) // _endthread()
  }
  scanner := bufio.NewScanner(f)
  if !scanner.Scan() {
    f.Close()
    LogIt("Room::IsRoomType - RoomType: not found")
    os.Exit(1)
  }
  Stuff = scanner.Text()
  if !scanner.Scan() {
    f.Close()
    LogIt("Room::IsRoomType - RoomType: not found")
    os.Exit(1)
  }
  Stuff = scanner.Text()
  if StrLeft(Stuff, 9) != "RoomType:" {
    f.Close()
    LogIt("Room::IsRoomType - RoomType: not found")
    os.Exit(1)
  }
  Stuff = StrGetWords(Stuff, 2)
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  f.Close()
  if StrIsNotWord(RoomType, Stuff) {
    return false
  }
  return true
}

// Show the room to the player
func ShowRoom(pDnode *Dnode) {
  if !OpenRoomFile(pDnode) {
    LogIt("Room::ShowRoom - Room does not exist")
    os.Exit(1) // _endthread()
  }
  ShowRoomName(pDnode)
  ShowRoomDesc(pDnode)
  ShowRoomExits(pDnode)
  CloseRoomFile()
  ShowPlayersInRoom(pDnode)
  ShowObjsInRoom(pDnode)
  ShowMobsInRoom(pDnode)
  pDnode.PlayerOut += "\r\n"
  CreatePrompt(pDnode.pPlayer)
  pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
}

// Close Room file
func CloseRoomFile() {
  if RoomFile != nil {
    RoomFile.Close()
  }
}

// Move followers
func MoveFollowers(pDnode *Dnode, ExitToRoomId string) {
  var pDnodeGrpMem *Dnode
  var i int

  for i = 1; i < GRP_LIMIT; i++ {
    if pDnode.pPlayer.pPlayerFollowers[i] == nil {
      return
    }
    pDnodeGrpMem = GetTargetDnode(pDnode.pPlayer.pPlayerFollowers[i].Name)
    if pDnodeGrpMem == nil {
      continue
    }
    if pDnode.pPlayer.RoomIdBeforeMove != pDnodeGrpMem.pPlayer.RoomId {
      continue
    }
    MovePlayer(pDnodeGrpMem, ExitToRoomId)
    ShowRoom(pDnodeGrpMem)
    MoveFollowers(pDnodeGrpMem, ExitToRoomId)
  }
}

// Go command - move the player
func MovePlayer(pDnode *Dnode, ExitToRoomId string) {
  var MoveMsg string

  pDnodeSrc = pDnode
  pDnodeTgt = pDnode
  if MudCmd != "flee" {
    MoveMsg = pDnode.PlayerName + " leaves."
    SendToRoom(pDnode.pPlayer.RoomId, MoveMsg)
  }
  pDnode.pPlayer.RoomIdBeforeMove = pDnode.pPlayer.RoomId
  pDnode.pPlayer.RoomId = ExitToRoomId
  // Osi("Rooms", ExitToRoomId) TODO: Remove debug line
  PlayerSave(pDnode.pPlayer)
  MoveMsg = pDnode.PlayerName + " arrives."
  SendToRoom(pDnode.pPlayer.RoomId, MoveMsg)
}

// Open Room file
func OpenRoomFile(pDnode *Dnode) bool {
  var RoomFileName string

  RoomFileName = ROOMS_DIR
  RoomFileName += pDnode.pPlayer.RoomId
  RoomFileName += ".txt"
  f, err := os.Open(RoomFileName)
  if err != nil {
    return false
  }
  RoomFile = f
  return true
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