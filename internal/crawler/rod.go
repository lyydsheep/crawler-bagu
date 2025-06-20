package crawler

import (
	"bytes"
	"crawler-bagu/pkg"
	"encoding/json"
	"errors"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

type Rod struct {
	b *rod.Browser
}

var (
	dir string
)

func (r *Rod) Consume(url string) error {
	log.Info().Msgf("consume url: %s.", url)
	r.b = rod.New().NoDefaultDevice().MustConnect().SlowMotion(10 * time.Second)
	defer r.b.Close()

	if err := SetCookiesFromString(r.b, os.Getenv("Cookie")); err != nil {
		log.Error().Msgf("SetCookiesFromString error: %v.", err)
		return err
	}

	// 访问页面
	page := r.b.MustPage(url)
	defer page.Close()

	// 获取标题
	_, element, _ := page.Has(`head > title`)
	title, _ := element.Text()
	d, err := pkg.Touch(title)
	if err != nil {
		log.Error().Msgf("Touch error: %v.", err)
		return err
	}
	dir = d
	log.Debug().Msgf("title is %s.", title)

	// 创建拦截器
	ch := make(chan struct{})
	go hiJack(page, ch)
	selector := `body > div.content-container > div.page-top > div:nth-child(2) > div > div > div > div > div.w-full.mb-12.lg\:w-\[calc\(100\%-21rem\)\] > div:nth-child(1) > div.sticky.bottom-2.left-0.right-0.rounded-lg.h-16.shadow.flex.items-center.justify-between.p-4.backdrop-blur-sm.border.bg-gray-100.bg-opacity-85.dark\:border-dark-nav.dark\:bg-dark-nav > div:nth-child(1) > button`
	// 判断是真题，循环点击
	for ok, err := check(page); ok; ok, err = check(page) {
		if err != nil {
			log.Error().Msgf("check error: %v.", err)
			return err
		}
		if err := click(selector, page, ch); err != nil {
			log.Error().Msgf("click error: %v.", err)
			return err
		}
		selector = `body > div.content-container > div.page-top > div:nth-child(2) > div > div > div > div > div.w-full.mb-12.lg\:w-\[calc\(100\%-21rem\)\] > div:nth-child(1) > div.sticky.bottom-2.left-0.right-0.rounded-lg.h-16.shadow.flex.items-center.justify-between.p-4.backdrop-blur-sm.border.bg-gray-100.bg-opacity-85.dark\:border-dark-nav.dark\:bg-dark-nav > div:nth-child(1) > button:nth-child(2)`
	}

	return nil
}

func check(page *rod.Page) (bool, error) {
	// 判断是否是真题
	exists, _, err := page.Has("svg[data-icon='building']")
	if err != nil {
		log.Error().Msgf("page.Has error: %v.", err)
		return true, err
	}
	return exists, nil
}

func click(selector string, page *rod.Page, ch chan struct{}) error {
	// 判断是否存在该按钮
	exists, _, err := page.Has(selector)
	if err != nil {
		panic(err)
	}

	if exists {
		log.Info().Msg("found the button, clicking...")
		page.MustElement(selector).MustClick()
		log.Info().Msgf("click successfully.")
		<-ch
		return nil
	}
	return errors.New("not found button")
}

func hiJack(page *rod.Page, ch chan struct{}) {
	router := page.HijackRequests()
	defer router.MustStop()
	targetUrl := "*/api/problems/*?groupId*"
	router.MustAdd(targetUrl, func(ctx *rod.Hijack) {
		ctx.MustLoadResponse()

		var resp Response
		if err := json.NewDecoder(bytes.NewReader([]byte(ctx.Response.Body()))).Decode(&resp); err != nil {
			log.Error().Msgf("json.NewDecoder error: %v.", err)
			return
		}

		// 保存题目和答案
		go save(resp)

		// 完成工作
		ch <- struct{}{}
	})
	router.Run()
}

func save(resp Response) {
	var (
		filename string
		ok       bool
		err      error
	)

	if filename, ok = checkExist(dir, resp.Data.GroupId, resp.Data.Id); ok {
		// 已存在
		log.Info().Msgf("file %s already exists.", filename)
		return
	}

	data, err := json.Marshal(resp.Data)
	if err != nil {
		log.Error().Msgf("json.Marshal error: %v", err)
		return
	}

	if err = saveAnswer(filename, data); err != nil {
		log.Error().Msgf("saveAnswer error: %v", err)
		return
	}

	log.Info().Msgf("saved answer to %s", filename)
}

func SetCookiesFromString(b *rod.Browser, cookieStr string) error {
	cookies, err := parseCookies(cookieStr)
	if err != nil {
		return err
	}
	b.MustSetCookies(cookies...)
	return nil
}

// parseCookies 解析 Cookie 字符串为 []*rod.Cookie
func parseCookies(cookieStr string) ([]*proto.NetworkCookie, error) {
	header := http.Header{}
	header.Add("Cookie", cookieStr)
	req := &http.Request{Header: header}
	log.Debug().Msgf("cookies: %v", req.Cookies())

	cookies := make([]*proto.NetworkCookie, 0)
	for _, c := range req.Cookies() {
		cookies = append(cookies, &proto.NetworkCookie{
			Name:   c.Name,
			Value:  c.Value,
			Domain: ".bagujing.com", // 根据实际域名设置
		})
	}
	return cookies, nil
}

func NewRod() Crawler {
	return &Rod{
		// 无头浏览器
		//b: rod.New().MustConnect(),
	}
}
