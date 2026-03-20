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

var RoomFile    *os.File
var RoomScanner *bufio.Scanner

// Get RoomId
func GetRoomId(RoomId string) string {
  var RoomFileName string

  RoomFileName  = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  var err error
  RoomFile, err = os.Open(RoomFileName)
  if err != nil {
    // No such file???, But there should be, This is bad!
    LogIt("Room::GetRoomId - Room does not exist")
    os.Exit(1)
  }
  // RoomId
  RoomScanner = bufio.NewScanner(RoomFile)
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  if StrLeft(Stuff, 7) != "RoomId:" {
    // Very bad, where did the RoomId go anyway?
    LogIt("Room::GetRoomId - RoomId: not found")
    os.Exit(1)
  }
  RoomId = StrGetWord(Stuff, 2)
  RoomFile.Close()
  return RoomId
}

// Get RoomName
func GetRoomName(RoomId string) string {
  var RoomFileName string
  var RoomName     string

  RoomFileName  = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  var err error
  RoomFile, err = os.Open(RoomFileName)
  if err != nil {
    // No such file???, But there should be, This is bad!
    LogIt("Room::GetRoomName - Room does not exist")
    os.Exit(1) // _endthread()
  }
  // RoomName
  RoomScanner = bufio.NewScanner(RoomFile)

  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  if StrLeft(Stuff, 9) != "RoomName:" {
    // Very bad, where did the RoomName go anyway?
    LogIt("Room::GetRoomName - RoomName: not found")
    os.Exit(1)
  }
  RoomName = StrGetWords(Stuff, 2)
  RoomName = StrTrimLeft(RoomName)
  RoomName = StrTrimRight(RoomName)
  RoomFile.Close()
  return RoomName
}

// Get the list of exits that mobiles are allowed to use
func GetValidMobRoomExits(RoomId string) string {
  var ExitToRoomId  string
  var RoomFileName  string
  var ValidMobExits string

  RoomFileName  = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  var err error
  RoomFile, err = os.Open(RoomFileName)
  if err != nil {
    // No such file???, But there should be, This is bad!
    LogIt("Room::GetValidMobRoomExits - Room does not exist")
    os.Exit(1)
  }
  ValidMobExits = ""
  Stuff = "Not Done"
  RoomScanner = bufio.NewScanner(RoomFile)
  for Stuff != "End of Exits" {
    // Loop - process all exits
    RoomScanner.Scan()
    Stuff = RoomScanner.Text()
    if StrLeft(Stuff, 13) == "ExitToRoomId:" {
      // An Exit has been found
      ExitToRoomId = StrGetWord(Stuff, 2)
      if ExitToRoomId == "VineyardPath382" {
        //Success = 100;
        var x int
        x = 0
        _ = x
      }
      if !IsRoomType(ExitToRoomId, "NoNPC") {
        // And it's a valid Mob Exit
        ValidMobExits += ExitToRoomId
        ValidMobExits += " "
      }
    }
  }
  ValidMobExits = StrTrimRight(ValidMobExits)
  RoomFile.Close()
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
    // If the file isn't there, then the Room does not exit, doh!
    LogIt("Room::IsExit - Room does not exist")
    os.Exit(1)
  }
  ExitLookup = StrGetWord(CmdStr, 2)
  ExitLookup = StrMakeLower(ExitLookup)
  ExitLookup = TranslateWord(ExitLookup)
  Stuff = "Not Done"
  RoomScanner = bufio.NewScanner(RoomFile)
  for Stuff != "End of Exits" {
    // Loop until Exit is found or end of file
    if !RoomScanner.Scan() {
      break
    }
    Stuff = RoomScanner.Text()
    if StrLeft(Stuff, 9) == "ExitName:" {
      // Ok, an Exit has been found
      ExitName = StrGetWord(Stuff, 2)
      ExitName = StrMakeLower(ExitName)
      ExitName = TranslateWord(ExitName)
      if ExitName == ExitLookup {
        // THE Exit has been found
        Found = true
        Stuff = "End of Exits"
      }
    }
  }
  if Found {
    // At this point we know that the command entered referred to a valid exit
    if IsSleeping() {
      // Player is sleeping, send msg, command is not done
      CloseRoomFile()
      return true
    }
    if pDnodeActor.pPlayer.Position == "sit" {
      // The player is sitting, abort the move
      CloseRoomFile()
      pDnodeActor.PlayerOut += "You must be standing before you can move."
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return true
    }
    if MudCmdIsExit == "go" {
      // Go command
      if pDnodeActor.pPlayer.pPlayerFollowers[0] != nil {
        // If player is following another player
        if MudCmd != "flee" {
          // And not fleeing, abort the move
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
        // Position to ExitToRoomId line
        if !RoomScanner.Scan() {
          break
        }
        Stuff = RoomScanner.Text()
      }
      ExitToRoomId = StrGetWord(Stuff, 2)
      MovePlayer(pDnodeActor, ExitToRoomId)
      CloseRoomFile()
      ShowRoom(pDnodeActor)
      if PlayerRoomHasNotBeenHere(pDnodeActor.pPlayer) {
        // Player has not been here (on their own), give some experience
        pDnodeActor.PlayerOut += "\r\n"
        pDnodeActor.PlayerOut += "&YYou gain experience by exploring!&N"
        pDnodeActor.PlayerOut += "\r\n"
        CreatePrompt(pDnodeActor.pPlayer)
        pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
        GainExperience(pDnodeActor, 25)
      }
      if MudCmd != "flee" {
        // Player is not fleeing
        MoveFollowers(pDnodeActor, ExitToRoomId)
      }
    } else {
      // MudCmd was not 'go'
      if MudCmdIsExit == "look" {
        // Look at an exit
        ShowRoomExitDesc()
        CloseRoomFile()
      }
    }
    // so command processor will exit properly
    return true
  } else {
    // At this point we know that the command entered did NOT referred to a valid exit
    CloseRoomFile()
    // so command processor will keep trying
    return false
  }
}

