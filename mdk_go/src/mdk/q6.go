/*
* @Author: Kaifu Wang, Mao Tang
* @Date:   2015-11-29 20:48:22
* @Last Modified by:   Kaifu Wang
* @Last Modified time: 2015-12-02 11:52:32
 */

package main

import (
	//"html/template"
	//"io/ioutil"
	// "bufio"
	// "database/sql"
	// "encoding/json"
	// "errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	// "github.com/gorilla/mux"
	// "log"
	// "math"
	// "math/big"
	"net/http"
	// "os"
	// "regexp"
	"strconv"
	// "strings"
	"time"
	//  "mime"
	//"reflect"
	// "sync"
	// "sync/atomic"
)

func q6(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // Nov 1
	params := parseForm(r)
	teamInfo := "MDK,1474-9673-1103\n"
	content := q6Query(params["tid"], params["seq"], params["opt"], params["tweetid"], params["tag"])

	message := teamInfo + content

	fmt.Fprintf(w, message)
	r.Body.Close()
	return
}

func q6Query(tid string, seq string, opt string, tweetid string, tag string) (content string) {
	if opt == "s" {
		db.Exec("UPDATE tw SET hashtags='' where tweet_id=" + tweetid)

		content = "0\n"

		if transacMap[tid] == nil {
			var t *transac
			tr := transac{tid: tid, seq: 0}
			t = &tr
			transacMap[tid] = t
		}

		transacMap[tid].seq = 1
		return
	}

	if opt == "e" && transacMap[tid] != nil {
		content = transacMap[tid].handleSeq(tid, "6", opt, tweetid, tag)
		transacMap[tid] = nil
		return
	}

	content = transacMap[tid].handleSeq(tid, seq, opt, tweetid, tag)
	return
}

func (t *transac) handleSeq(tid string, seq string, opt string, tweetid string, tag string) (content string) {
	if t == nil {
		// var t *transac
		tr := transac{tid: tid, seq: 0}
		t = &tr
		transacMap[tid] = &tr
	}
	t.Lock()
	s, _ := strconv.Atoi(seq)
	for t.seq != s {
		t.Unlock()
		time.Sleep(time.Millisecond * 10)
		t.Lock()
	}
	t.Unlock()
	
	if opt == "r" {
		var censoredText string
		var hashTag string
		err := db.QueryRow("SELECT censored_text,hashtags FROM tw WHERE tweet_id="+tweetid).Scan(&censoredText, &hashTag)
		if err != nil {
			fmt.Println(err)
		}
		content = decode(censoredText) + decode(hashTag) + "\n"
		t.increaseSeq()
		return
	}

	if opt == "a" {
		// _, err := db.Exec("UPDATE tw SET censored_text=concat(censored_text,'" + tag + "') where tweet_id=" + tweetid)
		_, err := db.Exec("UPDATE tw SET hashtags='" + tag + "' where tweet_id=" + tweetid)
		if err != nil {
			fmt.Println(err)
		}
		content = tag + "\n"
		t.increaseSeq()
		return
	}

	if opt == "e" {
		content = "0\n"
		return
	}
	return

}

func (t *transac) increaseSeq() {
	t.Lock()
	t.seq++
	t.Unlock()
}
