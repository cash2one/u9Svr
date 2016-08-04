requirejs.config({
	appDir: 'app',
	baseUrl: "static/scripts",
	dir: '.rjtmp',
	removeCombined: true,
	modules: [
	  {
	    name: 'common',
	    include: [
			'json',
			'mock',
			'jquery',
			'mockjax',
			'bootstrap',
			'paginator',
			'common/base'
	    ]
	  },
	  {
	    name: 'manager/user_statistic/main',
	    include: [
	      'moment',
	      'datepicker',
	      'spin'
	    ],
	    exclude: ['common']
	  },
	],
	paths: {
		json: 'lib/json2',
		mock: 'lib/mock',
		jquery: 'lib/jquery',
		mockjax: 'lib/jquery.mockjax',
		bootstrap:'lib/bootstrap',
		moment: 'lib/moment-with-locales',
		datepicker: 'lib/bootstrap-datetimepicker.min',
		paginator: 'lib/bootstrap-paginator.min',
		spin: 'lib/spin',
		//underscore: 'underscore',
		//backbone: 'backbone',
		//backboneLocalstorage: 'backbone.localStorage',
		//text: 'text',
	},
	shim: {
		'paginator': {
			deps: ['jquery']
		},
		'bootstrap':{
			deps:['jquery']
		}
	}
})






