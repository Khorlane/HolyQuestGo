### Commands

- Light command

- Shout command

- Transfer command (admin)

### World

- Doors

### Connections

- Enable banning of ip addresses

- Add command to disconnect players

- Sometimes player ain't there but game thinks they are, connection is
  still there, no player on the other end. Kill the connection, be sure
  player file is set to Online:No

### Mobiles

- Make Function Mobile::CountMob count wounded mobs

- Make wimpy mobs flee when almost dead

### Players

- When a player dies, reduce move points to zero

- Allow players to create a description of themselves

- Player CRUD dialog

### Objects

- Containers

- Clean up 'pObject = new Object' code, what if \<object\>.txt doesn't
  exist for some reason. A pointer to object is still returned???

### Shops

- Add shop specific messages for buy / sell success / failure like
  BuyNotExist BuySuccess BuyNotAfford BuyNotAfford SellSuccess
  SellNotExist

### Misc

- Externalize starting hitpoints, movepoints, starting room, greeting,
  etc. In other words, most if not all config.h stuff

- Complete move points implemetation

- Design and write weather system

- Design and write quest system

- Send message to all players when the watch changes

OMugs Done List

- 12/15/02 - Implement player levels

- 12/17/02 - Make mobiles fight back

- 12/18/02 - Assist command new functionality: allow player to control
  whether or not they can be assisted

- 12/18/02 - Look \<player\> command

- 12/18/02 - Look \<mobile\> command

- 12/19/02 - Enhance 'where' command, add where \<object\> and where
  \<player\>

- 12/19/02 - Afk command, add afk player state, display on who listing

- 12/19/02 - Spawning mobiles now, needs work though

- 12/20/02 - Announce to room when mobile spawns

- 12/30/02 - Develop an event que system (mob spawn event system)

- 01/06/03 - Make mobiles move

- 01/08/03 - Make mobiles flagged as NoMove, Not move

- 01/15/03 - Create all spawn events, not just the first evemt and then
  skip out

- 01/30/03 - Disconnect players during logon if no response after a
  predefined number of ticks

- 01/31/03 - Advance a character to a given level for example: advance
  kwam 18

- 02/24/03 - Determine and log player's ip address

- 03/04/03 - If two players are fighting a mob and one player dies or
  flees, make the mob switch targets

- 03/05/03 - Allow player to change their password

- 03/07/03 - Allow player to delete themselves

- 03/13/03 - Shops

- 03/18/03 - Award loot to player(s) when mob dies

- 03/26/03 - Write help for kill command

- 03/27/03 - Implement MOTD

- 04/10/03 - Use armor in damage calculations

- 05/08/03 - Implement skills

- 05/12/03 - During a fight, show player and mobile health percentage

- 05/15/03 - When a player deletes themselves, delete PlayerObj and
  PlayerEqu

- 05/15/03 - Make mobs regen health, when full health remove 'dot'
  notation & delete stats

<!-- -->

- 05/22/03 - Make aggro mobs attack when player enters the room

- 05/22/03 - Add 'all' or a count to sell command

- 05/28/03 - Get home directory instead of hard coding it

<!-- -->

- 05/29/03 - Log when a player levels or is 'advanced' a level or levels

- 05/29/03 - When creating a new player and the player is disconnected
  after entering their password, but they haven't entered their sex,
  they can log back in, but 'sex' is blank.

- 05/30/03 - Write 'who is online' program

- 06/03/03 - Added 'hail' command

- 06/26/03 - Command to list all socials

- 06/26/03 - Command to list all commands

- 06/26/03 - Command to list all help topics

- 06/26/03 - Emote command

- 07/09/03 - Drink command, thirst

- 07/09/03 - Eat command, hunger

- 07/11/03 - Chat command

- 10/07/03 - Reworked experience gain formulas

- 10/09/03 - Consider command

- 10/14/03 - Invisible command

- 10/24/03 - Rewrote the 'world' section of the Admin doc

- 10/28/03 - Reworked the mobile moving code to make 'spawn' rooms first
  and to limit the number of mobiles the code tries to move at each
  'tick'

- 11/05/03 - Dialogs for online mobile building

- 11/06/03 - Moved 'generate room' code into OMugs codebase

- 11/07/03 - Fixed the 'buy axe' bug

- 11/10/03 - Improved the 'train' command to allow a player to untrain a
  skill

- 11/11/03 - Merged 'line count' code into OMugs codebase

- 11/14/03 - Merged 'Osi - OMugs Script Interpreter' code into OMugs
  codebase

- 12/03/03 - Renamed IP connection functions to all start with 'Sock'
  and minor changes

- 12/10/03 - Merged 'room generator' into OMugs codebase

- 12/12/03 - Dialogs for online object building

- 12/15/03 - Drinking from a spring, a well, etc is implemeted

- 12/16/03 - Fix WhoIsOnline bug, needed to save player when afk status
  changes

- 12/17/03 - If the same command is entered 3 times in 1 second, send
  'NO SPAMMING'

- 12/18/03 - Allow players to reconnect

- 01/08/04 - Date, time, calendar

- 06/11/04 - Added current/total hitpoints and movepoints to status
  display

- 06/11/04 - Enhanced fight messages to say 'REALLY HURT', 'SMACK DOWN',
  etc
