// Typical Jmap response will look like
// map[id:834331 e:LON t:VOD l:246.54 l_fix:246.54 ltt:3:08PM GMT lt:Feb 27, 3:08PM GMT l_cur:GBX246.54 s:0 c:+1.29 c_fix:1.29 cp:0.52 cp_fix:0.52 ccol:chg]


package gaccts

import "encoding/json"
import "errors"
import "fmt"
import "io/ioutil"
import "net/http"

//import "money"

type Jmap map[string]interface{}

type Fetched struct {
	index int
	body []byte
	err error
}

func mapifyResponseOrDie(gepic string, response []byte) Jmap {
        var s []byte
        s = response

        var s1 []byte
        s1 =  s[6:len(s)-3]

        var dat Jmap
        err := json.Unmarshal(s1, &dat)
        if  err != nil {
                fmt.Println("ERR: Failed for gepic:", gepic)
                panic(err)
        }
        return dat

}





func getUrl(messages chan Fetched, idx int, url string) {
	//var bytes []byte
	f := Fetched{idx, []byte{}, nil}
	f.index = idx
	response, err := http.Get(url)
	if err != nil { f.err = err ; return }
		//fmt.Printf("%s", err)
		//os.Exit(1)
	//} else {
	defer response.Body.Close()
	// fmt.Println("C")
	body, err := ioutil.ReadAll(response.Body)
	f.body = body
	f.err  = err
	//fmt.Println(body)
	messages <- f 
/*
	if err != nil { return bytes, err }
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	}
*/
}

func getUrlsMasked( urls []string, mask []bool) (responses [][]byte, errs []error, err error) {
	//bytes := make

	size := len(urls)

	responses = make([][]byte, size)
	errs = make([]error, size)
	if len(urls) != len(mask) {
		e := errors.New("Slice length of urls and mask are different")
		return responses, errs, e
	}

	numUrls := 0
	messages := make(chan Fetched)
	//fetches := make([]Fetched, size)
	for i, url := range urls {
		if mask[i] {
			go getUrl(messages, i, url)
			numUrls++
		}
	}

	//fmt.Println("AA")

	// reassmble in correct order
	for i:= 0; i < numUrls; i++ {
		//fmt.Println("B")
		//f := Fetched{}
		f := <- messages
		//fmt.Println("D")
		idx := f.index
		responses[idx] = f.body
		errs[idx] = f.err
	}

	return 
}

func getUrls(urls []string) (responses [][]byte, errs []error, err error) {
	mask := make([]bool, len(urls))
	for i, _ := range mask { mask[i] = true }
	return getUrlsMasked(urls, mask)
}


func TestGofi() {
	gepics := []string{"LON:ULVR", "LON:VOD"}
	CreateGofiReport(gepics)
}



func GepicsToUrls( gepics []string) (urls []string) {
	base := "http://finance.google.com/finance/info?client=ig&q="
	urls = make([]string, len(gepics))
	for i, gepic := range gepics { urls[i] = base + gepic }
	return
}


func GetGepicsMasked( gepics []string, mask []bool) (jmaps []Jmap) {
	urls := GepicsToUrls(gepics)

	bodies, errs, err := getUrlsMasked(urls, mask)
	if err != nil { panic(err) }

	jmaps = make([]Jmap, len(gepics))
	//for i, jn := range jsons  {
	for i:= 0 ; i< len(gepics); i++ {
		//fmt.Println(gepics[i])
		if mask[i] {
			if errs[i] != nil { panic(errs[i]) }
			jmaps[i] = mapifyResponseOrDie(gepics[i], bodies[i])
		}
	}
	return
}



func CreateGofiReport(gepics []string) {
	urls := GepicsToUrls(gepics)
	bodies, errs, err := getUrls(urls)
	if err != nil { panic(err) }
	// TODO should probably use a lot of GetGepics()

	//      ULVR 2435.00  -12.00   -0.49
	hdr := "EPIC   PRICE    CHG+    CHG%"
	fmt.Println(hdr)
	for i := 0; i < len(gepics) ; i++ {
		//fmt.Println(i)
		if errs[i] != nil {
			fmt.Println("Error with gepic", gepics[i])
			//panic(err[i]) TODO reinstate
		}

		jmap := mapifyResponseOrDie(gepics[i], bodies[i])
		//q := money.Quote{}
		//money.BytesToQuote(bodies[i], &q)
		t := jmap["t"]
		l_fix :=jmap["l_fix"]
		c_fix := jmap["c_fix"]
		cp_fix := jmap["cp_fix"]

		fmt.Printf("%4.4s %7.7s %7.7s %7.7s\n", t, l_fix, c_fix, cp_fix)

	}

}
		
	
