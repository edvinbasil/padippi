package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func printUsage() {
	fmt.Println(
		`
    Usage
    -----
    classalerts <command> [options]

    commands are
    daily           :   send daily timetable
    printtt <oday>  :   print timetable for the given day
                        where oday is an offset specifier of the format
                        number  : current_day + offset
                        -number : current_day - offset
                        =number : offset(th) day
                        monday is 0. friday is 4. offset if bounded
    sendtt <hour>   :   sends the timetable for the given hour, for the current day

    Examples
    --------
    classalerts daily

    classalerts printtt +1
    classalerts printtt -4
    classalerts printtt =3

    classalerts sendtt 8

        `)
}

type timeTable map[string]hrTable
type hrTable [][]subject
type subject struct {
	slug   string
	roleid string
	link   string
}

func (tt timeTable) getHrNames(hr int, day int) string {
	daytt := tt[fmt.Sprint(hr)][day]
	hrlist := ""
	for _, item := range daytt {
		hrlist += item.slug + " "
	}
	return hrlist
}
func (tt timeTable) getTTDay(day string) string {
	cur_weekday := int(time.Now().Local().Weekday()) - 1
	lookupday := cur_weekday
	if strings.HasPrefix(day, "-") {
		offset, _ := strconv.Atoi(strings.TrimPrefix(day, "-"))
		lookupday = lookupday - offset
	} else if strings.HasPrefix(day, "=") {
		offset, _ := strconv.Atoi(strings.TrimPrefix(day, "="))
		lookupday = offset
	} else {
		offset, _ := strconv.Atoi(day)
		lookupday = lookupday + offset
	}
	if lookupday > 4 {
		lookupday = lookupday % 5
	} else if lookupday < 0 {
		lookupday = 5 + (lookupday % 5)
	}
	hours := []int{8, 9, 10, 11, 1, 2, 3, 4, 5}
	ttoftheday := ""
	for _, hr := range hours {
		ttoftheday += fmt.Sprint(hr, ": ", tt.getHrNames(hr, lookupday), "\n")
	}
	daymap := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	return fmt.Sprintf("TimeTable for %s\n\n%s", daymap[lookupday], ttoftheday)
}
func (tt timeTable) getAlertRoles(day int, hr int) string {
	daytt := tt[fmt.Sprint(hr)][day]
	var rlist []string
	for _, sub := range daytt {
		rlist = append(rlist, sub.roleid)
	}
	return strings.Join(rlist, " ")
}
func (tt timeTable) getAlertLinks(day int, hr int) string {
	daytt := tt[fmt.Sprint(hr)][day]
	var llist []string
	for _, sub := range daytt {
		if sub.link != "" {
			llist = append(llist, fmt.Sprint(sub.slug, ": ", sub.link, "\n"))
		}
	}
	return strings.Join(llist, "\n")
}

