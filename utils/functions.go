package utils

import (
	"fmt"
	"github.com/tkrajina/gpxgo/gpx"
	"github.com/twpayne/go-polyline"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

func DeleteFiles(path Gpx) error {
	gpxPath := "gpxs/" + path.StoragePath
	mapsPath := "maps/" + path.StoragePath + ".png"

	if _, err := os.Stat(gpxPath); err == nil {
		err := os.Remove(gpxPath)
		if err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	} else if os.IsNotExist(err) {
		return nil
	} else {
		return err
	}

	if _, err := os.Stat(mapsPath); err == nil {
		err := os.Remove(mapsPath)
		if err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	} else if os.IsNotExist(err) {
		return nil
	} else {
		return err
	}

	return nil
}

func SaveFile(file *multipart.FileHeader, storagePath string) error {
	// Open the file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create the directory if it doesn't exist
	if err := os.MkdirAll("gpxs", os.ModePerm); err != nil {
		return fmt.Errorf("errore creazione cartella gpxs: %w", err)
	}

	// Create the destination file
	storagePath = "gpxs/" + storagePath
	dst, err := os.Create(storagePath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy the file
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
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

func CreateMap(file *multipart.FileHeader, storagePath string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("errore apertura file: %w", err)
	}
	defer src.Close()

	gpxData, err := gpx.Parse(src)
	if err != nil {
		return fmt.Errorf("errore parsing GPX: %w", err)
	}

	var coords [][]float64
	var sumLat, sumLon float64
	var count int

	for _, track := range gpxData.Tracks {
		for _, segment := range track.Segments {
			for _, point := range segment.Points {
				lat := point.Latitude
				lon := point.Longitude
				sumLat += lat
				sumLon += lon
				coords = append(coords, []float64{lat, lon})
				count++
			}
		}
	}

	if count == 0 {
		return fmt.Errorf("invalid GPX file: no coordinates found")
	}

	centerLat := sumLat / float64(count)
	centerLon := sumLon / float64(count)
	path := string(polyline.EncodeCoords(coords))
	path = url.QueryEscape(path)

	const zoom = 13
	const bearing = 0
	const width = 800
	const height = 800
	var token = os.Getenv("MAPBOX_TOKEN")
	var filePath = "./maps/" + storagePath + ".png"

	url := fmt.Sprintf(
		"https://api.mapbox.com/styles/v1/mapbox/outdoors-v12/static/path-2+1a5fb4(%s)/%f,%f,%d,%d/%dx%d?access_token=%s",
		path,
		centerLon,
		centerLat,
		zoom,
		bearing,
		width,
		height,
		token,
	)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("mapbox error request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("mapbox error response: %s", resp.Status)
	}

	if err := os.MkdirAll("./maps", os.ModePerm); err != nil {
		return fmt.Errorf("error creating maps directory: %w", err)
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating map file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying map data: %w", err)
	}

	return nil
}
