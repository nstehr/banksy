package banksy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Server webhook server that listens for incoming PRs
type Server struct {
	labeller  *Labeller
	port      int
	hookToken string
	github    *github.Client
}

// NewServer instantiates a new Server struct
func NewServer(l *Labeller, port int, hookToken string, apiToken string, baseURL string) (*Server, error) {
	if hookToken == "" {
		return nil, errors.New("Token used to secure webhook must be specified")
	}
	if apiToken == "" {
		return nil, errors.New("Token used to make API calls must be specified")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	if len(baseURL) > 0 {
		githubURL, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		client.BaseURL = githubURL
		client.UploadURL = githubURL
	}

	return &Server{labeller: l, port: port, hookToken: hookToken, github: client}, nil
}

func (s *Server) webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Only POST methods accepted")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	payload, err := github.ValidatePayload(r, []byte(s.hookToken))
	if err != nil {
		log.Println("webhook token not valid")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	switch e := event.(type) {
	// right now only support pull request events
	case *github.PullRequestEvent:
		// only apply to certain actions
		if e.GetAction() != "opened" && e.GetAction() != "updated" && e.GetAction() != "reopened" {
			return
		}
		// kick off the PR labelling process, do it in a goroutine not to slow the webhook response
		go func() {
			ctx := context.Background()
			number := e.GetNumber()
			repo := e.GetRepo().GetName()
			owner := e.GetRepo().GetOwner().GetLogin()

			pr, _, err := s.github.PullRequests.Get(ctx, owner, repo, number)
			if err != nil {
				log.Println("Error retrieving PR: ", err)
			}
			ctx = context.Background()
			files, _, err := s.github.PullRequests.ListFiles(ctx, owner, repo, number, nil)
			if err != nil {
				log.Println("Error retrieving PR files: ", err)
			}
			labels := s.labeller.determineLabelling(pr, files)
			log.Printf("Applying the following labels: %v\n", labels)
			if len(labels) > 0 {
				labelCtx := context.Background()
				_, _, err := s.github.Issues.ReplaceLabelsForIssue(labelCtx, owner, repo, number, labels)
				if err != nil {
					log.Println("Error labelling issue", err)
				}
			}
		}()
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}

// Start starts listening for incoming webhooks
func (s *Server) Start() {
	http.HandleFunc("/banksy/hook", s.webhookHandler)
	log.Printf("Listening on :%d\n", s.port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)

}
