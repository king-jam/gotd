package slack

// validateaUser will validate userID against the userIDList
func (h slashCommandHandler) validateUser(userID string) bool {
	var userIDList = []string{
		// Hopkinton
		"U5SFY08HW", // Ethan K
		"U5SFZ590Q", // Val C
		"UGG0Y2W82", // Aman W
		"U5UAGKX4L", // Amy M
		"U5U133V3Q", // Geoff R
		"U5U0X61DM", // Joe G
		"U5U1DSEQ7", // Justin K
		"U61HFJ7V2", // Kranti U
		"WRJEJEEF9", // Minh N
		"UFDAJLGJU", // Viet D
		"U5V5T2DPZ", // Dale B
		"UGYDW6UJK", // Edgardo R
		"U5T9HLMAN", // James K
		"UEK11RZJP", // Sammie G
		"UFQQU5S7N", // Nicole R
		"U5Y4FJ8JK", // Erin B
		"UM4E99TM4", // Jim C
		"UKUAY9URK", // Jeff E
		"UKWLQ33L3", // Charles W
		"UK4UMCQDV", // Calvin C
		"UM9GXDZ5E", // Brett B
		// Cambridge
		"UDDRDQ5LH", // Amy S
		"U5SFUMV32", // Megan M
		"U74S52CLT", // Thinh N
		"U6VF3D8AW", // Xhe
		"UBFRHV5GW", // AK
	}

	for _, user := range userIDList {
		if userID == user {
			return true
		}
	}

	return false
}
