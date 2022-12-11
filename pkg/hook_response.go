package registry_ops

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
type HookResponse struct {
	PushData   pushData `json:"push_data"`
	Repository repository
}
