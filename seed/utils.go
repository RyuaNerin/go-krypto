package seed

func getB0(A uint32) uint32 { return 0xFF & (A >> 0) }
func getB1(A uint32) uint32 { return 0xFF & (A >> 8) }
func getB2(A uint32) uint32 { return 0xFF & (A >> 16) }
func getB3(A uint32) uint32 { return 0xFF & (A >> 24) }
