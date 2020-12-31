package lawn

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	endpointFmt = "https://github.com/users/%s/contributions"
	className   = "js-calendar-graph-svg"
)

// Client represents GitHub contributions graph scraper client.
type Client struct {
	cli *http.Client
}

// NewClient returns a new client.
func NewClient() *Client {
	cli := &http.Client{
		Timeout: 15 * time.Second,
	}
	return &Client{
		cli: cli,
	}
}

// Fetch fetches github contributions graph svg.
func (c *Client) Fetch(ctx context.Context, w io.Writer, username string) error {
	r, err := c.fetch(ctx, username)
	if err != nil {
		return err
	}
	defer r.Close()
	if err := c.parse(w, r); err != nil {
		return err
	}
	return nil
}

func (c *Client) fetch(ctx context.Context, username string) (io.ReadCloser, error) {
	uri := fmt.Sprintf(endpointFmt, username)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func (c *Client) parse(w io.Writer, r io.Reader) error {
	node, err := html.Parse(r)
	if err != nil {
		return err
	}
	svg, ok := c.parser(node)
	if !ok {
		return errors.New("svg image is not found")
	}
	n := c.formatSVG(svg)
	if err := html.Render(w, n); err != nil {
		return err
	}
	return nil
}

func (c *Client) parser(n *html.Node) (*html.Node, bool) {
	if n.Type == html.ElementNode && n.DataAtom == atom.Svg && c.hasClass(n.Attr, className) {
		return n, true
	}
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		if n, ok := c.parser(ch); ok {
			return n, ok
		}
	}
	return nil, false
}

func (c *Client) hasClass(attrs []html.Attribute, class string) bool {
	for _, attr := range attrs {
		if attr.Key != "class" {
			continue
		}

		return strings.Contains(attr.Val, class)
	}

	return false
}

func (c *Client) formatSVG(svg *html.Node) *html.Node {
	doc := &html.Node{
		Type: html.DocumentNode,
	}
	doc.AppendChild(&html.Node{
		Type: html.CommentNode,
		Data: `?xml version="1.0" standalone="no"?`,
	})
	doc.AppendChild(&html.Node{
		Type: html.DoctypeNode,
		Data: "svg",
		Attr: []html.Attribute{
			{
				Key: "public",
				Val: "-//W3C//DTD SVG 1.1//EN",
			},
			{
				Key: "system",
				Val: "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd",
			},
		},
	})
	svg.Parent.RemoveChild(svg)
	doc.AppendChild(svg)
	svg.Attr = append(svg.Attr, []html.Attribute{
		{
			Key: "xmlns",
			Val: "http://www.w3.org/2000/svg",
		},
		{
			Key: "version",
			Val: "1.1",
		},
	}...)
	defs := &html.Node{
		Type: html.ElementNode,
		Data: "defs",
	}
	svg.AppendChild(defs)
	style := &html.Node{
		Type:     html.ElementNode,
		Data:     "style",
		DataAtom: atom.Style,
		Attr: []html.Attribute{
			{
				Key: "type",
				Val: "text/css",
			},
		},
	}
	defs.AppendChild(style)
	style.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: `<![CDATA[
text.month { font-size: 10px; fill: #767676 }
text.wday { font-size: 9px; fill: #767676 }
svg {
	--color-calendar-graph-day-bg:#ebedf0;
	--color-calendar-graph-day-border:rgba(27,31,35,0.06);
	--color-calendar-graph-day-L1-bg:#9be9a8;
	--color-calendar-graph-day-L2-bg:#40c463;
	--color-calendar-graph-day-L3-bg:#30a14e;
	--color-calendar-graph-day-L4-bg:#216e39;
	--color-calendar-graph-day-L4-border:rgba(27,31,35,0.06);
	--color-calendar-graph-day-L3-border:rgba(27,31,35,0.06);
	--color-calendar-graph-day-L2-border:rgba(27,31,35,0.06);
	--color-calendar-graph-day-L1-border:rgba(27,31,35,0.06);
}
]]>`,
	})
	return doc
}
