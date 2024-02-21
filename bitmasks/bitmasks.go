package bitmasks

func Set(b, mask uint8) uint8    { return b | mask }
func Clear(b, mask uint8) uint8  { return b &^ mask }
func Toggle(b, mask uint8) uint8 { return b ^ mask }
func Has(b, mask uint8) bool     { return b&mask != 0 }
