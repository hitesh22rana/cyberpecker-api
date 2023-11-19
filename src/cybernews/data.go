package cybernews

type News struct {
	Id       string `json:"id"`
	Source   string `json:"source"`
	Headline string `json:"headline"`
	News     string `json:"news"`
	Link     string `json:"link"`
	Image    string `json:"image"`
}

type NewsFields struct {
	source   string
	url      string
	headline string
	news     string
	link     string
	image    string
}

var newsCategory = map[string][]NewsFields{
	"general": {
		{
			source:   "The Economic Times",
			url:      "https://ciosea.economictimes.indiatimes.com/news/next-gen-technologies",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     "article.desc figure a",
			image:    "article.desc figure a img",
		},
		{
			source:   "The Economic Times",
			url:      "https://telecom.economictimes.indiatimes.com/news/internet",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     "article.desc figure a",
			image:    "article.desc figure a img",
		},
		{
			source:   "The Economic Times",
			url:      "https://ciosea.economictimes.indiatimes.com/news/consumer-tech",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
	"dataBreach": {
		{
			source:   "The Hacker News",
			url:      "https://thehackernews.com/search/label/data%20breach",
			headline: "h2.home-title",
			news:     ".home-desc",
			link:     "a.story-link",
			image:    ".img-ratio img",
		},
		{
			source:   "The Economic Times",
			url:      "https://ciso.economictimes.indiatimes.com/news/data-breaches",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
	"cyberAttack": {
		{
			source:   "The Hacker News",
			url:      "https://thehackernews.com/search/label/Cyber%20Attack",
			headline: "h2.home-title",
			news:     ".home-desc",
			link:     "a.story-link",
			image:    ".img-ratio img",
		},
		{
			source:   "The Economic Times",
			url:      "https://ciso.economictimes.indiatimes.com/news/cybercrime-fraud",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
	"vulnerability": {
		{
			source:   "The Hacker News",
			url:      "https://thehackernews.com/search/label/Vulnerability",
			headline: "h2.home-title",
			news:     ".home-desc",
			link:     "a.story-link",
			image:    ".img-ratio img",
		},
		{
			source:   "The Economic Times",
			url:      "https://ciso.economictimes.indiatimes.com/news/vulnerabilities-exploits",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
	"malware": {
		{
			source:   "The Hacker News",
			url:      "https://thehackernews.com/search/label/Malware",
			headline: "h2.home-title",
			news:     ".home-desc",
			link:     "a.story-link",
			image:    ".img-ratio img",
		},
	},
	"security": {
		{
			source:   "The Economic Times",
			url:      "https://ciosea.economictimes.indiatimes.com/news/security",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
		{
			source:   "The Economic Times",
			url:      "https://telecom.economictimes.indiatimes.com/tag/hacking",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
	"cloud": {
		{
			source:   "The Economic Times",
			url:      "https://ciosea.economictimes.indiatimes.com/news/cloud-computing",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
	"bigData": {
		{
			source:   "The Economic Times",
			url:      "https://ciosea.economictimes.indiatimes.com/news/big-data",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
		{
			source:   "The Economic Times",
			url:      "https://ciosea.economictimes.indiatimes.com/news/data-center",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
	"research": {
		{
			source:   "The Economic Times",
			url:      "https://ciosea.economictimes.indiatimes.com/tag/research",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
	"socialMedia": {
		{
			source:   "The Economic Times",
			url:      "https://telecom.economictimes.indiatimes.com/search/social",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
	"corporate": {
		{
			source:   "The Economic Times",
			url:      "https://ciosea.economictimes.indiatimes.com/news/corporate",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
		{
			source:   "The Economic Times",
			url:      "https://telecom.economictimes.indiatimes.com/news/industry",
			headline: "article.desc div h3.heading",
			news:     "article.desc div p.desktop-view",
			link:     ".desc figure a",
			image:    ".desc figure a img",
		},
	},
}
