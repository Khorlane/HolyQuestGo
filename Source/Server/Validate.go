package server

import (
	"bufio"
	"os"
)

// LogValErr logs a validation error message with a file name.
func LogValErr(Message string, FileName string) {
	if StrGetLength(Message) > 50 {
		// Message is too long, chop it off
		Message = StrLeft(Message, 50)
	}
	Message = " - " + Message + " "
	for StrGetLength(Message) < 55 {
		Message += "-"
	}
	Message += "> "
	LogBuf = "ValErr"
	LogBuf += Message
	LogBuf += FileName
	LogIt(LogBuf)
	ValErr = true
}

// ValidateIt performs validation based on the specified type.
func ValidateIt(ValidationType string) bool {
	ValErr = false
	ValidationType = StrMakeLower(ValidationType)

	switch ValidationType {
		case "all":
			ValidateAll()
		case "mobiles":
			ValidateLibraryMobiles()
			ValidateLibraryWorldMobiles()
		case "objects":
			ValidateLibraryObjects()
			ValidateLibraryLoot()
			ValidateLibraryShops()
		case "rooms":
			ValidateLibraryRooms()
			ValidateLibraryShops()
			ValidateLibraryWorldMobiles()
	}

	if ValErr {
		LogBuf = "ValErr - Validation failed!!"
		LogIt(LogBuf)
	} else {
		LogBuf = "Validation successful!!"
		LogIt(LogBuf)
	}

	return ValErr
}

// ValidateAll performs all validation checks.
func ValidateAll() {
	ValidateLibraryLoot()
	ValidateLibraryMobiles()
	ValidateLibraryObjects()
	ValidateLibraryRooms()
	ValidateLibraryShops()
	ValidateLibraryWorldMobiles()
	ValidateRunningPlayers()
	ValidateRunningPlayersPlayerEqu()
	ValidateRunningPlayersPlayerObj()
	ValidateRunningRoomMob()
	ValidateRunningRoomObj()
}

// ValidateLibraryLoot validates the library loot.
func ValidateLibraryLoot() {
	var FileName     string
	var LineCount    int
	var Message      string
	var ObjectId     string
	var LootFileName string

	LogBuf = "Begin validation LibraryLoot"
	LogIt(LogBuf)

	if err := ChgDir(LOOT_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryLoot - Change directory to LOOT_DIR failed")
		return
	}

	// Get list of all LibraryLoot files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryLoot - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateLibraryLoot - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open player file
		LootFileName = entry.Name()
		LootFileName = LOOT_DIR + LootFileName
		LootFile, err := os.Open(LootFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateLibraryLoot - Open loot file failed")
			return
		}
		defer LootFile.Close()

		LineCount = 0
		Scanner := bufio.NewScanner(LootFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "" {
				break
			}

			LineCount++
			ObjectId = StrGetWord(Stuff, 3)

			// ObjectId
			ObjectIdFileName := OBJECTS_DIR + ObjectId + ".txt"
			if !FileExist(ObjectIdFileName) {
				// ObjectId file not found
				Message = "Object file '" + ObjectId + "' not found"
				FileName = LootFileName
				LogValErr(Message, FileName)
			}
		}
	}

	LogBuf = "Done validating LibraryLoot"
	LogIt(LogBuf)
}

// ValidateLibraryMobiles validates the library mobiles.
func ValidateLibraryMobiles() {
	var FieldName      string
	var FieldValue     string
	var FileName       string
	var LineCount      int
	var Message        string
	var MobileFileName string
	var MobileId       string

	LogBuf = "Begin validation LibraryMobiles"
	LogIt(LogBuf)

	if err := ChgDir(MOBILES_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryMobiles - Change directory to MOBILES_DIR failed")
		return
	}

	// Get list of all LibraryMobiles files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryMobiles - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateLibraryMobiles - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open mobile file
		MobileFileName = entry.Name()
		MobileId = StrLeft(MobileFileName, StrGetLength(MobileFileName)-4)
		MobileFileName = MOBILES_DIR + MobileFileName
		MobileFile, err := os.Open(MobileFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateLibraryMobiles - Open mobile file failed")
			return
		}
		defer MobileFile.Close()

		LineCount = 0
		Scanner := bufio.NewScanner(MobileFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "" {
				break
			}

			LineCount++
			FieldName = StrGetWord(Stuff, 1)
			FieldValue = StrGetWord(Stuff, 2)

			// MobileId field must be first
			if LineCount == 1 && FieldName != "MobileId:" {
				Message = "MobileId field must be the first field"
				FileName = MobileFileName
				LogValErr(Message, FileName)
			}

			// MobileId validation
			if FieldName == "MobileId:" && MobileId != FieldValue {
				Message = "MobileId must match file name"
				FileName = MobileFileName
				LogValErr(Message, FileName)
			}
		}
	}

	LogBuf = "Done validating LibraryMobiles"
	LogIt(LogBuf)
}

