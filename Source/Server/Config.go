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
	MAC                    = 300.0 // Maximum Armor Class
	MDRP                   = 60.0  // Maximum Damage Reduction Percent
	MGBP                   = 80.0  // Maximum Group Bonus Percent
	MGRP                   = 50.0  // Maximum Gain Reduction Percent
	MOB_ARM_PER_LEVEL      = 5     // Currently not used
	MOB_DMG_PER_LEVEL      = 9     // Damage per level
	MOB_EXP_PER_LEVEL      = 50    // Experience per level
	MOB_HPT_PER_LEVEL      = 31    // Hit points per level
	MOB_MOVE_PCT           = 30    // Percent chance mob will move
	PLAYER_DMG_HAND        = 5     // Bare hand damage
	PLAYER_DMG_PCT         = 75    // Random damage deduction
	PLAYER_EXP_PER_LEVEL   = 1000  // Experience per level factor
	PLAYER_HPT_PER_LEVEL   = 31    // Hit points per level
	PLAYER_LOOSES_EXP_LEVEL = 5    // Level when death causes XP loss
	PLAYER_SKILL_PER_LEVEL = 3     // Skill points gained per level

	// Tick durations
	MILLI_SECONDS_TO_SLEEP = 100   // Defines tick duration: 100 = 10 ticks per second
	EVENT_TICK             = 50    // 5 seconds
	FIGHT_TICK             = 17    // 1.7 seconds
	HUNGER_THIRST_TICK     = 540   // 54 seconds
	INPUT_TICK             = 300   // 30 seconds
	MOB_HEAL_TICK          = 900   // 90 seconds
	STATS_TICK             = 101   // 10.1 seconds
	WHO_IS_ONLINE_TICK     = 165   // 16.5 seconds

	// Time constants
	REAL_MINUTES_PER_HOUR = 4  // Game time advances 1 hour every 4 minutes
	HOURS_PER_DAY         = 24 // Game hours in each game day
	DAYS_PER_WEEK         = 7  // Game days in each game week
	DAYS_PER_MONTH        = 30 // Game weeks in each game month
	MONTHS_PER_YEAR       = 12 // Game months in each game year

	// Miscellaneous constants
	GRP_LIMIT           = 4   // Maximum group members
	HPT_GAIN_STAND      = 2   // Hit points gained while standing
	HPT_GAIN_SIT        = 5   // Hit points gained while sitting
	HPT_GAIN_SLEEP      = 9   // Hit points gained while sleeping
	HUNGER_THIRST_LEVEL = 3   // No hunger or thirst until this level
	MAX_INPUT_LENGTH    = 1024 // Maximum string length received from player
	MAX_ROOMS           = 8192 // Must be a multiple of 8
	MAX_ROOMS_CHAR      = MAX_ROOMS / 8 // 8 rooms represented by each byte
	SAFE_ROOM           = "JesseSquare8"
	START_ROOM          = "Welcome226"
	UNTRAIN_COST        = "250" // Cost to untrain a skill
	PORT_NBR             = 7373 // Listening on this port for connections

	// Directory paths
	GREETING_DIR           = "./Library/"
	HELP_DIR               = "./Library/"
	MOTD_DIR               = "./Library/"
	SOCIAL_DIR             = "./Library/"
	VALID_CMDS_DIR         = "./Library/"
	VALID_NAMES_DIR        = "./Library/"
	DAY_NAMES_DIR          = "./Library/Calendar/"
	DAY_OF_MONTH_DIR       = "./Library/Calendar/"
	HOUR_NAMES_DIR         = "./Library/Calendar/"
	MONTH_NAMES_DIR        = "./Library/Calendar/"
	LOOT_DIR               = "./Library/Loot/"
	MOBILES_DIR            = "./Library/Mobiles/"
	OBJECTS_DIR            = "./Library/Objects/"
	ROOMS_DIR              = "./Library/Rooms/"
	SCRIPTS_DIR            = "./Library/Scripts/"
	SHOPS_DIR              = "./Library/Shops/"
	SQL_DIR                = "./Library/Sql/"
	TALK_DIR               = "./Library/Talk/"
	WORLD_MAP_DIR          = "./Library/World/Map/"
	WORLD_MOBILES_DIR      = "./Library/World/Mobiles/"
	CONTROL_DIR            = "./Running/Control/"
	CONTROL_EVENTS_DIR     = "./Running/Control/Events/"
	CONTROL_MOB_INWORLD_DIR = "./Running/Control/Mobiles/InWorld/"
	CONTROL_MOB_NOMOVE_DIR = "./Running/Control/Mobiles/NoMove/"
	CONTROL_MOB_SPAWN_DIR  = "./Running/Control/Mobiles/Spawn/"
	LOG_DIR                = "./Running/Log/"
	PLAYER_DIR             = "./Running/Players/"
	PLAYER_EQU_DIR         = "./Running/Players/PlayerEqu/"
	PLAYER_OBJ_DIR         = "./Running/Players/PlayerObj/"
	PLAYER_ROOM_DIR        = "./Running/Players/PlayerRoom/"
	ROOM_MOB_DIR           = "./Running/RoomMob/"
	ROOM_OBJ_DIR           = "./Running/RoomObj/"
	MOB_PLAYER_DIR         = "./Running/Violence/MobPlayer/"
	PLAYER_MOB_DIR         = "./Running/Violence/PlayerMob/"
	MOB_STATS_ARM_DIR      = "./Running/Violence/MobStats/Armor/"
	MOB_STATS_ATK_DIR      = "./Running/Violence/MobStats/Attack/"
	MOB_STATS_DMG_DIR      = "./Running/Violence/MobStats/Damage/"
	MOB_STATS_DSC_DIR      = "./Running/Violence/MobStats/Desc1/"
	MOB_STATS_EXP_DIR      = "./Running/Violence/MobStats/ExpPoints/"
	MOB_STATS_HPT_DIR      = "./Running/Violence/MobStats/HitPoints/"
	MOB_STATS_LOOT_DIR     = "./Running/Violence/MobStats/Loot/"
	MOB_STATS_ROOM_DIR     = "./Running/Violence/MobStats/Room/"
	DOC_DIR                = "./Doc/"
	SOURCE_DIR             = "./Source/"
	TMP_DIR                = "./Tmp/"
	UTILITY_DIR            = "./Utility/"
	WEB_SITE_DIR           = "./WebSite/"
)

// Variables
var (
	Buf               string
	CmdStr            string
	CurrentLineNumber uint8
	ErrorCode         error
	HomeDir           string
	LogBuf            string
	MudCmd            string
	PACMN             float32 // Percent Armor Class Magic Number
	RoomID						string
	ScriptFileName    string
	StateConnections  bool
	StateRunning      bool
	StateStopping     bool
	Stuff             string
	TmpStr            string
	ValErr            bool
)

// DebugIt logs a message if the debug level is greater than or equal to the specified level.
func DEBUGIT(level int) {
	if DEBUG_LVL >= level {
		pc, _, _, _ := runtime.Caller(1) // Get the caller's program counter
		functionName := runtime.FuncForPC(pc).Name()
		LogIt(functionName)
	}
}