package rosco

import (
	"bufio"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"time"
)

const memsBaudRate = 9600

// MemsConnection contains all the configuration necessary
// to open a serial port
type MemsConnection struct {
	config     *serial.Config
	port       *serial.Port
	portReader *bufio.Reader
	portChan   chan []byte
	stateChan  chan error
	waiting    bool
	response   []byte
}

// NewMemsConnection returns a pointer to a Connect instance
func NewMemsConnection(portPath string) (*MemsConnection, error) {
	config := serial.Config{Name: portPath, Baud: memsBaudRate, ReadTimeout: time.Nanosecond}
	port, err := serial.OpenPort(&config)
	if err != nil {
		return nil, err
	}
	portReader := bufio.NewReader(port)
	stateChan := make(chan error)
	return &MemsConnection{config: &config, port: port,
		portReader: portReader,
		stateChan:  stateChan}, nil
}

// Start initializes a read loop that attempts to reconnect
// when the connection is broken
func (c *MemsConnection) Start() {
	for {
		select {
		case err := <-c.stateChan:
			if err != nil {
				fmt.Printf("Error connecting to %s", c.config.Name)
				go c.initialize()
			} else {
				fmt.Printf(" | Connection to %s reestablished!", c.config.Name)
			}
		}
	}
}

func (c *MemsConnection) initialize() {
	c.port.Close()
	for {
		time.Sleep(time.Second)
		port, err := serial.OpenPort(c.config)
		if err != nil {
			continue
		}
		c.port = port
		c.portReader = bufio.NewReader(port)
		c.stateChan <- nil
		return
	}
}

// Read loop from serial port
func (c *MemsConnection) readLoop() {
	for {
		response, err := c.portReader.ReadBytes('\n')
		// report the error
		if err != nil && err != io.EOF {
			c.stateChan <- err
			return
		}
		if len(response) > 0 {
			c.waiting = false
			c.response = response
			log.Printf("ECU: %x\r\n", response)
		}
	}
}

func (c *MemsConnection) Read() []byte {
	response, err := c.portReader.ReadBytes('\n')
	n := len(response)

	// report the error
	if err != nil && err != io.EOF {
		c.stateChan <- err
	}

	if n > 0 {
		c.waiting = false
		c.response = response[0:n]
		log.Printf("ECU < %x\r\n", response)
		return response
	}

	return nil
}

func (c *MemsConnection) Write(message []byte) {
	c.waiting = true
	_, err := c.port.Write(message)
	if err != nil {
		fmt.Printf("Error writing to serial port: %v ", err)
		c.waiting = false
	} else {
		log.Printf("FCR > %x\r\n", message)
	}
}
