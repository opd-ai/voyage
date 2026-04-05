package leaderboard

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

// Client errors.
var (
	ErrServerUnavailable = errors.New("leaderboard server unavailable")
	ErrSubmissionFailed  = errors.New("failed to submit entry")
	ErrQueryFailed       = errors.New("failed to query leaderboard")
	ErrInvalidResponse   = errors.New("invalid server response")
)

// ClientConfig holds configuration for the leaderboard client.
type ClientConfig struct {
	// ServerURL is the base URL of the leaderboard server.
	ServerURL string

	// Timeout for HTTP requests.
	Timeout time.Duration

	// RetryCount is the number of retries for failed requests.
	RetryCount int

	// RetryDelay is the delay between retries.
	RetryDelay time.Duration

	// LocalStorage for offline caching.
	LocalStorage *LocalStorage
}

// DefaultConfig returns a default client configuration.
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		ServerURL:  "https://api.voyage-game.example.com/leaderboard",
		Timeout:    10 * time.Second,
		RetryCount: 3,
		RetryDelay: 1 * time.Second,
	}
}

// Client provides leaderboard server communication.
type Client struct {
	config     *ClientConfig
	httpClient *http.Client
	mu         sync.RWMutex
	online     bool
}

// NewClient creates a new leaderboard client.
func NewClient(config *ClientConfig) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		online: false,
	}
}

// Submit submits an entry to the leaderboard server.
func (c *Client) Submit(entry *Entry) error {
	if err := entry.Validate(); err != nil {
		return err
	}

	if err := c.storeLocally(entry); err != nil {
		return err
	}

	return c.attemptServerSubmit(entry)
}

// storeLocally saves the entry to local storage if available.
func (c *Client) storeLocally(entry *Entry) error {
	if c.config.LocalStorage == nil {
		return nil
	}
	return c.config.LocalStorage.AddEntry(entry)
}

// attemptServerSubmit tries to submit entry to server, handling success/failure.
func (c *Client) attemptServerSubmit(entry *Entry) error {
	data, err := entry.Marshal()
	if err != nil {
		return err
	}

	if err := c.submitToServer(data); err != nil {
		c.setOnline(false)
		return nil // Local storage succeeded, server submission deferred
	}

	c.setOnline(true)
	c.clearPendingEntry(entry)
	return nil
}

// clearPendingEntry removes entry from pending list after successful server submit.
func (c *Client) clearPendingEntry(entry *Entry) {
	if c.config.LocalStorage != nil {
		c.config.LocalStorage.RemovePendingEntry(entry)
	}
}

// submitToServer sends entry data to the server with retries.
func (c *Client) submitToServer(data []byte) error {
	var lastErr error
	for i := 0; i <= c.config.RetryCount; i++ {
		if i > 0 {
			time.Sleep(c.config.RetryDelay)
		}

		req, err := http.NewRequest(
			http.MethodPost,
			c.config.ServerURL+"/submit",
			bytes.NewReader(data),
		)
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			return nil
		}
		lastErr = ErrSubmissionFailed
	}

	return lastErr
}

// Query retrieves leaderboard entries from the server.
func (c *Client) Query(opts QueryOptions) (*Board, error) {
	// Build query URL
	queryURL, err := c.buildQueryURL(opts)
	if err != nil {
		return nil, err
	}

	// Attempt server query
	board, err := c.queryServer(queryURL)
	if err != nil {
		c.setOnline(false)
		// Fall back to local storage
		if c.config.LocalStorage != nil {
			return c.queryLocal(opts), nil
		}
		return nil, ErrServerUnavailable
	}

	c.setOnline(true)
	return board, nil
}

// QueryOptions specifies query filters.
type QueryOptions struct {
	Seed   *int64          // Filter by specific seed
	Genre  *engine.GenreID // Filter by genre
	Limit  int             // Maximum entries to return
	Offset int             // Pagination offset
}

