//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Communication.go                                      *
// Usage:     Handles communication between server and clients      *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
  "bufio"
  "fmt"
  "os"
  "syscall"
  "time"
)

// Globals
var pDnodeOthers *Dnode
var ListenSocket  syscall.Handle
var ValidCmds     []string
var SockAddr      syscall.Sockaddr

const WSAEWOULDBLOCK syscall.Errno = 10035
const WSAEINTR       syscall.Errno = 10004

// Return pointer of target, if target in 'playing' state
func GetTargetDnode(TargetName string) *Dnode {
  var pDnodeLookup *Dnode
  var TargetFound   bool
  var LookupName    string

  TargetFound = false
  TargetName = StrMakeLower(TargetName)
  SetpDnodeCursorFirst()
  for !EndOfDnodeList() {
    pDnodeLookup = GetDnode()
    LookupName = StrMakeLower(pDnodeLookup.PlayerName)
    if TargetName == LookupName {
      if pDnodeLookup.PlayerStatePlaying {
        TargetFound = true
        break
      }
    }
    SetpDnodeCursorNext()
  }
  RepositionDnodeCursor()
  if TargetFound {
    return pDnodeLookup
  }
  return nil
}

// Check to see if player is fighting
func IsFighting() bool {
  var RandomNumber int
  var FightingMsg string

  if !pDnodeActor.PlayerStateFighting {
    return false
  }
  RandomNumber = GetRandomNumber(5)
  switch RandomNumber {
  case 1:
    FightingMsg = "You are fighting for your life!"
  case 2:
    FightingMsg = "Not now, your are fighting."
  case 3:
    FightingMsg = "No can do, you are fighting"
  case 4:
    FightingMsg = "You are busy swinging a weapon."
  case 5:
    FightingMsg = "NO!, now get back in the fight!."
  default:
    FightingMsg = "You are fighting."
  }
  pDnodeActor.PlayerOut += FightingMsg
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  return true
}

// Check to see if player is sleeping
func IsSleeping() bool {
  var RandomNumber int
  var SleepingMsg string

  if pDnodeActor.pPlayer.Position != "sleep" {
    return false
  }
  RandomNumber = GetRandomNumber(5)
  switch RandomNumber {
  case 1:
    SleepingMsg = "You must be dreaming."
  case 2:
    SleepingMsg = "You dream about doing something."
  case 3:
    SleepingMsg = "It's such a nice dream, please don't wake me."
  case 4:
    SleepingMsg = "Your snoring almost wakes you up."
  case 5:
    SleepingMsg = "Dream, dream, dreeeeaaaammmm, all I do is dream."
  default:
    SleepingMsg = "You must be dreaming."
  }
  pDnodeActor.PlayerOut += SleepingMsg
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  return true
}

// Send output to all players
func SendToAll(PlayerMsg, AllMsg string) {
  SetpDnodeCursorFirst()
  for !EndOfDnodeList() {
    pDnodeOthers = GetDnode()
    if pDnodeActor == pDnodeOthers {
      pDnodeActor.PlayerOut += PlayerMsg
    } else {
      if pDnodeOthers.PlayerStatePlaying {
        if pDnodeActor.PlayerStateInvisible {
          SetpDnodeCursorNext()
          continue
        }
        pDnodeOthers.PlayerOut += AllMsg
        CreatePrompt(pDnodeOthers.pPlayer)
        pDnodeOthers.PlayerOut += GetOutput(pDnodeOthers.pPlayer)
      }
    }
    SetpDnodeCursorNext()
  }
  RepositionDnodeCursor()
}

// Send output to other players in the same room as player
func SendToRoom(TargetRoomId, MsgText string) {
  var LookupRoomId string

  SetpDnodeCursorFirst()
  for !EndOfDnodeList() {
    pDnodeOthers = GetDnode()
    if pDnodeOthers.PlayerStatePlaying {
      LookupRoomId = pDnodeOthers.pPlayer.RoomId
      if pDnodeSrc != pDnodeOthers {
        if pDnodeTgt != pDnodeOthers {
          if TargetRoomId == LookupRoomId {
            if pDnodeOthers.pPlayer.Position != "sleep" {
              if pDnodeSrc != nil {
                if pDnodeSrc.PlayerStateInvisible {
                  SetpDnodeCursorNext()
                  continue
                }
              }
              if pDnodeOthers.PlayerStateInvisible {
                SetpDnodeCursorNext()
                continue
              }
              pDnodeOthers.PlayerOut += "\r\n"
              pDnodeOthers.PlayerOut += MsgText
              pDnodeOthers.PlayerOut += "&N"
              pDnodeOthers.PlayerOut += "\r\n"
              CreatePrompt(pDnodeOthers.pPlayer)
              pDnodeOthers.PlayerOut += GetOutput(pDnodeOthers.pPlayer)
            }
          }
        }
      }
    }
    SetpDnodeCursorNext()
  }
  RepositionDnodeCursor()
}

// Show players in a given room
func ShowPlayersInRoom(pDnode *Dnode) {
  var LookupRoomId string
  var TargetRoomId string

  TargetRoomId = pDnode.pPlayer.RoomId
  SetpDnodeCursorFirst()
  for !EndOfDnodeList() {
    pDnodeOthers = GetDnode()
    if pDnode != pDnodeOthers {
      if pDnodeOthers.PlayerStatePlaying {
        LookupRoomId = pDnodeOthers.pPlayer.RoomId
        if TargetRoomId == LookupRoomId {
          if pDnodeOthers.PlayerStateInvisible {
            SetpDnodeCursorNext()
            continue
          }
          pDnode.PlayerOut += "\r\n"
          pDnode.PlayerOut += "&W"
          pDnode.PlayerOut += pDnodeOthers.PlayerName
          pDnode.PlayerOut += " is "
          if pDnodeOthers.PlayerStateFighting {
            pDnode.PlayerOut += "engaged in a fight!"
          } else if pDnodeOthers.pPlayer.Position == "sleep" {
            pDnode.PlayerOut += "here, sound asleep."
          } else {
            pDnode.PlayerOut += pDnodeOthers.pPlayer.Position
            if pDnodeOthers.pPlayer.Position == "sit" {
              pDnode.PlayerOut += "t"
            }
            pDnode.PlayerOut += "ing here."
          }
          pDnode.PlayerOut += "&N"
        }
      }
    }
    SetpDnodeCursorNext()
  }
  RepositionDnodeCursor()
}

// Check for new connections
func SockCheckForNewConnections() {
  DEBUGIT(5)
  SockNewConnection()
}

// Close port
func SockClosePort(Port int) {
  DEBUGIT(1)
  Result := CloseSocket(ListenSocket)
  if Result != 0 {
    LogBuf = "Communication::~Communication - Error: closesocket"
    LogIt(LogBuf)
    // no _endthread in Go; exiting the goroutine is implicit
    return
  }
  Buf = fmt.Sprintf("%d", Port)
  LogBuf = "Closed port " + Buf
  LogIt(LogBuf)
}

// Open port
func SockOpenPort(Port int) {
  DEBUGIT(1)
  fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
  if err != nil {
    Buf = err.Error()
    LogBuf = "Communication::SockOpenPort - Error: initializing socket: " + Buf
    LogIt(LogBuf)
    PrintIt("Communication::SockOpenPort - Error: initializing socket")
    os.Exit(1)
  }
  ListenSocket = fd
  err = syscall.SetsockoptInt(ListenSocket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
  if err != nil {
    Buf = err.Error()
    LogBuf = "Communication::SockOpenPort - Error: setsockopt: SO_REUSEADDR: " + Buf
    LogIt(LogBuf)
    PrintIt("Communication::SockOpenPort - Error: setsockopt: SO_REUSEADDR")
    os.Exit(1)
  }
  err = syscall.SetsockoptInt(ListenSocket, syscall.SOL_SOCKET, syscall.SO_SNDBUF, 1)
  if err != nil {
    Buf = err.Error()
    LogBuf = "Communication::SockOpenPort - Error: setsockopt: SO_SNDBUF: " + Buf
    LogIt(LogBuf)
    PrintIt("Communication::SockOpenPort - Error: setsockopt: SO_SNDBUF")
    os.Exit(1)
  }
  linger := &syscall.Linger{Onoff: 0, Linger: 0}
  err = syscall.SetsockoptLinger(ListenSocket, syscall.SOL_SOCKET, syscall.SO_LINGER, linger)
  if err != nil {
    Buf = err.Error()
    LogBuf = "Communication::SockOpenPort - Error: setsockopt: SO_LINGER: " + Buf
    LogIt(LogBuf)
    PrintIt("Communication::SockOpenPort - Error: setsockopt: SO_LINGER")
    os.Exit(1)
  }
  sa := &syscall.SockaddrInet4{Port: Port}
  err = syscall.Bind(ListenSocket, sa)
  if err != nil {
    Buf = err.Error()
    LogBuf = "Communication::SockOpenPort - Error: bind: " + Buf
    LogIt(LogBuf)
    PrintIt("Communication::SockOpenPort - Error: bind")
    os.Exit(1)
  }
  err = syscall.SetNonblock(ListenSocket, true)
  if err != nil {
    Buf = err.Error()
    LogBuf = "Communication::SockOpenPort - Error: SetNonblock: " + Buf
    LogIt(LogBuf)
    PrintIt("Communication::SockOpenPort - Error: SetNonblock")
    os.Exit(1)
  }
  err = syscall.Listen(ListenSocket, 20)
  if err != nil {
    Buf = err.Error()
    LogBuf = "Communication::SockOpenPort - Error: listen: " + Buf
    LogIt(LogBuf)
    CloseSocket(ListenSocket)
    PrintIt("Communication::SockOpenPort - Error: listen")
    os.Exit(1)
  }
  Buf = fmt.Sprintf("%d", Port)
  LogBuf = "Listening on port " + Buf
  LogIt(LogBuf)
  CommandArrayLoad()
}

// Receive player input, check player status, send output
func SockRecv() {
  var ConnectionCount int
  var DnodeFdSave int

  DEBUGIT(5)
  SetpDnodeCursorFirst()
  for !EndOfDnodeList() {
    pDnodeActor = GetDnode()
    // Check connection status (logon timeout)
    if !pDnodeActor.PlayerStatePlaying {
      pDnodeActor.InputTick++
      if pDnodeActor.InputTick >= INPUT_TICK {
        pDnodeActor.PlayerStateBye = true
        Buf = fmt.Sprintf("%d", pDnodeActor.DnodeFd)
        LogBuf = "Time out during logon on descriptor " + Buf
        LogIt(LogBuf)
        pDnodeActor.PlayerOut += "\r\n"
        pDnodeActor.PlayerOut += "No input ... closing connection"
      }
    }
    // Receive / exception handling via non-blocking recv
    buf := make([]byte, MAX_INPUT_LENGTH)
    n, _, err := syscall.Recvfrom(pDnodeActor.DnodeFd, buf, 0)
    if err != nil {
      if err != WSAEWOULDBLOCK && err != WSAEINTR {
        // Exception: kick out connection
        pDnodeActor.PlayerStateBye = true
        if pDnodeActor.PlayerStatePlaying {
          pDnodeActor.PlayerStatePlaying = false
          if pDnodeActor.pPlayer != nil {
            PlayerSave(pDnodeActor.pPlayer)
          }
        }
      }
    } else {
      if n == 0 {
        // Disconnected
        pDnodeActor.PlayerStateBye = true
        if pDnodeActor.PlayerStatePlaying {
          pDnodeActor.PlayerStatePlaying = false
          if pDnodeActor.pPlayer != nil {
            PlayerSave(pDnodeActor.pPlayer)
          }
        }
      } else if n > 0 {
        // Got input
        pDnodeActor.PlayerInp += string(buf[:n])
        pDnodeActor.InputTick = 0
      }
    }
    // Banner for new connection
    if pDnodeActor.PlayerStateSendBanner {
      pDnodeActor.PlayerStateSendBanner = false
      pDnodeActor.PlayerStateLoggingOn = true
      pDnodeActor.PlayerStateWaitNewCharacter = true
      LogonGreeting()
      pDnodeActor.PlayerOut += "\r\n"
      pDnodeActor.PlayerOut += "Create a new character Y-N?"
      pDnodeActor.PlayerOut += "\r\n"
    }
    // Update player stats
    if pDnodeActor.PlayerStatePlaying {
      pDnodeActor.StatsTick++
      if pDnodeActor.StatsTick >= STATS_TICK {
        pDnodeActor.StatsTick = 0
        UpdatePlayerStats()
      }
      // Hunger & thirst
      pDnodeActor.HungerThirstTick++
      if pDnodeActor.HungerThirstTick >= HUNGER_THIRST_TICK {
        pDnodeActor.HungerThirstTick = 0
        pDnodeActor.pPlayer.Hunger++
        pDnodeActor.pPlayer.Thirst++
        if pDnodeActor.pPlayer.Level < HUNGER_THIRST_LEVEL {
          pDnodeActor.pPlayer.Hunger = 0
          pDnodeActor.pPlayer.Thirst = 0
        }
        if pDnodeActor.pPlayer.Admin {
          pDnodeActor.pPlayer.Hunger = 0
          pDnodeActor.pPlayer.Thirst = 0
        }
        if pDnodeActor.pPlayer.Hunger > 99 {
          pDnodeActor.pPlayer.Hunger = 100
          pDnodeActor.PlayerOut += "\r\n"
          pDnodeActor.PlayerOut += "You are extremely hungry!!!"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
        }
        if pDnodeActor.pPlayer.Thirst > 99 {
          pDnodeActor.pPlayer.Thirst = 100
          pDnodeActor.PlayerOut += "\r\n"
          pDnodeActor.PlayerOut += "You are extremely thirsty!!!"
          pDnodeActor.PlayerOut += "\r\n"

          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
        }
      }
    }
    // Handle fights
    if pDnodeActor.PlayerStateFighting {
      pDnodeActor.FightTick++
      if pDnodeActor.FightTick >= FIGHT_TICK {
        pDnodeActor.FightTick = 0
        Violence()
      }
    }
    // Is game stopping?
    if StateStopping {
      SetpDnodeCursorFirst()
      for !EndOfDnodeList() {
        pDnodeOthers = GetDnode()
        if !pDnodeOthers.PlayerStateBye {
          pDnodeOthers.PlayerStateBye = true
          pDnodeOthers.PlayerStatePlaying = false
          if pDnodeOthers.PlayerStatePlaying {
            PlayerSave(pDnodeOthers.pPlayer)
          }
          pDnodeOthers.PlayerOut += "\r\n"
          pDnodeOthers.PlayerOut += "Game is stopping ... Bye Bye!"
          pDnodeOthers.PlayerOut += "\r\n"
          LogBuf = pDnodeOthers.PlayerName
          LogBuf += " will be force disconnected"
          LogIt(LogBuf)
        }
        SetpDnodeCursorNext()
      }
      RepositionDnodeCursor()
    }
    // Send player output
    if StrGetLength(pDnodeActor.PlayerOut) > 0 {
      Color()
      SockSend(pDnodeActor.PlayerOut)
    }
    // Is player quitting?
    if pDnodeActor.PlayerStateBye {
      if !pDnodeActor.PlayerStateReconnecting {
        if pDnodeActor.pPlayer != nil {
          PlayerDestructor(pDnodeActor.pPlayer)
        }
      }
      DnodeFdSave = int(pDnodeActor.DnodeFd)
      if DeleteNode() {
        Buf = fmt.Sprintf("%d", DnodeFdSave)
        LogBuf = "Closed connection on descriptor " + Buf
        LogIt(LogBuf)
        ConnectionCount = GetDnodeCount()
        if ConnectionCount == 1 {
          if StateStopping {
            StateRunning = false
          }
        }
      }
      SetpDnodeCursorNext()
      continue
    }
    // Process player input
    if !StateStopping {
      LineFeedPosition := StrFindOneOf(pDnodeActor.PlayerInp, "\r\n")
      if LineFeedPosition > -1 {
        if pDnodeActor.PlayerName != "Ixaka" && pDnodeActor.PlayerName != "Kwam" {
          LogBuf = pDnodeActor.PlayerIpAddress
          LogBuf += " "
          if pDnodeActor.pPlayer != nil {
            LogBuf += pDnodeActor.pPlayer.RoomId
            LogBuf += " "
          }
          LogBuf += pDnodeActor.PlayerInp
          StrReplace(&LogBuf, "\r", " ")
          StrReplace(&LogBuf, "\n", " ")
          LogIt(LogBuf)
        }
        CommandParse()
      }
    }
    SetpDnodeCursorNext()
  }
}

// Replace or strip out color codes
func Color() {
  if !pDnodeActor.PlayerStatePlaying {
    return
  }
  sPlayerOut := pDnodeActor.PlayerOut
  if pDnodeActor.pPlayer.Color {
    StrReplace(&sPlayerOut, "&N", Normal)
    StrReplace(&sPlayerOut, "&K", BrightBlack)
    StrReplace(&sPlayerOut, "&R", BrightRed)
    StrReplace(&sPlayerOut, "&G", BrightGreen)
    StrReplace(&sPlayerOut, "&Y", BrightYellow)
    StrReplace(&sPlayerOut, "&B", BrightBlue)
    StrReplace(&sPlayerOut, "&M", BrightMagenta)
    StrReplace(&sPlayerOut, "&C", BrightCyan)
    StrReplace(&sPlayerOut, "&W", BrightWhite)
  } else {
    StrReplace(&sPlayerOut, "&N", "")
    StrReplace(&sPlayerOut, "&K", "")
    StrReplace(&sPlayerOut, "&R", "")
    StrReplace(&sPlayerOut, "&G", "")
    StrReplace(&sPlayerOut, "&Y", "")
    StrReplace(&sPlayerOut, "&B", "")
    StrReplace(&sPlayerOut, "&M", "")
    StrReplace(&sPlayerOut, "&C", "")
    StrReplace(&sPlayerOut, "&W", "")
  }
  pDnodeActor.PlayerOut = sPlayerOut
}

// Load command array
func CommandArrayLoad() {
  ValidCmdsFileName := VALID_CMDS_DIR + "ValidCommands.txt"
  file, err := os.Open(ValidCmdsFileName)
  if err != nil {
    LogBuf = "Communication::CommandArrayLoad - Open Valid Commands file failed (read)"
    LogIt(LogBuf)
    return
  }
  defer file.Close()
  ValidCmds = ValidCmds[:0]
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    ValidCmds = append(ValidCmds, line)
  }
  LogBuf = "Command array loaded"
  LogIt(LogBuf)
}

// Check command authorization, level, and validity
func CommandCheck(MudCmdChk string) string {
  var CommandCheckResult string
  var ValCmd string
  var ValCmdInfo string
  var WhoCanDo string

  CommandCheckResult = "Not Found"
  for _, ValidCmd := range ValidCmds {
    ValCmdInfo = ValidCmd
    ValCmd = StrGetWord(ValCmdInfo, 1)
    WhoCanDo = StrGetWord(ValCmdInfo, 2)
    if MudCmdChk == ValCmd {
      if WhoCanDo == "all" {
        CommandCheckResult = "Ok"
        break
      } else if WhoCanDo == "admin" {
        if pDnodeActor.pPlayer.Admin {
          CommandCheckResult = "Ok"
          break
        } else {
          CommandCheckResult = "NotOk"
        }
      } else if StrToInt(WhoCanDo) > pDnodeActor.pPlayer.Level {
        CommandCheckResult = "Level " + WhoCanDo
        break
      } else {
        CommandCheckResult = "Ok"
        break
      }
    }
  }
  if CommandCheckResult == "" {
    LogBuf = "Communication::CommandCheck - Broke!"
    LogIt(LogBuf)
    CommandCheckResult = "Not Found"
  }
  return CommandCheckResult
}

