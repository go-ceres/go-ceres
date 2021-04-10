package file

type Option func(f *fileSource)

func Unmarshal(unmarshal string) Option {
	return func(f *fileSource) {
		f.unmarshal = unmarshal
	}
}