// ValidateLibraryObjects validates the library objects.
func ValidateLibraryObjects() {
	var FieldName      string
	var FieldValue     string
	var FileName       string
	var LineCount      int
	var Message        string
	var ObjectFileName string
	var ObjectId       string

	LogBuf = "Begin validation LibraryObjects"
	LogIt(LogBuf)

	if err := ChgDir(OBJECTS_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryObjects - Change directory to OBJECTS_DIR failed")
		return
	}

	// Get list of all LibraryObjects files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryObjects - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateLibraryObjects - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open object file
		ObjectFileName = entry.Name()
		ObjectId = StrLeft(ObjectFileName, StrGetLength(ObjectFileName)-4)
		ObjectFileName = OBJECTS_DIR + ObjectFileName
		ObjectFile, err := os.Open(ObjectFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateLibraryObjects - Open object file failed")
			return
		}
		defer ObjectFile.Close()

		LineCount = 0
		Scanner := bufio.NewScanner(ObjectFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "" {
				break
			}

			LineCount++
			FieldName = StrGetWord(Stuff, 1)
			FieldValue = StrGetWord(Stuff, 2)

			// ObjectId field must be first
			if LineCount == 1 && FieldName != "ObjectId:" {
				Message = "ObjectId field must be the first field"
				FileName = ObjectFileName
				LogValErr(Message, FileName)
			}

			// ObjectId validation
			if FieldName == "ObjectId:" && ObjectId != FieldValue {
				Message = "ObjectId must match file name"
				FileName = ObjectFileName
				LogValErr(Message, FileName)
			}
		}
	}

	LogBuf = "Done validating LibraryObjects"
	LogIt(LogBuf)
}

// ValidateLibraryRooms validates the library rooms.
func ValidateLibraryRooms() {
	var FieldName    string
	var FieldValue   string
	var FileName     string
	var LineCount    int
	var Message      string
	var RoomFileName string
	var RoomId       string

	LogBuf = "Begin validation LibraryRooms"
	LogIt(LogBuf)

	if err := ChgDir(ROOMS_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryRooms - Change directory to ROOMS_DIR failed")
		return
	}

	// Get list of all LibraryRooms files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryRooms - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateLibraryRooms - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open room file
		RoomFileName = entry.Name()
		RoomId = StrLeft(RoomFileName, StrGetLength(RoomFileName)-4)
		RoomFileName = ROOMS_DIR + RoomFileName
		RoomFile, err := os.Open(RoomFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateLibraryRooms - Open room file failed")
			return
		}
		defer RoomFile.Close()

		LineCount = 0
		Scanner := bufio.NewScanner(RoomFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "End of Room" {
				break
			}

			LineCount++
			FieldName = StrGetWord(Stuff, 1)
			FieldValue = StrGetWord(Stuff, 2)

			// RoomId field must be first
			if LineCount == 1 && FieldName != "RoomId:" {
				Message = "RoomId field must be the first field"
				FileName = RoomFileName
				LogValErr(Message, FileName)
			}

			// RoomId validation
			if FieldName == "RoomId:" && RoomId != FieldValue {
				Message = "RoomId must match file name"
				FileName = RoomFileName
				LogValErr(Message, FileName)
			}
		}
	}

	LogBuf = "Done validating LibraryRooms"
	LogIt(LogBuf)
}

// ValidateLibraryShops validates the library shops.
func ValidateLibraryShops() {
	var FieldName    string
	var FieldValue   string
	var FileName     string
	var LineCount    int
	var Message      string
	var ObjectId     string
	var ShopFileName string

	LogBuf = "Begin validation LibraryShops"
	LogIt(LogBuf)

	if err := ChgDir(SHOPS_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryShops - Change directory to SHOPS_DIR failed")
		return
	}

	// Get list of all LibraryShops files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryShops - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateLibraryShops - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open shop file
		ShopFileName = entry.Name()
		ShopFileName = SHOPS_DIR + ShopFileName
		ShopFile, err := os.Open(ShopFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateLibraryShops - Open shop file failed")
			return
		}
		defer ShopFile.Close()

		LineCount = 0
		Scanner := bufio.NewScanner(ShopFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "End of Shop" {
				break
			}

			LineCount++
			FieldName = StrGetWord(Stuff, 1)
			FieldValue = StrGetWord(Stuff, 2)

			if FieldName != "Item:" {
				continue
			}

			// Item validation
			ObjectId = FieldValue
			ObjectIdFileName := OBJECTS_DIR + ObjectId + ".txt"
			if !FileExist(ObjectIdFileName) {
				Message = "Object file '" + ObjectId + "' not found"
				FileName = ShopFileName
				LogValErr(Message, FileName)
			}
		}
	}

	LogBuf = "Done validating LibraryShops"
	LogIt(LogBuf)
}

