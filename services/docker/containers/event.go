package containers

import (
	"context"
	"log"

	dc "github.com/Szent7/medovukha-core/services/docker"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

func EventStream(ctx context.Context, cli dc.IDockerClient, eventsCh chan<- events.Message, errCh chan<- error) {
	filters := filters.NewArgs()
	filters.Add("type", "container")

	msgs, errs := cli.Events(ctx, events.ListOptions{Filters: filters})

	// fmt.Println("Start EventStream")

	for {
		select {
		case <-ctx.Done():
			{
				// fmt.Println("Stop EventStream")
				return
			}
		case err := <-errs:
			{
				if err != nil {
					log.Printf("DockerEventStream error: %s\n", err.Error())
					errCh <- err
				}
			}
		case msg := <-msgs:
			{
				// fmt.Println("New message in EventStream")
				eventsCh <- msg
			}
		}
	}
}
