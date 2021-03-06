package resource

import (
	"path/filepath"
	"testing"

	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/xa0082249956/hugo/helpers"
	"github.com/xa0082249956/hugo/hugofs"
	"github.com/xa0082249956/hugo/media"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func newTestResourceSpec(assert *require.Assertions) *Spec {
	return newTestResourceSpecForBaseURL(assert, "https://example.com/")
}

func newTestResourceSpecForBaseURL(assert *require.Assertions, baseURL string) *Spec {
	cfg := viper.New()
	cfg.Set("baseURL", baseURL)
	cfg.Set("resourceDir", "resources")
	cfg.Set("contentDir", "content")

	imagingCfg := map[string]interface{}{
		"resampleFilter": "linear",
		"quality":        68,
		"anchor":         "left",
	}

	cfg.Set("imaging", imagingCfg)

	fs := hugofs.NewMem(cfg)

	s, err := helpers.NewPathSpec(fs, cfg)

	assert.NoError(err)

	spec, err := NewSpec(s, media.DefaultTypes)
	assert.NoError(err)
	return spec
}

func newTestResourceOsFs(assert *require.Assertions) *Spec {
	cfg := viper.New()
	cfg.Set("baseURL", "https://example.com")

	workDir, err := ioutil.TempDir("", "hugores")

	if runtime.GOOS == "darwin" && !strings.HasPrefix(workDir, "/private") {
		// To get the entry folder in line with the rest. This its a little bit
		// mysterious, but so be it.
		workDir = "/private" + workDir
	}

	cfg.Set("workingDir", workDir)
	cfg.Set("contentDir", filepath.Join(workDir, "content"))
	cfg.Set("resourceDir", filepath.Join(workDir, "res"))

	fs := hugofs.NewFrom(hugofs.Os, cfg)
	fs.Destination = &afero.MemMapFs{}

	s, err := helpers.NewPathSpec(fs, cfg)

	assert.NoError(err)

	spec, err := NewSpec(s, media.DefaultTypes)
	assert.NoError(err)
	return spec

}

func fetchSunset(assert *require.Assertions) *Image {
	return fetchImage(assert, "sunset.jpg")
}

func fetchImage(assert *require.Assertions, name string) *Image {
	spec := newTestResourceSpec(assert)
	return fetchImageForSpec(spec, assert, name)
}

func fetchImageForSpec(spec *Spec, assert *require.Assertions, name string) *Image {
	r := fetchResourceForSpec(spec, assert, name)
	assert.IsType(&Image{}, r)
	return r.(*Image)
}

func fetchResourceForSpec(spec *Spec, assert *require.Assertions, name string) Resource {
	src, err := os.Open(filepath.FromSlash("testdata/" + name))
	assert.NoError(err)

	assert.NoError(spec.BaseFs.ContentFs.MkdirAll(filepath.Dir(name), 0755))
	out, err := spec.BaseFs.ContentFs.Create(name)
	assert.NoError(err)
	_, err = io.Copy(out, src)
	out.Close()
	src.Close()
	assert.NoError(err)

	factory := func(s string) string {
		return path.Join("/a", s)
	}

	r, err := spec.NewResourceFromFilename(factory, name, name)
	assert.NoError(err)

	return r
}

func assertImageFile(assert *require.Assertions, fs afero.Fs, filename string, width, height int) {
	f, err := fs.Open(filename)
	if err != nil {
		printFs(fs, "", os.Stdout)
	}
	assert.NoError(err)
	defer f.Close()

	config, _, err := image.DecodeConfig(f)
	assert.NoError(err)

	assert.Equal(width, config.Width)
	assert.Equal(height, config.Height)
}

func assertFileCache(assert *require.Assertions, fs afero.Fs, filename string, width, height int) {
	assertImageFile(assert, fs, filepath.Join("_gen/images", filename), width, height)
}

func writeSource(t testing.TB, fs *hugofs.Fs, filename, content string) {
	writeToFs(t, fs.Source, filename, content)
}

func writeToFs(t testing.TB, fs afero.Fs, filename, content string) {
	if err := afero.WriteFile(fs, filepath.FromSlash(filename), []byte(content), 0755); err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}
}

func printFs(fs afero.Fs, path string, w io.Writer) {
	if fs == nil {
		return
	}
	afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			s := path
			if lang, ok := info.(hugofs.LanguageAnnouncer); ok {
				s = s + "\t" + lang.Lang()
			}
			if fp, ok := info.(hugofs.FilePather); ok {
				s += "\tFilename: " + fp.Filename() + "\tBase: " + fp.BaseDir()
			}
			fmt.Fprintln(w, "    ", s)
		}
		return nil
	})
}
