package main

import "os"
import "fmt"
import "net/url"
import "log"
import "bufio"
import "time"
import "regexp"
import "github.com/influxdb/influxdb/client"

const (
    MyPort        = 8086
    MyDB          = "reddit"
    MyMeasurement = "nb_users"
)

func writePoints(con *client.Client) {
    var (
        pts []client.Point
    )

    // 04/01/15:17:41:32 379 Guitar
    re := regexp.MustCompile("([0-9:/]+) ([0-9]+) ([a-zA-Z]+)")

/*
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	if err != nil {
	    log.Fatal(err)
	}
	defer file.Close()
*/

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		a := re.FindStringSubmatch(scanner.Text())
		if (len(a) == 0) {
			continue
		}
        // time.Parse is an interesting method ...
		t, err := time.Parse("01/02/06:15:04:05", a[1])
		if err != nil {
			fmt.Println(err);
		}
		v := a[2]
        sub := a[3]
        pts = append(pts, client.Point{
            Measurement: MyMeasurement,
            Tags: map[string]string{
                "subreddit": sub,
            },
            Time: t,
            Fields: map[string]interface{}{
                "value": v,
            },
            Precision: "s",
        })
	}

	if err := scanner.Err(); err != nil {
	    log.Fatal(err)
	}

    bps := client.BatchPoints{
        Points:          pts,
        Database:        MyDB,
        RetentionPolicy: "default",
    }
//	fmt.Println(bps);

	log.Printf("Pushing data ...")
    _, err := con.Write(bps)
    if err != nil {
        log.Fatal(err)
    }
}

func main() {

    // InfluxDB IP from the env
    myHost := os.Getenv("INFLUXDB_HOST")

    u, err := url.Parse(fmt.Sprintf("http://%s:%d", myHost, MyPort))
    if err != nil {
        log.Fatal(err)
    }

    conf := client.Config{
        URL:      *u,
        Username: os.Getenv("INFLUXDB_USER"),
        Password: os.Getenv("INFLUXDB_PWD"),
    }

    con, err := client.NewClient(conf)
    if err != nil {
        log.Fatal(err)
    }

    dur, ver, err := con.Ping()
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Happy as a Hippo! %v, %s", dur, ver)

	writePoints(con)
}
