package syntax

import (
	"batedit/token"
	"unicode"

	"github.com/nsf/termbox-go"
)

func Tokenize(rLine []rune) []token.Token {
	var tokens []token.Token
	var tok token.Token

	for i := 0; i < len(rLine); {
		ch := rLine[i]
		switch ch {
		case ' ':
			tok = token.Token{Content: "•", Color: termbox.ColorDarkGray}
			tokens = append(tokens, tok)
		case '/':
			if len(rLine) > 1 && i < len(rLine)-1 && rLine[i+1] == '/' {
				line := rLine[i:]
				tok = token.Token{Content: string(line), Color: termbox.ColorGreen}
				tokens = append(tokens, tok)
				i += len(line)
			} else {
				tok = token.Token{Content: string(ch), Color: termbox.ColorWhite}
				tokens = append(tokens, tok)
			}
		case '\t':
			tok = token.Token{Content: "→", Color: termbox.ColorDarkGray}
			tokens = append(tokens, tok)
		case '(', ')', '[', ']', '{', '}':
			tok = token.Token{Content: string(ch), Color: termbox.ColorLightYellow}
			tokens = append(tokens, tok)
		case '"':
			text := string(ch)
			if len(rLine) <= 1 || i == len(rLine)-1 {
				tok = token.Token{Content: text, Color: termbox.ColorYellow}
				tokens = append(tokens, tok)
			} else {
				oldPos := i
				i++
				ch = rLine[i]
				for ch != '"' {
					if i == len(rLine)-1 {
						break
					}
					i++
					ch = rLine[i]
				}
				text = string(rLine[oldPos : i+1])
				tok = token.Token{Content: text, Color: termbox.ColorYellow}
				tokens = append(tokens, tok)
			}
		case '\'':
			text := string(ch)
			if len(rLine) <= 1 || i == len(rLine)-1 {
				tok = token.Token{Content: text, Color: termbox.ColorYellow}
				tokens = append(tokens, tok)
			} else {
				oldPos := i
				i++
				ch = rLine[i]
				for ch != '\'' {
					if i == len(rLine)-1 {
						break
					}
					i++
					ch = rLine[i]
				}
				text = string(rLine[oldPos : i+1])
				tok = token.Token{Content: text, Color: termbox.ColorYellow}
				tokens = append(tokens, tok)
			}
		case '`':
			text := string(ch)
			if len(rLine) <= 1 || i == len(rLine)-1 {
				tok = token.Token{Content: text, Color: termbox.ColorYellow}
				tokens = append(tokens, tok)
			} else {
				oldPos := i
				i++
				ch = rLine[i]
				for ch != '`' {
					if i == len(rLine)-1 {
						break
					}
					i++
					ch = rLine[i]
				}
				text = string(rLine[oldPos : i+1])
				tok = token.Token{Content: text, Color: termbox.ColorYellow}
				tokens = append(tokens, tok)
			}
		default:
			if unicode.IsLetter(ch) {
				text := ""
				for unicode.IsLetter(ch) {
					text += string(ch)

					if i == len(rLine)-1 || !unicode.IsLetter(rLine[i+1]) {
						break
					}
					//fmt.Println(text)
					i++
					ch = rLine[i]
				}
				if token.IsGoKeyword(text) {
					color := token.GoKeyWords[text]
					tok = token.Token{Content: text, Color: color}
				} else {
					tok = token.Token{Content: text, Color: termbox.ColorWhite}
				}
				tokens = append(tokens, tok)
			} else {
				tok = token.Token{Content: string(ch), Color: termbox.ColorWhite}
				tokens = append(tokens, tok)
			}
		}
		i += 1
	}
	return tokens
}
