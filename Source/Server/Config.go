//*******************************************************************
// HolyQuest: Powered by the Online Multi-User Game Server (OMugs)  *
// File:      Config.go                                             *
// Usage:     Sets up configuration constants for the game          *
// This file is part of the HolyQuestGo project.                    *
// It is licensed under the Unlicense.                              *
// See the LICENSE file in the project root for details.            *
//*******************************************************************

package server

import (
  "runtime"
)

// Configuration constants
const (
  // Version information
  VERSION = "2025.11.28"

  // Debugging level
  DEBUG_LVL = 1 // 0 to turn off, 1 or more to turn on

  // Game mechanics constants
  MAC                     float64 = 300.0 // Maximum Armor Class
  MDRP                    float64 = 60.0  // Maximum Damage Reduction Percent
  MGBP                    float64 = 80.0  // Maximum Group Bonus Percent
  MGRP                    float64 = 50.0  // Maximum Gain Reduction Percent
  MOB_ARM_PER_LEVEL       int     = 5     // Currently not used
  MOB_DMG_PER_LEVEL       int     = 9     // Damage per level
  MOB_EXP_PER_LEVEL       int     = 50    // Experience per level
  MOB_HPT_PER_LEVEL       int     = 31    // Hit points per level
  MOB_MOVE_PCT            int     = 30    // Percent chance mob will move
  PLAYER_DMG_HAND         int     = 5     // Bare hand damage
  PLAYER_DMG_PCT          int     = 75    // Random damage deduction
  PLAYER_EXP_PER_LEVEL    int     = 1000  // Experience per level factor
  PLAYER_HPT_PER_LEVEL    int     = 31    // Hit points per level
  PLAYER_LOOSES_EXP_LEVEL int     = 5     // Level when death causes XP loss
  PLAYER_SKILL_PER_LEVEL  int     = 3     // Skill points gained per level

  // Tick durations
  MILLI_SECONDS_TO_SLEEP  int     = 100   // Defines tick duration: 100 = 10 ticks per second
  EVENT_TICK              int     = 50    // 5 seconds
  FIGHT_TICK              int     = 17    // 1.7 seconds
  HUNGER_THIRST_TICK      int     = 540   // 54 seconds
  INPUT_TICK              int     = 300   // 30 seconds
  MOB_HEAL_TICK           int     = 900   // 90 seconds
  STATS_TICK              int     = 101   // 10.1 seconds
  WHO_IS_ONLINE_TICK      int     = 165   // 16.5 seconds

  // Time constants
  REAL_MINUTES_PER_HOUR   int     = 4     // Game time advances 1 hour every 4 minutes
  HOURS_PER_DAY           int     = 24    // Game hours in each game day
  DAYS_PER_WEEK           int     = 7     // Game days in each game week
  DAYS_PER_MONTH          int     = 30    // Game weeks in each game month
  MONTHS_PER_YEAR         int     = 12    // Game months in each game year

  // Miscellaneous constants
  GRP_LIMIT               int     = 4     // Maximum group members
  HPT_GAIN_STAND          int     = 2     // Hit points gained while standing
  HPT_GAIN_SIT            int     = 5     // Hit points gained while sitting
  HPT_GAIN_SLEEP          int     = 9     // Hit points gained while sleeping
  HUNGER_THIRST_LEVEL     int     = 3     // No hunger or thirst until this level
  MAX_INPUT_LENGTH        int     = 1024  // Maximum string length received from player
  MAX_ROOMS               int     = 8192  // Must be a multiple of 8
  MAX_ROOMS_CHAR          int     = MAX_ROOMS / 8 // 8 rooms represented by each byte
  SAFE_ROOM               string  = "JesseSquare8"
  START_ROOM              string  = "Welcome226"
  UNTRAIN_COST            string  = "250" // Cost to untrain a skill
  PORT_NBR                int     = 7373 // Listening on this port for connections

  // Directory paths
  GREETING_DIR            string  = "./Library/"
  HELP_DIR                string  = "./Library/"
  MOTD_DIR                string  = "./Library/"
  SOCIAL_DIR              string  = "./Library/"
  VALID_CMDS_DIR          string  = "./Library/"
  VALID_NAMES_DIR         string  = "./Library/"
  DAY_NAMES_DIR           string  = "./Library/Calendar/"
  DAY_OF_MONTH_DIR        string  = "./Library/Calendar/"
  HOUR_NAMES_DIR          string  = "./Library/Calendar/"
  MONTH_NAMES_DIR         string  = "./Library/Calendar/"
  LOOT_DIR                string  = "./Library/Loot/"
  MOBILES_DIR             string  = "./Library/Mobiles/"
  OBJECTS_DIR             string  = "./Library/Objects/"
  ROOMS_DIR               string  = "./Library/Rooms/"
  SCRIPTS_DIR             string  = "./Library/Scripts/"
  SHOPS_DIR               string  = "./Library/Shops/"
  SQL_DIR                 string  = "./Library/Sql/"
  TALK_DIR                string  = "./Library/Talk/"
  WORLD_MAP_DIR           string  = "./Library/World/Map/"
  WORLD_MOBILES_DIR       string  = "./Library/World/Mobiles/"
  CONTROL_DIR             string  = "./Running/Control/"
  CONTROL_EVENTS_DIR      string  = "./Running/Control/Events/"
  CONTROL_MOB_INWORLD_DIR string  = "./Running/Control/Mobiles/InWorld/"
  CONTROL_MOB_NOMOVE_DIR  string  = "./Running/Control/Mobiles/NoMove/"
  CONTROL_MOB_SPAWN_DIR   string  = "./Running/Control/Mobiles/Spawn/"
  LOG_DIR                 string  = "./Running/Log/"
  PLAYER_DIR              string  = "./Running/Players/"
  PLAYER_EQU_DIR          string  = "./Running/Players/PlayerEqu/"
  PLAYER_OBJ_DIR          string  = "./Running/Players/PlayerObj/"
  PLAYER_ROOM_DIR         string  = "./Running/Players/PlayerRoom/"
  ROOM_MOB_DIR            string  = "./Running/RoomMob/"
  ROOM_OBJ_DIR            string  = "./Running/RoomObj/"
  MOB_PLAYER_DIR          string  = "./Running/Violence/MobPlayer/"
  PLAYER_MOB_DIR          string  = "./Running/Violence/PlayerMob/"
  MOB_STATS_ARM_DIR       string  = "./Running/Violence/MobStats/Armor/"
  MOB_STATS_ATK_DIR       string  = "./Running/Violence/MobStats/Attack/"
  MOB_STATS_DMG_DIR       string  = "./Running/Violence/MobStats/Damage/"
  MOB_STATS_DSC_DIR       string  = "./Running/Violence/MobStats/Desc1/"
  MOB_STATS_EXP_DIR       string  = "./Running/Violence/MobStats/ExpPoints/"
  MOB_STATS_HPT_DIR       string  = "./Running/Violence/MobStats/HitPoints/"
  MOB_STATS_LOOT_DIR      string  = "./Running/Violence/MobStats/Loot/"
  MOB_STATS_ROOM_DIR      string  = "./Running/Violence/MobStats/Room/"
  DOC_DIR                 string  = "./Doc/"
  SOURCE_DIR              string  = "./Source/"
  TMP_DIR                 string  = "./Tmp/"
  UTILITY_DIR             string  = "./Utility/"
  WEB_SITE_DIR            string  = "./WebSite/"
)

// Variables
var (
  Buf                     string
  CmdStr                  string
  CurrentLineNumber       uint8
  ErrorCode               error
  HomeDir                 string
  LogBuf                  string
  MudCmd                  string
  PACMN                   float64 // Percent Armor Class Magic Number
  RoomID                  string
  ScriptFileName          string
  StateConnections        bool
  StateRunning            bool
  StateStopping           bool
  Stuff                   string
  TmpStr                  string
  ValErr                  bool
)

// DebugIt logs a message if the debug level is greater than or equal to the specified level.
func DEBUGIT(level int) {
  if DEBUG_LVL >= level {
    pc, _, _, _ := runtime.Caller(1) // Get the caller's program counter
    functionName := runtime.FuncForPC(pc).Name()
    LogIt(functionName)
  }
}