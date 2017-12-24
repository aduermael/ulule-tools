package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	ulule "github.com/aduermael/ulule-api-client"
)

func main() {

	if len(os.Args) != 4 {
		log.Fatalln("wrong parameters, expecting:\nAPI_KEY USERNAME PROJECT_ID")
	}

	apiKey := os.Args[1]
	username := os.Args[2]
	projectIDStr := os.Args[3]

	ululeClient := ulule.ClientWithUsernameAndApiKey(username, apiKey)

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		log.Fatal(err)
	}

	project, err := ululeClient.GetProject(projectID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("----------------------------------")
	fmt.Println("Projet:", project.Slug)
	fmt.Println("----------------------------------")
	fmt.Println("Contreparties:")
	fmt.Println("----------------------------------")

	ticketsPerReward := make(map[int]int)

	rewardsByID := make(map[int]*ulule.Reward)

	for _, reward := range project.Rewards {
		rewardsByID[int(reward.ID)] = reward
	}

	var rewardIds []int
	for id := range rewardsByID {
		rewardIds = append(rewardIds, id)
	}
	sort.Ints(rewardIds)

	for _, id := range rewardIds {
		reward := rewardsByID[id]
		fmt.Printf("%.f (%d %s)\n", reward.ID, reward.Price, project.CurrencyDisplay)
		nbTickets := readInt("Combien de tickets par commande? ")
		ticketsPerReward[int(reward.ID)] = nbTickets
	}

	fmt.Println("----------------------------------")
	fmt.Println("Récupération des commandes:")
	fmt.Println("----------------------------------")

	offset := 0
	lastpage := false
	var orders []*ulule.Order

	// each ticket is an order ID
	tickets := make([]*ulule.Order, 0)

	nbOrdersByRewardID := make(map[int]int)
	nbInvalidOrdersByRewardID := make(map[int]int)

	for id := range rewardsByID {
		nbOrdersByRewardID[id] = 0
		nbInvalidOrdersByRewardID[id] = 0
	}

	for lastpage == false {

		orders, err, lastpage = ululeClient.GetProjectOrders(projectID, 100, offset)
		if err != nil {
			log.Fatal(err)
		}

		offset += len(orders)
		fmt.Printf("\r%d", offset)

		for _, order := range orders {

			if order.Status == ulule.OrderStatusPaymentDone {
				// Ulule only allows 1 item per order
				// just making sure of it...
				if len(order.Items) > 1 {
					fmt.Printf("\rorder %.f has more than one item (%d)\n", order.ID, len(order.Items))
					fmt.Println("order status:", order.StatusDisplay, "- total:", order.Total, project.CurrencyDisplay)
					for _, item := range order.Items {
						fmt.Println(item.Product, "- price:", item.UnitPrice, "- qty:", item.Quantity)
					}
					fmt.Printf("\r%d", offset)
				}

				for _, item := range order.Items {
					rewardID := item.Product

					if _, exists := ticketsPerReward[rewardID]; !exists {
						fmt.Println("no tickets for reward:", rewardID)
					} else if ticketsPerReward[rewardID] == 0 {
						fmt.Println("0 tickets for reward:", rewardID)
					}

					for i := 0; i < ticketsPerReward[rewardID]; i++ {
						tickets = append(tickets, order)
					}

					nbOrdersByRewardID[rewardID]++
				}
			} else { // status != ulule.OrderStatusPaymentDone
				// the numbers displayed on the project page include these
				// order statuses:
				// - OrderStatusPaymentDone
				// - OrderStatusInvalid
				// - OrderStatusCompleted
				// Tickets are created when status == OrderStatusPaymentDone
				// let's count invalid orders to make sure we land on our feet
				if order.Status == ulule.OrderStatusInvalid ||
					order.Status == ulule.OrderStatusCompleted {
					for _, item := range order.Items {
						rewardID := item.Product
						nbInvalidOrdersByRewardID[rewardID]++
					}
				}
			}
		}
	}

	fmt.Printf("\rTotal: %d\n", offset)
	fmt.Println("----------------------------------")

	for _, id := range rewardIds {
		reward := rewardsByID[id]
		fmt.Printf("%.f (%d %s) - commandes: %d (%d invalides, total: %d) - tickets: %d\n",
			reward.ID,
			reward.Price,
			project.CurrencyDisplay,
			nbOrdersByRewardID[id],
			nbInvalidOrdersByRewardID[id],
			nbOrdersByRewardID[id]+nbInvalidOrdersByRewardID[id],
			ticketsPerReward[id]*nbOrdersByRewardID[id])
	}

	fmt.Println("----------------------------------")

	fmt.Println("Nombre de tickets:", len(tickets))

	fmt.Println("----------------------------------")

	n := readInt("Combien de tickets gagnants? ")

	fmt.Println("----------------------------------")

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < n; i++ {
		// animation :)
		for j := 0; j < 50; j++ {
			index := rand.Intn(len(tickets) - 1)

			// erase line
			fmt.Printf("\r                                                                    ")

			name := ""

			firstName := strings.Title(strings.TrimSpace(strings.ToLower(tickets[index].User.FirstName)))
			lastName := strings.Title(strings.TrimSpace(strings.ToLower(tickets[index].User.LastName)))
			if len(lastName) > 1 {
				lastName = lastName[:1] + "."
			}

			name = strings.TrimSpace(firstName + " " + lastName)

			// first name and last name are optional fields, if the winner did
			// not complete them, just display "-anonyme-" instead
			if name == "" {
				name = "-anonyme-"
			}

			fmt.Printf("\r%s (%.f)", name, tickets[index].User.ID)
			time.Sleep(50 * time.Millisecond)
		}
		fmt.Printf("\n")
	}

	fmt.Println("----------------------------------")
}

func readInt(msg string) int {
	err := errors.New("")
	var i int

	for err != nil {
		fmt.Printf(msg)
		reader := bufio.NewReader(os.Stdin)
		str, _ := reader.ReadString('\n')
		str = strings.TrimSpace(str)
		i, err = strconv.Atoi(str)
		if err != nil {
			fmt.Println("nombre incorrect")
		}
	}
	return i
}
