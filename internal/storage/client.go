package storage

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/PouryaMansouri/BucketDesk/internal/profiles"
)

type Client struct {
	profile profiles.Profile
	s3      *s3.Client
}

type Folder struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
}

type Object struct {
	Key          string `json:"key"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	LastModified string `json:"lastModified,omitempty"`
	URL          string `json:"url"`
	MimeType     string `json:"mimeType"`
}

type ListResult struct {
	Prefix    string   `json:"prefix"`
	Folders   []Folder `json:"folders"`
	Objects   []Object `json:"objects"`
	NextToken string   `json:"nextToken,omitempty"`
	Limit     int32    `json:"limit"`
}

func New(profile profiles.Profile) *Client {
	return &Client{
		profile: profile,
		s3: s3.New(s3.Options{
			Region:                     profile.Region,
			BaseEndpoint:               aws.String(profile.Endpoint),
			UsePathStyle:               profile.PathStyle,
			Credentials:                aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(profile.AccessKey, profile.SecretKey, "")),
			RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
			ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
		}),
	}
}

func (c *Client) Test(ctx context.Context) (map[string]string, error) {
	if _, err := c.s3.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(c.profile.Bucket)}); err != nil {
		return nil, err
	}

	key := "_minio-manager/connection-tests/" + time.Now().Format("20060102150405") + "-" + randomHex(4) + ".txt"
	if _, err := c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.profile.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewBufferString("minio-manager connection test"),
		ContentType: aws.String("text/plain; charset=utf-8"),
	}); err != nil {
		return nil, err
	}

	_, _ = c.s3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.profile.Bucket),
		Key:    aws.String(key),
	})

	return map[string]string{
		"bucket":     c.profile.Bucket,
		"endpoint":   c.profile.Endpoint,
		"publicBase": c.profile.PublicBaseURL(),
	}, nil
}

func (c *Client) List(ctx context.Context, prefix string, token string, limit int32) (ListResult, error) {
	if limit < 1 || limit > 500 {
		limit = 100
	}

	input := &s3.ListObjectsV2Input{
		Bucket:            aws.String(c.profile.Bucket),
		Prefix:            aws.String(prefix),
		Delimiter:         aws.String("/"),
		MaxKeys:           aws.Int32(limit),
		ContinuationToken: optionalString(token),
	}

	output, err := c.s3.ListObjectsV2(ctx, input)
	if err != nil {
		return ListResult{}, err
	}

	result := ListResult{
		Prefix:  prefix,
		Folders: make([]Folder, 0, len(output.CommonPrefixes)),
		Objects: make([]Object, 0, len(output.Contents)),
		Limit:   limit,
	}

	if output.NextContinuationToken != nil {
		result.NextToken = *output.NextContinuationToken
	}

	for _, folder := range output.CommonPrefixes {
		if folder.Prefix == nil {
			continue
		}
		clean := strings.TrimRight(*folder.Prefix, "/")
		result.Folders = append(result.Folders, Folder{
			Name:   path.Base(clean),
			Prefix: *folder.Prefix,
		})
	}

	for _, item := range output.Contents {
		if item.Key == nil || *item.Key == prefix {
			continue
		}
		name := path.Base(*item.Key)
		object := Object{
			Key:      *item.Key,
			Name:     name,
			Size:     aws.ToInt64(item.Size),
			URL:      c.profile.ObjectURL(*item.Key),
			MimeType: mimeType(name),
		}
		if item.LastModified != nil {
			object.LastModified = item.LastModified.Format(time.RFC3339)
		}
		result.Objects = append(result.Objects, object)
	}

	return result, nil
}

func (c *Client) Upload(ctx context.Context, prefix string, filename string, contentType string, body io.Reader) (Object, error) {
	key := buildObjectKey(prefix, filename)
	if contentType == "" {
		contentType = mimeType(filename)
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return Object{}, err
	}

	_, err = c.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       aws.String(c.profile.Bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(data),
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("public, max-age=31536000"),
	})
	if err != nil {
		return Object{}, err
	}

	return Object{
		Key:      key,
		Name:     path.Base(key),
		Size:     int64(len(data)),
		URL:      c.profile.ObjectURL(key),
		MimeType: contentType,
	}, nil
}

func (c *Client) Delete(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	objects := make([]types.ObjectIdentifier, 0, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		objects = append(objects, types.ObjectIdentifier{Key: aws.String(key)})
	}

	if len(objects) == 0 {
		return nil
	}

	_, err := c.s3.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(c.profile.Bucket),
		Delete: &types.Delete{
			Objects: objects,
			Quiet:   aws.Bool(true),
		},
	})
	return err
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return aws.String(value)
}

func buildObjectKey(prefix string, filename string) string {
	filename = strings.TrimSpace(filename)
	filename = strings.ReplaceAll(filename, "\\", "-")
	filename = strings.ReplaceAll(filename, "/", "-")
	if filename == "" {
		filename = "file.bin"
	}

	ext := path.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	base = strings.Trim(strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
			return r
		}
		if r == ' ' {
			return '-'
		}
		return -1
	}, base), "-_")
	if base == "" {
		base = "file"
	}

	prefix = strings.Trim(strings.TrimSpace(prefix), "/")
	name := fmt.Sprintf("%s-%d-%s%s", base, time.Now().UnixMilli(), randomHex(4), strings.ToLower(ext))
	if prefix == "" {
		return name
	}
	return prefix + "/" + name
}

func randomHex(size int) string {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(data)
}

func mimeType(filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	if value := mime.TypeByExtension(ext); value != "" {
		return value
	}
	switch ext {
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
