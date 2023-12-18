package githuberrors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
)

const (
	SecondaryRateLimitMessage                 = `You have exceeded a secondary rate limit`
	SecondaryRateLimitDocumentationPathSuffix = `secondary-rate-limits`
	HeaderXRateLimitRemaining                 = "x-ratelimit-remaining"
	HeaderRetryAfter                          = "retry-after"
	HeaderXRateLimitReset                     = "x-ratelimit-reset"
)

func IsSecondaryRateLimitError(r *github.Response) (bool, *SecondaryRateLimitError) {
	var body *SecondaryRateLimitBody
	res := r.Response

	if isRateLimit, b := isSecondaryRateLimit(res); !isRateLimit {
		return false, nil
	} 
	body = b

	retryAfter := parseSecondaryLimitTime(res)
	return true, &SecondaryRateLimitError{
		Body:       body,
		RetryAfter: retryAfter,
		Response:   *res,
	}
}

type SecondaryRateLimitError struct {
	Response   http.Response
	Body       *SecondaryRateLimitBody
	RetryAfter *time.Time
}

func (s *SecondaryRateLimitError) Error() string {
	return fmt.Sprintf("%v %v: %d %v (%v)", s.Response.Request.Method, sanitizeURL(s.Response.Request.URL), s.Response.StatusCode, s.Body.Message, s.Body.DocumentURL)
}

type SecondaryRateLimitBody struct {
	Message     string `json:"message"`
	DocumentURL string `json:"documentation_url"`
}

// IsSecondaryRateLimit checks whether the response is a legitimate secondary rate limit.
// It checks the prefix of the message and the suffix of the documentation URL in the response body in case
// the message or documentation URL is modified in the future.
// https://docs.github.com/en/rest/overview/rate-limits-for-the-rest-api#about-secondary-rate-limits
func (s SecondaryRateLimitBody) IsSecondaryRateLimit() bool {
	return strings.HasPrefix(s.Message, SecondaryRateLimitMessage) || strings.HasSuffix(s.DocumentURL, SecondaryRateLimitDocumentationPathSuffix)
}

// isSecondaryRateLimit checks whether the response is a legitimate secondary rate limit.
// it is used to avoid handling primary rate limits and authentic HTTP Forbidden (403) responses.
func isSecondaryRateLimit(resp *http.Response) (bool, *SecondaryRateLimitBody) {
	if resp.StatusCode != http.StatusForbidden {
		return false, nil
	}

	if resp.Header == nil {
		return false, nil
	}

	// a primary rate limit
	if remaining, ok := httpHeaderIntValue(resp.Header, HeaderXRateLimitRemaining); ok && remaining == 0 {
		return false, nil
	}

	// an authentic HTTP Forbidden (403) response
	defer resp.Body.Close()
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, nil // unexpected error
	}

	// restore original body
	resp.Body = io.NopCloser(bytes.NewReader(rawBody))

	var body SecondaryRateLimitBody
	if err := json.Unmarshal(rawBody, &body); err != nil {
		return false, nil // unexpected error
	}
	if !body.IsSecondaryRateLimit() {
		return false, nil
	}

	return true, &body
}

func httpHeaderIntValue(header http.Header, key string) (int64, bool) {
	val := header.Get(key)
	if val == "" {
		return 0, false
	}
	asInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, false
	}
	return asInt, true
}

// parseSecondaryLimitTime parses the GitHub API response header,
// looking for the secondary rate limit as defined by GitHub API documentation.
// https://docs.github.com/en/rest/overview/resources-in-the-rest-api#secondary-rate-limits
func parseSecondaryLimitTime(resp *http.Response) *time.Time {
	if isRateLimit, _ := isSecondaryRateLimit(resp); !isRateLimit {
		return nil
	}

	if sleepUntil := parseRetryAfter(resp.Header); sleepUntil != nil {
		return sleepUntil
	}

	if sleepUntil := parseXRateLimitReset(resp); sleepUntil != nil {
		return sleepUntil
	}

	return nil
}

// parseRetryAfter parses the GitHub API response header in case a Retry-After is returned.
func parseRetryAfter(header http.Header) *time.Time {
	retryAfterSeconds, ok := httpHeaderIntValue(header, "retry-after")
	if !ok || retryAfterSeconds <= 0 {
		return nil
	}

	// per GitHub API, the header is set to the number of seconds to wait
	sleepUntil := time.Now().Add(time.Duration(retryAfterSeconds) * time.Second)

	return &sleepUntil
}

// parseXRateLimitReset parses the GitHub API response header in case a x-ratelimit-reset is returned.
// to avoid handling primary rate limits (which are categorized),
// we only handle x-ratelimit-reset in case the primary rate limit is not reached.
func parseXRateLimitReset(resp *http.Response) *time.Time {
	secondsSinceEpoch, ok := httpHeaderIntValue(resp.Header, HeaderXRateLimitReset)
	if !ok || secondsSinceEpoch <= 0 {
		return nil
	}

	// per GitHub API, the header is set to the number of seconds since epoch (UTC)
	sleepUntil := time.Unix(secondsSinceEpoch, 0)

	return &sleepUntil
}

// sanitizeURL redacts the client_secret parameter from the URL which may be
// exposed to the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("client_secret")) > 0 {
		params.Set("client_secret", "REDACTED")
		uri.RawQuery = params.Encode()
	}
	return uri
}