// ValidateLibraryWorldMobiles validates the library world mobiles.
func ValidateLibraryWorldMobiles() {
	var FieldName           string
	var FieldValue          string
	var FileName            string
	var LineCount           int
	var Message             string
	var MobileId            string
	var MobileIdFileName    string
	var RoomIdFileName      string
	var WorldMobileFileName string
	var WorldMobileName     string

	LogBuf = "Begin validation LibraryWorldMobiles"
	LogIt(LogBuf)

	if err := ChgDir(WORLD_MOBILES_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryWorldMobiles - Change directory to WORLD_MOBILES_DIR failed")
		return
	}

	// Get list of all LibraryWorldMobiles files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateLibraryWorldMobiles - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateLibraryWorldMobiles - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open world mobile file
		WorldMobileFileName = entry.Name()
		if WorldMobileFileName == "ReadMe.txt" {
			// Skip ReadMe files
			continue
		}
		WorldMobileName = StrLeft(WorldMobileFileName, StrGetLength(WorldMobileFileName)-4)
		MobileId = WorldMobileName
		WorldMobileFileName = WORLD_MOBILES_DIR + WorldMobileFileName
		WorldMobileFile, err := os.Open(WorldMobileFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateLibraryWorldMobiles - Open world mobile file failed")
			return
		}
		defer WorldMobileFile.Close()

		// MobileId
		MobileIdFileName = MOBILES_DIR + MobileId + ".txt"
		if !FileExist(MobileIdFileName) {
			// MobileId file not found
			Message = "Mobile file '" + MobileId + "' not found"
			FileName = WorldMobileFileName
			LogValErr(Message, FileName)
		}

		// Check file contents
		LineCount = 0
		Scanner := bufio.NewScanner(WorldMobileFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "" {
				break
			}

			LineCount++
			FieldName = StrGetWord(Stuff, 1)
			FieldValue = StrGetWord(Stuff, 2)

			// RoomId
			if FieldName == "RoomId:" {
				RoomIdFileName = ROOMS_DIR + FieldValue + ".txt"
				if !FileExist(RoomIdFileName) {
					// RoomId file not found
					Message = "Room file '" + FieldValue + "' not found"
					FileName = WorldMobileFileName
					LogValErr(Message, FileName)
				}
			}
		}
	}

	LogBuf = "Done validating LibraryWorldMobiles"
	LogIt(LogBuf)
}

// ValidateRunningPlayers validates the running players.
func ValidateRunningPlayers() {
	var FieldName      string
	var FieldValue     string
	var FileName       string
	var LineCount      int
	var Message        string
	var PlayerFileName string
	var PlayerName     string
	var RoomIdFileName string

	LogBuf = "Begin validation RunningPlayers"
	LogIt(LogBuf)

	if err := ChgDir(PLAYER_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateRunningPlayers - Change directory to PLAYER_DIR failed")
		return
	}

	// Get list of all RunningPlayers files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateRunningPlayers - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateRunningPlayers - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open player file
		PlayerFileName = entry.Name()
		PlayerName = StrLeft(PlayerFileName, StrGetLength(PlayerFileName)-4)
		PlayerFileName = PLAYER_DIR + PlayerFileName
		PlayerFile, err := os.Open(PlayerFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateRunningPlayers - Open player file failed")
			return
		}
		defer PlayerFile.Close()

		LineCount = 0
		Scanner := bufio.NewScanner(PlayerFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "" {
				break
			}

			LineCount++
			StrReplace(&Stuff, ":", " ")
			FieldName = StrGetWord(Stuff, 1)
			FieldValue = StrGetWord(Stuff, 2)

			// Name field must be first
			if LineCount == 1 {
				if FieldName != "Name" {
					// Must be Name
					Message = "Name field must be the first field"
					FileName = PlayerFileName
					LogValErr(Message, FileName)
				}
			}

			// Name validation
			if FieldName == "Name" {
				if PlayerName != FieldValue {
					// Name must match file name
					Message = "Name must match file name"
					FileName = PlayerFileName
					LogValErr(Message, FileName)
				}
			}

			// RoomId validation
			if FieldName == "RoomId" {
				RoomIdFileName = ROOMS_DIR + FieldValue + ".txt"
				if !FileExist(RoomIdFileName) {
					// RoomId file not found
					Message = "Room file '" + FieldValue + "' not found"
					FileName = PlayerFileName
					LogValErr(Message, FileName)
				}
			}
		}
	}

	LogBuf = "Done validating RunningPlayers"
	LogIt(LogBuf)
}

