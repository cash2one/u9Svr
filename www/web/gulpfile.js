// generated on 2016-07-18 using generator-webapp 2.1.0
const gulp = require('gulp');
const gulpLoadPlugins = require('gulp-load-plugins');
const browserSync = require('browser-sync');
const del = require('del');
const fileinclude  = require('gulp-file-include');
const $ = gulpLoadPlugins();
const reload = browserSync.reload;
const rjs = require('requirejs');
const prettify = require('gulp-prettify');

var srcDir = 'app';
var outDir = '.tmp';
var rjDir = '.rjtmp'

gulp.task('clean', del.bind(null, ['.tmp', 'dist', rjDir]));

gulp.task('lint', () => {
  return gulp.src(srcDir+'/**/*.js')
    .pipe(reload({stream: true, once: true}))
    .pipe($.eslint({fix: true}))
    .pipe($.eslint.format())
    .pipe($.if(!browserSync.active, $.eslint.failAfterError()))
    .pipe(gulp.dest(srcDir));
});

gulp.task('styles', () => {
  return gulp.src(srcDir+'/**/*.css')
    .pipe($.sourcemaps.init())
    .pipe($.autoprefixer({browsers: ['> 1%', 'last 2 versions', 'Firefox ESR']}))
    .pipe($.sourcemaps.write())
    .pipe(gulp.dest(outDir))
    .pipe(reload({stream: true}));
});

gulp.task('fonts', () => {
  return gulp.src(require('main-bower-files')('**/*.{eot,svg,ttf,woff,woff2}', function (err) {})
    .concat(outDir + '/static/fonts/**/*'))
    .pipe(gulp.dest(outDir + '/static/fonts'));
});

gulp.task('scripts', () => {
  console.log('begin copy script');
  gulp.src([
    'bower_components/json/json2.js',
    'bower_components/mockjs/dist/mock.js',
    'bower_components/jquery/dist/jquery.js',
    'bower_components/jquery-mockjax/dist/jquery.mockjax.js',
    'bower_components/bootstrap/dist/js/bootstrap.js',
    'bower_components/moment/min/moment-with-locales.js',
    'bower_components/eonasdan-bootstrap-datetimepicker/build/js/bootstrap-datetimepicker.min.js',
    'bower_components/spin.js/spin.js',
    'bower_components/bootstrap-paginator/build/bootstrap-paginator.min.js',
  ])
  .pipe(gulp.dest(srcDir + '/static/scripts/lib'));
  console.log('end copy script');
  return;
});

gulp.task('html', ['styles', 'fonts', 'scripts'], () => {
  return gulp.src(srcDir + '/**/*.html')
	  .pipe(fileinclude({prefix: '@@',basepath: '@file'}))
    .pipe($.useref({searchPath: [outDir, srcDir, '.']}))
    .pipe($.if('*.js', $.uglify()))
    .pipe($.if('*.css', $.cssnano({safe: true, autoprefixer: false})))
    .pipe(prettify({indent_inner_html: true, indent_size: 2}))
    //.pipe($.if('*.html', $.htmlmin({collapseWhitespace: true})))
	  .pipe(gulp.dest(outDir + ''));
});

gulp.task('preServe', () => {
  outDir = '.tmp';
});

gulp.task('serve', ['preServe', 'html', 'extras'], () => {
  gulp.src(srcDir + '/static/scripts/**')
  .pipe(gulp.dest(outDir + '/static/scripts'));

  var cleanList = [
    outDir + '/views/**/header.*',
    outDir + '/views/**/footer.*',
    outDir + '/views/common'
  ];
  paths = del.sync(cleanList);
  console.log('Deleted files and folders:\n', paths.join('\n'));

  var staticDir = outDir + '/static';
  browserSync({
    notify: false,
    port: 9000,
    server: {
      baseDir: [outDir],
      routes: {
        staticDir : 'static'
      }
    }
  });

  gulp.watch([srcDir + '/**/*.html']).on('change', reload);
  gulp.watch(srcDir + '/**/*.html', ['html']);
  gulp.watch(srcDir + '/**/*.css', ['styles']);
});

gulp.task('serve:dist', () => {
  outDir = 'dist';
  var staticDir = outDir + '/static';
  browserSync({
    notify: false,
    port: 9000,
    server: {
      baseDir: [outDir],
      routes: {
        staticDir : 'static'
      }
    }
  });
});

gulp.task('preBuild', () => {
  outDir = 'dist'
});

gulp.task('requirejs', () => {
  console.log('begin optimize requirejs');
  rjs.optimize({
    mainConfigFile:'require_config.js',
  }, function(buildResponse){
    console.log('requirejs:', buildResponse);
    console.log('end optimize requirejs');
  });
  return;
});

gulp.task('extras', () => {
  return gulp.src([srcDir + '/static/*.*'])
  .pipe(gulp.dest(outDir + '/static'));
});

gulp.task('build', ['preBuild', 'lint', 'html', 'requirejs', 'extras'], () => {
  scriptDir = '/static/scripts'
  rjScriptDir = rjDir + scriptDir;

  console.log('begin copy requrejs script');
  gulp.src(rjScriptDir + ['/lib/require.js'])
  .pipe(gulp.dest(outDir + scriptDir + '/lib'));

  gulp.src([rjScriptDir + '/**/*.js','!' + rjScriptDir + '/lib/*.js',])
  .pipe(gulp.dest(outDir + scriptDir));
  console.log('end copy requrejs script');


  console.log('begin optimize script');
  scriptOutDir = outDir + scriptDir;
  gulp.src(scriptOutDir + '/**/*.js')
    .pipe($.plumber())
    .pipe($.sourcemaps.init())
    //.pipe($.babel())
    .pipe($.sourcemaps.write('.'))
    .pipe(gulp.dest(scriptOutDir));
  console.log('end optimize script');


  console.log('begin optimize views');
  var htmlCleanList = [
    outDir + '/views/**/header.*',
    outDir + '/views/**/footer.*',
    //outDir + '/views/**/index.*',
    outDir + '/views/common',
  ];
  paths = del.sync(htmlCleanList);
  console.log('Deleted files and folders:\n', paths.join('\n'));
  console.log('end optimize views');


  return gulp.src(outDir + '/**/*').pipe($.size({title: 'build', gzip: true}));
});

gulp.task('default', ['clean'], () => {
  gulp.start('build');
});
