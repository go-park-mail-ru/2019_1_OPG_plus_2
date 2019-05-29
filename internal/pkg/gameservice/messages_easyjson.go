// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package gameservice

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice(in *jlexer.Lexer, out *Point) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "x":
			out.X = int(in.Int())
		case "y":
			out.Y = int(in.Int())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice(out *jwriter.Writer, in Point) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"x\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.X))
	}
	{
		const prefix string = ",\"y\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Y))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Point) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Point) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Point) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Point) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice(l, v)
}
func easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice1(in *jlexer.Lexer, out *GenericMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user":
			out.User = string(in.String())
		case "type":
			out.MType = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice1(out *jwriter.Writer, in GenericMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.User))
	}
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.MType))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GenericMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GenericMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GenericMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GenericMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice1(l, v)
}
func easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice2(in *jlexer.Lexer, out *GameMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "data":
			easyjson66c1e240Decode(in, &out.Data)
		case "user":
			out.User = string(in.String())
		case "type":
			out.MType = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice2(out *jwriter.Writer, in GameMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjson66c1e240Encode(out, in.Data)
	}
	{
		const prefix string = ",\"user\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.User))
	}
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.MType))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v GameMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v GameMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *GameMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *GameMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice2(l, v)
}
func easyjson66c1e240Decode(in *jlexer.Lexer, out *struct {
	Coords []Point `json:"coords"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "coords":
			if in.IsNull() {
				in.Skip()
				out.Coords = nil
			} else {
				in.Delim('[')
				if out.Coords == nil {
					if !in.IsDelim(']') {
						out.Coords = make([]Point, 0, 4)
					} else {
						out.Coords = []Point{}
					}
				} else {
					out.Coords = (out.Coords)[:0]
				}
				for !in.IsDelim(']') {
					var v1 Point
					(v1).UnmarshalEasyJSON(in)
					out.Coords = append(out.Coords, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson66c1e240Encode(out *jwriter.Writer, in struct {
	Coords []Point `json:"coords"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"coords\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Coords == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Coords {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}
func easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice3(in *jlexer.Lexer, out *EventData) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "event_type":
			out.EventType = string(in.String())
		case "event_data":
			if m, ok := out.EventData.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.EventData.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.EventData = in.Interface()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice3(out *jwriter.Writer, in EventData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"event_type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.EventType))
	}
	{
		const prefix string = ",\"event_data\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if m, ok := in.EventData.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.EventData.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.EventData))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice3(l, v)
}
func easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice4(in *jlexer.Lexer, out *ChatMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "data":
			easyjson66c1e240Decode1(in, &out.Data)
		case "user":
			out.User = string(in.String())
		case "type":
			out.MType = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice4(out *jwriter.Writer, in ChatMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjson66c1e240Encode1(out, in.Data)
	}
	{
		const prefix string = ",\"user\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.User))
	}
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.MType))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ChatMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ChatMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ChatMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ChatMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice4(l, v)
}
func easyjson66c1e240Decode1(in *jlexer.Lexer, out *struct {
	Text string `json:"text, string"`
}) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "text":
			out.Text = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson66c1e240Encode1(out *jwriter.Writer, in struct {
	Text string `json:"text, string"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"text\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Text))
	}
	out.RawByte('}')
}
func easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice5(in *jlexer.Lexer, out *BroadcastEventMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "data":
			(out.Data).UnmarshalEasyJSON(in)
		case "user":
			out.User = string(in.String())
		case "type":
			out.MType = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice5(out *jwriter.Writer, in BroadcastEventMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"data\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(in.Data).MarshalEasyJSON(out)
	}
	{
		const prefix string = ",\"user\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.User))
	}
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.MType))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v BroadcastEventMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v BroadcastEventMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson66c1e240Encode20191OPGPlus2InternalPkgGameservice5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *BroadcastEventMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *BroadcastEventMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson66c1e240Decode20191OPGPlus2InternalPkgGameservice5(l, v)
}
