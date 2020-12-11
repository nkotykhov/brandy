package rpm

import (
	"bytes"
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	LeadMagic = 0xedabeedb
	LeadSize  = 96
)

const (
	PackageTypeBinary int16 = 0
	PackageTypeSource int16 = 1
)

const (
	// Magic plus header structure version
	HeaderMagic = 0x8eade801
)

type HeaderTag int32

const (
	TagName HeaderTag = 1000 + iota
	TagVersion
	TagRelease
	TagEpoch
	TagSummary
	TagDescription
	TagBuildTime
	TagBuildHost
	TagInstallTime
	TagSize
	TagDistribution
	TagVendor
	TagGif
	TagXmp
	TagLicense
	TagPackager
	TagGroup
	TagChangelog
	TagSource
	TagPatch
	TagURL
	TagOS
	TagArch
	TagPrein
	TagPostin
	TagPreun
	TagPostun
	TagOldFilenames
	TagFileSizes
	TagFileStates
	TagFileModes
	TagFileUIDs
	TagFileGIDs
	TagFileRDevs
	TagFileMTimes
	TagFileMD5S
	TagFileLinkTos
	TagFileFlags
	TagRoot
	TagFileUsername
	TagFileGroupname
	TagExclude
	TagExclusive
	TagIcon
	TagSourceRPM
	TagFileVerifyFlags
	TagArchiveSize
	TagProvideName
	TagRequireFlags
	TagRequireName
	TagRequireVersion
	TagNoSource
	TagNoPatch
	TagConflictFlags
	TagConflictName
	TagConflictVersion
	TagDefaultPrefix
	TagBuildroot
	TagInstallPrefix
	TagExclusiveArch
	TagExclusiveOS
	TagEnclusiveArch
	TagEnclusiveOS
	TagAutoReqProv
	TagRPMVersion
	TagTriggerScripts
	TagTriggerName
	TagTriggerVersion
	TagTriggerFlags
	TagTriggerIndex
)

const(
	TagVerifyScript HeaderTag = 1079+iota
	TagChangelogTime
	TagChangelogName
	TagChangelogText
	TagBrokenMD5
	TagPrereq
	TagPreinProg
	TagPostinProg
	TagPreunProg
	TagPostunProg
	TagBuildArchs
	TagObsoleteName
	TagVerifyScriptProg
	TagTriggerScriptProg
	TagDocdir
	TagCookie
	TagFileDevices
	TagFileInodes
	TagFileLangs
	TagPrefixes
	TagInstPrefixes
	TagTriggerin
	TagTriggerun
	TagTriggerpostun
	TagAutoreq
	TagAutoprov
	TagCapability
	TagSourcePackage
	TagOldOrigFilenames
	TagBuildPrereq
	TagBuildRequires
	TagBuildConflicts
	TagBuildMacros
	TagProvideFlags
	TagProvideVersion
	TagObsoleteFlags
	TagObsoleteVersion
	TagDirIndexes
	TagBaseNames
	TagDirNames
	TagOrigDirIndexes
	TagOrigBaseNames
	TagOrigdirNames
	TagOptFlags
	TagDistURL
	TagPayloadFormat
	TagPayloadCompressor
	TagPayloadFlags
	TagInstallColor
	TagInstallTID
	TagRemoveTID
	TagSha1RHN
	TagRHNPlatform
	TagPlatform
	TagPatchesName
	TagPatchesFlags
	TagPatchesVersion
	TagCacheCtime
	TagCachePkgPath
	TagCachePkgSize
	TagCachePkgMtime
	TagFileColors
	TagFileClass
	TagClassDict
	TagFileDependsX
	TagFileDependsN
	TagDependsDict
	TagSourcePkgID
	TagFileContexts
	TagFsContexts
	TagReContexts
	TagPolicies

)

const (
	TagFilenames  = 5000
)

type HeaderDataType int32

const (
	DataTypeNull        HeaderDataType = 0
	DataTypeChar                       = 1
	DataTypeInt8                       = 2
	DataTypeInt16                      = 3
	DataTypeInt32                      = 4
	DataTypeInt64                      = 5
	DataTypeString                     = 6
	DataTypeBin                        = 7
	DataTypeStringArray                = 8

	DataTypeNotFound = -1
)

const HeaderIndexEntrySize = 16

type HeaderIndexEntry struct {
	Tag      HeaderTag
	DataType HeaderDataType
	Offset   int32
	Count    int32
}

