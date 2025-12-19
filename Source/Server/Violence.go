//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Violence.go                                           *
// Usage:     Handles violence-related actions and interactions      *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server
import (
  "math"
  "fmt"
  "os"
  "bufio"
)

// Calculate damage to mobile
func CalcDamageToMobile(Damage int, WeaponSkill int) int {
  var CalcDamage  int
  var Deduction   int
  var Percent     int
  var SkillFactor int

  SkillFactor = int(float64(WeaponSkill) / 1.66)
  Percent = GetRandomNumber(PLAYER_DMG_PCT - SkillFactor)
  Deduction = int(math.Ceil(float64(Damage) * (float64(Percent) / 100.0)))
  CalcDamage = Damage - Deduction
  return CalcDamage
}

// Calculate damage to player
func CalcDamageToPlayer(Damage int, PAC int) int {
  var CalcDamage int
  var Deduction  int
  var Percent    int

  Percent = GetRandomNumber(PLAYER_DMG_PCT)
  Deduction = int(math.Ceil(float64(Damage) * (float64(Percent) / 100.0)))
  CalcDamage = Damage - Deduction
  CalcDamage = int(math.Floor(float64(CalcDamage) - (float64(PAC)*PACMN*float64(CalcDamage))))
  return CalcDamage
}

// Calculate damage to player
func CalcHealthPct(HitPoints int, HitPointsMax int) string {
  var HealthPct string
  var Percent   int

  if HitPoints < 1 {
    Percent = 0
  } else {
    Percent = CalcPct(HitPoints, HitPointsMax)
  }
  if Percent > 75 {
    HealthPct = "&C"
  } else if Percent > 50 {
    HealthPct = "&Y"
  } else if Percent > 25 {
    HealthPct = "&M"
  } else {
    HealthPct = "&R"
  }
  Buf = fmt.Sprintf("%3d", int(Percent))
  TmpStr = Buf
  HealthPct += TmpStr
  HealthPct += "&N"
  return HealthPct
}

// Get mobile Armor
func GetMobileArmor(MobileId string) int {
  var MobileArmor            int
  var MobStatsArmorFile     *os.File
  var MobStatsArmorFileName  string

  MobStatsArmorFileName = MOB_STATS_ARM_DIR
  MobStatsArmorFileName += MobileId
  MobStatsArmorFileName += ".txt"
  f, err := os.Open(MobStatsArmorFileName)
  if err != nil {
    MobileArmor = 0
    return MobileArmor
  }
  MobStatsArmorFile = f
  Stuff = ""
  scanner := bufio.NewScanner(MobStatsArmorFile)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  MobStatsArmorFile.Close()
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  MobileArmor = StrToInt(Stuff)
  return MobileArmor
}

// Get mobile Attack
func GetMobileAttack(MobileId string) string {
  var MobileAttack            string
  var MobStatsAttackFile     *os.File
  var MobStatsAttackFileName  string

  MobStatsAttackFileName = MOB_STATS_ATK_DIR
  MobStatsAttackFileName += MobileId
  MobStatsAttackFileName += ".txt"
  f, err := os.Open(MobStatsAttackFileName)
  if err != nil {
    LogIt("Violence::GetMobileAttack - Open MobStatsAttack file failed (read)")
    os.Exit(1)
  }
  MobStatsAttackFile = f
  Stuff = ""
  scanner := bufio.NewScanner(MobStatsAttackFile)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  MobStatsAttackFile.Close()
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  MobileAttack = Stuff
  return MobileAttack
}

// Get mobile Damage
func GetMobileDamage(MobileId string) int {
  var MobileDamage           int
  var MobStatsDamageFile    *os.File
  var MobStatsDamageFileName string

  MobStatsDamageFileName = MOB_STATS_DMG_DIR
  MobStatsDamageFileName += MobileId
  MobStatsDamageFileName += ".txt"
  f, err := os.Open(MobStatsDamageFileName)
  if err != nil {
    LogIt("Violence::GetMobileDamage - Open MobStatsDamageFile file failed (read)")
    os.Exit(1)
  }
  MobStatsDamageFile = f
  Stuff = ""
  scanner := bufio.NewScanner(MobStatsDamageFile)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  MobStatsDamageFile.Close()
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  MobileDamage = StrToInt(Stuff)
  return MobileDamage
}