func main() {
	ml := subject{"ml", "<@&798101716755152906>", ""}
	cs := subject{"c.sec", "<@&798102007452794911>", ""}
	tida := subject{"ti.da", "<@&798102312073166868>", ""}
	ca := subject{"c.algebra", "<@&798102678542614529>", ""}
	fs := subject{"formal.s", "<@&798102849443069992>", ""}
	cg := subject{"c.graphics", "<@&798102960877076490>", ""}
	dm := subject{"dm", "<@&798103196396158986>", ""}
	ct := subject{"c.theory", "<@&798103275252744222>", ""}
	pa := subject{"p.algo", "<@&798103411156451348>", ""}
	compgeo := subject{"compgeo", "<@&798109507124330497>", ""}
	carch := subject{"c.arch", "<@&798419028490059796>", ""}
	geoinfo := subject{"geoinfo", "<@&798419761700405308>", "https://eduserver.example.com/mod/webexactivity/view.php?id=XXXXXX&action=joinmeeting"}

	// Apart from x, these correspond to the A-H slots given in the timetable
	// x is a placeholder for an empty slot
	x := []subject{}
	a := []subject{}
	b := []subject{}
	c := []subject{ml, cs}
	d := []subject{tida, geoinfo}
	e := []subject{ca, fs, compgeo}
	f := []subject{dm, cg, ct}
	g := []subject{pa}
	h := []subject{carch}

	// dont remember why these exist :)
	// prolly correspond to plus slots
	ap := a
	bp := b
	cp := c
	dp := d
	ep := e
	fp := f
	gp := g
	hp := h

	// dont remember what these were either. oops
	ea := e
	ga := g

	// This corresponds the the timetable in a table format
	// y axis correspond the the starting hr
	// x axis correspond the the day. [mon-fri]
	// eg:
	//		the following represents the timetable for the 8.00AM slots
	//
	//		"8":  {a, b, c, d, {ca, fs}},
	//
	//		mon - A slot,
	//		tue - b slot
	//		...
	//		fri - ca and fs slot
	tt := timeTable{
		"8":  {a, b, c, d, {ca, fs}},
		"9":  {f, g, a, b, c},
		"10": {d, e, f, g, {compgeo}},
		"11": {b, c, d, e, f},
		"1":  {g, ap, h, x, hp},
		"2":  {ep, fp, gp, cp, bp},
		"3":  {x, x, x, x, x},
		"4":  {x, x, x, x, x},
		"5":  {h, h, ea, ga, dp},
	}

	webhookurl := os.Getenv("DISCORD_WEBHOOK_URL")

	// TestURL
	//webhookurl := "https://discordapp.com/api/webhooks/XXXXXXXXXXXXXXXXXX/YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"

	// Actual URL
	// webhookurl := "https://discordapp.com/api/webhooks/XXXXXXXXXXXXXXXXXX/YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY"
	if webhookurl == "" {
		fmt.Println("No webhook URL found. Please set the DISCORD_WEBHOOK_URL env variable")
		os.Exit(2)
	}
	arglen := len(os.Args)
	if arglen < 2 {
		fmt.Println("No argument passed")
		printUsage()
		return
	}
	cmdopt := os.Args[1]
	if cmdopt == "daily" {
		// send webhook daily
		embed := []map[string]string{{
			"title":       "TimeTable for the day",
			"description": tt.getTTDay("0"),
			"color":       "5373784",
		}}

		m := map[string]interface{}{
			"username": "Padippi.go",
			"embeds":   embed,
		}
		mJson, _ := json.Marshal(m)
		contentReader := bytes.NewReader(mJson)
		req, _ := http.NewRequest("POST", webhookurl, contentReader)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, _ := client.Do(req)
		fmt.Println(resp.StatusCode)
	} else if cmdopt == "printtt" {
		if arglen != 3 {
			fmt.Println("No day given")
			return
		}
		day := os.Args[2]
		ttoftheday := tt.getTTDay(day)
		fmt.Println(ttoftheday)
	} else if cmdopt == "sendtt" {
		hr, _ := strconv.Atoi(os.Args[2])
		cur_weekday := int(time.Now().Local().Weekday()) - 1
		rolelist := tt.getAlertRoles(cur_weekday, hr)
		linklist := tt.getAlertLinks(cur_weekday, hr)
		// fmt.Println("Rolelist:", rolelist)
		// fmt.Println("LinkList:", linklist)
		alertmsg := ""
		if !(rolelist == "") {
			alertmsg = fmt.Sprintln("Class in 10 minutes", rolelist)
			if !(linklist == "") {
				alertmsg = fmt.Sprintf("%s\nLinks\n-------\n%s\n-----------------------",
					alertmsg, linklist)
			}
			// fmt.Println("Alert:", alertmsg)
			m := map[string]interface{}{
				"username": "Padippi.go",
				"content":  alertmsg,
			}
			mJson, _ := json.Marshal(m)
			contentReader := bytes.NewReader(mJson)
			req, _ := http.NewRequest("POST", webhookurl, contentReader)
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			resp, _ := client.Do(req)
			fmt.Println(resp.StatusCode)
		} else {
			fmt.Println("No classes boii")
		}

	} else {
		fmt.Println("Invalid Command")
		printUsage()
		return
	}
}
