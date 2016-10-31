package interfaces

import "gopkg.in/vmihailenco/msgpack.v2"

func (f File) MarshalMsgpack() ([]byte, error) {

	return msgpack.Marshal(
		f.FileID,
		f.BucketID,
		f.FileName,
		f.LuaScript,
		f.ContentType,
		f.Owners,
		f.IsPrivate,
		f.IsReadOnly,
		f.MetaData,
		f.StructuralData,
		f.RawData,
		f.CreatedAt,
		f.UpdatedAt,
	)
}

func (f *File) UnmarshalMsgpack(b []byte) error {
	return msgpack.Unmarshal(
		b,
		&f.FileID,
		&f.BucketID,
		&f.FileName,
		&f.LuaScript,
		&f.ContentType,
		&f.Owners,
		&f.IsPrivate,
		&f.IsReadOnly,
		&f.MetaData,
		&f.StructuralData,
		&f.RawData,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
}

func (f Bucket) MarshalMsgpack() ([]byte, error) {
	return msgpack.Marshal(
		f.BucketID,
		f.BucketName,
		f.Owners,
		f.LuaScript,
		f.MetaData,
		f.StructuralData,
		f.RawData,
		f.CreatedAt,
		f.UpdatedAt,

		f.MetaDataStoreName,
		f.StructuralDataStoreName,
		f.DataStoreName,
	)
}

func (f *Bucket) UnmarshalMsgpack(b []byte) error {
	return msgpack.Unmarshal(
		b,
		&f.BucketID,
		&f.BucketName,
		&f.Owners,
		&f.LuaScript,
		&f.MetaData,
		&f.StructuralData,
		&f.RawData,
		&f.CreatedAt,
		&f.UpdatedAt,

		&f.MetaDataStoreName,
		&f.StructuralDataStoreName,
		&f.DataStoreName,
	)
}
