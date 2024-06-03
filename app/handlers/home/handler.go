package home

import (
	"html/template"
	"net/http"
	"strings"

	lru "github.com/hashicorp/golang-lru/v2"
	"go.uber.org/zap"

	"gohome.4gophers.ru/gohome/pkg/logger"
)

type Handler struct {
	logger   *logger.Logger
	template *template.Template
	cache    *lru.Cache[string, Meta]
}

func (h *Handler) save(path string) {
	if strings.Contains(path, ".") {
		return
	}

	if strings.Count(path, "/") <= 1 {
		return
	}

	h.cache.Add(path, Meta{
		Name: "gohome.4gophers.ru" + path,
		Repo: "https://gitflic.ru/project" + path + ".git",
	})
}

func New(logger *logger.Logger) *Handler {
	l, _ := lru.New[string, Meta](1024)

	t, err := template.New("test").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	return &Handler{
		logger:   logger,
		template: t,
		cache:    l,
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("url path", zap.String("path", r.URL.Path))

	path := r.URL.Path

	h.save(path)

	data := h.cache.Values()

	if err := h.template.Execute(w, data); err != nil {
		h.logger.Error("error on render template", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type Meta struct {
	Name string
	Repo string
}

var tmpl = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">

	{{range .}}
	<meta name="go-import" content="{{.Name}} git {{.Repo}}">
	{{ end }}
    <title>Go Home</title>    

    <!-- Bootstrap core CSS -->
<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC" crossorigin="anonymous">


    <!-- Favicons -->
<link rel="apple-touch-icon" href="/docs/5.0/assets/img/favicons/apple-touch-icon.png" sizes="180x180">
<link rel="icon" href="/docs/5.0/assets/img/favicons/favicon-32x32.png" sizes="32x32" type="image/png">
<link rel="icon" href="/docs/5.0/assets/img/favicons/favicon-16x16.png" sizes="16x16" type="image/png">
<link rel="manifest" href="/docs/5.0/assets/img/favicons/manifest.json">
<link rel="mask-icon" href="/docs/5.0/assets/img/favicons/safari-pinned-tab.svg" color="#7952b3">
<link rel="icon" href="/docs/5.0/assets/img/favicons/favicon.ico">
<meta name="theme-color" content="#7952b3">


    <style>
      .bd-placeholder-img {
        font-size: 1.125rem;
        text-anchor: middle;
        -webkit-user-select: none;
        -moz-user-select: none;
        user-select: none;
      }

      @media (min-width: 768px) {
        .bd-placeholder-img-lg {
          font-size: 3.5rem;
        }
      }
    </style>


	<style>
	.icon-list {
		padding-left: 0;
		list-style: none;
	  }
	  .icon-list li {
		display: flex;
		align-items: flex-start;
		margin-bottom: .25rem;
	  }
	  .icon-list li::before {
		display: block;
		flex-shrink: 0;
		width: 1.5em;
		height: 1.5em;
		margin-right: .5rem;
		content: "";
		background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='%23212529' viewBox='0 0 16 16'%3E%3Cpath d='M8 0a8 8 0 1 1 0 16A8 8 0 0 1 8 0zM4.5 7.5a.5.5 0 0 0 0 1h5.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3a.5.5 0 0 0 0-.708l-3-3a.5.5 0 1 0-.708.708L10.293 7.5H4.5z'/%3E%3C/svg%3E") no-repeat center center / 100% auto;
	  }
	</style>
  </head>
  <body>
    
<div class="col-lg-8 mx-auto p-3 py-md-5">
  <header class="d-flex align-items-center pb-3 mb-5 border-bottom">
    <a href="/" class="d-flex align-items-center text-dark text-decoration-none">
      <span class="fs-4">Go Home</span>
    </a>
  </header>

  <main>
    <h1>Домен для Go пакетов</h1>
	<p></p>
    <p class="fs-5 col-md-10">
	  Очень легко устанавливать Go пакеты с гитхаба. Но все сложнее если библиотека лежит на gitflic.ru. Тут обычный go get не сработает.
	</p>
	<p class="fs-5 col-md-10">
	  Для этого я сделала <code>Go Home</code> с помощью которого пакеты можно устанавливать как <code>go get gohome.4gophers.ru/getapp/boosty</code> и установить пакет.
	</p>

	<p class="fs-5 col-md-10">
	  Для начала нужно создать пакет, название которого будет начинаться с gohome.4gophers.ru:

<pre>
go mod init gohome.4gophers.ru/kovardin/example
</pre>
	
	</p>
	<p class="fs-5 col-md-10">  
	  После этого создать репозиторий на gitflic: 
<pre>
https://gitflic.ru/project/kovardin/example
</pre>

	</p>

	<p class="fs-5 col-md-8">
	Теперь вы можете указывать в своем коде зависимость

<pre>
package main

import (
	"fmt"

	"gohome.4gophers.ru/kovardin/example"
)

func main() {
	fmt.Println(example.Hello("Artem"))
}
</pre>
	</p>

    <hr class="col-3 col-md-2 mb-5">

    <div class="row g-5">
      <div class="col-md-6">
        <h2>Ссылки</h2>
        <p>Как работает и что за giflic такой.</p>
        <ul class="icon-list">
          <li><a href="https://pkg.go.dev/cmd/go#hdr-Remote_import_paths" rel="noopener" target="_blank">Как работает импорт пакетов</a></li>
          <li><a href="https://gitflic.ru/user/kovardin" rel="noopener" target="_blank">Да кто это ваш giflic</a></li>
          <li><a href="https://t.me/kodikapusta" rel="noopener" target="_blank">Код и капуста</a></li>
        </ul>
      </div>

      <div class="col-md-6">
        <h2>Статьи</h2>
        <p>Еще немножко почитать на разные темы.</p>
        <ul class="icon-list">
          <li><a href="https://kovardin.ru/articles/ads/mediation/">Про мобильную медиацию</a></li>
          <li><a href="https://kovardin.ru/articles/mobile/boosty-android/">Как использовать boosty для мобильного приложения</a></li>
          <li><a href="https://kovardin.ru/articles/godot/mytacker-and-app-metrica/">Подключаем MyTracker и AppMetrica к игре на Godot</a></li>
          <li><a href="https://kovardin.ru/articles/godot/first-game/">Первая игра на Godot</a></li>
        </ul>
      </div>
    </div>
  </main>
  <footer class="pt-5 my-5 text-muted border-top">
    Ковардин Артем &middot; &copy; 2024
  </footer>
</div>


<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/js/bootstrap.bundle.min.js" integrity="sha384-MrcW6ZMFYlzcLA8Nl+NtUVF0sA7MsXsP1UyJoMp4YLEuNSfAP+JcXn/tWtIaxVXM" crossorigin="anonymous"></script>
      
  </body>
</html>


`