// Get mobile Desc1
func GetMobileDesc1(MobileId string) string {
  var MobileDesc1            string
  var MobStatsDesc1File     *os.File
  var MobStatsDesc1FileName  string

  MobStatsDesc1FileName = MOB_STATS_DSC_DIR
  MobStatsDesc1FileName += MobileId
  MobStatsDesc1FileName += ".txt"
  f, err := os.Open(MobStatsDesc1FileName)
  if err != nil {
    LogIt("Violence::GetMobileDesc1 - Open MobStatsDesc1 file failed (read)")
    os.Exit(1)
  }
  MobStatsDesc1File = f
  Stuff = ""
  scanner := bufio.NewScanner(MobStatsDesc1File)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  MobStatsDesc1File.Close()
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  MobileDesc1 = Stuff
  return MobileDesc1
}

// Get mobile Experience Points
func GetMobileExpPointsLevel(MobileId string) string {
  var MobStatsExpPointsFileName string

  MobStatsExpPointsFileName = MOB_STATS_EXP_DIR
  MobStatsExpPointsFileName += MobileId
  MobStatsExpPointsFileName += ".txt"
  f, err := os.Open(MobStatsExpPointsFileName)
  if err != nil {
    LogIt("Violence::GetMobileExpPointsLevel - Open MobStatsExpPointsFile file failed (read)")
    os.Exit(1)
  }
  Stuff = ""
  scanner := bufio.NewScanner(f)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  f.Close()
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  return Stuff
}

// Get mobile Hit Points
func GetMobileHitPoints(MobileId string) string {
  var MobStatsHitPointsFileName string
  var MobHitPoints              string

  MobStatsHitPointsFileName = MOB_STATS_HPT_DIR
  MobStatsHitPointsFileName += MobileId
  MobStatsHitPointsFileName += ".txt"
  f, err := os.Open(MobStatsHitPointsFileName)
  if err != nil {
    LogIt("Violence::WhackMobile - Open MobStatsHitPointsFile file failed (read)")
    os.Exit(1)
  }
  Stuff = ""
  scanner := bufio.NewScanner(f)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  f.Close()
  MobHitPoints = Stuff
  return MobHitPoints
}

// Get mobile Loot
func GetMobileLoot(MobileId string) string {
  var MobileLoot            string
  var MobStatsLootFile     *os.File
  var MobStatsLootFileName  string

  MobStatsLootFileName = MOB_STATS_LOOT_DIR
  MobStatsLootFileName += MobileId
  MobStatsLootFileName += ".txt"
  f, err := os.Open(MobStatsLootFileName)
  if err != nil {
    LogIt("Violence::GetMobileLoot - Open MobStatsLoot file failed (read)")
    os.Exit(1)
  }
  MobStatsLootFile = f
  Stuff = ""
  scanner := bufio.NewScanner(MobStatsLootFile)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  MobStatsLootFile.Close()
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  MobileLoot = Stuff
  return MobileLoot
}

// Get mobile Room
func GetMobileRoom(MobileId string) string {
  var MobileRoom            string
  var MobStatsRoomFile     *os.File
  var MobStatsRoomFileName  string

  MobStatsRoomFileName = MOB_STATS_ROOM_DIR
  MobStatsRoomFileName += MobileId
  MobStatsRoomFileName += ".txt"
  f, err := os.Open(MobStatsRoomFileName)
  if err != nil {
    LogIt("Violence::GetMobileRoom - Open MobStatsRoom file failed (read)")
    os.Exit(1)
  }
  MobStatsRoomFile = f
  Stuff = ""
  scanner := bufio.NewScanner(MobStatsRoomFile)
  if scanner.Scan() {
    Stuff = scanner.Text()
  }
  MobStatsRoomFile.Close()
  Stuff = StrTrimLeft(Stuff)
  Stuff = StrTrimRight(Stuff)
  MobileRoom = Stuff
  return MobileRoom
}

