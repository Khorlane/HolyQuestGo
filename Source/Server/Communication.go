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

func LogonGreeting() {
}

func UpdatePlayerStats() {
}

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
  pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
  pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
        pDnodeOthers.PlayerOut += GetPlayerOutput(pDnodeOthers.pPlayer)
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
              pDnodeOthers.PlayerOut += GetPlayerOutput(pDnodeOthers.pPlayer)
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

// Close a socket handle (C++ closesocket equivalent)
func CloseSocket(fd syscall.Handle) int {
  err := syscall.Closesocket(fd)
  if err != nil {
    return 1
  }
  return 0
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
          pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
        }
        if pDnodeActor.pPlayer.Thirst > 99 {
          pDnodeActor.pPlayer.Thirst = 100
          pDnodeActor.PlayerOut += "\r\n"
          pDnodeActor.PlayerOut += "You are extremely thirsty!!!"
          pDnodeActor.PlayerOut += "\r\n"

          CreatePrompt(pDnodeActor.pPlayer)
          pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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

// Return current time in milliseconds
func clock() int64 {
  return time.Now().UnixMilli()
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
      pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
          pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
  pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
  var Level int
  var LevelString string
  var PlayerName string
  var PlayerNameSave string
  var TargetName string
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
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
    return
  }
  // Get target Dnode pointer
  pDnodeTgt = GetTargetDnode(TargetName)
  if pDnodeTgt == nil {
    // Target player not found
    pDnodeActor.PlayerOut += TargetNameSave
    pDnodeActor.PlayerOut += " is not online.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
  pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
  pDnodeTgt.PlayerOut += GetPlayerOutput(pDnodeTgt.pPlayer)
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
  pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
}

// Assist command
func DoAssist() {
  var AssistMsg string
  var MobileId string
  var PlayerNameCheck string
  var TargetNameCheck string
  var TargetNameSave string
  var TargetNotHere bool

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
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
    return
  }
  //*********************
  //* Turning assist on *
  //*********************
  if TargetNameCheck == "on" {
    pDnodeActor.pPlayer.AllowAssist = true
    pDnodeActor.PlayerOut += "You are now accepting assists.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
    return
  }
  //**********************
  //* Turning assist off *
  //**********************
  if TargetNameCheck == "off" {
    pDnodeActor.pPlayer.AllowAssist = false
    pDnodeActor.PlayerOut += "You are now rejecting assists.\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
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
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
    return
  }
  if !pDnodeTgt.PlayerStateFighting {
    // Tell player that target is not fighting
    pDnodeActor.PlayerOut += pDnodeTgt.pPlayer.Name
    pDnodeActor.PlayerOut += " is not fighting."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
    return
  }
  if !pDnodeTgt.pPlayer.AllowAssist {
    // Tell player that target is not accepting assistance
    pDnodeActor.PlayerOut += pDnodeTgt.pPlayer.Name
    pDnodeActor.PlayerOut += " is not accepting assistance."
    pDnodeActor.PlayerOut += "\r\n"
    CreatePrompt(pDnodeActor.pPlayer)
    pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
    return
  }
  //******************
  //* Send to player *
  //******************
  pDnodeActor.PlayerOut += "You begin assisting "
  pDnodeActor.PlayerOut += pDnodeTgt.pPlayer.Name
  pDnodeActor.PlayerOut += ".\r\n"
  CreatePrompt(pDnodeActor.pPlayer)
  pDnodeActor.PlayerOut += GetPlayerOutput(pDnodeActor.pPlayer)
  //******************
  //* Send to target *
  //******************
  pDnodeTgt.PlayerOut += "\r\n"
  pDnodeTgt.PlayerOut += pDnodeActor.pPlayer.Name
  pDnodeTgt.PlayerOut += " begins assisting you.\r\n"
  CreatePrompt(pDnodeTgt.pPlayer)
  pDnodeTgt.PlayerOut += GetPlayerOutput(pDnodeTgt.pPlayer)
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

func DoBuy() {
  // TODO: implement DoBuy
}

func DoChat() {
  // TODO: implement DoChat
}

func DoColor() {
  // TODO: implement DoColor
}

func DoConsider() {
  // TODO: implement DoConsider
}

func DoDelete() {
  // TODO: implement DoDelete
}

func DoDestroy() {
  // TODO: implement DoDestroy
}

func DoDrink() {
  // TODO: implement DoDrink
}

func DoDrop() {
  // TODO: implement DoDrop
}

func DoEat() {
  // TODO: implement DoEat
}

func DoEmote() {
  // TODO: implement DoEmote
}

func DoEquipment() {
  // TODO: implement DoEquipment
}

func DoExamine() {
  // TODO: implement DoExamine
}

func DoFlee() {
  // TODO: implement DoFlee
}

func DoFollow(pDnode *Dnode, CmdStr string) {
  // TODO: implement DoFollow
}

func DoGet() {
  // TODO: implement DoGet
}

func DoGive() {
  // TODO: implement DoGive
}

func DoGo() {
  // TODO: implement DoGo
}

func DoGoTo() {
  // TODO: implement DoGoTo
}

func DoGoToArrive() {
  // TODO: implement DoGoToArrive
}

func DoGoToDepart() {
  // TODO: implement DoGoToDepart
}

func DoGroup() {
  // TODO: implement DoGroup
}

func DoGsay() {
  // TODO: implement DoGsay
}

func DoHail() {
  // TODO: implement DoHail
}

func DoHelp() {
  // TODO: implement DoHelp
}

func DoInventory() {
  // TODO: implement DoInventory
}

func DoInvisible() {
  // TODO: implement DoInvisible
}

func DoKill() {
  // TODO: implement DoKill
}

func DoList() {
  // TODO: implement DoList
}

func DoLoad() {
  // TODO: implement DoLoad
}

func DoLogon() {
  // TODO: implement DoLogon
}

func DoLook(CmdStr string) {
  // TODO: implement DoLook
}

func DoMoney() {
  // TODO: implement DoMoney
}

func DoMotd() {
  // TODO: implement DoMotd
}

func DoOneWhack() {
  // TODO: implement DoOneWhack
}

func DoPassword() {
  // TODO: implement DoPassword
}

func DoPlayed() {
  // TODO: implement DoPlayed
}

func DoQuit() {
  // TODO: implement DoQuit
}

func DoRefresh() {
  // TODO: implement DoRefresh
}

func DoRemove() {
  // TODO: implement DoRemove
}

func DoRestore(CmdStr string) {
  // TODO: implement DoRestore
}

func DoRoomInfo() {
  // TODO: implement DoRoomInfo
}

func DoSave() {
  // TODO: implement DoSave
}

func DoSay() {
  // TODO: implement DoSay
}

func DoSell() {
  // TODO: implement DoSell
}

func DoShow() {
  // TODO: implement DoShow
}

func DoSit() {
  // TODO: implement DoSit
}

func DoSleep() {
  // TODO: implement DoSleep
}

func DoStand() {
  // TODO: implement DoStand
}

func DoStatus() {
  // TODO: implement DoStatus
}

func DoStop() {
  // TODO: implement DoStop
}

func DoTell() {
  // TODO: implement DoTell
}

func DoTime() {
  // TODO: implement DoTime
}

func DoTitle() {
  // TODO: implement DoTitle
}

func DoTrain() {
  // TODO: implement DoTrain
}

func DoWake() {
  // TODO: implement DoWake
}

func DoWear() {
  // TODO: implement DoWear
}

func DoWhere() {
  // TODO: implement DoWhere
}

func DoWho() {
  // TODO: implement DoWho
}

func DoWield() {
  // TODO: implement DoWield
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

func Violence() {
  // Implementation goes here
}