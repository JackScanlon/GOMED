package trud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	expectedApiVersion string = "1"
)

var (
	releaseName = map[uint16]string{
		9:   "SnomedReadMap",
		101: "SnomedCTRelease",
	}
)

type ReleaseMetadata struct {
	ID                                 string    `json:"id"`
	Name                               string    `json:"name"`
	ReleaseDate                        string    `json:"releaseDate"`
	ArchiveFileURL                     string    `json:"archiveFileUrl"`
	ArchiveFileName                    string    `json:"archiveFileName"`
	ArchiveFileSizeBytes               int       `json:"archiveFileSizeBytes"`
	ArchiveFileSha256                  string    `json:"archiveFileSha256"`
	ArchiveFileLastModifiedTimestamp   time.Time `json:"archiveFileLastModifiedTimestamp"`
	ChecksumFileURL                    string    `json:"checksumFileUrl"`
	ChecksumFileName                   string    `json:"checksumFileName"`
	ChecksumFileSizeBytes              int       `json:"checksumFileSizeBytes"`
	ChecksumFileLastModifiedTimestamp  time.Time `json:"checksumFileLastModifiedTimestamp"`
	SignatureFileURL                   string    `json:"signatureFileUrl"`
	SignatureFileName                  string    `json:"signatureFileName"`
	SignatureFileSizeBytes             int       `json:"signatureFileSizeBytes"`
	SignatureFileLastModifiedTimestamp time.Time `json:"signatureFileLastModifiedTimestamp"`
	PublicKeyFileURL                   string    `json:"publicKeyFileUrl"`
	PublicKeyFileName                  string    `json:"publicKeyFileName"`
	PublicKeyFileSizeBytes             int       `json:"publicKeyFileSizeBytes"`
	PublicKeyID                        int       `json:"publicKeyId"`
}

type PackageMetadata struct {
	APIVersion string             `json:"apiVersion"`
	Releases   []*ReleaseMetadata `json:"releases"`
	HTTPStatus int                `json:"httpStatus"`
	Message    string             `json:"message"`
}

type Release struct {
	Name       string           `json:"Name"`
	URL        string           `json:"URL"`
	CategoryId uint16           `json:"CategoryId"`
	Metadata   *ReleaseMetadata `json:"Metadata"`
}

type ReleaseOpt func(*Release)

func WithURL(url string) ReleaseOpt {
	return func(r *Release) {
		r.URL = url
	}
}

func WithCategory(id uint16) ReleaseOpt {
	return func(r *Release) {
		r.CategoryId = id

		if name, ok := releaseName[id]; ok {
			r.Name = name
		} else {
			r.Name = ""
		}
	}
}

func WithMetadata(m *ReleaseMetadata) ReleaseOpt {
	return func(r *Release) {
		r.Metadata = m
	}
}

func NewRelease(opts ...ReleaseOpt) *Release {
	release := &Release{}
	for _, opt := range opts {
		opt(release)
	}

	return release
}

func (r *Release) HasRelease(directory string) (bool, error) {
	filepath := path.Join(directory, r.Metadata.Name)

	_, err := os.Stat(filepath)
	if err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func queryRelease(ctx context.Context, url string, release *Release) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var metadata PackageMetadata
	if err := json.NewDecoder(res.Body).Decode(&metadata); err != nil {
		return err
	}

	if metadata.APIVersion != expectedApiVersion {
		return fmt.Errorf("expected api version of %s but got %s", expectedApiVersion, metadata.APIVersion)
	}

	if len(metadata.Releases) < 1 {
		return fmt.Errorf("failed to find releases for url: %s", url)
	}

	release.Metadata = metadata.Releases[0]
	return nil
}

func getReleases(category Category, apiKey string) ([]*Release, error) {
	categoryIds := category.GetIds()
	categoryLen := len(categoryIds)
	if categoryLen < 1 {
		return nil, errors.New("invalid category, no known category id(s)")
	}

	ctx := context.Background()
	errs, _ := errgroup.WithContext(ctx)

	releases := make([]*Release, categoryLen)
	for i := 0; i < categoryLen; i++ {
		id := categoryIds[i]
		url := fmt.Sprintf(categoryUrl, apiKey, id)

		releases[i] = NewRelease(
			WithURL(url),
			WithCategory(id),
		)

		errs.Go(func() error {
			if err := queryRelease(ctx, url, releases[i]); err != nil {
				return fmt.Errorf("[gid: %d, cat: %d] failed to query package at <%s> with err: %v", i, id, url, err)
			}

			return nil
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return releases, nil
}