// Get MobPlayer MobileId
func GetMobPlayerMobileId(PlayerName string, i int) string {
  var j                  int
  var MobileId           string
  var MobPlayerFile     *os.File
  var MobPlayerFileName  string

  MobPlayerFileName = MOB_PLAYER_DIR
  MobPlayerFileName += PlayerName
  MobPlayerFileName += ".txt"
  f, err := os.Open(MobPlayerFileName)
  if err != nil {
    MobileId = "No more mobiles"
    return MobileId
  }
  MobPlayerFile = f
  MobileId = ""
  scanner := bufio.NewScanner(MobPlayerFile)
  for j = 1; j <= i; j++ {
    if scanner.Scan() {
      MobileId = scanner.Text()
    } else {
      break
    }
  }
  MobPlayerFile.Close()
  MobileId = StrTrimLeft(MobileId)
  MobileId = StrTrimRight(MobileId)
  if MobileId == "" {
    MobileId = "No more mobiles"
  }
  return MobileId
}

// Get PlayerMob MobileId
func GetPlayerMobMobileId(PlayerName string) string {
  var MobileId         string
  var PlayerMobFile   *os.File
  var PlayerMobFileName string

  PlayerMobFileName = PLAYER_MOB_DIR
  PlayerMobFileName += PlayerName
  PlayerMobFileName += ".txt"
  f, err := os.Open(PlayerMobFileName)
  if err != nil {
    LogIt("Violence::GetPlayerMobMobileId - Open PlayerMob file failed")
    os.Exit(1)
  }
  PlayerMobFile = f
  MobileId = ""
  scanner := bufio.NewScanner(PlayerMobFile)
  if scanner.Scan() {
    MobileId = scanner.Text()
  }
  PlayerMobFile.Close()
  MobileId = StrTrimLeft(MobileId)
  MobileId = StrTrimRight(MobileId)
  return MobileId
}

