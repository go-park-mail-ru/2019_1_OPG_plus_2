// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjson85f0d656Decode20191OPGPlus2InternalPkgModels(in *jlexer.Lexer, out *RoomsOnlineMessage) {
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
		case "rooms_online":
			if in.IsNull() {
				in.Skip()
				out.RoomsOnline = nil
			} else {
				in.Delim('[')
				if out.RoomsOnline == nil {
					if !in.IsDelim(']') {
						out.RoomsOnline = make([]RoomData, 0, 1)
					} else {
						out.RoomsOnline = []RoomData{}
					}
				} else {
					out.RoomsOnline = (out.RoomsOnline)[:0]
				}
				for !in.IsDelim(']') {
					var v1 RoomData
					(v1).UnmarshalEasyJSON(in)
					out.RoomsOnline = append(out.RoomsOnline, v1)
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
func easyjson85f0d656Encode20191OPGPlus2InternalPkgModels(out *jwriter.Writer, in RoomsOnlineMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"rooms_online\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.RoomsOnline == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.RoomsOnline {
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

// MarshalJSON supports json.Marshaler interface
func (v RoomsOnlineMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656Encode20191OPGPlus2InternalPkgModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RoomsOnlineMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656Encode20191OPGPlus2InternalPkgModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RoomsOnlineMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656Decode20191OPGPlus2InternalPkgModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RoomsOnlineMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656Decode20191OPGPlus2InternalPkgModels(l, v)
}
func easyjson85f0d656Decode20191OPGPlus2InternalPkgModels1(in *jlexer.Lexer, out *RoomPlayer) {
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
		case "username":
			out.Username = string(in.String())
		case "avatar":
			out.Avatar = string(in.String())
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
func easyjson85f0d656Encode20191OPGPlus2InternalPkgModels1(out *jwriter.Writer, in RoomPlayer) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"username\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"avatar\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Avatar))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RoomPlayer) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656Encode20191OPGPlus2InternalPkgModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RoomPlayer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656Encode20191OPGPlus2InternalPkgModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RoomPlayer) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656Decode20191OPGPlus2InternalPkgModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RoomPlayer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656Decode20191OPGPlus2InternalPkgModels1(l, v)
}
func easyjson85f0d656Decode20191OPGPlus2InternalPkgModels2(in *jlexer.Lexer, out *RoomDeletedMessage) {
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
		case "room_id":
			out.RoomId = string(in.String())
		case "status":
			out.Status = int(in.Int())
		case "message":
			out.Message = string(in.String())
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
func easyjson85f0d656Encode20191OPGPlus2InternalPkgModels2(out *jwriter.Writer, in RoomDeletedMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"room_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.RoomId))
	}
	{
		const prefix string = ",\"status\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Status))
	}
	{
		const prefix string = ",\"message\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RoomDeletedMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656Encode20191OPGPlus2InternalPkgModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RoomDeletedMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656Encode20191OPGPlus2InternalPkgModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RoomDeletedMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656Decode20191OPGPlus2InternalPkgModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RoomDeletedMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656Decode20191OPGPlus2InternalPkgModels2(l, v)
}
func easyjson85f0d656Decode20191OPGPlus2InternalPkgModels3(in *jlexer.Lexer, out *RoomData) {
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
		case "id":
			out.Id = string(in.String())
		case "players_num":
			out.PlayersNum = int(in.Int())
		case "players":
			if in.IsNull() {
				in.Skip()
				out.Players = nil
			} else {
				in.Delim('[')
				if out.Players == nil {
					if !in.IsDelim(']') {
						out.Players = make([]RoomPlayer, 0, 2)
					} else {
						out.Players = []RoomPlayer{}
					}
				} else {
					out.Players = (out.Players)[:0]
				}
				for !in.IsDelim(']') {
					var v4 RoomPlayer
					(v4).UnmarshalEasyJSON(in)
					out.Players = append(out.Players, v4)
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
func easyjson85f0d656Encode20191OPGPlus2InternalPkgModels3(out *jwriter.Writer, in RoomData) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"players_num\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.PlayersNum))
	}
	{
		const prefix string = ",\"players\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Players == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Players {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RoomData) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656Encode20191OPGPlus2InternalPkgModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RoomData) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656Encode20191OPGPlus2InternalPkgModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RoomData) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656Decode20191OPGPlus2InternalPkgModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RoomData) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656Decode20191OPGPlus2InternalPkgModels3(l, v)
}
