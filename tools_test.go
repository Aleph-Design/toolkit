package toolkit

import "testing"

func TestTools_RandomString(t *testing.T) {

		var testTools Tools

		str := testTools.RandomString(10)

		if len(str) != 10 {
			t.Error("Returned wrong string length!")
		}
}

/*
janhkila@imac app % cd .. 
janhkila@imac toolkit-project % cd toolkit	<== Make sure you're in 'toolkit' folder!
janhkila@imac toolkit % go test .
ok      github.com/Aleph-Design/toolkit 0.144s
janhkila@imac toolkit % 
*/
