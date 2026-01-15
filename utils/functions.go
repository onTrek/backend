package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/twpayne/go-polyline"
)

const bucketName = "NOME-TUO-BUCKET.appspot.com"

func DeleteFiles(client *storage.Client, gpxData Gpx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)

	gpxCloudPath := fmt.Sprintf("gpxs/%s.gpx", gpxData.StoragePath)
	mapCloudPath := fmt.Sprintf("maps/%s.png", gpxData.StoragePath)

	gpxObj := bucket.Object(gpxCloudPath)
	if err := gpxObj.Delete(ctx); err != nil {
		if err != storage.ErrObjectNotExist {
			return fmt.Errorf("impossibile eliminare GPX (%s): %w", gpxCloudPath, err)
		}
		log.Printf("File GPX %s non trovato, salto cancellazione.", gpxCloudPath)
	} else {
		log.Printf("File GPX %s eliminato correttamente.", gpxCloudPath)
	}

	mapObj := bucket.Object(mapCloudPath)
	if err := mapObj.Delete(ctx); err != nil {
		if err != storage.ErrObjectNotExist {
			return fmt.Errorf("impossibile eliminare Mappa (%s): %w", mapCloudPath, err)
		}
		log.Printf("File Mappa %s non trovato, salto cancellazione.", mapCloudPath)
	} else {
		log.Printf("File Mappa %s eliminato correttamente.", mapCloudPath)
	}

	return nil
}

func SaveFile(client *storage.Client, file *multipart.FileHeader, pathPrefix string, id string, extension string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer src.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	bucket := client.Bucket(bucketName)

	objectPath := fmt.Sprintf("%s/%s%s", pathPrefix, id, extension)

	if pathPrefix != "gpxs" && pathPrefix != "avatars" {
		return "", fmt.Errorf("invalid storage path: %s", pathPrefix)
	}

	obj := bucket.Object(objectPath)
	wc := obj.NewWriter(ctx)

	if pathPrefix == "gpxs" {
		wc.ContentType = "application/gpx+xml"
	} else {
		wc.ContentType = file.Header.Get("Content-Type")
	}

	wc.Metadata = map[string]string{
		"original-name": file.Filename,
	}

	if _, err = io.Copy(wc, src); err != nil {
		return "", fmt.Errorf("error uploading to cloud: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("error closing cloud writer: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)

	return publicURL, nil
}

func CalculateStats(file *multipart.FileHeader) (GPXStats, error) {
	var stats GPXStats

	src, err := file.Open()
	if err != nil {
		return GPXStats{}, fmt.Errorf("errore apertura file: %w", err)
	}
	defer src.Close()

	gpxData, err := gpx.Parse(src)
	if err != nil {
		return GPXStats{}, fmt.Errorf("errore parsing GPX: %w", err)
	}

	var totalDistance float64
	var totalAscent float64
	var totalDescent float64
	var startTime *time.Time
	var endTime *time.Time
	var maxElevation *float64
	var minElevation *float64

	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			points := segment.Points
			for i := 1; i < len(points); i++ {
				p1 := points[i-1]
				p2 := points[i]

				// Distanza in km
				totalDistance += p1.Distance2D(&p2) / 1000.0

				// Dislivello
				if p1.Elevation.NotNull() && p2.Elevation.NotNull() {
					diff := p2.Elevation.Value() - p1.Elevation.Value()
					if diff > 0 {
						totalAscent += diff
					} else {
						totalDescent += -diff
					}
				}

				// Altitudine max/min
				if p2.Elevation.NotNull() {
					elev := p2.Elevation.Value()
					if maxElevation == nil || elev > *maxElevation {
						maxElevation = &elev
					}
					if minElevation == nil || elev < *minElevation {
						minElevation = &elev
					}
				}

				// Timestamp
				if startTime == nil && !p1.Timestamp.IsZero() {
					startTime = &p1.Timestamp
				}
				if !p2.Timestamp.IsZero() {
					endTime = &p2.Timestamp
				}
			}
		}
	}

	var durationStr string
	if startTime != nil && endTime != nil {
		duration := endTime.Sub(*startTime)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		seconds := int(duration.Seconds()) % 60
		durationStr = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	} else {
		durationStr = "00:00:00"
	}

	// Popola i risultati
	stats.Km = math.Round(totalDistance*100) / 100 // 2 decimali
	stats.Ascent = int(math.Round(totalAscent))    // intero
	stats.Descent = int(math.Round(totalDescent))  // intero
	stats.Duration = durationStr
	if maxElevation != nil {
		stats.MaxAltitude = int(math.Round(*maxElevation)) // intero
	}
	if minElevation != nil {
		stats.MinAltitude = int(math.Round(*minElevation)) // intero
	}

	return stats, nil
}

