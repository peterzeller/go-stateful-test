package smallcheck

type Config struct {
	Depth int
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