type headerInfo struct {
	Magic    uint32
	Reserved uint32
	IndexCnt uint32
	Size     uint32
}

type Header struct {
	indexes []HeaderIndexEntry
	data    []byte
}

func ReadHeader(r io.Reader) (*Header, error) {
	var hInfo headerInfo
	if err := binary.Read(r, binary.BigEndian, &hInfo); err != nil {
		return nil, fmt.Errorf("error reading RPM header: %s", err.Error())
	}

	if hInfo.Magic != HeaderMagic {
		return nil, fmt.Errorf("error reading RPM header: bad header magic or version")
	}

	nIdx := int(hInfo.IndexCnt)
	indexTable := make([]byte, HeaderIndexEntrySize*nIdx)
	data := make([]byte, hInfo.Size)
	if _, err := io.ReadFull(r, indexTable); err != nil {
		return nil, fmt.Errorf("error reading RPM header index table: %s", err.Error())
	}
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, fmt.Errorf("error reading RPM header data: %s", err.Error())
	}

	header := &Header{}
	header.indexes = make([]HeaderIndexEntry, nIdx)
	header.data = data

	br := bytes.NewReader(indexTable)
	for i := 0; i < nIdx; i++ {
		if err := binary.Read(br, binary.BigEndian, &header.indexes[i]); err != nil {
			return nil, fmt.Errorf("error parsing index table: %s", err.Error())
		}
	}
	return header, nil
}

func (h Header) AvailableTags() ([]HeaderTag){
	tags := make([]HeaderTag, len(h.indexes))
	for i,idx := range h.indexes {
		tags[i] = idx.Tag
	}

	return tags

}
func (h Header) GetTag(t HeaderTag) (HeaderDataType, []byte, error) {
	for _, idx := range h.indexes {
		if idx.Tag == t {
			d, err := readData(&idx, h.data)
			return idx.DataType, d, err
		}
	}

	return DataTypeNotFound, nil, fmt.Errorf("error tag not found")
}


func (h Header) getTag(t HeaderTag) (*HeaderIndexEntry, []byte, error) {
	for _, idx := range h.indexes {
		if idx.Tag == t {
			d, err := readData(&idx, h.data)
			return &idx, d, err
		}
	}

	return nil, nil, fmt.Errorf("error tag not found")
}

func (h Header) GetString(t HeaderTag) (string, error) {
	idx, v, err := h.getTag(t)
	if err != nil {
		return "", err
	}
	if idx.DataType != DataTypeString {
		return "", fmt.Errorf("error invalid data type %d, %d", idx.DataType, DataTypeStringArray)
	}
	br := bufio.NewReader(bytes.NewReader(v))
	cstr, err := br.ReadBytes(0)
	if err != nil {
		return "", err
	}
	s := string(cstr[0:len(cstr)-1])
	return s, nil
}

func (h Header) GetStrings(t HeaderTag) ([]string, error) {
	idx, v, err := h.getTag(t)
	if err != nil {
		return nil, err
	}
	if idx.DataType != DataTypeStringArray {
		return nil, fmt.Errorf("error invalid data type %d, %d", idx.DataType, DataTypeStringArray)
	}
	s := make([]string, idx.Count)
	br := bufio.NewReader(bytes.NewReader(v))
	for i := 0; i < int(idx.Count); i++ {
		cstr, err := br.ReadBytes(0)
		if err != nil {
			return nil, err
		}
		s[i] = string(cstr[0:len(cstr)-1])
	}
	return s, nil
}

func readData(idx *HeaderIndexEntry, d []byte) ([]byte, error) {
	var end int
	switch idx.DataType {
	case DataTypeNull:
		return d[idx.Offset:idx.Offset], nil
	case DataTypeChar, DataTypeInt8:
		end = int(idx.Offset + idx.Count)
	case DataTypeInt16:
		end = int(idx.Offset + 2*idx.Count)
	case DataTypeInt32:
		end = int(idx.Offset + 4*idx.Count)
	case DataTypeInt64:
		end = int(idx.Offset + 8*idx.Count)
	case DataTypeString, DataTypeStringArray:
		cnt := idx.Count
		if idx.DataType == DataTypeString {
			cnt = 1
		}
		end = int(idx.Offset)
		for cnt > 0 {
			next := bytes.IndexByte(d[end:], 0)
			if next < 0 {
				return nil, fmt.Errorf("trancuted string")
			}
			end += next + 1
			cnt--
		}
	default:
		return nil, fmt.Errorf("invalid index referense")
	}

	return d[idx.Offset:end], nil
}
