package toolkit

import "crypto/rand"

/*
This package in the toolkit folder will hold code that's
used when working on modules while developing.
*/
/*
Tools is a type to instanciate this module.
Any variable of this type will have access to all methods
trough receiver *Tools.
*/
type Tools struct {}

/*
Generate a string of random characters of length N.
===================================================
filenames < 100 characters.
forbidden characters are: <>:/\|?* and ASCI 1 - 31
@sourceStr
-	The source for generated random characters in the
	returned output string.
*/
const sourceStr = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890+&$%!"

func (t *Tools) RandomString(n int) string {
	if n < 100 {
		res, src := make([]rune, n), []rune(sourceStr)

		for i := range res {

			p, _ := rand.Prime(rand.Reader, len(src))

			x, y := p.Uint64(), uint64(len(src))

			res[i] = src[x % y]
		}

		return string(res)

	} else {
		return "N must be less than 100!"
	}
}