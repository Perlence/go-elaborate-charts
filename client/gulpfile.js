var gulp = require('gulp')
var coffee = require('gulp-coffee')
var concat = require('gulp-concat')
var gulpif = require('gulp-if')
var less = require('gulp-less')
var minifyCss = require('gulp-minify-css')
var uglify = require('gulp-uglify')

var isRelease = process.env.RELEASE === '1'

// Compile LESS styles.
gulp.task('less', function () {
  var paths = [
    'bower_components/font-awesome/less',
    'bower_components/bootstrap/less',
    'bower_components/ladda-bootstrap/dist'
  ]
  return gulp.src('less/**')
    .pipe(less({paths: paths}))
    .pipe(gulpif(isRelease, minifyCss()))
    .pipe(gulp.dest('css'))
})

// Copy fonts from bower_components to ./fonts.
gulp.task('fonts', function () {
  var fonts = [
    'bower_components/font-awesome/fonts/**'
  ]
  return gulp.src(fonts)
    .pipe(gulp.dest('fonts'))
})

// Concat and uglify third-party packages.
gulp.task('components', function () {
  var components = [
    'bower_components/jquery/dist/jquery.js',
    'bower_components/bootstrap/dist/js/bootstrap.js',
    'bower_components/highstock-release/highstock.src.js',
    'bower_components/ladda-bootstrap/dist/spin.js',
    'bower_components/ladda-bootstrap/dist/ladda.js'
    // 'bower_components/history.js/scripts/bundled-uncompressed/html4+html5/jquery.history.js'
  ]
  return gulp.src(components)
    .pipe(concat('components.js'))
    .pipe(gulpif(isRelease, uglify()))
    .pipe(gulp.dest('js'))
})

// Compile client-side CoffeeScript.
gulp.task('coffee', function () {
  return gulp.src('coffee/**')
    .pipe(coffee({bare: true}))
    .pipe(gulpif(isRelease, uglify()))
    .pipe(gulp.dest('js'))
})

// Watch for changes in the folder with client-side scripts.
gulp.task('watch', function () {
  var watchable = [
    'coffee/**',
    'less/**'
  ]
  return gulp.watch(watchable, ['coffee'])
})

// Build client-side.
gulp.task('build', ['less', 'fonts', 'components'])

// Default task: Build client-side.
gulp.task('default', ['build'])
