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
  "math"  
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
    pDnode:            pDnodeActor,
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
  var PctOfExpPoints    int
  var AdjustedExpPoints int

  PctOfExpPoints = 100 - (PlayerLevel-MobileLevel-2)*20
  if PctOfExpPoints < 1 {
    PctOfExpPoints = 0
  }
  if PctOfExpPoints > 100 {
    PctOfExpPoints = 100
  }
  AdjustedExpPoints = (ExpPoints*PctOfExpPoints + 99) / 100
  return AdjustedExpPoints
}

// Calculate experience needed to obtain the next level
func CalcLevelExperience(Level int) float64 {
  var AddExp   float64
  var BaseExp  float64
  var TotalExp float64

  BaseExp = CalcLevelExperienceBase(Level)
  AddExp  = CalcLevelExperienceAdd(Level, BaseExp)
  TotalExp = BaseExp + AddExp
  return TotalExp
}

// Calculate additional experience points based on level and base experience
func CalcLevelExperienceAdd(Level int, BaseExp float64) float64 {
  var AddExp   float64
  var LogLevel float64

  LogLevel = math.Log10(float64(Level) + 20)
  AddExp   = math.Pow(BaseExp, LogLevel) * (float64(Level) / 10000.0)
  return AddExp
}

// Calculate the base experience points for a given level
func CalcLevelExperienceBase(Level int) float64 {
  // Recursive
  // Assuming PLAYER_EXP_PER_LEVEL = 1000
  // experience needed to get to level 5 is:
  // 14000 = 5 * 1000 + 4 * 1000 + 3 * 1000 + 2 * 1000
  if Level < 2 {
    return 0
  }
  return float64(Level*PLAYER_EXP_PER_LEVEL) + CalcLevelExperienceBase(Level-1)
}

// Is this a valid Player?
func IsPlayer(PlayerName string) bool {
  var PlayerFileName string

  PlayerFileName = PLAYER_DIR + PlayerName + ".txt"
  if FileExist(PlayerFileName) {
    return true
  } else {
    return false
  }
}

// Return player count
func GetCount() int {
  return PlayerCount
}

