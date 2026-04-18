package profiles

import (
	"errors"
	"net/url"
	"strings"
)

func (p Profile) Normalized() Profile {
	p.ID = strings.TrimSpace(p.ID)
	p.Name = strings.TrimSpace(p.Name)
	p.Endpoint = strings.TrimRight(strings.TrimSpace(p.Endpoint), "/")
	p.Region = strings.TrimSpace(p.Region)
	p.Bucket = strings.TrimSpace(p.Bucket)
	p.AccessKey = strings.TrimSpace(p.AccessKey)
	p.SecretKey = strings.TrimSpace(p.SecretKey)
	p.CDNURL = strings.TrimRight(strings.TrimSpace(p.CDNURL), "/")

	if p.Region == "" {
		p.Region = "us-east-1"
	}
	if p.Name == "" {
		p.Name = p.Bucket
	}
	if p.ID == "" {
		p.ID = strings.ToLower(strings.NewReplacer(" ", "-", "_", "-").Replace(p.Name))
	}
	return p
}

func (p Profile) Validate() error {
	if p.Endpoint == "" {
		return errors.New("endpoint is required")
	}
	if _, err := url.ParseRequestURI(p.Endpoint); err != nil {
		return errors.New("endpoint must be a valid URL")
	}
	if p.Bucket == "" {
		return errors.New("bucket is required")
	}
	if p.AccessKey == "" {
		return errors.New("access key is required")
	}
	if p.SecretKey == "" {
		return errors.New("secret key is required")
	}
	return nil
}

func (p Profile) PublicBaseURL() string {
	if p.CDNURL != "" {
		return p.CDNURL
	}

	endpoint, err := url.Parse(p.Endpoint)
	if err != nil || endpoint.Host == "" || p.Bucket == "" {
		return ""
	}

	endpointPath := strings.TrimRight(endpoint.Path, "/")
	if p.PathStyle || endpointPath != "" {
		return endpoint.Scheme + "://" + endpoint.Host + endpointPath + "/" + p.Bucket
	}

	return endpoint.Scheme + "://" + p.Bucket + "." + endpoint.Host
}

func (p Profile) ObjectURL(key string) string {
	base := p.PublicBaseURL()
	if base == "" {
		return key
	}
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(key, "/")
}

func (p Profile) Public() PublicProfile {
	return PublicProfile{
		ID:         p.ID,
		Name:       p.Name,
		Endpoint:   p.Endpoint,
		Region:     p.Region,
		Bucket:     p.Bucket,
		AccessKey:  p.AccessKey,
		HasSecret:  p.SecretKey != "",
		CDNURL:     p.CDNURL,
		PathStyle:  p.PathStyle,
		PublicBase: p.PublicBaseURL(),
	}
}
