var gulp = require('gulp')
var concat = require('gulp-concat')
var gulpif = require('gulp-if')
var inject = require('gulp-inject')
var less = require('gulp-less')
var minifyCss = require('gulp-minify-css')
var sourcemaps = require('gulp-sourcemaps')
var tsc = require('gulp-typescript')
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
    'bower_components/ladda-bootstrap/dist/ladda.js',

    'bower_components/moment/min/moment-with-locales.js',
    'bower_components/lodash/lodash.js',
    'bower_components/bluebird/js/browser/bluebird.js'
    // 'bower_components/history.js/scripts/bundled-uncompressed/html4+html5/jquery.history.js'
  ]
  return gulp.src(components)
    .pipe(concat('components.js'))
    .pipe(gulpif(isRelease, uglify()))
    .pipe(gulp.dest('js'))
})

gulp.task('declare', function () {
  var target = gulp.src('typings/app.d.ts')
  var sources = gulp.src('./ts/**/*.ts')
  return target
    .pipe(inject(sources, {
      starttag: '//{',
      endtag: '//}',
      transform: function (filepath) {
        return '/// <reference path="..' + filepath + '" />'
      }
    }))
    .pipe(gulp.dest('./typings'))
})

gulp.task('ts', ['declare'], function () {
  var sources = [
    'ts/**/*.ts',
    'typings/**/*.ts'
  ]
  var ts = gulp.src(sources)
    .pipe(sourcemaps.init())
    .pipe(tsc({
      target: 'ES5',
      declarationFiles: true,
      noExternalResolve: true,
      sortOutput: true
    }))
  return ts.js
    .pipe(concat('app.js'))
    .pipe(gulpif(isRelease, uglify()))
    .pipe(sourcemaps.write('js'))
    .pipe(gulp.dest('js'))
})

// Watch for changes in the folder with client-side scripts.
gulp.task('watch', function () {
  gulp.watch('less/**', ['less'])
  gulp.watch('ts/**', ['ts'])
})

// Build client-side.
gulp.task('build', ['less', 'fonts', 'components', 'ts'])

// Default task: Build client-side.
gulp.task('default', ['build'])