// Validate player name
func IsNameValid(Name string) bool {
  var NameIn              string
  var ValidNameFile      *os.File
  var ValidNamesFileName  string
  var err                 error

  ValidNamesFileName = VALID_NAMES_DIR + "ValidNames.txt"
  ValidNameFile, err = os.Open(ValidNamesFileName)
  if err != nil {
    LogBuf = "Player::IsNameValid - Error opening valid name file, it may not exist"
    LogIt(LogBuf)
    os.Exit(1)
  }
  defer ValidNameFile.Close()
  Name = StrMakeLower(Name)
  Scanner := bufio.NewScanner(ValidNameFile)
  if Scanner.Scan() {
    NameIn = Scanner.Text()
    NameIn = StrMakeLower(NameIn)
  }
  for NameIn != "" {
    if Name == NameIn {
      return true
    }
    if !Scanner.Scan() {
      break
    }
    NameIn = Scanner.Text()
    NameIn = StrMakeLower(NameIn)
  }
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
  pPlayer.Thirst -= Percent
  if pPlayer.Thirst <= 0 {
    // Not thirsty
    pPlayer.Thirst = 0
    pPlayer.Output  = "You are no longer thirsty, not even a little bit."
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Thirst < 20 {
    // A little bit thirsty
    pPlayer.Output  = "You are a little bit thirsty."
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Thirst < 40 {
    // Lip balm
    pPlayer.Output  = "You need some lip balm."
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Thirst < 60 {
    // Thirsty
    pPlayer.Output  = "You are thirsty."
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Thirst < 80 {
    // Parched
    pPlayer.Output  = "Your throat is parched!"
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Thirst < 100 {
    // Extreme thirst
    pPlayer.Output  = "You are extremely thirsty!!!"
    pPlayer.Output += "\r\n"
    return
  }
}

// Player eat
func Eat(pPlayer *Player, Percent int) {
  pPlayer.Hunger -= Percent
  if pPlayer.Hunger <= 0 {
    // Not hungry
    pPlayer.Hunger = 0
    pPlayer.Output  = "You are no longer hungry, not even a little bit."
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Hunger < 20 {
    // A little bit hungry
    pPlayer.Output  = "You are a little bit hungry."
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Hunger < 40 {
    // Stomach growling
    pPlayer.Output  = "Your stomach is growling."
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Hunger < 60 {
    // Hungry
    pPlayer.Output  = "You are hungry."
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Hunger < 80 {
    // Eat a horse
    pPlayer.Output  = "You could eat a horse!"
    pPlayer.Output += "\r\n"
    return
  } else if pPlayer.Hunger < 100 {
    // Extreme hunger
    pPlayer.Output  = "You are extremely hungry!!!"
    pPlayer.Output += "\r\n"
    return
  }
}

// Player gains some experience
func GainExperience(pDnode *Dnode, ExperienceToBeGained int) {
  var LevelExperience float64

  pDnode.pPlayer.Experience += float64(ExperienceToBeGained)
  LevelExperience = CalcLevelExperience(pDnode.pPlayer.Level + 1)
  if pDnode.pPlayer.Experience >= LevelExperience {
    // Player has enough experience to gain a level
    pDnode.PlayerOut += "&Y"
    pDnode.PlayerOut += "You gain a LEVEL!!!"
    pDnode.PlayerOut += "&N"
    pDnode.PlayerOut += "\r\n"
    pDnode.pPlayer.Level++
    LogBuf  = pDnode.PlayerName
    LogBuf += " has gained level "
    TmpStr = fmt.Sprintf("%d", pDnode.pPlayer.Level)
    LogBuf += TmpStr
    LogBuf += "!"
    LogIt(LogBuf)
  }
}

// Return the current output string for the player
func GetOutput(pPlayer *Player) string {
  return pPlayer.Output
}

// Get skill for current weapon
func GetWeaponSkill(pPlayer *Player) int {
  var WeaponSkill int

  WeaponSkill = 0
  pPlayer.WeaponType = StrMakeLower(pPlayer.WeaponType)
  switch (pPlayer.WeaponType) {
  case "axe":
    WeaponSkill = pPlayer.SkillAxe
  case "club":
    WeaponSkill = pPlayer.SkillClub
  case "dagger":
    WeaponSkill = pPlayer.SkillDagger
  case "hammer":
    WeaponSkill = pPlayer.SkillHammer
  case "spear":
    WeaponSkill = pPlayer.SkillSpear
  case "staff":
    WeaponSkill = pPlayer.SkillStaff
  case "sword":
    WeaponSkill = pPlayer.SkillSword
  }
  return WeaponSkill
}

// Parse player stuff
func ParsePlayerStuff(pPlayer *Player) {
  var Amount int
  var Name   string

  DEBUGIT(1)
  Name = pPlayer.Name
  if !OpenPlayerFile(Name, "Read") {
    LogBuf = "Player::Save - Error opening player file for read, Players directory may not exist"
    LogIt(LogBuf)
    os.Exit(1)
  }
  PlayerReadLine()
  for Stuff != "" {
    if StrLeft(Stuff, 5) == "Name:" {
      TmpStr = ""
    } else if StrLeft(Stuff, 9) == "Password:" {
      pPlayer.Password = StrRight(Stuff, StrGetLength(Stuff)-9)
    } else if StrLeft(Stuff, 6) == "Admin:" {
      TmpStr = StrRight(Stuff, StrGetLength(Stuff)-6)
      TmpStr = StrMakeLower(TmpStr)
      if TmpStr == "yes" {
        pPlayer.Admin = true
      } else {
        pPlayer.Admin = false
      }
    } else if StrLeft(Stuff, 4) == "AFK:" {
      TmpStr = ""
    } else if StrLeft(Stuff, 12) == "AllowAssist:" {
      TmpStr = StrRight(Stuff, StrGetLength(Stuff)-12)
      TmpStr = StrMakeLower(TmpStr)
      if TmpStr == "yes" {
        pPlayer.AllowAssist = true
      } else {
        pPlayer.AllowAssist = false
      }
    } else if StrLeft(Stuff, 11) == "AllowGroup:" {
      TmpStr = StrRight(Stuff, StrGetLength(Stuff)-11)
      TmpStr = StrMakeLower(TmpStr)
      if TmpStr == "yes" {
        pPlayer.AllowGroup = true
      } else {
        pPlayer.AllowGroup = false
      }
    } else if StrLeft(Stuff, 11) == "ArmorClass:" {
      pPlayer.ArmorClass = CalcPlayerArmorClass(pPlayer)
    } else if StrLeft(Stuff, 5) == "Born:" {
      pPlayer.Born = int64(StrToInt(StrRight(Stuff, StrGetLength(Stuff)-5)))
    } else if StrLeft(Stuff, 6) == "Color:" {
      TmpStr = StrRight(Stuff, StrGetLength(Stuff)-6)
      TmpStr = StrMakeLower(TmpStr)
      if TmpStr == "yes" {
        pPlayer.Color = true
      } else {
        pPlayer.Color = false
      }
    } else if StrLeft(Stuff, 11) == "Experience:" {
      pPlayer.Experience = float64(StrToInt(StrRight(Stuff, StrGetLength(Stuff)-11)))
    } else if StrLeft(Stuff, 11) == "GoToArrive:" {
      pPlayer.GoToArrive = StrRight(Stuff, StrGetLength(Stuff)-11)
      if pPlayer.Admin {
        if pPlayer.GoToArrive == "" {
          pPlayer.GoToArrive = "arrives!"
        }
      }
    } else if StrLeft(Stuff, 11) == "GoToDepart:" {
      pPlayer.GoToDepart = StrRight(Stuff, StrGetLength(Stuff)-11)
      if pPlayer.Admin {
        if pPlayer.GoToDepart == "" {
          pPlayer.GoToDepart = "leaves!"
        }
      }
    } else if StrLeft(Stuff, 10) == "HitPoints:" {
      pPlayer.HitPoints = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-10))
    } else if StrLeft(Stuff, 7) == "Hunger:" {
      pPlayer.Hunger = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-7))
    } else if StrLeft(Stuff, 10) == "Invisible:" {
      TmpStr = StrRight(Stuff, StrGetLength(Stuff)-10)
      TmpStr = StrMakeLower(TmpStr)
      if TmpStr == "yes" {
        pPlayer.Invisible = true
        pDnodeActor.PlayerStateInvisible = true
      } else {
        pPlayer.Invisible = false
        pDnodeActor.PlayerStateInvisible = false
      }
    } else if StrLeft(Stuff, 6) == "Level:" {
      pPlayer.Level = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-6))
    } else if StrLeft(Stuff, 11) == "MovePoints:" {
      pPlayer.MovePoints = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-11))
    } else if StrLeft(Stuff, 9) == "OneWhack:" {
      TmpStr = StrRight(Stuff, StrGetLength(Stuff)-9)
      TmpStr = StrMakeLower(TmpStr)
      if TmpStr == "yes" {
        pPlayer.OneWhack = true
      } else {
        pPlayer.OneWhack = false
      }
    } else if StrLeft(Stuff, 7) == "Online:" {
      TmpStr = ""
    } else if StrLeft(Stuff, 9) == "Position:" {
      pPlayer.Position = StrRight(Stuff, StrGetLength(Stuff)-9)
    } else if StrLeft(Stuff, 7) == "RoomId:" {
      pPlayer.RoomId = StrRight(Stuff, StrGetLength(Stuff)-7)
    } else if StrLeft(Stuff, 9) == "RoomInfo:" {
      TmpStr = StrRight(Stuff, StrGetLength(Stuff)-9)
      TmpStr = StrMakeLower(TmpStr)
      if TmpStr == "yes" {
        pPlayer.RoomInfo = true
      } else {
        pPlayer.RoomInfo = false
      }
    } else if StrLeft(Stuff, 4) == "Sex:" {
      pPlayer.Sex = StrRight(Stuff, StrGetLength(Stuff)-4)
      pPlayer.Sex = StrMakeUpper(pPlayer.Sex)
    } else if StrLeft(Stuff, 7) == "Silver:" {
      Amount = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-7))
      SetMoney(pPlayer, '+', Amount, "Silver")
    } else if StrLeft(Stuff, 9) == "SkillAxe:" {
      pPlayer.SkillAxe = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-9))
    } else if StrLeft(Stuff, 10) == "SkillClub:" {
      pPlayer.SkillClub = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-10))
    } else if StrLeft(Stuff, 12) == "SkillDagger:" {
      pPlayer.SkillDagger = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-12))
    } else if StrLeft(Stuff, 12) == "SkillHammer:" {
      pPlayer.SkillHammer = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-12))
    } else if StrLeft(Stuff, 11) == "SkillSpear:" {
      pPlayer.SkillSpear = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-11))
    } else if StrLeft(Stuff, 11) == "SkillStaff:" {
      pPlayer.SkillStaff = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-11))
    } else if StrLeft(Stuff, 11) == "SkillSword:" {
      pPlayer.SkillSword = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-11))
    } else if StrLeft(Stuff, 7) == "Thirst:" {
      pPlayer.Thirst = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-7))
    } else if StrLeft(Stuff, 11) == "TimePlayed:" {
      pPlayer.TimePlayed = int64(StrToInt(StrRight(Stuff, StrGetLength(Stuff)-11)))
    } else if StrLeft(Stuff, 6) == "Title:" {
      pPlayer.Title = StrRight(Stuff, StrGetLength(Stuff)-6)
    } else if StrLeft(Stuff, 13) == "WeaponDamage:" {
      pPlayer.WeaponDamage = StrToInt(StrRight(Stuff, StrGetLength(Stuff)-13))
    } else if StrLeft(Stuff, 12) == "WeaponDesc1:" {
      pPlayer.WeaponDesc1 = StrRight(Stuff, StrGetLength(Stuff)-12)
    } else if StrLeft(Stuff, 11) == "WeaponType:" {
      pPlayer.WeaponType = StrRight(Stuff, StrGetLength(Stuff)-11)
      pPlayer.WeaponType = StrMakeLower(pPlayer.WeaponType)
    } else {
      LogBuf = Name
      LogBuf += " has an unidentified player file field"
      LogIt(LogBuf)
      LogBuf = Stuff
      LogIt(LogBuf)
    }
    if pPlayer.WeaponType == "hand" {
      pPlayer.WeaponDamage = PLAYER_DMG_HAND
    }
    PlayerReadLine()
  }
  PlayerCloseFile()
}

