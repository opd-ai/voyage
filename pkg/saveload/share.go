package saveload

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"io"
	"math/big"
	"strings"

	"github.com/opd-ai/voyage/pkg/engine"
)

// ShareCode errors.
var (
	ErrInvalidShareCode   = errors.New("invalid share code format")
	ErrShareCodeTooShort  = errors.New("share code too short")
	ErrShareCodeCorrupted = errors.New("share code data corrupted")
	ErrInvalidChecksum    = errors.New("share code checksum mismatch")
)

// ShareVersion is the current share code format version.
const ShareVersion uint8 = 1

// DecisionType represents the type of player decision.
type DecisionType uint8

const (
	// DecisionMove represents a movement decision.
	DecisionMove DecisionType = iota
	// DecisionChoice represents an event choice.
	DecisionChoice
	// DecisionRest represents a rest action.
	DecisionRest
	// DecisionForage represents a forage action.
	DecisionForage
	// DecisionTrade represents a trade action.
	DecisionTrade
)

// Decision represents a single player decision during a run.
type Decision struct {
	Turn   int          // Turn number when decision was made
	Type   DecisionType // Type of decision
	Value  int          // Context-dependent value (direction, choice index, etc.)
	Target int          // Optional target (tile ID, NPC ID, etc.)
}

// RunData represents a complete run for sharing.
type RunData struct {
	// Core identification
	Seed  int64          `json:"seed"`
	Genre engine.GenreID `json:"genre"`
	Diffi int            `json:"difficulty"`

	// Decision sequence
	Decisions []Decision `json:"decisions"`

	// Final state (for verification)
	FinalTurn    int  `json:"finalTurn"`
	FinalX       int  `json:"finalX"`
	FinalY       int  `json:"finalY"`
	WonGame      bool `json:"won"`
	CrewSurvived int  `json:"crewSurvived"`
}

// NewRunData creates a new RunData for recording a run.
func NewRunData(seed int64, genre engine.GenreID, difficulty int) *RunData {
	return &RunData{
		Seed:      seed,
		Genre:     genre,
		Diffi:     difficulty,
		Decisions: make([]Decision, 0, 256),
	}
}

// RecordDecision adds a decision to the run data.
func (rd *RunData) RecordDecision(turn int, dtype DecisionType, value, target int) {
	rd.Decisions = append(rd.Decisions, Decision{
		Turn:   turn,
		Type:   dtype,
		Value:  value,
		Target: target,
	})
}

// SetFinalState sets the final game state.
func (rd *RunData) SetFinalState(turn, x, y int, won bool, crewSurvived int) {
	rd.FinalTurn = turn
	rd.FinalX = x
	rd.FinalY = y
	rd.WonGame = won
	rd.CrewSurvived = crewSurvived
}

// Base58 alphabet (Bitcoin-style, no 0OIl to avoid confusion)
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// ExportShareCode encodes the run data as a compact shareable code.
func (rd *RunData) ExportShareCode() (string, error) {
	// Encode to binary
	data, err := rd.encodeBinary()
	if err != nil {
		return "", err
	}

	// Compress
	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	if _, err := gz.Write(data); err != nil {
		return "", err
	}
	gz.Close()

	// Add checksum (simple XOR checksum)
	checksum := calculateChecksum(compressed.Bytes())

	// Prepend version and checksum
	final := make([]byte, 0, 2+compressed.Len())
	final = append(final, ShareVersion)
	final = append(final, checksum)
	final = append(final, compressed.Bytes()...)

	// Encode to base58
	return encodeBase58(final), nil
}

// ImportShareCode decodes a share code back to RunData.
func ImportShareCode(code string) (*RunData, error) {
	data, err := decodeShareCodeData(code)
	if err != nil {
		return nil, err
	}

	if err := validateShareCodeHeader(data); err != nil {
		return nil, err
	}

	decompressed, err := decompressShareData(data[2:])
	if err != nil {
		return nil, err
	}

	return decodeBinary(decompressed)
}

// decodeShareCodeData validates and decodes the base58 share code.
func decodeShareCodeData(code string) ([]byte, error) {
	code = strings.TrimSpace(code)
	if len(code) < 10 {
		return nil, ErrShareCodeTooShort
	}
	return decodeBase58(code)
}

// validateShareCodeHeader checks version and checksum in the decoded data.
func validateShareCodeHeader(data []byte) error {
	if len(data) < 3 {
		return ErrShareCodeCorrupted
	}
	if data[0] != ShareVersion {
		return ErrInvalidShareCode
	}
	if calculateChecksum(data[2:]) != data[1] {
		return ErrInvalidChecksum
	}
	return nil
}

// decompressShareData decompresses gzip-compressed share data.
func decompressShareData(compressedData []byte) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, ErrShareCodeCorrupted
	}
	defer gz.Close()

	decompressed, err := io.ReadAll(gz)
	if err != nil {
		return nil, ErrShareCodeCorrupted
	}
	return decompressed, nil
}

