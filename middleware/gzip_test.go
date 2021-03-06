package middleware

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aphistic/sweet"
	"github.com/efritz/response"
	"github.com/go-nacelle/nacelle"
	. "github.com/onsi/gomega"
)

type GzipSuite struct{}

func (s *GzipSuite) TestCompress(t sweet.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		resp := response.Respond([]byte(shakespeare))
		resp.AddHeader("Content-Length", fmt.Sprintf("%d", len(shakespeare)))
		return resp
	}

	wrapped, err := NewGzip().Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	r.Header.Add("Accept-Encoding", "deflate/gzip")

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	headers, body, err := response.Serialize(resp)
	Expect(err).To(BeNil())

	Expect(headers).To(Equal(http.Header{
		"Content-Encoding": []string{"gzip"},
		"Vary":             []string{"Accept-Encoding"},
		"Content-Type":     []string{"application/octet-stream"},
	}))

	content, err := readGzipContent(body)
	Expect(err).To(BeNil())
	Expect(string(content)).To(Equal(shakespeare))
}
func (s *GzipSuite) TestNoContent(t sweet.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		resp := response.Respond(nil)
		resp.SetStatusCode(http.StatusNoContent)
		return resp
	}

	wrapped, err := NewGzip().Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	r.Header.Add("Accept-Encoding", "deflate/gzip")

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	headers, body, err := response.Serialize(resp)
	Expect(err).To(BeNil())

	Expect(headers).To(Equal(http.Header{
		"Vary":         []string{"Accept-Encoding"},
		"Content-Type": []string{"application/octet-stream"},
	}))

	content, err := readGzipContent(body)
	Expect(err).To(BeNil())
	Expect(string(content)).To(BeEmpty())
}

func (s *GzipSuite) TestCompressClosesWrappedReader(t sweet.T) {
	reader := &closeWrapper{
		Reader: strings.NewReader(shakespeare),
	}

	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.Stream(reader)
	}

	wrapped, err := NewGzip().Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	r.Header.Add("Accept-Encoding", "deflate/gzip")

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	_, _, err = response.Serialize(resp)
	Expect(err).To(BeNil())
	Expect(reader.closed).To(BeTrue())
}

func (s *GzipSuite) TestNonGzipAcceptEncodingDoesNotCompress(t sweet.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		return response.Respond([]byte(shakespeare))
	}

	wrapped, err := NewGzip().Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	headers, body, err := response.Serialize(resp)
	Expect(err).To(BeNil())

	Expect(headers).To(BeEmpty())
	Expect(body).To(Equal([]byte(shakespeare)))
}

func (s *GzipSuite) TestExplicitEncodingDoesNotCompress(t sweet.T) {
	bare := func(ctx context.Context, r *http.Request, logger nacelle.Logger) response.Response {
		resp := response.Respond([]byte(shakespeare))
		resp.AddHeader("Content-Encoding", "text/plain")

		return resp
	}

	wrapped, err := NewGzip().Convert(bare)
	Expect(err).To(BeNil())

	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	r.Header.Add("Accept-Encoding", "deflate, gzip")

	resp := wrapped(context.Background(), r, nacelle.NewNilLogger())
	headers, body, err := response.Serialize(resp)
	Expect(err).To(BeNil())

	Expect(headers).To(Equal(http.Header{"Content-Encoding": []string{"text/plain"}}))
	Expect(body).To(Equal([]byte(shakespeare)))
}

func (s *GzipSuite) TestInvalidLevel(t sweet.T) {
	_, err := NewGzip(WithGzipLevel(400)).Convert(nil)
	Expect(err).To(MatchError("gzip: invalid compression level: 400"))
}

//
// Helpers

func readGzipContent(body []byte) ([]byte, error) {
	r, err := gzip.NewReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(r)
}

//
//

type closeWrapper struct {
	io.Reader
	closed bool
}

func (w *closeWrapper) Close() error {
	w.closed = true
	return nil
}

var shakespeare = `
    From fairest creatures we desire increase,
    That thereby beauty's rose might never die,
    But as the riper should by time decease,
    His tender heir might bear his memory:
    But thou, contracted to thine own bright eyes,
    Feed'st thy light'st flame with self-substantial fuel,
    Making a famine where abundance lies,
    Thyself thy foe, to thy sweet self too cruel.
    Thou that art now the world's fresh ornament
    And only herald to the gaudy spring,
    Within thine own bud buriest thy content
    And, tender churl, makest waste in niggarding.
    Pity the world, or else this glutton be,
    To eat the world's due, by the grave and thee.

    When forty winters shall beseige thy brow,
    And dig deep trenches in thy beauty's field,
    Thy youth's proud livery, so gazed on now,
    Will be a tatter'd weed, of small worth held:
    Then being ask'd where all thy beauty lies,
    Where all the treasure of thy lusty days,
    To say, within thine own deep-sunken eyes,
    Were an all-eating shame and thriftless praise.
    How much more praise deserved thy beauty's use,
    If thou couldst answer 'This fair child of mine
    Shall sum my count and make my old excuse,'
    Proving his beauty by succession thine!
    This were to be new made when thou art old,
    And see thy blood warm when thou feel'st it cold.

    Look in thy glass, and tell the face thou viewest
    Now is the time that face should form another;
    Whose fresh repair if now thou not renewest,
    Thou dost beguile the world, unbless some mother.
    For where is she so fair whose unear'd womb
    Disdains the tillage of thy husbandry?
    Or who is he so fond will be the tomb
    Of his self-love, to stop posterity?
    Thou art thy mother's glass, and she in thee
    Calls back the lovely April of her prime:
    So thou through windows of thine age shall see
    Despite of wrinkles this thy golden time.
    But if thou live, remember'd not to be,
    Die single, and thine image dies with thee.

    Unthrifty loveliness, why dost thou spend
    Upon thyself thy beauty's legacy?
    Nature's bequest gives nothing but doth lend,
    And being frank she lends to those are free.
    Then, beauteous niggard, why dost thou abuse
    The bounteous largess given thee to give?
    Profitless usurer, why dost thou use
    So great a sum of sums, yet canst not live?
    For having traffic with thyself alone,
    Thou of thyself thy sweet self dost deceive.
    Then how, when nature calls thee to be gone,
    What acceptable audit canst thou leave?
    Thy unused beauty must be tomb'd with thee,
    Which, used, lives th' executor to be.

    Those hours, that with gentle work did frame
    The lovely gaze where every eye doth dwell,
    Will play the tyrants to the very same
    And that unfair which fairly doth excel:
    For never-resting time leads summer on
    To hideous winter and confounds him there;
    Sap cheque'd with frost and lusty leaves quite gone,
    Beauty o'ersnow'd and bareness every where:
    Then, were not summer's distillation left,
    A liquid prisoner pent in walls of glass,
    Beauty's effect with beauty were bereft,
    Nor it nor no remembrance what it was:
    But flowers distill'd though they with winter meet,
    Leese but their show; their substance still lives sweet.
`