// ValidateRunningPlayersPlayerEqu validates the equipment of running players.
func ValidateRunningPlayersPlayerEqu() {
	var FileName          string
	var LineCount         int
	var Message           string
	var ObjectId          string
	var ObjectIdFileName  string
	var PlayerEquFileName string
	var WearPosition      string

	LogBuf = "Begin validation RunningPlayersPlayerEqu"
	LogIt(LogBuf)

	if err := ChgDir(PLAYER_EQU_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateRunningPlayersPlayerEqu - Change directory to PLAYER_EQU_DIR failed")
		return
	}

	// Get list of all RunningPlayersPlayerEqu files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateRunningPlayersPlayerEqu - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateRunningPlayersPlayerEqu - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open player equipment file
		PlayerEquFileName = entry.Name()
		PlayerEquFileName = PLAYER_EQU_DIR + PlayerEquFileName
		PlayerEquFile, err := os.Open(PlayerEquFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateRunningPlayersPlayerEqu - Open player equipment file failed")
			return
		}
		defer PlayerEquFile.Close()

		LineCount = 0
		Scanner := bufio.NewScanner(PlayerEquFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "" {
				break
			}

			LineCount++
			WearPosition = StrGetWord(Stuff, 1)
			ObjectId = StrGetWord(Stuff, 2)

			// Wear position validation
			if WearPosition < "01" {
				// Wear position must be 01 - 20
				Message = "Wear position < 01"
				FileName = PlayerEquFileName
				LogValErr(Message, FileName)
			}
			if WearPosition > "20" {
				// Wear position must be 01 - 20
				Message = "Wear position > 20"
				FileName = PlayerEquFileName
				LogValErr(Message, FileName)
			}

			// ObjectId validation
			ObjectIdFileName = OBJECTS_DIR + ObjectId + ".txt"
			if !FileExist(ObjectIdFileName) {
				// ObjectId file not found
				Message = "Object file '" + ObjectId + "' not found"
				FileName = PlayerEquFileName
				LogValErr(Message, FileName)
			}
		}
	}

	LogBuf = "Done validating RunningPlayersPlayerEqu"
	LogIt(LogBuf)
}

// ValidateRunningPlayersPlayerObj validates the objects of running players.
func ValidateRunningPlayersPlayerObj() {
	var FileName          string
	var LineCount         int
	var Message           string
	var ObjectId          string
	var ObjectIdFileName  string
	var PlayerObjFileName string

	LogBuf = "Begin validation RunningPlayersPlayerObj"
	LogIt(LogBuf)

	if err := ChgDir(PLAYER_OBJ_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateRunningPlayersPlayerObj - Change directory to PLAYER_OBJ_DIR failed")
		return
	}

	// Get list of all RunningPlayersPlayerObj files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateRunningPlayersPlayerObj - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateRunningPlayersPlayerObj - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open player object file
		PlayerObjFileName = entry.Name()
		PlayerObjFileName = PLAYER_OBJ_DIR + PlayerObjFileName
		PlayerObjFile, err := os.Open(PlayerObjFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateRunningPlayersPlayerObj - Open player object file failed")
			return
		}
		defer PlayerObjFile.Close()

		LineCount = 0
		Scanner := bufio.NewScanner(PlayerObjFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "" {
				break
			}

			LineCount++
			ObjectId = StrGetWord(Stuff, 2)

			// ObjectId validation
			ObjectIdFileName = OBJECTS_DIR + ObjectId + ".txt"
			if !FileExist(ObjectIdFileName) {
				// ObjectId file not found
				Message = "Object file '" + ObjectId + "' not found"
				FileName = PlayerObjFileName
				LogValErr(Message, FileName)
			}
		}
	}

	LogBuf = "Done validating RunningPlayersPlayerObj"
	LogIt(LogBuf)
}

