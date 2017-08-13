package main

import (
	"fmt"
	"strconv"

	owfs "github.com/dhiltgen/go-owfs.git"
)

type OneWireClient struct {
	client *owfs.OwfsClient
}

func NewOneWireClient(connection string) (*OneWireClient, error) {
	client, err := owfs.NewClient(connection)
	if err != nil {
		return nil, err
	}

	return &OneWireClient{
		client: &client,
	}, nil
}

func (c *OneWireClient) Collect() (*Metrics, error) {
	metrics := &Metrics{}

	sensors, err := c.client.Dir("/")
	if err != nil {
		return metrics, err
	}

	for _, sensor := range sensors {
		temp, err := c.getTempSensor(sensor)
		if err == nil {
			metrics.Temperatures = append(metrics.Temperatures, temp)
		}
	}

	return metrics, nil
}

func (c *OneWireClient) getTemperature(address string) (float64, error) {
	path := fmt.Sprintf("%s/temperature9", address)
	data, err := c.client.Read(path)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(string(data), 64)
}

func (c *OneWireClient) getTempSensor(address string) (TempSensor, error) {
	temp, err := c.getTemperature(address)
	if err != nil {
		return TempSensor{}, err
	}

	return TempSensor{
		Address:     address,
		Temperature: temp,
	}, nil
}
