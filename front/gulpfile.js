'use strict';

// Include Gulp & Tools We'll Use
var gulp = require('gulp');
var fs = require('fs');
var $ = require('gulp-load-plugins')();
var del = require('del');
var runSequence = require('run-sequence');
var path = require('path');
var pkg = require('./package.json');
var through = require('through2');
var swig = require('swig');
var browserSync = require('browser-sync');
var reload = browserSync.reload;

var AUTOPREFIXER_BROWSERS = [
  'ie >= 10',
  'ie_mob >= 10',
  'ff >= 30',
  'chrome >= 34',
  'safari >= 7',
  'opera >= 23',
  'ios >= 7',
  'android >= 4.4',
  'bb >= 10'
];


// Clean Output Directory
gulp.task('clean', del.bind(null, ['dist'], {dot: true}));

// Copy package manger and LICENSE files to dist
gulp.task('metadata', function () {
  return gulp.src(['package.json'])
    .pipe(gulp.dest('./dist'));
});

// Build Production Files, the Default Task
gulp.task('default', ['clean'], function (cb) {
  runSequence(['assets', 'pages'],cb);
});


// ***** Landing page tasks ***** //

/**
 * Site metadata for use with templates.
 * @type {Object}
 */
var site = {};

/**
 * Generates an HTML file based on a template and file metadata.
 */
function applyTemplate() {
  return through.obj(function(file, enc, cb) {
    var data = {
      site: site,
      page: file.page,
      content: file.contents.toString()
    };

    var templateFile = path.join(
        __dirname, 'templates', file.page.layout + '.html');
    var tpl = swig.compileFile(templateFile, {cache: false});
    file.contents = new Buffer(tpl(data), 'utf8');
    this.push(file);
    cb();
  });
}


/**
 * Generates an HTML file for each md file in _pages directory.
 */
gulp.task('pages', function() {
  return gulp.src(['pages/*.md'])
    .pipe($.frontMatter({property: 'page', remove: true}))
    .pipe($.marked())
    .pipe(applyTemplate())
    .pipe($.replace('$$version$$', pkg.version))
    /* Replacing code blocks class name to match Prism's. */
    .pipe($.replace('class="lang-', 'class="language-'))
    /* Translate html code blocks to "markup" because that's what Prism uses. */
    .pipe($.replace('class="language-html', 'class="language-markup'))
    .pipe($.rename(function(path) {
      if (path.basename !== 'index') {
        path.dirname = path.basename;
        path.basename = 'index';
      }
    }))
    .pipe(gulp.dest('dist'));
});

/**
 * Copies assets from MDL and _assets directory.
 */
gulp.task('assets', function () {
  return gulp.src([
      'assets/**/*',
    ])
    .pipe($.if(/\.js/i, $.replace('$$version$$', pkg.version)))
    .pipe($.if(/\.(svg|jpg|png)$/i, $.imagemin({
      progressive: true,
      interlaced: true
    })))
    .pipe($.if(/\.css$/i, $.autoprefixer(AUTOPREFIXER_BROWSERS)))
    .pipe($.if(/\.css$/i, $.csso()))
    .pipe($.if(/\.js$/i, $.uglify({preserveComments: 'some', sourceRoot: '.',
      sourceMapIncludeSources: true})))
    .pipe(gulp.dest('dist/assets'));
});

/**
 * Serves the landing page from "out" directory.
 */
gulp.task('serve:browsersync', ['default'], function () {
  browserSync({
    notify: false,
    server: {
      baseDir: ['dist']
    }
  });

  gulp.watch(['templates/**/*', "pages/**/*"], ['pages', reload]);
  gulp.watch(['assets/**/*'], ['assets', reload]);
});

gulp.task('serve', ['default'], function() {
  $.connect.server({
    root: 'dist',
    port: 5000,
    livereload: true
  });

  gulp.watch(['templates/**/*', "pages/**/*"], ['pages']);
  gulp.watch(['assets/**/*'], ['assets']);

  gulp.src('./dist/index.html')
    .pipe($.open('', {url: 'http://localhost:5000'}));
});
