//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Player.go                                             *
// Usage:     Manages player entities and their interactions        *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var PlayerFile   *os.File
var PlayerReader *bufio.Reader
var PlayerCount   int = 0

// Player represents a player in the game.
type Player struct {
  // Public variables
  pPlayerGrpMember [GRP_LIMIT]*Player
  pPlayerFollowers [GRP_LIMIT]*Player
  SessionTime      int64
  RoomIdBeforeMove string

  // Private variables
  pDnode           *Dnode
  Output            string
  PlayerRoomBitPos  int
  PlayerRoomBits    [8]bool
  PlayerRoomChar    byte
  PlayerRoomCharPos int
  PlayerRoomVector  []byte

  // Player file variables
  Name              string
  Password          string
  Admin             bool
  Afk               string
  AllowAssist       bool
  AllowGroup        bool
  ArmorClass        int
  Born              int64
  Color             bool
  Experience        float64
  GoToArrive        string
  GoToDepart        string
  HitPoints         int
  Hunger            int
  Invisible         bool
  Level             int
  MovePoints        int
  OneWhack          bool
  Online            string
  Position          string
  RoomId            string
  RoomInfo          bool
  Sex               string
  Silver            int
  SkillAxe          int
  SkillClub         int
  SkillDagger       int
  SkillHammer       int
  SkillSpear        int
  SkillStaff        int
  SkillSword        int
  Thirst            int
  TimePlayed        int64
  Title             string
  WeaponDamage      int
  WeaponDesc1       string
  WeaponType        string
}

// NewPlayer creates and initializes a new Player instance.
func PlayerConstructor() *Player {
  PlayerCount++
  return &Player{
    pPlayerGrpMember:  [GRP_LIMIT]*Player{},
    pPlayerFollowers:  [GRP_LIMIT]*Player{},
    SessionTime:       0,
    RoomIdBeforeMove:  "",
    pDnode:            nil,
    Output:            "",
    PlayerRoomBitPos:  0,
    PlayerRoomBits:    [8]bool{},
    PlayerRoomChar:    0,
    PlayerRoomCharPos: 0,
    PlayerRoomVector:  []byte{},
    Name:              "",
    Password:          "",
    Admin:             false,
    Afk:               "",
    AllowAssist:       false,
    AllowGroup:        false,
    ArmorClass:        0,
    Born:              0,
    Color:             false,
    Experience:        0.0,
    GoToArrive:        "",
    GoToDepart:        "",
    HitPoints:         0,
    Hunger:            0,
    Invisible:         false,
    Level:             0,
    MovePoints:        0,
    OneWhack:          false,
    Online:            "",
    Position:          "",
    RoomId:            "",
    RoomInfo:          false,
    Sex:               "",
    Silver:            0,
    SkillAxe:          0,
    SkillClub:         0,
    SkillDagger:       0,
    SkillHammer:       0,
    SkillSpear:        0,
    SkillStaff:        0,
    SkillSword:        0,
    Thirst:            0,
    TimePlayed:        0,
    Title:             "",
    WeaponDamage:      0,
    WeaponDesc1:       "",
    WeaponType:        "",
  }
}

// player destructor
func PlayerDestructor(pPlayer *Player) {
  PlayerCount--
}

// Calculate adjusted experience points
func CalcAdjustedExpPoints(PlayerLevel int, MobileLevel int, ExpPoints int) int {
	// Implement the logic for calculating adjusted experience points here.
	return 0
}

// Calculate experience needed to obtain the next level
func CalcLevelExperience(Level int) float64 {
  var Experience float64
  return Experience
}

// Calculate additional experience points based on level and base experience.
func CalcLevelExperienceAdd(Level int, BaseExp float64) float64 {
	// Implement the logic for calculating additional experience points here.
	return 0.0
}

// Calculate the base experience points for a given level.
func CalcLevelExperienceBase(Level int) float64 {
	return 0.0
}

//Is this a valid Player? 
func IsPlayer(PlayerName string) bool {
	return false
}

// Return player count
func GetCount() int {
	return PlayerCount
}

// Validate player name
func IsNameValid(Name string) bool {
	return false
}

// Create player prompt
func CreatePrompt(pPlayer *Player) {
  pPlayer.Output = "\r\n"
  hpStr := fmt.Sprintf("%d", pPlayer.HitPoints)
  pPlayer.Output += hpStr + "H "
  mpStr := fmt.Sprintf("%d", pPlayer.MovePoints)
  pPlayer.Output += mpStr + "M "
  pPlayer.Output += "> "
}

// Player drink
func Drink(pPlayer *Player, Percent int) {
	return
}

// Player eat
func Eat(pPlayer *Player, Percent int) {
	return
}

// Player gains some experience
func GainExperience(pDnode *Dnode, ExperienceToBeGained int) {
	return
}

