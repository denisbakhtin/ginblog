GIN-powered blog boilerplate
===============

This is a skeleton project that provides essentials most web blogs need - MVC pattern, user authorisation, orm, admin dashboard, javascript form validation, rss feeds, etc.

It consists of the following core components:

- GIN - A web microframework (with best performance atm) for Golang - https://github.com/gin-gonic/gin
- GIN middlewares [gin-csrf](https://github.com/utrack/gin-csrf), [gin/contrib/sessions](https://github.com/gin-gonic/contrib/tree/master/sessions)
- gorm - The orm library for go - http://gorm.io/
- Comments with oauth2 authentication
- logrus - advanced Go logger - https://github.com/Sirupsen/logrus
- Twitter Bootstrap 4 - popular HTML, CSS, JS framework for developing responsive, mobile first web projects - http://getbootstrap.com
- Webpack asset compiler
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

# Screenshots
## Home page
![](/public/images/screenshot_home.jpg)
## Blog post
![](/public/images/screenshot_post.jpg)
## Dashboard
![](/public/images/screenshot_dashboard.jpg)
## Ckeditor 5 WYSIWYG editor
![](/public/images/screenshot_editor.jpg)
## Custom 404, 405, 500 error pages
![](/public/images/screenshot_error.jpg)

# Usage
```
git clone https://github.com/denisbakhtin/ginblog.git
cd ginblog
go get .
```
Copy sample config `cp config/config.json.example config/config.json`, create postgresql database, modify config/config.json accordingly.
Install `npm`, `webpack`, run `npm install` in the project directory.

Type `go run main.go` to launch web server, `npm run build` to rebuild assets.

# Deployment
```
npm run build
make build
```
Upload `ginblog` binary `config`, `views` and `public` directory to your server.

# Project structure

`/config`

Contains application configuration file & go wrapper.

`/controllers`

All your controllers that serve defined routes.

`/models`

You database models.

`/public`

It contains all web-site static files

`/views`

Your views using standard `Go` template system.

`main.go`

This file starts your web application, contains routes definition & some custom middlewares.

# Make it your own

I assume you have followed installation instructions and you have `ginblog` installed in your `GOPATH` location.

Let's say you want to create `Amazing Website`. Add a new `GitHub` repository `https://github.com/denisbakhtin/amazingblog` (of course replace that with your own repo).

Prepare `ginblog`: delete its `.git` directory.

Issue:

```
rm -rf src/github.com/denisbakhtin/ginblog/.git
```

Replace all references of `github.com/denisbakhtin/ginblog` with `github.com/denisbakhtin/amazingblog`:

```
grep -rl 'github.com/denisbakhtin/ginblog' ./ | xargs sed -i 's/github.com\/denisbakhtin\/ginblog/github.com\/denisbakhtin\/amazingblog/g'
```

Move all files to the new location:

```
mv src/github.com/denisbakhtin/ginblog/ src/github.com/denisbakhtin/amazingblog
```

And push it to the corresponding repo:

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

For Continuous Development a good option is to install `fresh` - https://github.com/pilu/fresh
Then simply run `fresh` in the project directory.

To rebuild assets on change install run `npm run watch`.