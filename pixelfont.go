package pixelding

/*	+ - +
  	|   |
	+ - +
*/
const (
	SingleFrame = "\u250C\u2500\u2510\u2502*\u2502\u2514\u2500\u2518"
	DoubleFrame = "\u2554\u2550\u2557\u2551*\u2551\u255A\u2550\u255D"
	RoundFrame  = "\u256D\u2500\u256E\u2502*\u2502\u2570\u2500\u256F"
	BlockFrame  = "\u259B\u2580\u259C\u258C*\u2590\u2599\u2584\u259F"
	TextFrame   = "+-+|*|+-+"
)
const HBar = "\u2588\u258F\u258E\u258D\u258C\u258B\u258A\u2589"

const (
	Dot1x1Pattern   = 0B01010101  // - - - - - - - - - -
	Dot2x2Pattern   = 0B00110011  // --  --  --  --  --
	Dot4x4Pattern   = 0B00001111  // ----    ----    ----
	Dot1x3Pattern   = 0B00010001  // -   -   -   -   -   -
	Dot3x2x1Pattern = 0B00100111  // ---  -  ---  -  ---  -
	Dot6x2Pattern   = 0B00111111  // ------  ------  ------
	Dot7x1Pattern   = 0B01111111  // ------- ------- -------
	Dot5x1x1Pattern = 0B01011111  // ----- - ----- - ----- -
)

/*
const SingleFrame = [9]string{
	string(0x250C), string(0x2500), string(0x2510),
	string(0x2502), string(32), string(0x2502),
	string(0x2514), string(0x2500), string(0x2518)}
*/
//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) LoadStdStamp() *PixelStamp {
	StdStamp := PixelStamp{
		false, 0,
		[]uint64{
			0B00111111111111111111111111110011111111111111111111111100,
			0B01000000000000000000000000000111111111111111111111111110,
			0B10011111000000000000000000011000001100100111011000001111,
			0B10011000100000000000000000011001110100100111010011110111,
			0B10011000101000000000000010011001110100100011010011111111,
			0B10011000100000000000000010011001110100100101010011111111,
			0B10011111001010001001110010011001110100100110010011000111,
			0B10011000001001010010001010011001110100100111010011110111,
			0B10011000001000100011111010011001110100100111010011110111,
			0B10011000001001010010000010011001110100100111010011110111,
			0B10011000001010001001110011011000001100100111011000001111,
			0B01000000000000000000000000001111111111111111111111111110,
			0B00111111111111111111111111100111111111111111111111111100},
	}
	return &StdStamp
}