// Is this a valid room?
func IsRoom(RoomId string) bool {
  var RoomFileName string

  RoomFileName  = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  var err error
  RoomFile, err = os.Open(RoomFileName)
  if err == nil {
    RoomScanner = bufio.NewScanner(RoomFile)
    RoomScanner.Scan()
    Stuff = RoomScanner.Text()
    RoomFile.Close()
    if StrLeft(Stuff, 7) != "RoomId:" {
      LogIt("Room::IsRoom - RoomId: not found")
      os.Exit(1)
    }
    Stuff = StrGetWord(Stuff, 2)
    if Stuff != RoomId {
      return false
    }
    return true
  } else {
    return false
  }
}

// Is the room type valid?
func IsRoomType(RoomId string, RoomType string) bool {
  var RoomFileName string

  RoomFileName  = ROOMS_DIR
  RoomFileName += RoomId
  RoomFileName += ".txt"
  var err error
  RoomFile, err = os.Open(RoomFileName)
  if err != nil {
    // No such file???, But there should be, This is bad!
    LogIt("Room::IsRoomType - Room does not exist")
    os.Exit(1) // _endthread()
  }
  // RoomType
  RoomScanner = bufio.NewScanner(RoomFile)
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  if StrLeft(Stuff, 9) != "RoomType:" {
    // Very bad, where did the RoomType go anyway?
    LogIt("Room::IsRoomType - RoomType: not found")
    os.Exit(1)
  }
  Stuff = StrGetWords(Stuff, 2)
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  RoomFile.Close()
  if StrIsNotWord(RoomType, Stuff) {
    // No matching RoomType found
    return false
  }
  // Matching RoomType found
  return true
}

