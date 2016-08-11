// package sync gets data from Ulule's API to store everything
// in a (local) redis database. Each sync is identified by a name.
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

	// leave blank to sync them all
	// otherwise, only invalid orders from a previous sync will
	// be considered
	syncOnlyInvalids := false
	invalidOrdersSyncName, err := linenoise.Line("only invalid orders from sync> ")
	if err != nil {
		logrus.Fatal(err)
	}

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		logrus.Fatal(err)
	}

	invalidOrders := make(map[string]struct{})

	if invalidOrdersSyncName != "" {
		syncOnlyInvalids = true
		invalids, errInvalidOrders := redis.Strings(redis.Values(conn.Do("SMEMBERS", invalidOrdersSyncName+"_invalidOrders")))
		if errInvalidOrders != nil {
			logrus.Fatal(errInvalidOrders)
		}
		for _, invalid := range invalids {
			invalidOrders[invalid] = struct{}{}
		}
	}

	// store syncName
	_, err = conn.Do("SADD", "syncs", syncName)
	if err != nil {
		logrus.Fatal(err)
	}

	// get projects information
	project, err := ululeClient.GetProject(strconv.Itoa(projectId))
	if err != nil {
		logrus.Fatal(err)
	}

	err = conn.Send("MULTI")
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Println("sync project \"" + project.Slug + "\"")

	err = conn.Send("HMSET", syncName+"_project",
		"id", project.Id,
		"url", project.Url,
		"goal", project.Goal,
		"goalRaised", project.GoalRaised,
		"amountRaised", project.AmountRaised,
		"commentCount", project.CommentCount,
		"committed", project.Committed,
		"currency", project.Currency,
		"currencyDisplay", project.CurrencyDisplay,
		"dateEnd", project.DateEnd,
		"dateStart", project.DateStart,
		"finished", project.Finished,
		"slug", project.Slug,
		"supportersCount", project.SupportersCount,
		"timeZone", project.TimeZone,
		"nbrewards", len(project.Rewards),
	)
	if err != nil {
		logrus.Fatal(err)
	}

	for i, reward := range project.Rewards {
		index := strconv.Itoa(i)
		err = conn.Send("HMSET", syncName+"_project",
			"reward"+index+"_id", reward.Id,
			"reward"+index+"_available", reward.Available,
			"reward"+index+"_price", reward.Price,
			"reward"+index+"_stock", reward.Stock,
			"reward"+index+"_stockAvailable", reward.StockAvailable,
			"reward"+index+"_stockTaken", reward.StockTaken,
		)
		if err != nil {
			logrus.Fatal(err)
		}

		logrus.Println("sync reward \"" + strconv.Itoa(int(reward.Id)) + " (" + strconv.Itoa(reward.Price) + " " + project.CurrencyDisplay + ")\"")
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		logrus.Fatal(err)
	}

	// sync orders
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

	// TODO: get supporters as well to be able to anotate anonymous users
	// anonymous = user ids in orders - user ids in supporters
	// in other words, if a user id from orders can't be found in the
	// list of supporters, it means the user wants to be anonymous

	for _, order := range orders {

		if syncOnlyInvalids == true {
			_, ok := invalidOrders[strconv.Itoa(int(order.Id))]
			if ok == false {
				continue
			}
		}

		// don't save cancelled orders
		if order.Status != clientapi.OrderStatusCancelled {

			err := conn.Send("MULTI")
			if err != nil {
				logrus.Fatal(err)
			}

			err = conn.Send("SADD", syncName, strconv.Itoa(int(order.Id)))
			if err != nil {
				logrus.Fatal(err)
			}

			// order.User.FirstName & order.User.LastName are not reliable
			// as these entries are totally optional, even for paying supporters
			// it's better to rely on order.ShippingAddress.FirstName &
			// order.ShippingAddress.LastName. If order.ShippingAddress doesn't
			// exist, keep empty strings for first and last names.

			firstName := ""
			lastName := ""

			shippingAddress := &clientapi.Address{}
			if order.ShippingAddress != nil {
				shippingAddress = order.ShippingAddress
				firstName = shippingAddress.FirstName
				lastName = shippingAddress.LastName

				// get rid of extra lines in Address1
				shippingAddress.Address1 = strings.Trim(shippingAddress.Address1, " \n\r")
				shippingAddress.Address1 = strings.Replace(shippingAddress.Address1, "\n", " ", -1)
				shippingAddress.Address1 = strings.Replace(shippingAddress.Address1, "\r", " ", -1)
				shippingAddress.Address1 = strings.Replace(shippingAddress.Address1, "  ", " ", -1)

				// get rid of extra lines in Address2
				shippingAddress.Address2 = strings.Trim(shippingAddress.Address2, " \n\r")
				shippingAddress.Address2 = strings.Replace(shippingAddress.Address2, "\n", " ", -1)
				shippingAddress.Address2 = strings.Replace(shippingAddress.Address2, "\r", " ", -1)
				shippingAddress.Address2 = strings.Replace(shippingAddress.Address2, "  ", " ", -1)

				// Apparenlty, some people think address2 is a confirmation of address1
				// so they put the exact same thing on both sides
				// in this case, just replacing address2 by empty string
				if shippingAddress.Address2 == shippingAddress.Address1 {
					shippingAddress.Address2 = ""
				}
			}

			// quick format for first & last name

			// first name
			firstName = strings.ToLower(firstName)
			firstName = strings.Trim(firstName, " -")
			// split on ' '
			firstNameParts := strings.Split(firstName, " ")
			for i, part := range firstNameParts {
				if len(part) > 1 {
					firstNameParts[i] = strings.ToUpper(string(part[0])) + part[1:]
				} else if len(part) > 0 {
					firstNameParts[i] = strings.ToUpper(string(part[0]))
				}
			}
			firstName = strings.Join(firstNameParts, " ")
			// split on '-'
			firstNameParts = strings.Split(firstName, "-")
			for i, part := range firstNameParts {
				if len(part) > 1 {
					firstNameParts[i] = strings.ToUpper(string(part[0])) + part[1:]
				} else if len(part) > 0 {
					firstNameParts[i] = strings.ToUpper(string(part[0]))
				}
			}
			firstName = strings.Join(firstNameParts, "-")
			firstName = strings.Trim(firstName, " ")

			// last name
			lastName = strings.ToLower(lastName)
			lastName = strings.Trim(lastName, " -")
			// split on ' '
			lastNameParts := strings.Split(lastName, " ")
			for i, part := range lastNameParts {
				if len(part) > 1 {
					lastNameParts[i] = strings.ToUpper(string(part[0])) + part[1:]
				} else if len(part) > 0 {
					lastNameParts[i] = strings.ToUpper(string(part[0]))
				}
			}
			lastName = strings.Join(lastNameParts, " ")
			// split on '-'
			lastNameParts = strings.Split(lastName, "-")
			for i, part := range lastNameParts {
				if len(part) > 1 {
					lastNameParts[i] = strings.ToUpper(string(part[0])) + part[1:]
				} else if len(part) > 0 {
					lastNameParts[i] = strings.ToUpper(string(part[0]))
				}
			}
			lastName = strings.Join(lastNameParts, "-")
			lastName = strings.Trim(lastName, " ")

			// TODO: improve formating considering dashes and other name patterns

			billingAddress := shippingAddress
			if order.BillingAddress != nil {
				billingAddress = order.BillingAddress

				// get rid of extra lines in Address1
				billingAddress.Address1 = strings.Trim(billingAddress.Address1, " \n\r")
				billingAddress.Address1 = strings.Replace(billingAddress.Address1, "\n", " ", -1)
				billingAddress.Address1 = strings.Replace(billingAddress.Address1, "\r", " ", -1)
				billingAddress.Address1 = strings.Replace(billingAddress.Address1, "  ", " ", -1)

				// get rid of extra lines in Address2
				billingAddress.Address2 = strings.Trim(billingAddress.Address2, " \n\r")
				billingAddress.Address2 = strings.Replace(billingAddress.Address2, "\n", " ", -1)
				billingAddress.Address2 = strings.Replace(billingAddress.Address2, "\r", " ", -1)
				billingAddress.Address2 = strings.Replace(billingAddress.Address2, "  ", " ", -1)

				// Apparenlty, some people think address2 is a confirmation of address1
				// so they put the exact same thing on both sides
				// in this case, just replacing address2 by empty string
				if billingAddress.Address2 == billingAddress.Address1 {
					billingAddress.Address2 = ""
				}
			}

			// if len(order.Items) != 1 {
			// 	fmt.Println("items:", len(order.Items), "url:", order.Url, "|", order.Total, "|", order.StatusDisplay, "|", order.User.UserName)
			// }

			err = conn.Send("HMSET", syncName+"_order_"+strconv.Itoa(int(order.Id)),

				// repeating hash key in values for flexibility
				"orderId", strconv.Itoa(int(order.Id)),

				"email", order.User.Email,
				"firstName", firstName,
				"lastName", lastName,
				"name", order.User.Name,
				"username", order.User.UserName,
				"datejoined", order.User.DateJoined,
				"userurl", order.User.Url,
				"userid", order.User.Id,
				// TODO: add anonymous field, value: 0 or 1

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
				"shippingState", shippingAddress.State,
				"shippingPhoneNumber", shippingAddress.PhoneNumber,
				"shippingEntityName", shippingAddress.EntityName,

				"billingAddr1", billingAddress.Address1,
				"billingAddr2", billingAddress.Address2,
				"billingCity", billingAddress.City,
				"billingCountry", billingAddress.Country,
				"billingCode", billingAddress.PostalCode,
				"billingState", billingAddress.State,
				"billingPhoneNumber", billingAddress.PhoneNumber,
				"billingEntityName", billingAddress.EntityName,

				"nbItems", len(order.Items),
			)
			if err != nil {
				logrus.Fatal(err)
			}

			for i, item := range order.Items {
				index := strconv.Itoa(i)
				err = conn.Send("HMSET", syncName+"_order_"+strconv.Itoa(int(order.Id)),
					"item"+index+"_quantity", item.Quantity,
					"item"+index+"_unitprice", item.UnitPrice,
					"item"+index+"_product", item.Product,
				)
			}

			// Put invalid orders in a set.
			// Invalid orders are the ones that can't be sent as there's something
			// incorrect with the shipping address. Basically if some fields are empty
			// This allows to build a new list later based on invalid orders that
			// got fixed by contributors.

			if (shippingAddress.Address1 == "" && shippingAddress.Address2 == "") || firstName == "" || lastName == "" || shippingAddress.City == "" || shippingAddress.Country == "" || shippingAddress.PostalCode == "" {
				conn.Send("SADD", syncName+"_invalidOrders", strconv.Itoa(int(order.Id)))
			}

			_, err = conn.Do("EXEC")
			if err != nil {
				logrus.Fatal(err)
			}
		}
	}

	nb, err := conn.Do("SCARD", syncName)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Println("nb orders:", nb)
}
