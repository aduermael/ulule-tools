package main

import (
	"fmt"
	"github.com/GeertJohan/go.linenoise"
	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"ulule/clientapi"
	"ulule/credentials"
)

var ()

func main() {
	username, apikey := credentials.Get(linenoise.Line)
	ululeClient := clientapi.New(username, apikey)

	projects, err := ululeClient.GetProjects("created")
	if err != nil {
		logrus.Fatal(err)
	}
	for _, project := range projects {
		percentage := int(float32(project.AmountRaised) / float32(project.Goal) * 100.0)
		percentageStr := strconv.Itoa(percentage)
		fmt.Println(project.Id, "|", project.Slug, "|", project.AmountRaised, project.CurrencyDisplay, "|", percentageStr+"%")
	}

	projectIdStr, err := linenoise.Line("project id> ")
	if err != nil {
		logrus.Fatal(err)
	}

	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		logrus.Fatal(err)
	}

	// Allows to simply sync Ulule orders several
	// times under specific names
	syncName, err := linenoise.Line("sync name> ")
	if err != nil {
		logrus.Fatal(err)
	}

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		logrus.Fatal(err)
	}

	offset := 0
	limit := 1000

	lastPage := false
	orders := make([]*clientapi.Order, 0)

	for !lastPage {
		var newOrders []*clientapi.Order
		newOrders, err, lastPage = ululeClient.GetProjectOrders(projectId, limit, offset)
		if err != nil {
			logrus.Fatal(err)
		}
		offset += len(newOrders)
		orders = append(orders, newOrders...)
		logrus.Println("orders:", offset)
	}

	for _, order := range orders {

		// don't save cancelled orders
		if order.Status != clientapi.OrderStatusCancelled {
			_, err := conn.Do("SADD", syncName, strconv.Itoa(int(order.Id)))
			if err != nil {
				logrus.Fatal(err)
			}
		}

		// quick format for first & last name
		firstName := order.User.FirstName
		lastName := order.User.LastName

		firstName = strings.ToUpper(firstName)
		if len(firstName) > 1 {
			firstName = string(firstName[0]) + strings.ToLower(firstName[1:])
		}
		lastName = strings.ToUpper(lastName)
		if len(lastName) > 1 {
			lastName = string(lastName[0]) + strings.ToLower(lastName[1:])
		}

		order.User.FirstName = firstName
		order.User.LastName = lastName

		shippingAddress := &clientapi.Address{}
		if order.ShippingAddress != nil {
			shippingAddress = order.ShippingAddress
		}
		billingAddress := shippingAddress
		if order.BillingAddress != nil {
			billingAddress = order.BillingAddress
		}

		_, err := conn.Do("HMSET", syncName+"_order_"+strconv.Itoa(int(order.Id)),
			"email", order.User.Email,
			"firstName", order.User.FirstName,
			"lastName", order.User.LastName,
			"name", order.User.Name,
			"username", order.User.UserName,
			"datejoined", order.User.DateJoined,
			"userurl", order.User.Url,
			"userid", order.User.Id,

			"total", order.Total,
			"method", order.PaymentMethod,
			"status", order.Status,
			"statusDisplay", order.StatusDisplay,
			"url", order.Url,

			"shippingAddr1", shippingAddress.Address1,
			"shippingAddr2", shippingAddress.Address2,
			"shippingCity", shippingAddress.City,
			"shippingCountry", shippingAddress.Country,
			"shippingCode", shippingAddress.PostalCode,

			"billingAddr1", billingAddress.Address1,
			"billingAddr2", billingAddress.Address2,
			"billingCity", billingAddress.City,
			"billingCountry", billingAddress.Country,
			"billingCode", billingAddress.PostalCode,
		)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	nb, err := conn.Do("SCARD", syncName)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Println("nb orders:", nb)
}