//----------------------------------------------------------------------------------------------------------------------
func (p *PixelDING) LoadStdFont() *PixelFont {
	StdFont := PixelFont{
		false,
		map[int]PixelChar{
			32: {0, 0, 3, 0, 0, 0, 0, []uint64{0B000, 0B000, 0B000, 0B000, 0B000}},
			46: {0, 0, 3, 0, 0, 0, 0, []uint64{0B000, 0B000, 0B000, 0B000, 0B010}},
			44: {0, 0, 0, 0, 0, 0, 0, []uint64{0B000, 0B000, 0B000, 0B010, 0B100}},
			33: {0, 0, 0, 0, 0, 0, 0, []uint64{0B010, 0B010, 0B010, 0B000, 0B010}},
			40: {0, 0, 0, 0, 0, 0, 0, []uint64{0B001, 0B010, 0B010, 0B010, 0B001}},
			41: {0, 0, 0, 0, 0, 0, 0, []uint64{0B010, 0B001, 0B001, 0B001, 0B010}},
			91: {0, 0, 0, 0, 0, 0, 0, []uint64{0B011, 0B010, 0B010, 0B010, 0B011}},
			93: {0, 0, 0, 0, 0, 0, 0, []uint64{0B011, 0B001, 0B001, 0B001, 0B011}},
			42: {0, 0, 0, 0, 0, 0, 0, []uint64{0B00000, 0B00100, 0B11111, 0B01010, 0B00000}},
			43: {0, 0, 0, 0, 0, 0, 0, []uint64{0B000, 0B010, 0B111, 0B010, 0B000}},
			45: {0, 0, 0, 0, 0, 0, 0, []uint64{0B000, 0B000, 0B111, 0B000, 0B000}},
			47: {0, 0, 0, 0, 0, 0, 0, []uint64{0B001, 0B010, 0B010, 0B100, 0B100}},
			92: {0, 0, 0, 0, 0, 0, 0, []uint64{0B100, 0B010, 0B010, 0B001, 0B001}},
			61: {0, 0, 0, 0, 0, 0, 0, []uint64{0B000, 0B111, 0B000, 0B111, 0B000}},
			65: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01110, 0B10001, 0B11111, 0B10001, 0B10001}},
			66: {0, 0, 0, 0, 0, 0, 0, []uint64{0B11110, 0B10001, 0B11110, 0B10001, 0B11110}},
			67: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01110, 0B10001, 0B10000, 0B10001, 0B01110}},
			68: {0, 0, 0, 0, 0, 0, 0, []uint64{0B11110, 0B10001, 0B10001, 0B10001, 0B11110}},
			69: {0, 0, 0, 0, 0, 0, 0, []uint64{0B1111, 0B1000, 0B1110, 0B1000, 0B1111}},
			70: {0, 0, 0, 0, 0, 0, 0, []uint64{0B1111, 0B1000, 0B1110, 0B1000, 0B1000}},
			71: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01110, 0B10000, 0B10111, 0B10001, 0B01110}},
			72: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10001, 0B10001, 0B11111, 0B10001, 0B10001}},
			73: {0, 0, 0, 0, 0, 0, 0, []uint64{0B111, 0B010, 0B010, 0B010, 0B111}},
			74: {0, 0, 0, 0, 0, 0, 2, []uint64{0B0001, 0B0001, 0B0001, 0B1001, 0B0110}},
			75: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10001, 0B11110, 0B10100, 0B10010, 0B10001}},
			76: {0, 0, 0, 0, 0, 1, 0, []uint64{0B1000, 0B1000, 0B1000, 0B1000, 0B1111}},
			77: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10001, 0B11011, 0B10101, 0B10001, 0B10001}},
			78: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10001, 0B11001, 0B10101, 0B10011, 0B10001}},
			79: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01110, 0B10001, 0B10001, 0B10001, 0B01110}},
			80: {0, 0, 0, 0, 0, 0, 0, []uint64{0B11110, 0B10001, 0B11110, 0B10000, 0B10000}},
			81: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01110, 0B10001, 0B10001, 0B10010, 0B01101}},
			82: {0, 0, 0, 0, 0, 0, 0, []uint64{0B11110, 0B10001, 0B11110, 0B10010, 0B10001}},
			83: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01111, 0B10000, 0B01110, 0B00001, 0B11110}},
			84: {0, 0, 0, 0, 0, 2, 1, []uint64{0B11111, 0B00100, 0B00100, 0B00100, 0B00100}},
			85: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10001, 0B10001, 0B10001, 0B10001, 0B01110}},
			86: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10001, 0B10001, 0B10001, 0B01010, 0B00100}},
			87: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10001, 0B10001, 0B10101, 0B11011, 0B10001}},
			88: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10001, 0B01010, 0B00100, 0B01010, 0B10001}},
			89: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10001, 0B01010, 0B00100, 0B00100, 0B00100}},
			90: {0, 0, 0, 0, 0, 0, 0, []uint64{0B11111, 0B00010, 0B00100, 0B01000, 0B11111}},
			48: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01110, 0B10001, 0B10101, 0B10001, 0B01110}},
			49: {0, 0, 0, 0, 0, 0, 0, []uint64{0B010, 0B110, 0B010, 0B010, 0B111}},
			50: {0, 0, 0, 0, 0, 0, 0, []uint64{0B11110, 0B00001, 0B01110, 0B10000, 0B11111}},
			51: {0, 0, 0, 0, 0, 0, 0, []uint64{0B11110, 0B00001, 0B01110, 0B00001, 0B11110}},
			52: {0, 0, 0, 0, 0, 0, 0, []uint64{0B10010, 0B10010, 0B11111, 0B00010, 0B00010}},
			53: {0, 0, 0, 0, 0, 0, 0, []uint64{0B11111, 0B10000, 0B11110, 0B00001, 0B11110}},
			54: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01110, 0B10000, 0B11110, 0B10001, 0B01110}},
			55: {0, 0, 0, 0, 0, 0, 0, []uint64{0B11111, 0B00001, 0B00010, 0B00100, 0B01000}},
			56: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01110, 0B10001, 0B01110, 0B10001, 0B01110}},
			57: {0, 0, 0, 0, 0, 0, 0, []uint64{0B01110, 0B10001, 0B01111, 0B00001, 0B01110}},
		},
	}
	f := prepareFont(StdFont)
	return &f
}
