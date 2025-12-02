package unpackstring

import "errors"

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if s == "" {
		return "", nil
	}

	temp := []rune(s)
	var res []rune

	var prev rune
	hasPrev := false

	for i := 0; i < len(temp); {
		r := temp[i]

		switch {
		case r == '\\':
			// Экранируем следующую руну
			if i+1 >= len(temp) {
				return "", ErrInvalidString
			}
			next := temp[i+1]
			res = append(res, next)
			prev = next
			hasPrev = true
			i += 2

		case r >= '0' && r <= '9':
			// Цифра без экранирования — множитель повтора
			if !hasPrev {
				return "", ErrInvalidString
			}
			count := int(r - '0')
			if count == 0 {
				// Удаляем предыдущую руну
				if len(res) == 0 {
					return "", ErrInvalidString
				}
				res = res[:len(res)-1]
			} else {
				// Один prev уже есть, добавляем ещё count-1
				for j := 0; j < count-1; j++ {
					res = append(res, prev)
				}
			}
			i++

		default:
			// Обычный символ
			res = append(res, r)
			prev = r
			hasPrev = true
			i++
		}
	}

	return string(res), nil
}
