package main

import (
	//"html/template"
	//"io/ioutil"
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	//	"mime"
	//"reflect"
)

var keyStr = "8271997208960872478735181815578166723519929177896558845922250595511921395049126920528021164569045773"
var key = big.NewInt(0)
var sentimentMap = make(map[string]int)
var banMap = make(map[string]bool)

var db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/Tweets?charset=utf8mb4,utf8")
var pattern = regexp.MustCompile("[A-Za-z]+")

type Entry map[string]string
type Conversion struct {
	Text string
}

type transac struct {
	sync.Mutex
	tid string
	seq int
}

var transacMap = make(map[string]*transac)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Home</h1><div>Page</div>")
}

func phase1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>check</h1><div>phase1</div>")
	fmt.Println("answer")
}

func q1(w http.ResponseWriter, r *http.Request) {
	if checkParams(r, []string{"key", "message"}) == false {
		fmt.Fprintf(w, "lack of parameters")
		//message := "lack of parameters"
		return
	}
	params := parseForm(r)
	//fmt.Fprintf(w, "key = %s ; message = %s", params["key"], params["message"])
	decryptedM, err := decryptMessage(params["message"], params["key"], key)
	if err != nil {
		fmt.Fprintf(w, "error")
		//message := "wrong"
		return
	} else {
		t := time.Now()
		message := "MDK,1474-9673-1103\n" + t.Format("2006-01-02 15:04:05") + "\n" + decryptedM + "\n"
		fmt.Fprintf(w, message)
		return
	}

}

func checkParams(r *http.Request, params []string) bool {
	for _, param := range params {
		if trimSpace(r.FormValue(param)) == "" {
			fmt.Print(param)
			return false
		}
	}
	return true
}

func parseForm(r *http.Request) Entry {
	r.ParseForm()
	params := Entry{}
	for key, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			params[key] = trimSpace(value)
		}
	}
	return params
}

func trimSpace(para string) (ret string) {
	ret = strings.Trim(para, "\t\r\n ")
	return
}

func shiftChar(c int64, z int64) int64 {
	if 65 <= c && c <= 90 {
		return (c-39-z)%26 + 65
	} else if 97 <= c && c <= 122 {
		return (c-70-z)%26 + 97
	} else {
		return c // if not a letter return the same
	}
}

// TODO return two strings, the second one with the catched error if any, otherwise empty ?
func decryptMessage(m string, xyStr string, x *big.Int) (message string, err error) {
	xy := big.NewInt(0)
	_, success := xy.SetString(xyStr, 10) // DONE handle error if not a number
	if success == false {
		err = errors.New("key is not a number")
		return
	}
	y := new(big.Int)
	y.Quo(xy, x) // TODO handle if result not int
	// if y not int, create an erro and return
	z := 1 + y.Int64()%25
	n := int(math.Sqrt(float64(len(m)))) // TODO error if len(m) not a square number
	// Check if n*n == len(m), if not create an error and return
	decryptedM := make([]byte, len(m))

	// diagonalize and ceaser step, reverse
	j := 0
	for diag := 0; diag < 2*n-1; diag++ {
		s := 0
		if diag >= n {
			s = diag - n + 1
		}
		for i := s; i <= diag-s; i++ {
			decryptedM[j] = byte(shiftChar(int64(m[n*i+diag-i]), z))
			j++
		}
	}
	message = string(decryptedM)
	return
}

func testDecrypt(key *big.Int) {
	xy := "306063896731552281713201727176392168770237379582172677299123272033941091616817696059536783089054693601"
	m := "URYYBBJEX"
	decryptedM, err := decryptMessage(m, xy, key)
	if err != nil {
		fmt.Println(err)
	} else if decryptedM != "HELLOWORK" {
		fmt.Println("error decrypt: ", decryptedM)
	} else {
		fmt.Println("decrypt ok")
	}
}

