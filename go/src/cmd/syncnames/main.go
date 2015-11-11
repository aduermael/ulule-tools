package main

import (
	"github.com/GeertJohan/go.linenoise"
	"github.com/Sirupsen/logrus"
	"strconv"
	"ulule/clientapi"
	"ulule/credentials"
)

var ()

func main() {
	username, apikey := credentials.Get(linenoise.Line)
	ululeClient := clientapi.New(username, apikey)

	projectIdStr, err := linenoise.Line("project id> ")
	if err != nil {
		logrus.Fatal(err)
	}

	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		logrus.Fatal(err)
	}

	offset := 0
	limit := 100

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

}
