package config

// NOTE: This interface is in the config package due to a circular dependency if it is in the message package

// Message provides the generic interface to the different notification messages that the Sync Service can handle
type Message interface {
	// RequiresReload lets the Demuxer know that it needs to force a reload of the
	// sync configuration information before processing this message
	// Note: Should only return true once, the first time it is called
	RequireReload() bool

	// Match returns true if the message originates from the given core service
	// Returns (is_match, should_ignore) where should_ignore causes the demuxer to drop the message
	Match(populatedConfig *PopulatedCoreServiceConfiguration) (bool, bool)

	// Execute causes the message to be processed and any changes specified by it to be
	// implemented. The config given here is the config that was given to Match() and
	// resulted in a true response.
	Execute(populatedConfig *PopulatedCoreServiceConfiguration)

	// String provides a string representation of the message, used for log messages
	String() string
}
