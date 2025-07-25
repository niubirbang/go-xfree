package goxfree

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Installer struct {
	dir string
}

func NewInstaller(option Option) *Installer {
	return &Installer{
		dir: option.GetDir(),
	}
}

func (i *Installer) Run() error {
	return i.check(false)
}

func (i *Installer) Quit() error {
	return nil
}

func (i *Installer) check(force bool) error {
	if !i.existsGeoIP() || force {
		if err := i.downloadGeoIP(); err != nil {
			return err
		}
	}
	if !i.existsGeoSite() || force {
		if err := i.downloadGeoSite(); err != nil {
			return err
		}
	}
	if !i.existsCountry() || force {
		if err := i.downloadCountry(); err != nil {
			return err
		}
	}
	if !i.existsUI() || force {
		if err := i.downloadUI(); err != nil {
			return err
		}
	}
	return nil
}

func (i *Installer) getGeoIPPath() string {
	return path.Join(i.dir, "GeoIP.dat")
}
func (i *Installer) existsGeoIP() bool {
	_, err := os.Stat(i.getGeoIPPath())
	return err == nil
}
func (i *Installer) downloadGeoIP() error {
	log.Println("downloading geoip")
	geoipURL := "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.dat"
	if err := i.download(geoipURL, i.getGeoIPPath()); err != nil {
		log.Println("download geoip failed:", err)
		return err
	}
	log.Println("downloaded geoip")
	return nil
}

func (i *Installer) getGeoSitePath() string {
	return path.Join(i.dir, "GeoSite.dat")
}
func (i *Installer) existsGeoSite() bool {
	_, err := os.Stat(i.getGeoSitePath())
	return err == nil
}
func (i *Installer) downloadGeoSite() error {
	log.Println("downloading geosite")
	geositeURL := "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geosite.dat"
	if err := i.download(geositeURL, i.getGeoSitePath()); err != nil {
		log.Println("download geosite failed:", err)
		return err
	}
	log.Println("downloaded geosite")
	return nil
}

func (i *Installer) getCountryPath() string {
	return path.Join(i.dir, "country.mmdb")
}
func (i *Installer) existsCountry() bool {
	_, err := os.Stat(i.getCountryPath())
	return err == nil
}
func (i *Installer) downloadCountry() error {
	log.Println("downloading country")
	countryURL := "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/country.mmdb"
	if err := i.download(countryURL, i.getCountryPath()); err != nil {
		log.Println("download country failed:", err)
		return err
	}
	log.Println("downloaded country")
	return nil
}

func (i *Installer) getUIPath() string {
	return path.Join(i.dir, "ui")
}
func (i *Installer) existsUI() bool {
	_, err := os.Stat(i.getUIPath())
	return err == nil
}
func (i *Installer) downloadUI() error {
	log.Println("downloading ui")
	uiURL := "https://github.com/MetaCubeX/metacubexd/archive/refs/heads/gh-pages.zip"
	if err := i.downloadAndUnzip(uiURL, i.getUIPath()); err != nil {
		log.Println("download ui failed:", err)
		return err
	}
	log.Println("downloaded ui")
	return nil
}

func (i *Installer) remove(path string) error {
	err := os.RemoveAll(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (i *Installer) download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, resp.Body); err != nil {
		return err
	}
	return nil
}

func (i *Installer) downloadAndUnzip(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return err
	}
	getCommonRootPrefix := func(reader *zip.Reader) string {
		var prefix string
		for _, f := range reader.File {
			parts := strings.Split(f.Name, "/")
			if len(parts) < 2 {
				return ""
			}
			if prefix == "" {
				prefix = parts[0]
			} else if prefix != parts[0] {
				return ""
			}
		}
		return prefix
	}
	rootPrefix := getCommonRootPrefix(zr)
	for _, f := range zr.File {
		relPath := f.Name
		if rootPrefix != "" && strings.HasPrefix(f.Name, rootPrefix+"/") {
			relPath = strings.TrimPrefix(f.Name, rootPrefix+"/")
		}
		if relPath == "" {
			continue
		}
		fpath := filepath.Join(dest, relPath)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid path: %s", fpath)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