// Show the room to the player
func ShowRoom(pDnode *Dnode) {
  if !OpenRoomFile(pDnode) {
    // If the file isn't there, then the Room does not exit, doh!
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
  // Recursive
  var pDnodeGrpMem *Dnode
  var i             int

  for i = 1; i < GRP_LIMIT; i++ {
    if pDnode.pPlayer.pPlayerFollowers[i] == nil {
      // No followers or no more followers
      return
    }
    pDnodeGrpMem = GetTargetDnode(pDnode.pPlayer.pPlayerFollowers[i].Name)
    if pDnodeGrpMem == nil {
      // Follower is not online and/or not in 'playing' state
      continue
    }
    if pDnode.pPlayer.RoomIdBeforeMove != pDnodeGrpMem.pPlayer.RoomId {
      // Not in same room, can't follow
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
  // Leaves message
  if MudCmd != "flee" {
    // If player is not fleeing
    MoveMsg = pDnode.PlayerName + " leaves."
    SendToRoom(pDnode.pPlayer.RoomId, MoveMsg)
  }
  // Switch rooms
  pDnode.pPlayer.RoomIdBeforeMove = pDnode.pPlayer.RoomId
  pDnode.pPlayer.RoomId = ExitToRoomId
  // Osi("Rooms", ExitToRoomId) TODO: Remove debug line
  PlayerSave(pDnode.pPlayer)
  // Arrives message
  MoveMsg = pDnode.PlayerName + " arrives."
  SendToRoom(pDnode.pPlayer.RoomId, MoveMsg)
}

// Open Room file
func OpenRoomFile(pDnode *Dnode) bool {
  var RoomFileName string

  RoomFileName  = ROOMS_DIR
  RoomFileName += pDnode.pPlayer.RoomId
  RoomFileName += ".txt"
  var err error
  RoomFile, err = os.Open(RoomFileName)
  if err == nil {
    RoomScanner = bufio.NewScanner(RoomFile)
    return true
  } else {
    return false
  }
}

// Show the room description to the player
func ShowRoomDesc(pDnode *Dnode) {
  // RoomDesc
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  if StrLeft(Stuff, 9) != "RoomDesc:" {
    LogIt("Room::ShowRoomDesc - RoomDesc: not found")
    os.Exit(1)
  }
  // Room Description
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  for Stuff != "End of RoomDesc" {
    pDnode.PlayerOut += Stuff
    pDnode.PlayerOut += "\r\n"
    RoomScanner.Scan()
    Stuff = RoomScanner.Text()
  }
}

// Show exit description to player
func ShowRoomExitDesc() {
  // ExitDesc
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  if StrLeft(Stuff, 9) != "ExitDesc:" {
    LogIt("Room::ShowRoomExitDesc - ExitDesc: not found")
    os.Exit(1)
  }
  // Exit Description
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  for StrLeft(Stuff, 13) != "ExitToRoomId:" {
    pDnodeActor.PlayerOut += Stuff
    pDnodeActor.PlayerOut += "\r\n"
    RoomScanner.Scan()
    Stuff = RoomScanner.Text()
  }
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Show the room exits to the player
func ShowRoomExits(pDnode *Dnode) {
  var NoExits bool

  NoExits = true
  pDnode.PlayerOut += "&C"
  pDnode.PlayerOut += "Exits:"
  Stuff = "Not Done"
  for Stuff != "End of Exits" {
    RoomScanner.Scan()
    Stuff = RoomScanner.Text()
    if StrLeft(Stuff, 9) == "ExitName:" {
      NoExits = false
      Stuff = StrGetWord(Stuff, 2)
      pDnode.PlayerOut += " "
      pDnode.PlayerOut += Stuff
    }
  }
  if NoExits {
    pDnode.PlayerOut += " None"
  }
  pDnode.PlayerOut += "&N"
}

// Show the room name to the player
func ShowRoomName(pDnode *Dnode) {
  var RoomId   string
  var RoomType string
  var Terrain  string
  var RoomName string

  // RoomId
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  if StrLeft(Stuff, 7) != "RoomId:" {
    LogIt("Room::ShowRoomName - RoomId: not found")
    os.Exit(1)
  }
  RoomId = StrGetWord(Stuff, 2)
  if RoomId != pDnode.pPlayer.RoomId {
    LogIt("Room::ShowRoomName - RoomId mis-match")
    os.Exit(1)
  }
  // RoomType
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  if StrLeft(Stuff, 9) != "RoomType:" {
    LogIt("Room::ShowRoomName - RoomType: not found")
    os.Exit(1)
  }
  RoomType = StrGetWords(Stuff, 2)
  // Terrain
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  if StrLeft(Stuff, 8) != "Terrain:" {
    LogIt("Room::ShowRoomName - Terrain: not found")
    os.Exit(1)
  }
  Terrain = StrGetWord(Stuff, 2)
  // RoomName
  RoomScanner.Scan()
  Stuff = RoomScanner.Text()
  if StrLeft(Stuff, 9) != "RoomName:" {
    LogIt("Room::ShowRoomName - RoomName: not found")
    os.Exit(1)
  }
  RoomName = StrGetWords(Stuff, 2)
  RoomName = StrTrimLeft(RoomName)
  // Build player output
  pDnode.PlayerOut += "\r\n"
  pDnode.PlayerOut += "&C"
  pDnode.PlayerOut += RoomName
  pDnode.PlayerOut += "&N"
  if pDnode.pPlayer.RoomInfo {
    // Show hidden room info
    pDnode.PlayerOut += "&M"
    pDnode.PlayerOut += " ["
    pDnode.PlayerOut += "&N"
    pDnode.PlayerOut += RoomId
    pDnode.PlayerOut += " "
    pDnode.PlayerOut += Terrain
    pDnode.PlayerOut += " "
    pDnode.PlayerOut += RoomType
    pDnode.PlayerOut += "&M"
    pDnode.PlayerOut += "]"
    pDnode.PlayerOut += "&N"
  }
  pDnode.PlayerOut += "\r\n"
}