// Whack the mobile - do some damage!
func WhackMobile(MobileId string, DamageToMobile int, MobileDesc1 string, WeaponType string) string {
  var DamageAmount               string
  var DamageMagnitude            int
  var ExtraDamageMsg             string
  var MobileBeenWhacked          string
  var MobHealthPct               string
  var MobHealthPctNew            int
  var MobHealthPctOld            int
  var MobHitPointsInfo           string
  var MobHitPointsLeft           int
  var MobHitPointsTotal          int
  var MobStatsHitPointsFile     *os.File
  var MobStatsHitPointsFileName  string
  var WeaponAction               string

  MobHitPointsInfo  = GetMobileHitPoints(MobileId)
  MobHitPointsTotal = StrToInt(StrGetWord(MobHitPointsInfo, 1))
  MobHitPointsLeft  = StrToInt(StrGetWord(MobHitPointsInfo, 2))
  MobHealthPctOld   = CalcPct(MobHitPointsLeft, MobHitPointsTotal)
  MobHitPointsLeft  = MobHitPointsLeft - DamageToMobile
  if MobHitPointsLeft < 0 {
    MobHitPointsLeft = 0
  }
  MobStatsHitPointsFileName = MOB_STATS_HPT_DIR
  MobStatsHitPointsFileName += MobileId
  MobStatsHitPointsFileName += ".txt"
  f, err := os.OpenFile(MobStatsHitPointsFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
  if err != nil {
    LogIt("Violence::WhackMobile - Open MobStatsHitPointsFile file failed (write)")
    os.Exit(1)
  }
  MobStatsHitPointsFile = f
  Buf    = fmt.Sprintf("%d", MobHitPointsTotal)
  TmpStr = Buf
  Stuff  = TmpStr
  Stuff += " "
  Buf    = fmt.Sprintf("%d", MobHitPointsLeft)
  TmpStr = Buf
  Stuff += TmpStr
  fmt.Fprintln(MobStatsHitPointsFile, Stuff)
  MobStatsHitPointsFile.Close()
  MobHealthPct    = CalcHealthPct(MobHitPointsLeft, MobHitPointsTotal)
  WeaponType      = StrMakeLower(WeaponType)
  WeaponAction    = TranslateWord(WeaponType)
  Buf             = fmt.Sprintf("%d", DamageToMobile)
  DamageAmount    = Buf
  MobHealthPctNew = CalcPct(MobHitPointsLeft, MobHitPointsTotal)
  DamageMagnitude = MobHealthPctOld - MobHealthPctNew
  if DamageMagnitude > 70 {
    ExtraDamageMsg = "&RALMOST OBLITERATE&N"
  } else if DamageMagnitude > 60 {
    ExtraDamageMsg = "&RSMACK DOWN&N"
  } else if DamageMagnitude > 50 {
    ExtraDamageMsg = "&MSEVERLY WOUND&N"
  } else if DamageMagnitude > 40 {
    ExtraDamageMsg = "&YREALLY HURT&N"
  }
  if MobHitPointsLeft > 0 {
    MobileBeenWhacked = "alive"
    MobileBeenWhacked += " "
    MobileBeenWhacked += MobHealthPct
    MobileBeenWhacked += " "
    if DamageMagnitude > 40 {
      MobileBeenWhacked += "You"
      MobileBeenWhacked += " "
      MobileBeenWhacked += ExtraDamageMsg
      MobileBeenWhacked += " "
      MobileBeenWhacked += MobileDesc1
      MobileBeenWhacked += " "
      MobileBeenWhacked += "with your"
      MobileBeenWhacked += " "
      MobileBeenWhacked += WeaponAction
      MobileBeenWhacked += " "
      MobileBeenWhacked += "of"
      MobileBeenWhacked += " "
      MobileBeenWhacked += DamageAmount
      MobileBeenWhacked += " "
      MobileBeenWhacked += "points of damage."
    } else {
      MobileBeenWhacked += "You"
      MobileBeenWhacked += " "
      MobileBeenWhacked += WeaponAction
      MobileBeenWhacked += " "
      MobileBeenWhacked += MobileDesc1
      MobileBeenWhacked += " "
      MobileBeenWhacked += "for"
      MobileBeenWhacked += " "
      MobileBeenWhacked += DamageAmount
      MobileBeenWhacked += " "
      MobileBeenWhacked += "points of damage."
    }
  } else {
    MobileBeenWhacked = "dead"
    MobileBeenWhacked += " "
    if MobHealthPctOld == 100 && MobHealthPctNew == 0 {
      MobileBeenWhacked += "You"
      MobileBeenWhacked += " "
      MobileBeenWhacked += "&ROBLITERATED&N"
      MobileBeenWhacked += " "
      MobileBeenWhacked += MobileDesc1
      MobileBeenWhacked += " "
      MobileBeenWhacked += "with a"
      MobileBeenWhacked += " "
      MobileBeenWhacked += WeaponAction
      MobileBeenWhacked += " "
      MobileBeenWhacked += "that did"
      MobileBeenWhacked += " "
      MobileBeenWhacked += DamageAmount
      MobileBeenWhacked += " "
      MobileBeenWhacked += "points of damage."
    } else {
      MobileBeenWhacked += "You"
      MobileBeenWhacked += " "
      MobileBeenWhacked += "vanquish"
      MobileBeenWhacked += " "
      MobileBeenWhacked += MobileDesc1
      MobileBeenWhacked += " "
      MobileBeenWhacked += "with a"
      MobileBeenWhacked += " "
      MobileBeenWhacked += WeaponAction
      MobileBeenWhacked += " "
      MobileBeenWhacked += "that did"
      MobileBeenWhacked += " "
      MobileBeenWhacked += DamageAmount
      MobileBeenWhacked += " "
      MobileBeenWhacked += "points of damage."
    }
  }
  return MobileBeenWhacked
}

// Whack the player - do some damage!
func WhackPlayer(MobileDesc1 string, MobileAttack string, DamageToPlayer int) string {
  var PlayerBeenWhacked string

  MobileDesc1 = StrMakeFirstUpper(MobileDesc1)
  Buf = fmt.Sprintf("%d", DamageToPlayer)
  TmpStr = Buf
  PlayerBeenWhacked = MobileDesc1
  PlayerBeenWhacked += " "
  PlayerBeenWhacked += MobileAttack
  PlayerBeenWhacked += " "
  PlayerBeenWhacked += "you for"
  PlayerBeenWhacked += " "
  PlayerBeenWhacked += TmpStr
  PlayerBeenWhacked += " "
  PlayerBeenWhacked += "points of damage."
  return PlayerBeenWhacked
}