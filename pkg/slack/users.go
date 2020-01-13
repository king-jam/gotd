package slack

// validateaUser will validate userID against the userIDList
func (h slashCommandHandler) validateUser(userID string) bool {
	var userIDList = []string{
		// Hopkinton
		"WR3R1N810", // Ethan K
		"WRGJMAYGG", // Val C
		"WRG57T0QL", // Aman W
		"WR3R1UBT4", // Amy M
		"WR953FV3K", // Geoff R
		"WRG56TZCL", // Joe G
		"WRG56U7RS", // Justin K
		"WRGJMR8LQ", // Kranti U
		"WRJEJEEF9", // Minh N
		"WR3R2QZM0", // Viet D
		"WRG570H1A", // Dale B
		"WR8LBQQM7", // Edgardo R
		"WR8LAPMKK", // James K
		"WR8LBN249", // Sammie G
		"WR8LBPNV7", // Nicole R
		"WR54WA943", // Erin B
		"WRJQTFTFY", // Jim C
		"WRGJYNQ7P", // Jeff E
		"WRG5816CU", // Charles W
		"WR54X6ZKM", // Calvin C
		"WLYG93DAR", // Brett B
		// Cambridge
		"WRG57PEN8", // Amy S
		"WRG56LXEG", // Megan M
		"WLYC0PC2Y", // Thinh N
		"WM0MKG77H", // Xhe
		"WRBHB0JP9", // AK
	}

	for _, user := range userIDList {
		if userID == user {
			return true
		}
	}

	return false
}
