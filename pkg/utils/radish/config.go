package radish

// Configures the radish task manager so that different processes can utilize different
// asynchronous task processing resources depending on process compute constraints.
type Config struct {
	Workers    int    `default:"4" desc:"the number of workers to process tasks asynchronously"`
	QueueSize  int    `split_words:"true" default:"64" desc:"the number of async tasks to buffer in the queue before blocking"`
	ServerName string `split_words:"true" default:"radish" desc:"used to describe the radish service in the log"`
}

func (c Config) Validate() error {
	if c.Workers == 0 {
		return ErrNoWorkers
	}

	if c.ServerName == "" {
		return ErrNoServerName
	}

	return nil
}

func (c Config) IsZero() bool {
	return c.Workers == 0 && c.QueueSize == 0 && c.ServerName == ""
}
