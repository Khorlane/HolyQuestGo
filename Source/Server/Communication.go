package server

import (
	"bufio"
	"fmt"
	"os"
	"syscall"
)

// Globals
//lint:ignore U1000 This field is reserved for future use
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
            pDnodeActor_pPlayer_Save()
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
            pDnodeActor_pPlayer_Save()
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
          pDnodeActor_pPlayer_CreatePrompt()
          pDnodeActor.PlayerOut += pDnodeActor_pPlayer_GetOutput()
        }
        if pDnodeActor.pPlayer.Thirst > 99 {
          pDnodeActor.pPlayer.Thirst = 100
          pDnodeActor.PlayerOut += "\r\n"
          pDnodeActor.PlayerOut += "You are extremely thirsty!!!"
          pDnodeActor.PlayerOut += "\r\n"
          pDnodeActor_pPlayer_CreatePrompt()
          pDnodeActor.PlayerOut += pDnodeActor_pPlayer_GetOutput()
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
            pDnodeOthers_pPlayer_Save()
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