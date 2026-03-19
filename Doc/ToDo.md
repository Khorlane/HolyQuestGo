# OMugs – Online Multi-User Game Server

# Done List

# License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any means.

For more information, please refer to <http://unlicense.org/>

# Things Done

- [x] HQ-001 - 12/15/02 - Implement player levels
- [x] HQ-002 - 12/17/02 - Make mobiles fight back
- [x] HQ-003 - 12/18/02 - Assist command new functionality: allow player to control whether or not they can be assisted
- [x] HQ-004 - 12/18/02 - Look \<player\> command
- [x] HQ-005 - 12/18/02 - Look \<mobile\> command
- [x] HQ-006 - 12/19/02 - Enhance 'where' command, add where \<object\> and where \<player\>
- [x] HQ-007 - 12/19/02 - Afk command, add afk player state, display on who listing
- [x] HQ-008 - 12/19/02 - Spawning mobiles now, needs work though
- [x] HQ-009 - 12/20/02 - Announce to room when mobile spawns
- [x] HQ-010 - 12/30/02 - Develop an event que system (mob spawn event system)
- [x] HQ-011 - 01/06/03 - Make mobiles move
- [x] HQ-012 - 01/08/03 - Make mobiles flagged as NoMove, Not move
- [x] HQ-013 - 01/15/03 - Create all spawn events, not just the first evemt and then skip out
- [x] HQ-014 - 01/30/03 - Disconnect players during logon if no response after a predefined number of ticks
- [x] HQ-015 - 01/31/03 - Advance a character to a given level for example: advance kwam 18
- [x] HQ-016 - 02/24/03 - Determine and log player's ip address
- [x] HQ-017 - 03/04/03 - If two players are fighting a mob and one player dies or flees, make the mob switch targets
- [x] HQ-018 - 03/05/03 - Allow player to change their password
- [x] HQ-019 - 03/07/03 - Allow player to delete themselves
- [x] HQ-020 - 03/13/03 - Shops
- [x] HQ-021 - 03/18/03 - Award loot to player(s) when mob dies
- [x] HQ-022 - 03/26/03 - Write help for kill command
- [x] HQ-023 - 03/27/03 - Implement MOTD
- [x] HQ-024 - 04/10/03 - Use armor in damage calculations
- [x] HQ-025 - 05/08/03 - Implement skills
- [x] HQ-026 - 05/12/03 - During a fight, show player and mobile health percentage
- [x] HQ-027 - 05/15/03 - When a player deletes themselves, delete PlayerObj and PlayerEqu
- [x] HQ-028 - 05/15/03 - Make mobs regen health, when full health remove 'dot' notation & delete stats
- [x] HQ-029 - 05/22/03 - Make aggro mobs attack when player enters the room
- [x] HQ-030 - 05/22/03 - Add 'all' or a count to sell command
- [x] HQ-031 - 05/28/03 - Get home directory instead of hard coding it
- [x] HQ-032 - 05/29/03 - Log when a player levels or is 'advanced' a level or levels
- [x] HQ-033 - 05/29/03 - When creating a new player and the player is disconnected after entering their password, but they haven't entered their sex, they can log back in, but 'sex' is blank.
- [x] HQ-034 - 05/30/03 - Write 'who is online' program
- [x] HQ-035 - 06/03/03 - Added 'hail' command
- [x] HQ-036 - 06/26/03 - Command to list all socials
- [x] HQ-037 - 06/26/03 - Command to list all commands
- [x] HQ-038 - 06/26/03 - Command to list all help topics
- [x] HQ-039 - 06/26/03 - Emote command
- [x] HQ-040 - 07/09/03 - Drink command, thirst
- [x] HQ-041 - 07/09/03 - Eat command, hunger
- [x] HQ-042 - 07/11/03 - Chat command
- [x] HQ-043 - 10/07/03 - Reworked experience gain formulas
- [x] HQ-044 - 10/09/03 - Consider command
- [x] HQ-045 - 10/14/03 - Invisible command
- [x] HQ-046 - 10/24/03 - Rewrote the 'world' section of the Admin doc
- [x] HQ-047 - 10/28/03 - Reworked the mobile moving code to make 'spawn' rooms first and to limit the number of mobiles the code tries to move at each 'tick'
- [x] HQ-048 - 11/05/03 - Dialogs for online mobile building
- [x] HQ-049 - 11/06/03 - Moved 'generate room' code into OMugs codebase
- [x] HQ-050 - 11/07/03 - Fixed the 'buy axe' bug
- [x] HQ-051 - 11/10/03 - Improved the 'train' command to allow a player to untrain a skill
- [x] HQ-052 - 11/11/03 - Merged 'line count' code into OMugs codebase
- [x] HQ-053 - 11/14/03 - Merged 'Osi - OMugs Script Interpreter' code into OMugs codebase
- [x] HQ-054 - 12/03/03 - Renamed IP connection functions to all start with 'Sock' and minor changes
- [x] HQ-055 - 12/10/03 - Merged 'room generator' into OMugs codebase
- [x] HQ-056 - 12/12/03 - Dialogs for online object building
- [x] HQ-057 - 12/15/03 - Drinking from a spring, a well, etc is implemeted
- [x] HQ-058 - 12/16/03 - Fix WhoIsOnline bug, needed to save player when afk status changes
- [x] HQ-059 - 12/17/03 - If the same command is entered 3 times in 1 second, send 'NO SPAMMING'
- [x] HQ-060 - 12/18/03 - Allow players to reconnect
- [x] HQ-061 - 01/08/04 - Date, time, calendar
- [x] HQ-062 - 06/11/04 - Added current/total hitpoints and movepoints to status display
- [x] HQ-063 - 06/11/04 - Enhanced fight messages to say 'REALLY HURT', 'SMACK DOWN', etc
- [x] HQ-064 - 12/18/25 - Converted C++ OMugs project to HolyQuestGo project
- [x] HQ-065 - 03/18/26 - Update AI context and instructions (commit: cd0f830)
- [x] HQ-066 - 03/18/26 - Enhance chat_context.txt (commit: 5febb37)
- [x] HQ-067 - 03/18/26 - Align Room.go with Room.cpp (commit: f14fe7f)
- [x] HQ-068 - 03/18/26 - Refine chat context for session continuity (commit: d33d229)
- [x] HQ-069 - 03/18/26 - Align Shop.go with Shop.cpp (commit: 5e81197)
- [x] HQ-070 - 03/18/26 - Document OMugs to HolyQuestGo file mapping (commit: 321ae3d)
- [x] HQ-071 - 03/18/26 - Reformat ToDo history and add Features backlog file (commit: b6d56fc)
- [x] HQ-072 - 03/18/26 - Align Social.go with Social.cpp (commit: 712c95d)
- [x] HQ-073 - 03/18/26 - Align BigDog.go with BigDog.cpp (commit: a15cbc8)
- [x] HQ-074 - 03/18/26 - Align Calendar.go with Calendar.cpp (commit: 3697ef6)
- [x] HQ-075 - 03/18/26 - Align Color.go with Color.h (commit: f10f908)
- [ ] HQ-076 - 03/18/26 - Align Communication.go with Communication.cpp
- [x] HQ-077 - 03/18/26 - Align Config.go with Config.h (commit: f78104a)
- [ ] HQ-078 - 03/18/26 - Align Descriptor.go with Descriptor.cpp
- [ ] HQ-079 - 03/18/26 - Align Dnode.go with Dnode.cpp
- [x] HQ-080 - 03/18/26 - Align Help.go with Help.cpp (commit: 3c7aabd)
- [x] HQ-081 - 03/18/26 - Align Log.go with Log.cpp (commit: 13bd13b)
- [ ] HQ-082 - 03/18/26 - Align Mobile.go with Mobile.cpp
- [ ] HQ-083 - 03/18/26 - Align Object.go with Object.cpp
- [ ] HQ-084 - 03/18/26 - Align Player.go with Player.cpp
- [x] HQ-085 - 03/18/26 - Align Utility.go with Utility.cpp (commit: c6a653c)
- [x] HQ-086 - 03/18/26 - Align Validate.go with Tools\\Validate.cpp (commit: 8e0f766)
- [ ] HQ-087 - 03/18/26 - Align Violence.go with Violence.cpp
- [x] HQ-088 - 03/18/26 - Temporary fix to World.go to remove compile error (commit: aba50ca)
- [ ] HQ-089 - 03/18/26 - Align World.go with World.cpp
- [x] HQ-090 - 03/18/26 - Fix ValErr - Object file 'BearHide' not found -------------------> Bear.txt (commit: pending)
- [ ] HQ-091 - 03/18/26 - Apply strategy used to correct HQ-090 correct similar issues in Validate.go
