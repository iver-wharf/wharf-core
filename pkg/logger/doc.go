// Package logger contains logging types and functions in a memory-efficient and
// fast manner.
//
// This package contains abstractions for the different layers of logging:
// logger.Logger, logger.Event, logger.Sink, and logger.Context.
//
// The Logger interface is what you use to create new log events.
//
// The Event interface is a single log event. It contains methods to populate it
// with data in a type-safe way. A single Event can contain any number of
// log Context used to construct the different log messages.
//
// The Sink interface is used to create a new concrete log Context specific to
// that Sink. Examples are console sinks that either produce JSON- or
// pretty-formatted logs.
//
// The Context interface is created by the Sink and a concrete implementation
// holds data about the log Event for that particular Sink. Things like
// formatted fields, ready to be written out to STDOUT.
package logger
