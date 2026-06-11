package jobs

type Queue struct {
	channel chan *DeploymentJob
}

func NewQueue() *Queue {
	return &Queue{
		channel: make(chan *DeploymentJob, 100),
	}
}

func (q *Queue) Enqueue(
	job *DeploymentJob,
) {
	q.channel <- job
}

func (q *Queue) Jobs() <-chan *DeploymentJob {
	return q.channel
}