// Save player stuff
func PlayerSave(pPlayer *Player) {
  DEBUGIT(1)
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
func SetMoney(pPlayer *Player, PlusMinus byte, Amount int, Metal string) {
  if PlusMinus == '-' {
    Amount = Amount * -1
  }

  if Metal == "Silver" {
    pPlayer.Silver = pPlayer.Silver + Amount
  }
}

// Show player money
func ShowMoney(pPlayer *Player) {
  TmpStr = fmt.Sprintf("%d", pPlayer.Silver)
  pPlayer.Output = "Silver: " + TmpStr + "\r\n"
}

func ShowStatus(pPlayer *Player) {
  var Exp1 string
  var Exp2 string

  pPlayer.Output = "\r\n"
  // Name
  pPlayer.Output += "Name:         "
  pPlayer.Output += pPlayer.Name
  pPlayer.Output += "\r\n"
  // Sex
  pPlayer.Output += "Sex:          "
  pPlayer.Output += pPlayer.Sex
  pPlayer.Output += "\r\n"
  // Level
  TmpStr = fmt.Sprintf("%d", pPlayer.Level)
  pPlayer.Output += "Level:        "
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
  // Hit Points
  pPlayer.Output += "Hit Points:   "
  TmpStr = fmt.Sprintf("%d", pPlayer.HitPoints)
  pPlayer.Output += TmpStr
  pPlayer.Output += "/"
  TmpStr = fmt.Sprintf("%d", pPlayer.Level*PLAYER_HPT_PER_LEVEL)
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
  // Current Experience and Experience needed for next level
  TmpStr = fmt.Sprintf("%15.0f", pPlayer.Experience)
  TmpStr = StrLeft(TmpStr, StrFindFirstChar(TmpStr, '.'))
  Exp1 = FormatCommas(TmpStr)

  TmpStr = fmt.Sprintf("%15.0f", CalcLevelExperience(pPlayer.Level+1))
  TmpStr = StrLeft(TmpStr, StrFindFirstChar(TmpStr, '.'))
  Exp2 = FormatCommas(TmpStr)
  for StrGetLength(Exp1) < StrGetLength(Exp2) {
    Exp1 = StrInsertChar(Exp1, 0, ' ')
  }
  pPlayer.Output += "Experience:   "
  pPlayer.Output += Exp1
  pPlayer.Output += "\r\n"
  pPlayer.Output += "Next level:   "
  pPlayer.Output += Exp2
  pPlayer.Output += "\r\n"
  // Armor Class
  TmpStr = fmt.Sprintf("%d", pPlayer.ArmorClass)
  pPlayer.Output += "Armor Class:  "
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
  // Color
  if pPlayer.Color {
    TmpStr = "On"
  } else {
    TmpStr = "Off"
  }
  pPlayer.Output += "Color:        "
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
  // AllowGroup
  if pPlayer.AllowGroup {
    TmpStr = "On"
  } else {
    TmpStr = "Off"
  }
  pPlayer.Output += "Allow Group:  "
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
  // AllowAssist
  if pPlayer.AllowAssist {
    TmpStr = "On"
  } else {
    TmpStr = "Off"
  }
  pPlayer.Output += "Allow Assist: "
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
  // Position
  pPlayer.Output += "Position:     "
  TmpStr = pPlayer.Position
  StrReplace(&TmpStr, "s", "S")
  if TmpStr == "Sit" {
    TmpStr += "t"
  }
  TmpStr += "ing"
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
  // Silver
  TmpStr = fmt.Sprintf("%d", pPlayer.Silver)
  pPlayer.Output += "Silver:       "
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
  // Hunger
  TmpStr = fmt.Sprintf("%d", pPlayer.Hunger)
  pPlayer.Output += "Hunger:       "
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
  // Thirst
  TmpStr = fmt.Sprintf("%d", pPlayer.Thirst)
  pPlayer.Output += "Thirst:       "
  pPlayer.Output += TmpStr
  pPlayer.Output += "\r\n"
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
  DEBUGIT(5)
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
func PlayerRoomHasNotBeenHere(pPlayer *Player) bool {
  var Char       byte
  var CharPos    int
  var CharStr    string
  var RoomNbr    int
  var RoomNbrStr string

  if len(pPlayer.PlayerRoomVector) == 0 {
    PlayerRoomStringRead(pPlayer)
  }
  // Get RoomNbr from RoomId
  CharPos = StrGetLength(pPlayer.RoomId) - 1
  Char = StrGetAt(pPlayer.RoomId, CharPos)
  for Char >= '0' && Char <= '9' {
    CharStr = ""
    CharStr = CharStr + string(Char)
    RoomNbrStr = StrInsert(RoomNbrStr, 0, CharStr)
    CharPos--
    Char = StrGetAt(pPlayer.RoomId, CharPos)
  }
  RoomNbr = StrToInt(RoomNbrStr)
  // Has player been here?
  pPlayer.PlayerRoomCharPos = int(math.Ceil(float64(RoomNbr)/8.0)) - 1
  pPlayer.PlayerRoomBitPos = RoomNbr - (pPlayer.PlayerRoomCharPos*8) - 1
  pPlayer.PlayerRoomChar = pPlayer.PlayerRoomVector[pPlayer.PlayerRoomCharPos]
  PlayerRoomCharToBitsConvert(pPlayer)

  if pPlayer.PlayerRoomBits[pPlayer.PlayerRoomBitPos] {
    // Player has been here
    return false
  } else {
    // Player has not been here
    pPlayer.PlayerRoomBits[pPlayer.PlayerRoomBitPos] = true
    PlayerRoomBitsToCharConvert(pPlayer)
    pPlayer.PlayerRoomVector[pPlayer.PlayerRoomCharPos] = pPlayer.PlayerRoomChar
    PlayerRoomStringWrite(pPlayer)
    return true
  }
}

// Convert from PlayerRoom char to PlayerRoom bits 
func PlayerRoomCharToBitsConvert(pPlayer *Player) {
  var BitPos    int
  var Char      int
  var Remainder int

  Char = int(pPlayer.PlayerRoomChar)
  for i := 0; i < len(pPlayer.PlayerRoomBits); i++ {
    pPlayer.PlayerRoomBits[i] = false
  }
  for BitPos = 7; BitPos >= 0; BitPos-- {
    Remainder = int(Char) - int(math.Pow(2, float64(BitPos)))
    if Remainder > -1 {
      pPlayer.PlayerRoomBits[BitPos] = true
      Char = Remainder
    }
  }
}

// Convert from PlayerRoom bits to PlayerRoom char
func PlayerRoomBitsToCharConvert(pPlayer *Player) {
  var BitPos int
  var Char   int

  Char = 0
  for BitPos = 7; BitPos >= 0; BitPos-- {
    if pPlayer.PlayerRoomBits[BitPos] {
      Char += int(math.Pow(2, float64(BitPos)))
    }
  }
  pPlayer.PlayerRoomChar = byte(Char)
}

// Read PlayerRoomVector from disk
func PlayerRoomStringRead(pPlayer *Player) {
  var PlayerRoomString string
  var BitsetFileName string

  BitsetFileName = PLAYER_ROOM_DIR
  BitsetFileName += pPlayer.Name
  BitsetFileName += ".txt"

  if !FileExist(BitsetFileName) {
    for pPlayer.PlayerRoomCharPos = 0; pPlayer.PlayerRoomCharPos < MAX_ROOMS_CHAR; pPlayer.PlayerRoomCharPos++ {
      pPlayer.PlayerRoomVector = append(pPlayer.PlayerRoomVector, byte(0x00))
    }
    PlayerRoomStringWrite(pPlayer)
    return
  }

  pPlayer.PlayerRoomVector = []byte{}
  Data, _ := os.ReadFile(BitsetFileName)
  PlayerRoomString = string(Data)
  pPlayer.PlayerRoomVector = []byte(PlayerRoomString)
}

// Write PlayerRoomVector to disk
func PlayerRoomStringWrite(pPlayer *Player) {
  var BitsetFileName string

  BitsetFileName = PLAYER_ROOM_DIR
  BitsetFileName += pPlayer.Name
  BitsetFileName += ".txt"
  _ = os.WriteFile(BitsetFileName, pPlayer.PlayerRoomVector, 0o644)
}