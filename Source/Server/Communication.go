package server

import (
	"fmt"
	"os"
	"syscall"
)

// Globals
//lint:ignore U1000 This field is reserved for future use
var pDnodeOthers *Dnode
var ListenSocket  syscall.Handle
var ValidCmds     []string
var SockAddr syscall.Sockaddr

const WSAEWOULDBLOCK syscall.Errno = 10035
const WSAEINTR       syscall.Errno = 10004

func LogonGreeting() {
}

func UpdatePlayerStats() {
}

func SendToRoom(RoomId string, Message string) {
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

// Color processing placeholder
func Color() {
}

 // Load command array
func CommandArrayLoad() {
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