// --------------------------------Query 2--------------------------------------
func q2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // Nov 1
	// if checkParams(r, []string{"userid", "tweet_time"}) == false {
	// 	fmt.Fprintf(w, "lack of parameters")
	// 	return
	// }

	params := parseForm(r)

	tweetId, score, text := q2Query(params["userid"], params["tweet_time"])
	teamInfo := "MDK,1474-9673-1103\n"
	message := teamInfo + tweetId + ":" + score + ":" + text + "\n"
	fmt.Fprintf(w, message)
	r.Body.Close()
	return
}

func q2Query(userid string, time string) (tweetId string, scoreStr string, censoredText string) {
	time = strings.Replace(time, "+", " ", 1)
	err := db.QueryRow("SELECT tweet_id,censored_text,score FROM tw WHERE created_at='"+time+"'"+" AND user_id="+userid).Scan(&tweetId, &censoredText, &scoreStr)
	errorHandle(err)
	censoredText = decode(censoredText)
	return
}

func q3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // Nov 1
	// if checkParams(r, []string{"userid", "start_date", "end_date", "n"}) == false {
	// 	fmt.Fprintf(w, "lack of parameters")
	// 	return
	// }
	params := parseForm(r)
	teamInfo := "MDK,1474-9673-1103\n"
	content := q3Query(params["userid"], params["start_date"], params["end_date"], params["n"])
	message := teamInfo + content

	fmt.Fprintf(w, message)
	r.Body.Close()
	return
}

func q3Query(userid string, start_date string, end_date string, n string) (content string) {
	content = "Positive Tweets\n"
	row1, err := db.Query("SELECT DATE(created_at) as date,impact,tweet_id,censored_text FROM tw WHERE user_id=" + userid + " AND DATE(created_at) BETWEEN '" + start_date + "' AND '" + end_date + "' AND impact > 0 ORDER BY impact DESC,tweet_id LIMIT " + n)

	//log.Println("Finished query")

	defer row1.Close()
	if err != nil {
		log.Println(err)
		errorHandle(err)
	}

	date := ""
	impact := ""
	tweet_id := ""
	text := ""
	for row1.Next() {
		row1.Scan(&date, &impact, &tweet_id, &text)
		if err != nil {
			log.Println(err)
			errorHandle(err)
		}
		//fmt.Print("text = " + text)
		text = decode(text)
		content += date + "," + impact + "," + tweet_id + "," + text + "\n"
		//log.Print(text + " " + tweet_id)
	}

	content += "\nNegative Tweets\n"
	row2, err := db.Query("SELECT DATE(created_at) as date,impact,tweet_id,censored_text FROM tw WHERE user_id=" + userid + " AND DATE(created_at) BETWEEN '" + start_date + "' AND '" + end_date + "' AND impact < 0 ORDER BY impact,tweet_id LIMIT " + n)
	defer row2.Close()
	if err != nil {
		log.Println(err)
		errorHandle(err)
	}
	date = ""
	impact = ""
	tweet_id = ""
	text = ""
	for row2.Next() {
		row2.Scan(&date, &impact, &tweet_id, &text)
		if err != nil {
			log.Println(err)
			errorHandle(err)
		}
		text = decode(text)
		content += date + "," + impact + "," + tweet_id + "," + text + "\n"
		//log.Print(text + " " + tweet_id)
	}

	return

}

func q4(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // Nov 1
	// if checkParams(r, []string{"hashtag", "n"}) == false {
	// 	fmt.Fprintf(w, "lack of parameters")
	// 	return
	// }
	params := parseForm(r)
	teamInfo := "MDK,1474-9673-1103\n"
	content := q4Query(params["hashtag"], params["n"])
	message := teamInfo + content
	fmt.Fprintf(w, message)
	r.Body.Close()
	return

}