// buildQueryURL constructs the query URL with parameters.
func (c *Client) buildQueryURL(opts QueryOptions) (string, error) {
	base, err := url.Parse(c.config.ServerURL + "/query")
	if err != nil {
		return "", err
	}

	params := url.Values{}
	if opts.Seed != nil {
		params.Set("seed", strconv.FormatInt(*opts.Seed, 10))
	}
	if opts.Genre != nil {
		params.Set("genre", string(*opts.Genre))
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

// queryServer fetches entries from the server.
func (c *Client) queryServer(queryURL string) (*Board, error) {
	resp, err := c.httpClient.Get(queryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrQueryFailed
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return Unmarshal(data)
}

// queryLocal queries the local storage.
func (c *Client) queryLocal(opts QueryOptions) *Board {
	local := c.config.LocalStorage.GetBoard()
	entries := filterLocalEntries(local, opts)
	return buildBoardFromSlice(entries, opts)
}

// filterLocalEntries returns entries from local board matching the query options.
func filterLocalEntries(local *Board, opts QueryOptions) []*Entry {
	if opts.Seed != nil && opts.Genre != nil {
		return local.GetBySeedAndGenre(*opts.Seed, *opts.Genre)
	}
	if opts.Seed != nil {
		return local.GetBySeed(*opts.Seed)
	}
	if opts.Genre != nil {
		return local.GetByGenre(*opts.Genre)
	}
	return local.GetAll()
}

// buildBoardFromSlice creates a Board from a slice applying offset and limit.
func buildBoardFromSlice(entries []*Entry, opts QueryOptions) *Board {
	start, end := calculatePaginationBounds(len(entries), opts.Offset, opts.Limit)
	result := NewBoard()
	for _, e := range entries[start:end] {
		_ = result.Add(e)
	}
	return result
}

// calculatePaginationBounds computes start and end indices for pagination.
func calculatePaginationBounds(total, offset, limit int) (start, end int) {
	start = offset
	if start > total {
		start = total
	}
	end = total
	if limit > 0 && start+limit < end {
		end = start + limit
	}
	return start, end
}

// SyncPending attempts to sync pending entries to the server.
func (c *Client) SyncPending() (int, error) {
	if c.config.LocalStorage == nil {
		return 0, nil
	}

	pending := c.config.LocalStorage.GetPendingEntries()
	synced := 0

	for _, entry := range pending {
		data, err := entry.Marshal()
		if err != nil {
			continue
		}

		if err := c.submitToServer(data); err != nil {
			// Stop on first failure - server is likely unavailable
			return synced, err
		}

		c.config.LocalStorage.RemovePendingEntry(entry)
		synced++
	}

	c.setOnline(true)
	return synced, nil
}

// FetchGlobal fetches the global leaderboard and merges into local storage.
func (c *Client) FetchGlobal(limit int) error {
	opts := QueryOptions{Limit: limit}
	board, err := c.queryServer(func() string {
		url, _ := c.buildQueryURL(opts)
		return url
	}())
	if err != nil {
		return err
	}

	if c.config.LocalStorage != nil {
		c.config.LocalStorage.MergeBoard(board)
	}

	return nil
}

// GetReplayableSeed returns a top-scoring seed for replay.
func (c *Client) GetReplayableSeed(genre engine.GenreID) (int64, *Entry, error) {
	g := genre
	opts := QueryOptions{Genre: &g, Limit: 1}
	board, err := c.Query(opts)
	if err != nil {
		return 0, nil, err
	}

	seed, entry := board.GetReplayableSeed(genre)
	return seed, entry, nil
}

// IsOnline returns whether the client has server connectivity.
func (c *Client) IsOnline() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.online
}

// setOnline updates the online status.
func (c *Client) setOnline(online bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.online = online
}

// CheckConnectivity tests server availability.
func (c *Client) CheckConnectivity() bool {
	resp, err := c.httpClient.Get(c.config.ServerURL + "/health")
	if err != nil {
		c.setOnline(false)
		return false
	}
	resp.Body.Close()
	c.setOnline(resp.StatusCode == http.StatusOK)
	return c.online
}

// SubmitResponse represents the server response to a submission.
type SubmitResponse struct {
	Success bool   `json:"success"`
	Rank    int    `json:"rank,omitempty"`
	Message string `json:"message,omitempty"`
}

// ParseSubmitResponse parses the server response.
func ParseSubmitResponse(data []byte) (*SubmitResponse, error) {
	var resp SubmitResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
