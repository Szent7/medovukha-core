package containers

import (
	"context"
	"log"
	dc "medovukha/services/docker"
	"time"

	"github.com/docker/docker/api/types/events"
)

type (
	StreamState int
	StreamEvent int
)

const (
	Stopped StreamState = iota
	Starting
	Running
	Stopping
)

const (
	ClientConnected StreamEvent = iota
	ClientDisconnected
	StreamExited
)

type DockerStreamManager struct {
	cli dc.IDockerClient

	eventCh chan<- events.Message

	state   StreamState
	clients int

	ctx    context.Context
	cancel context.CancelFunc

	input chan StreamEvent
}

func (m *DockerStreamManager) OnClientConnected() {
	m.input <- ClientConnected
}

func (m *DockerStreamManager) OnClientDisconnected() {
	m.input <- ClientDisconnected
}

func (m *DockerStreamManager) loop() {
	for event := range m.input {
		switch m.state {

		case Stopped:
			{
				if event == ClientConnected {
					m.clients++
					m.start()
				}
			}

		case Starting:
			{
				if event == ClientConnected {
					m.clients++
				}
				if event == ClientDisconnected {
					m.clients--
				}
			}

		case Running:
			{
				switch event {
				case ClientConnected:
					{
						m.clients++
					}
				case ClientDisconnected:
					{
						m.clients--
						if m.clients == 0 {
							time.AfterFunc(5*time.Second, func() {
								m.input <- ClientDisconnected
							})
						}

					}
				case StreamExited:
					{
						m.state = Stopped
					}
				}
			}

		case Stopping:
			{
				if event == StreamExited {
					m.state = Stopped
					if m.clients > 0 {
						m.start()
					}
				}
			}
		}
	}
}

func (m *DockerStreamManager) start() {
	m.state = Starting
	m.ctx, m.cancel = context.WithCancel(context.Background())

	go func() {
		m.state = Running
		err := dockerEventStream(m.ctx, m.cli, m.eventCh)
		if err != nil {
			log.Printf("DockerStreamManager error: %s\n", err.Error())
		}
		m.input <- StreamExited
	}()
}

func (m *DockerStreamManager) stop() {
	m.state = Stopping
	m.cancel()
}

func NewDockerStreamManager(cli dc.IDockerClient, eventCh chan<- events.Message) *DockerStreamManager {
	manager := &DockerStreamManager{
		cli:     cli,
		eventCh: eventCh,
		state:   Stopped,
		input:   make(chan StreamEvent, 10),
	}

	go manager.loop()
	return manager
}

func dockerEventStream(ctx context.Context, cli dc.IDockerClient, eventsCh chan<- events.Message) error {
	msgs, errs := cli.Events(ctx, events.ListOptions{})

	for {
		select {
		case <-ctx.Done():
			{
				return nil
			}
		case err := <-errs:
			{
				if err != nil {
					log.Printf("DockerEventStream error: %s\n", err.Error())
					return err
				}
			}
		case msg := <-msgs:
			{
				eventsCh <- msg
			}
		}
	}
}
