# OMugs – Online Multi-User Game Server
# Coding Guide


## Table of Contents

- [Introduction](#introduction)
- [Code Overview](#code-overview)
  - [Style](#style)
  - [Organization](#organization)
  - [Communications](#communications)
  - [Variables](#variables)
  - [Go Source Files](#go-source-files)
- [Coding Conventions](#coding-conventions)
  - [Variable Names](#variable-names)
  - [!Terse](#terse)
  - [Functions](#functions)
- [The Codebase](#the-codebase)
  - [BigDog](#bigdog)
  - [OpenPort(int port)](#openportint-port)
  - [InitDescriptor](#initdescriptor)
  - [CheckForNewConnections](#checkfornewconnections)
  - [NewConnection](#newconnection)
  - [SockRecv](#sockrecv)
  - [SockSend](#socksend)
  - [CommandParse](#commandparse)
  - [pDnodeActor](#pdnodeactor)
  - [SendToRoom](#sendtoroom)
- [Fighting](#fighting)
  - [Directories](#directories)
  - [Files](#files)
- [Armor Class](#armor-class)
  - [Definitions](#definitions)
  - [Damage To Player](#damage-to-player)
- [How To Add A New Command](#how-to-add-a-new-command)
  - [Command Design](#command-design)
  - [Command Table](#command-table)
  - [Define Function](#define-function)
  - [Command Parsing](#command-parsing)
  - [Message and Prompt](#message-and-prompt)
  - [Do Checking](#do-checking)
- [Osi – OMugs Script Interpreter](#osi--omugs-script-interpreter)
  - [Major Components](#major-components)
  - [Front End](#front-end)
  - [Back End](#back-end)

# License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any means.

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

The OMugs codebase development environment is Visual Studio Code with an
execution environment of Microsoft Windows or Linux. This codebase is the
the result of converting from the orignal C++ code to Go.

## Style

There is a distinct *style* used when writing code for OMugs. Hallmarks
of that style include, consist indentation, no tabs, readable function and
variable names using PascalCase.

## Organization

The server performs three major services. First, establish a TCP/IP
connection and begin listening on the specified port. Second, connect
and disconnect players. Third, accept and process player commands. The
first service is a one-time event, whereas the second and third services
are repeated, each in turn, as long as the server is active.

## Communications

OMugs is a telnet server and uses a small subset of Socket functions to
communicate with the players. The two most often used functions are
*recv* and *send*. Sends and receives are done using nonblocking
function calls thereby enabling multiple sockets with simultaneous
pending operations.

## Variables

OMugs functionality centers around sending and receiving text data.
Therefore a significant percentage of the variables are strings. There
are a number of custom string functions to assist with string handling.

OMugs maintains a persistent state by not storing anything about the
state of the server in memory. Everything is externalized to text files
and then read from the file when it is needed. If the server crashed,
upon restart the state of the server is in the external text files ready
for use.

## Go Source Files

The Go source file are, in general, organized around major components of
the server. Player.go, Mobile.go, Room.go, etc.

# Coding Conventions

A consistent look is critical to the development and maintainability of
the OMugs codebase. The OMugs coding motto: KIASAP - Keep It As Simple
As Possible. OMugs code will not be the most efficient, terse, tight
code. It will however, be as maintainable as possible. In spite of the
KIASAP motto, there will be portions of code that even it's creator will
be baffled as to what the code does, even when its only been a few days
since the code was written coding it. These occurrences should be minimized
as much as possible.

## Variable Names

Variable names should be as readable as possible. For example, use a
name like 'PlayerStateWaitNameConfirmation' as opposed to something like
'PlStNmCn' or worse 'plstnmcn'. Each word is capitalized and does not
use special characters to separate words. Naming variables using a
single letter such as 'n' is acceptable, but only in limited situations
like:

int n

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

## Functions

Function names follow the same convention as variable names.

Return from functions as soon as possible. For example, in the DoLook
function there is a return as soon as it is determined that the Look
command cannot be completed.


# The Codebase

This section outlines the flow of the code and is intended as a way to
'get into' the code as opposed to explaining everything in great detail.

BigDog

🡪 OpenPort  
🡪 InitDescriptor  
🡪 CheckForNewConnections  
> 🡪 NewConnection  
>  
> 🡪 Dnode constructor  
🡪 RecvIt  
> 🡪 SendIt  
>  
> 🡪 CommandParse  

## BigDog

Yep, this is it, *THE* BIG DOG. This function is the origin of all that
is OMugs. An old sage once said, "If you can't run with the big dogs,
stay on the porch".

## OpenPort(int port)

Begins listening on the assigned port for new connections.

## InitDescriptor

Initializes the pointers used to manage a linked list of connections.

## CheckForNewConnections

Check for and accepts new connections

## NewConnection

This function handles new connections and establishes a nonblocking
socket.


## SockRecv

The code in this function is all contained within a loop that checks
each connection. See CheckForNewConnections for the 'loop through the
connections code'. It primary function is to receive player input, send
output to players, and call the command parser.

## SockSend

The code in this function could have just stayed in SockRecv, but the
nesting was getting too deep, so it was split out.


## CommandParse

This function parses player commands and then calls the appropriate
function. The following lines can be used to divide the command parser
into logical sections.

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
command (if any), and pDnode3 refers to others.

## SendToRoom

This function is used to send messages to players that are in the same
room as the player who issued a given command.

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
is a text file called 'ValidCommands.txt' in the '\Library\
directory. Set the authorization for the command to 'admin' or 'all'.

## Define Function

For the purposes of this discussion assume you're implementing the 'eat'
command. All command functions are named 'Do\<Command\>' where
\<Command\> is the command used by the player with only the first letter
capitalized. So, the function name is DoEat.

Add the DoEat function to Communication.go in alphabetical order, right
after the DoDrop function.

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
CommandParse calls the appropriate function.

## Message and Prompt

The 'eat' command should send a message to the player and all players in
the room. Something like:  
Ixaka sees: You take a big bite from a red apple.  
Players in the room see: Ixaka takes a big bite from a red apple.

## Do Checking

Some of the more common validation checks are: command syntax, command
parameters, and player position. Eat is very similar to DoDrop command,
so you can copy the parts needed from the DoDrop function. In the validation
code, the player position, command syntax, object existence, and object type are checked.

After all the validation is completed, you know that the player position
is Ok and that the command, object, and object type are all valid. Now,
we finally get to do something, like effect the player's hunger level,
remove the item from the player's inventory, and send messages to the
player and all players in the room.

# Osi – OMugs Script Interpreter

Osi allows the game developer to enhance the virtual world by adding
triggers to rooms and NPCs. For example, when a player gives an NPC a
particular object, the NPC could give the player a special object in
return. **Not implemented in this version of OMugs**.

## Major Components

Osi is a recursive descendant parser containing a front end to interpret
the script source code and a back end to execute the translated script.
 Each time a player enters a room or hails an
NPC, Osi is called and check for the existance of a script file for that
room or NPC. If no script is found, Osi exits. The Buffer, Scanner,
Token, Parser, and Symbol make up the front end. Icode is common to the
front end and back end. Executor and RunStack make up the back end.