// encodeBinary converts RunData to a compact binary format.
func (rd *RunData) encodeBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write seed (8 bytes)
	if err := binary.Write(buf, binary.LittleEndian, rd.Seed); err != nil {
		return nil, err
	}

	// Write genre (1 byte encoded)
	buf.WriteByte(encodeGenre(rd.Genre))

	// Write difficulty (1 byte)
	buf.WriteByte(byte(rd.Diffi))

	// Write decision count (2 bytes)
	if err := binary.Write(buf, binary.LittleEndian, uint16(len(rd.Decisions))); err != nil {
		return nil, err
	}

	// Write decisions (compact encoding)
	for _, d := range rd.Decisions {
		// Pack turn delta and type into 2 bytes
		// Format: [10 bits turn delta][3 bits type][3 bits value]
		// For larger values, use extended encoding
		if d.Turn < 1024 && d.Value < 8 {
			packed := uint16(d.Turn&0x3FF) | (uint16(d.Type&0x7) << 10) | (uint16(d.Value&0x7) << 13)
			if err := binary.Write(buf, binary.LittleEndian, packed); err != nil {
				return nil, err
			}
		} else {
			// Extended encoding: marker (0xFFFF) + full data (C-005)
			if err := binary.Write(buf, binary.LittleEndian, uint16(0xFFFF)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, uint16(d.Turn)); err != nil {
				return nil, err
			}
			buf.WriteByte(byte(d.Type))
			if err := binary.Write(buf, binary.LittleEndian, int16(d.Value)); err != nil {
				return nil, err
			}
			if err := binary.Write(buf, binary.LittleEndian, int16(d.Target)); err != nil {
				return nil, err
			}
		}
	}

	// Write final state
	if err := binary.Write(buf, binary.LittleEndian, uint16(rd.FinalTurn)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, int16(rd.FinalX)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, int16(rd.FinalY)); err != nil {
		return nil, err
	}
	if rd.WonGame {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}
	buf.WriteByte(byte(rd.CrewSurvived))

	return buf.Bytes(), nil
}

// decodeBinary converts binary data back to RunData.
func decodeBinary(data []byte) (*RunData, error) {
	if len(data) < 14 {
		return nil, ErrShareCodeCorrupted
	}

	buf := bytes.NewReader(data)
	rd := &RunData{}

	if err := decodeRunHeader(buf, rd); err != nil {
		return nil, ErrShareCodeCorrupted
	}
	if err := decodeDecisions(buf, rd); err != nil {
		return nil, ErrShareCodeCorrupted
	}
	if err := decodeFinalState(buf, rd); err != nil {
		return nil, ErrShareCodeCorrupted
	}

	return rd, nil
}

// decodeRunHeader reads the seed, genre, and difficulty from the buffer.
func decodeRunHeader(buf *bytes.Reader, rd *RunData) error {
	if err := binary.Read(buf, binary.LittleEndian, &rd.Seed); err != nil {
		return err
	}

	genreByte, _ := buf.ReadByte()
	rd.Genre = decodeGenre(genreByte)

	diffiByte, _ := buf.ReadByte()
	rd.Diffi = int(diffiByte)
	return nil
}

// decodeDecisions reads all decisions from the buffer.
func decodeDecisions(buf *bytes.Reader, rd *RunData) error {
	var decisionCount uint16
	if err := binary.Read(buf, binary.LittleEndian, &decisionCount); err != nil {
		return err
	}

	rd.Decisions = make([]Decision, 0, decisionCount)
	for i := uint16(0); i < decisionCount; i++ {
		d, err := decodeDecision(buf)
		if err != nil {
			return err
		}
		rd.Decisions = append(rd.Decisions, d)
	}
	return nil
}

// decodeDecision reads a single decision from the buffer.
func decodeDecision(buf *bytes.Reader) (Decision, error) {
	var packed uint16
	if err := binary.Read(buf, binary.LittleEndian, &packed); err != nil {
		return Decision{}, err
	}

	if isExtendedEncoding(packed) {
		return decodeExtendedDecision(buf)
	}
	return decodeCompactDecision(packed), nil
}

// isExtendedEncoding checks if the packed value indicates extended encoding.
// Uses exact marker value 0xFFFF to avoid false positives (C-005).
func isExtendedEncoding(packed uint16) bool {
	return packed == 0xFFFF
}

// decodeExtendedDecision reads a decision using extended encoding format.
// The marker (0xFFFF) has already been read by decodeDecision, so we just
// read the remaining data directly (C-005).
func decodeExtendedDecision(buf *bytes.Reader) (Decision, error) {
	var turn uint16
	var dtype byte
	var value, target int16
	if err := binary.Read(buf, binary.LittleEndian, &turn); err != nil {
		return Decision{}, err
	}
	dtype, _ = buf.ReadByte()
	if err := binary.Read(buf, binary.LittleEndian, &value); err != nil {
		return Decision{}, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &target); err != nil {
		return Decision{}, err
	}

	return Decision{
		Turn:   int(turn),
		Type:   DecisionType(dtype),
		Value:  int(value),
		Target: int(target),
	}, nil
}