// ValidateRunningRoomMob validates the mobiles in running rooms.
func ValidateRunningRoomMob() {
	var FileName         string
	var LineCount        int
	var Message          string
	var MobileId         string
	var MobileIdFileName string
	var PositionOfDot    int
	var RoomId           string
	var RoomIdFileName   string
	var RoomMobFileName  string

	LogBuf = "Begin validation RunningRoomMob"
	LogIt(LogBuf)

	if err := ChgDir(ROOM_MOB_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateRunningRoomMob - Change directory to ROOM_MOB_DIR failed")
		return
	}

	// Get list of all RunningRoomMob files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateRunningRoomMob - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateRunningRoomMob - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open room mob file
		RoomMobFileName = entry.Name()
		if RoomMobFileName == "ReadMe.txt" {
			// Skip ReadMe files
			continue
		}
		RoomId = StrLeft(RoomMobFileName, StrGetLength(RoomMobFileName)-4)
		RoomMobFileName = ROOM_MOB_DIR + RoomMobFileName
		RoomMobFile, err := os.Open(RoomMobFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateRunningRoomMob - Open room mob file failed")
			return
		}
		defer RoomMobFile.Close()

		// RoomId validation
		RoomIdFileName = ROOMS_DIR + RoomId + ".txt"
		if !FileExist(RoomIdFileName) {
			// RoomId file not found
			Message = "Room file '" + RoomId + "' not found"
			FileName = RoomMobFileName
			LogValErr(Message, FileName)
		}

		// Check file contents
		LineCount = 0
		Scanner := bufio.NewScanner(RoomMobFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "" {
				break
			}

			LineCount++
			MobileId = StrGetWord(Stuff, 2)
			PositionOfDot = StrFindFirstChar(MobileId, '.')
			if PositionOfDot > 1 {
				// Mobile is hurt
				MobileId = StrLeft(MobileId, PositionOfDot)
			}

			// MobileId validation
			MobileIdFileName = MOBILES_DIR + MobileId + ".txt"
			if !FileExist(MobileIdFileName) {
				// MobileId file not found
				Message = "Mobile file '" + MobileId + "' not found"
				FileName = RoomMobFileName
				LogValErr(Message, FileName)
			}
		}
	}

	LogBuf = "Done validating RunningRoomMob"
	LogIt(LogBuf)
}

// ValidateRunningRoomObj validates the objects in running rooms.
func ValidateRunningRoomObj() {
	var FileName         string
	var LineCount        int
	var Message          string
	var ObjectId         string
	var ObjectIdFileName string
	var RoomId           string
	var RoomIdFileName   string
	var RoomObjFileName  string

	LogBuf = "Begin validation RunningRoomObj"
	LogIt(LogBuf)

	if err := ChgDir(ROOM_OBJ_DIR); err != nil {
		// Change directory failed
		LogIt("ValidateRunningRoomObj - Change directory to ROOM_OBJ_DIR failed")
		return
	}

	// Get list of all RunningRoomObj files
	if err := ChgDir(HomeDir); err != nil {
		// Change directory failed
		LogIt("ValidateRunningRoomObj - Change directory to HomeDir failed")
		return
	}

	entries, err := os.ReadDir("./")
	if err != nil {
		LogIt("ValidateRunningRoomObj - ReadDir failed")
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Open RunningRoomObj file
		RoomObjFileName = entry.Name()
		if RoomObjFileName == "ReadMe.txt" {
			// Skip ReadMe files
			continue
		}
		RoomId = StrLeft(RoomObjFileName, StrGetLength(RoomObjFileName)-4)
		RoomObjFileName = ROOM_OBJ_DIR + RoomObjFileName
		RoomObjFile, err := os.Open(RoomObjFileName)
		if err != nil {
			// File does not exist - Very bad!
			LogIt("ValidateRunningRoomObj - Open room object file failed")
			return
		}
		defer RoomObjFile.Close()

		// RoomId validation
		RoomIdFileName = ROOMS_DIR + RoomId + ".txt"
		if !FileExist(RoomIdFileName) {
			// RoomId file not found
			Message = "Room file '" + RoomId + "' not found"
			FileName = RoomObjFileName
			LogValErr(Message, FileName)
		}

		// Check file contents
		LineCount = 0
		Scanner := bufio.NewScanner(RoomObjFile)
		for Scanner.Scan() {
			Stuff := Scanner.Text()
			if Stuff == "" {
				break
			}

			LineCount++
			ObjectId = StrGetWord(Stuff, 2)

			// ObjectId validation
			ObjectIdFileName = OBJECTS_DIR + ObjectId + ".txt"
			if !FileExist(ObjectIdFileName) {
				// ObjectId file not found
				Message = "Object file '" + ObjectId + "' not found"
				FileName = RoomObjFileName
				LogValErr(Message, FileName)
			}
		}
	}

	LogBuf = "Done validating RunningRoomObj"
	LogIt(LogBuf)
}