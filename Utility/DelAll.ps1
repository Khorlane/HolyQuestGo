# Delete all files generated during when running the HolyQuestGo server
# This is a cleanup script to remove all temporary and generated files
# so that the server can be started fresh without any leftover data.
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\MobPlayer\*'          -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\PlayerMob\*'          -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\MobStats\Armor\*'     -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\MobStats\Attack\*'    -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\MobStats\Damage\*'    -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\MobStats\Desc1\*'     -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\MobStats\ExpPoints\*' -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\MobStats\HitPoints\*' -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\MobStats\Loot\*'      -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Violence\MobStats\Room\*'      -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Control\RoomMobMove.txt'       -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Control\RoomMobMoveTemp.txt'   -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Control\RoomMobList.txt'       -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Control\RoomMobListTemp.txt'   -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Control\Events\*'              -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Control\Mobiles\InWorld\*'     -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Control\Mobiles\NoMove\*'      -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\Control\Mobiles\Spawn\*'       -Force -ErrorAction SilentlyContinue
Remove-Item 'C:\Projects\HolyQuestGo\Running\RoomMob\*'                     -Force -ErrorAction SilentlyContinue
pause