package main

import (
	hazelcast "github.com/hazelcast/hazelcast-go-client"
	
	//"reflect"
	"context"
	"errors"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	//"github.com/golang-collections/collections"
	"time"
	//"github.com/hazelcast/hazelcast-go-client/serialization"
	//"strconv"
	//"github.com/hazelcast/hazelcast-go-client/config"
    //"github.com/hazelcast/hazelcast-go-client/config/property"
    
)


var httpClient = &http.Client{}
var key = "key"

func main() {
	/*ctx := context.TODO()
	client, _ := hazelcast.StartNewClient(ctx)
	m, _ := client.GetMap(ctx, "map")
	for i := 0; i < 1000; i++ {
		key := strconv.Itoa(i)
		m.Set(ctx, key, "message" + key)

	}
	//m.Destroy(ctx)
	client.Shutdown(ctx)*/
	manageConnections()
	
}

func manageConnections() {
	ctx := context.Background()
	go updateWithoutLock("5701", ctx)
	go updateWithoutLock("5702", ctx)
	updateWithoutLock("5703", ctx)
}

func updateWithoutLock(port string, ctx context.Context) {
	//ctx := context.TODO()
	//config := hazelcast.Config{}
	config := hazelcast.NewConfig()
	config.Cluster.Network.SetAddresses("localhost:" + port)	
	client, _ := hazelcast.StartNewClientWithConfig(ctx, config)
	m, err := client.GetMap(ctx, "foo")
	if err != nil {
		log.Fatal(err)
	}
	
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	fmt.Println("got map on port " + port)
	m.Put(ctx, key, 0)
	for i := 0; i < 1000; i++ {		
		content, err := m.Get(ctx, key)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}
			panic(err)
		}
		newVal := content.(int64) + 1
		time.Sleep(10 * time.Millisecond)
		oldVal, _ := m.Put(ctx, key, newVal)
		fmt.Println("№ ", i, "updated successfully on port "+ port, " value is: ", newVal, "old value is: ", oldVal, " old must be ", content)
	}
	client.Shutdown(ctx)
}

func managePessimisticLock() {
	ctx := context.Background()
	go updateWithPessimisticLock("5701", ctx)
	go updateWithPessimisticLock("5702", ctx)
	updateWithPessimisticLock("5703", ctx)
}

func updateWithPessimisticLock(port string, ctx context.Context) {
	client := getClient(port, ctx)
	testMap, err := client.GetMap(ctx, "foo")
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	fmt.Println("got map on port " + port)
	testMap.Put(ctx, key, 0)
	for i := 0; i < 1000; i++ {
		_ = testMap.Lock(ctx, key)
		fmt.Println("№ ", i, " locked on port " + port)
		value, _ := testMap.Get(ctx, key)
		fmt.Println("№ ", i, "oldValue ", value, " retrieved on port " + port)
		newVal := value.(int64) + 1
		time.Sleep(10 * time.Millisecond)
		oldVal, _ := testMap.Put(ctx, key, newVal)
		err := testMap.Unlock(ctx, key)
		if err != nil {
			fmt.Println("cannot write on port " + port)
		} else {
			fmt.Println("№ ", i, " updated successfully on port "+port, " value is: ", newVal, "old value is: ", oldVal, "old must be ", value)
		}
	}
	testMap.Destroy(ctx)
}
func manageOptimisticLock() {
	ctx := context.Background()
	go updateWithOptimisticLock("5701", ctx)
	go updateWithOptimisticLock("5702", ctx)
	updateWithOptimisticLock("5703", ctx)
}

func updateWithOptimisticLock(port string, ctx context.Context) {
	client := getClient(port, ctx)
	testMap, err := client.GetMap(ctx, "map")
	if err != nil {
		log.Fatal(err)
	}
	
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	
	testMap.Put(ctx, key, 0)
	fmt.Println("got map on port " + port)
	for i := 0; i < 1000; i++ {
		for {
			value, _ := testMap.Get(ctx, key)
			fmt.Println("oldvalue ", value, " retrieved on port "+port)
			newVal := value.(int64) + 1
			time.Sleep(20 * time.Millisecond)
			isReplaced, _ := testMap.ReplaceIfSame(ctx, key, value, newVal)
			if isReplaced {
				fmt.Println("№ ", i, " port: "+port, " || value updated successfully from ", value, " to ", newVal)
				break
			} else {
				fmt.Println("№ ", i, " port: "+port, " || value changed during transaction. Try again")
			}
		}
	}
}

func manageQueue() {
	ctx := context.Background()
	go readFromQueue("5702", ctx)
	go readFromQueue("5703", ctx)
	writeToQueue("5701", ctx)
}

func writeToQueue(port string, ctx context.Context) {
	customQueue := getQueue(port, ctx)
	maxSize := 10
	for i := 0; i < 100; i++ {
		err := customQueue.Put(ctx,i)
		if err != nil {
			fmt.Println(err)
		}
		sizeQueue, _ := customQueue.Size(ctx)
		fmt.Println("Size is ", sizeQueue)
		if sizeQueue >= maxSize{
			break}
		//fmt.Println(i, "added")
		time.Sleep(100)
	}
	customQueue.Destroy(ctx)
}

func readFromQueue(port string, ctx context.Context) {
	customQueue := getQueue(port, ctx)
	for {
		index, _ := customQueue.Take(ctx)
		fmt.Println(index, "is consumed on port "+port)
		time.Sleep(100)
	}
}

func getQueue(port string, ctx context.Context) *hazelcast.Queue {
	client := getClient(port, ctx)
	customQueue, _ := client.GetQueue(ctx, "customQueue")
	return customQueue
}

func getClient(port string, ctx context.Context) *hazelcast.Client {
	config := hazelcast.NewConfig()
	config.Cluster.Network.SetAddresses("localhost:" + port)
	client, _ := hazelcast.StartNewClientWithConfig(ctx, config)
	return client
}