func q4Query(hashtag string, n string) (content string) {
	content = ""
	row, err := db.Query("select h2.c,tw.text,h2.u,DATE(h2.date) from (select h.c,h1.user_id,h.u,h.date from (select count(*) as c,MIN(created_at) as date,GROUP_CONCAT(DISTINCT user_id ORDER BY user_id) as u from hashtags where hashtag = '" + hashtag + "' GROUP BY DATE(created_at) ORDER BY c DESC,DATE(created_at) LIMIT " + n + ") as h, hashtags as h1 where h1.hashtag='" + hashtag + "' and h1.created_at=h.date) as h2,tw where tw.user_id=h2.user_id and tw.created_at=h2.date")
	if err != nil {
		log.Println(err)
		errorHandle(err)
	}
	defer row.Close()
	count := ""
	datetime := ""
	users := ""
	text := ""
	for row.Next() {
		row.Scan(&count, &text, &users, &datetime)
		text = decode(text)
		content += datetime + ":" + count + ":" + users + ":" + text + "\n"
		// log.Print(text + " " + tweetId)
	}
	return
}

func q5(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // Nov 1
	// if checkParams(r, []string{"hashtag", "n"}) == false {
	// 	fmt.Fprintf(w, "lack of parameters")
	// 	return
	// }
	params := parseForm(r)
	teamInfo := "MDK,1474-9673-1103\n"
	content := q5Query(params["userid_min"], params["userid_max"])
	message := teamInfo + content + "\n"
	fmt.Fprintf(w, message)
	r.Body.Close()
	return

}

func q5Query(start string, end string) (count string) {
	err := db.QueryRow("select max.c-min.c from (select counts+total as c from usum where user_id="+end+") as max,(select total as c from usum where user_id="+start+") as min").Scan(&count)
	errorHandle(err)
	return
}

/*func encode(tag string) (encoded string) {
	con := Conversion{
		Text: tag,
	}
	b, err := json.Marshal(con);

}*/

func decode(oldStr string) (newStr string) {
	oldStr = strings.Replace(oldStr, "%", "%%", -1)
	jsonStr := "[{\"Text\":\"" + oldStr + "\"}]"
	var jsonText = []byte(jsonStr)
	var con []Conversion
	err = json.Unmarshal(jsonText, &con)
	if err != nil {
		fmt.Println("error:", err)
	}
	newStr = con[0].Text
	return
}

func parseSenti() {
	path := "../exFiles/afinn.txt"
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		word, score := s[0], s[1]
		value, err := strconv.Atoi(score)
		if err != nil {
			continue
		}
		sentimentMap[word] = value
		// log.Print(sentimentMap[word])
	}
}

func parseBanned() {
	path := "../exFiles/banned.txt"
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := trimSpace(scanner.Text())

		decoded := ""
		for _, r := range word {
			w := string(r)
			if r >= 'A' {
				if r <= 'M' {
					w = string(r + 13)
				} else if r >= 'N' && r <= 'Z' {
					w = string(r - 13)
				} else if r <= 'm' && r >= 'a' {
					w = string(r + 13)
				} else if r >= 'n' && r <= 'z' {
					w = string(r - 13)
				}
			}

			decoded += w
			// fmt.Println(i, r, string(decoded))
		}

		banMap[decoded] = true
	}
}

func errorHandle(err error) {
	if err != nil {
		panic(err)
	}
}

// -----------------------------------------------------------------------------

func main() {

	//load sentiment score file and censored text file
	parseSenti()
	parseBanned()

	// TODO put key in a file
	key.SetString(keyStr, 10)
	db.SetMaxOpenConns(1000000)
	db.SetMaxIdleConns(1000000)
	r := mux.NewRouter()
	r.HandleFunc("/", viewHandler)
	r.HandleFunc("/getw", phase1).Methods("GET")
	r.HandleFunc("/q1", q1).Methods("GET")
	r.HandleFunc("/q2", q2).Methods("GET")
	r.HandleFunc("/q3", q3).Methods("GET")
	r.HandleFunc("/q4", q4).Methods("GET")
	r.HandleFunc("/q5", q5).Methods("GET")
	r.HandleFunc("/q6", q6).Methods("GET")
	http.Handle("/", r)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}
