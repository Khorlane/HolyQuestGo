# OMugs – Online Multi-User Game Server
# Coding Guide

## Table of Contents

- [Introduction](#introduction)
- [Code Overview](#code-overview)
	- [Style](#style)
	- [Organization](#organization)
	- [Communications](#communications)
	- [Variables](#variables)
	- [Classes](#classes)
- [Coding Conventions](#coding-conventions)
	- [Variable Names](#variable-names)
	- [!Terse](#terse)
	- [Functions](#functions)
	- [Class Template](#class-template)
- [The Codebase](#the-codebase)
	- [BigDog](#bigdog)
	- [OpenPort(int port)](#openportint-port)
	- [InitDescriptor](#initdescriptor)
	- [CheckForNewConnections](#checkfornewconnections)
	- [NewConnection](#newconnection)
	- [Dnode(int SocketHandle)](#dnodeint-sockethandle)
	- [SockRecv](#sockrecv)
	- [SockSend](#socksend)
	- [CommandParse](#commandparse)
	- [pDnodeActor](#pdnodeactor)
	- [SendToRoom](#sendtoroom)
- [Fighting](#fighting)
	- [Directories](#directories)
	- [Files](#files)
	- [Pseudo Code](#pseudo-code)
- [Armor Class](#armor-class)
	- [Definitions](#definitions)
	- [Damage To Player](#damage-to-player)
- [How To Add A New Command](#how-to-add-a-new-command)
	- [Command Design](#command-design)
	- [Command Table](#command-table)
	- [Define Function](#define-function)
	- [Command Parsing](#command-parsing)
	- [Message and Prompt](#message-and-prompt)
	- [Do Something](#do-something)
- [Osi – OMugs Script Interpreter](#osi--omugs-script-interpreter)
	- [Major Components](#major-components)
	- [Front End](#front-end)
	- [Back End](#back-end)

# License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

For more information, please refer to <http://unlicense.org/>

Author  
Stephen L Bryant

Revision 1.0 April 22, 2001
Revision 1.1 December 4, 2002
Revision 1.2 December 22, 2025

Revision tracking after December 22, 2025 is maintained via git.

# Introduction

This document is intended to give a 'would be' OMugs coder a head start
on understanding the OMugs codebase. Most of the code is NOT discussed
here, only selected portions of the codebase are discussed here. What
about the rest? Go read the code!. There are several code snippets in
this document. These snippets are not intended to be complete or be the
actual code in the current version of OMugs. They are used to explain
key portions of the code. As of this writing, *the author has written
all OMugs code* and this document is prerequisite reading and
understanding for any new coders and a reminder to the author.

# Code Overview

The OMugs codebase development platform and running environment is
Microsoft Windows (WIN) and Microsoft Visual C++ (MSVC). The design
takes advantage of WIN and MSVC components whenever it is considered to
make the OMugs codebase more developer friendly. Being portable or not
being portable is not a coding criteria, if it is portable that is fine,
and if it is not portable that is fine too.

## Style

There is a distinct *style* used when writing code for OMugs. Hallmarks
of that style include, consist indentation, braces are always on a
separate line, severely limited use of macros, no function overloading,
no pointers to functions, limited use of pointers, readable function and
variable names, statements like IF, FOR, WHILE are always followed by a
set of braces, even if there is only one statement inside the braces.

## Organization

The code is divided into two major components. This first is the windows
aspect, which is used to start the server. The second is the server,
which runs as thread, thereby allowing control to return to the OMugs
window component while the OMugs server component continues to run.

The windows component was generated using MSVC application wizard and is
mostly left unchanged. The following comment lines indicate changes to
the generated code:

/\* slb code begin block \*/

/\* slb code end block \*/

The windows component calls a function call BigDog, and thus the server
doth begin to run.

The server component performs three services. First, establish a TCP/IP
connection and begin listening on the specified port. Second, connect
and disconnect players. Third, accept and process player commands. The
first service is a one-time event, whereas the second and third services
are repeated, each in turn, as long as the server is active.

## Communications

OMugs is a telnet server and uses a small subset of WinSock functions to
communicate with the players. The two most often used functions are
*recv* and *send*. Sends and receives are done using nonblocking
function calls thereby enabling multiple sockets with simultaneous
pending operations.

## Variables

OMugs functionality centers around sending and receiving text data.
Therefore a significant percentage of the variables are c strings. MSVC
provides a MFC called CString that has numerous useful functions to
manipulate strings. OMugs relies heavily on CString functions. Most
variables are CString, int, or bool. All external data is stored in
standard windows text files and can be easily viewed and changed using
the Windows Notepad program.

As one reviews the OMugs codebase, one might be struck by the fact that
there is little use of such contracts as pointers, structures, linked
lists, arrays, binary trees, ques, stacks, etc. There are actually a few
linked lists and arrays, but those are definitely exceptions.

OMugs maintains a persistent state by not storing anything about the
state of the server in memory. Everything is externalized to text files
and then read from the file when it is needed. If the server crashed,
upon restart the state of the server is in the external text files ready
for use.

## Classes

Classes are used as means of organizing the code. Classes are not
derived from other classes. OMugs is viewed as one large program with
classes providing a means by which to group functions. Classes are not
designed to plug into another application, so instead of calling a
function to set the player name in the Player class, just set it
directly. For example:

pDnodeActor-\>pPlayer-\>Name = pDnodeActor-\>PlayerName;

# Coding Conventions

A consistent look is critical to the development and maintainability of
the OMugs codebase. The OMugs coding motto: KIASAP - Keep It As Simple
As Possible. OMugs code will not be the most efficient, terse, tight
code. It will however, be as maintainable as possible. In spite of the
KIASAP motto, there will be portions of code that even it's creator will
be baffled as to what the code does, only months or even days after
coding it. These occurrences should be minimized as much as possible.

## Variable Names

Variable names should be as readable as possible. For example, use a
name like 'PlayerStateWaitNameConfirmation' as opposed to something like
'PlStNmCn' or worse 'plstnmcn'. Each word is capitalized and does not
use special characters to separate words. Naming variables using a
single letter such as 'n' is acceptable, but only in limited situations
like:

int n;

...

n = WordList.Find(Word);

if (n == -1)

...

## !Terse

Minimize the terseness of the code. Terse code can make it more
difficult to debug and/or explore the code. For example, by coding

n = WordList.Find(Word);

if (n == -1)

instead of

if (WordList.Find(Word) == -1)

it is possible to see the result of the 'Find' function when using the
MSVC interactive debugger.

Avoid writing complex if statements like

if (condition1 && condition2 \|\| (condition4 && condition5))

Use white space

x = y + z

instead of

x=y+z

Use macros only for value substitution NOT code substitution. Macros are
limited to substitutions like

\#define GRP_LIMIT 4

\#define HELP_FILE_NAME ".\\Library\\"

## Functions

Function names follow the same convention as variable names. In general
variables should be passed by value not by reference. Therefore code

void FirstFunction(int FirstValue)

instead of

void FirstFunction(int \*FirstValue)

This works well because variables being passed are either a member of
the class or the variable is not changed by the function. An example of
an exception to this rule is

\*pDnodeActor because it is a pointer.

Return from functions as soon as possible. For example, in the DoLook
function there is a return as soon as it is determined that the Look
command cannot be completed.

void Communication::DoLook(Dnode \*pDnodeActor, CString CmdStr)

{

Dnode \*pDnode2;

Mobile \*pMobile;

bool IsPlayer;

CString TargetName;

CString TmpStr;

if (IsSleeping(pDnodeActor))

{ // Player is sleeping, send msg, command is not done

return;

}

. . . .

}

## Class Template

The following is the code template used to define all classes. The
PLAYER class is defined in this template.

/\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

\* OMUGS - Online Multi-User Game Server \*

\* File: Player.h \*

\* Usage: Define Player class \*

\* Author: Steve Bryant (stevebryant@www.holyquest.org) \*

\* \*

\* This program belongs to Stephen L. Bryant. \*

\* It is considered a trade secret and is not to be \*

\* divulged or used by parties who have not received \*

\* written authorization from the owner. \*

\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*/

\#ifndef PLAYER_H

\#define PLAYER_H

/\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

\* Includes \*

\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*/

\#include "Config.h"

/\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

\* Define Player class \*

\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*/

class Player

{

// Friend classes

// Public functions static

public:

// Public functions

public:

// Private functions

private:

// Public variables

public:

// Private variables static

private:

// Private variables

private:

// File variables

public:

};

\#endif

# The Codebase

This section outlines the flow of the code and is intended as a way to
'get into' the code as opposed to explaining everything in great detail.
It is best used with the OMugs MSVC project and frequent use of the
'Find in Files' option. The following is an outline showing the flow of
the code at a high level.

BigDog

🡪 OpenPort

🡪 InitDescriptor

🡪 CheckForNewConnections

> 🡪 NewConnection
>
> 🡪 Dnode object constructor

🡪 RecvIt

> 🡪 SendIt
>
> 🡪 CommandParse

## BigDog

Yep, this is it, *THE* BIG DOG. This function is the origin of all that
is OMugs. An old sage once said, "If you can't run with the big dogs,
stay on the porch". So here is what BigDog looked like at some point
during development. Although this is not exactly what you currently find
in the codebase, this version has the critical components.

1 void BigDog()

2 {

3 int MilliSecondsToSleep;

4

5 srand((unsigned)time(NULL));

6 StateRunning = true;

7 MilliSecondsToSleep = 250;

8 Log::OpenFileLog();

9 Log::LogIt("OMugs has started");

10 Communication::OpenPort(PORT_NBR);

11 Descriptor::InitDescriptor();

12 while (StateRunning)

13 {

14 Sleep(MilliSecondsToSleep);

15 if (!StateStopping)

16 {

17 Communication::CheckForNewConnections();

18 }

19 Communication::RecvIt();

20 }

21 Log::LogIt("OMugs has stopped");

22 Log::CloseFileLog();

23 Communication::ClosePort();

24 }

Line 5 – Uses time to seed the random number function

Line 8 – Opens the Log

Line 10 – Begins listening on port specified

Line 11 – Initializes the list of connections

Line 12 thru 20 – The main server loop

Line 14 – Sleep, so this process doesn't hog the cpu

Line 17 – Checks for new connections and creates new connections as
needed

Line 19 – Receives and sends text to and from the players and
disconnects players.

Line 22 – Closes the Log

## OpenPort(int port)

Begins listening on the assigned port for new connections.

## InitDescriptor

Initializes the pointers used to manage a linked list of connections.

## CheckForNewConnections

The following code can be found very near to top of
CheckForNewConnections. This is the code that is used in numerous places
throughout OMugs to loop through the connections.

1 Descriptor::SetpDnodeCursorFirst();

2 while (!Descriptor::EndOfDnodeList())

3 {

4 pDnodeActor = Descriptor::GetDnode();

5 FD_SET(pDnodeActor-\>DnodeFd, &InpSet);

6 FD_SET(pDnodeActor-\>DnodeFd, &OutSet);

7 FD_SET(pDnodeActor-\>DnodeFd, &ExcSet);

8 Descriptor::SetpDnodeCursorNext();

9 }

Line 1 – Point to the first connection in the list

Line 2 – Is this the end of the list

Line 4 – Get pointer to the current Dnode object (which is the
connection)

Line 5 thru 7 – Do stuff with the current connection (this can be
anything)

Line 8 – Advance the pointer to the next connection

Near the end of CheckForNewConnections is the follow line of code, which
will add the new connection to the list of connections.

NewConnection(ListenSocket);

## NewConnection

This function handles new connections and establishes a nonblocking
socket. Near the bottom of this function are the following two lines of
code

1 pDnodeActor = new Dnode(SocketHandle);

2 FdList-\>AppendIt(pDnodeActor);

Line 1 – Constructs a new Dnode object, which used to communicate with
players

Line 2 – Appends this new connection to the list of connections

## Dnode(int SocketHandle)

The Dnode class constructor initializes the variables that are used to
communicate with the players and track the 'state' of the connection.

1 Count++;

2 DnodeFd = SocketHandle;

3 pPlayer = NULL;

4 PlayerInp = "";

5 PlayerName = "";

6 PlayerNewCharacter = "";

7 PlayerOut = "";

8 PlayerPassword = "";

9 PlayerWrongPasswordCount = 0;

Line 1 – Keeps a count of the number of connections

Line 2 – Assigns the socket handle to the DnodeFd variable

Line 3 – Pointer to the player object

Line 4 – Player input is kept here until it is processed (see function
RecvIt)

Line 5 – At this point in the code there is no player object, so store
player name here

Line 6 – A 'Y' or 'N' switch

Line 7 – Player output is kept here until it is processed (see function
SendIt)

Line 8 – Player's response to password question

Line 9 – Count of attempts to enter correct password, 3 chances to get
it right

## SockRecv

The code in this function is all contained within a loop that checks
each connection. See CheckForNewConnections for the 'loop through the
connections code'. It primary function is to receive player input, send
output to players, and call the command parser.

1 ::recv(pDnodeActor-\>DnodeFd, InpStr, MAX_INPUT_LENGTH-1, 0);

2 SendIt(pDnodeActor, BufferOut);

3 CommandParse(pDnodeActor);

Line 1 – Receives player input and places it in InpStr

Line 2 – Calls the send routine

Line 3 – Calls the command parser, if there is input

## SockSend

The code in this function could have just stayed in SockRecv, but the
nesting was getting too deep, so it was split out.

1 Written = ::send(pDnodeActor-\>DnodeFd, arg, Length, 0);

Line 1 – Send output to player, The variable 'Written' is the number of
characters the '::send' was able to send. 'Written' is compared against
the total bytes that should have been sent and used to preserve the
unsent characters to be sent via the next '::send'.

## CommandParse

This function parses player commands and then calls the appropriate
function. The following lines can be used to divide the command parser
into logical sections.

1 /\* Isolate a single player command \*/

2 /\* Player logon \*/

3 /\* Separate the command from the command string \*/

4 /\* COLOR command \*/

5 /\* SOCIAL command \*/

6 /\* Bad command \*/

1 – Each player command must end with a line return. Due to nonblocking
communication between a player and the server, input must be accumulated
until a line feed is found. Then it is assumed that a complete command
has been found and is ready to be processed.

2 – When a player first connects to OMugs, they must reply to a series
of prompts to complete the logon process.

3 – This section of code separates the one word command from the rest of
the player input. Some commands require additional translations. For
example, a player may enter 'n' instead of 'go north'. The command
processor does not understand 'n', but does understand 'go north'.

4 – This begins the section of code that determines what command was
entered and then calls the appropriate function to process that command.

5 – If the command does not match any of the previous commands, then
maybe it is a social command, like smile, laugh, grin, frown, etc.
Social commands are stored in an external file and therefore can be
added or deleted without recompiling the code.

6 – The command parser gives up trying to figure out what the player
entered.

## pDnodeActor

This variable points to the connection being processed and has a pointer
the corresponding player object. Throughout the code pDnodeActor refers
to the player who entered the command, pDnode2 is the target of the
command (if any), and pDnode3 refers to others. Some common references
using pDnodeActor are:

pDnodeActor-\>PlayerOut Output to be sent to player

pDnodeActor-\>pPlayer-\>CreatePrompt() Player class function that
creates the prompt

## SendToRoom

This function is used to send messages to players that are in the same
room as the player who issued a given command. For example, player Ixaka
issues the following command

'say Hi!'. As a result all players in the same room as Ixaka with see
the following message: Ixaka says: Hi!

When Ixaka issues the 'say' command, the function SendToRoom is called.
The SendToRoom function requires the following parameters:

CString TargetRoomId, Dnode \*pDnodeActor, Dnode \*pDnode2, CString
MsgText

TargetRoomId is the RoomId of the originating player

\*pDnodeActor is the pointer to the originating player

\*pDnode2 is the pointer to the target player

MsgText is the message to be send to the players in the TargetRoomId

\*pDnodeActor and \*pDnode2 are used to control who actually gets the
message and these pointers can be NULL, meaning all players will see the
message. SendToRoom therefore has several variations:

SendToRoom variation 1

In Communication::DoGive the player receives a message, the target
player receives a different message, and everyone else in the room
receives a third variation of the message. The originating player and
target player **should not** get a message from the SendToRoom function,
so the following code is used:

SendToRoom(pDnodeActor-\>pPlayer-\>RoomId, pDnodeActor, pDnode2,
GiveMsg);

This causes the message to be sent everyone in the room, except the
originating player and the target player.

SendToRoom variation 2

In Communication::DoDrop the player receives a message and everyone else
in the room receives a variation of the message. Only the origination
player **should not** get a message from the SendToRoom function, so the
following code is used:

SendToRoom(pDnodeActor-\>pPlayer-\>RoomId, pDnodeActor, pDnodeActor,
DropMsg);

This causes the message to be sent to everyone in the room, except the
originating player.

SendToRoom variation 3

In Communication::DoHail everyone in the room should receive the message
when an NPC responds to a 'hail', so the following code is used:

SendToRoom(RoomId, NULL, NULL, MobileMsg);

This causes the message to be sent to everyone in the room.

# Fighting

In keeping with the persistent world theme, all information concerning
fights is kept in external text files. This results in a somewhat
complex set of directories, files, and the relationships between them.
The directories and files used to track all fights are found in the
\Running\Fight directory. When a fight begins, any player involved is
placed in the FIGHTING state. A player cannot 'do' certain things while
in the FIGHTING state. For example, 'sleeping' is not allowed while
fighting. Players cannot fight each other, so only mobile / player
fights need to be tracked. A player can attack only one mobile at a time
and a mobile can attack only one player at a time. Players and mobiles
always automatically attack whomever or whatever attacked them unless
they are already in a fight with someone or something else. A player or
mobile cannot switch targets during a fight, the fight is to the death
or until the player or mobile flee the scene.

## Directories

Fights range from simple, such as Zed attacks small rat and small rat
fights back, to complex, such as:

Zed attacks small rat

Small rat fights back and begins biting Zed

Zeke attacks small rat

Large rat goes aggro on Zeke and starts biting him

In a room with several players and several mobiles, things can get
complicated. Three directories are used to track who is fighting whom.

**MobPlayer** tracks which player a mobile is attacking

**MobStats** contains 'fight stats' copied from \Library\Mobiles

**PlayerMob** tracks which mobile a player is attacking

## Files

Files, in the above-mentioned directories, are created when a fight
begins and are deleted when the fight ends. Directories 'MobPlayer' and
'PlayerMob' both contain files which are named using the player's name,
like Zed.txt. The existence of a file in one of these directories
corresponds directly with the player being in the FIGHTING state.

The MobStats directory contains one file for each mobile engage in a
fight. The statistics, copied from \Library\Mobiles, include: (see Admin
document for explanations of these)

HitPoints decremented each time the mobile is damaged

Armor used to calculate damage done by a player

Attack determines fight message contents

Damage used to calculate damage done to a player

ExpPoints experience awarded to a player when the mobile dies

## Pseudo Code

This is the pseudo code for specific events during the life of a fight.
The 'fight code' processes the information found in the \Fight directory
once during each player's turn, if the player is in the FIGHTING state.
This is pseudo code describes the events that start and stop a fight.
The names of mobiles and players in these examples are made up and may
or may not exist in the \Players and \Library\Mobiles directories. These
are only three sample cases, there are more possibilities.

Case 1 – Zed issues the kill command, as in 'kill rat'

Search room for a rat and select the first rat found, i.e. smallrat

Remove rat from room

Assign rat the next mobile number

> (use \Running\Control\NextMobileNumber.txt)

Create \Running\Fight\PlayerMob\Zed.txt

containing one line -\> SmallRat.7

Create \Running\Fight\MobPlayer\Zed.txt

containing one line -\> SmallRat.7

Create \Running\Fight\MobStats\SmallRat.7.txt

containing the following lines,

copied from \Library\Mobiles\SmallRat.txt,

HitPoints: 5

> Armor: 1
>
> Attack: Bite
>
> Damage: 5
>
> ExpPoints: 25

Now the fight code will use these files to cause Zed and the small rat
to whack each other.

Case 2 – Zeke is exploring the forest and a black wolf attacks without
provocation. This code required in this situation is same as case 1,
except the names are changed and the fight is started automatically due
to the player entering a room containing an 'aggro' mobile.

Select which player the wolf is to attack, in this case, Zeke

Remove wolf from room

Assign wolf the next mobile number

> (use \Running\Control\NextMobileNumber.txt)

Create \Running\Fight\PlayerMob\Zeke.txt

containing one line -\> BlackWolf.101

Create \Running\Fight\MobPlayer\Zeke.txt

containing one line -\> BlackWolf.101

Create \Running\Fight\MobStats\BlackWolf.101.txt

containing the following lines,

copied from \Library\Mobiles\BlackWolf.txt,

HitPoints: 25

> Armor: 10
>
> Attack: Bite
>
> Damage: 25
>
> ExpPoints: 100

Now the fight code will use these files to cause Zeke and the black wolf
to whack each other.

Case 3 – Zeke is lost in the desert and they stumble upon a den of
vipers. The vipers are aggro and they attack Zeke. Assume that there are
two vipers.

Select a target for each viper (let's say Zeke)

Remove the vipers from the room

Assign each viper a number (let's use 98 and 99)

Create \Running\Fight\PlayerMob\Zeke.txt

containing one line -\> Viper.98

Create \Running\Fight\MobPlayer\Zeke.txt

containing one line -\> Viper.98

Create \Running\Fight\MobPlayer\Zeke.txt

containing one line -\> Viper.99

Create \Running\Fight\MobStats\Viper.98.txt

> and create \Running\Fight\MobStats\Viper.99.txt

both containing the following lines,

copied from \Library\Mobiles\Viper.txt,

HitPoints: 50

> Armor: 10
>
> Attack: Bite
>
> Damage: 50
>
> ExpPoints: 125

Now the fight code will use these files to cause Zeke to whack viper.98
and both viper.98 and viper.99 to whack Zeke.

# Armor Class

See admin.doc section of the same name for additional information
concerning armor class.

## Definitions

The Player's Armor Class (PAC) Armor class is the sum of the 'armor
value' of all objects worn by the player. Maximum Armor Class (MAC) and
Maximum Damage Reduction Percent (MDRP) are specified in Config.h and
used in BigDog.cpp to calculate the Player Armor Class Magic Number
(PACMN). When the PACMN is multiplied the player's armor class, the
result is the Damage Reduction Percent (DRP). This bit of algebraic
manipulation reduces the number calculations required during each fight
iteration.

PACMN calculation:

MAC = 300

MDRP = 60

PACMN = 1 / MAC \* MDRP / 100

PACMN = 1 / 300 \* 60 / 100 = .002

## Damage To Player

The damage taken by the player (DTP) during a fight is calculated as
follows:

DTP = floor(Damage – (PAC \* PACMN \* Damage) )

# How To Add A New Command

This section describes the tasks and the proper sequence required to add
a new command to OMugs. Following this procedure should minimize the
risk of disruption to the players.

## Command Design

Decide what the command is going to 'do'. This may seem obvious, but it
is important the think through what effect the command is going to have
and how will that effect be accomplished. For example, the 'get' command
takes an item from the room and places in the player's inventory. How is
that accomplished?

All commands result in some text being sent to the player. Here are some
questions you should consider:

- What is going to be sent to the player issuing the command?

- Does the command have a target player? If so, what do they see?

- What do other players in the room see?

- What do all players logged on see?

This table shows the possible message relationships between the player
issuing the command and the rest of the players.

|             |            |            |                    |                        |
|-------------|:----------:|:----------:|:------------------:|:----------------------:|
| **Command** | **Player** | **Target** | **Others in Room** | **Others not in room** |
| who         |     Y      |     N      |         N          |           N            |
| tell        |     Y      |     Y      |         N          |           N            |
| give        |     Y      |     Y      |         Y          |           N            |
| say         |     Y      |     N      |         Y          |           N            |
| ooc         |     Y      |     Y      |         Y          |           Y            |

## Command Table

Add the command, in alphabetical order, to the command table. This table
is a text file called 'ValidCommands.txt' in the '\Running\Control\\
directory. Set the authorization for the command to 'admin' like the
'advance' command. This will allow you to test the command without
affecting the players. When ready to release the command the players,
just change 'admin' to 'all' or a number, indicating the level the
player must be before using this command. Of course, if this is an
administrator only command, no further changes are required.

## Define Function

For the purposes of this discussion assume you're implementing the 'eat'
command. All command functions are named 'Do\<Command\>' where
\<Command\> is the command used by the player with only the first letter
capitalized. So, the function name is DoEat.

Add the function definition to Communication.h by finding the other
commands and inserting DoEat in alphabetical order. All command
functions are passed \*pDnodeActor. Only pass CmdStr if the command does
not expect or care if the player typed anything after the command (see
afk command).

void static DoDrop(Dnode \*pDnodeActor, CString CmdStr);

**void static DoEat(Dnode \*pDnodeActor, CString CmdStr);**

void static DoEmote(Dnode \*pDnodeActor, CString CmdStr);

Add the DoEat function to Communication.cpp in alphabetical order, right
after the DoDrop function.

/\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

\* Eat command \*

\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*/

void Communication::DoEat(Dnode \*pDnodeActor, CString CmdStr)

{

}

Now you should be able to compile OMugs without any errors. This is a
good practice because this confirms that the compiler is happy with your
changes so far and an unhappy compiler is a *bad thing.* If you fire up
OMugs, logon, and type 'eat' you will get a message produced by the
CommandParse function saying that the command is not available.

## Command Parsing

The CommandParse function finds the next new line in
pDnodeActor-\>PlayerInp and assigns that string to CmdStr and removes
that string from pDnodeActor-\>PlayerInp. For example, if
pDnodeActor-\>PlayerInp contained the following:

'eat apple\nsay Yummy\n'

that would mean the player had typed 'eat apple' and then 'say Yummy'
before the eat command was processed. CommandParse would assign 'eat
apple' to CmdStr and remove 'eat apple\n' from pDnodeActor-\>PlayerInp
leaving 'say Yummy\n' to be parsed later.

After validating the command and making sure it is not a 'social',
CommandParse calls the appropriate function. The code is simple, just
add the following in alphabetical order:

/\* EAT command \*/

if (MudCmd == "eat")

{

DoEat(pDnodeActor, CmdStr);

return;

}

Compile OMugs and test the eat command. You should not get any messages
this time, in fact you don't even get a prompt. Each time a player
enters a command, they expect something to happen and they expect to get
a message of some sort and a prompt.

## Message and Prompt

The 'eat' command should send a message to the player and all players in
the room. Something like:

Ixaka sees: You take a big bite from a red apple.

Players in the room see: Ixaka takes a big bite from a red apple.

This test code will send a message to the player and a prompt. Add the
code to the DoEat function, compile, and test.

pDnodeActor-\>PlayerOut += "The eat command is not ready, starve!";

pDnodeActor-\>PlayerOut += "\r\n";

pDnodeActor-\>pPlayer-\>CreatePrompt();

pDnodeActor-\>PlayerOut += pDnodeActor-\>pPlayer-\>GetOutput();

## Do Something

So far the new 'eat' command doesn't do very much. It only sends a
message and a prompt. But, DoEat has a firm foundation upon which the
rest of the function can be built. DoEat needs code that will do
something, like modify the player's hunger level and remove the item
eaten from their inventory. Most commands require some validation before
they do something, AFK is an example of an exception.

Some of the more common validation checks are: command syntax, command
parameters, and player position. Eat is very similar to DoDrop command,
so you can copy the parts needed from the DoDrop function. Be sure to
remove the 'test' code first. In the validation code, the player
position, command syntax, object existence, and object type are checked.

//\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

//\* Validate command \*

//\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

if (IsSleeping(pDnodeActor))

{ // Player is sleeping, send msg, command is not done

return;

}

if (Utility::WordCount(CmdStr) == 1)

{ // Invalid command format

pDnodeActor-\>PlayerOut += "Eat what?";

pDnodeActor-\>PlayerOut += "\r\n";

pDnodeActor-\>pPlayer-\>CreatePrompt();

pDnodeActor-\>PlayerOut += pDnodeActor-\>pPlayer-\>GetOutput();

return;

}

//\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

//\* Does player have object? \*

//\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

TmpStr = Utility::GetWord(CmdStr, 2);

ObjectName = TmpStr;

TmpStr.MakeLower();

pObject = Object::IsObjInPlayerInv(pDnodeActor, TmpStr);

if (!pObject)

{

pDnodeActor-\>PlayerOut += "You don't have a(n) ";

pDnodeActor-\>PlayerOut += ObjectName;

pDnodeActor-\>PlayerOut += " in your inventory.\r\n";

pDnodeActor-\>pPlayer-\>CreatePrompt();

pDnodeActor-\>PlayerOut += pDnodeActor-\>pPlayer-\>GetOutput();

return;

}

//\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

//\* Is object food? \*

//\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*\*

pObject-\>Type.MakeLower();

if (pObject-\>Type != "food")

{

pDnodeActor-\>PlayerOut += "You can't eat ";

pDnodeActor-\>PlayerOut += pObject-\>Desc1;

pDnodeActor-\>PlayerOut += ".";

pDnodeActor-\>PlayerOut += "\r\n";

pDnodeActor-\>pPlayer-\>CreatePrompt();

pDnodeActor-\>PlayerOut += pDnodeActor-\>pPlayer-\>GetOutput();

return;

}

After all the validation is completed, you know that the player position
is Ok and that the command, object, and object type are all valid. Now,
we finally get to do something, like effect the player's hunger level,
remove the item from the player's inventory, and send messages to the
player and all players in the room.

//\*\*\*\*\*\*\*\*\*\*\*\*\*\*

//\* Eat object \*

//\*\*\*\*\*\*\*\*\*\*\*\*\*\*

// Send messages

pDnodeActor-\>PlayerOut += "You eat ";

pDnodeActor-\>PlayerOut += pObject-\>Desc1;

pDnodeActor-\>PlayerOut += ".";

pDnodeActor-\>PlayerOut += "\r\n";

EatMsg = pDnodeActor-\>PlayerName;

EatMsg += " eats ";

EatMsg += pObject-\>Desc1;

EatMsg += ".";

SendToRoom(pDnodeActor-\>pPlayer-\>RoomId, pDnodeActor, pDnodeActor,
EatMsg);

// Eat and remove object from player's inventory

pDnodeActor-\>pPlayer-\>Eat(pDnodeActor, pObject-\>FoodPct);

pDnodeActor-\>PlayerOut += pDnodeActor-\>pPlayer-\>GetOutput();

Object::RemoveObjFromPlayerInv(pDnodeActor, pObject-\>ObjectId, 1);

// Clean up and give prompt

delete pObject;

pDnodeActor-\>pPlayer-\>CreatePrompt();

pDnodeActor-\>PlayerOut += pDnodeActor-\>pPlayer-\>GetOutput();

Player::Eat is a new function that must be added to the Player class to
actually change the player's hunger level. Exploration and understanding
of that function is left as an exercise for the reader.

# Osi – OMugs Script Interpreter

Osi allows the game developer to enhance the virtual world by adding
triggers to rooms and NPCs. For example, when a player gives an NPC a
particular object, the NPC could give the player a special object in
return.

## Major Components

Osi is a recursive descendant parser containing a front end to interpret
the script source code and a back end to execute the translated script.
The Osi claseses are Buffer, Scanner, Token, Parser, Symbol, Icode,
Executor, and RunStack. Each time a player enters a room or hails an
NPC, Osi is called and check for the existance of a script file for that
room or NPC. If no script is found, Osi exits. The Buffer, Scanner,
Token, Parser, and Symbol classes make up the front end. The Icode class
is common to the front end and back end. The Executor and RunStack
classes make up the back end.

## Front End

In the Osi front end, the classes are created as follows:

Osi creates Parser

pParser = new Parser(ScriptFileName);

The Parser creates Scanner, Token , Symbol, and Icode

pScanner = new Scanner(ScriptFileName);

pToken = new Token;

pSymbolRoot = new Symbol("");

pIcode = new Icode;

The Scanner creates the Buffer

pBuffer = new Buffer(ScriptFileName);

After initialization, ch points to the first character in the script
file, thereby 'the pump is primed'. Then Parser calls Scanner which
determines the type of token to be parsed. Scanner calls Token which in
turn calls Buffer to read the script file. As the Parser acquires the
information needed, it builds Symbol and Icode. Below is a diagram
showing the flow of data from the the script file to Icode.

Symbol

\|

Script -\> Buffer -\> Token -\> Scanner -\> Parser -\> Icode

## Back End

In the Osi back end, the classes are created as follows:

Osi creates Executor

pExecutor = new Executor(pSymbolRoot, pIcode);

The Executor creates RunStack

pRunStack = new RunStack;

The Executor reads and executes the Icode using Symbol to store values
and RunStack to store intermediate results.
