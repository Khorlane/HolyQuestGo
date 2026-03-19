# OMugs – Online Multi-User Game Server

# Feature List

# License

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any means.

For more information, please refer to <http://unlicense.org/>

# Things To Do

## Commands

- Light command

- Shout command

- Transfer command (admin)

## World

- Doors

## Connections

- Enable banning of ip addresses

- Add command to disconnect players

- Sometimes player ain't there but game thinks they are, connection is
  still there, no player on the other end. Kill the connection, be sure
  player file is set to Online:No

## Mobiles

- Make Function Mobile::CountMob count wounded mobs

- Make wimpy mobs flee when almost dead

## Players

- When a player dies, reduce move points to zero

- Allow players to create a description of themselves

- Player CRUD dialog

## Objects

- Containers

- Clean up 'pObject = new Object' code, what if \<object\>.txt doesn't
  exist for some reason. A pointer to object is still returned???

## Shops

- Add shop specific messages for buy / sell success / failure like
  BuyNotExist BuySuccess BuyNotAfford BuyNotAfford SellSuccess
  SellNotExist

## Misc

- Externalize starting hitpoints, movepoints, starting room, greeting,
  etc. In other words, most if not all config.h stuff

- Complete move points implemetation

- Design and write weather system

- Design and write quest system

- Send message to all players when the watch changes

- Mobile Armor is not implemented yet
GetMobileArmor()
... Old code
if err != nil {
  // Mobile Armor is not implemented, so for now, we just return zero
  MobileArmor = 0
  return MobileArmor
  // This code is currently unreachable, on purpose.
  LogIt("Violence::GetArmor - Open MobStatsArmorFile file failed (read)")
  os.Exit(1)
}
... New code to stop Go from complaining about unreachable code
if err != nil {
  // Mobile Armor is not implemented, so for now, we just return zero
  MobileArmor = 0
  return MobileArmor
}
