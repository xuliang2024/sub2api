package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

func requestedOpenAIImageDimensions(size string) (int, int, bool) {
	upstreamSize := normalizeOpenAIImageUpstreamSize(size)
	width, height, ok := parseImageBillingDimensions(upstreamSize)
	if !ok || width <= 0 || height <= 0 {
		return 0, 0, false
	}
	return width, height, true
}

func normalizeOpenAIImagesResultDimensions(results []openAIResponsesImageResult, requestedSize string) []openAIResponsesImageResult {
	width, height, ok := requestedOpenAIImageDimensions(requestedSize)
	if !ok {
		return results
	}
	out := make([]openAIResponsesImageResult, len(results))
	for i, result := range results {
		out[i] = normalizeOpenAIImageResultDimensions(result, width, height)
	}
	return out
}

func normalizeOpenAIImageResultDimensions(result openAIResponsesImageResult, width int, height int) openAIResponsesImageResult {
	if width <= 0 || height <= 0 || strings.TrimSpace(result.Result) == "" {
		return result
	}
	targetSize := fmt.Sprintf("%dx%d", width, height)
	raw, err := base64.StdEncoding.DecodeString(result.Result)
	if err != nil {
		result.Size = targetSize
		return result
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil {
		result.Size = targetSize
		return result
	}
	if cfg.Width == width && cfg.Height == height {
		result.Size = targetSize
		return result
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		result.Size = targetSize
		return result
	}
	resized := resizeImageToExactDimensions(img, width, height)
	encoded, outputFormat, err := encodeOpenAIImageResult(resized, format, result.OutputFormat)
	if err != nil {
		result.Size = targetSize
		return result
	}
	result.Result = base64.StdEncoding.EncodeToString(encoded)
	result.Size = targetSize
	result.OutputFormat = outputFormat
	return result
}

func resizeImageToExactDimensions(src image.Image, targetWidth int, targetHeight int) image.Image {
	bounds := src.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()
	if srcWidth <= 0 || srcHeight <= 0 || targetWidth <= 0 || targetHeight <= 0 {
		return src
	}

	crop := bounds
	if srcWidth*targetHeight > srcHeight*targetWidth {
		cropWidth := srcHeight * targetWidth / targetHeight
		if cropWidth > 0 && cropWidth < srcWidth {
			left := bounds.Min.X + (srcWidth-cropWidth)/2
			crop = image.Rect(left, bounds.Min.Y, left+cropWidth, bounds.Max.Y)
		}
	} else if srcWidth*targetHeight < srcHeight*targetWidth {
		cropHeight := srcWidth * targetHeight / targetWidth
		if cropHeight > 0 && cropHeight < srcHeight {
			top := bounds.Min.Y + (srcHeight-cropHeight)/2
			crop = image.Rect(bounds.Min.X, top, bounds.Max.X, top+cropHeight)
		}
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, crop, draw.Over, nil)
	return dst
}

func encodeOpenAIImageResult(img image.Image, decodedFormat string, outputFormat string) ([]byte, string, error) {
	format := strings.ToLower(strings.TrimSpace(outputFormat))
	if format == "" {
		format = strings.ToLower(strings.TrimSpace(decodedFormat))
	}
	var buf bytes.Buffer
	switch format {
	case "jpg", "jpeg":
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95}); err != nil {
			return nil, "", err
		}
		return buf.Bytes(), "jpeg", nil
	default:
		if err := png.Encode(&buf, img); err != nil {
			return nil, "", err
		}
		return buf.Bytes(), "png", nil
	}
}