// decodeCompactDecision unpacks a decision from compact 16-bit format.
func decodeCompactDecision(packed uint16) Decision {
	return Decision{
		Turn:  int(packed & 0x3FF),
		Type:  DecisionType((packed >> 10) & 0x7),
		Value: int((packed >> 13) & 0x7),
	}
}

// decodeFinalState reads the final game state from the buffer.
func decodeFinalState(buf *bytes.Reader, rd *RunData) error {
	var finalTurn uint16
	var finalX, finalY int16
	if err := binary.Read(buf, binary.LittleEndian, &finalTurn); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &finalX); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &finalY); err != nil {
		return err
	}
	rd.FinalTurn = int(finalTurn)
	rd.FinalX = int(finalX)
	rd.FinalY = int(finalY)

	wonByte, _ := buf.ReadByte()
	rd.WonGame = wonByte == 1

	crewByte, _ := buf.ReadByte()
	rd.CrewSurvived = int(crewByte)
	return nil
}

// encodeGenre converts a genre ID to a byte.
func encodeGenre(g engine.GenreID) byte {
	switch g {
	case engine.GenreFantasy:
		return 0
	case engine.GenreScifi:
		return 1
	case engine.GenreHorror:
		return 2
	case engine.GenreCyberpunk:
		return 3
	case engine.GenrePostapoc:
		return 4
	default:
		return 0
	}
}

// decodeGenre converts a byte back to genre ID.
func decodeGenre(b byte) engine.GenreID {
	switch b {
	case 0:
		return engine.GenreFantasy
	case 1:
		return engine.GenreScifi
	case 2:
		return engine.GenreHorror
	case 3:
		return engine.GenreCyberpunk
	case 4:
		return engine.GenrePostapoc
	default:
		return engine.GenreFantasy
	}
}

// calculateChecksum computes a simple XOR checksum.
func calculateChecksum(data []byte) byte {
	var checksum byte
	for _, b := range data {
		checksum ^= b
	}
	return checksum
}

// encodeBase58 encodes bytes to a base58 string.
func encodeBase58(data []byte) string {
	// Convert to big integer
	num := new(big.Int).SetBytes(data)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)

	var result []byte
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		result = append([]byte{base58Alphabet[mod.Int64()]}, result...)
	}

	// Add leading '1's for leading zero bytes
	for _, b := range data {
		if b == 0 {
			result = append([]byte{'1'}, result...)
		} else {
			break
		}
	}

	return string(result)
}

// decodeBase58 decodes a base58 string to bytes.
func decodeBase58(s string) ([]byte, error) {
	num := big.NewInt(0)
	base := big.NewInt(58)

	for _, c := range s {
		idx := strings.IndexRune(base58Alphabet, c)
		if idx < 0 {
			return nil, ErrInvalidShareCode
		}
		num.Mul(num, base)
		num.Add(num, big.NewInt(int64(idx)))
	}

	result := num.Bytes()

	// Add leading zeros
	for _, c := range s {
		if c == '1' {
			result = append([]byte{0}, result...)
		} else {
			break
		}
	}

	return result, nil
}

// GetDecisionCount returns the number of decisions recorded.
func (rd *RunData) GetDecisionCount() int {
	return len(rd.Decisions)
}

// GetDecisionAt returns the decision at a specific index.
func (rd *RunData) GetDecisionAt(index int) (Decision, bool) {
	if index < 0 || index >= len(rd.Decisions) {
		return Decision{}, false
	}
	return rd.Decisions[index], true
}

// Replay provides an iterator for replaying decisions.
type Replay struct {
	runData *RunData
	index   int
}

// NewReplay creates a new replay iterator from run data.
func NewReplay(rd *RunData) *Replay {
	return &Replay{
		runData: rd,
		index:   0,
	}
}

// HasNext returns true if there are more decisions to replay.
func (r *Replay) HasNext() bool {
	return r.index < len(r.runData.Decisions)
}

// Next returns the next decision and advances the iterator.
func (r *Replay) Next() (Decision, bool) {
	if !r.HasNext() {
		return Decision{}, false
	}
	d := r.runData.Decisions[r.index]
	r.index++
	return d, true
}

// Peek returns the next decision without advancing.
func (r *Replay) Peek() (Decision, bool) {
	if !r.HasNext() {
		return Decision{}, false
	}
	return r.runData.Decisions[r.index], true
}

// Reset restarts the replay from the beginning.
func (r *Replay) Reset() {
	r.index = 0
}

// Position returns the current position in the replay.
func (r *Replay) Position() int {
	return r.index
}

// Remaining returns the number of decisions remaining.
func (r *Replay) Remaining() int {
	return len(r.runData.Decisions) - r.index
}
