package smallcheck

type Config struct {
	// maximum depth to explore
	Depth int
	// print the logs after every run, not just for failing runs
	PrintAllLogs bool
	// print the logs directly, not just after a run
	PrintLiveLogs bool
}

func setDefaults(cfg Config) Config {
	if cfg.Depth == 0 {
		cfg.Depth = 5
	}
	return cfg
}

type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Failed() bool
	Logf(format string, args ...interface{})
}
