package main

import (
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"sort"
	"strconv"
	"strings"
	"ulule/clientapi"
	// "ulule/credentials"
)

var ()

func main() {

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		logrus.Fatal(err)
	}

	values, err := redis.Values(conn.Do("SMEMBERS", "syncs"))
	if err != nil {
		logrus.Fatal(err)
	}

	for _, value := range values {
		b, ok := value.([]byte)
		if ok {
			fmt.Println(string(b))
		}
	}

	syncName, err := linenoise.Line("use sync> ")
	if err != nil {
		logrus.Fatal(err)
	}

	// display project summary
	values, err = redis.Values(conn.Do("HGETALL", syncName+"_project"))
	if err != nil {
		logrus.Fatal("HGETALL "+syncName+"_project error:", err)
	}
	stringMap, err := redis.StringMap(values, nil)
	if err != nil {
		logrus.Fatal("redis.StringMap err:", err)
	}
	// logrus.Printf("%+v", stringMap)

	fmt.Println("project:", stringMap["slug"])
	fmt.Println("rewards:")
	nbRewards, err := strconv.Atoi(stringMap["nbrewards"])
	if err != nil {
		logrus.Fatal(err)
	}

	for i := 0; i < nbRewards; i++ {
		index := strconv.Itoa(i)
		id := stringMap["reward"+index+"_id"]
		price := stringMap["reward"+index+"_price"]
		fmt.Println(id, price)
	}

	fmt.Println("commands: countries [rewards], exit")

	for {

		cmd, err := linenoise.Line("> ")
		if err != nil {
			logrus.Fatal(err)
		}

		parts := strings.Split(cmd, " ")
		if len(parts) > 0 {
			cmd = parts[0]
		}
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

		switch cmd {
		case "exit":
			break
		case "countries":
			displayCountries(syncName, conn, args)
		case "export":
			exportPaidOrders(syncName, conn, args)

		}
	}
}

func exportPaidOrders(syncName string, conn redis.Conn, rewardIDs []string) {
	values, err := redis.Values(conn.Do("SMEMBERS", syncName))
	if err != nil {
		logrus.Fatal(err)
	}

	conn.Send("MULTI")

	// tmp
	nbDisplayed := 0

	for _, value := range values {
		b, ok := value.([]byte)
		if ok {
			orderID := string(b)
			conn.Send("HGETALL", syncName+"_order_"+orderID)
		}
	}

	values, err = redis.Values(conn.Do("EXEC"))
	if err != nil {
		logrus.Fatal(err)
	}

	for _, value := range values {
		stringMap, _ := redis.StringMap(value, nil)

		nbItems, _ := strconv.Atoi(stringMap["nbItems"])
		if nbItems == 0 {
			continue
		}

		accept := true
		// filter reward ids
		if rewardIDs != nil && len(rewardIDs) > 0 {
			accept = false
			// HACK: considering there's only one
			// single item per order (id: 0)
			rewardID := stringMap["item0_product"]
			//logrus.Println(rewardID)

			// if rewardID == "177453" {
			// 	logrus.Printf("%#v", stringMap)
			// }

			for _, id := range rewardIDs {
				if id == rewardID {
					accept = true
					break
				}
			}
		}

		if accept {
			paymentStatus, _ := strconv.Atoi(stringMap["status"])
			accept = clientapi.OrderStatus(paymentStatus) == clientapi.OrderStatusPaymentDone
		}

		if accept {

			//logrus.Println(stringMap["statusDisplay"])
			// logrus.Println(stringMap)

			logrus.Println(stringMap["firstName"]+" "+stringMap["lastName"], "|",
				stringMap["shippingAddr1"], "|",
				stringMap["shippingAddr2"], "|",
				stringMap["shippingCity"], "|",
				stringMap["shippingCode"], "|",
				stringMap["shippingCountry"], "|",
				stringMap["email"], "|")

			nbDisplayed++
			if nbDisplayed >= 10 {
				break
			}
			// if _, exist := countries[stringMap["shippingCountry"]]; !exist {
			// 	countries[stringMap["shippingCountry"]] = 1
			// } else {
			// 	countries[stringMap["shippingCountry"]]++
			// }

			// nbContributions++
		}
	}

}

func displayCountries(syncName string, conn redis.Conn, rewardIDs []string) {
	values, err := redis.Values(conn.Do("SMEMBERS", syncName))
	if err != nil {
		logrus.Fatal(err)
	}

	conn.Send("MULTI")

	for _, value := range values {
		b, ok := value.([]byte)
		if ok {
			orderID := string(b)
			conn.Send("HGETALL", syncName+"_order_"+orderID)
		}
	}

	values, err = redis.Values(conn.Do("EXEC"))
	if err != nil {
		logrus.Fatal(err)
	}

	countries := make(map[string]int)
	nbContributions := 0

	for _, value := range values {
		stringMap, _ := redis.StringMap(value, nil)

		nbItems, _ := strconv.Atoi(stringMap["nbItems"])
		if nbItems == 0 {
			continue
		}

		accept := true
		// filter reward ids
		if rewardIDs != nil && len(rewardIDs) > 0 {
			accept = false
			// HACK: considering there's only one
			// single item per order (id: 0)
			rewardID := stringMap["item0_product"]
			//logrus.Println(rewardID)

			// if rewardID == "177453" {
			// 	logrus.Printf("%#v", stringMap)
			// }

			for _, id := range rewardIDs {
				if id == rewardID {
					accept = true
					break
				}
			}
		}

		if accept {
			paymentStatus, _ := strconv.Atoi(stringMap["status"])
			accept = clientapi.OrderStatus(paymentStatus) == clientapi.OrderStatusPaymentDone ||
				clientapi.OrderStatus(paymentStatus) == clientapi.OrderStatusInvalid
		}

		if accept {
			//logrus.Println(stringMap["statusDisplay"])

			if _, exist := countries[stringMap["shippingCountry"]]; !exist {
				countries[stringMap["shippingCountry"]] = 1
			} else {
				countries[stringMap["shippingCountry"]]++
			}

			nbContributions++
		}
	}

	pl := make(PairList, len(countries))
	i := 0

	for country, count := range countries {
		if country == "" {
			country = "none"
		}
		// fmt.Println(country, ":", count)
		pl[i] = Pair{country, count}
		i++
	}

	sort.Sort(sort.Reverse(pl))

	for _, pair := range pl {
		fmt.Println(pair.Key, ":", pair.Value)
	}
	fmt.Println("contributions:", nbContributions)
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