// Return the current output string for the player
func GetOutput(pPlayer *Player) string {
  return pPlayer.Output
}

// Get skill for current weapon
func GetWeaponSkill() int {
	return 0
}

// Parse player stuff
func ParsePlayerStuff() {
	return
}

// Save player stuff
func PlayerSave(pPlayer *Player) {
  if !OpenPlayerFile(pPlayer.Name, "Write") {
    LogBuf = "Player::Save - Error opening player file for write, Players directory may not exist"
    LogIt(LogBuf)
    return
  }
  // Name
  Stuff = "Name:" + pPlayer.Name
  PlayerWriteLine(Stuff)
  // Password
  Stuff = "Password:" + pPlayer.Password
  PlayerWriteLine(Stuff)
  // Admin
  if pPlayer.Admin {
    Stuff = "Admin:Yes"
  } else {
    Stuff = "Admin:No"
  }
  PlayerWriteLine(Stuff)
  // AFK
  if pPlayer.pDnode.PlayerStateAfk {
    Stuff = "AFK:Yes"
  } else {
    Stuff = "AFK:No"
  }
  PlayerWriteLine(Stuff)
  // AllowAssist
  if pPlayer.AllowAssist {
    Stuff = "AllowAssist:Yes"
  } else {
    Stuff = "AllowAssist:No"
  }
  PlayerWriteLine(Stuff)
  // AllowGroup
  if pPlayer.AllowGroup {
    Stuff = "AllowGroup:Yes"
  } else {
    Stuff = "AllowGroup:No"
  }
  PlayerWriteLine(Stuff)
  // ArmorClass - save only, ParsePlayerStuff calls CalcPlayerArmorClass
  TmpStr = fmt.Sprintf("%d", pPlayer.ArmorClass)
  Stuff = "ArmorClass:" + TmpStr
  PlayerWriteLine(Stuff)
  // Born
  TmpStr = fmt.Sprintf("%d", pPlayer.Born)
  Stuff = "Born:" + TmpStr
  PlayerWriteLine(Stuff)
  // Color
  if pPlayer.Color {
    Stuff = "Color:Yes"
  } else {
    Stuff = "Color:No"
  }
  PlayerWriteLine(Stuff)
  // Experience
  TmpStr = fmt.Sprintf("%15.0f", pPlayer.Experience)
  Stuff = "Experience:" + TmpStr
  PlayerWriteLine(Stuff)
  // GoToArrive
  Stuff = "GoToArrive:" + pPlayer.GoToArrive
  PlayerWriteLine(Stuff)
  // GoToDepart
  Stuff = "GoToDepart:" + pPlayer.GoToDepart
  PlayerWriteLine(Stuff)
  // HitPoints
  TmpStr = fmt.Sprintf("%d", pPlayer.HitPoints)
  Stuff = "HitPoints:" + TmpStr
  PlayerWriteLine(Stuff)
  // Hunger
  TmpStr = fmt.Sprintf("%d", pPlayer.Hunger)
  Stuff = "Hunger:" + TmpStr
  PlayerWriteLine(Stuff)
  // Invisible
  if pPlayer.Invisible {
    Stuff = "Invisible:Yes"
  } else {
    Stuff = "Invisible:No"
  }
  PlayerWriteLine(Stuff)
  // Level
  TmpStr = fmt.Sprintf("%d", pPlayer.Level)
  Stuff = "Level:" + TmpStr
  PlayerWriteLine(Stuff)
  // MovePoints
  TmpStr = fmt.Sprintf("%d", pPlayer.MovePoints)
  Stuff = "MovePoints:" + TmpStr
  PlayerWriteLine(Stuff)
  // OneWhack
  if pPlayer.OneWhack {
    Stuff = "OneWhack:Yes"
  } else {
    Stuff = "OneWhack:No"
  }
  PlayerWriteLine(Stuff)
  // Online
  if pPlayer.pDnode.PlayerStatePlaying {
    Stuff = "Online:Yes"
  } else {
    Stuff = "Online:No"
  }
  PlayerWriteLine(Stuff)
  // Position
  Stuff = "Position:" + pPlayer.Position
  PlayerWriteLine(Stuff)
  // RoomId
  Stuff = "RoomId:" + pPlayer.RoomId
  PlayerWriteLine(Stuff)
  // RoomInfo
  if pPlayer.RoomInfo {
    Stuff = "RoomInfo:Yes"
  } else {
    Stuff = "RoomInfo:No"
  }
  PlayerWriteLine(Stuff)
  // Sex
  Stuff = "Sex:" + pPlayer.Sex
  PlayerWriteLine(Stuff)
  // Silver
  TmpStr = fmt.Sprintf("%d", pPlayer.Silver)
  Stuff = "Silver:" + TmpStr
  PlayerWriteLine(Stuff)
  // SkillAxe
  TmpStr = fmt.Sprintf("%d", pPlayer.SkillAxe)
  Stuff = "SkillAxe:" + TmpStr
  PlayerWriteLine(Stuff)
  // SkillClub
  TmpStr = fmt.Sprintf("%d", pPlayer.SkillClub)
  Stuff = "SkillClub:" + TmpStr
  PlayerWriteLine(Stuff)
  // SkillDagger
  TmpStr = fmt.Sprintf("%d", pPlayer.SkillDagger)
  Stuff = "SkillDagger:" + TmpStr
  PlayerWriteLine(Stuff)
  // SkillHammer
  TmpStr = fmt.Sprintf("%d", pPlayer.SkillHammer)
  Stuff = "SkillHammer:" + TmpStr
  PlayerWriteLine(Stuff)
  // SkillSpear
  TmpStr = fmt.Sprintf("%d", pPlayer.SkillSpear)
  Stuff = "SkillSpear:" + TmpStr
  PlayerWriteLine(Stuff)
  // SkillStaff
  TmpStr = fmt.Sprintf("%d", pPlayer.SkillStaff)
  Stuff = "SkillStaff:" + TmpStr
  PlayerWriteLine(Stuff)
  // SkillSword
  TmpStr = fmt.Sprintf("%d", pPlayer.SkillSword)
  Stuff = "SkillSword:" + TmpStr
  PlayerWriteLine(Stuff)
  // Thirst
  TmpStr = fmt.Sprintf("%d", pPlayer.Thirst)
  Stuff = "Thirst:" + TmpStr
  PlayerWriteLine(Stuff)
  // TimePlayed
  if pPlayer.pDnode.PlayerStatePlaying {
    // Don't update TimePlayed if player is not 'playing'
    pPlayer.TimePlayed += GetTimeSeconds() - pPlayer.SessionTime
    pPlayer.SessionTime = GetTimeSeconds()
  }
  TmpStr = fmt.Sprintf("%d", pPlayer.TimePlayed)
  Stuff = "TimePlayed:" + TmpStr
  PlayerWriteLine(Stuff)
  // Title
  Stuff = "Title:" + pPlayer.Title
  PlayerWriteLine(Stuff)
  // WeaponDamage
  TmpStr = fmt.Sprintf("%d", pPlayer.WeaponDamage)
  Stuff = "WeaponDamage:" + TmpStr
  PlayerWriteLine(Stuff)
  // WeaponDesc1
  Stuff = "WeaponDesc1:" + pPlayer.WeaponDesc1
  PlayerWriteLine(Stuff)
  // WeaponType
  Stuff = "WeaponType:" + pPlayer.WeaponType
  PlayerWriteLine(Stuff)
  // Done
  PlayerCloseFile()
}

