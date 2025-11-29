package dat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"unicode/utf16"
)

var (
	// VDataMagic is the 8-byte sequence separating fixed and variable data.
	VDataMagic = []byte{0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb, 0xbb}
	// NullTerminator is the sequence ending strings.
	NullTerminator = []byte{0x00, 0x00, 0x00, 0x00}
)

// DatFile represents a parsed .datc64 file.
type DatFile struct {
	RowCount     uint32
	RowLength    int
	DataFixed    []byte
	DataVariable []byte
}

// Read parses a .datc64 file from the given reader.
func Read(r io.Reader) (*DatFile, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read dat file: %w", err)
	}

	if len(data) < 4 {
		return nil, fmt.Errorf("file too small")
	}

	rowCount := binary.LittleEndian.Uint32(data[0:4])

	// Find the boundary (0xbb...)
	// The boundary starts after the 4-byte row count.
	// We search for the magic sequence.
	// According to dat2json, we need to find an aligned sequence.

	boundaryIdx := -1
	searchStart := 4
	for {
		idx := bytes.Index(data[searchStart:], VDataMagic)
		if idx == -1 {
			break
		}

		absoluteIdx := searchStart + idx
		// Check alignment if rowCount > 0
		// The boundary offset (relative to data start, excluding 4 bytes header?)
		// dat2json: file.subarray(INT_ROWCOUNT) -> findAlignedSequence
		// idx is relative to searchStart (which is 4).
		// So absoluteIdx is the index in `data`.
		// The fixed data size is absoluteIdx - 4.
		// If rowCount > 0, fixed data size must be divisible by rowCount.

		fixedDataSize := absoluteIdx - 4
		if rowCount == 0 || fixedDataSize%int(rowCount) == 0 {
			boundaryIdx = absoluteIdx
			break
		}

		searchStart = absoluteIdx + 1
	}

	if boundaryIdx == -1 {
		return nil, fmt.Errorf("variable data section boundary not found")
	}

	fixedDataSize := boundaryIdx - 4
	rowLength := 0
	if rowCount > 0 {
		rowLength = fixedDataSize / int(rowCount)
	}

	return &DatFile{
		RowCount:     rowCount,
		RowLength:    rowLength,
		DataFixed:    data[4:boundaryIdx],
		DataVariable: data[boundaryIdx:],
	}, nil
}

// ReadBool reads a boolean at the given offset.
func (f *DatFile) ReadBool(offset int) (bool, error) {
	if offset >= len(f.DataFixed) {
		return false, fmt.Errorf("offset out of bounds")
	}
	return f.DataFixed[offset] != 0, nil
}

// ReadInt32 reads a 32-bit integer at the given offset.
func (f *DatFile) ReadInt32(offset int) (int32, error) {
	if offset+4 > len(f.DataFixed) {
		return 0, fmt.Errorf("offset out of bounds")
	}
	return int32(binary.LittleEndian.Uint32(f.DataFixed[offset : offset+4])), nil
}

// ReadUint32 reads a 32-bit unsigned integer at the given offset.
func (f *DatFile) ReadUint32(offset int) (uint32, error) {
	if offset+4 > len(f.DataFixed) {
		return 0, fmt.Errorf("offset out of bounds")
	}
	return binary.LittleEndian.Uint32(f.DataFixed[offset : offset+4]), nil
}

// ReadInt64 reads a 64-bit integer at the given offset.
func (f *DatFile) ReadInt64(offset int) (int64, error) {
	if offset+8 > len(f.DataFixed) {
		return 0, fmt.Errorf("offset out of bounds")
	}
	return int64(binary.LittleEndian.Uint64(f.DataFixed[offset : offset+8])), nil
}

// ReadUint64 reads a 64-bit unsigned integer at the given offset.
func (f *DatFile) ReadUint64(offset int) (uint64, error) {
	if offset+8 > len(f.DataFixed) {
		return 0, fmt.Errorf("offset out of bounds")
	}
	return binary.LittleEndian.Uint64(f.DataFixed[offset : offset+8]), nil
}

// ReadFloat32 reads a 32-bit float at the given offset.
func (f *DatFile) ReadFloat32(offset int) (float32, error) {
	if offset+4 > len(f.DataFixed) {
		return 0, fmt.Errorf("offset out of bounds")
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(f.DataFixed[offset : offset+4])), nil
}

// ReadFloat64 reads a 64-bit float at the given offset.
func (f *DatFile) ReadFloat64(offset int) (float64, error) {
	if offset+8 > len(f.DataFixed) {
		return 0, fmt.Errorf("offset out of bounds")
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(f.DataFixed[offset : offset+8])), nil
}

// ReadString reads a string pointer at offset, then reads the string from DataVariable.
func (f *DatFile) ReadString(offset int) (string, error) {
	if offset+8 > len(f.DataFixed) {
		return "", fmt.Errorf("offset out of bounds")
	}
	// Pointers are 8 bytes
	ptr := binary.LittleEndian.Uint64(f.DataFixed[offset : offset+8])
	if ptr == 0xfefefefefefefefe { // Null check (heuristic)
		return "", nil
	}
	return f.readStringAt(int(ptr))
}

func (f *DatFile) readStringAt(offset int) (string, error) {
	if offset >= len(f.DataVariable) {
		return "", nil
	}

	// dat2json logic:
	// 1. Find 4-byte null terminator (0x00 0x00 0x00 0x00)
	// 2. Ensure (end - offset) is even (UTF-16 alignment)

	terminator := []byte{0, 0, 0, 0}
	end := offset

	for {
		idx := bytes.Index(f.DataVariable[end:], terminator)
		if idx == -1 {
			// No terminator found, return what we have? Or error?
			// For now, let's assume the string goes to the end if no terminator found (shouldn't happen in valid files)
			end = len(f.DataVariable)
			break
		}

		foundAt := end + idx
		// Check alignment: The string length (foundAt - offset) must be even for UTF-16
		if (foundAt-offset)%2 == 0 {
			end = foundAt
			break
		}

		// If not aligned, continue searching after the first byte of the found sequence
		end = foundAt + 1
	}

	// Decode UTF-16LE
	stringBytes := f.DataVariable[offset:end]
	u16s := make([]uint16, len(stringBytes)/2)
	for i := 0; i < len(u16s); i++ {
		u16s[i] = binary.LittleEndian.Uint16(stringBytes[i*2 : i*2+2])
	}

	return string(utf16.Decode(u16s)), nil
}

// ReadArray reads an array header at offset.
// Returns the number of elements and the offset of the first element in DataVariable.
func (f *DatFile) ReadArray(offset int) (length int, varOffset int, err error) {
	if offset+16 > len(f.DataFixed) {
		return 0, 0, fmt.Errorf("offset out of bounds")
	}
	length64 := binary.LittleEndian.Uint64(f.DataFixed[offset : offset+8])
	varOffset64 := binary.LittleEndian.Uint64(f.DataFixed[offset+8 : offset+16])

	// Check if varOffset is within bounds (if length > 0)
	if length64 > 0 && int(varOffset64) >= len(f.DataVariable) {
		return int(length64), int(varOffset64), fmt.Errorf("variable offset out of bounds: %d", varOffset64)
	}

	return int(length64), int(varOffset64), nil
}
