package plugin

const (
	MetadataPrefix  = "plugin."
	MetadataVersion = MetadataPrefix + "version"
	MetadataName    = MetadataPrefix + "name"

	MetadataBaseDirectory = "basedir"
	MetadataRuntime       = "runtime"

	MetadataRepository = "repository"
	MetadataTag        = "tag"
)

type ProcessMetadata struct {
	BaseDirectory string
	Version       string
	Runtime       Runtime
}

type ImageMetadata struct {
	Repository string
	Tag        string
}

type Metadata struct {
	Name    string
	Image   *ImageMetadata
	Process *ProcessMetadata
}

func MapMetadata(properties map[string]string) *Metadata {

	return &Metadata{
		Name: properties[MetadataName],
	}
}

func MapProcessMetadata(properties map[string]string) *ProcessMetadata {
	baseDirectory, exists := properties[MetadataBaseDirectory]
	if !exists {
		return nil
	}
	runtime, exists := properties[MetadataRuntime]
	if !exists {
		return nil
	}
	version, exists := properties[MetadataVersion]
	if !exists {
		return nil
	}
	return &ProcessMetadata{
		BaseDirectory: baseDirectory,
		Version:       version,
		Runtime:       Runtime(runtime),
	}
}

func MapImageMetadata(properties map[string]string) *ImageMetadata {
	repository, exists := properties[MetadataRepository]
	if !exists {
		return nil
	}
	tag, exists := properties[MetadataTag]
	if !exists {
		tag = ""
	}
	return &ImageMetadata{
		Repository: repository,
		Tag:        tag,
	}
}
