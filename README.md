GIN-powered blog boilerplate
===============

*Update in progress*

Provides essentials that most web blogs need - MVC pattern, user authorisation, SQL db migration, admin dashboard, javascript form validation, rss feeds, etc.

It consists of the following core components:

- GIN - A web microframework (with best performance atm) for Golang - https://github.com/gin-gonic/gin
- GIN middlewares [gin-csrf](https://github.com/utrack/gin-csrf), [gin/contrib/sessions](https://github.com/gin-gonic/contrib/tree/master/sessions)
- gorm - The orm library for go - http://gorm.io/
- Comments with oauth2 authentication
- logrus - advanced Go logger - https://github.com/Sirupsen/logrus
- Twitter Bootstrap 4 - popular HTML, CSS, JS framework for developing responsive, mobile first web projects - http://getbootstrap.com
- Gulp asset compiler
- Parsley JS - form validation - http://parsleyjs.org
- CKEditor 5 with image upload - https://ckeditor.com/ckeditor-5/
- bluemonday - html sanitizer (for excerpts, etc) - https://github.com/microcosm-cc/bluemonday
- RSS feeds - https://github.com/gorilla/feeds
- sitemap - XML sitemap for search engines - https://github.com/denisbakhtin/sitemap
- gocron - periodic task launcher (for sitemap generation, etc) - https://github.com/jasonlvhit/gocron

# TODO (May be)
- Site search with Postgresql full text search (okish for most websites) - http://www.postgresql.org/docs/9.4/static/textsearch-intro.html
- Social plugins (share, like buttons)
- Auto posting previews to social walls

# Screenshots (some may be outdated)
## Home page
![](/public/images/screenshot_home.jpg)
## Dashboard
![](/public/images/screenshot_dashboard.jpg)
## Markdown editor
![](/public/images/screenshot_markdown.jpg)
## Fancy 404, 405, 500 error pages
![](/public/images/screenshot_error.jpg)

# Usage
```
git clone https://github.com/denisbakhtin/ginblog.git
cd ginblog
go get .
```
Copy sample config `cp config/config.json.example config/config.json`, create postgresql database, modify config/config.json accordingly.
Install `npm`, `gulp`, run `npm install` in the project directory.

Type `go run main.go` to launch web server, `gulp` to rebuild assets.

# Deployment
```
make build
```
Upload `ginblog` binary `config`, `views` and `public` directory to your server.

# Project structure

`/config`

Contains application configuration file.

`/controllers`

All your controllers that serve defined routes.

`/models`

You database models.

`/public`

It contains all web-site static files

`/system`

Core functions and structs.

`/views`

Your views using standard `Go` template system.

`main.go`

This file starts your web application, contains routes definition & some custom middlewares.

# Make it your own

I assume you have followed installation instructions and you have `ginblog` installed in your `GOPATH` location.

Let's say I want to create `Amazing Website`. I create new `GitHub` repository `https://github.com/denisbakhtin/amazingblog` (of course replace that with your own repository).

Now I have to prepare `ginblog`. First thing is that I have to delete its `.git` directory.

I issue:

```
rm -rf src/github.com/denisbakhtin/ginblog/.git
```

Then I want to replace all references from `github.com/denisbakhtin/amazingblog` to `github.com/denisbakhtin/amazingblog`:

```
grep -rl 'github.com/denisbakhtin/ginblog' ./ | xargs sed -i 's/github.com\/denisbakhtin\/ginblog/github.com\/denisbakhtin\/amazingblog/g'
```

Now I have to move all `ginblog` files to the new location:

```
mv src/github.com/denisbakhtin/ginblog/ src/github.com/denisbakhtin/amazingblog
```

And push it to my new repository at `GitHub`:

```
cd src/github.com/denisbakhtin/amazingblog
git init
git add --all .
git commit -m "Amazing Blog First Commit"
git remote add origin https://github.com/denisbakhtin/amazingblog.git
git push -u origin master
```

You can now go back to your `GOPATH` and check if everything is ok:

```
go install github.com/denisbakhtin/amazingblog
```

And that's it.

# Continuous Development

For Continuous Development a good option is `reflex` - https://github.com/cespare/reflex
Then simply run `reflex -c reflex.conf` in the project directory.

Or run `realize s` (works on Windows unlike reflex) - https://github.com/oxequa/realize.

To rebuild assets on change install `npm`, run `npm install` and then `npm run watch`. Run `npm run build` to build assets for production.