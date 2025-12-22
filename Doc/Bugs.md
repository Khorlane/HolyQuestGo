# OMugs – Online Multi-User Game Server
# Bug Report

# License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

For more information, please refer to <http://unlicense.org/>

# Bugs

<table style="width:100%;">
<colgroup>
<col style="width: 14%" />
<col style="width: 28%" />
<col style="width: 43%" />
<col style="width: 13%" />
</colgroup>
<tbody>
<tr>
<td><p><strong>Date</strong></p>
<p><strong>Reported</strong></p></td>
<td><strong>Bug</strong></td>
<td style="text-align: left;">Steps to recreate</td>
<td><p><strong>Date</strong></p>
<p><strong>Fixed</strong></p></td>
</tr>
<tr>
<td>03/18/2004</td>
<td style="text-align: left;">When you delete a character, OMugs
abends</td>
<td style="text-align: left;"><ol type="1">
<li><p>Delete a character</p></li>
</ol></td>
<td>03/18/2004</td>
</tr>
<tr>
<td><p>11/05/2003</p>
<p>Alex</p></td>
<td style="text-align: left;">When buying 'an axe' from a shop and 'a
blunt axe' is listed before 'an axe', you cannot buy the axe. You always
get the blunt axe.</td>
<td style="text-align: left;"><ol start="2" type="1">
<li><p>Go to weapon shop</p></li>
<li><p>issue command 'buy axe'</p></li>
<li><p>You will get a blunt axe</p></li>
</ol></td>
<td>11/07/2003</td>
</tr>
<tr>
<td><p>03/10/2003</p>
<p>slb</p></td>
<td style="text-align: left;">Each time player logs on, the weapon type
returns to 'hand'</td>
<td style="text-align: left;"><ol start="5" type="1">
<li><p>Logon, wield sword, logoff</p></li>
<li><p>Logon, kill rat</p></li>
<li><p>Message says 'slap' rat instead of 'slash' rat</p></li>
</ol></td>
<td>03/25/2003</td>
</tr>
<tr>
<td><p>03/24/2003</p>
<p>slb</p></td>
<td style="text-align: left;">Numerous problems with wear command</td>
<td style="text-align: left;"><ol start="8" type="1">
<li><p>Wear helmet</p></li>
<li><p>Wear helmet</p></li>
<li><p>Eq check shows that player is wearing two helmets</p></li>
</ol></td>
<td><p>03/24/2003</p>
<p>slb</p></td>
</tr>
<tr>
<td><p>103/17/2003</p>
<p>slb</p></td>
<td style="text-align: left;">'Played' command reports an extremely high
amount of time played, when a player has just been created.</td>
<td style="text-align: left;"><p>1) Create a new player</p>
<p>2) Type 'played'</p>
<p>3) You've played shows 'days' of time played. This is obviously
wrong.</p></td>
<td><p>03/18/2003</p>
<p>slb</p></td>
</tr>
<tr>
<td><p>03/04/2003</p>
<p>slb</p></td>
<td style="text-align: left;">Assist is always 'off' when a player logs
on, even though they set it to 'on' last time they logged on.</td>
<td style="text-align: left;"><p>1) Logon Steve</p>
<p>2) Type 'assist on'</p>
<p>3) Type 'status' and check 'assist'</p>
<p>4) Type 'quit'</p>
<p>5) Logon Steve</p>
<p>6) Type 'status' and check 'assist'</p>
<p>7) Assist is 'off' again, should be 'on'</p></td>
<td><p>03/04/2003</p>
<p>slb</p></td>
</tr>
<tr>
<td><p>12/13/2002</p>
<p>slb</p></td>
<td style="text-align: left;">Social (like bow) is not replace $P with
player’s name or $T with target player’s name</td>
<td style="text-align: left;"><p>1) Logon on two players</p>
<p>2) Player 1 types ‘bow’</p>
<p>3) Player 2 sees “$P bows respectfully.”</p></td>
<td><p>12/172002</p>
<p>slb</p></td>
</tr>
<tr>
<td><p>12/04/2002</p>
<p>slb</p></td>
<td style="text-align: left;">Sometimes when the first player logs on,
two connections are created.</td>
<td style="text-align: left;"><p>1) Start Omugs</p>
<p>2) Logon Steve</p></td>
<td></td>
</tr>
<tr>
<td><p>04/04/2002</p>
<p>slb</p></td>
<td style="text-align: left;">Same player shows twice on 'who' list.
Same player can logon more than once. There can be more than one Steve
logged on at the same time.</td>
<td style="text-align: left;"><p>1) Logon Steve</p>
<p>2) Connect using Chris, but say new player = 'y'</p>
<p>3) Enter Chris for name</p>
<p>4) Receive name already used message</p>
<p>5) Disconnect</p>
<p>6) Logon Chris normal</p></td>
<td><p>04/09/2002</p>
<p>slb</p></td>
</tr>
</tbody>
</table>
