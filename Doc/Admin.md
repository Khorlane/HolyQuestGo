# OMugs – Online Multi-User Game Server
# Administration Guide

## Table of Contents

- [Introduction](#introduction)
    - [Terminology](#terminology)
    - [Creating a virtual world](#creating-a-virtual-world)
    - [Critical virtual world components](#critical-virtual-world-components)
    - [Other virtual world components](#other-virtual-world-components)
    - [The library](#the-library)
    - [The library files](#the-library-files)
- [Rooms](#rooms)
    - [Room specification](#room-specification)
    - [A note about room exits](#a-note-about-room-exits)
    - [Room example 1](#room-example-1)
    - [Room example 2](#room-example-2)
- [Doors](#doors)
    - [Door specification](#door-specification)
    - [Door side 1 example](#door-side-1-example)
    - [Door side 2 example](#door-side-2-example)
    - [Valid door states and transitions](#valid-door-states-and-transitions)
        - [Door states](#door-states)
        - [Door state transitions](#door-state-transitions)
        - [Current door status](#current-door-status)
- [NPCs](#npcs)
    - [NPC specification](#npc-specification)
    - [NPC example 1](#npc-example-1)
    - [NPC example 2](#npc-example-2)
- [Objects](#objects)
    - [Object specification](#object-specification)
    - [Object type specifications](#object-type-specifications)
    - [Object examples](#object-examples)
- [Shops](#shops)
    - [Shop specification](#shop-specification)
    - [Shop example](#shop-example)
- [Loot](#loot)
    - [Loot specification](#loot-specification)
    - [Loot example](#loot-example)
- [Talk](#talk)
    - [Talk specification](#talk-specification)
    - [Talk example](#talk-example)
- [World](#world)
    - [NPC spawn specifications](#npc-spawn-specifications)
    - [NPC spawn example](#npc-spawn-example)
    - [Object spawn specifications](#object-spawn-specifications)
- [Socials](#socials)
    - [Social specification](#social-specification)
    - [Social examples](#social-examples)
- [Help](#help)
    - [Help specification](#help-specification)
    - [Help example](#help-example)
- [Player commands](#player-commands)
    - [Afk](#afk)
    - [Assist](#assist)
    - [Buy](#buy)
    - [Chat](#chat)
    - [Close](#close)
    - [Color](#color)
    - [Consider](#consider)
    - [Delete](#delete)
    - [Drink](#drink)
    - [Destroy](#destroy)
    - [Drop](#drop)
    - [Eat](#eat)
    - [Emote](#emote)
    - [Equipment](#equipment)
    - [Examine](#examine)
    - [Flee](#flee)
    - [Follow](#follow)
    - [Get](#get)
    - [Give](#give)
    - [Go](#go)
    - [Group](#group)
    - [Gsay](#gsay)
    - [Hail](#hail)
    - [Help](#help-1)
    - [Inventory](#inventory)
    - [Kill](#kill)
    - [List](#list)
    - [Lock](#lock)
    - [Look](#look)
    - [Money](#money)
    - [Motd](#motd)
    - [Open](#open)
    - [Password](#password)
    - [Played](#played)
    - [Quit](#quit)
    - [Remove](#remove)
    - [Save](#save)
    - [Say](#say)
    - [Sell](#sell)
    - [Show](#show)
    - [Sit](#sit)
    - [Sleep](#sleep)
    - [Stand](#stand)
    - [Status](#status)
    - [Tell](#tell)
    - [Time](#time)
    - [Title](#title)
    - [Train](#train)
    - [Unlock](#unlock)
    - [Wear](#wear)
    - [Wake](#wake)
    - [Who](#who)
    - [Wield](#wield)
- [Administrator commands](#administrator-commands)
    - [Advance](#advance)
    - [GoTo](#goto)
    - [GoToArrive](#gotoarrive)
    - [GoToDepart](#gotodepart)
    - [Invisible](#invisible)
    - [Load](#load)
    - [Refresh](#refresh)
    - [Restore](#restore)
    - [RoomInfo](#roominfo)
    - [Stop](#stop)
    - [Where](#where)
- [World building](#world-building)
    - [Armor class](#armor-class)
    - [Damage Reduction Percentage calculation](#damage-reduction-percentage-calculation)
    - [Combat example](#combat-example)

# License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any means.

For more information, please refer to <http://unlicense.org/>

Author  
Stephen L Bryant  

Revision 1.0 April    22, 2001  
Revision 1.1 December  4, 2002  
Revision 1.2 October  23, 2003  
Revision 1.3 December 22, 2025

Revision tracking after December 22, 2025 is maintained via git.

# Introduction

## Terminology

A Virtual World is a representation of the real world using a computer.
The main components of a virtual world are very similar to the real
world but they are given different names. Instead of real world Places,
People, and Things, there are virtual world Rooms, NPCs, and Objects.
There is an additional and very important component needed for creating
a Virtual World, imagination.

A Place is a Room, but not always a room that would be found in a house.
Rooms represent any place in the world, a bedroom, street, city square,
road, etc. A room does not have a size. It is given dimension only
verbally in the Room Description.

People are NPCs which stands for Non Player Character, in addition to
people, they are used to represent any living creature. Examples of NPCs
are: sentries, guards, sheep, dogs, cats, and spiders.

A Thing is an Object. Some objects can be acquired by players, such as:
a candle, waterskin, sword, pair of pants, or a helmet. Other objects
can be used. A player may drink from a spring, but the spring cannot be
taken. There are also objects which are just for atmosphere, such as: a
statue, table, or a cart.

The last, but certainly the most important component of a successful
virtual world is players. Without players, well, what's the point?

Rooms, NPCs, and Objects all have attributes that are used to make them
unique. These attributes must be coded correctly in order for the Room,
NPC, or Object to have the desired effect in the creation of the Virtual
World. This document describes those attributes and occasionally
suggests appropriate ways to use them. The following is a cross
reference of virtual world terms, not all of which are used in this
document. In general terms these are equivalents, but there can be
distinctions such as a 'Player' usually refers to person at the keyboard
and 'Character' refers to the in-game representation being controlled by
the player.

Room ⬄ Place ⬄ Location  
NPC ⬄ Mobile ⬄ Mob  
Object ⬄ Item ⬄ Thing  
Player ⬄ Character ⬄ Toon  

## Creating a virtual world

A virtual world is experienced by Players. Yes, the ultra-valuable
Player. Nothing of any significance happens in a virtual world without
players. Once a player logs on and begins to play, their interest must
be captured or else they logoff and might never return. An interesting
story is a must for a successful virtual world. So, get a story or make
one up. It may be helpful to assume the role of a director of a play.
What are the scenes, who are the actors, and what are the props? Scenes
might include towns, meeting places, forests, mountains, and places of
worship. Actors in a town might include citizens of the town, guards,
shopkeepers, dogs, cats, and rats. Props might include all armor,
weapons, food, statues, springs, bags, and jewelry. Drawing a map and
noting the location of major scenes, who is involved, and what props
they might need will help ensure a memorable world for players to
experience.

## Critical virtual world components

A virtual world is made of rooms that are connected to each other. A
player explores a virtual world by moving from room to room. Each room
has a name, description, and exits. A room map is essential to creating
a cohesive world.

NPCs are placed in rooms and may move around in the world or remain
stationary. Each NPCs has a description and is given attributes like hit
points and damage. Roaming NPCs will move randomly about the world,
although their movement may be limited by creating rooms which do not
allow NPCs.

Objects are placed in a room or given to an NPC. Players acquire objects
by finding them through exploration, killing NPCs, buying them from NPC
shopkeepers, or getting them from other players.

The initial state of NPCs and objects within the virtual world is
defined using randomizing factors. This ensures that the same NPC and/or
object does not reappear in the same place on predetermined schedule.

## Other virtual world components

There are a several other components that make up a virtual world. First
among these are the socials. Socials are used by players to communicate
in a non-verbal manner with other players. These should be carefully
defined so that they enhance the atmosphere of the virtual world as
opposed to detracting from it. For example, a common command in virtual
worlds is bow. But maybe bow does not fit the atmosphere and should be
replaced with kisscheeks because in the given virtual world people do
not bow, they kiss each other on the cheeks.

Help is also very important and should also fit the atmosphere of the
virtual world. Help usually includes examples and this is where care
must be taken to maintain the atmosphere.

## The library

The virtual world is defined using a set of text files located in the
LIBRARY directory. Each component is located within the LIBRARY
directory and each component has a special format which must be followed
carefully. Remember, it is a computer program that reads these files.
Computer programs tend to be stupid and upon encountering something they
do not understand, they either quit working or spit out some error
message and then try to continue, neither of which is desirable. So, be
very careful to create each file according the specifications in this
document. For the most part, UPPER and lower case have no significance
and anything may be typed in all UPPER case, all lower case, or Mixed
Case. Mixed Case is preferred because it increases readability. The
exception relates to the text displayed to the player. This text will be
displayed 'as is' including spacing, line breaks, and blank lines. A
prime example of this is room descriptions.

## The library files

| File/Directory | Description |
|---|---|
| `\Library\Doors\<RoomId><RoomId>.txt`   | Door files |
| `\Library\Loot\<NPCtype>.txt`           | Loot files |
| `\Library\Mobiles\<MobileId>.txt`       | NPC files |
| `\Library\Objects\<ObjectId>.txt`       | Object files |
| `\Library\Rooms\<RoomId>.txt`           | Room files |
| `\Library\Scripts\Rooms\<RoomId>.txt`   | Scripts - room |
| `\Library\Shops\<RoomId>.txt`           | Shop files |
| `\Library\Sql\<DescriptiveName>.txt`    | Sql statements |
| `\Library\Talk\<NPC type>.txt`          | Talk files |
| `\Library\World\Map\<MapDatabase>.mdb`  | zMapper Access database |
| `\Library\World\Mobiles\<MobileId>.txt` | Specifies NPC spawn information |
| `\Library\World\Objects\<ObjectId.txt>` | Specifies Object spawn information |
| `\Library\Greeting.txt`                 | Player logon greeting |
| `\Library\Help.txt`                     | Help |
| `\Library\Motd.txt`                     | Message of the day |
| `\Library\Social.txt`                   | Socials |
| `\Library\ValidCommands.txt`            | Commands control file |
| `\Library\ValidNames.txt`               | Valid player names |


# Rooms

Rooms are used to define the virtual world. Nothing can exist if there
are no rooms, everything is 'in' some room within the virtual world. The
term 'room' is used regardless of whether the place is a bedroom,
street, river, field, forest, etc. Rooms should be connected in
'geographically' sensible manner. For example, you would not want to
connect an upstairs bedroom directly to the street below. It would make
more sense to have the bedroom connect to a hallway, then some stairs,
then out the door and onto the street.

Rooms are specified in the \Library\Rooms directory as `<RoomId>.txt.
For example, ArmorShop.txt would contain specifications for RoomId
ArmorShop.

## Room specification

A room is specified as follows:  
-----------------------------
```
RoomId:       xxxx  
RoomType:     xxxx  
Terrain:      xxxx  
RoomName:     xxxx  
RoomDesc:     xxxx  
End of RoomDesc  
ExitName:     xxxx  
ExitDesc:     xxxx  
ExitToRoomId: <RoomId>  
End of Exits  
End of Room  
```
-----------------------------

Room Id  
This must be unique within the whole virtual world.

RoomType  
Describes the type of room using one or more of the following codes separated by spaces:
| Code | Meaning |
|---|---|
| None    | Nothing special about this room |
| Dark    | Player must be carrying a light to see |
| Drink   | Player may drink from `<RoomName\ |
| NoFight | No fighting |
| NoNPC   | NPC are not allowed to enter this room from another room |

Terrain  
The type of terrain effects number of movement points used when a player
enters a room. It also effects the rate of food and water consumption.  
Terrain must be **<u>one</u>** of the following:  
Inside  
Street  
Road  
Field  
Forest  
Swamp  
Desert  
Hill  
Mountain

RoomName  
This is the room's name and is displayed to players upon entering a
room.

RoomDesc  
This is the room's description and is displayed to players upon entering
a room. It should be as descriptive as possible, but should not tell the
player how he/she feels. Room descriptions should, in general, be longer
than one line. The room description is free form and is displayed to the
player 'as is' including blank lines and spaces.  

End of RoomDesc  
This line indicates the end of the room description.

ExitName  
This text is the name of the exit displayed to the player. A player will
use this name to exit the room or manipulate the door, if a door is
specified. ExitName must be unique for this room.

ExitDesc  
This text will be displayed to the player when using the command 'look
ExitName'. The text should start on the next line.

ExitToRoomId  
This is the unique id of the room to which the player will be moved when
using the command  
'go ExitName' which can be shortened to `n` to go north.

End of Exits  
This line indicates the end of all exit information. Even if there are
no exits, this line is required

End of Room  
This line indicates the end of all room information

### A note about room exits
Exits allow a player to move from room to room. A room may have from
zero exits to ?? exits, where ?? is a ridiculously high number. Exits
specifications are a repeating, orderly, and logical group. For each
exit, code all the 'exit' specifications. An exit from one room is an
entrance to another room, so exits/entrances must match. For example,
assume two rooms say and Entry and a Dining Room. The Entry would have
an exit to the Dining room and the Dining Room would have an exit to the
Entry.

## Room example 1  
The specifications for an example room \Library\Rooms\Entry.txt follows:  
```
RoomId:   Entry
RoomType: NoFight
Terrain:  Inside
RoomName: Entry
RoomDesc:
A stand with a flower arrangement setting on it adorns
the corner to your left. A chandelier hangs from the two
story ceiling casting away any shadows. A long flight of
stairs leads upstairs to the bedrooms.
End of RoomDesc
ExitName: DiningRoom
ExitDesc: You see the dining room with table and chairs.
ExitToRoomId: DiningRoom
End of Exits
End of Room
```

## Room example 2
The specifications for an example room \Library\Rooms\DiningRoom.txt follows:
```
RoomId: DiningRoom
RoomType: NoFight NoNPC
Terrain: Inside
RoomName: Dining Room
Desc:
A large table sits under a candle chandelier. Three pictures
are hung on one wall. On another wall is a magnolia wreath.
End of RoomDesc
ExitName: Entry
ExitDesc: You see the Entry.
ExitToRoomId: Entry
End of Exits
End of Room
```

These two examples define two rooms both have terrain type of 'inside'.
RoomId Entry is named Entry and RoomId DinningRoom is named Dining Room.
Fighting is not allowed in either room and no NPCs are not allowed in
the Dinning Room.

# Doors

Doors are defined in terms of the relationship between Rooms and Exits.
So, in Room Example 2 above, if you wanted a door between the DiningRoom
and the Entry, the DoorId would be DiningRoomEntry. Now, there is also a
door when going from the Entry to the DiningRoom which is in fact the
other side of the same door, but with a DoorId of EntryDiningRoom. A
door specification is actually a pair of specifications, one for each
side of the door. Opening/closing one side of the door automatically
cause the other side to be opened/closed, also unlocking/locking one
side of the door automatically cause the other side to be
unlocked/locked. It is possible to have a door that can be
unlocked/locked from only one side by having one side have a key and
other side not have a key.

Doors are specified in the \Library\Doors directory as
`<RoomId><RoomId>.txt`. For example, a door between the DiningRoom and
the Entry is defined by two entries, (1) DiningRoomEntry.txt would
contain one side door specifications and (2) EntryDiningRoom.txt would
contain the specifications for the other side of the door.

## Door specification

A door is defined as a pair of these specifications:

```
DoorId: `<RoomId><RoomId>
DoorName: `<name>
InitialState: Open | Closed
LockInitialState: NoLock | Locked | Unlocked
KeyObjectId: None | `<ObjectId>
ResetInterval: `<time>
```

DoorId  
The first specification is the concatenation of RoomId and an
ExitToRoomId. The second specification is a concatenation of the same
ExitToRoomId and the same RoomId.

DoorName  
The name that the player must use to refer to this particular door when
opening, closing, locking, or unlocking the door. Typically this will be
simply 'door', but it could be 'gate' or anything else.

InitialState  
This controls whether the door is open or closed when the games starts
and is specified as 'Open' or 'Closed'. This is the state the door will
normally be in when encountered by a player.

ResetInterval  
This controls how long before the door is returned to the InitialState
after a player opens/closes the door. It also controls how long before
the lock is returned to the LockInitialState after being locked/unlocked
by a player. ResetInterval `<time>` is specified as a series of 7
integer numbers, each separated from the next by a space and in the
following sequence:  
seconds  
minutes  
hours  
days  
weeks  
months  
years

LockInitialState  
Specifies the state of the lock when the game starts.
NoLock The door has no lock and cannot be locked or unlocked.
Locked The door has a lock which will normally be locked
Unlocked The door has a lock which will normally be unlocked

KeyObjectId  
The ObjectId of the key which fits the lock. The player must have the
key in their inventory when they issue the lock/unlock command. If
LockInitialState is Locked or Unlocked, then a valid ObjectId must be
specified. The object must, of course, be defined as type 'key'. If
LockInitialState is NoLock, then 'None' must be specified for the
KeyObjectId.

## Door side 1 example

The specifications for an example door \Library\Doors\DiningRoomEntry.txt follows:
```
DoorId: DiningRoomEntry
DoorName: Door
InitialState: Closed
ResetInterval: 0 2 0 0 0 0 0
LockInitialState: Locked
KeyObjectId: HouseKey
```

## Door side 2 example

The specifications for an example door \Library\Doors\EntryDiningRoom.txt follows:
```
DoorId: EntryDiningRoom
DoorName: Door
InitialState: Closed
ResetInterval: 0 2 0 0 0 0 0
LockInitialState: Locked
KeyObjectId: HouseKey
```
The two examples define a door between the DiningRoom and the Entry that
is normally closed and locked, will return to closed and locked state
after 2 minutes of being opened, and requires the player to have an
object in their inventory with an ObjectId of HouseKey to unlock or lock
the door.

Some of the values between the two side of a door may be different, such
as DoorName, LockInitialState, and KeyObjectid. This is not recommended,
but could be used in special cases such as creating a door that can be
only be unlocked/locked from one side by setting KeyObjectId to 'None'
on one side of the door.

## Valid door states and transitions

### Door states

Open NoLock  
Open Unlocked  
Closed NoLock  
Closed Unlocked  
Closed Locked

### Door state transitions

OpenNoLock 🡪 Close 🡪 ClosedNoLock  
OpenUnlocked 🡪 Close 🡪 CloseUnlocked  
ClosedNoLock 🡪 Open 🡪 OpenNoLock  
ClosedUnlocked 🡪 Open 🡪 OpenUnlocked  
ClosedUnlocked 🡪 Lock 🡪 ClosedLocked  
ClosedLocked 🡪 Unlock 🡪 ClosedUnlocked  

### Current door status

When the game is running. the current status of the doors is tracked in
the \Running\Doors directory. The primary door tracks the open/closed
status of both the primary and secondary doors as well as the status of
the lock on that side of the door. The secondary door tracks only the
status of the lock on that side of the door. Possible file content
combinations are show below:

DiningRoomEntry.txt  
Open|Closed NoLock|Locked|Unlocked

EntryDinigRoom.txt  
Open|Closed NoLock|Locked|Unlocked

# NPCs

NPCs represent everything that is living with the virtual world, other
than players. When building your virtual world, it might be helpful to
categorize NPCs. For example: enemies, shopkeepers, guards, citizens,
hostile animals, and friendly animals. The server does not distinguish
between these categories, but as you create the NPCs you would give each
category similar attributes. Shopkeepers and guards should have high
hitpoints and hit very hard. Killing citizens and friendly animals
should cause the player's faction for that NPC to go down very quickly.
Be sure to keep the NPCs within the realms of believability for the
virtual world you are creating. If you are creating a historical virtual
world set in the time of the Roman Empire, you would not have dragons,
skeletons, goblins, and the like. It would be appropriate to have Huns,
Goths, Greeks, Egyptians, and the like.

NPCs are specified in the \Library\Mobiles directory as
`<MobileId>.txt`. For example, SmallRat.txt would contain specifications
for Mobile SmallRat.

## NPC specification

An NPC is specified as follows:

```
MobileId: xxxx
Names: xxxx xxxx xxxx
Sex: xxxx
Desc1: xxxx
Desc2: xxxx
Desc3: xxxx
End Desc3
Action: xxxx
Faction: xxx
Level: 99
HitPoints: 99999
Armor: 999
Attack: xxxx
Damage: 99999
ExpPoints: 99999
Loot: <NPCtype>
Talk: <NPCtype>
```

MobileId  
This must be unique within the whole virtual world. Players can use this
id to identify the NPC.

Names  
A list of names, separated by spaces, that a player can use to identify
the NPC.

Sex  
The sex of the NPC. It must be **<u>one</u>** of the following codes:
| Code | Meaning |
|---|---|
| F | Female |
| M | Male |
| N | Neutral |

Desc1 
The description of the NPC used when the NPC takes some action. For
example, a short description of "the Mangy Mutt" would result in
messages such as "The Mangy Mutt leaves south." and "The Mangy Mutt
bites you for 5 points of damage.". Desc1 should never end with a
punctuation mark because it will be used as part of sentences such as
those above. If Desc1 starts a sentence, the first letter of the first
word will be capitalized,.

Desc2  
The description displayed when the NPC is just hanging out. For example,
"The Mangy Mutt is here, running around and barking." Desc2 should be a
complete sentence including the appropriate punctuation and
capitalization.

Desc3  
The description displayed when a player looks at the NPC by typing `look <mobile>`.
Desc3 can be multi-line and will be displayed to the player
'as is' including blank lines and spaces.

End of Mobile Desc3  
This line indicates the end of the Desc3 specification.

Action  
Action is used modify the behavior of the NPC. Upper and lowercase must
be specified exactly as shown. If no action code is required, specify
only 'Action:', thereby indicating a blank action code. If multiple
codes are selected, separate the codes by one space. The following codes
may be selected:

| Code    | Meaning |
| ---     | --- |
| None    | just an ordinary mobile |
| Aggro   | will attack a player without warning and without considering faction|
| Faction | will determine whether or not to attack player based on factions |
| Destroy | will destroy any object laying around that can be taken |
| Help    | will help other NPCs in the same room that are being attacked| 
| NoMove  | will not move |
| Wimpy   | will flee when about to die |

Faction  
Use one of the following to describe the NPC's faction. Faction is used
to determine whether or not an NPC that is sensitive to factions (Action
set to Faction) will attack a player without provocation. It is
required, but will be ignored if the 'Action' codes do not include
'Faction'. It must be **<u>one</u>** of the following:  
Good  
Lawful  
Neutral  
Lawless  
Evil  

Level  
This number is the level of the NPC. Level is used when creating an NPC
to establish HitPoints, Armor, Damage, and ExpPoints, each of which can
be modified.

HitPoints  
This number is added to the NPC’s base HitPoints. This allows the NPCs
to be created that are more difficult than the standard NPC for a given
level. Negative numbers are also allowed, thus making an NPC easier than
normal. This number should be zero for most NPCs.

Armor  
This number is added to the NPC’s base Armor and determines the ability
of the NPC to avoid damage. Not implemented.

Attack  
This specifies the type of attack the NPC will use during combat to
produce the appropriate message. It is converted to lower case because
it is used in the middle of a sentence. For example, The Mangy Mutt
*bites* you for 5 points of damage. Attack must be **<u>one</u>** of the following:  
bites  
claws  
crushes  
hits  
mauls  
pierces  
punches  
slashes  
stabs  
stings  
thrashes  

Damage  
This number is added to the NPC’s base Damage. Damage is the amount of
damage the NPC can do per round. This is not the exact amount of damage,
a random number of damage points will be added or subtracted from this
number for each round of combat to obtain the actual amount of damage.
Negative numbers are also allowed, thus making an NPC easier than
normal. This number should be zero for most NPCs.

ExpPoints  
This number is added to the NPC’s base ExpPoints. ExpPoints is the
number of experience points used to award experience to players when the
NPC is killed. Negative numbers are also allowed, thus making an NPC
give less experience than normal. This number should be zero for most
NPCs.

Loot  
Specify `<NPCtype>` like dog, bee, snake, etc. This MUST match exactly
an existing NPCtype located in \Library\Loot..

Talk  
Specify `<NPCtype>` like dog, bee, snake, etc. This MUST match exactly
an existing NPCtype located in \Library\Talk..

## NPC example 1

The specifications for an example NPC \Library\Mobiles\MangyMutt.txt follows:
```
MobileId:  MangyMutt
Names:     Mangy Mutt
Sex:       Neutral
Desc1:     The Mangy Mutt
Desc2:     The Mangy Mutt is here, running around and barking.
Desc3:
The Mutt is a small dog that smells awful, his fur is
matted, and he is very skinny.
End Desc3
Action:    Destroy
Faction:   Neutral
Level:     1
HitPoints: 5
Armor:     0
Attack:    Bite
Damage:    1
ExpPoints: 5
Loot:      Dog
Talk:      Dog
```
This NPC is a Mangy Mutt that will pick any objects that it finds laying
around. It is a level one NPC with 5 extra hit points, no extra armor
points, 1 extra damage point when it bites a player. When a player kills
it, the player receives the base experience points for a level 1 NPC,
plus 5 extra experience points. A loot and talk file named 'Dog.txt'
already exist in the \Library\Loot and \Library\Talk directories.

## NPC example 2

The specifications for an example NPC \Library\Mobiles\SandScorpion.txt follows:  
```
MobileId:  SandScorpion
Names:     sand scorpion
Sex:       Neutral
Desc1:     a sand scorpion
Desc2:     A sand scorpion raises its tail and stands its ground.
Desc3: The sand scorpion is dangerous and won't hesitate to sting you!
Action:    Aggro
Faction:   Neutral
Level:     3
HitPoints: 10
Armor:     2
Attack:    Sting
Damage:    5
ExpPoints: 25
Loot:      Scorpion
Talk:      None
```
This NPC is a sand scorpion that will attack without provocation. It is
a level three NPC with 10 extra hit points and it does an extra 5 points
damage when it stings a player. When a player kills it, the player
receives and extra 25 experience points. A loot and talk file named
'Scorpion.txt' already exist in the \Library\Loot directory and
\Library\Talk must contain a file name 'None.txt' which would be empty.

Objects are specified in the \Library\Objects directory as
`<ObjectId>.txt`. For example, RedApple.txt would contain specifications
for Object RedApple.

# Objects

## Object specification

An object is specified as follows:
```
ObjectId: xxxxx
Names: xxxx xxxx xxxx
Desc1: xxxx
Desc2: xxxx
Desc3: xxxx
End Desc3
Weight: 99
Cost: 99
Type: xxxx
<Object type specifications>
```

ObjectId  
This must be unique within the whole virtual world.

Names  
A list of names, separated by spaces, that a player uses to identify the
object. At least one of these names must be unique within the whole
virtual world, for example the ObjectId.

Desc1  
The description displayed when the object is used. For example, a short
description of "a long, green stick" would result in messages such as
"The Mangy Mutt picks up the long, green stick.". Desc1 should be kept
short and should never end with a punctuation mark because it will be
inserted into the middle of sentences.

Desc2  
The description displayed when the object is just there. For example, "A
folded robe is lying here". Desc2 should be kept short and should end
with appropriate punctuation.

Desc3  
The description displayed for the object when a player examines the
object by typing 'examine `<object>`. Desc3 should be a fairly
descriptive and should typically be more than one line

End Desc3  
This line indicates the end of Desc3.

Weight  
The weight of the object. The amount of weight players can carry is
limited. The stronger they are, the more they can carry. Weight can be
zero for very small objects. This is a whole number that is zero or
greater.

Cost  
The value used to calculate the price when a player is buying or selling
the object in a shop. The value is specified as silver pieces.

Type  
Specifies the object type and must be **<u>one</u>** of the following
codes:

| Code      | Meaning |
| ---       | --- |
| Armor     | protection during fights, clothing, and accessories |
| Container | a place to store things |
| Drink     | satisfies thirst |
| Food      | satisfies hunger |
| Junk      | object of very little value |
| Key       | yep, it’s a key alright |
| Light     | useful in dark places |
| NoTake    | may not be taken |
| Treasure  | gems, objects of great value, unique objects |
| Weapon    | useful in combat |

## Object type specifications

All object types are listed below with explanations. Some object types
require additional specifications that are unique to that object type.
For example, an object type of 'Weapon' requires additional
specifications of 'WeaponType' and 'WeaponDamage'.

Armor  
ArmorValue:999  
This value is added to the player's armor class when the object is worn
and is subtracted from the player's armor class when the object is
removed. As a player increases their armor class, the amount of damage
taken during fights is reduced. See the World Building section for more
about armor class.

ArmorWear:xxxx  
Where xxxx is **<u>one</u>** of the following:  
Head  
Ear  
Neck  
Shoulders  
Chest  
Back  
Arms  
Wrist  
Hands  
Finger  
Shield  
Waist  
Legs  
Ankle  
Feet  

Note: The player may also specify 'left' or 'right' when 'wearing' one of these objects.  
Ear  
Wrist  
Finger  
Ankle  

Container  
Capacity:999  
The sum of the weight of all objects in the container cannot exceed the capacity.  

Drink  
DrinkPct:999  
The player's thirst will be decreased by this percentage when they drink
this item. Drink items are completely used up in one use. Pct ranges
from -100 to +100. Salt water is an example of an item that would have a
negative percent thereby increasing thirst.

Food  
FoodPct:999  
The player's hunger will be decreased by this percentage when they eat
this item. Food items are completely used up in one use. Pct ranges from
-100 to +100. Moldy bread is an example of an item that would have a
negative percent thereby increasing hunger.

Junk  
This is just what is says, Junk! Not worth very much.

Key  
A key’s object id should match at least one door's KeyObjId
specification, unless of course you want confuse players. No additional
specification is required.

Light  
Hours:99
Number of hours this light will burn.  

NoTake  
This object can be used to enhance the virtual world. Example are:
statues, tables, carvings.

Treasure  
Players will enjoy find these objects of great value. These objects
should be worth a significant amount of money and should not be easy to
obtain. No additional specification is required.

Weapon  
WeaponType:xxx  
Specify the type of weapon as **<u>one</u>** of the following:
| Type   | Message |
| ---    | --- |
| Axe    | Chop
| Club   | Bash
| Dagger | Stab
| Hammer | Pound
| Spear  | Pierce
| Staff  | Whack
| Sword  | Slash

WeaponDamage:999  
The amount of damage the weapon can do per round. This is not the exact
amount of damage, a random number will be subtracted from this number
for each round of combat to obtain the actual amount of damage.

## Object examples

The specifications for two example objects follows:  
Object example 1 \Library\Objects\ShortSword.txt  
```
ObjectId : ShortSword
Names:     Short Sword ShortSword
Desc1:     a short sword
Desc2:     A short sword with a short, shiny, and sharp blade is here.
Desc3:
The short sword has a short, shiny, sharp blade.
End Desc3
Weight:       3
Cost:         5
Type:         Weapon
WeaponType:   Sword
WeaponDamage: 10
```
This object's id is 'shortsword' and it can be referenced by the names
'shortsword', 'short', or 'sword'. It weighs 3 units and its cost is 5
units. It is a sword type weapon that does maximum damage of 10 points.

Object example 2 \Library\Objects\TrainingSkullCap.txt  
```
ObjectId:    TrainingSkullCap
ObjectNames: Training Skull Cap TrainingSkullCap
ObjectDesc1: a training skull cap
ObjectDesc2: A training skullcap is lying on the ground.
ObjectDesc3:
This is a training skullcap. It doesn't offer much protection, but
some protection is better than none.
End Desc3
Weight:     1
Cost:       2
ObjectType: Armor
ArmorValue: 1
ArmorWear:  Head
```
This object's id is 'trainingskullcap' and it can be referenced by the
names 'trainingskullcap', 'training', 'skull', 'cap' and 'skullcap'. It
weighs 1 units and its cost is 2 units. It is a piece of armor with a
value of 1 that can be worn on the head.

# Shops

Shops are where players can buy and sell items. Shops should be
scattered through out the virtual world, but of course should be found
mainly in cities and towns with an occasional outpost. Only one shop can
exist in a given room and each shop should have a shopkeeper NPC.
Typically, each shop specializes in a particular type of item. The shops
you create might include a weapon shop, armor shop, jewelry shop, or a
bazaar that buys and sells a mixture of items.

Shops are specified in the \Library\Shops directory as `<RoomId>.txt`.
For example, ArmorShop.txt would contain specifications for Shop
ArmorShop.

## Shop specification

Each shop is specified in a separate file whose name MUST match exactly
the name of an existing RoomId in the \Library\Rooms directory. The shop
codes are specified as follows:
```
<shop welcome message>
Item: <ObjectId>
.  
.  
.  
Item: <ObjectId>
End of Items  
End of Shop  
```

`<shop welcome message>`  
The message given to a player when the player issues the 'list' command.  

Item  
Item that this shop will buy and sell. ObjectId MUST match exactly an
ObjectId that exists in the \Library\Objects directory. This
specification is repeated for each item bought and sold in this shop.  

End of Items  
This line indicates the end of the list of items.

End of Shop  
This line indicates the end of the shop specifications.

## Shop example

The specifications for an example shop \Library\Shops\WeaponShop follows:  
```
Welcome to my weapon shop
Item: RustySword
Item: BluntAxe
Item: KindlingAxe
Item: SmallClub
Item: TinyDagger
Item: SimpleHammer
Item: CrackedSpear
Item: LongStick
End of Items
End of Shop
```
The message "Welcome to my weapon shop" will be displayed when the
player issues the 'list' command. This shop sells a variety of weapons.

# Loot

Loot refers to items given to player's upon the defeat of an enemy. Loot
items should 'make sense' for a given NPC. For example, loot for a rat
could be a rat ear or two, a rat tail, and several whiskers. Rat loot
would not include items like a shepherd's staff, a loaf of bread, or a
helmet.

Loot is specified in the \Library\Loot directory as `<NPCtype>.txt`. For
example, Rat.txt would contain specifications for rat loot.

## Loot specification

Each NPCtype is specified in a separate file whose name MUST match
exactly the 'loot' specification in an NPC specification. The loot is
specified as follows:  

`<count>` `<percent>` `<ObjectId>`  

`<count>`  
The NPC has this many of the specified object . When more than 1 is
specified, the object is considered `<count>` times for looting.

`<percent>`  
The chance the player has of getting this object.  

`<ObjectId>`  
The ObjectId of the object. This MUST match exactly an ObjectId
specified in the \Library\Objects directory.

## Loot example

The specifications for an example loot file \Library\Loot\Bee.txt follows:  
```
1 50 BeeStinger
2 65 BeeWing
```
The first line specifies that a bee will have 1 BeeStringer and the
player will have a 50% chance of getting that object for each bee
killed. The second specifies that a bee will have 2 wings and the player
will have a 65% chance of getting that object for each bee killed.

# Talk

Talk enables NPCs to respond when a player 'hails' them by issuing the
'hail' command. Talk specifications are typically grouped by NPCtype.
For example, thief for all thieves, guard for all guards, merchant for
all merchants.

Talk is specified in the \Library\Talk directory as `<NPCtype>.txt. For
example, Thief.txt would contain specifications for any thief speaking..

## Talk specification

Each NPCtype is specified in a separate file whose name MUST match
exactly the 'talk' specification in an NPC specification. Talk is
specified as follows:
```
Maximum Message Number <count>
Message <sequence number>
<message text>
End of Message
End of Messages
```
Maximum Message Number `<count>`  
`<count>` is simply the number of messages. Each time a player hails an
NPC, a random number is generated between 1 and `<count>` to determine
which message will be used to respond to the player.  

Message `<sequence number>`  
The first message must be sequence number 1, the second message must be
sequence number two, etc.  

`<message text>`  
This message, when selected, will be displayed 'as is' to the player
including spaces and blank lines.

End of Message  
This line indicates the end of a message.

End of Messages  
This line indicates the end of all messages.

## Talk example

The specifications for an example talk file \Library\Loot\Thief.txt follows:  
```
Maximum Message Number 2
Message 1
I will relieve you of the burden of living!
End of Message
Message 2
I will steal your heart!
End of Message
End of Messages
```

# World

The world files are used to 'spawn' NPCs. Spawning is the process by
which NPCs enter the world and it controls how many bees, snakes,
thieves, etc there are in the world at any given time and how often they
return after being killed. Objects could also be spawned, but currently
that functionality is not implemented.

## NPC spawn specifications

NPC spawning is specified in the \Library\World\Mobiles directory as
`<MobileId>.txt`. For example, Viper.txt would contain specifications
for spawning vipers.

Each NPC to be spawned is specified in a separate file whose name MUST
match exactly an existing NPC in the \Library\Mobiles directory. NPC
spawning is specified as follows:
```
MaxInWorld: <count>
RoomId: <RoomId>
Interval: <time>
```
MaxInWorld  
This controls the maximum number of the specified NPC that can exist in
the world at the same time.

RoomId  
The RoomId into which the NPC will spawn. This RoomId MUST exist in the
\Library\Rooms directory.

Interval  
This determines the amount of time that will elapse between the time an
NPC is killed and when they spawn again. Interval `<time>` is specified
as a series of 7 integer numbers, each separated from the next by a
space and in the following sequence:  
seconds
minutes
hours
days
weeks
months
years

## NPC spawn example

The specifications for an example NPC spawn file \Library\World\Mobiles\SmallRat.txt follows:  
```
MaxInWorld: 30
RoomId: RatSpawnRoom220
Interval: 0 2 0 0 0 0 0
```
This specification will cause 30 small rats to be spawned into RoomId
RatSpawnRoom220 two minutes after the game is started. Additionally,
when OMugs determines that the number of small rats in the world has
fallen below 30, the appropriate number of small rats are scheduled for
spawning 2 minutes later in order to bring the small rat count back up
to 30,

## Object spawn specifications

The object spawning functionality is not implemented.

# Socials

Socials provide a means for players to communicate with each other
without 'saying' something. Socials provide another opportunity for the
virtual world builder to be creative and have some fun. Some common
socials are smile, wave, point, and laugh.

It is important to maintain the atmosphere of the game when creating
socials. For example, a virtual world set in the modern day United
States probably would not include a 'curtsey' social, but a 'curtsey'
would be fit nicely in a virtual world set in medieval times.

Socials are specified in the \Library\Socials.txt file.

## Social specification

Socials have two forms, the two line form and the seven line form. The
two line form is for socials that do not have a target. The seven line
form is for socials that may have a target.

A two line social is specified as follows:
```
Social : xxxx
MinPos : xxxx
Lines : 2
Message to player
Message to all players in the room
```

A seven line social is specified as follows:
```
Social : xxxx
MinPos : xxxx
Lines : 7
Message to player when no target is given
Message to all players in the room when no target is given
Message to player when target is player
Message to player when target is not in the room
Message to player when target is given
Message to target when target is given
Message to all players in the room when target is given
```

Social  
The unique name of the social. This is used by players just like any
command and must be only **<u>one</u>** word.

MinPos  
This is the minimum position the player must be in to execute the
social. The order of the positions are from lowest to highest: sleep,
sit, stand. This must be **<u>one</u>** of those three values.

Lines  
This is the number of lines in the social and it must be either a '2' or
a '7'. A two line social cannot have a target, but a 7 line social may
or may not have a target.

Messages  
The messages may contain any combination of the following substitution
codes:

| Code | Is replaced with |
| ---  | --- |
| $P               | player's name |
| $T               | target player's name |
| $pHeShe          | he or she, based on the player's sex |
| $pHimHer         | him or her, based on the player's sex |
| $pHisHers        | his or hers, based on the player's sex |
| $pHimselfHerself | himself or herself, based on the player's sex |
| $tHeShe          | he or she, based on the target player's sex |
| $tHimHer         | him or her, based on the target player's sex |
| $tHisHers        | his or her, based on the target player's sex |
| $tHimselfHerself | himself or herself, based on the target player's sex |

Great care should be taken in creating these messages. It can be
difficult to get all these message so that they make sense to the
recipient of the message. Write the messages out first without using the
substitution codes can help avoid grammatical errors. Pretend to be the
player, the target player, and the other players in the room as see if
the message you would receive in each case make sense.

## Social examples

The specifications for two example socials, the first is a two line
social and the second is a seven line social:

Social example 1, two line social:
```
Social : grin
MinPos : sit
Lines : 2
A broad grin spreads across your face.
$P grins broadly.
```
Player must be sitting or standing to grin (cannot be asleep). The grin
social does not accept a target. Assume Zeke grins, and that Zak and Zeb
are bystanders in the room. Here are the messages each player would
receive in response to Zeke's grin social:
```
Zeke types
grin
Zeke sees
A broad grin spreads across your face.
Zak and Zeb see
Zeke grins broadly.
```

Social example 2, seven line social:
```
Social : smile
MinPos : sit
Lines : 7
You smile.
$P smiles.
Smiling at yourself??
The recipient of your pleasantry is nowhere in sight.
You smile at $T.
$P smiles at you.
$P smiles at $T.
```
Player must be sitting or standing to smile (cannot be sleeping). The
smile social may or may not have target. Assume Zeke smiles at Zeesey,
and that Zak and Zeb are bystanders in the room. Here are the messages
each player would receive in response to Zeke's smile social:
```
Zeke types
Smile Zeesey
Zeke sees
You smile at Zessey.
Zeesey sees
Zeke smiles at you.
Zak and Zeb see
Zeke smiles at Zeesey.
```

Social example 3, a seven line social:
```
Social : dance
MinPos : stand
Lines : 7
You dance for joy!
$P dances for joy!
Dancing with yourself, very peculiar.
Your dance partner doesn't seem to be here.
You clasp $T and begin to dance with $tHimHer.
$P clasps you. You and $pHimHer begin to dance.
$P clasps $T and begins to dance with $tHimHer.
```
You must standing to dance. The dance social may or may not have a
target. Assume Zeke a male player and he wants to dance with a female
named Zeesey, and that Zak and Zeb are bystanders in the room. Here are
the messages each player would receive in response to Zeke's dance
social.
```
Zeke types
Dance Zeesey
Zeke sees
You clasp Zessey and begin to dance with her.
Zeesey sees
Zeke clasps you. You and him begin to dance.
Zak and Zeb see
Zeke clasps Zeesey and begins to dance with her.
```

# Help

Help is where players will look to learn about the game, the more
completely you cover the game in help topics, the easier it will be for
players to enjoy the world you have created. At a minimum, there should
be a help topic for every player command. There should also be help on
general questions, for example a 'travel' topic would describe how a
player travels and all the factors that might effect their travels and
include references to specific commands like 'go'.

It is important to maintain the atmosphere of the game within the help
topics. For example, pay close attention to way in which each command is
presented and the content of the command examples.

## Help specification

Help is specified as follows:

```
Help:xxxx
Xxxx xxxx xxxx
Xxxx xxxx xxxx
Related help:xxxx
```

Help  
This is the unique help keyword. Leaving this blank for **<u>one</u>**
help entry will cause that entry to be displayed if the player types
only 'help'. Everything from this line, down to and including 'Related
help', is displayed to the player 'as is', so care should be taken to
format the help text in an easy to read manner.

Related help  
This is where related help topics are listed and is required.

## Help example

The specifications for example help entry 1 follows:
```
Help:Go
Go is used to move from where you are now
into an adjacent room.
Usage:Go <exit>
Example: Go N
Related help: Exits
```
This is the help entry for the 'Go' command. It begins with 'Help:Go'
and ends with 'Related help: Exists'

The specifications for example help entry 2 follows:
```
Help:
Type 'help' followed by the command you wish
to get help with, like 'help move' or 'help newbie'.
Related help: Newbie
```
This is the help entry will be displayed when a player enters 'help' by
itself.

# Player commands

Player commands and parameters are **<u>not</u>** case sensitive. Each
command has a specific usage that includes the follow components:
```
command - The command
CAPITAL - Indicates parameters that, when used, must be spelled exactly as indicated
{ }     - Optional parameters are enclosed within braces
|       - Separates a list of possible parameter values
< >     - Specifies a value to be supplied by the player
```
Some examples: 
The Inventory command's usage is "inventory", which means this command
has no parameters.

The Color command's usage is "color {ON | OFF}", which means the player
may enter just 'color' or 'color on' or color off'.

The Give command's usage is "give `<object>` `<target>`", which means
the player must enter "give" followed by the name of an object and the
name of the player to whom the player wishes to give the object.

The Group command's usage is " Usage: group LIST | ON | OFF | NONE |
`<player>`", which means the player must enter "group" followed by
"LIST" or "ON" or "OFF" or "NONE" or the name of a player.

## Afk

Usage : afk

Toggles the afk, away from keyboard, flag on and off.

## Assist

Usage : assist `<player>`

Allows one player to assist another player in killing a mobile.

## Buy

Usage : buy `<object>`

May only be used in a shop. Allows a player to buy items from a shop.

## Chat

Usage : chat `<message>`

Whole game communication channel for out-of-character chatter.

## Close

Usage : close `<name>`

Allows a player to close a door, gate, grate, etc. Name must match the
DoorName specified in the primary door specificaton.

## Color

Usage : color {ON | OFF}

Controls whether or not text is displayed with or without color.
Entering 'color' by itself will inform the player whether color is on or
off.

## Consider

Usage : consider `<NPC>`

Allows a player to get a hint about their chances of winning a fight the
NPC.

## Delete

Usage : delete `<player>` `<password>` MAKE IT SO

Allows a player to delete themselves. The player is deleted and
disconnected immediately.

## Drink

Usage : drink `<object>` 

Player drinks the object specified only if the object is a drink object.
Each drink object drunk reduces the player's thirst. Thirst is
represented on a scale of 0 to 100, 0 being no thirst and 100 being
maximum thirst. A player that has a thirst of 0 will gain hit points
faster than a player with a thirst of 100.

## Destroy

Usage : destroy `<object>` 

Allows a player to 'delete' an object from their inventory.

## Drop

Usage : drop `<object>` 

Allows a player to remove an object from their inventory and drop it on
the ground.

## Eat

Usage : eat `<object>` 

Player eats the object specified only if the object is a food object.
Each food object eaten reduces the player's hunger. Hunger is
represented on a scale of 0 to 100, 0 being no hunger and 100 being
maximum hunger. A player that has a hunger of 0 will gain hit points
faster than a player with a hunger of 100.

## Emote

Usage : emote `<message>` 

Message is sent to all players in the same room as the originating
player. The `<message>` is prefixed by the originating player's name. So
if Abbi types 'em thinks Ixaka is silly!' everyone in the room would see
'Abbi thinks Ixaka is silly!'. Players can use emotes in a to express
themselves in a variety of ways. Here's some examples (assume a male
player named Ixaka and a female player named Abbi):

(Ixaka) em rubs his chin, lost in thought

Ixaka rubs his chin, lost in thought.

(Ixaka) em is confused.

Ixaka is confused.

(Ixaka) em approaches Abbi for a closer look.

Ixaka approaches Abbi for a closer look.

(Abbi) em is unsure of Ixaka's intentions and moves away.

Abbi is unsure of Ixaka's intentions and moves away.

## Equipment

Usage : equipment

Alternate usage : eq

Lists objects that the player is wearing and wielding.

## Examine

Usage : examine `<object>` 

Allow a player to get more information about an object.

## Flee

Usage : flee `<direction>` 

Allows a player to flee from a fight in the `<direction>` indicated. The
`<direction>` must be a standard room exit such as North, South, East,
or West.

## Follow

Usage : follow LIST | NONE | `<player>` 

Allows players to automatically follow each other. Before players can
follow, they must be grouped with the player they are going to follow.

LIST will list the player being followed and any players who are
following

NONE will cause the player to stop following

`<player>` causes the player to begin following `<player>`

## Get

Usage : get `<object>`

Allows a player to pick up items.

## Give

Usage : give `<object>` `<target>`

Allows a player to give an object to another player.

## Go

Usage : go `<exit>`

Alternate usage : {go} `<direction>`

Player leaves the current room via the exit specified. The alternate
usage allows a player to type only the direction, such as: north, south,
east, west, northeast, southeast, northwest, southwest. These can be
abbreviated as n s e w ne se nw sw.

## Group

Usage : group LIST | ON | OFF | NONE | `<player>`

Provides the ability for players to form groups.

LIST Lists the group members

ON Player is accepting requests to group

OFF Player is rejecting requests to group

NONE Player leaves the group immediately

`<player>` Allows a player to form a group and/or add new players to the
group

## Gsay

Usage : gsay `<message>`

Send a message to all member of the group.

## Hail

Usage : hail `<npc>`

Hailing certain NPCs will cause the NPC to reply to the player.

## Help

Usage : help {word}

Help by itself will return a general help screen. Help followed by
**<u>one</u>** word will display the help topic matching the word.

## Inventory

Usage : inventory

Alternate usage : i

Displays a list of items that the player is currently carrying.

## Kill

Usage : kill `<mobile>`

Start a fight between the player and the mobile.

## List

Usage : list

May only be used in a shop. Returns a list of items that a player may
buy and sell while in the shop.

## Lock

Usage : lock `<name>`

Allows a player to lock a door, gate, grate, etc. Name must match the
DoorName specified in the primary door specificaton.

## Look

Usage : look {`<exit>` | `<mobile>`}

Alternate usage : l {`<exit>` | `<mobile>`}

Look by itself, returns the room name, room description, and exits. Look
{exit} will return a description of what is to be seen when looking at
an exit.

## Money

Usage : money

Displays how much money a player posses.

## Motd

Usage : motd

Displays the Message of the Day.

## Open

Usage : open `<name>`

Allows a player to open a door, gate, grate, etc. Name must match the
DoorName specified in the primary door specificaton.

## Password

Usage : password `<password>` `<newpassword>` `<newpassword>`

Allows a player to change their password. Where `<password>` is their
current password and `<newpassword>` must be entered twice.

## Played

Usage : played

Displays player's age, time played, and birthday.

## Quit

Usage : quit

Player leaves the game immediately.

## Remove

Usage : remove `<object>`

Allows a player to remove an object they are wearing or wielding and
place it in their inventory.

## Save

Usage : save

Allows a player to save their player information to disk.

## Say

Usage : say `<message>`

Message is sent to all players in the same room as the originating
player. This command is the primary means of player to player
communication, it allows player to talk to each other.

## Sell

Usage : sell `<object>`

May only be used in a shop. Allows a player to sell items to a shop.

## Show

Usage : show COMMANDS | HELP | SOCIALS

Shows the player a list of valid commands, help topics, or socials.

## Sit

Usage : sit

Player sits down.

## Sleep

Usage : sleep

Player goes to sleep.

## Stand

Usage : stand

Player stands up.

## Status

Usage : status

Displays information and settings specific to the player.

## Tell

Usage : tell `<player>` `<message>`

Sends a private message to a player.

## Time

Usage : time

Display the current date and time according to the computer on which the
game server is running and the current in-game date and time.

## Title

Usage : title {`<text>` | NONE}

Sets player's title to the `<text>`. If `<text>` is not specified, the
player's current title is displayed. Specifying 'none' will remove the
player's title. Note: text must be less than 26 characters.

## Train

Usage : train {`<weapon type>`}

Increases the player's skill by one point for the specified type of
weapon. Player's are awarded skill points each time they achieve a new
level. If `<weapon type>` is not entered, the player's unused skill
points and a summary of their skills will be displayed.

## Unlock

Usage : unlock `<name>`

Allows a player to unlock a door, gate, grate, etc. Name must match the
DoorName specified in the primary door specificaton.

## Wear

Usage : wear `<object>`

Allows the player to wear various pieces of clothing, jewelry, and armor
which are currently in their inventory.

## Wake

Usage : wake

Player awakens.

## Who

Usage : who

Returns a list of players currently online.

## Wield

Usage : wield `<object>`

Allows a player to use the object as a weapon. The object must be a
weapon, of course.

# Administrator commands

Administrator commands follow the same syntax rules as player commands.
The server will only recognize these commands if the player entering the
command has the 'Admin' flag set to 'yes'.

## Advance

Usage : advance `<player>` `<level>`

Advances `<player>` to the specified `<level>` and adjusts the player's
hitpoints and other statistics as necessary. This command can also be
used to demote a player.

## GoTo

Usage : goto `<RoomId>`

Immediately go to the room specified.

## GoToArrive

Usage : gotoarrrive {`<text>` | none}

Sets player's arrival message to the `<text>`. If `<text>` is not
specified, the player's current arrival message is displayed. Specifying
'none' will remove the player's arrival message. Note: text must be less
than 61 characters.

## GoToDepart

Usage : gotodepart {`<text>` | NONE}

Sets player's depart message to the `<text>`. If `<text>` is not
specified, the player's current depart message is displayed. Specifying
'none' will remove the player's depart message. Note: text must be less
than 61 characters not counting color codes.

## Invisible

Usage : Invisible {ON | OFF}

Player becomes invisible to all other players.

## Load

Usage : load OBJ{ect} | MOB{ile} `<id>`

Loads the object specified by `<id>` into the players inventory or loads
the mobile specified by `<id>` into the room.

## Refresh

Usage : refresh `<target>`

Refreshes the 'in memory array' specified by `<target>`. Currently the
only valid `<target>` is 'commands'.

## Restore

Usage : restore `<player>`

Restores a player's hit points to the maximum value.

## RoomInfo

Usage : roominfo ON | OFF

When on, displays hidden room information next the room name.

## Stop

Usage : stop

Stops the game immediately, but first all players are saved and forced
logged off.

## Where

Usage : where `<target>`

Specify a player, mobile, or object name in place of `<target>`. The
where command always searches all player, mobile, and object locations
for the specified `<target>`. A unique name must be specified for
target. For example, 'where rat' will return nothing, if the only rats
are SmallRat,

LargeRat, and SewerRat. You would instead use 'where smallrat', 'where
largerat', or 'where sewerrat'.

# World building

This section is intended to provide additional information concerning
critical components of a well rounded virtual world.

## Armor class

Armor class is the sum of the 'armor value' of all objects worn by the
player. Armor class should reflect the amount and quality of armor worn
by a player. The primary purpose of armor is to reduce the amount of
damage taken during a fight. The following list gives the major factors
used to determine the effectiveness of armor.

1.  Armor objects must be created with appropriate 'cost' and 'armor
    value'.

2.  A player wearing the 'best' armor will still take damage during a
    fight.

3.  A maximum value for armor class must be established.

It is difficult to give much advice concerning item \#1, but it might
help to design sets of armor for beginner, intermediate, and advanced
players. The 'cost' of a piece of armor should reflect its 'armor
value', meaning the better the armor, the more it cost.

Armor reduces damage, so item \#2 is defined by establishing the Maximum
Damage Reduction Percent or (MDRP). The MDRP is specified as number
ranging from 1 to 99. For example, setting MDRP to 60, means the damage
taken by a fully armored player is reduced by 60%.

Item #3 is defined as Maximum Armor Class or MAC. Care must be taken
when designing armor objects so that a player's armor cannot exceed the
MAC value. On the other hand, making this value to large could render
the player's armor virtually ineffective.

The player's armor class, MAC, and MDRP are used to determine the
player's Damage Reduction Percentage or DRP which in turn is used to
reduce the amount damage taken during a fight. The following describes
the calculation of DRP.

### Damage Reduction Percentage calculation

Term definitions:
```
POP  Percent Of Protection, ranges from 0% to 100% protected
PAC  Player's Armor Class, ranges from 0 to MAC
MAC  Maximum Armor Class, your choice
MDRP Maximum Damage Reduction Percent, your choice
DRP  Damage Reduction Percent, what we are calculating
```

Calculations:  
POP = PAC / MAC * 100  
DRP = POP * MDRP / 100  

Example:  
PAC = 12 (Add up player's armor values)  
MAC = 300 (your choice)  
MDRP = 60% (your choice)  
POP = 12 / 300 * 100 = 4%  
DRP = 4 * 60 / 100 = 2.4%  

In this example, the player's armor provides 4% protection which in turn
will reduce the damage taken by 2.4%.

### Combat example

Let's assume Ixaka is fighting a viper.

Given:  
Ixaka's armor class is 17 (PAC)  
Maximum armor class is 300 (MAC)  
Maximum damage reduction percent is 60 (MDRP)  
Viper bites Ixaka for 15 points of damage (POD)  

Damage to Ixaka (DTI):
```
DTI = floor(POD – DRP)
DTI = POD - POP * MDRP / 100
DTI = POD - ((PAC / MAC * 100) * MDRP / 100)
DTI = 15 – ((17 / 300 * 100) * 60 / 100)
DTI = 15 – 3.39%
DTI = 14.49
DTI = floor(14.49) = 14
```