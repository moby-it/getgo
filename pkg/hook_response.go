package docker_ops

type pushData struct {
	Pusher    string
	Pushed_At int
	Tag       string
}
type repository struct {
	Namespace string
	Name      string
	RepoName  string `json:"repo_name"`
}

// / Stripped Webhook JSON post data of Dockerhub. See https://docs.docker.com/docker-hub/webhooks for more.
type HookResponse struct {
	PushData   pushData `json:"push_data"`
	Repository repository
}
