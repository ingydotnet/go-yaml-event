package yaml

import (
	"fmt"
	"io"
)

// EventType represents the type of a YAML parser event
type EventType int

const (
	EventNone EventType = iota
	EventStreamStart
	EventStreamEnd
	EventDocumentStart
	EventDocumentEnd
	EventAlias
	EventScalar
	EventSequenceStart
	EventSequenceEnd
	EventMappingStart
	EventMappingEnd
)

func (e EventType) String() string {
	switch e {
	case EventStreamStart:
		return "STREAM-START"
	case EventStreamEnd:
		return "STREAM-END"
	case EventDocumentStart:
		return "DOCUMENT-START"
	case EventDocumentEnd:
		return "DOCUMENT-END"
	case EventAlias:
		return "ALIAS"
	case EventScalar:
		return "SCALAR"
	case EventSequenceStart:
		return "SEQUENCE-START"
	case EventSequenceEnd:
		return "SEQUENCE-END"
	case EventMappingStart:
		return "MAPPING-START"
	case EventMappingEnd:
		return "MAPPING-END"
	default:
		return "NONE"
	}
}

// Event represents a YAML parser event
type Event struct {
	Type       EventType
	Value      string
	Anchor     string
	Tag        string
	Style      yaml_style_t
	Implicit   bool
	StartMark  Mark
	EndMark    Mark
	HeadComment []byte
	LineComment []byte
	FootComment []byte
	TailComment []byte
}

// StyleString returns a human-readable representation of the style
func (e *Event) StyleString() string {
	switch e.Type {
	case EventScalar:
		switch yaml_scalar_style_t(e.Style) {
		case yaml_PLAIN_SCALAR_STYLE:
			return "plain"
		case yaml_SINGLE_QUOTED_SCALAR_STYLE:
			return "single-quoted"
		case yaml_DOUBLE_QUOTED_SCALAR_STYLE:
			return "double-quoted"
		case yaml_LITERAL_SCALAR_STYLE:
			return "literal"
		case yaml_FOLDED_SCALAR_STYLE:
			return "folded"
		default:
			return "any"
		}
	case EventSequenceStart:
		switch yaml_sequence_style_t(e.Style) {
		case yaml_FLOW_SEQUENCE_STYLE:
			return "flow"
		case yaml_BLOCK_SEQUENCE_STYLE:
			return "block"
		default:
			return "any"
		}
	case EventMappingStart:
		switch yaml_mapping_style_t(e.Style) {
		case yaml_FLOW_MAPPING_STYLE:
			return "flow"
		case yaml_BLOCK_MAPPING_STYLE:
			return "block"
		default:
			return "any"
		}
	default:
		return ""
	}
}

// Mark represents a position in the YAML input stream
type Mark struct {
	Index  int
	Line   int
	Column int
}

// Parser provides a high-level interface for parsing YAML streams
type Parser struct {
	parser yaml_parser_t
	done   bool
}

// NewParser creates a new YAML parser reading from the given reader
func NewParser(reader io.Reader) (*Parser, error) {
	var p Parser
	if !yaml_parser_initialize(&p.parser) {
		return nil, fmt.Errorf("failed to initialize YAML parser")
	}
	yaml_parser_set_input_reader(&p.parser, reader)
	return &p, nil
}

// Next returns the next event in the YAML stream
func (p *Parser) Next() (*Event, error) {
	if p.done {
		return nil, nil
	}

	var yamlEvent yaml_event_t
	if !yaml_parser_parse(&p.parser, &yamlEvent) {
		if p.parser.error != yaml_NO_ERROR {
			return nil, fmt.Errorf("parser error: %v", p.parser.problem)
		}
		p.done = true
		return nil, nil
	}

	event := &Event{
		StartMark: Mark{
			Index:  int(yamlEvent.start_mark.index),
			Line:   int(yamlEvent.start_mark.line),
			Column: int(yamlEvent.start_mark.column),
		},
		EndMark: Mark{
			Index:  int(yamlEvent.end_mark.index),
			Line:   int(yamlEvent.end_mark.line),
			Column: int(yamlEvent.end_mark.column),
		},
		HeadComment: yamlEvent.head_comment,
		LineComment: yamlEvent.line_comment,
		FootComment: yamlEvent.foot_comment,
		TailComment: yamlEvent.tail_comment,
	}

	switch yamlEvent.typ {
	case yaml_STREAM_START_EVENT:
		event.Type = EventStreamStart
	case yaml_STREAM_END_EVENT:
		event.Type = EventStreamEnd
		p.done = true
	case yaml_DOCUMENT_START_EVENT:
		event.Type = EventDocumentStart
		event.Implicit = yamlEvent.implicit
	case yaml_DOCUMENT_END_EVENT:
		event.Type = EventDocumentEnd
		event.Implicit = yamlEvent.implicit
	case yaml_ALIAS_EVENT:
		event.Type = EventAlias
		event.Anchor = string(yamlEvent.anchor)
	case yaml_SCALAR_EVENT:
		event.Type = EventScalar
		event.Value = string(yamlEvent.value)
		event.Anchor = string(yamlEvent.anchor)
		event.Tag = string(yamlEvent.tag)
		event.Implicit = yamlEvent.implicit
		event.Style = yaml_style_t(yamlEvent.scalar_style())
	case yaml_SEQUENCE_START_EVENT:
		event.Type = EventSequenceStart
		event.Anchor = string(yamlEvent.anchor)
		event.Tag = string(yamlEvent.tag)
		event.Implicit = yamlEvent.implicit
		event.Style = yaml_style_t(yamlEvent.sequence_style())
	case yaml_SEQUENCE_END_EVENT:
		event.Type = EventSequenceEnd
	case yaml_MAPPING_START_EVENT:
		event.Type = EventMappingStart
		event.Anchor = string(yamlEvent.anchor)
		event.Tag = string(yamlEvent.tag)
		event.Implicit = yamlEvent.implicit
		event.Style = yaml_style_t(yamlEvent.mapping_style())
	case yaml_MAPPING_END_EVENT:
		event.Type = EventMappingEnd
	}

	yaml_event_delete(&yamlEvent)
	return event, nil
}

// Close releases the parser resources
func (p *Parser) Close() {
	yaml_parser_delete(&p.parser)
}