func CreateMap(file *multipart.FileHeader, client *storage.Client, storagePath string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("errore apertura file: %w", err)
	}
	defer src.Close()

	gpxData, err := gpx.Parse(src)
	if err != nil {
		return "", fmt.Errorf("errore parsing GPX: %w", err)
	}

	var coords [][]float64
	var urlPath string

	minLat, maxLat := 90.0, -90.0
	minLon, maxLon := 180.0, -180.0

	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				lat := point.Latitude
				lon := point.Longitude

				if lat < minLat {
					minLat = lat
				}
				if lat > maxLat {
					maxLat = lat
				}
				if lon < minLon {
					minLon = lon
				}
				if lon > maxLon {
					maxLon = lon
				}
				coords = append(coords, []float64{lat, lon})
			}
		}
	}

	if len(coords) == 0 {
		return "", fmt.Errorf("invalid GPX file: no coordinates found")
	}

	centerLat := (minLat + maxLat) / 2
	centerLon := (minLon + maxLon) / 2

	width := 1280.0
	height := 960.0
	padding := 0.1

	latDiff := maxLat - minLat
	lonDiff := maxLon - minLon

	zoomLat := math.Log2(360 / (latDiff * (1 + padding)))
	zoomLon := math.Log2(360 / (lonDiff * (1 + padding)))
	zoom := math.Min(zoomLat, zoomLon)
	zoom = math.Min(zoom, 20)
	zoom = math.Max(zoom, 1)

	path := string(polyline.EncodeCoords(coords))
	path = url.QueryEscape(path)

	token := os.Getenv("MAPBOX_TOKEN")

	start := coords[0]
	startMarker := fmt.Sprintf("%f,%f", start[1], start[0])
	startMarker = url.QueryEscape(startMarker)

	end := coords[len(coords)-1]
	endMarker := fmt.Sprintf("%f,%f", end[1], end[0])
	endMarker = url.QueryEscape(endMarker)

	if HaversineDistance(start[0], start[1], end[0], end[1]) > 60 {
		urlPath = fmt.Sprintf(
			"https://api.mapbox.com/styles/v1/mapbox/outdoors-v12/static/path-2+1a5fb4(%s),pin-s-pitch+003a1d(%s),pin-s-racetrack+003a1d(%s)/%f,%f,%.2f,0/%.0fx%.0f?access_token=%s",
			path,
			startMarker,
			endMarker,
			centerLon,
			centerLat,
			zoom,
			width,
			height,
			token,
		)
	} else {
		urlPath = fmt.Sprintf(
			"https://api.mapbox.com/styles/v1/mapbox/outdoors-v12/static/path-2+1a5fb4(%s),pin-s-pitch+003a1d(%s)/%f,%f,%.2f,0/%.0fx%.0f?access_token=%s",
			path,
			startMarker,
			centerLon,
			centerLat,
			zoom,
			width,
			height,
			token,
		)
	}

	resp, err := http.Get(urlPath)
	if err != nil {
		return "", fmt.Errorf("mapbox error request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return "", fmt.Errorf("mapbox error response: %s", resp.Status)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*50)
	defer cancel()

	bucket := client.Bucket(bucketName)

	objectPath := fmt.Sprintf("maps/%s.png", storagePath)

	wc := bucket.Object(objectPath).NewWriter(ctx)
	wc.ContentType = "image/png"
	wc.Metadata = map[string]string{
		"source": "mapbox-static",
	}

	if _, err = io.Copy(wc, resp.Body); err != nil {
		return "", fmt.Errorf("errore upload map su firebase: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("errore chiusura writer firebase: %w", err)
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectPath)
	return publicURL, nil
}

func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371e3

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func FindFileByID(id string) (string, error) {

	files, err := os.ReadDir("avatars")
	if err != nil {
		return "", fmt.Errorf("error reading profile directory: %w", err)
	}

	for _, file := range files {
		if file.Name() == id+".jpg" || file.Name() == id+".jpeg" || file.Name() == id+".png" {
			return "avatars/" + file.Name(), nil
		}
	}
	return "", nil

}
