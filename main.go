package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const TIME_FMT = "2006-01-02 15:04"

type RecordsList []time.Time

func (a RecordsList) Len() int {
	return len(a)
}
func (a RecordsList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a RecordsList) Less(i, j int) bool {
	return !a[i].After(a[j])
}

type DailyRecord map[string]RecordsList // key: date, value:date-time
type Checker struct {
	Records map[string]DailyRecord // key: name, value: DailyRecord
}

func NewChecker() Checker {
	var o Checker
	o.Records = make(map[string]DailyRecord)
	return o
}

func (o *Checker) ReadFile(name string) {
	// Use iconv to convert csv file from gbk to utf-8 as input stream.
	cmd := exec.Command("iconv", "-f", "gbk", "-t", "utf-8", name)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// Read records from csv.
	r := csv.NewReader(stdout)
	items, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	for i, item := range items {
		name := item[1]
		if name == "姓名" {
			continue
		}
		dailyRecord, ok := o.Records[name]
		if !ok {
			dailyRecord = make(DailyRecord)
			o.Records[name] = dailyRecord
		}

		date_time := item[3]
		date := strings.Fields(date_time)[0]
		dateTime, err := time.Parse(TIME_FMT, date_time)
		if err != nil {
			log.Fatal(err)
		}
		dailyRecord[date] = append(dailyRecord[date], dateTime)
		if i == 13 {
			list := dailyRecord[date]
			sort.Sort(list)
			fmt.Println(list)
		}
	}
}

func (o *Checker) Process() {
	for _, daily := range o.Records {
		for _, records := range daily {
			sort.Sort(records)
		}
	}
}

func (o *Checker) WriteResultRaw() error {
	f, err := os.Create("haha2.xls")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f)
	w.Write([]string{"编号", "姓名", "年龄"})
	w.Write([]string{"1", "张三", "23"})
	w.Write([]string{"2", "李四", "24"})
	w.Write([]string{"3", "王五", "25"})
	w.Write([]string{"4", "赵六", "26"})
	w.Flush()
	return nil
}

func main() {
	var (
		checker Checker = NewChecker()
	)
	checker.ReadFile("InOutData.csv")
	//
	fmt.Println("ok")
}
