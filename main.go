package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	serve()
}

func serve() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/tax", tax)
	log.Fatal(http.ListenAndServe(":9002", nil))
}

func tax(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	monthly, _ := strconv.Atoi(v.Get("month"))

	var buf bytes.Buffer
	if monthly == 0 {
		buf.WriteString("字段名： 月薪 month, 每月扣减 minus, 年终奖 bonus")

		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		return
	}

	deduct, _ := strconv.Atoi(v.Get("minus"))
	bonus, _ := strconv.Atoi(v.Get("bonus"))
	start := 5000

	buf.WriteString(fmt.Sprintf("每月工资: %d，起征点: %d，每月扣减金额为: %d \n", monthly, start, deduct))

	var tax float64
	var accuTotal float64

	for i := 1; i <= 12; i++ {
		accu := (monthly - start - deduct) * i
		accuTotal = rate(accu)
		buf.WriteString(fmt.Sprintf("%d月: %.2f \n", i, accuTotal-tax))
		tax += accuTotal - tax
	}

	accu := (monthly-start-deduct)*12 + bonus
	before := monthly*12 + bonus
	accuTotal = rate(accu)

	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintln("年终奖: ", bonus))
	buf.WriteString(fmt.Sprintln("年终奖扣税: ", accuTotal-tax))

	tax += accuTotal - tax
	after := float64(before) - tax

	buf.WriteString(fmt.Sprintln("累计应纳税所得:", accu))
	buf.WriteString(fmt.Sprintln("总共扣税:", tax))
	buf.WriteString(fmt.Sprintln("年薪税前:", before))
	buf.WriteString(fmt.Sprintln("年薪税后:", after))

	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

func rate(m int) float64 {
	money := float64(m)
	switch {
	case money <= 36000:
		return money * 0.03
	case money > 36000 && money <= 144000:
		return money*0.1 - 2520 //速算数
	case money > 144000 && money <= 300000:
		return money*0.2 - 16920
	case money > 300000 && money <= 42000:
		return money*0.25 - 31920
	case money > 420000 && money <= 660000:
		return money*0.30 - 52920
	case money > 660000 && money <= 960000:
		return money*0.35 - 85920
	default:
		return money*0.45 - 181920
	}
}
