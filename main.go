package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jessevdk/go-flags"
)

var sess *session.Session = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var s3svc *s3.S3 = s3.New(sess)

var opts struct {
	DryRun bool   `short:"D" long:"dry-run" description:"Report what objects would be restored without restoring them"`
	Prefix string `short:"p" long:"prefix" description:"only look for deleted objects under the given prefix"`
}

func processMarkers(markers []*s3.DeleteMarkerEntry, bucket string, dryRun bool) {
	for _, marker := range markers {
		if *marker.IsLatest {
			if dryRun {
				fmt.Printf("Would have restored %s\n", *marker.Key)
			} else {
				fmt.Printf("Restoring %s\n", *marker.Key)
				s3svc.DeleteObject(&s3.DeleteObjectInput{
					Bucket:    &bucket,
					Key:       marker.Key,
					VersionId: marker.VersionId,
				})
			}
		}
	}
}

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	if len(args) < 1 {
		panic("An s3 bucket is required")
	}

	bucket := args[0]

	objects, err := s3svc.ListObjectVersions(&s3.ListObjectVersionsInput{
		Bucket: &bucket,
		Prefix: &opts.Prefix,
	})

	if err != nil {
		panic(err)
	}

	processMarkers(objects.DeleteMarkers, bucket, opts.DryRun)

	keyMarker := objects.NextKeyMarker
	for keyMarker != nil {
		objects, err := s3svc.ListObjectVersions(&s3.ListObjectVersionsInput{
			Bucket:    &bucket,
			Prefix:    &opts.Prefix,
			KeyMarker: keyMarker,
		})

		if err != nil {
			panic(err)
		}
		processMarkers(objects.DeleteMarkers, bucket, opts.DryRun)
		keyMarker = objects.NextKeyMarker
	}
}