// Command parsing
func CommandParse() {
  var BadCommandMsg string
  var CmdStrLength int
  var CommandCheckResult string
  var MudCmdChk string
  var MudCmdOk bool
  var PositionOfNewline int
  var RandomNumber int

  //**************************
  // Get next command string *
  //**************************
  CmdStr = pDnodeActor.PlayerInp
  CmdStrLength = StrGetLength(CmdStr)
  PositionOfNewline = StrFindOneOf(CmdStr, "\r\n")
  if PositionOfNewline < 0 {
    // No newline found, skip out
    return
  }
  CmdStr = StrLeft(CmdStr, PositionOfNewline)
  pDnodeActor.PlayerInp = StrRight(pDnodeActor.PlayerInp, CmdStrLength-PositionOfNewline)
  pDnodeActor.PlayerInp = StrTrimLeft(pDnodeActor.PlayerInp)
  if CmdStr == "" {
    // Player hit enter without typing anything
    if !pDnodeActor.PlayerStateLoggingOn {
      // Player is not logging on
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  //***************
  // Player logon *
  //***************
  if pDnodeActor.PlayerStateLoggingOn {
    // Player just connected and needs to logon
    DoLogon()
    return
  }
  //*************
  // Get MudCmd *
  //*************
  MudCmd = StrGetWord(CmdStr, 1)
  MudCmd = StrMakeLower(MudCmd)
  // Translate 'n' into 'go north'
  MudCmd = TranslateWord(MudCmd)
  if StrCountWords(MudCmd) == 2 {
    // Re-get MudCmd. In the case of 'go north', MudCmd is 'go'
    CmdStr = MudCmd
    MudCmd = StrGetWord(CmdStr, 1)
    MudCmd = StrMakeLower(MudCmd)
  }
  // Check for spamming
  if MudCmd != "go" {
    // 'go' command is ok
    pDnodeActor.CmdName3 = pDnodeActor.CmdName2
    pDnodeActor.CmdName2 = pDnodeActor.CmdName1
    pDnodeActor.CmdName1 = MudCmd
    pDnodeActor.CmdTime3 = pDnodeActor.CmdTime2
    pDnodeActor.CmdTime2 = pDnodeActor.CmdTime1
    pDnodeActor.CmdTime1 = clock()
    if pDnodeActor.CmdName1 == pDnodeActor.CmdName2 {
      // Command same as last command
      if pDnodeActor.CmdName1 == pDnodeActor.CmdName3 {
        // Command same as last two commands
        if pDnodeActor.CmdTime1-pDnodeActor.CmdTime3 < 1000 {
          // Stop spamming
          pDnodeActor.PlayerOut += "&R"
          pDnodeActor.PlayerOut += "NO SPAMMING!!"
          pDnodeActor.PlayerOut += "&N"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          pDnodeActor.PlayerInp = ""
          return
        }
      }
    }
  }
  //****************
  // Check command *
  //****************
  MudCmdOk = false
  MudCmdChk = MudCmd
  CommandCheckResult = CommandCheck(MudCmdChk)
  if CommandCheckResult == "Ok" {
    // Mud command is Ok for this player
    MudCmdOk = true
  } else if StrGetWord(CommandCheckResult, 1) == "Level" {
    // Level restriction on command
    pDnodeActor.PlayerOut += "You must attain level "
    pDnodeActor.PlayerOut += StrGetWord(CommandCheckResult, 2)
    pDnodeActor.PlayerOut += " before you can use that command."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //******************
  //* SOCIAL command *
  //******************
  if !MudCmdOk {
    // Not validated yet, maybe cmd is a social
    if IsSocial() {
      // Yep, it was a social
      return
    }
  }
  //**************
  // Bad command *
  //**************
  if !MudCmdOk {
    // Not a valid cmd and it is not a social
    RandomNumber = GetRandomNumber(5)
    switch RandomNumber {
    case 1:
      BadCommandMsg = "How's that?"
    case 2:
      BadCommandMsg = "You try to give a command, but fail."
    case 3:
      BadCommandMsg = "Hmmm, making up commands?"
    case 4:
      BadCommandMsg = "Ehh, what's that again?"
    case 5:
      BadCommandMsg = "Feeling creative?"
    default:
      BadCommandMsg = "Your command is not clear."
    }
    pDnodeActor.PlayerOut += BadCommandMsg
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //**********************
  //* Process the MudCmd *
  //**********************
  // ADVANCE command
  if MudCmd == "advance" {
    DoAdvance()
    return
  }
  // AFK command
  if MudCmd == "afk" {
    DoAfk()
    return
  }
  // ASSIST command
  if MudCmd == "assist" {
    DoAssist()
    return
  }
  // BUY command
  if MudCmd == "buy" {
    DoBuy()
    return
  }
  // CHAT command
  if MudCmd == "chat" {
    DoChat()
    return
  }
  // COLOR command
  if MudCmd == "color" {
    DoColor()
    return
  }
  // CONSIDER command
  if MudCmd == "consider" {
    DoConsider()
    return
  }
  // DELETE command
  if MudCmd == "delete" {
    DoDelete()
    return
  }
  // DESTROY command
  if MudCmd == "destroy" {
    DoDestroy()
    return
  }
  // DRINK command
  if MudCmd == "drink" {
    DoDrink()
    return
  }
  // DROP command
  if MudCmd == "drop" {
    DoDrop()
    return
  }
  // EAT command
  if MudCmd == "eat" {
    DoEat()
    return
  }
  // EMOTE command
  if MudCmd == "emote" {
    DoEmote()
    return
  }
  // EQUIPMENT command
  if MudCmd == "equipment" {
    DoEquipment()
    return
  }
  // EXAMINE command
  if MudCmd == "examine" {
    DoExamine()
    return
  }
  // FLEE command
  if MudCmd == "flee" {
    DoFlee()
    return
  }
  // FOLLOW command
  if MudCmd == "follow" {
    DoFollow(pDnodeActor, CmdStr)
    return
  }
  // GET command
  if MudCmd == "get" {
    DoGet()
    return
  }
  // GIVE command
  if MudCmd == "give" {
    DoGive()
    return
  }
  // GO command
  if MudCmd == "go" {
    DoGo()
    return
  }
  // GOTOARRIVE command
  if MudCmd == "gotoarrive" {
    DoGoToArrive()
    return
  }
  // GOTODEPART command
  if MudCmd == "gotodepart" {
    DoGoToDepart()
    return
  }
  // GOTO command
  if MudCmd == "goto" {
    DoGoTo()
    return
  }
  // GROUP command
  if MudCmd == "group" {
    DoGroup()
    return
  }
  // GSAY command
  if MudCmd == "gsay" {
    DoGsay()
    return
  }
  // HAIL command
  if MudCmd == "hail" {
    DoHail()
    return
  }
  // HELP command
  if MudCmd == "help" {
    DoHelp()
    return
  }
  // INVENTORY command
  if MudCmd == "inventory" {
    DoInventory()
    return
  }
  // INVISIBLE command
  if MudCmd == "invisible" {
    DoInvisible()
    return
  }
  // KILL command
  if MudCmd == "kill" {
    DoKill()
    return
  }
  // LIST command
  if MudCmd == "list" {
    DoList()
    return
  }
  // LOAD command
  if MudCmd == "load" {
    DoLoad()
    return
  }
  // LOOK command
  if MudCmd == "look" {
    DoLook(CmdStr)
    return
  }
  // MONEY command
  if MudCmd == "money" {
    DoMoney()
    return
  }
  // MOTD command
  if MudCmd == "motd" {
    DoMotd()
    return
  }
  // OneWhack command
  if MudCmd == "onewhack" {
    DoOneWhack()
    return
  }
  // PASSWORD command
  if MudCmd == "password" {
    DoPassword()
    return
  }
  // PLAYED command
  if MudCmd == "played" {
    DoPlayed()
    return
  }
  // QUIT command
  if MudCmd == "quit" {
    DoQuit()
    return
  }
  // REFRESH command
  if MudCmd == "refresh" {
    DoRefresh()
    return
  }
  // REMOVE command
  if MudCmd == "remove" {
    DoRemove()
    return
  }
  // RESTORE command
  if MudCmd == "restore" {
    DoRestore(CmdStr)
    return
  }
  // ROOMINFO command
  if MudCmd == "roominfo" {
    DoRoomInfo()
    return
  }
  // SAVE command
  if MudCmd == "save" {
    DoSave()
    return
  }
  // SAY command
  if MudCmd == "say" {
    DoSay()
    return
  }
  // SELL command
  if MudCmd == "sell" {
    DoSell()
    return
  }
  // SHOW command
  if MudCmd == "show" {
    DoShow()
    return
  }
  // SIT command
  if MudCmd == "sit" {
    DoSit()
    return
  }
  // SLEEP command
  if MudCmd == "sleep" {
    DoSleep()
    return
  }
  // STAND command
  if MudCmd == "stand" {
    DoStand()
    return
  }
  // STATUS command
  if MudCmd == "status" {
    DoStatus()
    return
  }
  // STOP command
  if MudCmd == "stop" {
    DoStop()
    return
  }
  // TELL command
  if MudCmd == "tell" {
    DoTell()
    return
  }
  // TIME command
  if MudCmd == "time" {
    DoTime()
    return
  }
  // TITLE command
  if MudCmd == "title" {
    DoTitle()
    return
  }
  // TRAIN command
  if MudCmd == "train" {
    DoTrain()
    return
  }
  // WAKE command
  if MudCmd == "wake" {
    DoWake()
    return
  }
  // WEAR command
  if MudCmd == "wear" {
    DoWear()
    return
  }
  // WHERE command
  if MudCmd == "where" {
    DoWhere()
    return
  }
  // WHO command
  if MudCmd == "who" {
    DoWho()
    return
  }
  // WIELD command
  if MudCmd == "wield" {
    DoWield()
    return
  }
  pDnodeActor.PlayerOut += "Command is valid, but not implemented at this time."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Log it
  MudCmd = StrMakeFirstUpper(MudCmd)
  LogBuf = MudCmd
  LogBuf += " is in command array, but Do"
  LogBuf += MudCmd
  LogBuf += " is not coded."
  LogIt(LogBuf)
}

// Advance command
func DoAdvance() {
  var Level          int
  var LevelString    string
  var PlayerName     string
  var PlayerNameSave string
  var TargetName     string
  var TargetNameSave string

  DEBUGIT(1)
  PlayerName = pDnodeActor.PlayerName
  TargetName = StrGetWord(CmdStr, 2)
  PlayerNameSave = PlayerName
  TargetNameSave = TargetName
  PlayerName = StrMakeLower(PlayerName)
  TargetName = StrMakeLower(TargetName)
  Level = StrToInt(StrGetWord(CmdStr, 3))
  LevelString = fmt.Sprintf("%d", Level)
  if TargetName == "" {
    // No name given
    pDnodeActor.PlayerOut += "Advance who?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  // Get target Dnode pointer
  pDnodeTgt = GetTargetDnode(TargetName)
  if pDnodeTgt == nil {
    // Target player not found
    pDnodeActor.PlayerOut += TargetNameSave
    pDnodeActor.PlayerOut += " is not online.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if Level == pDnodeTgt.pPlayer.Level {
    // Advance to same level ... that's just plain silly
    pDnodeActor.PlayerOut += TargetNameSave
    pDnodeActor.PlayerOut += " is already at level "
    pDnodeActor.PlayerOut += LevelString
    pDnodeActor.PlayerOut += "."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if Level == 0 {
    // Advance to level 0 ... not valid
    pDnodeActor.PlayerOut += TargetNameSave
    pDnodeActor.PlayerOut += " cannot be advanced to level "
    pDnodeActor.PlayerOut += LevelString
    pDnodeActor.PlayerOut += "."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TargetNameSave = pDnodeTgt.PlayerName
  // Level up! or Down :(
  LogBuf = pDnodeTgt.PlayerName
  if Level > pDnodeTgt.pPlayer.Level {
    // Level up!
    LogBuf += " has been advanced to level "
  } else {
    // Level down :(
    LogBuf += " has been demoted to level "
  }
  LogBuf += LevelString
  LogBuf += " by "
  LogBuf += pDnodeActor.PlayerName
  LogIt(LogBuf)
  // Send message to player
  pDnodeActor.PlayerOut += TargetNameSave
  pDnodeActor.PlayerOut += " is now level "
  pDnodeActor.PlayerOut += LevelString
  pDnodeActor.PlayerOut += "."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Send message to target
  pDnodeTgt.PlayerOut += "\r\n"
  pDnodeTgt.PlayerOut += "&Y"
  pDnodeTgt.PlayerOut += PlayerNameSave
  if Level > pDnodeTgt.pPlayer.Level {
    // Level up!
    pDnodeTgt.PlayerOut += " has advanced you to level "
  } else {
    // Level down :(
    pDnodeTgt.PlayerOut += " has DEMOTED you to level "
  }
  pDnodeTgt.PlayerOut += LevelString
  pDnodeTgt.PlayerOut += "!"
  pDnodeTgt.PlayerOut += "&N"
  pDnodeTgt.PlayerOut += "\r\n"
  // Make it so
  pDnodeTgt.pPlayer.Level = Level
  pDnodeTgt.pPlayer.Experience = CalcLevelExperience(Level)
  PlayerSave(pDnodeTgt.pPlayer)
  // Prompt
  CreatePrompt(pDnodeTgt.pPlayer)
  pDnodeTgt.PlayerOut += GetOutput(pDnodeTgt.pPlayer)
  // Restore the player as a bonus to being advanced
  DoRestore("restore " + pDnodeTgt.pPlayer.Name)
}

// Afk command
func DoAfk() {
  if pDnodeActor.PlayerStateAfk {
    // Player returning from AFK
    pDnodeActor.PlayerStateAfk = false
    pDnodeActor.PlayerOut += "You are no longer Away From Keyboard"
    pDnodeActor.PlayerOut += "\r\n"
  } else {
    // Player going AFK
    pDnodeActor.PlayerOut += "You are now Away From Keyboard"
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerStateAfk = true
  }
  PlayerSave(pDnodeActor.pPlayer)
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Assist command
func DoAssist() {
  var AssistMsg       string
  var MobileId        string
  var PlayerNameCheck string
  var TargetNameCheck string
  var TargetNameSave  string
  var TargetNotHere   bool

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) < 2 {
    // No object or target
    pDnodeActor.PlayerOut += "Assist whom?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  PlayerNameCheck = pDnodeActor.PlayerName
  TargetNameCheck = StrGetWord(CmdStr, 2)
  TargetNameSave  = TargetNameCheck
  PlayerNameCheck = StrMakeLower(PlayerNameCheck)
  TargetNameCheck = StrMakeLower(TargetNameCheck)
  if PlayerNameCheck == TargetNameCheck {
    // Player is trying to assist themself
    pDnodeActor.PlayerOut += "You can't assist youself.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*********************
  //* Turning assist on *
  //*********************
  if TargetNameCheck == "on" {
    pDnodeActor.pPlayer.AllowAssist = true
    pDnodeActor.PlayerOut += "You are now accepting assists.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //**********************
  //* Turning assist off *
  //**********************
  if TargetNameCheck == "off" {
    pDnodeActor.pPlayer.AllowAssist = false
    pDnodeActor.PlayerOut += "You are now rejecting assists.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //******************
  //* Is target OK ? *
  //******************
  TargetNotHere = false
  pDnodeTgt = GetTargetDnode(TargetNameCheck)
  if pDnodeTgt == nil {
    // Target is not online and/or not in 'playing' state
    TargetNotHere = true
  } else {
    // Target is online and playing
    if pDnodeActor.pPlayer.RoomId != pDnodeTgt.pPlayer.RoomId {
      // Target is not in the same room
      TargetNotHere = true
    }
  }
  if TargetNotHere {
    // Tell player that target is not here
    pDnodeActor.PlayerOut += TargetNameSave
    pDnodeActor.PlayerOut += " is not here.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if !pDnodeTgt.PlayerStateFighting {
    // Tell player that target is not fighting
    pDnodeActor.PlayerOut += pDnodeTgt.pPlayer.Name
    pDnodeActor.PlayerOut += " is not fighting."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if !pDnodeTgt.pPlayer.AllowAssist {
    // Tell player that target is not accepting assistance
    pDnodeActor.PlayerOut += pDnodeTgt.pPlayer.Name
    pDnodeActor.PlayerOut += " is not accepting assistance."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //******************
  //* Send to player *
  //******************
  pDnodeActor.PlayerOut += "You begin assisting "
  pDnodeActor.PlayerOut += pDnodeTgt.pPlayer.Name
  pDnodeActor.PlayerOut += ".\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  //******************
  //* Send to target *
  //******************
  pDnodeTgt.PlayerOut += "\r\n"
  pDnodeTgt.PlayerOut += pDnodeActor.pPlayer.Name
  pDnodeTgt.PlayerOut += " begins assisting you.\r\n"
  CreatePrompt(pDnodeTgt.pPlayer)
  pDnodeTgt.PlayerOut += GetOutput(pDnodeTgt.pPlayer)
  //****************
  //* Send to room *
  //****************
  AssistMsg = pDnodeActor.PlayerName
  AssistMsg += " begins assisting "
  AssistMsg += pDnodeTgt.pPlayer.Name
  AssistMsg += "."
  pDnodeSrc = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, AssistMsg)
  //**************************
  //* Make the assist happen *
  //**************************
  MobileId = GetPlayerMobMobileId(pDnodeTgt.PlayerName)
  CreatePlayerMob(pDnodeActor.PlayerName, MobileId)
  pDnodeActor.PlayerStateFighting = true
}

// Buy command
func DoBuy() {
  var Cost       int
  var Desc1      string
  var ObjectId   string
  var ObjectName string
  var RoomId     string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  RoomId = pDnodeActor.pPlayer.RoomId
  if !IsShop(RoomId) {
    // Room is not a shop
    pDnodeActor.PlayerOut += "Find a shop."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  ObjectName = StrGetWord(CmdStr, 2)
  if ObjectName == "" {
    // No object given
    pDnodeActor.PlayerOut += "Buy what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if StrCountWords(CmdStr) > 2 {
    // Buy command not only takes 1 object
    pDnodeActor.PlayerOut += "The buy command must be followed by only one word."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pObject = nil
  IsShopObj(RoomId, ObjectName) // Sets pObject
  if pObject == nil {
    // Object not in shop for player to buy
    pDnodeActor.PlayerOut += "That item cannot be bought here."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  ObjectId = pObject.ObjectId
  Desc1 = pObject.Desc1
  Cost = pObject.Cost
  pObject = nil
  if pDnodeActor.pPlayer.Silver < Cost {
    // Player cannot afford item
    pDnodeActor.PlayerOut += "You cannot afford that item."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*******************
  //* Buy the object  *
  //*******************
  // Add object to player's inventory
  AddObjToPlayerInv(pDnodeActor, ObjectId)
  // Player receives some money
  SetMoney(pDnodeActor.pPlayer,'-', Cost, "Silver")
  // Send messages
  Buf = fmt.Sprintf("%d", Cost)
  TmpStr = Buf
  pDnodeActor.PlayerOut += "You buy "
  pDnodeActor.PlayerOut += Desc1
  pDnodeActor.PlayerOut += " for "
  pDnodeActor.PlayerOut += TmpStr
  pDnodeActor.PlayerOut += " piece(s) of silver."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Chat command
func DoChat() {
  var AllMsg    string
  var ChatMsg   string
  var PlayerMsg string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  ChatMsg = StrGetWords(CmdStr, 2)
  if StrGetLength(ChatMsg) < 1 {
    // Player did not enter any chat
    pDnodeActor.PlayerOut += "You start to chat, but, about what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*************
  //* Chat away *
  //*************
  PlayerMsg = "&G"
  PlayerMsg += "You Chat: "
  PlayerMsg += ChatMsg
  PlayerMsg += "&N"
  PlayerMsg += "\r\n"
  AllMsg = "\r\n"
  AllMsg += "&G"
  AllMsg += pDnodeActor.PlayerName
  AllMsg += " chats: "
  AllMsg += ChatMsg
  AllMsg += "&N"
  AllMsg += "\r\n"
  SendToAll(PlayerMsg, AllMsg)
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Color command
func DoColor() {
  DEBUGIT(1)
  TmpStr = StrGetWord(CmdStr, 2)
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr == "on" {
    // Turn color on
    pDnodeActor.pPlayer.Color = true
    PlayerSave(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += "You will now see &RP&Gr&Ye&Bt&Mt&Cy&N &RC&Go&Yl&Bo&Mr&Cs&N.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if TmpStr == "off" {
    // Turn color off
    pDnodeActor.pPlayer.Color = false
    PlayerSave(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += "Color is &Moff&N.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pDnodeActor.pPlayer.Color {
    // Color is on
    pDnodeActor.PlayerOut += "&CColor&N is &Mon&N.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  } else {
    // Color is off
    pDnodeActor.PlayerOut += "Color is off.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  }
}

// Consider command
func DoConsider() {
  var pMobile        *Mobile
  var HintMsg         string
  var LevelDiff       int
  var MobileName      string
  var PlayerName      string
  var PlayerNameCheck string
  var Target          string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) < 2 {
    // No target
    pDnodeActor.PlayerOut += "Consider whom or what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if StrCountWords(CmdStr) > 2 {
    // Two many targets
    pDnodeActor.PlayerOut += "Hmm, you are confused. Try 'help consider'."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  PlayerName = pDnodeActor.PlayerName
  PlayerNameCheck = PlayerName
  PlayerNameCheck = StrMakeLower(PlayerNameCheck)
  Target = StrGetWord(CmdStr, 2)
  MobileName = Target
  Target = StrMakeLower(Target)
  if Target == PlayerNameCheck {
    // Trying to consider self
    pDnodeActor.PlayerOut += "Consider yourself considered!"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if IsPlayer(Target) {
    // Trying to consider another player
    pDnodeActor.PlayerOut += "Why consider another player? Player killing is not allowed."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pMobile = IsMobInRoom(Target)
  if pMobile == nil {
    // Target mobile is not here
    pDnodeActor.PlayerOut += "There doesn't seem to be a(n) "
    pDnodeActor.PlayerOut += MobileName
    pDnodeActor.PlayerOut += " here.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  // Build message based on level difference
  LevelDiff = pMobile.Level - pDnodeActor.pPlayer.Level
  if LevelDiff < -6 {
    HintMsg = "&GDon't bother&N"
  } else if LevelDiff == -6 {
    HintMsg = "&GWay too easy&N"
  } else if LevelDiff == -5 {
    HintMsg = "&GVery easy&N"
  } else if LevelDiff == -4 {
    HintMsg = "&CEasy&N"
  } else if LevelDiff == -3 {
    HintMsg = "&CNo problem&N"
  } else if LevelDiff == -2 {
    HintMsg = "&CA worthy opponent&N"
  } else if LevelDiff == -1 {
    HintMsg = "&YYou might win&N"
  } else if LevelDiff == 0 {
    HintMsg = "&YTough fight&N"
  } else if LevelDiff == 1 {
    HintMsg = "&YLots of luck&N"
  } else if LevelDiff == 2 {
    HintMsg = "&RBad idea&N"
  } else if LevelDiff > 2 {
    HintMsg = "&RYOU ARE MAD&N"
  }
  //*****************
  //* Send messages *
  //*****************
  // Send message to player
  pDnodeActor.PlayerOut += "You consider "
  pDnodeActor.PlayerOut += pMobile.Desc1
  pDnodeActor.PlayerOut += "."
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += HintMsg
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  pMobile = nil
}

// Delete command
func DoDelete() {
  var AllMsg         string
  var Name           string
  var Password       string
  var Phrase         string
  var PlayerFileName string
  var PlayerMsg      string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) < 3 {
    pDnodeActor.PlayerOut += "You must provide your name and password."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  Name = StrGetWord(CmdStr, 2)
  Password = StrGetWord(CmdStr, 3)
  Phrase = StrGetWords(CmdStr, 4)
  if Name != pDnodeActor.PlayerName {
    pDnodeActor.PlayerOut += "Your name was not entered correctly."
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "Upper and lowercase letters must match."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if Password != pDnodeActor.pPlayer.Password {
    pDnodeActor.PlayerOut += "Your password was not entered correctly."
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "Upper and lowercase letters must match."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if Phrase == "" {
    pDnodeActor.PlayerOut += "If you really want delete your character,"
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "add the phrase: MAKE IT SO"
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "to the end of the command."
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "Like this: delete <name> <password> MAKE IT SO"
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "and you will be immediately DELETED!!"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if Phrase != "MAKE IT SO" {
    pDnodeActor.PlayerOut += "Ok, it seems that you provided a phrase"
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "but it is _not_ 'MAKE IT SO'."
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "The command must be like this: delete <name> <password> MAKE IT SO"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*****************
  //* Delete player *
  //*****************
  pDnodeActor.PlayerStateBye = true
  pDnodeActor.PlayerStatePlaying = false
  GrpLeave()
  LogBuf = pDnodeActor.PlayerName
  LogBuf += " issued the DELETE command"
  LogIt(LogBuf)
  // Delete Player file
  PlayerFileName = PLAYER_DIR
  PlayerFileName += pDnodeActor.PlayerName
  PlayerFileName += ".txt"
  err := Remove(PlayerFileName)
  if err != nil {
    LogBuf = "Communication::DoDelete - Failed to remove Player file: " + PlayerFileName + ". Error: " + err.Error()
    LogIt(LogBuf)
    os.Exit(1)
  }
  // Delete PlayerEqu file
  PlayerFileName = PLAYER_EQU_DIR
  PlayerFileName += pDnodeActor.PlayerName
  PlayerFileName += ".txt"
  _ = Remove(PlayerFileName)
  // Delete PlayerObj file
  PlayerFileName = PLAYER_OBJ_DIR
  PlayerFileName += pDnodeActor.PlayerName
  PlayerFileName += ".txt"
  _ = Remove(PlayerFileName)
  // Delete PlayerRoom file
  PlayerFileName = PLAYER_ROOM_DIR
  PlayerFileName += pDnodeActor.PlayerName
  PlayerFileName += ".txt"
  _ = Remove(PlayerFileName)
  // Send messages
  pDnodeSrc = pDnodeActor
  pDnodeTgt = nil
  PlayerMsg = "You have been DELETED!!!"
  PlayerMsg += "\r\n"
  AllMsg = "\r\n"
  AllMsg += pDnodeActor.PlayerName
  AllMsg += " has just DELETED $pHimselfHerself."
  AllMsg += "\r\n"
  pDnodeSrc = pDnodeActor
  pDnodeTgt = nil
  AllMsg = PronounSubstitute(AllMsg)
  pDnodeSrc = nil
  pDnodeTgt = nil
  SendToAll(PlayerMsg, AllMsg)
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Destroy command
func DoDestroy() {
  var ObjectName string
  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) == 1 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Destroy what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //****************************
  //* Does player have object? *
  //****************************
  TmpStr = StrGetWord(CmdStr, 2)
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  pObject = nil
  IsObjInPlayerInv(TmpStr) // Sets pObject
  if pObject == nil {
    pDnodeActor.PlayerOut += "You don't have a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += ".\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //******************
  //* Destroy object *
  //******************
  // Remove object from player's inventory
  RemoveObjFromPlayerInv(pObject.ObjectId, 1)
  // Send messages
  pDnodeActor.PlayerOut += "You help make the world cleaner by destroying "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += ".\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  pObject = nil
}

// Drink command
func DoDrink() {
  var DrinkMsg string
  var ObjectName string
  var RoomId string
  var RoomName string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrGetWord(CmdStr, 2) == "from" {
    // Toss out 'from', just extra verbage for player's benefit
    CmdStr = StrDelete(CmdStr, 5, 5)
  }
  if StrCountWords(CmdStr) == 1 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Drink from what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*************************
  //* Is this a drink room? *
  //*************************
  RoomId = pDnodeActor.pPlayer.RoomId
  if IsRoomType(RoomId, "Drink") {
    // Room contains something to drink
    RoomName = GetRoomName(RoomId)
    TmpStr = StrGetWord(CmdStr, 2)
    TmpStr = StrMakeLower(TmpStr)
    RoomName = StrMakeLower(RoomName)
    if StrIsWord(TmpStr, RoomName) {
      //*****************
      //* Player drinks *
      //*****************
      // Send messages
      pDnodeActor.PlayerOut += "You take a drink of clear &Ccool&N water."
      pDnodeActor.PlayerOut += "\r\n"
      DrinkMsg = pDnodeActor.PlayerName
      DrinkMsg += " takes a drink of clear &Ccool&N water."
      pDnodeSrc = pDnodeActor
      pDnodeTgt = pDnodeActor
      SendToRoom(pDnodeActor.pPlayer.RoomId, DrinkMsg)
      // Update thirst
      Drink(pDnodeActor.pPlayer, 20)
      PlayerSave(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      // Prompt
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  //****************************
  //* Does player have object? *
  //****************************
  TmpStr = StrGetWord(CmdStr, 2)
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  pObject = nil
  IsObjInPlayerInv(TmpStr) // Sets pObject
  if pObject == nil {
    // Object not found in player's inventory
    pDnodeActor.PlayerOut += "You don't have a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += " in your inventory.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*******************
  //* Is object drink? *
  //*******************
  pObject.Type = StrMakeLower(pObject.Type)
  if pObject.Type != "drink" {
    // Object is not a drink
    pDnodeActor.PlayerOut += "You can't drink "
    pDnodeActor.PlayerOut += pObject.Desc1
    pDnodeActor.PlayerOut += "."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //****************
  //* Drink object *
  //****************
  // Send messages
  pDnodeActor.PlayerOut += "You drink from "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += "."
  pDnodeActor.PlayerOut += "\r\n"
  DrinkMsg = pDnodeActor.PlayerName
  DrinkMsg += " drinks from "
  DrinkMsg += pObject.Desc1
  DrinkMsg += "."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, DrinkMsg)
  // Drink and remove object from player's inventory
  Drink(pDnodeActor.pPlayer, pObject.DrinkPct)
  PlayerSave(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  RemoveObjFromPlayerInv(pObject.ObjectId, 1)
  // Clean up and give prompt
  pObject = nil
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Drop command
func DoDrop() {
  var DropMsg string
  var ObjectName string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) == 1 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Drop what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //****************************
  //* Does player have object? *
  //****************************
  TmpStr = StrGetWord(CmdStr, 2)
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  pObject = nil
  IsObjInPlayerInv(TmpStr) // Sets pObject
  if pObject == nil {
    pDnodeActor.PlayerOut += "You don't have a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += " in your inventory.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //***************
  //* Drop object *
  //***************
  // Remove object from player's inventory
  RemoveObjFromPlayerInv(pObject.ObjectId, 1)
  // Send messages
  pDnodeActor.PlayerOut += "You drop "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += "."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  DropMsg = pDnodeActor.PlayerName
  DropMsg += " drops "
  DropMsg += pObject.Desc1
  DropMsg += "."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, DropMsg)
  // Add object to room
  AddObjToRoom(pDnodeActor.pPlayer.RoomId, pObject.ObjectId)
  pObject = nil
}

// Eat command
func DoEat() {
  var EatMsg string
  var ObjectName string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) == 1 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Eat what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //****************************
  //* Does player have object? *
  //****************************
  TmpStr = StrGetWord(CmdStr, 2)
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  pObject = nil
  IsObjInPlayerInv(TmpStr) // Sets pObject
  if pObject == nil {
    pDnodeActor.PlayerOut += "You don't have a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += " in your inventory.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*******************
  //* Is object food? *
  //*******************
  pObject.Type = StrMakeLower(pObject.Type)
  if pObject.Type != "food" {
    pDnodeActor.PlayerOut += "You can't eat "
    pDnodeActor.PlayerOut += pObject.Desc1
    pDnodeActor.PlayerOut += "."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //**************
  //* Eat object *
  //**************
  // Send messages
  pDnodeActor.PlayerOut += "You eat "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += "."
  pDnodeActor.PlayerOut += "\r\n"
  EatMsg = pDnodeActor.PlayerName
  EatMsg += " eats "
  EatMsg += pObject.Desc1
  EatMsg += "."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, EatMsg)
  // Eat and remove object from player's inventory
  Eat(pDnodeActor.pPlayer, pObject.FoodPct)
  PlayerSave(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  RemoveObjFromPlayerInv(pObject.ObjectId, 1)
  // Clean up and give prompt
  pObject = nil
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Emote command
func DoEmote() {
  var EmoteMsg string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  EmoteMsg = StrGetWords(CmdStr, 2)
  if StrGetLength(EmoteMsg) < 1 {
    // Player did not enter anything to say
    pDnodeActor.PlayerOut += "You try to show emotion, but fail."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*******************
  //* Emote something *
  //*******************
  pDnodeActor.PlayerOut += "&C"
  pDnodeActor.PlayerOut += pDnodeActor.PlayerName
  pDnodeActor.PlayerOut += " "
  pDnodeActor.PlayerOut += EmoteMsg
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  TmpStr = "&C"
  TmpStr += pDnodeActor.PlayerName
  TmpStr += " "
  TmpStr += EmoteMsg
  TmpStr += "&N"
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, TmpStr)
}

// Equipment command
func DoEquipment() {
  DEBUGIT(1)
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  ShowPlayerEqu(pDnodeActor)
}

// Examine command
func DoExamine() {
  var ObjectFound bool
  var ObjectName string
  var ObjectType string
  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) == 1 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Examine what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //***********************************************
  //* Ok, object, object ... where is the object? *
  //***********************************************
  ObjectFound = false
  TmpStr = StrGetWord(CmdStr, 2)
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  // Check room
  pObject = nil
  IsObjInRoom(TmpStr) // Sets pObject
  if pObject != nil {
    // Object is in the room
    ObjectFound = true
  } else {
    // Check player inventory
    pObject = nil
    IsObjInPlayerInv(TmpStr) // Sets pObject
    if pObject != nil {
      // Object is in player's inventory
      ObjectFound = true
    } else {
      // Check player equipment
      pObject = nil
      IsObjInPlayerEqu(TmpStr) // Sets pObject
      if pObject != nil {
        // Object is in player's equipment
        ObjectFound = true
      }
    }
  }
  if !ObjectFound {
    // Object can't be found
    pDnodeActor.PlayerOut += "There doesn't seem to be a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += " here.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //************************************
  //* Ojbect was found, now examine it *
  //************************************
  // Send messages
  pDnodeActor.PlayerOut += "You examine "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += "."
  pDnodeActor.PlayerOut += "\r\n"
  // Examine object
  pDnodeActor.PlayerOut += "Object type: "
  pDnodeActor.PlayerOut += pObject.Type
  pDnodeActor.PlayerOut += "\r\n"
  ObjectType = pObject.Type
  ObjectType = StrMakeLower(ObjectType)
  if ObjectType == "weapon" {
    // Object is a weapon
    pDnodeActor.PlayerOut += "Weapon type: "
    pDnodeActor.PlayerOut += pObject.WeaponType
    pDnodeActor.PlayerOut += "\r\n"
  }
  ExamineObj(pObject.ObjectId)
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  pObject = nil
}

// Flee command
func DoFlee() {
  var CandidateCount int
  var CandidateList string
  var CandidateTarget int
  var FleeMsg string
  var MobileId string
  var MobileIdSave string
  var MobPlayerFile *os.File
  var MobPlayerFileName string
  var MudCmdIsExit string
  var PlayerName1 string
  var PlayerName2 string
  var RoomIdBeforeFleeing string
  var Target string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if !pDnodeActor.PlayerStateFighting {
    pDnodeActor.PlayerOut += "You can't flee, you aren't fighting.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TmpStr = StrGetWord(CmdStr, 2)
  if TmpStr == "" {
    // No direction given
    pDnodeActor.PlayerOut += "Aimless fleeing is not allowed.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  RoomIdBeforeFleeing = pDnodeActor.pPlayer.RoomId
  MudCmdIsExit = "go"
  sMudCmdIsExit := MudCmdIsExit
  if !IsExit(sMudCmdIsExit) {
    // Direction given is not valid
    pDnodeActor.PlayerOut += "Flee where?\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //************************************
  //* Player has been moved, they fled *
  //************************************
  // Let everyone in room know they fled
  FleeMsg += "&R"
  FleeMsg += pDnodeActor.PlayerName
  FleeMsg += " has fled for $pHisHers life!!"
  FleeMsg += "&N"
  FleeMsg = PronounSubstitute(FleeMsg)
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(RoomIdBeforeFleeing, FleeMsg)
  // Let player know they succeeded in fleeing
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "&R"
  pDnodeActor.PlayerOut += "You have sucessfully fled for your life!!"
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  //***********
  //* Cleanup *
  //***********
  pDnodeActor.PlayerStateFighting = false
  if pDnodeActor.pPlayer.pPlayerFollowers[0] != nil {
    // Player is following someone, stop following
    DoFollow(pDnodeActor, "follow none")
  }
  PlayerName1 = pDnodeActor.PlayerName
  // Get mobile id for mob that fleeing player was fighting
  MobileIdSave = GetPlayerMobMobileId(PlayerName1)
  // Delete PlayerMob file -- player is no longer attacking mob
  DeletePlayerMob(PlayerName1)
  // See if a mob is whacking player
  MobPlayerFileName = MOB_PLAYER_DIR
  MobPlayerFileName += PlayerName1
  MobPlayerFileName += ".txt"
  var err error
  MobPlayerFile, err = os.Open(MobPlayerFileName)
  if err != nil {
    // If MobPlayer does not exist, then no mob is fighting player
    return
	}
	_ = MobPlayerFile.Close()
  //***************************
  //* Make mob switch targets *
  //***************************
  // Delete fighting mobiles from MobPlayer file
  SetpDnodeCursorFirst()
  for !EndOfDnodeList() {
    // Loop thru all connections
    pDnodeOthers = GetDnode()
    if pDnodeOthers.PlayerStateFighting {
      // Players who are fighting
      if RoomIdBeforeFleeing == pDnodeOthers.pPlayer.RoomId {
        // In the same room
        PlayerName2 = pDnodeOthers.PlayerName
        MobileId = GetPlayerMobMobileId(PlayerName2)
        DeleteMobPlayer(PlayerName1, MobileId)
        if MobileId == MobileIdSave {
          // Add player to candidate list for MobileIdSave
          CandidateList += PlayerName2
          CandidateList += " "
        }
      }
    }
    SetpDnodeCursorNext()
  }
  // Re-position pDnodeCursor
  RepositionDnodeCursor()
  // Put mobiles that are not fighting back in room
  PutMobBackInRoom(PlayerName1, RoomIdBeforeFleeing)
  // Player is gone, so delete MobPlayer completely
  DeleteMobPlayer(PlayerName1, "file")
  // Select a new target for MobileIdSave
  if StrGetLength(CandidateList) == 0 {
    // No available target for MobileIdSave
    return
  }
  CandidateCount = StrCountWords(CandidateList)
  CandidateTarget = GetRandomNumber(CandidateCount)
  Target = StrGetWord(CandidateList, CandidateTarget)
  CreateMobPlayer(Target, MobileIdSave)
}

// Follow command
func DoFollow(pDnode *Dnode, CmdStr1 string) {
  var pDnodeGrpLdr *Dnode
  var pDnodeGrpMem *Dnode
  var i int
  var j int
  var Target string
  var TargetInGroup bool

  DEBUGIT(1)
  CmdStr = CmdStr1
  i = 0
  j = 0
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  Target = StrGetWord(CmdStr, 2)
  TmpStr = Target
  TmpStr = StrMakeLower(TmpStr)
  if Target == "" {
    // Follow with no target
    pDnode.PlayerOut += "Who would you like to follow?\r\n"
    CreatePrompt(pDnode.pPlayer)
    pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
    return
  }
  //******************
  //* Stop following *
  //******************
  if TmpStr == "none" {
    // Player wants to stop following
    if pDnode.pPlayer.pPlayerFollowers[0] == nil {
      // Player is not following anyone
      pDnode.PlayerOut += "Ok. But alas, you were not following anyone"
      pDnode.PlayerOut += "\r\n"
      CreatePrompt(pDnode.pPlayer)
      pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
      return
    }
    for i = 1; i < GRP_LIMIT; i++ {
      // Search thru player-being-followed's list of followers
      if pDnode.pPlayer == pDnode.pPlayer.pPlayerFollowers[0].pPlayerFollowers[i] {
        // Found player, set to null
        pDnode.pPlayer.pPlayerFollowers[0].pPlayerFollowers[i] = nil
        j = i
        break
      }
    }
    // Compact the list of followers, so new followers are at the end
    for i = j; i < GRP_LIMIT-1; i++ {
      pDnode.pPlayer.pPlayerFollowers[0].pPlayerFollowers[i] = pDnode.pPlayer.pPlayerFollowers[0].pPlayerFollowers[i+1]
      pDnode.pPlayer.pPlayerFollowers[0].pPlayerFollowers[i+1] = nil
    }
    pDnode.PlayerOut += "You stop following "
    pDnode.PlayerOut += pDnode.pPlayer.pPlayerFollowers[0].Name
    pDnode.PlayerOut += ".\r\n"
    CreatePrompt(pDnode.pPlayer)
    pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
    pDnode.pPlayer.pPlayerFollowers[0] = nil
    return
  }
  //******************
  //* List followers *
  //******************
  if TmpStr == "list" {
    if pDnode.pPlayer.pPlayerFollowers[0] == nil {
      pDnode.PlayerOut += "You are not following anyone.\r\n"
    } else {
      pDnode.PlayerOut += "You are following "
      pDnode.PlayerOut += pDnode.pPlayer.pPlayerFollowers[0].Name
      pDnode.PlayerOut += ".\r\n"
    }
    for i = 1; i < GRP_LIMIT; i++ {
      if pDnode.pPlayer.pPlayerFollowers[i] != nil {
        if i == 1 {
          pDnode.PlayerOut += "Followers\r\n"
          pDnode.PlayerOut += "---------\r\n"
        }
        pDnode.PlayerOut += pDnode.pPlayer.pPlayerFollowers[i].Name
        pDnode.PlayerOut += "\r\n"
      }
    }
    CreatePrompt(pDnode.pPlayer)
    pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
    return
  }
  //*******************************
  //* Player is already following *
  //*******************************
  if pDnode.pPlayer.pPlayerFollowers[0] != nil {
    pDnode.PlayerOut += "You are already following "
    pDnode.PlayerOut += pDnode.pPlayer.pPlayerFollowers[0].Name
    pDnode.PlayerOut += ".\r\n"
    CreatePrompt(pDnode.pPlayer)
    pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
    return
  }
  //******************
  //* Target online? *
  //******************
  pDnodeTgt = GetTargetDnode(Target)
  if pDnodeTgt == nil {
    // Target not online
    pDnode.PlayerOut += Target
    pDnode.PlayerOut += " is not online.\r\n"
    CreatePrompt(pDnode.pPlayer)
    pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
    return
  }
  //***************************
  //* Can not follow yourself *
  //***************************
  if pDnode == pDnodeTgt {
    // Player trying to follow themself
    pDnode.PlayerOut += "You can not follow yourself.\r\n"
    CreatePrompt(pDnode.pPlayer)
    pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
    return
  }
  //*************************************
  //* Player and target must be grouped *
  //*************************************
  if pDnode.pPlayer.pPlayerGrpMember[0] == nil {
    // Player is not in a group
    pDnode.PlayerOut += "You must be in a group before you can follow."
    pDnode.PlayerOut += "\r\n"
    CreatePrompt(pDnode.pPlayer)
    pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
    return
  }
  pDnodeGrpLdr = GetTargetDnode(pDnode.pPlayer.pPlayerGrpMember[0].Name)
  TargetInGroup = false
  for i = 0; i < GRP_LIMIT; i++ {
    // For each member of leader's group including the leader
    if pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i] == nil {
      // Done looping through group members
      break
    }
    // Get group member's dnode
    pDnodeGrpMem = GetTargetDnode(pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i].Name)
    if pDnodeGrpMem == pDnodeTgt {
      // Target member found
      TargetInGroup = true
    }
  }
  if !TargetInGroup {
    // Target is not grouped with player
    pDnode.PlayerOut += Target
    pDnode.PlayerOut += " is not in the group.\r\n"
    CreatePrompt(pDnode.pPlayer)
    pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
    return
  }
  //************************************
  //* Player can not follow a follower *
  //************************************
  for i = 1; i < GRP_LIMIT; i++ {
    if pDnodeTgt.pPlayer == pDnode.pPlayer.pPlayerFollowers[i] {
      pDnode.PlayerOut += "Can not! "
      pDnode.PlayerOut += pDnodeTgt.pPlayer.Name
      pDnode.PlayerOut += " is following YOU.\r\n"
      CreatePrompt(pDnode.pPlayer)
      pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
      return
    }
  }
  //***********************
  //* Follow - make it so *
  //***********************
  // Message to player
  pDnode.PlayerOut += "You begin to follow "
  pDnode.PlayerOut += pDnodeTgt.PlayerName
  pDnode.PlayerOut += ".\r\n"
  CreatePrompt(pDnode.pPlayer)
  pDnode.PlayerOut += GetOutput(pDnode.pPlayer)
  // Message to player being followed
  pDnodeTgt.PlayerOut += "\r\n"
  pDnodeTgt.PlayerOut += pDnode.PlayerName
  pDnodeTgt.PlayerOut += " begins to follow you."
  pDnodeTgt.PlayerOut += "\r\n"
  CreatePrompt(pDnodeTgt.pPlayer)
  pDnodeTgt.PlayerOut += GetOutput(pDnodeTgt.pPlayer)
  // Make 'em follow
  pDnode.pPlayer.pPlayerFollowers[0] = pDnodeTgt.pPlayer
  for i = 1; i < GRP_LIMIT; i++ {
    // Loop through target's list of followers to find an empty slot
    if pDnodeTgt.pPlayer.pPlayerFollowers[i] == nil {
      // Found a slot for the new follower
      break
    }
  }
  pDnodeTgt.pPlayer.pPlayerFollowers[i] = pDnode.pPlayer
}

// Get command
func DoGet() {
  var GetMsg string
  var ObjectName string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) == 1 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Get what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //**********************************************
  //* See if object is in room and can be gotten *
  //**********************************************
  TmpStr = StrGetWord(CmdStr, 2)
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  pObject = nil
  IsObjInRoom(TmpStr) // Sets pObject
  if pObject == nil {
    pDnodeActor.PlayerOut += "There doesn't seem to be a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += " here.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pObject.Type == "NoTake" && !pDnodeActor.pPlayer.Admin {
    // Only administrators can 'get' a 'notake' object
    pDnodeActor.PlayerOut += "You may not take that."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    pObject = nil
    return
  }
  //**********************
  //* So take the object *
  //**********************
  // Remove object from room
  RemoveObjFromRoom(pObject.ObjectId)
  // Send messages
  pDnodeActor.PlayerOut += "You get "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += ".\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  GetMsg = pDnodeActor.PlayerName
  GetMsg += " gets "
  GetMsg += pObject.Desc1
  GetMsg += "."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, GetMsg)
  // Add object to player's inventory
  AddObjToPlayerInv(pDnodeTgt, pObject.ObjectId)
  pObject = nil
}

// Give command
func DoGive() {
  var GiveMsg string
  var ObjectName string
  var PlayerName string
  var TargetName string
  var TargetNotHere bool

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) < 2 {
    // No object or target
    pDnodeActor.PlayerOut += "Give what and to whom?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if StrCountWords(CmdStr) < 3 {
    // No target
    pDnodeActor.PlayerOut += "Give it to whom?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //****************************
  //* Does player have object? *
  //****************************
  TmpStr = StrGetWord(CmdStr, 2)
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  pObject = nil
  IsObjInPlayerInv(TmpStr) // Sets pObject
  if pObject == nil {
    // Player does not have object
    pDnodeActor.PlayerOut += "You don't have a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += " in your inventory.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //****************
  //* Is target Ok *
  //****************
  TargetNotHere = false
  TmpStr = StrGetWord(CmdStr, 3)
  TargetName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  PlayerName = pDnodeActor.PlayerName
  PlayerName = StrMakeLower(PlayerName)
  if PlayerName == TmpStr {
    // Player is trying to give something to themself
    pDnodeActor.PlayerOut += "Giving something to youself is just plain silly!\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pDnodeTgt = GetTargetDnode(TmpStr)
  if pDnodeTgt == nil {
    // Target is not online and/or not in 'playing' state
    TargetNotHere = true
  } else {
    // Target is online and playing
    if pDnodeActor.pPlayer.RoomId != pDnodeTgt.pPlayer.RoomId {
      // Target is not in the same room
      TargetNotHere = true
    }
  }
  if TargetNotHere {
    // Tell player that target is not here
    pDnodeActor.PlayerOut += TargetName
    pDnodeActor.PlayerOut += " is not here.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pDnodeTgt.pPlayer.Position == "sleep" {
    // Tell player that target is sleeping
    pDnodeActor.PlayerOut += pDnodeTgt.pPlayer.Name
    pDnodeActor.PlayerOut += " is sleeping."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //******************
  //* Send to player *
  //******************
  pDnodeActor.PlayerOut += "You give "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += " to "
  pDnodeActor.PlayerOut += pDnodeTgt.pPlayer.Name
  pDnodeActor.PlayerOut += ".\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  //******************
  //* Send to target *
  //******************
  pDnodeTgt.PlayerOut += "\r\n"
  pDnodeTgt.PlayerOut += pDnodeActor.pPlayer.Name
  pDnodeTgt.PlayerOut += " gives you "
  pDnodeTgt.PlayerOut += pObject.Desc1
  pDnodeTgt.PlayerOut += ".\r\n"
  CreatePrompt(pDnodeTgt.pPlayer)
  pDnodeTgt.PlayerOut += GetOutput(pDnodeTgt.pPlayer)
  //****************
  //* Send to room *
  //****************
  GiveMsg = pDnodeActor.PlayerName
  GiveMsg += " gives "
  GiveMsg += pObject.Desc1
  GiveMsg += " to "
  GiveMsg += pDnodeTgt.pPlayer.Name
  GiveMsg += "."
  pDnodeSrc = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, GiveMsg)
  //*****************************
  //* Transfer object ownership *
  //*****************************
  RemoveObjFromPlayerInv(pObject.ObjectId, 1)
  AddObjToPlayerInv(pDnodeTgt, pObject.ObjectId)
  pObject = nil
}

// Go command
func DoGo() {
  var MudCmdIsExit string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if pDnodeActor.pPlayer.Position == "sit" {
    // Player is sitting
    pDnodeActor.PlayerOut += "You must be standing before you can move."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TmpStr = StrGetWord(CmdStr, 2)
  if TmpStr == "" {
    // No direction given
    pDnodeActor.PlayerOut += "Aimless wandering is not allowed.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //***************
  //* Try to move *
  //***************
  MudCmdIsExit = "go"
  var sMudCmdIsExit string
  sMudCmdIsExit = MudCmdIsExit
  if IsExit(sMudCmdIsExit) {
    // Player has been moved
    return
  }
  //********************************
  //* Direction given is not valid *
  //********************************
  pDnodeActor.PlayerOut += "Go where?\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Goto command
func DoGoTo() {
  var GoToMsg string
  var RoomId string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  RoomId = StrGetWord(CmdStr, 2)
  if RoomId == "" {
    pDnodeActor.PlayerOut += "A destination is needed.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if !IsRoom(RoomId) {
    pDnodeActor.PlayerOut += "Go to where?\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*************
  //* GoTo Room *
  //*************
  // Send GoTo departure message
  GoToMsg = pDnodeActor.pPlayer.GoToDepart
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, GoToMsg)
  // GoTo room
  pDnodeActor.pPlayer.RoomId = RoomId
  DoLook("")
  // Send GoTo arrival message
  GoToMsg = pDnodeActor.pPlayer.GoToArrive
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, GoToMsg)
}

// GoToArrive command
func DoGoToArrive() {
  var GoToArrive string

  DEBUGIT(1)
  TmpStr = StrGetWord(CmdStr, 2)
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr == "" { // Player entered 'gotoarrive' by itself
    if pDnodeActor.pPlayer.GoToArrive == "" { // Player has no arrival message
      pDnodeActor.PlayerOut += "You do not have an arrival message"
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    } else { // Show player's arrival message
      pDnodeActor.PlayerOut += "Your arrival message is: "
      pDnodeActor.PlayerOut += pDnodeActor.pPlayer.GoToArrive
      pDnodeActor.PlayerOut += "&N" // In case arrive msg is messed up
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  if TmpStr == "none" { // Player entered 'gotoarrive none'
    if pDnodeActor.pPlayer.GoToArrive == "" { // Player has no arrival message
      pDnodeActor.PlayerOut += "You did not have an arrival message and you still do not have an arrival message"
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    } else {
      pDnodeActor.pPlayer.GoToArrive = ""
      pDnodeActor.PlayerOut += "Your arrival message has been removed.\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  TmpStr = StrGetWords(CmdStr, 2)
  GoToArrive = TmpStr
  // Strip out color codes so arrival message length can be checked
  StrReplace(&TmpStr, "&N", "")
  StrReplace(&TmpStr, "&K", "")
  StrReplace(&TmpStr, "&R", "")
  StrReplace(&TmpStr, "&G", "")
  StrReplace(&TmpStr, "&Y", "")
  StrReplace(&TmpStr, "&B", "")
  StrReplace(&TmpStr, "&M", "")
  StrReplace(&TmpStr, "&C", "")
  StrReplace(&TmpStr, "&W", "")
  if StrGetLength(TmpStr) > 60 {
    pDnodeActor.PlayerOut += "Arrival message must be less than 61 characters, color codes do not count.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pDnodeActor.pPlayer.GoToArrive = GoToArrive
  pDnodeActor.PlayerOut += "Your arrival message has been set.\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// GoToDepart command
func DoGoToDepart() {
  var GoToDepart string

  DEBUGIT(1)
  TmpStr = StrGetWord(CmdStr, 2)
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr == "" { // Player entered 'gotodepart' by itself
    if pDnodeActor.pPlayer.GoToDepart == "" { // Player has no departure message
      pDnodeActor.PlayerOut += "You do not have a departure message"
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    } else { // Show player's departure message
      pDnodeActor.PlayerOut += "Your departure message is: "
      pDnodeActor.PlayerOut += pDnodeActor.pPlayer.GoToDepart
      pDnodeActor.PlayerOut += "&N" // In case depart msg is messed up
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  if TmpStr == "none" { // Player entered 'gotodepart none'
    if pDnodeActor.pPlayer.GoToDepart == "" { // Player has no departure message
      pDnodeActor.PlayerOut += "You did not have an departure message and you still do not have an departure message"
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    } else {
      pDnodeActor.pPlayer.GoToDepart = ""
      pDnodeActor.PlayerOut += "Your departure message has been removed.\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  TmpStr = StrGetWords(CmdStr, 2)
  GoToDepart = TmpStr
  // Strip out color codes so arrival message length can be checked
  StrReplace(&TmpStr, "&N", "")
  StrReplace(&TmpStr, "&K", "")
  StrReplace(&TmpStr, "&R", "")
  StrReplace(&TmpStr, "&G", "")
  StrReplace(&TmpStr, "&Y", "")
  StrReplace(&TmpStr, "&B", "")
  StrReplace(&TmpStr, "&M", "")
  StrReplace(&TmpStr, "&C", "")
  StrReplace(&TmpStr, "&W", "")
  if StrGetLength(TmpStr) > 60 {
    pDnodeActor.PlayerOut += "Departure message must be less than 61 characters, color codes do not count.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pDnodeActor.pPlayer.GoToDepart = GoToDepart
  pDnodeActor.PlayerOut += "Your departure message has been set.\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

 // Group command
func DoGroup() {
  var (
    pDnodeGrpLdr    *Dnode // Group leader
    i               int
    j               int
    GrpFull         bool
    PlayerNameCheck string
    TargetNameCheck string
    TargetNameSave  string
  )

  DEBUGIT(1)
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  PlayerNameCheck = pDnodeActor.PlayerName
  TargetNameCheck = StrGetWord(CmdStr, 2)
  TargetNameSave = TargetNameCheck
  PlayerNameCheck = StrMakeLower(PlayerNameCheck)
  TargetNameCheck = StrMakeLower(TargetNameCheck)
  //************************
  //* Group with no target *
  //************************
  if StrGetLength(TargetNameCheck) < 1 {
    // No target given
    pDnodeActor.PlayerOut += "Group with whom?\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*************************
  //* List members of group *
  //*************************
  if TargetNameCheck == "list" {
    // List group
    if pDnodeActor.pPlayer.pPlayerGrpMember[0] == nil {
      // Player is not in a group
      pDnodeActor.PlayerOut += "You are not in a group.\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
    // Player is in a group, show members
    pDnodeActor.PlayerOut += pDnodeActor.pPlayer.pPlayerGrpMember[0].Name
    pDnodeActor.PlayerOut += " \r\n"
    j = StrGetLength(pDnodeActor.pPlayer.pPlayerGrpMember[0].Name)
    for i = 1; i < j+1; i++ {
      pDnodeActor.PlayerOut += "-"
    }
    for i = 1; i < GRP_LIMIT; i++ {
      // List group members
      if pDnodeActor.pPlayer.pPlayerGrpMember[0].pPlayerGrpMember[i] != nil {
        pDnodeActor.PlayerOut += " \r\n"
        pDnodeActor.PlayerOut += pDnodeActor.pPlayer.pPlayerGrpMember[0].pPlayerGrpMember[i].Name
      }
    }
    pDnodeActor.PlayerOut += " \r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //***********************
  //* Turning grouping on *
  //***********************
  if TargetNameCheck == "on" {
    pDnodeActor.pPlayer.AllowGroup = true
    pDnodeActor.PlayerOut += "You are now accepting requests to group.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //************************
  //* Turning grouping off *
  //************************
  if TargetNameCheck == "off" {
    pDnodeActor.pPlayer.AllowGroup = false
    pDnodeActor.PlayerOut += "You are now rejecting requests to group.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*******************************
  //* Player is leaving the group *
  //*******************************
  if TargetNameCheck == "none" {
    if pDnodeActor.pPlayer.pPlayerGrpMember[0] == nil {
      // Player is not in a group
      pDnodeActor.PlayerOut += "You are not in a group.\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
    GrpLeave()
    return
  }
  //*********************************************************
  //* Trying to create a new group or trying to add members *
  //*********************************************************
  if TargetNameCheck == PlayerNameCheck {
    // Trying to group with self
    if pDnodeActor.pPlayer.pPlayerGrpMember[0] == nil {
      // Player is not in a group
      pDnodeActor.PlayerOut += "One is a lonely number, try grouping with another player.\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
    // Player is in a group
    pDnodeActor.PlayerOut += "One is a lonely number, but wait, you are already in group.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pDnodeGrpLdr = GetTargetDnode(TargetNameCheck)
  if pDnodeGrpLdr == nil {
    // New group member ... not online
    pDnodeActor.PlayerOut += TargetNameSave
    pDnodeActor.PlayerOut += " is not online.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pDnodeActor.pPlayer.pPlayerGrpMember[0] == pDnodeGrpLdr.pPlayer {
    // Player is trying to group with their group's leader
    pDnodeActor.PlayerOut += "You are already grouped with "
    pDnodeActor.PlayerOut += pDnodeGrpLdr.PlayerName
    pDnodeActor.PlayerOut += ", who is the leader.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pDnodeActor.pPlayer.RoomId != pDnodeGrpLdr.pPlayer.RoomId {
    // New group member is not in same room as leader
    pDnodeActor.PlayerOut += pDnodeGrpLdr.PlayerName
    pDnodeActor.PlayerOut += " is not here.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pDnodeActor.pPlayer.pPlayerGrpMember[0] != nil {
    // Player is in a group
    if pDnodeActor.pPlayer != pDnodeActor.pPlayer.pPlayerGrpMember[0] {
      // But is not the leader, only leader may add members
      pDnodeActor.PlayerOut += "You are not the leader. Leader must add members.\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  if pDnodeGrpLdr.pPlayer.pPlayerGrpMember[0] != nil {
    // New group member already in a group
    pDnodeActor.PlayerOut += pDnodeGrpLdr.PlayerName
    pDnodeActor.PlayerOut += " is already in a group.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  GrpFull = true
  for i = 0; i < GRP_LIMIT; i++ {
    // Is group full 
    if pDnodeActor.pPlayer.pPlayerGrpMember[i] == nil {
      // Found an empty member slot
      j = i
      GrpFull = false
      break
    }
  }
  if GrpFull {
    // Group is full
    pDnodeActor.PlayerOut += "Your group is full, maximum of "
    Buf = fmt.Sprintf("%d", GRP_LIMIT)
    TmpStr = Buf
    TmpStr = StrTrimLeft(TmpStr)
    TmpStr = StrTrimRight(TmpStr)
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += " members allowed.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if !pDnodeGrpLdr.pPlayer.AllowGroup {
    // New group member is not accepting group requests
    pDnodeActor.PlayerOut += pDnodeGrpLdr.PlayerName
    pDnodeActor.PlayerOut += " is not accepting requests to group.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //***********************************
  //* Ok, done checking ... group 'em *
  //***********************************
  if pDnodeActor.pPlayer.pPlayerGrpMember[0] == nil {
    // Forming new group
    pDnodeActor.pPlayer.pPlayerGrpMember[0] = pDnodeActor.pPlayer
    pDnodeGrpLdr.pPlayer.pPlayerGrpMember[0] = pDnodeActor.pPlayer
    pDnodeActor.pPlayer.pPlayerGrpMember[1] = pDnodeGrpLdr.pPlayer
    pDnodeActor.PlayerOut += "You have formed a new group with "
    pDnodeActor.PlayerOut += pDnodeGrpLdr.PlayerName
    pDnodeActor.PlayerOut += " as your first member.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  } else {
    // Adding a member to an existing group
    pDnodeActor.pPlayer.pPlayerGrpMember[i] = pDnodeGrpLdr.pPlayer
    pDnodeGrpLdr.pPlayer.pPlayerGrpMember[0] = pDnodeActor.pPlayer
    pDnodeActor.PlayerOut += pDnodeGrpLdr.PlayerName
    pDnodeActor.PlayerOut += " has been added to your group.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  }
  pDnodeGrpLdr.PlayerOut += "\r\n"
  pDnodeGrpLdr.PlayerOut += "You have joined a group, "
  pDnodeGrpLdr.PlayerOut += pDnodeActor.PlayerName
  pDnodeGrpLdr.PlayerOut += " is the leader.\r\n"
  CreatePrompt(pDnodeGrpLdr.pPlayer)
  pDnodeGrpLdr.PlayerOut += GetOutput(pDnodeGrpLdr.pPlayer)
}

// Gsay command
func DoGsay() {
  var pDnodeGrpLdr *Dnode // Group leader
  var pDnodeGrpMem *Dnode // Group members
  var GsayMsg       string
  var i             int

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if pDnodeActor.pPlayer.pPlayerGrpMember[0] == nil {
    // Player is not in a group
    pDnodeActor.PlayerOut += "You are not in a group.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  GsayMsg = StrGetWords(CmdStr, 2)
  if StrGetLength(GsayMsg) < 1 {
    // Player typed gsay but did not type a message
    pDnodeActor.PlayerOut += "Are you trying to say something to the group?\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*****************
  //* Send the gsay *
  //*****************
  pDnodeActor.PlayerOut += "&C"
  pDnodeActor.PlayerOut += "You say to the group: "
  pDnodeActor.PlayerOut += GsayMsg
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Get group leader
  pDnodeGrpLdr = GetTargetDnode(pDnodeActor.pPlayer.pPlayerGrpMember[0].Name)
  for i = 0; i < GRP_LIMIT; i++ {
    // For each member of leader's group including the leader
    if pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i] == nil {
      // Done looping through group members
      return
    }
    // Send gsay to other group members
    pDnodeGrpMem = GetTargetDnode(pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i].Name)
    if pDnodeActor == pDnodeGrpMem {
      // Do not send gsay to originating player
      continue
    }
    pDnodeGrpMem.PlayerOut += "\r\n"
    pDnodeGrpMem.PlayerOut += "&C"
    pDnodeGrpMem.PlayerOut += pDnodeActor.PlayerName
    pDnodeGrpMem.PlayerOut += " says to the group: "
    pDnodeGrpMem.PlayerOut += GsayMsg
    pDnodeGrpMem.PlayerOut += "&N"
    pDnodeGrpMem.PlayerOut += "\r\n"
    CreatePrompt(pDnodeGrpMem.pPlayer)
    pDnodeGrpMem.PlayerOut += GetOutput(pDnodeGrpMem.pPlayer)
  }
}

// Hail command
func DoHail() {
  var pMobile        *Mobile
  var HailMsg         string
  var MobileMsg       string
  var PlayerName      string
  var PlayerNameCheck string
  var RoomId          string
  var Target          string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) < 2 {
    // No target
    pDnodeActor.PlayerOut += "You need a target."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if StrCountWords(CmdStr) > 2 {
    // Two many targets
    pDnodeActor.PlayerOut += "Only one target at a time."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  PlayerName = pDnodeActor.PlayerName
  RoomId = pDnodeActor.pPlayer.RoomId
  PlayerNameCheck = PlayerName
  PlayerNameCheck = StrMakeLower(PlayerNameCheck)
  Target = StrGetWord(CmdStr, 2)
  Target = StrMakeLower(Target)
  if Target == PlayerNameCheck {
    // Trying to kill self
    pDnodeActor.PlayerOut += "Hailing yourself is just plain silly."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pMobile = IsMobInRoom(Target)
  if pMobile == nil {
    // Target mobile is not here
    pDnodeActor.PlayerOut += "Try hailing an NPC that is in this room."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*****************
  //* Send messages *
  //*****************
  // Send message to player
  pDnodeActor.PlayerOut += "&W"
  pDnodeActor.PlayerOut += "You hail "
  pDnodeActor.PlayerOut += pMobile.Desc1
  pDnodeActor.PlayerOut += "!"
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Send message to room
  HailMsg = "&W"
  HailMsg += PlayerName
  HailMsg += " hails "
  HailMsg += pMobile.Desc1
  HailMsg += "."
  HailMsg += "&N"
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(RoomId, HailMsg)
  MobileMsg = MobTalk(pMobile)
  // Stip last \r\n from message, SendToRoom adds a \r\n
  MobileMsg = StrLeft(MobileMsg, StrGetLength(MobileMsg)-2)
  pDnodeSrc = nil
  pDnodeTgt = nil
  SendToRoom(RoomId, MobileMsg)
  pMobile = nil
}

// Help command
func DoHelp() {
  DEBUGIT(1)
  if IsHelp() {
    // Help was found and sent to player
    return
  }
  // No help entry found
  CmdStr = "help notfound"
  if IsHelp() {
    // Help notfound entry was found and sent to player
    return
  }
  // Ok, if we are here, then there really is no help
  pDnodeActor.PlayerOut += "No help topic found.\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Inventory command
func DoInventory() {
  DEBUGIT(1)
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  ShowPlayerInv()
}

// Invisible command
func DoInvisible() {
  DEBUGIT(1)
  TmpStr = StrGetWord(CmdStr, 2)
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr == "on" {
    // Turn Invisible on
    pDnodeActor.pPlayer.Invisible = true
    pDnodeActor.PlayerStateInvisible = true
    PlayerSave(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += "You are invisible.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if TmpStr == "off" {
    // Turn Invisible off
    pDnodeActor.pPlayer.Invisible = false
    pDnodeActor.PlayerStateInvisible = false
    PlayerSave(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += "You are visible.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pDnodeActor.pPlayer.Invisible {
    // Invisible is on
    pDnodeActor.PlayerOut += "Invisible is on.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  } else {
    // Invisible is off
    pDnodeActor.PlayerOut += "Invisible is off.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  }
}

// Kill command
func DoKill() {
  var pMobile        *Mobile
  var KillMsg         string
  var MobileId        string
  var MobileName      string
  var PlayerName      string
  var PlayerNameCheck string
  var RoomId          string
  var Target          string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if pDnodeActor.pPlayer.Position == "sit" {
    // Player is sitting
    pDnodeActor.PlayerOut += "You must stand up before starting a fight!"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  RoomId = pDnodeActor.pPlayer.RoomId
  if IsRoomType(RoomId, "NoFight") {
    // No fighting is allowed in this room
    pDnodeActor.PlayerOut += "You are not allowed to fight here."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pDnodeActor.PlayerStateFighting {
    // Player is fighting
    pDnodeActor.PlayerOut += "You can only fight one opponent at a time!"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if StrCountWords(CmdStr) < 2 {
    // No target
    pDnodeActor.PlayerOut += "You need a target."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if StrCountWords(CmdStr) > 2 {
    // Two many targets
    pDnodeActor.PlayerOut += "Only one target at a time."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  PlayerName = pDnodeActor.PlayerName
  PlayerNameCheck = PlayerName
  PlayerNameCheck = StrMakeLower(PlayerNameCheck)
  Target = StrGetWord(CmdStr, 2)
  MobileName = Target
  Target = StrMakeLower(Target)
  if Target == PlayerNameCheck {
    // Trying to kill self
    pDnodeActor.PlayerOut += "That would be just awful."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if IsPlayer(Target) {
    // Trying to kill another player
    pDnodeActor.PlayerOut += "Don't even think about it, player killing is not allowed."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pMobile = IsMobInRoom(Target)
  if pMobile == nil {
    // Target mobile is not here
    pDnodeActor.PlayerOut += "There doesn't seem to be a(n) "
    pDnodeActor.PlayerOut += MobileName
    pDnodeActor.PlayerOut += " here.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //**************
  //* One Whack! *
  //**************
  if pDnodeActor.pPlayer.OneWhack {
    // Admininstrator has set OneWhack to yes
    pDnodeActor.PlayerStateFighting = false
    pDnodeActor.PlayerOut += "&R"
    pDnodeActor.PlayerOut += "One WHACK and "
    pDnodeActor.PlayerOut += pMobile.Desc1
    pDnodeActor.PlayerOut += " is dead!\r\n"
    pDnodeActor.PlayerOut += "&N"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    // Send kill message to room
    KillMsg = "&R"
    KillMsg += PlayerName
    KillMsg += " kills "
    KillMsg += pMobile.Desc1
    KillMsg += " with one WHACK!"
    KillMsg += "&N"
    pDnodeSrc = pDnodeActor
    pDnodeTgt = pDnodeActor
    SendToRoom(RoomId, KillMsg)
    // Remove mobile from room
    RemoveMobFromRoom(RoomId, pMobile.MobileId)
    pMobile = nil
    return
  }
  //*****************
  //* Send messages *
  //*****************
  // Send message to player
  pDnodeActor.PlayerOut += "&R"
  pDnodeActor.PlayerOut += "You start a fight with "
  pDnodeActor.PlayerOut += pMobile.Desc1
  pDnodeActor.PlayerOut += "!"
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Send message to room
  KillMsg = "&R"
  KillMsg += PlayerName
  KillMsg += " starts a fight with "
  KillMsg += pMobile.Desc1
  KillMsg += "!"
  KillMsg += "&N"
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(RoomId, KillMsg)
  //*****************
  //* Start a fight *
  //*****************
  if !pMobile.Hurt {
    //  Mobile not hurt
    GetNextMobNbr()
    CreateMobStatsFile(RoomId)
    MobileId = pMobile.MobileId
    RemoveMobFromRoom(RoomId, MobileId)
    MobileId = pMobile.MobileId + "." + pMobile.MobNbr
  } else {
    // Mobile is hurt
    MobileId = pMobile.MobileId + "." + pMobile.MobNbr
    RemoveMobFromRoom(RoomId, MobileId)
  }
  UpdateMobInWorld(MobileId, "add") // Keep Mob InWorld count correct
  // Set player and mobile to fight
  CreatePlayerMob(PlayerName, MobileId)
  CreateMobPlayer(PlayerName, MobileId)
  pMobile = nil
  pDnodeActor.PlayerStateFighting = true
}

// List command
func DoList() {
  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if !IsShop(pDnodeActor.pPlayer.RoomId) {
    // Room is not a shop
    pDnodeActor.PlayerOut += "Find a shop."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  ListObjects(pDnodeActor.pPlayer.RoomId, &pDnodeActor.PlayerOut)
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Load command
func DoLoad() {
  var pMobile   *Mobile
  var LoadMsg   string
  var MobileId  string
  var ObjectId  string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if StrCountWords(CmdStr) == 1 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Load what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if StrCountWords(CmdStr) != 3 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Usage: load obj{ect}|mob{ile} <target>"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TmpStr = StrGetWord(CmdStr, 2)
  TmpStr = StrMakeLower(TmpStr)
  TmpStr = TranslateWord(TmpStr)
  if StrIsNotWord(TmpStr, "object mobile") {
    // obj or mob must be specified
    pDnodeActor.PlayerOut += "2nd parm must be obj{ect}|mob{ile}\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //***************
  //* Load object *
  //***************
  if TmpStr == "object" {
    // Loading an object
    ObjectId = StrGetWord(CmdStr, 3)
    pObject = nil
    IsObject(ObjectId) // Sets pObject
    if pObject == nil {
      // Object does not exist
      pDnodeActor.PlayerOut += "Object not found.\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
    pObject = nil
    AddObjToPlayerInv(pDnodeActor, ObjectId)
    pDnodeActor.PlayerOut += "Load successful\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*****************
  //* Load a mobile *
  //*****************
  if TmpStr == "mobile" {
    // Loading an mobile
    MobileId = StrGetWord(CmdStr, 3)
    MobileId = StrMakeLower(MobileId)
    pMobile = IsMobValid(MobileId)
    if pMobile == nil {
      // Mobile does not exist
      pDnodeActor.PlayerOut += "Mobile not found.\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
    AddMobToRoom(pDnodeActor.pPlayer.RoomId, MobileId)
    SpawnMobileNoMove(MobileId)
    pDnodeActor.PlayerOut += "Load successful\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    LoadMsg = pMobile.Desc1
    LoadMsg += " suddenly appears!"
    pDnodeSrc = pDnodeActor
    pDnodeTgt = pDnodeActor
    SendToRoom(pDnodeActor.pPlayer.RoomId, LoadMsg)
    pMobile = nil
    return
  }
}

// Logon
func DoLogon() {
  DEBUGIT(1)
  if pDnodeActor.PlayerStateWaitNewCharacter {
    // New character 'y-n' prompt
    pDnodeActor.PlayerStateWaitNewCharacter = false
    LogonWaitNewCharacter()
    return
  }
  if pDnodeActor.PlayerStateWaitName {
    // Name prompt
    pDnodeActor.PlayerStateWaitName = false
    LogonWaitName()
    return
  }
  if pDnodeActor.PlayerStateWaitNameConfirmation {
    // Name confirmation prompt
    pDnodeActor.PlayerStateWaitNameConfirmation = false
    LogonWaitNameConfirmation()
    return
  }
  if pDnodeActor.PlayerStateWaitPassword {
    // Password prompt
    pDnodeActor.PlayerStateWaitPassword = false
    LogonWaitPassword()
    return
  }
  if pDnodeActor.PlayerStateWaitMaleFemale {
    // Sex prompt
    pDnodeActor.PlayerStateWaitMaleFemale = false
    LogonWaitMaleFemale()
    return
  }
}

// Look command
func DoLook(CmdStr1 string) {
  var pMobile         *Mobile
  var IsPlayer         bool
  var MudCmdIsExit     string
  var TargetName       string

  DEBUGIT(1)
  CmdStr = CmdStr1
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  TmpStr = StrGetWord(CmdStr, 2)
  //*****************
  //* Just looking? *
  //*****************
  if TmpStr == "" {
    // Just look
    ShowRoom(pDnodeActor)
    return
  }
  //**********************
  //* Is it a room exit? *
  //**********************
  MudCmdIsExit = "look"
  if IsExit(MudCmdIsExit) {
    // Look room exit
    return
  }
  //*******************
  //* Is it a player? *
  //*******************
  IsPlayer = true
  TargetName = TmpStr
  TargetName = StrMakeLower(TargetName)
  pDnodeTgt = GetTargetDnode(TargetName)
  if pDnodeTgt == nil {
    // Target is not online and/or not in 'playing' state
    IsPlayer = false
  } else {
    // Target is online and playing
    if pDnodeActor.pPlayer.RoomId != pDnodeTgt.pPlayer.RoomId {
      // Target is not in the same room
      IsPlayer = false
    }
  }
  if IsPlayer {
    // Show player
    ShowPlayerEqu(pDnodeTgt)
    return
  }
  //*******************
  //* Is it a mobile? *
  //*******************
  pMobile = IsMobInRoom(TargetName)
  if pMobile != nil {
    // Player is looking at a mob
    TmpStr = StrMakeFirstLower(pMobile.Desc1)
    pDnodeActor.PlayerOut += "You look at "
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "."
    pDnodeActor.PlayerOut += "\r\n"
    ExamineMob(pMobile.MobileId)
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    pMobile = nil
    return
  }
  // Nothing found to look at
  pDnodeActor.PlayerOut += "If it's an object, use examine, otherwise <shrug>."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Money command
func DoMoney() {
  DEBUGIT(1)
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  ShowMoney(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Motd command
func DoMotd() {
  var MotdFile     *os.File
  var MotdFileName string

  DEBUGIT(1)
  // Read Motd file
  MotdFileName = MOTD_DIR
  MotdFileName += "Motd"
  MotdFileName += ".txt"
  MotdFile, err := os.Open(MotdFileName)
  if err != nil {
    LogBuf = "Communication::DoMotd - Open Motd file failed (read)"
    LogIt(LogBuf)
    return
  }
  Scanner := bufio.NewScanner(MotdFile)
  for Scanner.Scan() {
    Stuff = Scanner.Text()
    if Stuff == "End of Motd" {
      break
    }
    Stuff += "\r\n"
    pDnodeActor.PlayerOut += Stuff
  }
  MotdFile.Close()
  if pDnodeActor.PlayerStatePlaying {
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  }
}

// OneWhack command
func DoOneWhack() {
  DEBUGIT(1)
  TmpStr = StrGetWord(CmdStr, 2)
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr == "on" {
    // Turn OneWhack on
    pDnodeActor.pPlayer.OneWhack = true
    PlayerSave(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += "You will obliterate enemies with One Whack.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if TmpStr == "off" {
    // Turn OneWhack off
    pDnodeActor.pPlayer.OneWhack = false
    PlayerSave(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += "You will now have to fight for you life.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pDnodeActor.pPlayer.OneWhack {
    // OneWhack is on
    pDnodeActor.PlayerOut += "OneWhack is on.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  } else {
    // OneWhack is off
    pDnodeActor.PlayerOut += "OneWhack is off.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  }
}

// Password command
func DoPassword() {
  var Password     string
  var NewPassword1 string
  var NewPassword2 string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) != 4 {
    pDnodeActor.PlayerOut += "Password command requires: "
    pDnodeActor.PlayerOut += "Password NewPassword NewPassword"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  Password = StrGetWord(CmdStr, 2)
  NewPassword1 = StrGetWord(CmdStr, 3)
  NewPassword2 = StrGetWord(CmdStr, 4)
  if Password != pDnodeActor.pPlayer.Password {
    pDnodeActor.PlayerOut += "Password does not match current password."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if NewPassword1 != NewPassword2 {
    pDnodeActor.PlayerOut += "New passwords do not match."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*******************
  //* Change password *
  //*******************
  pDnodeActor.pPlayer.Password = NewPassword1
  PlayerSave(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += "Your password has been changed."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Played command
func DoPlayed() {
  var Days          int64
  var Hours         int64
  var Minutes       int64
  var Seconds       int64
  var n             int64
  var BirthDay      string
  var PlayerAge     string
  var TimePlayed    string
  var BornSec       int64
  var NowSec        int64
  var PlayerAgeSec  int64
  var TimePlayedSec int64

  DEBUGIT(1)
  PlayerSave(pDnodeActor.pPlayer) // Save() updates TimePlayed
  NowSec = time.Now().Unix()
  BornSec = pDnodeActor.pPlayer.Born
  PlayerAgeSec = NowSec - BornSec
  TimePlayedSec = pDnodeActor.pPlayer.TimePlayed
  // Birthday
  BirthDay = time.Unix(BornSec, 0).Format(time.ANSIC)
  // Age
  n = PlayerAgeSec
  Days = n / (24 * 3600)
  n = n % (24 * 3600)
  Hours = n / 3600
  n %= 3600
  Minutes = n / 60
  n %= 60
  Seconds = n
  PlayerAge = fmt.Sprintf("Your age: %d days, %d hours, %d minutes, %d seconds", Days, Hours, Minutes, Seconds)
  // TimePlayed
  n = TimePlayedSec
  Days = n / (24 * 3600)
  n = n % (24 * 3600)
  Hours = n / 3600
  n %= 3600
  Minutes = n / 60
  n %= 60
  Seconds = n
  TimePlayed = fmt.Sprintf("You've played: %d days, %d hours, %d minutes, %d seconds", Days, Hours, Minutes, Seconds)
  pDnodeActor.PlayerOut += PlayerAge
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += TimePlayed
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "Your birthday is: "
  pDnodeActor.PlayerOut += BirthDay
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Quit command
func DoQuit() {
  var AllMsg    string
  var PlayerMsg string

  DEBUGIT(1)
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  pDnodeActor.PlayerStateBye     = true
  pDnodeActor.PlayerStatePlaying = false
  PlayerSave(pDnodeActor.pPlayer)
  GrpLeave()
  LogBuf = pDnodeActor.PlayerName
  LogBuf += " issued the QUIT command"
  LogIt(LogBuf)
  PlayerMsg = "\r\n"
  PlayerMsg += "Bye Bye!"
  PlayerMsg += "\r\n"
  AllMsg = "\r\n"
  AllMsg += pDnodeActor.PlayerName
  AllMsg += " has left the game."
  AllMsg += "\r\n"
  SendToAll(PlayerMsg, AllMsg)
}

// Refresh command
func DoRefresh() {
  DEBUGIT(1)
  if StrCountWords(CmdStr) > 2 {
    // Invalid command format
    pDnodeActor.PlayerOut += "You may refresh only one thing at time."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TmpStr = StrGetWord(CmdStr, 2)
  if StrGetLength(TmpStr) == 0 {
    // Player did not provide an target to be refreshed
    pDnodeActor.PlayerOut += "Refresh what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr == "commands" {
    // Refresh commands array
    CommandArrayLoad()
    pDnodeActor.PlayerOut += "Commands have been refreshed."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  // Refresh target not valid
  pDnodeActor.PlayerOut += "Refresh what??"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Remove command
func DoRemove() {
  var ObjectName string
  var RemoveMsg  string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) > 2 {
    // Invalid command format, like 'remove shirt pants'
    pDnodeActor.PlayerOut += "You may remove only one item at time."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TmpStr = StrGetWord(CmdStr, 2)
  if StrGetLength(TmpStr) == 0 {
    // Player did not provide an object to be removed
    pDnodeActor.PlayerOut += "Remove what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  // Get pointer to object
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  pObject = nil
  IsObjInPlayerEqu(TmpStr) // Sets pObject
  if pObject == nil {
    // Object not in equipment
    pDnodeActor.PlayerOut += "You don't have a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += " equipped.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  // Decrease player's ArmorClass
  pDnodeActor.pPlayer.ArmorClass -= pObject.ArmorValue
  // Remove object from player's equipment
  RemoveObjFromPlayerEqu(pObject.ObjectId)
  // Send remove message to player
  pDnodeActor.PlayerOut += "You remove "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += ".\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Send remove message to room
  RemoveMsg = pDnodeActor.PlayerName
  RemoveMsg += " removes "
  RemoveMsg += pObject.Desc1
  RemoveMsg += "."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, RemoveMsg)
  // Add object to player's inventory
  AddObjToPlayerInv(pDnodeTgt, pObject.ObjectId)
  TmpStr = pObject.Type
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr == "weapon" {
    // Now player has no weapon
    pDnodeActor.pPlayer.WeaponDamage = PLAYER_DMG_HAND
    pDnodeActor.pPlayer.WeaponDesc1 = "a pair of bare hands"
    pDnodeActor.pPlayer.WeaponType = "Hand"
  }
  pObject = nil
}

// Restore command
func DoRestore(CmdStr string) {
  var PlayerName     string
  var TargetName     string
  var TargetNameSave string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  PlayerName = pDnodeActor.PlayerName
  TargetName = StrGetWord(CmdStr, 2)
  TargetNameSave = TargetName
  PlayerName = StrMakeLower(PlayerName)
  TargetName = StrMakeLower(TargetName)
  if StrGetLength(TargetName) < 1 {
    // No target, assume self
    TargetName = PlayerName
  }
  if TargetName == PlayerName {
    // Admin is restore themself
    pDnodeActor.PlayerOut += "You restore yourself!\r\n"
    pDnodeActor.pPlayer.HitPoints = pDnodeActor.pPlayer.Level * PLAYER_HPT_PER_LEVEL
    pDnodeActor.pPlayer.Hunger = 0
    pDnodeActor.pPlayer.Thirst = 0
    PlayerSave(pDnodeActor.pPlayer)
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pDnodeTgt = GetTargetDnode(TargetName)
  if pDnodeTgt == nil {
    // Tell player ... not found
    pDnodeActor.PlayerOut += TargetNameSave
    pDnodeActor.PlayerOut += " is not online.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  // Restore the player
  pDnodeTgt.pPlayer.HitPoints = pDnodeTgt.pPlayer.Level * PLAYER_HPT_PER_LEVEL
  pDnodeTgt.pPlayer.Hunger = 0
  pDnodeTgt.pPlayer.Thirst = 0
  PlayerSave(pDnodeTgt.pPlayer)
  //****************************
  //* Send the restore message *
  //****************************
  PlayerName = pDnodeActor.PlayerName
  TargetName = pDnodeTgt.PlayerName
  // Send restore message to player
  pDnodeActor.PlayerOut += "You have restored "
  pDnodeActor.PlayerOut += TargetName
  pDnodeActor.PlayerOut += "."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Send restore message to target
  pDnodeTgt.PlayerOut += "\r\n"
  pDnodeTgt.PlayerOut += "&Y"
  pDnodeTgt.PlayerOut += PlayerName
  pDnodeTgt.PlayerOut += " has restored you!"
  pDnodeTgt.PlayerOut += "&N"
  pDnodeTgt.PlayerOut += "\r\n"
  CreatePrompt(pDnodeTgt.pPlayer)
  pDnodeTgt.PlayerOut += GetOutput(pDnodeTgt.pPlayer)
}

// RoomInfo command
func DoRoomInfo() {
  DEBUGIT(1)
  TmpStr = StrGetWord(CmdStr, 2)
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr == "on" {
    pDnodeActor.pPlayer.RoomInfo = true
    pDnodeActor.PlayerOut += "You will now see hidden room info.\r\n"
    PlayerSave(pDnodeActor.pPlayer)
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if TmpStr == "off" {
    pDnodeActor.pPlayer.RoomInfo = false
    pDnodeActor.PlayerOut += "You will no longer see hidden room info.\r\n"
    PlayerSave(pDnodeActor.pPlayer)
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pDnodeActor.PlayerOut += "Try on or off.\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Save command
func DoSave() {
  DEBUGIT(1)
  PlayerSave(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += "Saved!"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Say command
func DoSay() {
  var SayMsg string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  SayMsg = StrGetWords(CmdStr, 2)
  if StrGetLength(SayMsg) < 1 {
    // Player did not enter anything to say
    pDnodeActor.PlayerOut += "You try to speak, but no words come out of your mouth."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if pDnodeActor.PlayerStateInvisible {
    // Player can't speak while invisible
    pDnodeActor.PlayerOut += "No talking while you are invisible."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*****************
  //* Say something *
  //*****************
  pDnodeActor.PlayerOut += "&W"
  pDnodeActor.PlayerOut += "You say: "
  pDnodeActor.PlayerOut += SayMsg
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  TmpStr = "&W"
  TmpStr += pDnodeActor.PlayerName
  TmpStr += " says: "
  TmpStr += SayMsg
  TmpStr += "&N"
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, TmpStr)
}

// Sell command
func DoSell() {
  var Cost         int
  var Desc1        string
  var InvCountInt  int
  var InvCountStr  string
  var ObjectId     string
  var ObjectName   string
  var RoomId       string
  var SellCountInt int
  var SellCountStr string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  RoomId = pDnodeActor.pPlayer.RoomId
  if !IsShop(RoomId) {
    // Room is not a shop
    pDnodeActor.PlayerOut += "Find a shop."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  ObjectName = StrGetWord(CmdStr, 2)
  if ObjectName == "" {
    // No object given
    pDnodeActor.PlayerOut += "Sell what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pObject = nil
  IsObjInPlayerInv(ObjectName) // Sets pObject
  if pObject == nil {
    // Player doesn't have object
    pDnodeActor.PlayerOut += "You don't have a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += " in your inventory."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  ObjectId = pObject.ObjectId
  Desc1 = pObject.Desc1
  Cost = pObject.Cost
  InvCountStr = pObject.Count
  pObject = nil
  IsShopObj(RoomId, ObjectName) // Sets pObject
  if pObject == nil {
    // Player cannot sell object to this shop
    pDnodeActor.PlayerOut += "That item cannot be sold here."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pObject = nil
  //********************
  //* Check sell count *
  //********************
  InvCountInt  = StrToInt(InvCountStr)
  SellCountStr = StrGetWord(CmdStr, 3)
  SellCountStr = StrMakeLower(SellCountStr) // In case player typed 'all'
  if SellCountStr == "" {
    // Player did not specify a sell count
    SellCountInt = 1
  } else {
    // Player might be selling more than 1 item
    if SellCountStr == "all" {
      // Player is selling all items
      SellCountInt = InvCountInt
    } else {
      // Compare player InvCount against SellCountInt
      SellCountInt = StrToInt(SellCountStr)
      if SellCountInt <= 0 {
        // Player entered a negative or zero amount to sell
        pDnodeActor.PlayerOut += "You cannot sell less that 1 item"
        pDnodeActor.PlayerOut += "\r\n"
        CreatePrompt(pDnodeActor.pPlayer)
        pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
        return
      } else if SellCountInt > InvCountInt {
        // Player is trying sell more than they have
        pDnodeActor.PlayerOut += "You don't have that many "
        pDnodeActor.PlayerOut += ObjectName
        pDnodeActor.PlayerOut += "s"
        pDnodeActor.PlayerOut += "\r\n"
        CreatePrompt(pDnodeActor.pPlayer)
        pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
        return
      }
    }
  }
  //*******************
  //* Sell the object *
  //*******************
  // Remove object from player's inventory
  RemoveObjFromPlayerInv(ObjectId, SellCountInt)
  // Player receives some money
  Cost = Cost * SellCountInt
  SetMoney(pDnodeActor.pPlayer,'+', Cost, "Silver")
  // Send messages
  Buf = fmt.Sprintf("%d", SellCountInt)
  TmpStr = Buf
  TmpStr = "(" + TmpStr + ") "
  pDnodeActor.PlayerOut += TmpStr
  pDnodeActor.PlayerOut += "You sell "
  pDnodeActor.PlayerOut += Desc1
  pDnodeActor.PlayerOut += " for "
  Buf = fmt.Sprintf("%d", Cost)
  TmpStr = Buf
  pDnodeActor.PlayerOut += TmpStr
  pDnodeActor.PlayerOut += " piece(s) of silver."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Show command
func DoShow() {
  var CommandCheckResult string
  var HelpFileName       string
  var HelpText           string
  var SocialFileName     string
  var SocialText         string
  var MudCmdChk          string
  var ValCmdInfo         string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if StrCountWords(CmdStr) > 2 {
    // Invalid command format
    pDnodeActor.PlayerOut += "You may show only one thing at time."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TmpStr = StrGetWord(CmdStr, 2)
  if StrGetLength(TmpStr) == 0 {
    // Player did not provide a target to be shown
    pDnodeActor.PlayerOut += "Show what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TmpStr = StrMakeLower(TmpStr)
  if StrIsNotWord(TmpStr, "commands socials help") {
    // Show target not valid
    pDnodeActor.PlayerOut += "Show what??"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*****************
  //* Show commands *
  //*****************
  if TmpStr == "commands" {
    pDnodeActor.PlayerOut += "Commands"
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "--------"
    pDnodeActor.PlayerOut += "\r\n"
    for _, ValidCmd := range ValidCmds {
      // For each string in the ValidCmds slice
      ValCmdInfo = ValidCmd
      MudCmdChk = StrGetWord(ValCmdInfo, 1)
      CommandCheckResult = CommandCheck(MudCmdChk)
      if CommandCheckResult == "Ok" {
        // Mud command is Ok for this player
        pDnodeActor.PlayerOut += MudCmdChk
        pDnodeActor.PlayerOut += "\r\n"
      } else if StrGetWord(CommandCheckResult, 1) == "Level" {
        // Mud command is Ok for this player, but level restricted
        pDnodeActor.PlayerOut += MudCmdChk
        pDnodeActor.PlayerOut += " acquired at level("
        pDnodeActor.PlayerOut += StrGetWord(CommandCheckResult, 2)
        pDnodeActor.PlayerOut += ")"
        pDnodeActor.PlayerOut += "\r\n"
      }
    }
  }
  //*************
  //* Show help *
  //*************
  if TmpStr == "help" {
    pDnodeActor.PlayerOut += "Help Topics"
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "-----------"
    pDnodeActor.PlayerOut += "\r\n"
    HelpFileName = HELP_DIR
    HelpFileName += "Help.txt"
    HelpFile, err := os.Open(HelpFileName)
    if err != nil {
      // Help file open failed
      pDnodeActor.PlayerOut += "No help is available, you are on your own!"
      pDnodeActor.PlayerOut += "\r\n"
    } else {
      // Help file is open
      Scanner := bufio.NewScanner(HelpFile)
      if Scanner.Scan() {
        HelpText = Scanner.Text() // Skip first line
      }
      if Scanner.Scan() {
        HelpText = Scanner.Text()
      }
      for HelpText != "End of Help" {
        // Read the whole file
        if StrLeft(HelpText, 5) == "Help:" {
          // Found a help topic
          pDnodeActor.PlayerOut += StrRight(HelpText, StrGetLength(HelpText)-5)
          pDnodeActor.PlayerOut += "\r\n"
        }
        if !Scanner.Scan() {
          break
        }
        HelpText = Scanner.Text()
      }
      HelpFile.Close()
    }
  }
  //****************
  //* Show socials *
  //****************
  if TmpStr == "socials" {
    pDnodeActor.PlayerOut += "Socials"
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "-------"
    pDnodeActor.PlayerOut += "\r\n"
    SocialFileName = SOCIAL_DIR
    SocialFileName += "Social.txt"
    SocialFile, err := os.Open(SocialFileName)
    if err != nil {
      // Social file open failed
      pDnodeActor.PlayerOut += "No socials are available, how boring!"
      pDnodeActor.PlayerOut += "\r\n"
    } else {
      // Social file is open
      Scanner := bufio.NewScanner(SocialFile)
      for Scanner.Scan() {
        SocialText = Scanner.Text()
        if SocialText == "End of Socials" {
          break
        }
        if StrLeft(SocialText, 9) == "Social : " {
          // Found a help topic
          pDnodeActor.PlayerOut += StrRight(SocialText, StrGetLength(SocialText)-9)
          pDnodeActor.PlayerOut += "\r\n"
        }
      }
      SocialFile.Close()
    }
  }
  // Prompt
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Sit command
func DoSit() {
  var SitMsg string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if pDnodeActor.pPlayer.Position == "sit" {
    // Player is already sitting
    pDnodeActor.PlayerOut += "You are already sitting down."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //************
  //* Sit down *
  //************
  pDnodeActor.pPlayer.Position = "sit"
  PlayerSave(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += "You sit down."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  SitMsg = pDnodeActor.PlayerName + " sits down."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, SitMsg)
}

// Sleep command
func DoSleep() {
  var SleepMsg string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if pDnodeActor.pPlayer.Position != "sit" {
    // Player must sitting before sleeping
    pDnodeActor.PlayerOut += "You must be sitting before you can sleep."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //***************
  //* Fall asleep *
  //***************
  pDnodeActor.pPlayer.Position = "sleep"
  PlayerSave(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += "You go to sleep."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  SleepMsg = pDnodeActor.PlayerName + " goes to sleep."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, SleepMsg)
}

// Stand command
func DoStand() {
  var StandMsg string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  if pDnodeActor.pPlayer.Position == "stand" {
    // Player is already standing
    pDnodeActor.PlayerOut += "You are already standing."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //************
  //* Stand up *
  //************
  pDnodeActor.pPlayer.Position = "stand"
  PlayerSave(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += "You stand up."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  StandMsg = pDnodeActor.PlayerName + " stands up."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, StandMsg)
}

// Status command
func DoStatus() {
  DEBUGIT(1)
  ShowStatus(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Stop command
func DoStop() {
  var GoGoGoFileName string
  var StopItFileName string

  DEBUGIT(1)
  StateStopping = true
  LogBuf = pDnodeActor.PlayerName
  LogBuf += " issued the STOP command"
  LogIt(LogBuf)
  pDnodeActor.PlayerOut += "Stop command issued!\r\n"
  GoGoGoFileName = CONTROL_DIR
  GoGoGoFileName += "GoGoGo"
  StopItFileName = CONTROL_DIR
  StopItFileName += "StopIt"
  Rename(GoGoGoFileName, StopItFileName)
}

// Tell command
func DoTell() {
  var PlayerName     string
  var TargetName     string
  var TargetNameSave string
  var TellMsg        string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  PlayerName = pDnodeActor.PlayerName
  TargetName = StrGetWord(CmdStr, 2)
  TargetNameSave = TargetName
  PlayerName = StrMakeLower(PlayerName)
  TargetName = StrMakeLower(TargetName)
  if TargetName == PlayerName {
    pDnodeActor.PlayerOut += "Seems silly to tell yourself something!\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if StrGetLength(TargetName) < 1 {
    pDnodeActor.PlayerOut += "Tell who?\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  TellMsg = StrGetWords(CmdStr, 3)
  if StrGetLength(TellMsg) < 1 {
    pDnodeActor.PlayerOut += "Um, tell "
    pDnodeActor.PlayerOut += TargetNameSave
    pDnodeActor.PlayerOut += " what?"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pDnodeTgt = GetTargetDnode(TargetName)
  if pDnodeTgt == nil {
    // Tell player ... not found
    pDnodeActor.PlayerOut += TargetNameSave
    pDnodeActor.PlayerOut += " is not online.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*****************
  //* Send the tell *
  //*****************
  PlayerName = pDnodeActor.PlayerName
  TargetName = pDnodeTgt.PlayerName
  // Send tell message to player
  pDnodeActor.PlayerOut += "&M"
  pDnodeActor.PlayerOut += "You tell "
  pDnodeActor.PlayerOut += TargetName
  pDnodeActor.PlayerOut += ": "
  pDnodeActor.PlayerOut += TellMsg
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Send tell message to target
  pDnodeTgt.PlayerOut += "&M"
  pDnodeTgt.PlayerOut += "\r\n"
  pDnodeTgt.PlayerOut += PlayerName
  pDnodeTgt.PlayerOut += " tells you: "
  pDnodeTgt.PlayerOut += TellMsg
  pDnodeTgt.PlayerOut += "&N"
  pDnodeTgt.PlayerOut += "\r\n"
  CreatePrompt(pDnodeTgt.pPlayer)
  pDnodeTgt.PlayerOut += GetOutput(pDnodeTgt.pPlayer)
}

// Time command
func DoTime() {
  var DisplayCurrentTime string

  DEBUGIT(1)
  // Server time
  DisplayCurrentTime = time.Now().Format("2006-01-02 15:04:05")
  pDnodeActor.PlayerOut += "Current server time is: "
  pDnodeActor.PlayerOut += DisplayCurrentTime
  pDnodeActor.PlayerOut += "\r\n"
  // Game time
  pDnodeActor.PlayerOut += "Current game time is: "
  pDnodeActor.PlayerOut += GetTime()
  pDnodeActor.PlayerOut += "\r\n"
  // Prompt
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Title command
func DoTitle() {
  var Title string

  DEBUGIT(1)
  TmpStr = StrGetWord(CmdStr, 2)
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr == "" {
    // Player entered 'title' by itself
    if pDnodeActor.pPlayer.Title == "" {
      // Player has no title
      pDnodeActor.PlayerOut += "You do not have a title"
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    } else {
      // Show player's title
      pDnodeActor.PlayerOut += "Your title is: "
      pDnodeActor.PlayerOut += pDnodeActor.pPlayer.Title
      pDnodeActor.PlayerOut += "&N" // In case title is messed up
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  if TmpStr == "none" {
    // Player entered 'title none'
    if pDnodeActor.pPlayer.Title == "" {
      // Player has no title
      pDnodeActor.PlayerOut += "You did not have a title and you still do not have a title"
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    } else {
      pDnodeActor.pPlayer.Title = ""
      pDnodeActor.PlayerOut += "Your title has been removed.\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  TmpStr = StrGetWords(CmdStr, 2)
  Title = TmpStr
  // Strip out color codes so Title length can be checked
  StrReplace(&TmpStr, "&N", "")
  StrReplace(&TmpStr, "&K", "")
  StrReplace(&TmpStr, "&R", "")
  StrReplace(&TmpStr, "&G", "")
  StrReplace(&TmpStr, "&Y", "")
  StrReplace(&TmpStr, "&B", "")
  StrReplace(&TmpStr, "&M", "")
  StrReplace(&TmpStr, "&C", "")
  StrReplace(&TmpStr, "&W", "")
  if StrGetLength(TmpStr) > 40 {
    pDnodeActor.PlayerOut += "Title must be less than 41 characters, color codes do not count.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pDnodeActor.pPlayer.Title = Title
  PlayerSave(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += "Your title has been set.\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Train command
func DoTrain() {
  var IncreaseDecrease     int
  var MinusSign            string
  var SkillPointsUsed      int
  var SkillPointsRemaining int
  var UnTrainCost          string
  var WeaponType           string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if IsFighting() {
    // Player is fighting, send msg, command is not done
    return
  }
  // Get command words
  WeaponType = StrGetWord(CmdStr, 2)
  MinusSign = StrGetWord(CmdStr, 3)
  UnTrainCost = StrGetWord(CmdStr, 4)
  WeaponType = StrMakeLower(WeaponType)
  // Calculate skill points used and remaining
  SkillPointsUsed = 0
  SkillPointsUsed += pDnodeActor.pPlayer.SkillAxe
  SkillPointsUsed += pDnodeActor.pPlayer.SkillClub
  SkillPointsUsed += pDnodeActor.pPlayer.SkillDagger
  SkillPointsUsed += pDnodeActor.pPlayer.SkillHammer
  SkillPointsUsed += pDnodeActor.pPlayer.SkillSpear
  SkillPointsUsed += pDnodeActor.pPlayer.SkillStaff
  SkillPointsUsed += pDnodeActor.pPlayer.SkillSword
  SkillPointsRemaining = PLAYER_SKILL_PER_LEVEL*pDnodeActor.pPlayer.Level - SkillPointsUsed
  // Do some more checking
  if StrCountWords(CmdStr) > 4 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Train command syntax error, try'er again."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if WeaponType != "" {
    // WeaponType specified
    if StrIsNotWord(WeaponType, "axe club dagger hammer spear staff sword") {
      // But it was the invalid
      pDnodeActor.PlayerOut += "Please specify a valid weapon type."
      pDnodeActor.PlayerOut += "\r\n"
      return
    } else {
      // Player is trying to train, check skill points remaining
      if SkillPointsRemaining == 0 && MinusSign == "" {
        // No skill points available
        pDnodeActor.PlayerOut += "Sorry, you have no skill points remaining."
        pDnodeActor.PlayerOut += "\r\n"
        CreatePrompt(pDnodeActor.pPlayer)
        pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
        return
      }
    }
  }
  if MinusSign != "" {
    // Third word was specified
    if MinusSign != "-" {
      // But it was not a minus sign
      pDnodeActor.PlayerOut += "If third parameter is specified, it must be a minus sign."
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    }
  }
  if MinusSign == "-" {
    // Player wishes to untrain a skill
    if UnTrainCost == "" {
      // But did not specify the cost
      pDnodeActor.PlayerOut += "Untraining a skill will cost "
      pDnodeActor.PlayerOut += UNTRAIN_COST
      pDnodeActor.PlayerOut += " silver"
      pDnodeActor.PlayerOut += "\r\n"
      pDnodeActor.PlayerOut += "Specify this amount after the minus sign."
      pDnodeActor.PlayerOut += "\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      return
    } else {
      // Cost specified, but was it the correct cost?
      if UnTrainCost != UNTRAIN_COST {
        // Wrong cost
        pDnodeActor.PlayerOut += "You must specify the correct amount! "
        pDnodeActor.PlayerOut += UNTRAIN_COST
        pDnodeActor.PlayerOut += " silver please."
        pDnodeActor.PlayerOut += "\r\n"
        CreatePrompt(pDnodeActor.pPlayer)
        pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
        return
      } else {
        // Player has specified everything correctly for untrain, check money
        if pDnodeActor.pPlayer.Silver < StrToInt(UNTRAIN_COST) {
          // Not enough money
          pDnodeActor.PlayerOut += "You do not have "
          pDnodeActor.PlayerOut += UNTRAIN_COST
          pDnodeActor.PlayerOut += " pieces of silver!"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          return
        }
        // Check skills
        if WeaponType == "axe" && pDnodeActor.pPlayer.SkillAxe == 0 {
          // No axe skill
          pDnodeActor.PlayerOut += "You have no "
          pDnodeActor.PlayerOut += WeaponType
          pDnodeActor.PlayerOut += " skill!"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          return
        }
        if WeaponType == "club" && pDnodeActor.pPlayer.SkillClub == 0 {
          // No club skill
          pDnodeActor.PlayerOut += "You have no "
          pDnodeActor.PlayerOut += WeaponType
          pDnodeActor.PlayerOut += " skill!"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          return
        }
        if WeaponType == "dagger" && pDnodeActor.pPlayer.SkillDagger == 0 {
          // No dagger skill
          pDnodeActor.PlayerOut += "You have no "
          pDnodeActor.PlayerOut += WeaponType
          pDnodeActor.PlayerOut += " skill!"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          return
        }
        if WeaponType == "hammer" && pDnodeActor.pPlayer.SkillHammer == 0 {
          // No hammer skill
          pDnodeActor.PlayerOut += "You have no "
          pDnodeActor.PlayerOut += WeaponType
          pDnodeActor.PlayerOut += " skill!"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          return
        }
        if WeaponType == "spear" && pDnodeActor.pPlayer.SkillSpear == 0 {
          // No spear skill
          pDnodeActor.PlayerOut += "You have no "
          pDnodeActor.PlayerOut += WeaponType
          pDnodeActor.PlayerOut += " skill!"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          return
        }
        if WeaponType == "staff" && pDnodeActor.pPlayer.SkillStaff == 0 {
          // No staff skill
          pDnodeActor.PlayerOut += "You have no "
          pDnodeActor.PlayerOut += WeaponType
          pDnodeActor.PlayerOut += " skill!"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          return
        }
        if WeaponType == "sword" && pDnodeActor.pPlayer.SkillSword == 0 {
          // No sword skill
          pDnodeActor.PlayerOut += "You have no "
          pDnodeActor.PlayerOut += WeaponType
          pDnodeActor.PlayerOut += " skill!"
          pDnodeActor.PlayerOut += "\r\n"
          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
          return
        }
      }
    }
  }
  if WeaponType == "" {
    // Show player's skills summary
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += "Skills Summary"
    pDnodeActor.PlayerOut += "\r\n"
    // Axe
    pDnodeActor.PlayerOut += "Axe:    "
    Buf = fmt.Sprintf("%3d", pDnodeActor.pPlayer.SkillAxe)
    TmpStr = Buf
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "\r\n"
    // Club
    pDnodeActor.PlayerOut += "Club:   "
    Buf = fmt.Sprintf("%3d", pDnodeActor.pPlayer.SkillClub)
    TmpStr = Buf
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "\r\n"
    // Dagger
    pDnodeActor.PlayerOut += "Dagger: "
    Buf = fmt.Sprintf("%3d", pDnodeActor.pPlayer.SkillDagger)
    TmpStr = Buf
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "\r\n"
    // Hammer
    pDnodeActor.PlayerOut += "Hammer: "
    Buf = fmt.Sprintf("%3d", pDnodeActor.pPlayer.SkillHammer)
    TmpStr = Buf
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "\r\n"
    // Spear
    pDnodeActor.PlayerOut += "Spear:  "
    Buf = fmt.Sprintf("%3d", pDnodeActor.pPlayer.SkillSpear)
    TmpStr = Buf
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "\r\n"
    // Staff
    pDnodeActor.PlayerOut += "Staff:  "
    Buf = fmt.Sprintf("%3d", pDnodeActor.pPlayer.SkillStaff)
    TmpStr = Buf
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "\r\n"
    // Sword
    pDnodeActor.PlayerOut += "Sword:  "
    Buf = fmt.Sprintf("%3d", pDnodeActor.pPlayer.SkillSword)
    TmpStr = Buf
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "\r\n"
    // Skill points used
    pDnodeActor.PlayerOut += "Skill points used:      "
    Buf = fmt.Sprintf("%4d", SkillPointsUsed)
    TmpStr = Buf
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "\r\n"
    // Skill points remaining
    pDnodeActor.PlayerOut += "Skill points remaining: "
    Buf = fmt.Sprintf("%4d", SkillPointsRemaining)
    TmpStr = Buf
    pDnodeActor.PlayerOut += TmpStr
    pDnodeActor.PlayerOut += "\r\n"
    // Prompt
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  if MinusSign != "-" {
    // Training a skill
    IncreaseDecrease = +1
  } else {
    // UnTraining a skill
    IncreaseDecrease = -1
  }
  // Ok, so train or untrain them already
  if WeaponType == "axe" {
    pDnodeActor.pPlayer.SkillAxe += IncreaseDecrease
  } else if WeaponType == "club" {
    pDnodeActor.pPlayer.SkillClub += IncreaseDecrease
  } else if WeaponType == "dagger" {
    pDnodeActor.pPlayer.SkillDagger += IncreaseDecrease
  } else if WeaponType == "hammer" {
    pDnodeActor.pPlayer.SkillHammer += IncreaseDecrease
  } else if WeaponType == "spear" {
    pDnodeActor.pPlayer.SkillSpear += IncreaseDecrease
  } else if WeaponType == "staff" {
    pDnodeActor.pPlayer.SkillStaff += IncreaseDecrease
  } else if WeaponType == "sword" {
    pDnodeActor.pPlayer.SkillSword += IncreaseDecrease
  }
  if MinusSign != "-" {
    // Training
    pDnodeActor.PlayerOut += "You have improved your "
    pDnodeActor.PlayerOut += WeaponType
    pDnodeActor.PlayerOut += " skill!"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  } else {
    // UnTraining
    pDnodeActor.PlayerOut += "Your "
    pDnodeActor.PlayerOut += WeaponType
    pDnodeActor.PlayerOut += " skill has decreased at a cost of "
    pDnodeActor.PlayerOut += UNTRAIN_COST
    pDnodeActor.PlayerOut += " silver."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    pDnodeActor.pPlayer.Silver -= StrToInt(UNTRAIN_COST)
  }
  PlayerSave(pDnodeActor.pPlayer)
}

// Wake command
func DoWake() {
  var WakeMsg string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if pDnodeActor.pPlayer.Position != "sleep" {
    pDnodeActor.PlayerOut += "You are already awake."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //***********
  //* Wake up *
  //***********
  pDnodeActor.pPlayer.Position = "sit"
  PlayerSave(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += "You awake and sit up."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  WakeMsg = pDnodeActor.PlayerName + " awakens and sits up."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, WakeMsg)
}

// Wear command
func DoWear() {
  var ObjectName   string
  var WearFailed   bool
  var WearMsg      string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) == 1 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Wear what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  // Get pointer to object
  TmpStr = StrGetWord(CmdStr, 2)
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  pObject = nil
  IsObjInPlayerInv(TmpStr) // Sets pObject
  if pObject == nil {
    // Player does not have object in inventory
    pDnodeActor.PlayerOut += "You don't have a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += ".\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  pObject.Type = StrMakeLower(pObject.Type)
  if pObject.Type != "armor" {
    // Player can't wear stuff that is NOT armor
    pDnodeActor.PlayerOut += "You can't wear "
    pDnodeActor.PlayerOut += pObject.Desc1
    pDnodeActor.PlayerOut += ".\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    pObject = nil
    return
  }
  // Handle wear positions that require left or right
  if StrIsWord(pObject.WearPosition, "ear wrist finger ankle") {
    // Object must be worn using left and right
    TmpStr = StrGetWord(CmdStr, 3)
    TmpStr = StrMakeLower(TmpStr)
    if StrIsNotWord(TmpStr, "left right") {
      // Player did not specify left or right
      pDnodeActor.PlayerOut += "You must specify left or right"
      pDnodeActor.PlayerOut += ".\r\n"
      CreatePrompt(pDnodeActor.pPlayer)
      pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
      pObject = nil
      return
    }
    pObject.WearPosition += TmpStr
  }
  //***************
  //* Wear object *
  //***************
  // Add object to player's equipment
  WearFailed = AddObjToPlayerEqu(pObject.WearPosition, pObject.ObjectId)
  if WearFailed {
    // Already wearing an object in that wear position
    pDnodeActor.PlayerOut += "You fail to wear "
    pDnodeActor.PlayerOut += pObject.Desc1
    pDnodeActor.PlayerOut += "."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    pObject = nil
    return
  }
  // Increase player's ArmorClass
  pDnodeActor.pPlayer.ArmorClass += pObject.ArmorValue
  // Remove object from player's inventory
  RemoveObjFromPlayerInv(pObject.ObjectId, 1)
  // Send messages
  pDnodeActor.PlayerOut += "You wear "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += "."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  WearMsg = pDnodeActor.PlayerName + " wears " + pObject.Desc1 + "."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, WearMsg)
  pObject = nil
}

// Where command
func DoWhere() {
  var SearchId string

  DEBUGIT(1)
  if StrCountWords(CmdStr) != 2 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Nothing given to search for!"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  SearchId = StrGetWord(CmdStr, 2)
  SearchId = StrMakeLower(SearchId)
  // Find Players
  pDnodeTgt = GetTargetDnode(SearchId)
  if pDnodeTgt != nil {
    // Target is online and in 'playing' state
    pDnodeActor.PlayerOut += "\r\n"
    pDnodeActor.PlayerOut += pDnodeTgt.PlayerName
    pDnodeActor.PlayerOut += " is in "
    pDnodeActor.PlayerOut += pDnodeTgt.pPlayer.RoomId
    pDnodeActor.PlayerOut += "\r\n"
  } else if IsMobValid(SearchId) != nil {
    // Find Mobiles
    WhereMob(SearchId)
  } else {
    pObject = nil
    IsObject(SearchId) // Sets pObject
    if pObject != nil {
      // Find Objects
      WhereObj(SearchId)
    } else {
      // Could not find it
      pDnodeActor.PlayerOut += "\r\n"
      pDnodeActor.PlayerOut += "No "
      pDnodeActor.PlayerOut += SearchId
      pDnodeActor.PlayerOut += " found."
      pDnodeActor.PlayerOut += "\r\n"
    }
  }
  // Prompt
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Who command
func DoWho() {
  var DisplayName  string
  var DisplayLevel string

  DEBUGIT(1)
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "&C"
  pDnodeActor.PlayerOut += "Players online"
  pDnodeActor.PlayerOut += "&N"
  pDnodeActor.PlayerOut += "\r\n"
  pDnodeActor.PlayerOut += "--------------"
  pDnodeActor.PlayerOut += "\r\n"
  // List all players
  SetpDnodeCursorFirst()
  for !EndOfDnodeList() {
    // Loop thru all connections
    pDnodeOthers = GetDnode()
    if pDnodeOthers.PlayerStatePlaying {
      // Who are 'playing'
      if pDnodeOthers.PlayerStateInvisible {
        // Player is invisible, no show
        SetpDnodeCursorNext()
        continue
      }
      DisplayName = fmt.Sprintf("%-15s", pDnodeOthers.PlayerName)
      DisplayLevel = fmt.Sprintf("%3d", pDnodeOthers.pPlayer.Level)
      pDnodeActor.PlayerOut += DisplayName
      pDnodeActor.PlayerOut += " "
      pDnodeActor.PlayerOut += DisplayLevel
      pDnodeActor.PlayerOut += " "
      if pDnodeOthers.PlayerStateAfk {
        // Player is AFK
        pDnodeActor.PlayerOut += "&Y"
        pDnodeActor.PlayerOut += "AFK "
        pDnodeActor.PlayerOut += "&N"
      } else {
        // Player is not AFK
        pDnodeActor.PlayerOut += "    "
      }
      pDnodeActor.PlayerOut += pDnodeOthers.pPlayer.Title
      pDnodeActor.PlayerOut += "&N"
      pDnodeActor.PlayerOut += "\r\n"
    }
    SetpDnodeCursorNext()
  }
  // Re-position pDnodeCursor
  RepositionDnodeCursor()
  // Create player prompt
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Wield command
func DoWield() {
  var ObjectName  string
  var WieldFailed bool
  var WieldMsg    string

  DEBUGIT(1)
  //********************
  //* Validate command *
  //********************
  if IsSleeping() {
    // Player is sleeping, send msg, command is not done
    return
  }
  if StrCountWords(CmdStr) == 1 {
    // Invalid command format
    pDnodeActor.PlayerOut += "Wield what?"
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*************************
  //* Get pointer to object *
  //*************************
  TmpStr = StrGetWord(CmdStr, 2)
  ObjectName = TmpStr
  TmpStr = StrMakeLower(TmpStr)
  pObject = nil
  IsObjInPlayerInv(TmpStr) // Sets pObject
  if pObject == nil {
    // Player does not have object in inventory
    pDnodeActor.PlayerOut += "You don't have a(n) "
    pDnodeActor.PlayerOut += ObjectName
    pDnodeActor.PlayerOut += ".\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //*******************
  //* Is it a weapon? *
  //*******************
  TmpStr = pObject.Type
  TmpStr = StrMakeLower(TmpStr)
  if TmpStr != "weapon" {
    // Player is trying to wield something that is not a weapon
    pDnodeActor.PlayerOut += "Try wielding a weapon."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    return
  }
  //********************
  //* Wield the weapon *
  //********************
  // Add object to player's equipment
  WieldFailed = AddObjToPlayerEqu(pObject.WearPosition, pObject.ObjectId)
  if WieldFailed {
    // Already wielding a weapon
    pDnodeActor.PlayerOut += "You are already wielding a weapon"
    pDnodeActor.PlayerOut += ".\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
    pObject = nil
    return
  }
  // Remove object from player's inventory
  RemoveObjFromPlayerInv(pObject.ObjectId, 1)
  // Send messages
  pDnodeActor.PlayerOut += "You wield "
  pDnodeActor.PlayerOut += pObject.Desc1
  pDnodeActor.PlayerOut += "."
  pDnodeActor.PlayerOut += "\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Send wear message to room
  WieldMsg = pDnodeActor.PlayerName
  WieldMsg += " wields "
  WieldMsg += pObject.Desc1
  WieldMsg += "."
  pDnodeSrc = pDnodeActor
  pDnodeTgt = pDnodeActor
  SendToRoom(pDnodeActor.pPlayer.RoomId, WieldMsg)
  // Set player's weapon info
  pDnodeActor.pPlayer.WeaponDamage = pObject.WeaponDamage
  pDnodeActor.pPlayer.WeaponDesc1 = pObject.Desc1
  pDnodeActor.pPlayer.WeaponType = pObject.WeaponType
  PlayerSave(pDnodeActor.pPlayer)
  pObject = nil
}

// Groups - Calculate group experience, if any
func GrpExperience(MobileExpPoints int, MobileLevel int) {
  var pDnode         *Dnode
  var pDnodeGrpLdr   *Dnode // Group leader
  var pDnodeGrpMem   *Dnode // Group members
  var pPlayer        *Player
  var ExpPoints       int
  var fGrpLimit       float64
  var fGrpMemberCount float64
  var GainLoose       string
  var GrpMemberCount  int
  var i               int
  var LevelTotal      float64
  var PlayerExpPct    int

  // Count group members
  GrpMemberCount = 0
  LevelTotal = 0
  for i = 0; i < GRP_LIMIT; i++ {
    // For each group member
    pPlayer = pDnodeActor.pPlayer.pPlayerGrpMember[0].pPlayerGrpMember[i]
    if pPlayer != nil {
      // Found a member
      GrpMemberCount++
      LevelTotal += float64(pPlayer.Level)
    }
  }
  // Award experience
  fGrpMemberCount = float64(GrpMemberCount)
  fGrpLimit = float64(GRP_LIMIT)
  MobileExpPoints += int(float64(MobileExpPoints) * (fGrpMemberCount / fGrpLimit * MGBP / 100.0))
  pDnodeGrpLdr = GetTargetDnode(pDnodeActor.pPlayer.pPlayerGrpMember[0].Name)
  for i = 0; i < GRP_LIMIT; i++ {
    // For each member of leader's group including the leader
    if pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i] == nil {
      // Done looping through group members
      return
    }
    // Get group member's Dnode
    pDnodeGrpMem = GetTargetDnode(pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i].Name)
    if pDnodeActor.pPlayer.RoomId != pDnodeGrpMem.pPlayer.RoomId {
      // Group members who are not in the same room do not get exp
      continue
    }
    // Calculate experience
    if MobileExpPoints >= 0 {
      // Player gains xp
      PlayerExpPct = int(float64(pDnodeGrpMem.pPlayer.Level) / LevelTotal * 100.0)
      ExpPoints = int(float64(MobileExpPoints) * (float64(PlayerExpPct) / 100.0))
      ExpPoints = CalcAdjustedExpPoints(pDnodeGrpMem.pPlayer.Level, MobileLevel, ExpPoints)
      GainLoose = "gain"
    } else {
      // Player loses xp
      ExpPoints = MobileExpPoints * pDnodeGrpMem.pPlayer.Level
      GainLoose = "loose"
    }
    // Send experience message
    if ExpPoints >= 0 {
      // Player gains xp
      Buf = fmt.Sprintf("%d", ExpPoints)
      TmpStr = Buf
    } else {
      // Player looses xp
      Buf = fmt.Sprintf("%d", ExpPoints*-1)
      TmpStr = Buf
    }
    pDnodeGrpMem.PlayerOut += "\r\n"
    pDnodeGrpMem.PlayerOut += "&Y"
    pDnodeGrpMem.PlayerOut += "You "
    pDnodeGrpMem.PlayerOut += GainLoose
    pDnodeGrpMem.PlayerOut += " "
    pDnodeGrpMem.PlayerOut += TmpStr
    pDnodeGrpMem.PlayerOut += " points of Group Experience!"
    pDnodeGrpMem.PlayerOut += "&N"
    pDnodeGrpMem.PlayerOut += "\r\n"
    // Gain some experience
    pDnode = pDnodeGrpMem
    GainExperience(pDnode, ExpPoints)
    // Save player
    PlayerSave(pDnodeGrpMem.pPlayer)
    // Player prompt
    CreatePrompt(pDnodeGrpMem.pPlayer)
    pDnodeGrpMem.PlayerOut += GetOutput(pDnodeGrpMem.pPlayer)
  }
}

// Groups - Player is leaving the group
func GrpLeave() {
  if pDnodeActor.pPlayer.pPlayerGrpMember[0] == nil {
    // Player is not in a group
    return
  }
  if pDnodeActor.pPlayer == pDnodeActor.pPlayer.pPlayerGrpMember[0] {
    // Player is group leader, disband the whole group
    GrpLeaveLeader()
  } else {
    // Player is a group member
    GrpLeaveMember()
  }
}

// Groups - Leader is leaving - Disband the whole group
func GrpLeaveLeader() {
  var pDnodeGrpMem *Dnode // Other group members
  var i             int
  var j             int

  // Player is group leader, disband the whole group
  for i = 1; i < GRP_LIMIT; i++ {
    // For each group member
    if pDnodeActor.pPlayer.pPlayerGrpMember[i] != nil {
      // Get member's dnode before member's pointer is nulled
      pDnodeGrpMem = GetTargetDnode(pDnodeActor.pPlayer.pPlayerGrpMember[i].Name)
      // Null member's leader pointer
      pDnodeActor.pPlayer.pPlayerGrpMember[i].pPlayerGrpMember[0] = nil
      // Null member's pointer
      pDnodeActor.pPlayer.pPlayerGrpMember[i] = nil
      // Let the group members know that group is disbanded
      pDnodeGrpMem.PlayerOut += "\r\n"
      pDnodeGrpMem.PlayerOut += "The group has been disbanded.\r\n"
      CreatePrompt(pDnodeGrpMem.pPlayer)
      pDnodeGrpMem.PlayerOut += GetOutput(pDnodeGrpMem.pPlayer)
      // Member now has no group, remove any followers of this member
      for j = 0; j < GRP_LIMIT; j++ {
        pDnodeGrpMem.pPlayer.pPlayerFollowers[j] = nil
        pDnodeActor.pPlayer.pPlayerFollowers[j]  = nil
      }
    }
  }
  // Leader now has no group, remove any followers
  for i = 0; i < GRP_LIMIT; i++ {
    pDnodeActor.pPlayer.pPlayerFollowers[i]  = nil
  }
  // Complete the disbanding of the whole group
  pDnodeActor.pPlayer.pPlayerGrpMember[0] = nil
  pDnodeActor.PlayerOut += "Your group has been disbanded.\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
}

// Groups - Member is leaving - Remove them from the group
func GrpLeaveMember() {
  var pDnodeGrpLdr *Dnode // Group leader
  var pDnodeGrpMem *Dnode // Other group members
  var i             int
  var j             int
  var GrpEmpty      bool

  if pDnodeActor.pPlayer.pPlayerFollowers[0] != nil {
    // Player is following someone
    DoFollow(pDnodeActor, "follow none")
  }
  for i = 1; i < GRP_LIMIT; i++ {
    // Loop thru leader's member list to find player
    if pDnodeActor.pPlayer.pPlayerGrpMember[0].pPlayerGrpMember[i] == pDnodeActor.pPlayer {
      // Remove player from leader's group
      pDnodeActor.pPlayer.pPlayerGrpMember[0].pPlayerGrpMember[i] = nil
      j = i // Save player's subscript
      break
    }
  }
  // Complete the disbanding and let the player know
  pDnodeActor.PlayerOut += "You have left the group.\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetOutput(pDnodeActor.pPlayer)
  // Let group leader know when a member leaves the group
  pDnodeGrpLdr = GetTargetDnode(pDnodeActor.pPlayer.pPlayerGrpMember[0].Name)
  pDnodeGrpLdr.PlayerOut += "\r\n"
  pDnodeGrpLdr.PlayerOut += pDnodeActor.PlayerName
  pDnodeGrpLdr.PlayerOut += " has left your group.\r\n"
  GrpEmpty = true
  for i = 1; i < GRP_LIMIT; i++ {
    // For each member of leader's group
    if pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i] != nil {
      // Let other group members know that player has left
      GrpEmpty = false
      pDnodeGrpMem = GetTargetDnode(pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i].Name)
      pDnodeGrpMem.PlayerOut += "\r\n"
      pDnodeGrpMem.PlayerOut += pDnodeActor.PlayerName
      pDnodeGrpMem.PlayerOut += " has left the group.\r\n"
      CreatePrompt(pDnodeGrpMem.pPlayer)
      pDnodeGrpMem.PlayerOut += GetOutput(pDnodeGrpMem.pPlayer)
      if pDnodeActor.pPlayer == pDnodeGrpMem.pPlayer.pPlayerFollowers[0] {
        // Another group member was following player
        DoFollow(pDnodeGrpMem, "follow none")
      }
    }
  }
  pDnodeActor.pPlayer.pPlayerGrpMember[0] = nil
  if GrpEmpty {
    // Player was the last in the group, let the leader know
    pDnodeGrpLdr.pPlayer.pPlayerGrpMember[0] = nil
    pDnodeGrpLdr.PlayerOut += "Your group has disbanded.\r\n"
    // Leader has no group, remove any followers
    for i = 0; i < GRP_LIMIT; i++ {
      pDnodeGrpLdr.pPlayer.pPlayerFollowers[i] = nil
    }
  }
  CreatePrompt(pDnodeGrpLdr.pPlayer)
  pDnodeGrpLdr.PlayerOut += GetOutput(pDnodeGrpLdr.pPlayer)
  // Compact the list of members, so new members are at the end
  for i = j; i < GRP_LIMIT-1; i++ { // j is subscript of member who is leaving
    pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i] = pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i+1]
    pDnodeGrpLdr.pPlayer.pPlayerGrpMember[i+1] = nil
  }
  // When a member leaves a group, remove any followers
  for i = 0; i < GRP_LIMIT; i++ {
    pDnodeActor.pPlayer.pPlayerFollowers[i] = nil
  }
}

// Logon greeting
func LogonGreeting() {
}

// Logon wait male female
func LogonWaitMaleFemale() {
	return
}

// Logon wait name
func LogonWaitName() {
	return
}

// Logon wait name confirmation
func LogonWaitNameConfirmation() {
	return
}

// Logon wait new character
func LogonWaitNewCharacter() {
	return
}

// Logon wait password
func LogonWaitPassword() {
	return
}

// Reposition Dnode cursor
func RepositionDnodeCursor() {
  SetpDnodeCursorFirst()
  for !EndOfDnodeList() {
    pDnodeOthers = GetDnode()
    if pDnodeActor == pDnodeOthers {
      break
    }
    SetpDnodeCursorNext()
  }
}

// New connection
func SockNewConnection() {
  DEBUGIT(1)
  nfd, sa, err := syscall.Accept(ListenSocket)
  if err != nil {
    if err == WSAEWOULDBLOCK {
      return
    }
    Buf = err.Error()
    LogBuf = "Communication::SockNewConnection - Error: accept: " + Buf
    LogIt(LogBuf)
    return
  }
  SockAddr = sa
  IpAddress := ""
  if sa4, ok := SockAddr.(*syscall.SockaddrInet4); ok {
    IpAddress = fmt.Sprintf("%d.%d.%d.%d", sa4.Addr[0], sa4.Addr[1], sa4.Addr[2], sa4.Addr[3])
  }
  err = syscall.SetNonblock(nfd, true)
  if err != nil {
    Buf = err.Error()
    LogBuf = "Communication::SockNewConnection - Error: SetNonblock: " + Buf
    LogIt(LogBuf)
    return
  }
  Buf = fmt.Sprintf("%d", nfd)
  TmpStr = Buf
  LogBuf = "New connection with socket handle "
  LogBuf += TmpStr
  LogBuf += " and address "
  LogBuf += IpAddress
  LogIt(LogBuf)
  pDnodeActor = DnodeConstructor(nfd, IpAddress)
  AppendIt()
  StateConnections = true
}

// Send message
func SockSend(arg string) {
  DEBUGIT(5)
  if arg == "" {
    return
  }
  buf := []byte(arg)
  err := syscall.Sendto(pDnodeActor.DnodeFd, buf, 0, nil)
  if err != nil {
    LogBuf = "SockSend - Error: " + err.Error()
    LogIt(LogBuf)
    return
  }
  pDnodeActor.PlayerOut = ""
}

// Update player statistics that are 'tick' dependant
func UpdatePlayerStats() {
}

// Violence, as in ... WHACK 'em!
func Violence() {
}

// Mobile's turn to do some damage
func ViolenceMobile() {
	return
}

// Mobile has died
func ViolenceMobileDied(MobileBeenWhacked string, MobileDesc1 string, MobileId string) {
	return
}

// Hand out the loot
func ViolenceMobileLoot(Loot string) {
	return
}

// Hand out the loot - for real this time
func ViolenceMobileLootHandOut(Loot string) bool {
	return false
}

// More mobiles to fight?
func ViolenceMobileMore() {
	return
}

// Player's turn to do some damage
func ViolencePlayer() {
	return
}

// Player has died, sad but true
func ViolencePlayerDied(MobileDesc1 string) {
	return
}

// Return current time in milliseconds
func clock() int64 {
  return time.Now().UnixMilli()
}

// Close a socket handle (C++ closesocket equivalent)
func CloseSocket(fd syscall.Handle) int {
  err := syscall.Closesocket(fd)
  if err != nil {
    return 1
  }
  return 0
}