// Manipulate player money
func SetMoney(PlusMinus byte, Amount int, Metal string) {
	return
}

// Show player money
func ShowMoney(pPlayer *Player) {
	return
}

// Show player status
func ShowStatus(pPlayer *Player) {
	return
}

// Close player file
func PlayerCloseFile() {
  PlayerFile.Close()
}

// Open player file
func OpenPlayerFile(Name string, Mode string) bool {
  PlayerFileName := PLAYER_DIR + Name + ".txt"
  if Mode == "Read" {
    f, err := os.Open(PlayerFileName)
    if err != nil {
      return false
    }
    PlayerFile = f
    PlayerReader = bufio.NewReader(PlayerFile)
    return true
  } else if Mode == "Write" {
    f, err := os.Create(PlayerFileName)
    if err != nil {
      return false
    }
    PlayerFile = f
    PlayerReader = bufio.NewReader(PlayerFile)
    return true
  } else {
    LogBuf = "Player::OpenFile - Mode is not 'Read' or 'Write'"
    LogIt(LogBuf)
    os.Exit(1)
    return false
  }
}

// Read a line from player file
func PlayerReadLine() {
  line, err := PlayerReader.ReadString('\n')
  if err != nil {
    Stuff = ""
    return
  }
  Stuff = strings.TrimRight(line, "\r\n")
}

// Write a line to player file
func PlayerWriteLine(Stuff string) {
  Stuff = Stuff + "\n"
  if PlayerFile != nil {
    _, _ = PlayerFile.WriteString(Stuff)
    PlayerFile.Sync()
  }
}
// Check whether or not player has been in the current room
func PlayerRoomHasNotBeenHere() bool {
	return false
}

// Convert from PlayerRoom char to PlayerRoom bits
func PlayerRoomCharToBitsConvert() {
	return
}

// Convert from PlayerRoom bits to PlayerRoom char
func PlayerRoomBitsToCharConvert() {
	return
}

// Read PlayerRoomVector from disk
func PlayerRoomStringRead() {
	return
}

// Write PlayerRoomVector to disk
func PlayerRoomStringWrite() {
	return
}