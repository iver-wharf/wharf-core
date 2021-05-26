// Package consolejson is a concrete implementation of the logger.Sink and
// logger.Context used for outputting JSON-formatted log lines.
//
// Its speed and memory footprint is what gains it its edge over the alternative
// consolepretty Sink, and is optimally combined with log interpreters such as
// Kibana.
package consolejson
