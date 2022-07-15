package postalcodeSearch

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var kanaConv = unicode.SpecialCase{
	// ひらがなをカタカナに変換
	unicode.CaseRange{
		0x3041, // Lo: ぁ
		0x3093, // Hi: ん
		[unicode.MaxCase]rune{
			0x30a1 - 0x3041, // UpperCase でカタカナに変換
			0,               // LowerCase では変換しない
			0x30a1 - 0x3041, // TitleCase でカタカナに変換
		},
	},
	// カタカナをひらがなに変換
	unicode.CaseRange{
		0x30a1, // Lo: ァ
		0x30f3, // Hi: ン
		[unicode.MaxCase]rune{
			0,               // UpperCase では変換しない
			0x3041 - 0x30a1, // LowerCase でひらがなに変換
			0,               // TitleCase では変換しない
		},
	},
}

func HankakuKatakanaToKatakana(hankakuKatakana string) string {
	return norm.NFKC.String(hankakuKatakana)
}

func KatakanaToHiragana(katakana string) string {
	return strings.ToLowerSpecial(kanaConv, katakana)
}
