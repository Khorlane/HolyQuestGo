package server

import (
  "syscall"
)

var DnodeCount   int = 0

// Descriptor Node structure
type Dnode struct {
//lint:ignore U1000 This field is reserved for future use
  pDnodeNext                     *Dnode
//lint:ignore U1000 This field is reserved for future use
  pDnodePrev                     *Dnode
  pPlayer                        *Player
  CmdName1                        string
  CmdName2                        string
  CmdName3                        string
  CmdTime1                        int64
  CmdTime2                        int64
  CmdTime3                        int64
  DnodeFd                         syscall.Handle
  FightTick                       int
  HungerThirstTick                int
  InputTick                       int
  PlayerInp                       string
  PlayerIpAddress                 string
  PlayerName                      string
  PlayerNewCharacter              string
  PlayerOut                       string
  PlayerPassword                  string
  PlayerStateAfk                  bool
  PlayerStateBye                  bool
  PlayerStateFighting             bool
  PlayerStateInvisible            bool
  PlayerStateLoggingOn            bool
  PlayerStatePlaying              bool
  PlayerStateReconnecting         bool
  PlayerStateSendBanner           bool
  PlayerStateWaitMaleFemale       bool
  PlayerStateWaitName             bool
  PlayerStateWaitNameConfirmation bool
  PlayerStateWaitNewCharacter     bool
  PlayerStateWaitPassword         bool
  PlayerWrongPasswordCount        int
  StatsTick                       int
}

var pDnodeActor  *Dnode
//lint:ignore U1000 This field is reserved for future use
var pDnodeSrc    *Dnode
//lint:ignore U1000 This field is reserved for future use
var pDnodeTgt    *Dnode
var pDnodeCursor *Dnode
var pDnodeHead   *Dnode

// Create and initialize a new Dnode instance
func DnodeConstructor(SocketHandle syscall.Handle, IpAddress string) *Dnode {
  DnodeCount = DnodeCount + 1
  return &Dnode{
    pDnodeNext:                      nil,
    pDnodePrev:                      nil,
    pPlayer:                         nil,
    CmdName1:                        "",
    CmdName2:                        "",
    CmdName3:                        "",
    CmdTime1:                        0,
    CmdTime2:                        0,
    CmdTime3:                        0,
    DnodeFd:                         SocketHandle,
    FightTick:                       0,
    HungerThirstTick:                0,
    InputTick:                       0,
    PlayerInp:                       "",
    PlayerIpAddress:                 IpAddress,
    PlayerName:                      "Player name unknown",
    PlayerNewCharacter:              "",
    PlayerOut:                       "",
    PlayerPassword:                  "",
    PlayerStateAfk:                  false,
    PlayerStateBye:                  false,
    PlayerStateFighting:             false,
    PlayerStateInvisible:            false,
    PlayerStateLoggingOn:            false,
    PlayerStatePlaying:              false,
    PlayerStateReconnecting:         false,
    PlayerStateSendBanner:           true,
    PlayerStateWaitMaleFemale:       false,
    PlayerStateWaitName:             false,
    PlayerStateWaitNameConfirmation: false,
    PlayerStateWaitNewCharacter:     false,
    PlayerStateWaitPassword:         false,
    PlayerWrongPasswordCount:        0,
    StatsTick:                       0,
  }
}

// Dnode destructor
func DnodeDestructor(pDnode *Dnode) {
  DnodeCount--
  pDnode.pDnodePrev.pDnodeNext = pDnode.pDnodeNext
  pDnode.pDnodeNext.pDnodePrev = pDnode.pDnodePrev
}


// Return Dnode count
func GetDnodeCount() int {
  return DnodeCount
}