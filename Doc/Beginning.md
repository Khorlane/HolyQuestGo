# OMugs – Online Multi-User Game Server
# Project Beginning

## Table of Contents

- [License](#license)
- [Introduction](#introduction)
- [Origin](#origin)
	- [CircleMUD](#circlemud)
	- [HolyQuest](#holyquest)
	- [World Editor](#world-editor)
	- [zMUD](#zmud)
	- [Licensing](#licensing)
- [Desire](#desire)
- [OMugs](#omugs)
	- [Definition](#definition)
	- [Combat Support](#combat-support)
	- [Socializing](#socializing)
- [Conclusion](#conclusion)

# License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any means.

For more information, please refer to <http://unlicense.org/>

Author  
Stephen L Bryant 

Revision 1.0 December  5, 2002
Revision 1.1 December 22, 2025

Revision tracking after December 22, 2025 is maintained via git.

# Introduction

This document describes the origin of OMugs which is a 'mud server'. In
order to understand some parts of this document, it is recommended that
the reader educate their self concerning muds and mud servers via the
internet. Use a website like [www.google.com](http://www.google.com) and
search the internet for 'mud faq', 'mudding', or just about any term you
find in this document.

# Origin

The origin of the Online Multi-User Game Server, OMugs, goes back to the
first months of the year 1999. It was then that I first participated in
an activity know as 'mudding'. The first mud I played was called 'The
Inquisition' and was, at that time, a combat oriented mud. I enjoyed
playing The Inquisition and became interested in the computer program
that facilitated so much fun. I learned that the program was called a
Mud Server and set out to explore this new (to me) breed of programs.

## CircleMUD

My explorations, via the internet, soon took me to the CircleMUD
website, [www.circlemud.org](http://www.circlemud.org). Almost all mud
servers are written and run on a unix platform. I prefer a windows
platform, and was pleased to see that CircleMUD had been ported to
windows and someone had been so kind as to make a Microsoft Visual C++
workspace for CircleMUD. I promptly downloaded CircleMUD and began
exploring the code.

Of course, the CircleMUD server being a program and me being a
programmer, I immediately wanted to change something. The first thing to
change was to add the '.txt' extension on the file name for the log so
that I could double-click the file and it would automatically open the
file using notepad. I went on to make a total of 39 changes to the
CircleMUD code, carefully documenting each change.

For reasons described later in this document, I abandoned CircleMUD. It
is a great server and as of this writing continues to have a very
dedicated development staff, an active e-mail list, and numerous active
muds using it.

## HolyQuest

The game that I created using the CircleMUD server was named HolyQuest
by my daughter, Sherry. My son, Chris, became interested in the game and
assisted with ideas for the world of HolyQuest. HolyQuest's theme is
ancient Israel around the time of King David. Creating HolyQuest
required the complete removal of the 'stock' CircleMud world and
replacing it with my own original world. In June of 2000 HolyQuest was
placed online and anyone in the world with internet access could
connect, although the world of HolyQuest only consisted of Jerusalem and
Snake Rock.

## World Editor

The building of the HolyQuest world, albeit small, was quite time
consuming. So I tried several world editor programs and even wrote a
couple of my own. The drive for obtaining or creating a world editor
came from the fact that a CircleMUD world had a fair degree of
complexity and used numerous codes, which I had to constantly look up in
the building document.

The world editor programs I tried only assisted with the creation of
rooms. The architecture of a CircleMUD world is oriented around 'areas'.
An object, say a candle, had to be defined in each area in which you
wanted a candle. There was no editor that I found that had a unified
world view interface.

## zMUD

I was not particularly enamored with the method CircleMUD used to create
in-game mobiles and objects. This led to several attempts to create a
means for loading mobiles and objects using zMUD, a mud client program
that can be downloaded from [www.zuggsoft.com](http://www.zuggsoft.com).

CircleMUD uses numbers to identify each room, mobile, and object. This
made it difficult (for me) to create the files that controlled where
mobiles and objects were created in the game. I was constantly having to
look up the number for a particular mobile or object.

## Licensing

CircleMUD is a derivative of DikuMud and therefore comes with a license
that basically says you cannot make money running a mud using the
CircleMUD server.

# Desire

For me, CirclueMUD had limitations and I had a desire for something
different. In fact HolyQuest, using the CircleMUD server, was already
'different'. In HolyQuest a player can only be a warrior and human. This
is a significant departure from the mud norm of many classes like
cleric, magician, wizard, warrior, paladin, thief, etc and many races
like elf, troll, ogre, human, orc, etc.

My desire for something different can be defined examining the things
that, for lack of a better phrase, I didn't like about CircleMUD.

<u>CircleMUD</u> <u>What I wanted</u>

Magic No magic

Rooms, Mobiles, Objects identified by a number Rooms, Mobiles, Objects
identified by text

Zones or Areas Expand the world without regard to a numbering scheme

Bit vectors and codes e.g. type flag 9 is armor Use text e.g. type:
armor

License, can't ever make any money The opportunity to make a little
money

The only way for me to get what I wanted was to write my own mud server
from scratch. Now in the world of mudding this is definitely not a new
idea. The mudding world is strewn with corpses of many mud servers that
started out with same goal as mine, to create something different. But
nonetheless I continued to pursue my desire for something different.

# OMugs

There was just one thing stopping me from writing my own mud server, I
had no idea how to code a telnet server application. A mud server must
be able to open and close internet connections as players logon and
logoff. So I was stuck with an idea and no way to implement it.

One day as I was surfing the net in my quest for more knowledge about
muds and mud servers, I found a site that offered a subscription to
mudding magazine and I promptly subscribed. In the first issue was an
article on creating you own mud from scratch accompanied by a code
listing in C++. This was not by any means a complete mud, just the part
required to open and close internet connections, allowing players to log
on and off.

Eureka! This was exactly what I needed to started my own mud server
codebase. I studied the code and purchased a book, Windows Sockets
Network Programming by Bob Quinn and Dave Shute, so I could decipher the
WinSock commands. Having obtained the knowledge necessary, I began work
on my own mud server codebase.

The creators of most mud server codebases have given their creation a
name. Some examples are: CircleMUD, DikuMud, Copper, Merc, and ROM. I
needed a name. One day on a trip to the Columbia SC zoo my wife, Dawn,
and I brainstormed naming ideas for the new mud server. We made a list
of words that might be used to make an acronym. Some of the words on the
list where: online, virtual, world, adventure, game, server, text,
text-based, internet, multi-user. We came up with many acronyms. I
decided on OtbMugs, Online Text-Based Multi-User Game Server and the
codebase existed for several months as OtbMugs. But I came to like the
name less and less and decided to drop the 'tb' (Text-Based) part
leaving OMugs, Online Multi-User Game Server, a name that Dawn had said
she liked on the day we were brainstorming names.

## Definition

OMugs is a combat oriented mud server and one of the goals of the game
is to gain levels. A player begins at level 1 and after some time
playing the game, the player will advance to level 2. This is called
leveling and is rewarded by giving the player a better chance of
defeating an opponent that was previously too difficult.

Socializing is an important aspect of any mud. If the game consisted of
only fighting, it could quickly become boring for many players. The
ability to 'talk' to other players is very important to the success of a
mud. OMugs supports talking to people in-room where everyone in the room
can 'hear' what each other is saying, as well as private messages, and
messages that everyone playing the game can hear.

Combat and socializing lead another very important factor, social
standing or status. A player's social standing on a mud is important.
Just as a BMW can be a status symbol, a rare and mighty sword can be a
status symbol on a mud. Levels can be viewed as ranks in the military.
Although there is no chain of command, a level 50 player knows more
about the game than the new player starting out at level 1.

## Combat Support

The following example of a play session serves to explain the OMugs
functionalities required to support combat.

A new player logs onto the game as Zeke. After looking around, he finds
that he is in a city with a sword in his hand and that he has no armor
and no money. A rat wanders by and Zeke decides to whack it. The rat is
vanquished and Zeke gets a rat ear and tail. After exploring the city
for a while he discovers a bake shop and sells the rat ear and tail to
the baker for 3 silver pieces each. Zeke heads for the Armor shop he
discovered while exploring the city and buys a leather tunic from
armorer for 5 silver pieces. Zeke continues to reduce the rat population
and soon he sees a message saying that he has gained a level. Now Zeke
is level 2, has upgraded his armor, and decides to venture outside the
city walls. Upon exiting the city gates, he is jumped by thief and
fights for this life. He is victorious and gets 10 gold pieces from the
thief and a pair of leather gloves. He puts the leather gloves on and
continues his adventures. This cycle of fighting, getting money, selling
loot, and upgrading his equipment, and leveling will continue throughout
Zeke's game life.

## Socializing

Players enjoy talking and playing together and OMugs provides several
channels of communication. Forms of communication include talking to
other players in the same room, or in your group, or to all players on
the mud, or a private message to a single player. There are also
non-verbal communications called socials and emotes. Socials and emotes
are only seen by players in the same room. Some socials are: smile,
wave, bow, and cheer. An emote (emotion) allows a player to type a
message and have their name placed before the message. If Zeke typed
'emote scratches his head.', the other players in the same room would
see a message saying 'Zeke scratches his head.'

# Conclusion

OMugs, as of this writing, has 30 player commands, 7 administrator
commands, and 9,534 lines of C++ code. There are many commands and
features yet to be implemented. This is just the beginning…..

Stephen L. Bryant

December 5, 2002
