package consumer

func WithUpdater(u Updater) Options {
	return func(c *Consumer) {
		c.updater = u
	}
}
