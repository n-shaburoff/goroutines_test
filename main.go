package main

import (
	"context"
	"fmt"
	cfg "github.com/n-shaburoff/gourutines_test/config"
	"github.com/n-shaburoff/gourutines_test/resources"
	"github.com/spf13/viper"
	"sync"
	"time"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	var config cfg.Configurations

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	ch := make(chan resources.Data)
	wg := sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	go func(ctx context.Context) {
		defer wg.Done()
		go sender(config.Sender.Count, ch)
		go receiver(ch)
		for {
			select {
			case <-ctx.Done():
				cancel()
				return
			}
		}
	}(ctx)

	wg.Wait()
}

func sender(n time.Duration, c chan resources.Data) {
	tick := time.NewTicker(n)
	for {
		select {
		case <-tick.C:
			val := resources.Data{
				Timestamp: time.Now().Format(time.Stamp),
				Msg:       "Hi, I'm sender",
			}
			c <- val
		}
	}
}

func receiver(c chan resources.Data) {
	tick := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-tick.C:
			fmt.Printf("%v \n", <-c)
		}
	}
}
