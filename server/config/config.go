package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Config defines Server configuration interface
type Config interface {
	Address() string
	Port() uint
	ListenAddr() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
	ShutDownTimeout() time.Duration

	TaskDelay() time.Duration
	TotalWorkers() uint
	QueueSize() uint
}

// DefaultConfig defines default server configuration
type DefaultConfig struct{}

var _ Config = &DefaultConfig{}

func (c *DefaultConfig) Address() string {
	return "0.0.0.0"
}

func (c *DefaultConfig) Port() uint {
	return 8080
}

func (c *DefaultConfig) ListenAddr() string {
	return fmt.Sprintf("%v:%v", c.Address(), c.Port())
}

func (c *DefaultConfig) ReadTimeout() time.Duration {
	return 10 * time.Second
}

func (c *DefaultConfig) WriteTimeout() time.Duration {
	return 10 * time.Second
}

func (c *DefaultConfig) IdleTimeout() time.Duration {
	return 10 * time.Second
}

func (c *DefaultConfig) ShutDownTimeout() time.Duration {
	return 30 * time.Second
}

func (c *DefaultConfig) TaskDelay() time.Duration {
	return 5 * time.Second
}

func (c *DefaultConfig) TotalWorkers() uint {
	return uint(runtime.NumCPU())
}

func (c *DefaultConfig) QueueSize() uint {
	return 10000
}

type config struct {
	addr            string
	port            uint
	readTimeout     time.Duration
	writeTimeout    time.Duration
	idleTimeout     time.Duration
	shutDownTimeout time.Duration

	taskDelay    time.Duration
	totalWorkers uint
	queueSize    uint
}

func parseTimeout(timeout string) (time.Duration, error) {
	parsedTimeout, err := strconv.ParseInt(timeout, 10, 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(parsedTimeout) * time.Second, nil
}

func parseEnvTimeout(def time.Duration, envName string) (time.Duration, error) {
	timeout, ok := os.LookupEnv(envName)
	if ok {
		parsedTimeout, err := parseTimeout(timeout)
		if err != nil {
			return 0, fmt.Errorf("Cannot parse %v '%v': %w", envName, timeout, err)
		}
		return time.Duration(parsedTimeout) * time.Second, nil
	}
	return def, nil
}

// New handles Server configuration from env variables
func New() (Config, error) {
	c := &config{}
	err := c.parseEnv()
	if err != nil {
		return nil, err
	}

	// cli flags will override ENV vars
	err = c.parseFlags()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *config) parseEnv() error {
	def := DefaultConfig{}

	// addr
	c.addr = def.Address()
	rawAddr, ok := os.LookupEnv("HASH_SERVER_ADDR")
	if ok {
		c.addr = rawAddr
	}

	// port
	c.port = def.Port()
	rawPort, ok := os.LookupEnv("HASH_SERVER_PORT")
	if ok {
		parsedPort, err := strconv.ParseUint(rawPort, 10, 64)
		if err != nil {
			return fmt.Errorf("Cannot parse HASH_SERVER_PORT '%v': %w", rawPort, err)
		}
		if parsedPort < 1 || parsedPort > 65535 {
			return fmt.Errorf("Invalid HASH_SERVER_PORT value '%v', valid range [1..65535]", parsedPort)
		}
		c.port = uint(parsedPort)
	}

	// timeouts
	var err error
	c.readTimeout, err = parseEnvTimeout(def.ReadTimeout(), "HASH_SERVER_READ_TIMEOUT")
	if err != nil {
		return err
	}

	c.writeTimeout, err = parseEnvTimeout(def.WriteTimeout(), "HASH_SERVER_WRITE_TIMEOUT")
	if err != nil {
		return err
	}

	c.idleTimeout, err = parseEnvTimeout(def.IdleTimeout(), "HASH_SERVER_IDLE_TIMEOUT")
	if err != nil {
		return err
	}

	c.shutDownTimeout, err = parseEnvTimeout(def.ShutDownTimeout(), "HASH_SERVER_SHUTDOWN_TIMEOUT")
	if err != nil {
		return err
	}

	c.taskDelay, err = parseEnvTimeout(def.TaskDelay(), "HASH_TASK_DELAY")
	if err != nil {
		return err
	}

	c.totalWorkers = def.TotalWorkers()
	rawTotalWorkers, ok := os.LookupEnv("HASH_TOTAL_WORKERS")
	if ok {
		totalWorkers, err := strconv.ParseUint(rawTotalWorkers, 10, 64)
		if err != nil {
			return fmt.Errorf("Cannot parse HASH_TOTAL_WORKERS '%v': %w", rawTotalWorkers, err)
		}
		if totalWorkers < 1 {
			return fmt.Errorf("Invalid HASH_TOTAL_WORKERS value '%v', should be greater than 0", totalWorkers)
		}
		c.totalWorkers = uint(totalWorkers)
	}

	c.queueSize = def.QueueSize()
	rawQueueSize, ok := os.LookupEnv("HASH_QUEUE_SIZE")
	if ok {
		queueSize, err := strconv.ParseUint(rawQueueSize, 10, 64)
		if err != nil {
			return fmt.Errorf("Cannot parse HASH_QUEUE_SIZE '%v': %w", rawQueueSize, err)
		}
		if c.queueSize <= 0 {
			return fmt.Errorf("Invalid HASH_QUEUE_SIZE value '%v', should be greater than 0", queueSize)
		}
		c.queueSize = uint(queueSize)
	}

	return nil
}

func (c *config) parseFlags() error {
	var taskDelay uint
	var totalWorkers uint
	var queueSize uint

	flag.UintVar(&taskDelay, "delay", uint(c.TaskDelay().Seconds()), "number of seconds to delay hash task")
	flag.UintVar(&totalWorkers, "workers", uint(c.TotalWorkers()), "number of workers in the hash pool")
	flag.UintVar(&queueSize, "queue-size", uint(c.QueueSize()), "hash pool queue size")
	flag.Parse()

	c.taskDelay = time.Duration(taskDelay) * time.Second

	if totalWorkers == 0 {
		return fmt.Errorf("Number of workers should be greater than 0")
	} else if totalWorkers > 0 {
		c.totalWorkers = totalWorkers
	}

	if queueSize == 0 {
		return fmt.Errorf("Queue size should be greater than 0")
	} else if queueSize > 0 {
		c.queueSize = queueSize
	}

	return nil
}

func (c *config) Address() string {
	return c.addr
}

func (c *config) Port() uint {
	return c.port
}

func (c *config) ListenAddr() string {
	return fmt.Sprintf("%v:%v", c.Address(), c.Port())
}

func (c *config) ReadTimeout() time.Duration {
	return c.readTimeout
}

func (c *config) WriteTimeout() time.Duration {
	return c.writeTimeout
}

func (c *config) IdleTimeout() time.Duration {
	return c.idleTimeout
}

func (c *config) ShutDownTimeout() time.Duration {
	return c.shutDownTimeout
}

func (c *config) TaskDelay() time.Duration {
	return c.taskDelay
}

func (c *config) TotalWorkers() uint {
	return c.totalWorkers
}

func (c *config) QueueSize() uint {
	return c.queueSize
}
