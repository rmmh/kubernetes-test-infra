package editdistance

/*

func init() {
	rand.Seed(time.Now().UnixNano())
}

//var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var letterRunes = []rune("abc")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	n := 0
	for {
		aLen := rand.Intn(20) + 5
		bLen := rand.Intn(20) + 5
		a := RandStringRunes(aLen)
		b := RandStringRunes(bLen)
		brDist := editdistance.BerghelRoachDistance(a, b, 5)
		//lvDist := editdistance.LevenshteinDistance(a, b, 5)
		lvDist := brDist
		if brDist > 5 {
			if lvDist <= 5 {
				panic("FUCK")
			}
		} else if brDist != lvDist {
			panic("UG")
		}
		if brDist < 5 {
			n++
			if n&1023 == 0 {
				fmt.Println(a, b, brDist, lvDist)
			}
			if brDist != lvDist {
				panic("uguu")
			}
		}
	}
	//fmt.Println("WERUPU", editdistance.BerghelRoachDistance("foo", "football", 89))
}
*/
