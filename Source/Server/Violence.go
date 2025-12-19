//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Violence.go                                           *
// Usage:     Handles violence-related actions and interactions      *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

// CalcDamageToMobile calculates damage done to a mobile.
func CalcDamageToMobile(Damage int, WeaponSkill int) int {
  return 0
}

// CalcDamageToPlayer calculates damage done to a player.
func CalcDamageToPlayer(Damage int, PAC int) int {
  return 0
}

// CalcHealthPct returns a string representing HP percentage.
func CalcHealthPct(HitPoints int, HitPointsMax int) string {
  return ""
}

// GetMobileArmor returns the armor value for a mobile.
func GetMobileArmor(MobileId string) int {
  return 0
}

// GetMobileAttack returns the attack description for a mobile.
func GetMobileAttack(MobileId string) string {
  return ""
}

// GetMobileDamage returns the damage value for a mobile.
func GetMobileDamage(MobileId string) int {
  return 0
}

// GetMobileDesc1 returns the primary description for a mobile.
func GetMobileDesc1(MobileId string) string {
  return ""
}

// GetMobileExpPointsLevel returns the XP and level info for a mobile.
func GetMobileExpPointsLevel(MobileId string) string {
  return ""
}

// GetMobileHitPoints returns the hit points info for a mobile.
func GetMobileHitPoints(MobileId string) string {
  return ""
}

// GetMobileLoot returns loot information for a mobile.
func GetMobileLoot(MobileId string) string {
  return ""
}

// GetMobileRoom returns the RoomId for a mobile.
func GetMobileRoom(MobileId string) string {
  return ""
}

// GetMobPlayerMobileId returns the MobileId for a player's Nth mobile.
func GetMobPlayerMobileId(PlayerName string, i int) string {
  return ""
}

// GetPlayerMobMobileId returns the MobileId of the mobile a player is fighting.
func GetPlayerMobMobileId(PlayerName string) string {
  return ""
}

// WhackMobile applies damage to a mobile and returns a result string.
func WhackMobile(PlayerName string, DamageToMobile int, MobileDesc1 string, WeaponType string) string {
  return ""
}

// WhackPlayer applies damage to a player and returns a result string.
func WhackPlayer(MobileDesc1 string, MobileAttack string, DamageToPlayer int) string {
  return ""
}