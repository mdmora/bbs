package gui

import (
	"fmt"
	"net"
	"net/http"
	"path/filepath"
)

var (
	listener net.Listener
	quit     chan struct{}
)

func OpenWebInterface(g *Gateway) (string, error) {
	// Get host.
	host := fmt.Sprintf("127.0.0.1:%d", g.config.WebGUIPort())

	quit = make(chan struct{})
	appLoc, e := filepath.Abs(g.config.WebGUIDir())
	if e != nil {
		return "", e
	}

	listener, e = net.Listen("tcp", host)
	if e != nil {
		return "", e
	}
	go serve(listener, NewServeMux(g, appLoc), quit)
	return fmt.Sprintf("%s://%s", "http", host), nil
}

func serve(listener net.Listener, mux *http.ServeMux, q chan struct{}) {
	for {
		if e := http.Serve(listener, mux); e != nil {
			select {
			case <-q:
				return
			default:
			}
			continue
		}
	}
}

// Close closes the http service.
func Close() {
	if quit != nil {
		// must close quit first
		close(quit)
		listener.Close()
		listener = nil
	}
}

// NewServeMux creates a http.ServeMux with handlers registered.
func NewServeMux(g *Gateway, appLoc string) *http.ServeMux {
	// Register objects.
	api := NewAPI(g)

	// Prepare mux.
	mux := http.NewServeMux()

	if g.config.WebGUIEnable() {
		mux.Handle("/", http.FileServer(http.Dir(appLoc)))
	}

	mux.HandleFunc("/api/get_stats", api.GetStats)
	mux.HandleFunc("/api/connections/get_all", api.GetConnections)
	mux.HandleFunc("/api/connections/new", api.AddConnection)
	mux.HandleFunc("/api/connections/remove", api.RemoveConnection)

	mux.HandleFunc("/api/get_subscription", api.GetSubscription)
	mux.HandleFunc("/api/get_subscriptions", api.GetSubscriptions)
	mux.HandleFunc("/api/subscribe", api.Subscribe)
	mux.HandleFunc("/api/unsubscribe", api.Unsubscribe)

	mux.HandleFunc("/api/users/get_current", api.GetCurrentUser)
	mux.HandleFunc("/api/users/set_current", api.SetCurrentUser)
	mux.HandleFunc("/api/users/get_masters", api.GetMasterUsers)
	mux.HandleFunc("/api/users/new_master", api.NewMasterUser)
	mux.HandleFunc("/api/users/get_all", api.GetUsers)
	mux.HandleFunc("/api/users/new", api.NewUser)
	mux.HandleFunc("/api/users/remove", api.RemoveUser)

	mux.HandleFunc("/api/get_boards", api.GetBoards)
	mux.HandleFunc("/api/new_board", api.NewBoard)
	mux.HandleFunc("/api/get_threads", api.GetThreads)
	mux.HandleFunc("/api/new_thread", api.NewThread)
	mux.HandleFunc("/api/get_threadpage", api.GetThreadPage)
	mux.HandleFunc("/api/get_posts", api.GetPosts)
	mux.HandleFunc("/api/new_post", api.NewPost)
	mux.HandleFunc("/api/import_thread", api.ImportThread)
	mux.HandleFunc("/api/remove_board", api.RemoveBoard)
	mux.HandleFunc("/api/remove_thread", api.RemoveThread)
	mux.HandleFunc("/api/remove_post", api.RemovePost)

	mux.HandleFunc("/api/get_thread_votes", api.GetVotesForThread)
	mux.HandleFunc("/api/get_post_votes", api.GetVotesForPost)
	mux.HandleFunc("/api/add_thread_vote", api.AddVoteForThread)
	mux.HandleFunc("/api/add_post_vote", api.AddVoteForPost)

	mux.HandleFunc("/api/hex/get_threadpage", api.GetThreadPageAsHex)
	mux.HandleFunc("/api/hex/get_threadpage/tp_ref", api.GetThreadPageWithTpRefAsHex)
	mux.HandleFunc("/api/hex/new_thread", api.NewThreadWithHex)
	mux.HandleFunc("/api/hex/new_post", api.NewPostWithHex)

	mux.HandleFunc("/api/tests/new_filled_board", api.TestNewFilledBoard)

	return mux
}
