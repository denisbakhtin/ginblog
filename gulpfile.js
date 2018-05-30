var gulp = require("gulp"),
    sass = require("gulp-sass"),
    autoprefixer = require("gulp-autoprefixer"),
    notify = require("gulp-notify"),
    concat = require("gulp-concat"),
    gzip = require("gulp-gzip"),
    del = require("del")


// Compile SCSS files to CSS
gulp.task("scss", function () {
    //Delete our old css files
    del(["public/css/**/*"])

    //compile hashed css files
    gulp.src("assets/scss/**/*.scss")
        .pipe(sass({
            //outputStyle: "compressed",
            includePaths: [
                "assets/scss",
                "node_modules/bootstrap/scss",
                "node_modules/select2/dist/css",
                "node_modules/@fortawesome/fontawesome-free-webfonts/scss",
            ]
        }).on("error", notify.onError(function (error) {
            return "Error: " + error.message;
        })))
        .pipe(autoprefixer({
            browsers: ["last 20 versions"]
        }))
        .pipe(gulp.dest("public/css"))
})

// images
gulp.task("images", function () {
    del(["public/images/**/*"])
    gulp.src("assets/images/**/*")
        .pipe(gulp.dest("public/images"))
})

// javascript
gulp.task("js", function () {
    //del(["public/js/**/*"])
    gulp.src("node_modules/ckeditor/**/*")
        .pipe(gulp.dest("public/js/ckeditor"))
    gulp.src([
            "node_modules/jquery/dist/jquery.slim.js",
            "node_modules/popper.js/dist/umd/popper.js",
            "node_modules/parsleyjs/dist/parsley.js",
            "node_modules/bootstrap/dist/js/bootstrap.js",
            "node_modules/jquery-slimscroll/jquery.slimscroll.js",
            "node_modules/select2/dist/js/select2.js",
            "assets/js/application.js",
        ])
        .pipe(concat("application.js"))
        .pipe(gulp.dest("public/js"))
})

// fonts
gulp.task('fonts', function () {
    del(["public/fonts/**/*"])
    gulp.src('node_modules/@fortawesome/fontawesome-free-webfonts/webfonts/**.*')
        .pipe(gulp.dest('public/fonts'));
});

// gzip
gulp.task('gzip', function () {
    gulp.src('public/js/*.js')
        .pipe(gzip())
        .pipe(gulp.dest('public/js'));
    gulp.src('public/css/*.css')
        .pipe(gzip())
        .pipe(gulp.dest('public/css'));
});

// Watch asset folder for changes
gulp.task("watch", ["fonts", "scss", "images", "js"], function () {
    gulp.watch("assets/scss/**/*", ["scss"])
    gulp.watch("assets/images/**/*", ["images"])
    gulp.watch("assets/js/**/*", ["js"])
})

// Set watch as default task
gulp.task("default", ["fonts", "scss", "images", "js", "gzip"])